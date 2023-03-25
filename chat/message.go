package chat

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	DialogId  string    `json:"dialog_id"`
	Timestamp time.Time `json:"timestamp"`
	Role      string    `json:"sender"`
	Content   string    `json:"content"`
}
