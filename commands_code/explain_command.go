package commands_code

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
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
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.UrlsFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *ExplainCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	urlsPaths := cliContext.StringSlice("url")
	googleDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "Let me show you code files."
	userPrompt := "Please explain this code for me in plain text English. For each section of explanation add in front also lines of code that you are explaining. Add short summary at start of explanation. At the start add file name"
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urlsPaths, googleDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
