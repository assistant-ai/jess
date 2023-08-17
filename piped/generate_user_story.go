package piped

import (
	"encoding/json"
	"fmt"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	"strings"
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
		//context := cliContext.String("context")
		err := error(nil)
		initialUserTopic := cliContext.String("prompt")
		// TODO add validation if file is empty or not, here could be issue that related if it is emprty or not, if it is empty then there should nt be ech for tilda
		//inputFilePath := cliContext.String("input_file")
		outputFolder, err := utils.ExpandTilde(cliContext.String("output_folder"))
		if err != nil {
			return err
		}

		_, err = utils.IsValidPath(outputFolder)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		//initialPrompt, err := c.Command.PreparePromptForDoubleAction(cliContext)

		if err != nil {
			log.Errorf("Error while sending message: %v", err)
			return err
		}

		utils.PrintlnCyan("USER PROMPT:\n" + initialUserTopic + "\n\n")
		// Section that responsible for generating user story file
		//detailedUserStoryAnswer, err := jess_cli.ExecutePrompt(llmClient, initialPrompt, context)
		//filePathForPromptOutput := outputFolder + "/00_user_story.txt"
		//err = utils.AnswersOutput(filePathForPromptOutput, detailedUserStoryAnswer)

		// Section that responsible for generating sub tasks files
		utils.PrintlnCyan("\nCollecting list of sub tasks for generated user story...\n")
		path_to_file_withUser_story := "/Users/nik/GitHub/jess/temp_experiments/us/00_user_story.txt"

		finalPrompt, err := prompt.FilePromptBuilder("give me list of the subtasks from this file. return only list and nothing more. return it as json. it should be in next format: {\n\t  \"subtasks\": []}", []string{path_to_file_withUser_story}, []string{}, []string{}, "")
		resultOfSubTasksCollections, err := jess_cli.ExecutePrompt(llmClient, finalPrompt, "")
		if err != nil {
			log.Errorf("Error while sending message: %v", err)
			return err
		}

		var JsonMap map[string][]string
		err = json.Unmarshal([]byte(resultOfSubTasksCollections), &JsonMap)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		subtasks := JsonMap["subtasks"]

		sizeOfSubTasks := len(subtasks[0:3])

		// generate list of test cases for provided user story
		utils.PrintlnCyan("\nGenerating basic test cases for: " + strings.ToUpper(initialUserTopic) + "\n")
		fileNameForBasicTestCases := outputFolder + "/" + fmt.Sprintf("%02d", sizeOfSubTasks+1) + "_BasicTestCases.txt"
		// TODO need to implement of using file with detailes if it was provided to the command. Right now it is hard to handle issues with empty files
		promptForbugHunting, err := prompt.FilePromptBuilder(text.BUG_HUNTER_PROMPT, []string{}, []string{}, []string{}, initialUserTopic)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		temporaryResultBugHunting, err := jess_cli.ExecutePrompt(llmClient, promptForbugHunting, "")
		err = utils.AnswersOutput(fileNameForBasicTestCases, temporaryResultBugHunting)

		// Iterate over subtasks and print each one
		utils.PrintlnCyan("\nGenerating subtasks for: " + strings.ToUpper(initialUserTopic) + "\n")
		for idx, subtask := range subtasks[0:3] {
			fileNameForSubTask := outputFolder + "/" + fmt.Sprintf("%02d", idx+1) + "_" + utils.ReplaceSpacesWithUnderscores(subtask) + ".txt"
			fmt.Println("\n", idx+1, " / ", sizeOfSubTasks, " - ", subtask)
			// TODO need to implement of using file with detailes if it was provided to the command. Right now it is hard to handle issues with empty files
			promptForSubTasks, err := prompt.FilePromptBuilder(text.TECH_TASK_PROMPT, []string{}, []string{}, []string{}, subtask)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			temporaryResultForSubTask, err := jess_cli.ExecutePrompt(llmClient, promptForSubTasks, "")
			err = utils.AnswersOutput(fileNameForSubTask, temporaryResultForSubTask)
		}

		return nil
	}

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

func (c *GenerateDetailedUserStoryCommand) PreparePromptForDoubleAction(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input_file")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := text.USER_STORY_PROMPT
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}