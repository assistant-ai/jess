package commands_cv

import (
	"encoding/json"
	"fmt"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/prompt_storage/cv_helper"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	_ "github.com/assistant-ai/prompt-tools"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"strings"
)

type JobPosition struct {
	PositionName string   `json:"position"`
	Requirements []string `json:"requirements"`
}

func getPromptForGenerateListOfRequirements(positionUrl []string) (string, error) {

	initialString := cv_helper.CV_ReqirementsCollectorPrompt

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Position:").
		AddUrls(positionUrl).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil

}

func getPromptForGenRecommendations(cvFilePath string, listOfReqs string) (string, error) {

	initialString := cv_helper.CV_reccomendationPrompt

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Users CV:").
		AddFile(cvFilePath).
		AddTextToPrompt("List of requirements to the position:").
		AddTextToPrompt(listOfReqs).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil

}

type CvRecommendationCommand struct{}

func (c *CvRecommendationCommand) Name() string {
	return "cv_help"
}

func (c *CvRecommendationCommand) Usage() string {
	return "Get cv requirements from provided URL and generate recommendations for user how to improve his cv. Main ide that CV should be in txt format. at first it would generate file with requirements and also generate file with recommendation how to improve CV."
}

func (c *CvRecommendationCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input_file",
			Aliases:  []string{"i"},
			Usage:    "[Mandatory] Add cv file path for further analysis. Only plain text files is supported.",
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output_folder",
			Aliases:  []string{"o"},
			Usage:    "[Mandatory] Folder to store generated files.",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Usage:    "[Mandatory] Url to position description. Url should be open and be visible without ant authorization.",
			Required: true,
		},
	}
}

func (c *CvRecommendationCommand) ExecAction(llmClient *client.Client, cliContext *cli.Context) cli.ActionFunc {
	//initialize section
	cvFile, err := utils.ExpandTilde(cliContext.String("input_file"))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	utils.IsValidPath(cvFile)
	urls := cliContext.StringSlice("url")
	providedOutputFolder, err := utils.ExpandTilde(cliContext.String("output_folder"))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	utils.IsValidPath(providedOutputFolder)
	utils.CreateFolderIfNotExists(providedOutputFolder)
	listOfReqsPrompt, err := getPromptForGenerateListOfRequirements(urls)

	if err != nil {
		panic(err)
		return nil
	}
	uuidObj, _ := uuid.NewUUID()
	randomContextId := uuidObj.String()
	var jobPosition JobPosition

	// Collecting requirements section
	println("\nCollecting requirements from URL...\n")
	collectedListOfRequirements, err := jess_cli.ExecutePrompt(llmClient, listOfReqsPrompt, randomContextId)
	err = json.Unmarshal([]byte(collectedListOfRequirements), &jobPosition)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	outputFolderWithPosition := providedOutputFolder + "/" + utils.ReplaceSpacesWithUnderscores(jobPosition.PositionName)
	utils.CreateFolderIfNotExists(outputFolderWithPosition)
	pathOutputToReqs := outputFolderWithPosition + "/" + "00_reqs_output.txt"
	err = utils.AnswersOutput(pathOutputToReqs, collectedListOfRequirements)
	if err != nil {
		panic(err)
		return nil
	}

	listOfRequirements := strings.Join(jobPosition.Requirements, ", ")
	recommendationsPrompt, err := getPromptForGenRecommendations(cvFile, listOfRequirements)
	if err != nil {
		panic(err)
		return nil
	}
	println("\nGenerating recommendation based on provided CV...\n")
	recommendation, err := jess_cli.ExecutePrompt(llmClient, recommendationsPrompt, randomContextId)
	pathToRecommendationOutputFile := outputFolderWithPosition + "/" + "01_recommendation.txt"
	err = utils.AnswersOutput(pathToRecommendationOutputFile, recommendation)
	if err != nil {
		panic(err)
		return nil
	}

	return nil
}
