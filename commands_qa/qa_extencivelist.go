package commands_qa

import (
	"fmt"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QaExtensiveListCommand struct{}

func (c *QaExtensiveListCommand) Name() string {
	return "e_check_list"
}

func (c *QaExtensiveListCommand) Usage() string {
	return "Generate check list of test cases"
}

func (c *QaExtensiveListCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
		commands_common.PromptFlag(),
	}
}

func (c *QaExtensiveListCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	userInput := cliContext.StringSlice("prompt")
	typeOfTestMsg := fmt.Sprintf("\"Type of test:\"%s", userInput)
	initialString := qa_helper.QA_ExtensiveListOfTestCases

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt(typeOfTestMsg).
		AddTextToPrompt("task description:").
		AddUrls(urls).
		AddFiles(files).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
