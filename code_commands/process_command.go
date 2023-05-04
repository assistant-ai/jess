package code_commands

import (
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

type ProcessCommand struct{}

func (c *ProcessCommand) Name() string {
	return "process"
}

func (c *ProcessCommand) Usage() string {
	return "Process a files"
}

func (c *ProcessCommand) Flags() []cli.Flag {
	return []cli.Flag{
		PromptFlag(),
		InputFilesFlag(),
		ContextFlag(),
		OutputFlag(),
		UrlsFlag(),
		GoogleDriveFilesFlag(),
	}
}

func (c *ProcessCommand) PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urlsList := cliContext.StringSlice("url")
	googleDriveFiles := cliContext.StringSlice("gdrive")

	prePrompt := "Let me show you files and/or and than I will show you my prompt to use, it might include questions/asks about the files."
	finalPrompt, err := FilePromptBuilder(prePrompt, filePaths, urlsList, googleDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
