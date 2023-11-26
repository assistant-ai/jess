package commands_qa

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt_storage/qa_helper"
	_ "github.com/assistant-ai/prompt-tools"
	prompttools "github.com/assistant-ai/prompt-tools"
	"github.com/urfave/cli/v2"
)

type QaRecommendationCommand struct{}

func (c *QaRecommendationCommand) Name() string {
	return "recommendations"
}

func (c *QaRecommendationCommand) Usage() string {
	return "Get high level recommendations"
}

func (c *QaRecommendationCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.UrlsFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *QaRecommendationCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	urls := cliContext.StringSlice("url")
	files := cliContext.StringSlice("input")
	initialString := qa_helper.QA_BasicReccomendation

	finalPrompt, err := prompttools.CreateInitialPrompt(initialString).
		AddTextToPrompt("Task description:").
		AddFiles(files).
		AddUrls(urls).
		GenerateFinalPrompt()

	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
