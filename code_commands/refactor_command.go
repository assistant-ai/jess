package code_commands

import (
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

type RefactorCommand struct{}

func (c *RefactorCommand) Name() string {
	return "refactor"
}

func (c *RefactorCommand) Usage() string {
	return "Refacotr file for me"
}

func (c *RefactorCommand) Flags() []cli.Flag {
	return []cli.Flag{
		InputFileFlag(),
		ContextFlag(),
		OutputFlag(),
	}
}

func (c *RefactorCommand) PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error) {
	filePath := cliContext.String("input")
	filePaths := []string{filePath}
	urls := []string{}
	gDriveFiles := []string{}
	prePrompt := "Let me show you code file."
	userPrompt := "Please refactor it using the best practices of the language that is used. Ideally final verison should be as readable as possible. You can refactor any way you want as long as public APIs are no changed. Feel free to extract anything that is needed to a helper function/methods entities. Your final output should be full content of the file end-to-end without any text before/after, I will use your output and override my original file."
	finalPrompt, err := FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
