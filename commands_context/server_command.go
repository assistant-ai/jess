package commands_context

import (
	"github.com/assistant-ai/jess/rest"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineServeCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   "serve",
		Usage:  "start REST service",
		Action: handleServeAction(llmClient),
		Flags:  serveFlags(),
	}
}

func handleServeAction(llmClient *client.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		rest.StartRest(llmClient)
		return nil
	}
}

func serveFlags() []cli.Flag {
	return []cli.Flag{}
}
