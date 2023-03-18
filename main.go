package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/b0noi/go-utils/gcp"
)

func main() {
	apiKey, err := gcp.AccessSecretVersion("projects/16255416068/secrets/gpt3-secret/versions/1")
	if err != nil {
		panic(err)
	}

	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "hello"},
		},
		"max_tokens": 2000,
		"n":          1,
		"model":      "gpt-4",
	})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
