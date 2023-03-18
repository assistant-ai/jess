package chat

import (
	"time"
)

type ChatMessage struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Role      string    `json:"sender"`
	Content   string    `json:"content"`
}
