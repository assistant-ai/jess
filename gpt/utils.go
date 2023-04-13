package gpt

import (
	"unicode/utf8"

	"github.com/assistent-ai/client/model"
)

func messageSize(message model.Message) int {
	totalSize := 0
	totalSize += utf8.RuneCountInString(message.Content)
	totalSize += utf8.RuneCountInString(message.Role)
	totalSize += 20 // for timestamp
	return totalSize
}

func reverseMessages(messages []model.Message) []model.Message {
	length := len(messages)
	for i := 0; i < length/2; i++ {
		j := length - 1 - i
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages
}

func TrimMessages(messages []model.Message, maxContentSize int) []model.Message {
	totalSize := 0
	for _, message := range messages {
		totalSize += messageSize(message)
	}

	if totalSize <= maxContentSize {
		return messages
	}

	trimmedMessages := make([]model.Message, 0)
	remainingSize := 0
	for _, message := range reverseMessages(messages) {
		messageSize := utf8.RuneCountInString(message.Content)
		if remainingSize+messageSize <= maxContentSize {
			remainingSize += messageSize
			trimmedMessages = append(trimmedMessages, message)
		} else {
			break
		}
	}

	return trimmedMessages
}
