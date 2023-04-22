package utils

import (
	"github.com/assistant-ai/jess/model"
)

func TrimMessages(messages []model.Message, maxUserMessagesSize int) []model.Message {
	if len(messages) < maxUserMessagesSize*2 {
		return messages
	}
	startIndex := len(messages) - maxUserMessagesSize
	return messages[startIndex:]
}
