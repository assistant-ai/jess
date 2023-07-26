package commands_config

import (
	"bufio"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/utils"
	"github.com/sirupsen/logrus"
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
		configStructureMap := utils.GetConfigMap()
		listOfDefaultConfigParameters := utils.GetListOfConfigFields()
		utils.PrintlnRed("W A R N I N G\n\n CHANGING CONFIGURATION:\n ++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		for _, element := range listOfDefaultConfigParameters {
			SetConfigElementWithNewValue(configStructureMap[element])
		}

		utils.PrintlnGreen("Configuration changed successfully")
		os.Exit(0)
		return nil
	}
}

func SetConfigElementWithNewValue(config utils.Config) {
	nameInConfigFile := config.GetNameInConfigFile()
	suggestedValues := config.GetSuggestedValues()
	recommendedMSG := generateRecommendationValuesMSG(suggestedValues)
	overviewMSG := generateOverviewMSG(config)
	nameInAppConfig := config.GetName()
	setupValueWithKeyInputInvitation(nameInConfigFile, recommendedMSG, overviewMSG, nameInAppConfig)
}

func generateOverviewMSG(config utils.Config) string {
	nameInConfigFile := config.GetNameInConfigFile()
	appConfigParameter := config.GetName()
	descriptionStr := config.GetShortDescription()
	descriptionMSG := "DESCRIPTION: \n\n parameter  >>>  " + nameInConfigFile + "   <<<  from config file is used to set  >>>   " + appConfigParameter + "   <<<  in application configuration.\n " + descriptionStr
	return descriptionMSG
}

func generateRecommendationValuesMSG(suggestedValues []string) string {
	listOfSuggestedValues := ""
	var recommendedMSG string
	if len(suggestedValues) == 1 {
		listOfSuggestedValues = suggestedValues[0]
		if listOfSuggestedValues == "" {
			recommendedMSG = ""
		}
	} else {
		listOfSuggestedValues = strings.Join(
			suggestedValues[:],
			", ",
		)
		recommendedMSG = "USE ONLY: " + listOfSuggestedValues
	}
	return recommendedMSG
}

func setupValueWithKeyInputInvitation(configKey string, recommendationMsg string, descriptionMsg string, appConfigName string) {
	supportiveMsg := recommendationMsg + "\n " + descriptionMsg + "\n\n [ for skip press Enter ]"
	inviteSetupMsg := "\n Please type new " + strings.ToUpper(configKey) + " you want to use: "
	utils.PrintlnYellow(inviteSetupMsg)
	utils.PrintlnCyan(supportiveMsg)
	utils.PrintCyanInvite()
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		println("error while reading input")
	}
	newValue := scanner.Text()
	if newValue != "" {

		println("Config key is: " + configKey)
		println("App config name is: " + appConfigName)
		oldConfig, err := utils.LoadConfig(configFilePath)
		if err != nil {
			println("error while loading old config file during changing value")
			logrus.Panic(err)
		}
		utils.PrintFieldValue(oldConfig, appConfigName, "OLD")
		viper.SetConfigFile(configFilePath)
		viper.ReadInConfig()
		viper.Set(configKey, newValue)
		viper.WriteConfig()
		newConfig, err := utils.LoadConfig(configFilePath)
		if err != nil {
			println("error while loading updated config file during changing value")
			logrus.Panic(err)
		}
		utils.PrintFieldValue(newConfig, appConfigName, "NEW")
		msgPath := "value was changed in config file: " + configFilePath
		utils.PrintlnYellow(msgPath)
	}
}
