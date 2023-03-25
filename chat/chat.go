package chat

import (
	"bufio"
	"fmt"
	"os"

	"github.com/assistent-ai/client/gpt"
)

func StartChat() {
	// Create a new scanner to read messages from the user
	scanner := bufio.NewScanner(os.Stdin)

	// Loop forever, reading messages from the user and sending them to the GPT API
	for {
		// Print a prompt to the user
		fmt.Print("You: ")

		// Read a line of text from the user
		if !scanner.Scan() {
			// If there was an error reading input, break out of the loop
			break
		}

		// Send the user's message to the GPT API and print the response
		response, err := gpt.MessageGptInDefaultConversation(scanner.Text())
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		fmt.Printf("Bot: %s\n", response)
	}

	// If we've reached the end of input, print a goodbye message
	fmt.Println("Goodbye!")
}
