package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/b0noi/go-utils/v2/gcp"
)

func MessageGptInDefaultConversation(message string) (string, error) {
	apiKey, err := gcp.AccessSecretVersion("projects/16255416068/secrets/gpt3-secret/versions/1")
	if err != nil {
		return "", nil
	}

	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
		"max_tokens": 2000,
		"n":          1,
		"model":      "gpt-4",
	})

	if err != nil {
		return "", nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// var response map[string]interface{}
	var response GptChatCompletionMessage
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
