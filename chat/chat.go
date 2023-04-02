package chat

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/assistent-ai/client/gpt"
	"github.com/assistent-ai/client/model"
)

func StartChat() {
	// Create a new scanner to read messages from the user
	scanner := bufio.NewScanner(os.Stdin)
	messages := make([]model.Message, 0)

	// Loop forever, reading messages from the user and sending them to the GPT API
	for {
		// Print a prompt to the user
		fmt.Print("You: ")

		// Read a line of text from the user
		if !scanner.Scan() {
			// If there was an error reading input, break out of the loop
			break
		}

		newMessage := model.Message{
			ID:        "", // You can assign a new ID here
			DialogId:  "", // You can assign a new DialogId here
			Timestamp: time.Now(),
			Role:      "user",
			Content:   scanner.Text(),
		}
		messages = append(messages, newMessage)
		messages, err := gpt.Message(messages)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		lastMessage := messages[len(messages)-1]

		// Format the timestamp
		formattedTimestamp := lastMessage.Timestamp.Format("2006-01-02 15:04:05")

		// Print the last message
		fmt.Printf("[%s] %s: %s\n", formattedTimestamp, lastMessage.Role, lastMessage.Content)
	}

	// If we've reached the end of input, print a goodbye message
	fmt.Println("Goodbye!")
}
