package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/assistant-ai/jess/db"
	"github.com/assistant-ai/jess/gpt"
	"github.com/assistant-ai/jess/model"
	"github.com/assistant-ai/jess/utils"
	"github.com/google/uuid"
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

func ShowMessages(messages []model.Message) {
	for _, message := range messages {
		formattedTimestamp := message.Timestamp.Format(model.TimestampFormattingTemplate)
		fmt.Printf("[%s] %s: %s\n", formattedTimestamp, message.Role, message.Content)
	}
}

func readNewMessage(scanner *bufio.Scanner, dialogId string) (model.Message, error) {
	fmt.Print("You: ")

	if !scanner.Scan() {
		return model.Message{}, fmt.Errorf("Error reading input")
	}
	msgUUID, err := uuid.NewRandom()
	if err != nil {
		return model.Message{}, err
	}
	msgId := msgUUID.String()

	return model.Message{
		ID:        msgId,
		DialogId:  dialogId,
		Timestamp: time.Now(),
		Role:      model.UserRoleName,
		Content:   scanner.Text(),
	}, nil
}

func StartChat(dialogId string, ctx *model.AppContext) error {
	if dialogId == "" {
		dialogId = model.DefaultDialogId
	}
	quit := make(chan bool)

	scanner := bufio.NewScanner(os.Stdin)
	messages := make([]model.Message, 0)
	if dialogId != model.RandomDialogId {
		messages, err := db.GetMessagesByDialogID(dialogId)
		if err != nil {
			return err
		}
		ShowMessages(messages)
	}

	for {
		newMessage, err := readNewMessage(scanner, dialogId)
		if err != nil {
			break
		}

		messages = append(messages, newMessage)
		go AnimateThinking(quit)

		over, err := gpt.IsDialogOver(messages, ctx)
		if err != nil && over {
			break
		}
		messages = utils.TrimMessages(messages, 8)
		messages, err = gpt.Message(messages, dialogId, ctx)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		if _, err = db.StoreMessage(newMessage); err != nil {
			return err
		}
		lastMessage := messages[len(messages)-1]
		if _, err = db.StoreMessage(lastMessage); err != nil {
			return err
		}
		quit <- true

		fmt.Printf("\n\n\n%s: %s\n", lastMessage.Role, lastMessage.Content)
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
