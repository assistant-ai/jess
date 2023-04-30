package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/commands"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	apiKeyFilePath := filepath.Join(os.Getenv("HOME"), ".open-ai.key")
	app, err := setupApp(&apiKeyFilePath)
	if err != nil {
		cli.Exit(err, 1)
		panic(err)
	}

	err = app.Run(os.Args)
	if err != nil {
		cli.Exit(err, 1)
		panic(err)
	}
}

func setupApp(apiKeyFilePath *string) (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "jessica"
	app.Usage = "Jessica is an AI assistent."
	app.Version = version

	commands, err := defineCommands(apiKeyFilePath)
	if err != nil {
		return nil, err
	}
	app.Commands = commands

	return app, nil
}

func defineCommands(apiKeyFilePath *string) ([]*cli.Command, error) {
	gpt, err := initGptClient(*apiKeyFilePath)
	if err != nil {
		return nil, err
	}
	processCommand := commands.JessCommand{
		Command: &commands.ProcessCommand{},
	}
	return []*cli.Command{
		defineDialogCommand(apiKeyFilePath),
		processCommand.DefineCommand(gpt),
		defineCodeCommand(gpt),
	}, nil
}

func defineDialogCommand(apiKeyFilePath *string) *cli.Command {
	return &cli.Command{
		Name:   "dialog",
		Usage:  "Manage dialogs",
		Action: handleDialogAction(apiKeyFilePath),
		Flags:  dialogFlags(),
	}
}

func handleDialogAction(apiKeyFilePath *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c.Bool("list") {
			handleDialogList()
		} else if id := c.String("context"); id != "" {
			handleDialogContinue(id, apiKeyFilePath)
		} else if id := c.String("show"); id != "" {
			handleDialogShow(id)
		} else if id := c.String("delete"); id != "" {
			handleDialogDelete(id)
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

func handleDialogList() {
	contextIds, err := db.GetContextIDs()
	if err != nil {
		cli.Exit(err, 1)
	} else {
		jess_cli.PrintContextIDs(contextIds)
	}
}

func handleDialogContinue(id string, apiKeyFilePath *string) {
	fmt.Println("Starting a new conversation...")
	ctx, err := initGptClient(*apiKeyFilePath)
	if err != nil {
		cli.Exit(err, 1)
	}
	err = jess_cli.StartChat(id, ctx)
	if err != nil {
		cli.Exit(err, 1)
	}
}

func handleDialogShow(id string) {
	messages, err := db.GetMessagesByContextID(id)
	if err != nil {
		cli.Exit(err, 1)
	} else {
		jess_cli.ShowMessages(messages)
	}
}

func handleDialogDelete(id string) {
	err := db.RemoveContext(id)
	if err != nil {
		cli.Exit(err, 1)
	}
}

func defineCodeCommand(gpt *gpt.GptClient) *cli.Command {
	questionCommand := commands.JessCommand{
		Command: &commands.QuestionCommand{},
	}
	explainCommand := commands.JessCommand{
		Command: &commands.ExplainCommand{},
	}
	refactorCommand := commands.JessCommand{
		Command: &commands.RefactorCommand{},
	}
	return &cli.Command{
		Name:  "code",
		Usage: "Actions to take with code",
		Subcommands: []*cli.Command{
			questionCommand.DefineCommand(gpt),
			explainCommand.DefineCommand(gpt),
			refactorCommand.DefineCommand(gpt),
		},
	}
}

func initGptClient(openAiKeyFilePath string) (*gpt.GptClient, error) {
	b, err := os.ReadFile(openAiKeyFilePath)
	if err != nil {
		return nil, err
	}
	client := gpt.NewDefaultGptClient(strings.ReplaceAll(string(b), "\n", ""))
	client.DefaultContext = `Your name is Jessica, but everyone call you Jess. You are AI assitent for software developers to help them with their code: explain/refactor/answer questions. Mostly you used as CLI tool, but not only.
	
When replying, consider information gaps and ask for clarification if needed. 
Limit this to avoid excess. 
Decide when to answer directly. 
Assume basic knowledge. 
Concise over politeness.`
	return client, nil
}
