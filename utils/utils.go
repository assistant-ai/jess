package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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

func Println_green(msg string) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), msg, string(colorReset))
}

func Println_red(msg string) {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	fmt.Println(string(colorRed), msg, string(colorReset))
}

func Println_cyan(msg string) {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"
	fmt.Println(string(colorCyan), msg, string(colorReset))
}
func Printf_thinking_yellow(msg rune) {
	//colorReset := "\033[0m"
	//colorYellow := "\033[33m"
	fmt.Printf("\r\033[33mThinking %c \033[0m", msg)
}

func Println_yellow(msg string) {
	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	fmt.Println(string(colorYellow), msg, string(colorReset))
}

func Println_purple(msg string) {
	colorReset := "\033[0m"
	colorPurple := "\033[35m"
	fmt.Println(string(colorPurple), msg, string(colorReset))
}
