package prompt

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

func GenerateMultiFileProcessPrompt(filePaths []string, userPrompt string) (string, error) {
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
Let me show you files and than I will show you my prompt to use, it might include questions/asks about the files.
File List:
{{- range .Files }}
\nFile path: {{ .Path }}
Content:
{{ .Content }}
{{- end }}

User Prompt: {{ .UserPrompt }}\n
`

	tpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var prompt strings.Builder
	data := struct {
		Files      []File
		UserPrompt string
	}{
		files,
		userPrompt,
	}

	if err := tpl.Execute(&prompt, data); err != nil {
		return "", err
	}

	return prompt.String(), nil
}
