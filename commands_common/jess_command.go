package commands_common

import (
	"fmt"
	"os"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

type BaseCommand interface {
	Flags() []cli.Flag
	PreparePrompt(gpt *gpt.GptClient, cliContext *cli.Context) (string, error)
	Name() string
	Usage() string
}

type JessCommand struct {
	Command BaseCommand
}

func (c *JessCommand) DefineCommand(gpt *gpt.GptClient) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.handleAction(gpt),
		Flags:  c.Command.Flags(),
	}
}

func (c *JessCommand) handleAction(gpt *gpt.GptClient) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		context := cliContext.String("context")
		output := cliContext.String("output")
		finalPrompt, err := c.Command.PreparePrompt(gpt, cliContext)
		if err != nil {
			return err
		}
		quit := make(chan bool)
		go jess_cli.AnimateThinking(quit)
		answer, err := gpt.SendMessageWithContextDepth(finalPrompt, context, 0, false)
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
