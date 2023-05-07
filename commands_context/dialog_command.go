package commands_context

import (
	"errors"
	"fmt"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineDialogCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "dialog",
		Usage:  "Manage dialogs",
		Action: handleDialogAction(gpt),
		Flags:  dialogFlags(),
	}
}

func handleDialogAction(gpt *gpt.GptClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Bool("list") {
			HandleDialogList()
		} else if id := c.String("context"); id != "" {
			handleDialogContinue(id, gpt)
		} else if id := c.String("show"); id != "" {
			handleDialogShow(id)
		} else if id := c.String("delete"); id != "" {
			HandleDialogDelete(id)
		} else {
			return errors.New("no required key provided")
		}
		return nil
	}
}

func dialogFlags() []cli.Flag {
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
		&cli.StringFlag{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "show dialog with the context id",
		},
		&cli.StringFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete dialog with the context id",
		},
	}
}

func HandleDialogList() {
	contextIds, err := db.GetContextIDs()
	jess_cli.HandleError(err)

	jess_cli.PrintContextIDs(contextIds)
}

func handleDialogContinue(id string, gpt *gpt.GptClient) {
	fmt.Println("Starting a new conversation...")

	err := jess_cli.StartChat(id, gpt)
	jess_cli.HandleError(err)
}

func handleDialogShow(id string) {
	messages, err := db.GetMessagesByContextID(id)
	jess_cli.HandleError(err)

	jess_cli.ShowMessages(messages)
}

func HandleDialogDelete(id string) {
	err := db.RemoveContext(id)
	jess_cli.HandleError(err)
}
