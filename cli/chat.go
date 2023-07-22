package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/assistant-ai/jess/utils"
	"os"
	"text/template"
	"time"

	"github.com/assistant-ai/jess/model"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"

	"github.com/sirupsen/logrus"
)

func AnimateThinking(quit chan bool) {
	spinChars := []rune{'-', '\\', '|', '/'}
	i := 0
	for {
		select {
		case <-quit:
			return
		default:
			utils.PrintfThinkingYellow(spinChars[i])
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

func ShowContext(context string) {
	fmt.Printf("Context message is: %s\n", context)
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

func StartChat(rawContextId string, llmClient *client.Client, logger *logrus.Logger) error {
	contextId := rawContextId
	if rawContextId == "" {
		contextId = db.RandomContextId
	}
	quit := make(chan bool)

	scanner := bufio.NewScanner(os.Stdin)

	if rawContextId != "" {
		messages := make([]db.Message, 0)
		messages, err := db.GetMessagesByContextID(contextId)
		if err != nil {
			return err
		}
		ShowMessages(messages)
	}

	for {
		utils.PrintlnCyan(">>> You:")
		if !scanner.Scan() {
			return fmt.Errorf("Error reading input")
		}
		newMessage := scanner.Text()
		if utils.IfAnswerInFinishingArray(newMessage) {
			if utils.AskForConfirmation() {
				utils.PrintlnCyan("Buy, buy! I'll miss you")
				break
			} else {
				utils.PrintlnCyan("Ok, let's continue")
				continue
			}
		}
		if newMessage == "end" {
			break
		}
		go AnimateThinking(quit)

		logger.WithFields(logrus.Fields{
			"message":   newMessage,
			"contextId": contextId,
		}).Debug("About to send a message")

		if llmClient == nil {
			logger.Fatal("llmClient is null, this should never be the case here")
		}
		response, err := llmClient.SendMessage(newMessage, contextId)
		quit <- true
		if err != nil {
			logger.Error("error from llmClient in chat dialog")
			return err
		}

		responseMsg := "\n Jess: " + response
		utils.PrintlnPurple(responseMsg)
	}

	fmt.Println("Goodbye!")
	return nil
}

func GeneratePromptForFile(input model.FileInput) (string, error) {
	fileTemplate := `
I am going to give you instructions and the content of the file, you have to output new version of the file.
What you tell me I will directly put in the file replacing everything that I showed to you, so you have to print FULL content of the file, even if you changing small part of it.
Do not forget to put correct new lines (\n) where required and correctly format the file.
Here is the instructions what to do with the file: {{.UserMessage}}
And here is file content:
{{.FileContent}}
`
	return generatePrompt(input, fileTemplate)
}
