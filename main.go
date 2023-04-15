package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/assistent-ai/client/chat"
	jess_cli "github.com/assistent-ai/client/cli"
	"github.com/assistent-ai/client/db"
	"github.com/assistent-ai/client/gpt"
	"github.com/assistent-ai/client/model"
	"github.com/urfave/cli/v2"
)

const Version = "2"

func main() {
	apiKeyFilePath := ""
	defaultFilePath := filepath.Join(os.Getenv("HOME"), ".open-ai.key")
	app := cli.NewApp()
	app.Name = "jessica"
	app.Usage = "Jessica is an AI assistent."
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "key-file",
			Value:       defaultFilePath,
			Usage:       "Path to the text file containing the API key",
			Destination: &apiKeyFilePath,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "dialog",
			Usage: "Manage dialogs",
			Action: func(c *cli.Context) error {
				if c.Bool("list") {
					dialogIds, err := db.GetDialogIDs()
					if err != nil {
						cli.Exit(err, 1)
					} else {
						jess_cli.PrintDialogIDs(dialogIds)
					}
					return nil
				} else if c.String("continue") != "" {
					id := c.String("continue")
					if id == "" {
						cli.Exit(errors.New("please provide dialog id"), 1)
					} else {
						// Replace with your actual logic to start a new conversation
						fmt.Println("Starting a new conversation...")
						ctx, err := initContext(apiKeyFilePath)
						if err != nil {
							cli.Exit(err, 1)
						}
						err = chat.StartChat(id, ctx)
						if err != nil {
							cli.Exit(err, 1)
						}
					}
					return nil
				} else if c.String("show") != "" {
					id := c.String("show")
					if id == "" {
						cli.Exit(errors.New("please provide dialog id"), 1)
					} else {
						messages, err := db.GetMessagesByDialogID(id)
						if err != nil {
							cli.Exit(err, 1)
						} else {
							chat.ShowMessages(messages)
						}
					}
					return nil
				} else if c.String("delete") != "" {
					id := c.String("delete")
					if id == "" {
						cli.Exit(errors.New("please provide dialog id"), 1)
					} else {
						err := db.RemoveMessagesByDialogId(id)
						if err != nil {
							cli.Exit(err, 1)
						}
					}
					return nil
				}
				return errors.New("no required key provided")
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all dialogs",
				},
				&cli.StringFlag{
					Name:    "continue",
					Aliases: []string{"c"},
					Usage:   "continue dialog with the id",
				},
				&cli.StringFlag{
					Name:    "show",
					Aliases: []string{"s"},
					Usage:   "show dialog with the id",
				},
				&cli.StringFlag{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete dialog with the id",
				},
			},
		},
		{
			Name:  "file",
			Usage: "Process file input",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Aliases:  []string{"p"},
					Usage:    "path to the file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "input",
					Aliases:  []string{"i"},
					Usage:    "prompt input to pass with the file",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				filePath := c.String("path")
				prompt := c.String("input")
				ctx, _ := initContext(apiKeyFilePath)
				b, _ := os.ReadFile(filePath)
				input := model.FileInput{
					UserMessage: prompt,
					FileContent: string(b),
				}
				gptPrompt, _ := chat.GeneratePromptForFile(input)
				answer, _ := gpt.RandomMessage(gptPrompt, ctx)
				fmt.Println(answer)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		cli.Exit(err, 1)
	}
}

func initContext(openAiKeyFilePath string) (*model.AppContext, error) {
	b, err := os.ReadFile(openAiKeyFilePath) // just pass the file name
	if err != nil {
		return nil, err
	}
	return &model.AppContext{
		OpenAiKey: strings.ReplaceAll(string(b), "\n", ""),
	}, nil
}
