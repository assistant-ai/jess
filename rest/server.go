package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
)

func StartRest(gpt *gpt.GptClient) {
	r := mux.NewRouter()

	r.HandleFunc("/contexts", getContextList).Methods("GET")
	r.HandleFunc("/contexts/{contextId}/messages", getMessagesInContext).Methods("GET")
	r.HandleFunc("/contexts/{contextId}/messages", func (w http.ResponseWriter, r *http.Request) {sendMessageToChatGPT(w, r, gpt)}).Methods("POST")

	cors := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),
        handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
    )

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}

func getContextList(w http.ResponseWriter, r *http.Request) {
	contextIds, err := db.GetContextIDs()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contextIds)
}

func getMessagesInContext(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contextId := vars["contextId"]
	messages, err := db.GetMessagesByContextID(contextId)
	if err != nil {
		return
	}

	// Implement the logic to fetch messages in a given context (using contextId).
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func sendMessageToChatGPT(w http.ResponseWriter, r *http.Request, gpt *gpt.GptClient) {
	vars := mux.Vars(r)
	contextId := vars["contextId"]

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	
	response, err := gpt.SendMessage(message, contextId)
	if err != nil {
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
