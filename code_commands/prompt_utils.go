package code_commands

import (
	"io/ioutil"
	"strings"
	"text/template"

	"net/http"
	"net/url"
	"errors"
	"github.com/go-shiori/go-readability"
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

func extractReadableTextFromURL(urlString string) (string, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download the page")
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", err
	}

	return article.TextContent, nil
}

func FilePromptBuilder(prePrompt string, filePaths []string, urls []string, userPrompt string) (string, error) {
	var files []File

	// Read file contents and populate the files slice
	for _, filePath := range filePaths {
		fileContent, err := readFileContent(filePath)
		if err != nil {
			return "", err
		}
		files = append(files, File{filePath, fileContent})
	}

	for _, url := range urls {
		urlContent, err := extractReadableTextFromURL(url)
		if err != nil {
			return "", err
		}
		files = append(files, File{url, urlContent})
	}

	// Template string
	templateStr := `
{{ .PrePrompt }}
File/Url List:
{{- range .Files }}
File/Url path: {{ .Path }}
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
