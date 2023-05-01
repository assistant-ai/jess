package context_commands

import (
	"errors"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineContextCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "context",
		Usage:  "Manage contexts",
		Action: handleContextAction(gpt),
		Flags:  contextFlags(),
	}
}

func handleContextAction(gpt *gpt.GptClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Bool("list") {
			HandleDialogList()
		} else {
			id := c.String("context")
			if c.Bool("show") {
				handleContextShow(id)
			} else if c.Bool("delete") {
				HandleDialogDelete(id)
			} else if prompt := c.String("prompt"); prompt != "" {
				handleContextSet(id, prompt)
			} else {
				return errors.New("no required key provided")
			}
		}
		return nil
	}
}

func contextFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all dialogs",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "continue dialog with the contextn id",
		},
		&cli.BoolFlag{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "show dialog with the context id",
		},
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete dialog with the context id",
		},
		&cli.StringFlag{
			Name:    "prompt",
			Aliases: []string{"p"},
			Usage:   "set dialog context prompt to this value",
		},
	}
}

func handleContextShow(contextId string) {
	message, err := db.GetContextMessage(contextId)
	messages := []db.Message{message}
	jess_cli.HandleError(err)

	jess_cli.ShowMessages(messages)
}

func handleContextSet(contextId string, prompt string) {
	err := db.UpdateContext(contextId, prompt)
	jess_cli.HandleError(err)
}
