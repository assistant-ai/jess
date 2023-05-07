package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineTextCommand(gpt *gpt.GptClient) *cli.Command {
	mailCommand := commands_common.JessCommand{
		Command: &MailCommand{},
	}
	return &cli.Command{
		Name:  "text",
		Usage: "Actions to take with text",
		Subcommands: []*cli.Command{
			mailCommand.DefineCommand(gpt),
		},
	}
}
