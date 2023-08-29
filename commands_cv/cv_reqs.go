package commands_cv

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/cv_helper"
	_ "github.com/assistant-ai/prompt-tools"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type CvReqirementsCommand struct{}

func (c *CvReqirementsCommand) Name() string {
	return "reqs"
}

func (c *CvReqirementsCommand) Usage() string {
	return "Get cv requirements for user from provided URL or file"
}

func (c *CvReqirementsCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *CvReqirementsCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	initialString := cv_helper.CV_ReqirementsCollectorPrompt

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Position:").
		AddUrls(urls).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
