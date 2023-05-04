package code_commands

import (
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

type ExplainCommand struct{}

func (c *ExplainCommand) Name() string {
	return "explain"
}

func (c *ExplainCommand) Usage() string {
	return "Explain code for me"
}

func (c *ExplainCommand) Flags() []cli.Flag {
	return []cli.Flag{
		InputFilesFlag(),
		ContextFlag(),
		OutputFlag(),
		UrlsFlag(),
		GoogleDriveFilesFlag(),
	}
}

func (c *ExplainCommand) PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	urlsPaths := cliContext.StringSlice("url")
	googleDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "Let me show you code files."
	userPrompt := "Please explain this code for me in plain English."
	finalPrompt, err := FilePromptBuilder(prePrompt, filePaths, urlsPaths, googleDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
