package commands_qa

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	_ "github.com/assistant-ai/prompt-tools"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QASwaggerCurlCommand struct{}

func (c *QASwaggerCurlCommand) Name() string {
	return "swagg_curl"
}

func (c *QASwaggerCurlCommand) Usage() string {
	return "Create cmd CURl code for testing API based on swagger test cases"
}

func (c *QASwaggerCurlCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *QASwaggerCurlCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	initialString := qa_helper.QA_SwagerCurl

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Swagger json and test cases:").
		AddFiles(files).
		AddUrls(urls).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
