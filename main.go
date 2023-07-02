package main

import (
	"fmt"
	"github.com/assistant-ai/jess/commands_config"
	"os"
	"strings"

	"github.com/assistant-ai/jess/auto"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/commands_code"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/commands_context"
	"github.com/assistant-ai/jess/commands_text"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/assistant-ai/llmchat-client/palm"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	app, err := setupApp()
	jess_cli.HandleError(err)

	err = app.Run(os.Args)
	jess_cli.HandleError(err)
}

func setupApp() (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "jessica"
	app.Usage = "Jessica is an AI assistent."
	app.Version = version

	config, err := utils.LoadConfig("")
	if err != nil {
		return nil, err
	}
	logger := utils.SetupLogger(config)
	logger.Debug("Logger been initialized")
	commands, err := defineCommands(config, logger)
	if err != nil {
		return nil, err
	}
	app.Commands = commands

	return app, nil
}

func defineCommands(config *utils.AppConfig, logger *logrus.Logger) ([]*cli.Command, error) {
	llmClient, err := initClient(config, logger)
	if err != nil {
		return nil, err
	}

	processCommand := commands_common.JessCommand{
		Command: &commands_code.ProcessCommand{},
	}

	commands := []*cli.Command{
		commands_context.DefineDialogCommand(llmClient, logger),
		commands_context.DefineContextCommand(llmClient),
		processCommand.DefineCommand(llmClient),
		commands_code.DefineCodeCommand(llmClient),
		commands_text.DefineTextCommand(llmClient),
		commands_config.DefineTestCommand(llmClient, config),
		commands_config.DefineConfig1Command(config),
		commands_context.DefineServeCommand(llmClient),
		auto.DefineAutoCommand(llmClient, logger),
	}

	return commands, nil
}

func initClient(config *utils.AppConfig, logger *logrus.Logger) (*client.Client, error) {
	var llmClient *client.Client
	var err error
	modelName := strings.ToLower(config.ModelName)
	fmt.Printf("Model that is used for this task: %s\n", modelName)
	logger.WithFields(logrus.Fields{
		"config.ModelName": config.ModelName,
	}).Debug("Creating client")
	if modelName == "gpt4" {
		llmClient, err = gpt.NewDefaultGptClientFromFile(config.OpenAiApiKeyPath, nil)
		if err != nil {
			return nil, err
		}
	} else if modelName == "gpt3" {
		llmClient, _ = gpt.NewGptClientFromFile(config.OpenAiApiKeyPath, 3, gpt.ModelGPT3Turbo, db.RandomContextId, 4000, nil)
	} else if modelName == "palm" {
		if config.GCPProjectId == "" {
			errorText := "model is PaLM but GCP Project ID is null"
			logger.Error(errorText)
			return nil, fmt.Errorf(errorText)
		}
		if config.ServiceAccountKeyPath == "" {
			errorText := "model is PaLM but GCP service acout json path is null"
			logger.Error(errorText)
			return nil, fmt.Errorf(errorText)
		}
		llmClient, err = palm.NewPalmClient(config.GCPProjectId, config.ServiceAccountKeyPath)
		if err != nil {
			return nil, err
		}
	} else {
		logger.WithFields(logrus.Fields{
			"config.ModelName": config.ModelName,
		}).Error("Model is not specified")
		utils.Println_red("Try to use next command to fix model error")
		utils.Println_yellow("jess config -c 'id'")
		return nil, nil
	}
	llmClient.DefaultContext = `Your name is Jessica, but everyone call you Jess. You are AI assitent for software developers to help them with their code: explain/refactor/answer questions. Mostly you used as CLI tool, but not only.

When replying, consider information gaps and ask for clarification if needed.
Limit this to avoid excess.
Decide when to answer directly.
Assume basic knowledge.
Concise over politeness.`
	return llmClient, nil
}
