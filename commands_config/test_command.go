package commands_config

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/urfave/cli/v2"
)

type TestCommand struct{}

func (c *TestCommand) Name() string {
	return "test"
}

func (c *TestCommand) Usage() string {
	return "Check if everything is configured fine and you have access to all required resources"
}

// TODO rebuild this command after changing promt builder
func (c *TestCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.ContextFlag(),
	}
}

// TODO rebuild this command after changing promt builder
func (c *TestCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	prePrompt := "User is just want to be sure if chat gpt is available. You should just reply with model of chat gpt that you are using."
	finalPrompt := prePrompt
	return finalPrompt, nil
}
