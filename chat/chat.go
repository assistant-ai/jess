package chat

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"bytes"
	"text/template"

	"github.com/assistent-ai/client/gpt"
	"github.com/assistent-ai/client/model"

	"github.com/assistent-ai/client/db"
	"github.com/google/uuid"
)

func thinkingAnimation(quit chan bool) {
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

func GeneratePromptForFile(input model.FileInput) (string, error) {
	tmpl := `
I am going to give you instructions and the content of the file, you have to output new verison of the file.
What you tell me I will directly put in the file replacing everything that I showed to you, so you have to print FULL content of the file, even if you changing small part of it.
Do not forget to put correct new lines (\n) where required and correclty format the file.
Here is the instructions what to do with the file: {{.UserMessage}}
And here is file content:
{{.FileContent}}
`
	// Parse the template
	template, err := template.New("fileTemplate").Parse(tmpl)
	if err != nil {
		return "", err
	}

	// Execute the template and write the output to a buffer
	var output bytes.Buffer
	err = template.Execute(&output, input)
	if err != nil {
		return "", err
	}

	// Print the resulting string
	return output.String(), nil
}

func ShowMessages(messages []model.Message) {
	for _, message := range messages {
		formattedTimestamp := message.Timestamp.Format(model.TimestampFormattingTemplate)
		fmt.Printf("[%s] %s: %s\n", formattedTimestamp, message.Role, message.Content)
	}
}

func StartChat(dialogId string, ctx *model.AppContext) error {
	if dialogId == "" {
		dialogId = model.DefaultDialogId
	}
	quit := make(chan bool)

	// Create a new scanner to read messages from the user
	scanner := bufio.NewScanner(os.Stdin)
	messages := make([]model.Message, 0)
	if (dialogId != model.RandomDialogId) {
		messages, err := db.GetMessagesByDialogID(dialogId)
		if err != nil {
			return err
		}
		ShowMessages(messages)
	}

	for {
		// Print a prompt to the user
		fmt.Print("You: ")

		// Read a line of text from the user
		if !scanner.Scan() {
			// If there was an error reading input, break out of the loop
			break
		}
		msgUUID, err := uuid.NewRandom()
		msgId := msgUUID.String()

		newMessage := model.Message{
			ID:        msgId,
			DialogId:  dialogId,
			Timestamp: time.Now(),
			Role:      model.UserRoleName,
			Content:   scanner.Text(),
		}
		messages = append(messages, newMessage)
		go thinkingAnimation(quit)
		over, err := gpt.IsDialogOver(messages, ctx)
		if err != nil && over {
			break
		}
		messages := gpt.TrimMessages(messages, 8)
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

		// Print the last message
		fmt.Printf("\n\n\n%s: %s\n", lastMessage.Role, lastMessage.Content)
	}

	// If we've reached the end of input, print a goodbye message
	fmt.Println("Goodbye!")
	return nil
}
