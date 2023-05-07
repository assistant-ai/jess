package context_commands

import (
	"github.com/assistant-ai/jess/rest"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineServeCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "serve",
		Usage:  "start REST service",
		Action: handleServeAction(gpt),
		Flags:  serveFlags(),
	}
}

func handleServeAction(gpt *gpt.GptClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		rest.StartRest(gpt)
		return nil
	}
}

func serveFlags() []cli.Flag {
	return []cli.Flag{}
}
