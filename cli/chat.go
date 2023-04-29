package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/assistant-ai/jess/model"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
)

func AnimateThinking(quit chan bool) {
	spinChars := []rune{'-', '\\', '|', '/'}
	i := 0
	for {
		select {
		case <-quit:
			return
		default:
			fmt.Printf("\rThinking %c", spinChars[i])
			time.Sleep(100 * time.Millisecond)

			i++
			if i == len(spinChars) {
				i = 0
			}
		}
	}
}

func generatePrompt(input model.FileInput, tmplFile string) (string, error) {
	template, err := template.New("fileTemplate").Parse(tmplFile)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = template.Execute(&output, input)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func ShowMessages(messages []db.Message) {
	for _, message := range messages {
		formattedTimestamp := message.Timestamp.Format(model.TimestampFormattingTemplate)
		fmt.Printf("[%s] %s: %s\n", formattedTimestamp, message.Role, message.Content)
	}
}

func readNewMessage(scanner *bufio.Scanner) (string, error) {
	fmt.Print("You: ")

	if !scanner.Scan() {
		return "", fmt.Errorf("Error reading input")
	}
	return scanner.Text(), nil
}

func StartChat(contextId string, gpt *gpt.GptClient) error {
	quit := make(chan bool)

	scanner := bufio.NewScanner(os.Stdin)
	messages := make([]db.Message, 0)
	messages, err := db.GetMessagesByContextID(contextId)
	if err != nil {
		return err
	}
	ShowMessages(messages)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			return fmt.Errorf("Error reading input")
		}
		newMessage := scanner.Text()
		if newMessage == "end" {
			break
		}
		if err != nil {
			break
		}
		go AnimateThinking(quit)

		response, err := gpt.SendMessage(newMessage, contextId)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		quit <- true

		fmt.Printf("\n\n\nJess: %s\n\n\n", response)
	}

	fmt.Println("Goodbye!")
	return nil
}

func GeneratePromptForFile(input model.FileInput) (string, error) {
	fileTemplate := `
I am going to give you instructions and the content of the file, you have to output new verison of the file.
What you tell me I will directly put in the file replacing everything that I showed to you, so you have to print FULL content of the file, even if you changing small part of it.
Do not forget to put correct new lines (\n) where required and correclty format the file.
Here is the instructions what to do with the file: {{.UserMessage}}
And here is file content:
{{.FileContent}}
`
	return generatePrompt(input, fileTemplate)
}
