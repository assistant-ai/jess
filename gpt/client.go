package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/assistent-ai/client/model"
	"github.com/google/uuid"
)

func IsDialogOver(messages []model.Message, ctx *model.AppContext) (bool, error) {
	newMessage := createNewMessage(model.SystemRoleName, "Based on the last response from the user, is this dialog over? Please respond with true/false only")
	messages = append(messages, newMessage)

	requestBody, err := prepareGPT4RequestBody(messages, ModelGPT3Turbo)
	if err != nil {
		return false, err
	}

	response, err := sendGPTRequest(requestBody, ctx)
	if err != nil {
		return false, err
	}

	result, err := strconv.ParseBool(response.Choices[0].Message.Content)
	if err != nil {
		return false, err
	}
	return result, nil
}

func RandomMessage(message string, ctx *model.AppContext) (string, error) {
	newMessage := createNewMessage(model.UserRoleName, message)
	messages := []model.Message{newMessage}

	response, err := Message(messages, model.RandomDialogId, ctx)
	if err != nil {
		return "", err
	}
	return response[1].Content, nil
}

func Message(messages []model.Message, dialogId string, ctx *model.AppContext) ([]model.Message, error) {
	requestBody, err := prepareGPT4RequestBody(messages, ModelGPT4)
	if err != nil {
		return nil, err
	}

	response, err := sendGPTRequest(requestBody, ctx)
	if err != nil {
		return nil, err
	}

	return addGPT4Response(response, messages, dialogId)
}

func createNewMessage(role, content string) model.Message {
	uuidMsg, _ := uuid.NewUUID()
	idMsg := uuidMsg.String()

	return model.Message{
		ID:        idMsg,
		DialogId:  "",
		Timestamp: time.Now(),
		Role:      role,
		Content:   content,
	}
}

func sendGPTRequest(requestBody []byte, ctx *model.AppContext) (*GptChatCompletionMessage, error) {
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.OpenAiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GptChatCompletionMessage
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("Error response from GPT")
	}

	return &response, nil
}

func addGPT4Response(response *GptChatCompletionMessage, messages []model.Message, dialogId string) ([]model.Message, error) {
	gpt4Text := response.Choices[0].Message.Content
	newMessage := createNewMessage(model.AssistentRoleNeam, gpt4Text)
	newMessage.DialogId = dialogId
	messages = append(messages, newMessage)

	return messages, nil
}

func prepareGPT4RequestBody(messages []model.Message, model GPTModel) ([]byte, error) {
	gptMessages := convertMessagesToMaps(messages)

	requestBody, err := json.Marshal(map[string]interface{}{
		"messages":   gptMessages,
		"max_tokens": 2000,
		"n":          1,
		"model":      model,
	})

	if err != nil {
		return nil, err
	}

	return requestBody, nil
}

func convertMessagesToMaps(messages []model.Message) []map[string]string {
	gptMessages := make([]map[string]string, len(messages))

	for i, message := range messages {
		formattedTimestamp := message.Timestamp.Format("2006-01-02 15:04:05")
		combinedContent := fmt.Sprintf("%s: %s", formattedTimestamp, message.Content)

		gptMessages[i] = map[string]string{
			"role":    message.Role,
			"content": combinedContent,
		}
	}

	return gptMessages
}
