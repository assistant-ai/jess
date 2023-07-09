package utils

import (
	"strings"
)

type GPTModel struct {
	Name      string `json:"name"`
	MaxTokens int    `json:"max_tokens"`
}

var gpt3turbo = &GPTModel{Name: "gpt-3.5-turbo", MaxTokens: 4000}
var gpt3TurboBig = &GPTModel{Name: "gpt-3.5-turbo-16k", MaxTokens: 16000}
var gpt4 = &GPTModel{Name: "gpt-4", MaxTokens: 8000}
var gpt4Big = &GPTModel{Name: "gpt-4-32k", MaxTokens: 32000}

func GetModels() map[string]*GPTModel {
	Models := make(map[string]*GPTModel)
	Models["gpt3Turbo"] = gpt3turbo
	Models["gpt3TurboBig"] = gpt3TurboBig
	Models["gpt4"] = gpt4
	Models["gpt4Big"] = gpt4Big
	return Models
}

func GetListOfModels() []string {
	modelsMap := GetModels()
	keys := make([]string, 0, len(modelsMap))
	for key := range modelsMap {
		keys = append(keys, key)
	}
	return keys
}

func NotInList(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return false
		}
	}
	return true
}

func IsModelGpt(modelName string) bool {
	contains := strings.Contains(strings.ToLower(modelName), "gpt")
	return contains
}

func IsModelGPTValid(modelName string, listOfModels []string) bool {
	if NotInList(modelName, listOfModels) {
		errorText := "Model name is not in the list of valid models\n It is CASE sensitive"
		yourModelMsg := "You entered : " + modelName
		validModelsMsg := "List of valid models :" + strings.Join(listOfModels, ", ")
		changeModelMsg := "Change model name in config.json file\n or use: \n jess config -c \"id\" \n or \n Deleting ~/.jess/config.json file - also might help"
		PrintlnRed(errorText)
		PrintlnYellow(yourModelMsg)
		PrintlnYellow(validModelsMsg)
		PrintlnYellow(changeModelMsg)
		return false
	}
	if !IsModelGpt(modelName) {
		return false
	}
	return true
}
