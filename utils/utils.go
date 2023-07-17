package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/b0noi/go-utils/v2/fs"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
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
