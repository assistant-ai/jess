package utils

import (
	"fmt"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/b0noi/go-utils/v2/fs"
	"path/filepath"
	"sort"
	"strings"
)

type Config struct {
	name                  string
	nameInConfigFile      string
	defaultValue          string
	shortDescription      string
	suggestedValues       []string
	errorMessage          string
	recommendationMessage string
	note                  string
}

func (c *Config) GetName() string {
	return c.name
}

func (c *Config) GetNameInConfigFile() string {
	return c.nameInConfigFile
}

func (c *Config) GetDefaultValue() string {
	return c.defaultValue
}

func (c *Config) GetShortDescription() string {
	return c.shortDescription
}

func (c *Config) GetSuggestedValues() []string {
	return c.suggestedValues
}

func (c *Config) GetErrorMessage() string {
	return c.errorMessage
}

func (c *Config) GetRecommendationMessage() string {
	return c.recommendationMessage
}

func (c *Config) GetNote() string {
	return c.note
}

var DefaultConfigStructureMap = make(map[string]Config)

var DefaultConfigFolderPath string
var DefaultConfigFilePath string

// TODO think how to delete it
//var err = error(nil)

func init() {
	DefaultConfigFolderPath, _ = fs.MaybeCreateProgramFolder("jess")
	//if err != nil {
	//	log.Fatalf("can't create config folder: %v", err)
	//}
	DefaultConfigFilePath = filepath.Join(DefaultConfigFolderPath, "config.yaml")

	listOfModels := gpt.GetListOfModels()
	msgForSetupModels := "[IMPORTANT] Use only next models: \n" + strings.Join(listOfModels, "\n")

	logLevelConfig := Config{name: "LogLevel",
		nameInConfigFile:      "log_level",
		defaultValue:          "INFO",
		shortDescription:      "Log level. used for logging and providing addition information to various sources. By default it print everything in terminal",
		suggestedValues:       []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		errorMessage:          "",
		recommendationMessage: "",
		note:                  ""}

	modelNameConfig := Config{name: "ModelName",
		nameInConfigFile:      "model",
		defaultValue:          "gpt3Turbo",
		shortDescription:      "Model name, is used for generating text. By default it is gpt3Turbo. You can use any model from the recommendation list",
		suggestedValues:       listOfModels,
		errorMessage:          "",
		recommendationMessage: msgForSetupModels,
		note:                  ""}

	OpenAiApiKeyPathConfig := Config{name: "OpenAiApiKeyPath",
		nameInConfigFile:      "openai.openai_api_key_path",
		defaultValue:          filepath.Join(DefaultConfigFolderPath, "open-ai.key"),
		shortDescription:      "Place where your open ai key is stored. By default it is stored in ~/.jess/open-ai.key",
		suggestedValues:       []string{"~/.jess/open-ai.key", "/home/user/.jess/open-ai.key"},
		errorMessage:          "",
		recommendationMessage: "try to use your home path, don't use spaces in path. don't use relative path",
		note:                  ""}

	JessEmailAccountConfig := Config{name: "JessEmailAccount",
		nameInConfigFile:      "gcp.jessica_mail_account",
		defaultValue:          "",
		shortDescription:      "Jessica mail account, that would be used for various purposes. No default values for this field, due to security reasons",
		suggestedValues:       []string{""},
		errorMessage:          "",
		recommendationMessage: "",
		note:                  ""}

	JessServiceAccountKeyPathConfig := Config{name: "ServiceAccountKeyPath",
		nameInConfigFile:      "gcp.service_account_key_path",
		defaultValue:          filepath.Join(DefaultConfigFolderPath, "jess-service-account.json"),
		shortDescription:      "This is a path service account key is used to connect jess to gcp for using google documents API. No default values for this field, due to security reasons. By default it is stored in ~/.jess/service-account-key.json, however this file doesn't exist by default",
		suggestedValues:       []string{""},
		errorMessage:          "",
		recommendationMessage: "",
		note:                  ""}

	GCPProjectIdConfig := Config{name: "GCPProjectId",
		nameInConfigFile:      "gcp.gcp_project_id",
		defaultValue:          "",
		shortDescription:      "This is a project id that is used to connect jess to GCP projects API for usage PALM model. No default values for this field, due to security reasons",
		suggestedValues:       []string{""},
		errorMessage:          "",
		recommendationMessage: "",
		note:                  ""}

	DefaultConfigStructureMap[logLevelConfig.name] = logLevelConfig
	DefaultConfigStructureMap[modelNameConfig.name] = modelNameConfig
	DefaultConfigStructureMap[OpenAiApiKeyPathConfig.name] = OpenAiApiKeyPathConfig
	DefaultConfigStructureMap[JessEmailAccountConfig.name] = JessEmailAccountConfig
	DefaultConfigStructureMap[JessServiceAccountKeyPathConfig.name] = JessServiceAccountKeyPathConfig
	DefaultConfigStructureMap[GCPProjectIdConfig.name] = GCPProjectIdConfig
}

// this is mainly service method that should be used for debugging
func PrintConfig() {
	for _, config := range DefaultConfigStructureMap {
		fmt.Println("Name: ", config.GetName())
		fmt.Println("------ NameInConfigFile: \t\t", config.GetNameInConfigFile())
		fmt.Println("------ DefaultValue: \t\t\t", config.GetDefaultValue())
		fmt.Println("------ ShortDescription: \t\t", config.GetShortDescription())
		fmt.Println("------ SuggestedValues: \t\t", config.GetSuggestedValues())
		fmt.Println("------ ErrorMessage: \t\t\t", config.GetErrorMessage())
		fmt.Println("------ RecommendationMessage: \t\t", config.GetRecommendationMessage())
		fmt.Println("------ Note: \t\t\t", config.GetNote())
		fmt.Println("")
	}
}

// Got get config name for initiating config
func GetConfigMap() map[string]Config {
	return DefaultConfigStructureMap
}

func GetListOfConfigFields() []string {
	var listOfConfigNames []string
	for _, config := range DefaultConfigStructureMap {
		listOfConfigNames = append(listOfConfigNames, config.GetName())
	}
	sort.Strings(listOfConfigNames)
	return listOfConfigNames
}

func GetDefaultConfigFilePath() string {
	return DefaultConfigFilePath
}

func GetDefaultConfigFolderPath() string {
	return DefaultConfigFolderPath
}
