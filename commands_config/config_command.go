package commands_config

import (
	"bufio"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

const configFilePath = "/Users/nik/.jess/config.yaml"

func DefineConfigCommand(config *utils.AppConfig) *cli.Command {
	return &cli.Command{
		Name:   "config",
		Usage:  "Check if everything is configured fine and you have access to all required resources",
		Flags:  ConfigFlags(),
		Action: ConfigAction(),
	}
}

// TODO rebuild this command after changing promt builder
func ConfigFlags() []cli.Flag {
	return []cli.Flag{
		commands_common.ContextFlag(),
	}
}

// TODO rebuild this command after changing promt builder
func ConfigAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		listOfModels := gpt.GetListOfModels()
		msgForSetupModels := "[IMPORTANT] Use only next models: \n" + strings.Join(listOfModels, "\n")
		setupValue("model", msgForSetupModels)
		setupValue("openai_api_key_path", "")
		setupValue("log_level", "")
		setupValue("gcp.service_account_key_path", "")
		utils.PrintlnGreen("Configuration changed successfully")
		os.Exit(0)
		return nil
	}
}

func setupValue(configKey string, msg string) {
	msgForInput := msg + "\nPlease print new " + strings.ToUpper(configKey) + " you want to use:\n [ for skip press enter ]"
	utils.PrintlnCyan(msgForInput)
	utils.PrintCyanInvite()
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		println("error while reading input")
	}
	newValue := scanner.Text()
	if newValue != "" {
		viper.SetConfigFile(configFilePath)
		viper.ReadInConfig()
		viper.Set(configKey, newValue)
		viper.WriteConfig()
	}
}
