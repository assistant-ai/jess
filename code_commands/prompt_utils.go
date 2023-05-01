package code_commands

import (
	"io/ioutil"
	"strings"
	"text/template"
)

// File structure to store file name and content
type File struct {
	Path    string
	Content string
}

// Reads the content of the file
func readFileContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func FilePromptBuilder(prePrompt string, filePaths []string, userPrompt string) (string, error) {
	var files []File

	// Read file contents and populate the files slice
	for _, filePath := range filePaths {
		fileContent, err := readFileContent(filePath)
		if err != nil {
			return "", err
		}
		files = append(files, File{filePath, fileContent})
	}

	// Template string
	templateStr := `
{{ .PrePrompt }}
File List:
{{- range .Files }}
File path: {{ .Path }}
Content:
{{ .Content }}
{{- end }}

User Prompt: {{ .UserPrompt }}
`

	tpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var prompt strings.Builder
	data := struct {
		Files      []File
		UserPrompt string
		PrePrompt  string
	}{
		files,
		userPrompt,
		prePrompt,
	}

	if err := tpl.Execute(&prompt, data); err != nil {
		return "", err
	}

	return prompt.String(), nil
}
