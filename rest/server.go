package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"embed"
	"io"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
)

//go:embed client.html
var clientHtml embed.FS

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
	// Serve embedded client.html file on root route
	r.HandleFunc("/", serveClientHtml)

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}

func serveClientHtml(w http.ResponseWriter, r *http.Request) {
	// Get the contents of the client.html file
	file, err := clientHtml.Open("client.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set "Content-Type" header and write the contents of the file
	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
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

	// Updated: Parse JSON from the request body
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	message := requestBody["message"]
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
