package cli

import (
	"fmt"

	"github.com/assistant-ai/jess/gpt"
	"github.com/assistant-ai/jess/model"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

func DefineExplainCommand(ctx *model.AppContext) *cli.Command {
	return &cli.Command{
		Name:   "explain",
		Usage:  "Explain the code",
		Action: handleExplainAction(ctx),
		Flags:  explainFlags(),
	}
}

func explainFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "input files",
		},
		&cli.StringFlag{
			Name:    "context",
			Aliases: []string{"c"},
			Usage:   "context id to store this to",
		},
	}
}

func handleExplainAction(ctx *model.AppContext) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePaths := c.StringSlice("input")
		context := c.String("context")
		finalPrompt, err := prompt.GenerateMultiFileProcessPrompt(filePaths, "Explain on the high level in basic terms to me what is going on in these files. I should be able to understand overall what the project is about and understand components that are used in it so if I need to change something or implement a feature I will know where to look what to do.")
		if err != nil {
			return err
		}
		quit := make(chan bool)
		go AnimateThinking(quit)
		answer, err := gpt.SendStringMessage(finalPrompt, context, ctx)
		quit <- true
		if err != nil {
			return err
		}
		fmt.Println("\n\n" + answer)
		return nil
	}
}
