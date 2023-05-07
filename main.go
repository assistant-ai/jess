package main

import (
	"os"
	"path/filepath"
	"strings"

	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/commands_code"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/commands_context"
	"github.com/assistant-ai/jess/commands_text"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	apiKeyFilePath := filepath.Join(os.Getenv("HOME"), ".open-ai.key")
	app, err := setupApp(&apiKeyFilePath)
	jess_cli.HandleError(err)

	err = app.Run(os.Args)
	jess_cli.HandleError(err)
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

	processCommand := commands_common.JessCommand{
		Command: &commands_code.ProcessCommand{},
	}

	commands := []*cli.Command{
		commands_context.DefineDialogCommand(gpt),
		commands_context.DefineContextCommand(gpt),
		processCommand.DefineCommand(gpt),
		commands_code.DefineCodeCommand(gpt),
		commands_text.DefineTextCommand(gpt),
		commands_context.DefineServeCommand(gpt),
	}

	return commands, nil
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
