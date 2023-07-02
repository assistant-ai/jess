package commands_config

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
	"os"
)

func DefineTestCommand(llmClient *client.Client, config *utils.AppConfig) *cli.Command {
	testCommand := commands_common.JessCommand{
		Command: &TestCommand{},
	}
	//configCommand := ConfigCommand{
	//	Command: &ConfigCommand{},
	//}

	//configCommand := &ConfigCommand{}

	modelNameMsg := "Model name:" + config.ModelName
	logInfoMsg := "Log level:" + config.LogLevel
	openaiKeyPath := "OpenAI API key path:" + config.OpenAiApiKeyPath
	yourKey, _ := os.ReadFile(config.OpenAiApiKeyPath)
	maskedKey := string(yourKey[0:5]) + "..." + string(yourKey[len(yourKey)-5:len(yourKey)])
	msgMaskedKey := "OpenAI API key stored in file: " + maskedKey
	utils.Println_green("Your current configuration:")
	utils.Println_yellow(logInfoMsg)
	utils.Println_yellow(modelNameMsg)
	utils.Println_yellow(openaiKeyPath)
	utils.Println_yellow(string(msgMaskedKey))

	return &cli.Command{
		Name:  "test",
		Usage: "Actions to take with config and check system availability",
		Subcommands: []*cli.Command{
			testCommand.DefineCommand(llmClient),
			//configCommand.DefineCommand(config),
		},
	}
}
