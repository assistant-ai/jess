package commands_qa

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QaListToCasesCommand struct{}

func (c *QaListToCasesCommand) Name() string {
	return "list2cases"
}

func (c *QaListToCasesCommand) Usage() string {
	return "Generate test cases from list"
}

func (c *QaListToCasesCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *QaListToCasesCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	initialString := qa_helper.QA_ListToCases

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("check list of tests:").
		AddUrls(urls).
		AddFiles(files).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
