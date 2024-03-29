package commands_code

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

type QuestionCommand struct{}

func (c *QuestionCommand) Name() string {
	return "question"
}

func (c *QuestionCommand) Usage() string {
	return "Ask question about code files"
}

func (c *QuestionCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.PromptFlag(),
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.UrlsFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *QuestionCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urlsList := cliContext.StringSlice("url")
	googleDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "Let me show you code files and than I will show you my question for the code in this files. Identify the main language of these files, and imagine you are an expert in that programming language. And provide answer with example of the code and best practices for that code where possible."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urlsList, googleDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
