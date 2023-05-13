package commands_context

import (
	"fmt"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"

	"github.com/sirupsen/logrus"
)

func DefineDialogCommand(llmClient *client.Client, logger *logrus.Logger) *cli.Command {
	return &cli.Command{
		Name:   "dialog",
		Usage:  "Manage dialogs",
		Action: handleDialogAction(llmClient, logger),
		Flags:  dialogFlags(),
	}
}

func handleDialogAction(llmClient *client.Client, logger *logrus.Logger) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Bool("list") {
			HandleDialogList()
		} else if id := c.String("context"); id != "" {
			handleDialogContinue(id, llmClient, logger)
		} else if id := c.String("show"); id != "" {
			handleDialogShow(id)
		} else if id := c.String("delete"); id != "" {
			HandleDialogDelete(id)
		} else {
			handleDialogContinue("", llmClient, logger)
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

func handleDialogContinue(id string, llmClient *client.Client, logger *logrus.Logger) {
	fmt.Println("Starting a new conversation...")

	err := jess_cli.StartChat(id, llmClient, logger)
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
