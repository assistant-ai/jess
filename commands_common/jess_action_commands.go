package commands_common

import (
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

type ExecActionCommand interface {
	Flags() []cli.Flag
	Name() string
	Usage() string
	ExecAction(cliContext *client.Client, llmClient *cli.Context) cli.ActionFunc
}

func (c *JessActionCommand) DefineCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.PipedAction(llmClient),
		Flags:  c.Command.Flags(),
	}
}

func (c *JessActionCommand) PipedAction(llmClient *client.Client) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		c.Command.ExecAction(llmClient, cliContext)
		return nil
	}
}

type JessActionCommand struct {
	Command ExecActionCommand
}
