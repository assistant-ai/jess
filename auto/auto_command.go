package auto

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"

	"github.com/sirupsen/logrus"
)

func DefineAutoCommand(llmClient *client.Client, logger *logrus.Logger) *cli.Command {
	return &cli.Command{
		Name:   "auto",
		Usage:  "do requested commands end-to-end",
		Action: handleAutoAction(llmClient, logger),
		Flags:  autoFlags(),
	}
}

func handleAutoAction(llmClient *client.Client, logger *logrus.Logger) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		id := c.String("context")
		prompt := c.String("prompt")
		input := c.String("input")
		return StartProcess(prompt, input, id, llmClient, logger)
	}
}

func autoFlags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFileFlag(),
		commands_common.ContextFlag(),
		commands_common.PromptFlag(),
	}
}
