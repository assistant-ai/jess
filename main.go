package main

import (
	"github.com/assistent-ai/client/chat"
	"github.com/assistent-ai/client/model"
)

func main() {
	// showAll()
	// testDB()
	chat.StartChat(model.DefaultDialogId)
}
