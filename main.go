package main

import (
	"time"

	"github.com/assistent-ai/client/chat"
	"github.com/assistent-ai/client/db"
	"github.com/assistent-ai/client/model"
)

func testDB() {
	db.DeleteMessage("test1")
	index, err := db.GetIndex()
	if err != nil {
		panic(err)
	}
	message := &model.Message{
		ID:        "test1",
		DialogId:  "default",
		Timestamp: time.Now(),
		Role:      "user",
		Content:   "hello",
	}
	if err != nil {
		panic(err)
	}
	index.Index(message.ID, message)
}

func main() {
	// testDB()
	chat.StartChat(model.DefaultDialogId)
}
