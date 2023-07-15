package utils

import (
	"github.com/b0noi/go-utils/v2/fs"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type AppConfig struct {
	OpenAiApiKeyPath      string
	ServiceAccountKeyPath string
	ModelName             string
	GCPProjectId          string
	LogLevel              string
}

func LoadConfig(configPath string) (*AppConfig, error) {
	folderPath, err := fs.MaybeCreateProgramFolder("jess")
	if err != nil {
		return nil, err
	}
	if configPath == "" {
		configPath = filepath.Join(folderPath, "config.yaml")
	}

	openAiApiKeyFilePath := filepath.Join(folderPath, "open-ai.key")
	// TODO figure out how to get the service account key file from the user
	// TODO add instructions on how to store and add it
	jessServiceAccountFilePath := filepath.Join(folderPath, "jess-service-account.json")

	if !IfConfigFileExists(configPath) {
		file, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		viper.SetConfigFile(configPath)
		viper.SetDefault("openai.openai_api_key_path", openAiApiKeyFilePath)
		viper.SetDefault("gcp.service_account_key_path", jessServiceAccountFilePath)
		viper.SetDefault("model", "gpt3Turbo")
		viper.SetDefault("gcp.gcp_project_id", "")
		viper.SetDefault("log_level", "INFO")
		viper.WriteConfig()
		PrintlnYellow("Config file created at " + configPath)
	}

	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		OpenAiApiKeyPath:      viper.GetString("openai.openai_api_key_path"),
		ServiceAccountKeyPath: viper.GetString("gcp.service_account_key_path"),
		ModelName:             viper.GetString("model"),
		GCPProjectId:          viper.GetString("gcp.gcp_project_id"),
		LogLevel:              viper.GetString("log_level"),
	}

	if appConfig.OpenAiApiKeyPath != openAiApiKeyFilePath {
		openAiApiKeyFilePath = appConfig.OpenAiApiKeyPath
	}

	if !IfFileWithAPiKeyExists(openAiApiKeyFilePath) {
		PrintlnRed("OpenAI API key file not found at " + openAiApiKeyFilePath)
		PrintlnRed("Please create a file with your OpenAI API key at")
		PrintlnRed("You can get your API key at https://beta.openai.com/account/api-keys")
		PrintlnRed("Then run `echo YOUR_OPEN_AI_API_TOKEN > " + openAiApiKeyFilePath + "`")
		os.Exit(1)
	}

	return appConfig, nil
}
