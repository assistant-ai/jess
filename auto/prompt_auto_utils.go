package auto

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenerateMemoryPrompt(userAsk string, memory string, prevCommands string) (string, error) {
	promptData := struct {
		UserAsk         string
		Memory          string
		OperationName   string
		OperationResult string
		PrevActions     string
		NextPrompt      string
	}{
		userAsk,
		memory,
		"memory " + memory,
		"memory saved",
		prevCommands,
		"what would be the next step? Highly consider ls command unless you want to ask about the file, path of which you see in the memory",
	}
	tmpl, err := template.New("StepPrompt").Parse(StepPromptTemplate)
	if err != nil {
		return "", err
	}
	var prompt strings.Builder
	if err := tmpl.Execute(&prompt, promptData); err != nil {
		return "", err
	}
	return prompt.String(), nil
}

func GenerateCatPrompt(userAsk string, memory, filePath string, prevCommands string) (string, error) {
	fileContent, err := readFileContent(filePath)
	if err != nil {
		return "", err
	}
	promptData := struct {
		UserAsk         string
		Memory          string
		OperationName   string
		OperationResult string
		PrevActions     string
		NextPrompt      string
	}{
		userAsk,
		memory,
		"cat " + filePath,
		fileContent,
		prevCommands,
		"next, store learning from this file that you will need to finish user task in your memory by using memory command",
	}
	tmpl, err := template.New("StepPrompt").Parse(StepPromptTemplate)
	if err != nil {
		return "", err
	}
	var prompt strings.Builder
	if err := tmpl.Execute(&prompt, promptData); err != nil {
		return "", err
	}
	return prompt.String(), nil
}

// Reads the content of the file
func readFileContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GenerateLsPrompt(userAsk string, memory string, projectRootPath string, prevCommands string) (string, error) {
	files, error := listOfFiles(projectRootPath)
	if error != nil {
		return "", error
	}
	promptData := struct {
		UserAsk         string
		Memory          string
		OperationName   string
		OperationResult string
		PrevActions     string
		NextPrompt      string
	}{
		userAsk,
		memory,
		"ls",
		listOfFilesToString(files),
		prevCommands,
		"what would be the next step?",
	}
	tmpl, err := template.New("StepPrompt").Parse(StepPromptTemplate)
	if err != nil {
		return "", err
	}
	var prompt strings.Builder
	if err := tmpl.Execute(&prompt, promptData); err != nil {
		return "", err
	}
	return prompt.String(), nil
}

func listOfFilesToString(files []string) string {
	var result strings.Builder
	for _, file := range files {
		result.WriteString(file)
		result.WriteString("\n")
	}
	return result.String()
}

func listOfFiles(rootPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, ".git") || strings.Contains(path, ".idea") || strings.Contains(path, ".vscode") || strings.Contains(path, ".DS_Store") {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
