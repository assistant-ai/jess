package utils

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type AppConfig struct {
	OpenAiApiKeyPath      string
	ServiceAccountKeyPath string
	ModelName             string
	GCPProjectId          string
	LogLevel              string
	JessEmailAccount      string
}

var singletonJessConfig *AppConfig

var configPath string
var folderPath string

func init() {
	folderPath = GetDefaultConfigFolderPath()
	configPath = GetDefaultConfigFilePath()

	if !IfConfigFileExists(configPath) {
		file, err := os.Create(configPath)
		if err != nil {
			log.Fatalf("can't create config file: %v", err)
			panic(err)
		}
		defer file.Close()
		viper.SetConfigFile(configPath)

		configStructureMap := GetConfigMap()
		listOfDefaultConfigParameters := GetListOfConfigFields()

		for _, element := range listOfDefaultConfigParameters {
			nameInConfigFile := configStructureMap[element].nameInConfigFile
			defaultValue := configStructureMap[element].defaultValue
			viper.SetDefault(nameInConfigFile, defaultValue)
		}

		viper.WriteConfig()
		PrintlnRed("There was no config file by default path: " + configPath)
		PrintlnRed("Config was file created successfully at " + configPath)
	}
}

// TODO split this function into two: one for loading the config and one for creating the config file
func LoadConfig(configPath string) (*AppConfig, error) {

	// TODO I would suggest to tup Always default folder
	if configPath == "" {
		configPath = GetDefaultConfigFilePath()
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	singletonJessConfig = &AppConfig{
		OpenAiApiKeyPath:      viper.GetString("openai.openai_api_key_path"),
		ServiceAccountKeyPath: viper.GetString("gcp.service_account_key_path"),
		JessEmailAccount:      viper.GetString("gcp.jessica_mail_account"),
		ModelName:             viper.GetString("model"),
		GCPProjectId:          viper.GetString("gcp.gcp_project_id"),
		LogLevel:              viper.GetString("log_level"),
	}

	//TODO extract this check somewhere else
	CheckExistingOfOpenAPIKey(configPath)
	return singletonJessConfig, nil
}

func CheckExistingOfOpenAPIKey(configPath string) {
	if !IfFileWithAPiKeyExists(singletonJessConfig.OpenAiApiKeyPath) {
		PrintlnRed("OpenAI API key file not found at " + singletonJessConfig.OpenAiApiKeyPath)
		PrintlnRed("According to the config file at " + configPath + " your API key file should be at " + singletonJessConfig.OpenAiApiKeyPath)
		PrintlnRed("Please create a file with your OpenAI API key at. Or change the path in the config file")
		PrintlnRed("You can get your API key at https://beta.openai.com/account/api-keys")
		PrintlnRed("Then run `echo YOUR_OPEN_AI_API_TOKEN > " + singletonJessConfig.OpenAiApiKeyPath + "`")
		os.Exit(1)
	}
}

func GetConfig() *AppConfig {
	return singletonJessConfig
}
