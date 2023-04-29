package cli

import (
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineQuestionCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "question",
		Usage:  "Ask questions about the code",
		Action: HandleProcessAction(gpt),
		Flags:  ProcessFlags(),
	}
}
