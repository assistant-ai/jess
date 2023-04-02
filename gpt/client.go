package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/assistent-ai/client/model"
	"github.com/b0noi/go-utils/v2/gcp"
)

func Message(messages []model.Message) ([]model.Message, error) {
	apiKey, err := gcp.AccessSecretVersion("projects/16255416068/secrets/gpt3-secret/versions/1")
	if err != nil {
		return nil, err
	}

	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := prepareGPT4RequestBody(messages)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// var response map[string]interface{}
	var response GptChatCompletionMessage
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return addGPT4Response(response, messages)
}

func addGPT4Response(response GptChatCompletionMessage, messages []model.Message) ([]model.Message, error) {
	// Assume we're only getting 1 response, so we use the first choice
	gpt4Text := response.Choices[0].Message.Content

	// Create a new chat.Message with the GPT-4 response
	newMessage := model.Message{
		ID:        "", // You can assign a new ID here
		DialogId:  "", // You can assign a new DialogId here
		Timestamp: time.Now(),
		Role:      "assistent",
		Content:   gpt4Text,
	}

	messages = append(messages, newMessage)

	// Append the new message to the input messages slice
	return messages, nil
}

func prepareGPT4RequestBody(messages []model.Message) ([]byte, error) {
	// Create a new slice to hold message maps
	gptMessages := make([]map[string]string, len(messages))

	// Iterate through the input messages
	for i, message := range messages {
		// Convert the timestamp to a human-readable format
		formattedTimestamp := message.Timestamp.Format("2006-01-02 15:04:05")

		// Combine the content with the timestamp
		combinedContent := fmt.Sprintf("%s: %s", formattedTimestamp, message.Content)

		// Add the message to the gptMessages slice
		gptMessages[i] = map[string]string{
			"role":    message.Role,
			"content": combinedContent,
		}
	}

	// Marshal the request body for GPT-4
	requestBody, err := json.Marshal(map[string]interface{}{
		"messages":   gptMessages,
		"max_tokens": 2000,
		"n":          1,
		"model":      "gpt-4",
	})

	if err != nil {
		return nil, err
	}

	return requestBody, nil
}
