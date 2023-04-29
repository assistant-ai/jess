package cli

import (
	"fmt"
	"os"

	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineProcessCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   "process",
		Usage:  "Do the process actions",
		Action: HandleProcessAction(gpt),
		Flags:  ProcessFlags(),
	}
}

func ProcessFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "prompt",
			Aliases: []string{"p"},
			Usage:   "prompt to suppy with file",
		},
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
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output file path, if not specificed stdout will be used",
		},
	}
}

func HandleProcessAction(gpt *gpt.GptClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		filePaths := c.StringSlice("input")
		userPrompt := c.String("prompt")
		context := c.String("context")
		output := c.String("output")
		finalPrompt, err := prompt.GenerateMultiFileProcessPrompt(filePaths, userPrompt)
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
		if output != "" {
			err = os.WriteFile(output, []byte(answer), 0644)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("\n\n" + answer)
		}
		return nil
	}
}