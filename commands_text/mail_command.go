package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

type MailCommand struct{}

func (c *MailCommand) Name() string {
	return "mail"
}

func (c *MailCommand) Usage() string {
	return "re-write my email for me"
}

func (c *MailCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFileFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.PromptFlag(),
	}
}

func (c *MailCommand) PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error) {
	filePath := cliContext.String("input")
	userPrompt := cliContext.String("prompt")
	filePaths := []string{filePath}
	urls := []string{}
	gDriveFiles := []string{}
	prePrompt := "User going to provide you with mail text as well as some context to it. Figure out which mail user wants to update yourself. You should re-write mail without chaning the intent or meaning but make it as clear as possible, as concrete as possible. User might provide additional requirenments."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
