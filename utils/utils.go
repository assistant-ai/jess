package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/b0noi/go-utils/v2/fs"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

func ExtractTextFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	readableText := doc.Find("body").Text()
	return readableText, nil
}

func AnswersOutput(output string, answer string) error {
	if output != "" {
		err := os.WriteFile(output, []byte(answer), 0644)
		if err != nil {
			return err
		}
	} else {
		PrintlnCyan(answer)
		return nil
	}
	return nil
}

func IsGitRepository(folderPath string) error {
	folderPath, err := ExpandTilde(folderPath)
	if err != nil {
		return err
	}
	gitDir := filepath.Join(folderPath, ".git")
	_, err = os.Stat(gitDir)
	if err != nil {
		log.Errorf("Look like provided folder is not git repository. give us another folder: %v", err)
		return err
	}
	return nil
}

// TODO reuse this function to fix other error with path
func ExpandTilde(path string) (string, error) {
	if path[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, nil
}

func isFolderPath(path string) (bool, error) {
	path, err := ExpandTilde(path)
	if err != nil {
		return false, err
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("%s does not exist", path)
		}
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func IsValidPath(path string) (bool, error) {
	isFolder, err := isFolderPath(path)
	if err != nil {
		return false, err
	}
	if !isFolder {
		return false, fmt.Errorf("%s is not a folder", path)
	}
	return true, nil
}

// TODO fix error handling and returnig erorrs
func IfFileWithAPiKeyExists(apiKeyPath string) bool {
	exists, _ := fs.PathExists(apiKeyPath)
	return exists
}

// TODO fix error handling and returnig erors
func IfConfigFileExists(configPath string) bool {
	exists, err := fs.PathExists(configPath)
	if err != nil {
		log.Fatalf("%s Config file not found at: %s", err, configPath)
	}
	return exists
}

func IsServiceAccountJsonFileExists(serviceAccountKeyPath string) (bool, error) {
	exists, err := fs.PathExists(serviceAccountKeyPath)
	if err != nil {
		log.Fatalf("%s Config file not found at: %s", err, serviceAccountKeyPath)
	}
	return exists, nil
}

func PrintPrompt(showPrompt bool, basicPrompt string, finalPrompt string) {
	if showPrompt {
		PrintlnYellow("BASIC_PROMPT:\n\n" + basicPrompt + "\n\n")
		PrintlnYellow("FINAL_PROMPT:\n\n" + finalPrompt + "\n\n")
	}
}

func getApiKeyFromFile(OpenAiApiKeyPath string) string {
	apiKey, err := os.ReadFile(OpenAiApiKeyPath)
	if err != nil {
		log.Fatalf("%s Error reading OpenAI API key file: %s", err, OpenAiApiKeyPath)
	}
	return string(apiKey)
}

func GetMaskedApiKey(OpenAiApiKeyPath string) string {
	apiKey := getApiKeyFromFile(OpenAiApiKeyPath)
	maskedKey := string(apiKey[0:5]) + "..." + string(apiKey[len(apiKey)-5:len(apiKey)])
	return maskedKey
}

func isValidURL(input string) bool {
	urlPattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
	regex := regexp.MustCompile(urlPattern)
	return regex.MatchString(input)
}

func PrintlnGreen(msg string) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), msg, string(colorReset))
}

func PrintlnRed(msg string) {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	fmt.Println(string(colorRed), msg, string(colorReset))
}

func PrintlnCyan(msg string) {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"
	fmt.Println(string(colorCyan), msg, string(colorReset))
}

func PrintCyanInvite() {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"
	fmt.Print(string(colorCyan), " >>> ", string(colorReset))
}

func PrintfThinkingYellow(msg rune) {
	fmt.Printf("\r\033[33mThinking %c \033[0m", msg)
}

func PrintlnYellow(msg string) {
	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	fmt.Println(string(colorYellow), msg, string(colorReset))
}

func PrintlnPurple(msg string) {
	colorReset := "\033[0m"
	colorPurple := "\033[35m"
	fmt.Println(string(colorPurple), msg, string(colorReset))
}
