package commands_config

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/utils"
	"github.com/urfave/cli/v2"
)

type TestCommand struct{}

func (c *TestCommand) Name() string {
	return "test"
}

func (c *TestCommand) Usage() string {
	return "Check if everything is configured fine and you have access to all required resources"
}

// TODO rebuild this command after changing promt builder
func (c *TestCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.ContextFlag(),
	}
}

// TODO rebuild this command after changing promt builder
func (c *TestCommand) PreparePrompt(cliContext *cli.Context) (string, error) {

	config := utils.GetConfig()
	msgConfigFolder := "Config folder:" + utils.GetDefaultConfigFolderPath()
	msgConfigPath := "Config folder:" + utils.GetDefaultConfigFilePath()
	modelNameMsg := "Model name:" + config.ModelName
	logInfoMsg := "Log level:" + config.LogLevel
	openaiKeyPath := "OpenAI API key path:" + config.OpenAiApiKeyPath
	maskedKey := utils.GetMaskedApiKey(config.OpenAiApiKeyPath)
	msgMaskedKey := "OpenAI API key stored in file: " + maskedKey
	// TODO update showing this message using utils.PrintConfig - it should show all config parameters
	utils.PrintlnGreen("Your current configuration:")
	utils.PrintlnYellow(msgConfigFolder)
	utils.PrintlnYellow(msgConfigPath)
	utils.PrintlnYellow(logInfoMsg)
	utils.PrintlnYellow(modelNameMsg)
	utils.PrintlnYellow(openaiKeyPath)
	utils.PrintlnYellow(string(msgMaskedKey))

	prePrompt := "User is just want to be sure if chat gpt is available. You should just reply with model of chat gpt that you are using."
	finalPrompt := prePrompt
	return finalPrompt, nil
}
