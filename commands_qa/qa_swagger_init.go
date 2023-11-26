package commands_qa

import (
	"fmt"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	_ "github.com/assistant-ai/prompt-tools"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QaSwaggerInitCommand struct{}

func (c *QaSwaggerInitCommand) Name() string {
	return "swagg_init"
}

func (c *QaSwaggerInitCommand) Usage() string {
	return "Get high level recommendations for swagger API"
}

func (c *QaSwaggerInitCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
		commands_common.PromptFlag(),
	}
}

func (c *QaSwaggerInitCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	userInput := cliContext.StringSlice("prompt")
	typeOfTestMsg := fmt.Sprintf("path for creating tests: %s", userInput)
	initialString := qa_helper.QA_swagger_check_list

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Swagger json:").
		AddFiles(files).
		AddUrls(urls).
		AddTextToPrompt(typeOfTestMsg).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
