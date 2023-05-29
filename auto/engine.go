package auto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/sirupsen/logrus"
)

func StartProcess(userAsk string, rootPath string, contextId string, llmClient *client.Client, logger *logrus.Logger) error {
	memory := ""
	exist, err := db.CheckIfContextExists(contextId)
	if err != nil {
		return err
	}
	if !exist {
		err = db.CreateContext(contextId, SystemContext)
		if err != nil {
			return err
		}
	} else {
		err = db.UpdateContext(contextId, SystemContext)
		if err != nil {
			return err
		}
	}
	prevCommnads := ""
	messageToSend, err := GenerateLsPrompt(userAsk, memory, rootPath, prevCommnads)
	prevCommnads = "ls"
	if err != nil {
		return err
	}
	counter := 0
	for {
		if counter > 6 {
			fmt.Println("stop right there")
			break
		}
		counter++
		messageToSend = SystemContext + "\n" + messageToSend
		logger.Debugln(messageToSend)
		response, err := llmClient.SendMessage(messageToSend, contextId)
		if err != nil {
			return err
		}
		logger.Debugln(response)
		var cmd Command
		err = json.Unmarshal([]byte(response), &cmd)
		if err != nil {
			return err
		}
		if cmd.Action == "end" {
			fmt.Println("Fuck yeah!")
			break
		} else if cmd.Action == "ls" {
			fmt.Println("Jess is asking about ls")
			messageToSend, err = GenerateLsPrompt(userAsk, memory, rootPath, prevCommnads)

			prevCommnads = prevCommnads + "\nls"
			if err != nil {
				return err
			}
		} else if cmd.Action == "cat" {
			fmt.Println("Jess is asking about content of the file: " + cmd.Path)
			messageToSend, err = GenerateCatPrompt(userAsk, memory, cmd.Path, prevCommnads)
			prevCommnads = prevCommnads + "\ncat " + cmd.Path
			if err != nil {
				return err
			}
		} else if cmd.Action == "memory" {
			fmt.Println("Jess is asking update memeory with: " + cmd.Context)
			memory = cmd.Context
			messageToSend, err = GenerateMemoryPrompt(userAsk, memory, prevCommnads)
			prevCommnads = prevCommnads + "\nmemory"
			if err != nil {
				return err
			}
		} else if cmd.Action == "update" {
			fmt.Println("Jess is updating file: " + cmd.Path)
			prevCommnads = prevCommnads + "\nupdate " + cmd.Path
			err = replaceFileWithContent(cmd.Path, cmd.Context)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("WTF?!?!")
			break
		}
	}
	return nil
}

func replaceFileWithContent(filePath string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filePath, data, 0644)
	return err
}
