package commands_qa

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QaCheckListCommand struct{}

func (c *QaCheckListCommand) Name() string {
	return "check_list"
}

func (c *QaCheckListCommand) Usage() string {
	return "Generate check list of test cases"
}

func (c *QaCheckListCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *QaCheckListCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	initialString := qa_helper.QA_GenerateListOfTestCases

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("task description:").
		AddUrls(urls).
		AddFiles(files).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
