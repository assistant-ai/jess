package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/db"
	"github.com/assistant-ai/jess/gpt"
	"github.com/assistant-ai/jess/model"
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
	ctx, err := initContext(*apiKeyFilePath)
	if err != nil {
		return nil, err
	}
	return []*cli.Command{
		defineDialogCommand(apiKeyFilePath),
		jess_cli.DefineProcessCommand(ctx),
		defineCodeCommand(ctx),
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
	dialogIds, err := db.GetDialogIDs()
	if err != nil {
		cli.Exit(err, 1)
	} else {
		jess_cli.PrintDialogIDs(dialogIds)
	}
}

func handleDialogContinue(id string, apiKeyFilePath *string) {
	fmt.Println("Starting a new conversation...")
	ctx, err := initContext(*apiKeyFilePath)
	if err != nil {
		cli.Exit(err, 1)
	}
	err = jess_cli.StartChat(id, ctx)
	if err != nil {
		cli.Exit(err, 1)
	}
}

func handleDialogShow(id string) {
	messages, err := db.GetMessagesByDialogID(id)
	if err != nil {
		cli.Exit(err, 1)
	} else {
		jess_cli.ShowMessages(messages)
	}
}

func handleDialogDelete(id string) {
	err := db.RemoveMessagesByDialogId(id)
	if err != nil {
		cli.Exit(err, 1)
	}
}

func defineCodeCommand(ctx *model.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "code",
		Usage: "Actions to take with code",
		Subcommands: []*cli.Command{
			jess_cli.DefineQuestionCommand(ctx),
			jess_cli.DefineExplainCommand(ctx),
		},
	}
}

func handleCodeApplyAction(apiKeyFilePath *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePath := c.String("input")
		ctx, err := initContext(*apiKeyFilePath)
		if err != nil {
			return err
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		var fileList []string

		err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Ignore hidden files
			if strings.Contains(path, ".git") {
				return nil
			}

			// Append the relative path
			relPath, _ := filepath.Rel(cwd, path)
			fileList = append(fileList, relPath)
			return nil
		})

		if err != nil {
			panic(err)
		}

		// Print the list of files
		// fmt.Printf("Here is the list of files:\n%s\n", strings.Join(fileList, "\n"))

		gptPrompt := fmt.Sprintf("Here is the content of the input file:\n---\n%s\n---\nHere is the list of files I see around:\n---\n%s\n---\nWhat would be the next step to take to implement this? In your answer first line should be one of the following: show, file or unknow. Show means that you need to see specific file(s) in order to make actions. In this case next line(s) should be list of files you need to see to do the action. file means that you can apply action to files and next line should be file content that you created with the applied changes. Unknown means that you can not express next steps with show or file commands. Next like in this case should show why it is not possible to do it.", string(fileContent), strings.Join(fileList, "\n"))
		// fmt.Println(gptPrompt)
		// // for {
		answer, err := gpt.RandomMessage(gptPrompt, ctx)
		if err != nil {
			return err
		}
		fmt.Println(answer)
		// // Parse GPT's response and take appropriate actions

		// // }
		return nil
	}
}

func defineApplyCommand(apiKeyFilePath *string) *cli.Command {
	return &cli.Command{
		Name:   "apply",
		Usage:  "apply chaing of thoughts",
		Flags:  applyFlags(),
		Action: handleFileAction(apiKeyFilePath),
	}
}

func applyFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "path to input file",
			Required: true,
		},
	}
}

func codeFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "path to input file",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "prompt",
			Aliases: []string{"p"},
			Usage:   "prompt input to pass with the file",
		},
		&cli.BoolFlag{
			Name:    "refactor",
			Aliases: []string{"r"},
			Usage:   "refactor file by applying best practices",
		},
		&cli.BoolFlag{
			Name:    "explain",
			Aliases: []string{"e"},
			Usage:   "explain to me what is going on here",
		},
		&cli.BoolFlag{
			Name:    "override",
			Aliases: []string{"o"},
			Usage:   "write output in the same input file instead of stdout",
		},
	}
}

func handleFileAction(apiKeyFilePath *string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePath := c.String("input")
		prompt := c.String("prompt")
		refactor := c.Bool("refactor")
		explain := c.Bool("explain")
		override := c.Bool("override")
		if refactor {
			prompt = "Refactor following file, extract code, de-duplicate, apply all best practices that you can think off that would be valuable here and would improve readability"
		} else if explain {
			prompt = "Please explain to me in simple words what this code do, on the high level what you think it is doing and per public method/class/function, whatvere you can to help me to understand it better"
		}
		ctx, err := initContext(*apiKeyFilePath)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		input := model.FileInput{
			UserMessage: prompt,
			FileContent: string(b),
		}
		gptPrompt, err := jess_cli.GeneratePromptForFile(input)
		if err != nil {
			return err
		}
		quit := make(chan bool)
		go jess_cli.AnimateThinking(quit)
		answer, err := gpt.RandomMessage(gptPrompt, ctx)
		if err != nil {
			return err
		}
		quit <- true
		if override {
			err = os.WriteFile(filePath, []byte(answer), 0644)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("\n\n" + answer)
		}
		return nil
	}
}

func initContext(openAiKeyFilePath string) (*model.AppContext, error) {
	b, err := os.ReadFile(openAiKeyFilePath)
	if err != nil {
		return nil, err
	}
	return &model.AppContext{
		OpenAiKey: strings.ReplaceAll(string(b), "\n", ""),
	}, nil
}
