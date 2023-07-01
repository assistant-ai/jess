package utils

import (
	"os"
	"path/filepath"

	"github.com/b0noi/go-utils/v2/fs"
	"github.com/spf13/viper"
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
	jessServiceAccountFilePath := filepath.Join(folderPath, "jess-service-account.json")
	viper.SetDefault("openai.openai_api_key_path", openAiApiKeyFilePath)
	viper.SetDefault("gcp.service_account_key_path", jessServiceAccountFilePath)
	viper.SetDefault("model", "GPT3")
	viper.SetDefault("gcp.gcp_project_id", "")
	viper.SetDefault("log_level", "INFO")

	exists, err := fs.PathExists(configPath)

	if err != nil {
		return nil, err
	}

	if !exists {
		file, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
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

	return appConfig, nil
}
