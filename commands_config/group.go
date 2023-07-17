package commands_config

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineTestCommand(llmClient *client.Client, config *utils.AppConfig) *cli.Command {
	testCommand := commands_common.JessCommand{
		Command: &TestCommand{},
	}

	return &cli.Command{
		Name:  "test",
		Usage: "Actions to take with config and check system availability",
		Subcommands: []*cli.Command{
			testCommand.DefineCommand(llmClient),
		},
	}
}
