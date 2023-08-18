package piped

import (
	"encoding/json"
	"fmt"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GenerateDetailedUserJessCommand struct {
	Command DoubleBaseCommand
}

func (c *GenerateDetailedUserJessCommand) DefineCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.handleActionForCommit(llmClient),
		Flags:  c.Command.Flags(),
	}
}

func (c *GenerateDetailedUserJessCommand) handleActionForCommit(llmClient *client.Client) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		err := error(nil)
		initialUserTopic := cliContext.String("prompt")
		outputFolder, err := utils.ExpandTilde(cliContext.String("output_folder"))
		if err != nil {
			return err
		}
		_, err = utils.IsValidPath(outputFolder)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		filePathForPromptOutput := outputFolder + "/00_user_story.txt"

		// Section for generation of user story
		utils.PrintlnCyan("USER PROMPT:\n" + initialUserTopic + "\n\n")
		err = generateUserStory(cliContext, c, llmClient, filePathForPromptOutput)
		if err != nil {
			return err
		}

		// Section that responsible for generating sub-tasks files
		utils.PrintlnCyan("\nCollecting list of sub tasks for generated user story...\n")
		subtasks, err := getListOfSubTasks(filePathForPromptOutput, llmClient)
		if err != nil {
			return err
		}

		sizeOfSubTasks := len(subtasks)
		//generate list of test cases for provided user story
		utils.PrintlnCyan("\nGenerating basic TEST CASES for: " + strings.ToUpper(initialUserTopic) + "\n")
		err = generateTestCases(outputFolder, sizeOfSubTasks, llmClient, cliContext, c)
		if err != nil {
			fmt.Println("Error:", err)
		}

		//Iterate over subtasks and print each one
		utils.PrintlnCyan("\nGenerating SUB_TASKS for: " + strings.ToUpper(initialUserTopic) + "\n")
		err = generateTechTasks(subtasks, outputFolder, llmClient, cliContext, c)
		if err != nil {
			return err
		}

		return nil
	}

}

func generateTechTasks(subtasks []string, outputFolder string, llmClient *client.Client, cliContext *cli.Context, c *GenerateDetailedUserJessCommand) error {
	sizeOfSubTasks := len(subtasks)
	var wg sync.WaitGroup
	for idx, subtask := range subtasks {
		wg.Add(1)
		go func(idx int, subtask string, llmClient *client.Client, cliContext *cli.Context, c *GenerateDetailedUserJessCommand) {
			fileNameForSubTask := outputFolder + "/" + fmt.Sprintf("%02d", idx+1) + "_" + utils.ReplaceSpacesWithUnderscores(subtask) + ".txt"
			fmt.Println("\n", strconv.Itoa(idx+1), " / ", sizeOfSubTasks, " - ", subtask)
			PROMPT := text.TECH_TASK_PROMPT + "\nTech task: " + subtask
			promptForSubTasks, _ := c.Command.PreparePromptForDoubleAction(cliContext, PROMPT)
			uuidObj, _ := uuid.NewUUID()
			time.Sleep(1 * time.Second)
			temporaryResultForSubTask, _ := jess_cli.ExecutePrompt(llmClient, promptForSubTasks, uuidObj.String())
			_ = utils.AnswersOutput(fileNameForSubTask, temporaryResultForSubTask)
			wg.Done()
		}(idx, subtask, llmClient, cliContext, c)
	}
	wg.Wait()
	return nil
}

func generateTestCases(outputFolder string, sizeOfSubTasks int, llmClient *client.Client, cliContext *cli.Context, c *GenerateDetailedUserJessCommand) error {
	fileNameForBasicTestCases := outputFolder + "/" + fmt.Sprintf("%02d", sizeOfSubTasks+1) + "_BasicTestCases.txt"
	bugHunterPrompt, err := c.Command.PreparePromptForDoubleAction(cliContext, text.BUG_HUNTER_PROMPT)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	temporaryResultBugHunting, err := jess_cli.ExecutePrompt(llmClient, bugHunterPrompt, uuidObj.String())
	err = utils.AnswersOutput(fileNameForBasicTestCases, temporaryResultBugHunting)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func generateUserStory(cliContext *cli.Context, c *GenerateDetailedUserJessCommand, llmClient *client.Client, filePathForPromptOutput string) error {
	initialPrompt, err := c.Command.PreparePromptForDoubleAction(cliContext, text.USER_STORY_PROMPT)
	if err != nil {
		log.Errorf("Error while sending message: %v", err)
		return err
	}
	// generate random context for user story in case if it is empty it just return  prompt
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	detailedUserStoryAnswer, err := jess_cli.ExecutePrompt(llmClient, initialPrompt, uuidObj.String())
	err = utils.AnswersOutput(filePathForPromptOutput, detailedUserStoryAnswer)
	if err != nil {
		log.Errorf("Error while sending message: %v", err)
		return err
	}
	return nil
}

func getListOfSubTasks(filePathForPromptOutput string, llmClient *client.Client) ([]string, error) {
	getSubTasksPrompt, err := prompt.FilePromptBuilder("give me list of the subtasks from this file. return only list and nothing more. return it as json. it should be in next format: {\n\t  \"subtasks\": []}", []string{filePathForPromptOutput}, []string{}, []string{}, "")
	resultOfSubTasksCollections, err := jess_cli.ExecutePrompt(llmClient, getSubTasksPrompt, "")
	if err != nil {
		log.Errorf("Error while sending message: %v", err)
		return nil, err
	}
	var JsonMap map[string][]string
	err = json.Unmarshal([]byte(resultOfSubTasksCollections), &JsonMap)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	subtasks := JsonMap["subtasks"]
	return subtasks, nil
}

type GenerateDetailedUserStoryCommand struct{}

func (c *GenerateDetailedUserStoryCommand) Name() string {
	return "big_user_story"
}

func (c *GenerateDetailedUserStoryCommand) Usage() string {
	return "Generate detailed user story, with extra subtasks and test cases"
}

func (c *GenerateDetailedUserStoryCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[Mandatory] USER STORY - Add information about user story you want to generate. Just give summary in one sentence. you want to provide as much detais as possible",
			Value:    "",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "input_file",
			Aliases:  []string{"i"},
			Usage:    "[Optional] Input file that could contain additional details for users story",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "output_folder",
			Aliases:  []string{"o"},
			Usage:    "[Mandatory] folder to store generated files",
			Required: true,
		},
	}
}

func (c *GenerateDetailedUserStoryCommand) PreparePromptForDoubleAction(cliContext *cli.Context, storedPrompt string) (string, error) {
	filePaths := cliContext.StringSlice("input_file")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := storedPrompt
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
