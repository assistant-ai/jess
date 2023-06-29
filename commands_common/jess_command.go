package commands_common

import (
	"github.com/assistant-ai/jess/utils"
	"os"
	"strings"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

type BaseCommand interface {
	Flags() []cli.Flag
	PreparePrompt(cliContext *cli.Context) (string, error)
	Name() string
	Usage() string
}

type JessCommand struct {
	Command BaseCommand
}

func (c *JessCommand) DefineCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.handleAction(llmClient),
		Flags:  c.Command.Flags(),
	}
}

func (c *JessCommand) handleAction(llmClient *client.Client) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		context := cliContext.String("context")
		output := cliContext.String("output")
		finalPrompt, err := c.Command.PreparePrompt(cliContext)
		if err != nil {
			return err
		}
		quit := make(chan bool)
		go jess_cli.AnimateThinking(quit)
		answer, err := llmClient.SendMessageWithContextDepth(finalPrompt, context, 0, false)
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
			//fmt.Println("\n\n" + answer)
			//TODO wrap changes with a function

			answer = strings.Replace(answer, "File/Url List:", "\033[33mFile/Url List:\033[32m", -1)
			answer = strings.Replace(answer, "File/Url path:", "\033[33mFile/Url path:\033[32m", -1)
			answer = strings.Replace(answer, "User Prompt:", "\033[33mUser Prompt:\033[32m", -1)
			answer = strings.Replace(answer, "Content:", "\033[33mUpdated content:\033[32m", -1)
			utils.Println_green("\n\n" + answer)
		}
		return nil
	}
}
