package chat

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/assistent-ai/client/gpt"
	"github.com/assistent-ai/client/model"

	"github.com/assistent-ai/client/db"
	"github.com/google/uuid"
)

func StartChat(dialogId string) error {
	if dialogId == "" {
		dialogId = model.DefaultDialogId
	}

	// Create a new scanner to read messages from the user
	scanner := bufio.NewScanner(os.Stdin)
	index, err := db.GetIndex()
	if err != nil {
		return err
	}
	defer index.Close()
	messages, err := db.GetMessagesByDialogId(dialogId, index)
	if err != nil {
		return err
	}
	for _, message := range messages {
		formattedTimestamp := message.Timestamp.Format(model.TimestampFormattingTemplate)
		fmt.Printf("[%s] %s: %s\n", formattedTimestamp, message.Role, message.Content)
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
		over, err := gpt.IsDialogOver(messages)
		if err != nil {
			return err
		}
		if over {
			break
		}
		messages, err := gpt.Message(messages, dialogId)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		if err = index.Index(newMessage.ID, newMessage); err != nil {
			return err
		}
		lastMessage := messages[len(messages)-1]
		if err = index.Index(lastMessage.ID, lastMessage); err != nil {
			return err
		}

		// Print the last message
		fmt.Printf("%s: %s\n", lastMessage.Role, lastMessage.Content)
	}

	// If we've reached the end of input, print a goodbye message
	fmt.Println("Goodbye!")
	return nil
}
