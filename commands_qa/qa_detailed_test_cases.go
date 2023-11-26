package commands_qa

import (
	"encoding/json"
	"fmt"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	"github.com/assistant-ai/llmchat-client/client"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type TestCase struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

type TestCases struct {
	TestCases []TestCase `json:"testCases"`
}

func PromtForConvertCheckListToJSON(inputFileWithTextDescription []string) (string, error) {
	initialPrompt := qa_helper.QA_convertCheckListToJson
	//description := fmt.Sprintf("test cases check list: %s", inputFileWithTextDescription)

	finalPrompt, err := prompttools.CreateInitialPrompt(initialPrompt).
		AddTextToPrompt("test cases check list:").
		AddFiles(inputFileWithTextDescription).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}

type QADetailedTestVasesCommand struct{}

func (c *QADetailedTestVasesCommand) Name() string {
	return "chl2tc"
}

func (c *QADetailedTestVasesCommand) Usage() string {
	return "generate detailed test cases from list"
}

func (c *QADetailedTestVasesCommand) Flags() []cli.Flag {
	return []cli.Flag{
		//&cli.StringFlag{
		//	Name:     "input_file",
		//	Aliases:  []string{"i"},
		//	Usage:    "[Mandatory] Add cv file path for further analysis. Only plain text files is supported.",
		//	Value:    "",
		//	Required: true,
		//},
		commands_common.InputFilesFlag(),
		&cli.StringFlag{
			Name:     "output_folder",
			Aliases:  []string{"o"},
			Usage:    "[Mandatory] Folder to store generated files.",
			Required: false,
		},
	}
}

func (c *QADetailedTestVasesCommand) ExecAction(llmClient *client.Client, cliContext *cli.Context) cli.ActionFunc {
	//initialize section
	//inputFile, err := utils.ExpandTilde(cliContext.String("input_file"))
	inputFile := cliContext.StringSlice("input")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return nil
	//}
	promptConverter, err := PromtForConvertCheckListToJSON(inputFile)
	if err != nil {
		panic(err)
		return nil
	}
	uuidObj, _ := uuid.NewUUID()
	randomContextId := uuidObj.String()

	//var checklistDescriptipin JobPosition
	// Collecting requirements section
	println("\nCollecting requirements from URL...\n")
	collectedListOfRequirements, err := jess_cli.ExecutePrompt(llmClient, promptConverter, randomContextId)
	if err != nil {
		panic(err)
		return nil
	}
	//fmt.Println(collectedListOfRequirements)

	var testCasesData TestCases
	if err := json.Unmarshal([]byte(collectedListOfRequirements), &testCasesData); err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	var listOfPrompts []string

	// Iterate over the test cases and print the titles and details
	for _, testCase := range testCasesData.TestCases {
		//fmt.Println("Title:", testCase.Title)
		//fmt.Println("Details:", testCase.Details)
		//fmt.Println()
		prompt := fmt.Sprintf("%s,title:%s, details:%s", qa_helper.QA_convertTitleToTestCase, testCase.Title, testCase.Details)
		listOfPrompts = append(listOfPrompts, prompt)
	}
	result, err := jess_cli.ExecutePrompt(llmClient, listOfPrompts[0], randomContextId)
	fmt.Println(result)

	//err := json.Unmarshal(collectedListOfRequirements, &jsonData)
	//if err != nil {
	//	fmt.Println("Error decoding JSON:", err)
	//	return
	//}

	// Print the decoded JSON data
	//fmt.Printf("%+v\n", jsonData)

	return nil
}
