package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineTextCommand(llmClient *client.Client) *cli.Command {
	mailCommand := commands_common.JessCommand{
		Command: &MailCommand{},
	}
	grammarCommand := commands_common.JessCommand{
		Command: &GrammarCommand{},
	}
	return &cli.Command{
		Name:  "text",
		Usage: "Actions to take with text",
		Subcommands: []*cli.Command{
			mailCommand.DefineCommand(llmClient),
			grammarCommand.DefineCommand(llmClient),
		},
	}
}
