package piped

import (
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefinePipedCommand(llmClient *client.Client) *cli.Command {

	doublePromptCommand := DoubleJessCommand{
		Command: &DoublePromptCommand{},
	}
	generateCommitMsg := GenerateCommitJessCommand{
		Command: &GenerateCommitCommand{},
	}
	return &cli.Command{
		Name:  "pipe",
		Usage: "Actions that would be executed in multiple calls by piping the output of one prompt as the input to another.",
		Subcommands: []*cli.Command{
			doublePromptCommand.DefineCommand(llmClient),
			generateCommitMsg.DefineCommand(llmClient),
		},
	}
}
