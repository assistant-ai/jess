package cli

import (
	"fmt"

	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineExplainCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "explain",
		Usage:  "Explain the code",
		Action: handleExplainAction(gpt),
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

func handleExplainAction(gpt *gpt.GptClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePaths := c.StringSlice("input")
		context := c.String("context")
		finalPrompt, err := prompt.GenerateMultiFileProcessPrompt(filePaths, "Explain on the high level in basic terms to me what is going on in these files. I should be able to understand overall what the project is about and understand components that are used in it so if I need to change something or implement a feature I will know where to look what to do.")
		if err != nil {
			return err
		}
		quit := make(chan bool)
		go AnimateThinking(quit)
		answer, err := gpt.SendMessage(finalPrompt, context)
		quit <- true
		if err != nil {
			return err
		}
		fmt.Println("\n\n" + answer)
		return nil
	}
}
