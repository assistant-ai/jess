package commands_context

import (
	"errors"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/urfave/cli/v2"
)

func DefineContextCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   "context",
		Usage:  "Manage contexts",
		Action: handleContextAction(llmClient),
		Flags:  contextFlags(),
	}
}

func handleContextAction(llmClient *client.Client) func(c *cli.Context) error {
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
			Usage:   "List all dialogs",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "Continue dialog with the context id",
		},
		&cli.BoolFlag{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show dialog with the context id",
		},
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete dialog with the context id",
		},
		&cli.StringFlag{
			Name:    "prompt",
			Aliases: []string{"p"},
			Usage:   "Set dialog context prompt to this value",
		},
	}
}

func handleContextShow(contextId string) {
	message, err := db.GetContextMessage(contextId)
	jess_cli.HandleError(err)
	jess_cli.ShowContext(message)
}

func handleContextSet(contextId string, prompt string) {
	err := db.UpdateContext(contextId, prompt)
	jess_cli.HandleError(err)
}
