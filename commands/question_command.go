package commands

import (
	"github.com/assistant-ai/llmchat-client/gpt"
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
		PromptFlag(),
		InputFilesFlag(),
		ContextFlag(),
		OutputFlag(),
	}
}

func (c *QuestionCommand) PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	prePrompt := "Let me show you code files and than I will show you my question for the code in this files, please answer with example of the code where possible."
	finalPrompt, err := FilePromptBuilder(prePrompt, filePaths, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
