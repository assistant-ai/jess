package auto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func StartProcess(userAsk string, rootPath string, contextId string, llmClient *client.Client, logger *logrus.Logger) error {
	memory := ""
	exist, err := db.CheckIfContextExists(contextId)
	contextSet := false
	if contextId != "" {
		logger.Debug("Context id is set to: " + contextId)
		contextSet = true
	} else {
		contextId = uuid.New().String()
		logger.Debug("Context id is set to random: " + contextId)
		defer func() {
			if !contextSet {
				db.RemoveContext(contextId)
			}
		}()
	}
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
	ctxFromDb, _ := db.GetContextMessage(contextId)
	logger.Debug("just in case, context in DB: " + ctxFromDb)
	prevCommnads := ""
	messageToSend, err := GenerateLsPrompt(userAsk, memory, rootPath, prevCommnads)
	prevCommnads = "ls"
	if err != nil {
		return err
	}
	counter := 0
	errEounter := 0
	for {
		if counter > 20 {
			fmt.Println("stop right there")
			break
		}
		counter++
		logger.Debugln("Sending message: " + messageToSend)
		logger.Debugln("With contextId: " + contextId)
		response, err := llmClient.SendMessageWithContextDepth(messageToSend, contextId, 1, false)
		if err != nil {
			return err
		}
		logger.Debugln(response)
		var cmd Command
		err = json.Unmarshal([]byte(response), &cmd)
		if err != nil {
			errEounter++
			if errEounter > 3 {
				return err
			}
			fmt.Println("Jess have not responded with proper JSON, trying to extract it.")
			err = tryToExtractJsonMessage(response, &cmd, llmClient)
			if err != nil {
				fmt.Println("Still nothing:(. Trying again..")
				continue
			}
			errEounter = 0
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
			memory = memory + "\n" + cmd.Context
			messageToSend, err = GenerateMemoryPrompt(userAsk, memory, prevCommnads)
			prevCommnads = prevCommnads + "\nmemory"
			if err != nil {
				return err
			}
		} else if cmd.Action == "update" || cmd.Action == "new" {
			fmt.Println("Jess is updating file: " + cmd.Path)
			prevCommnads = prevCommnads + "\nupdate " + cmd.Path
			err = replaceFileWithContent(cmd.Path, cmd.Context)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("WTF?!?!: " + cmd.Action)
			break
		}
	}
	return nil
}

func tryToExtractJsonMessage(message string, cmd *Command, llmClient *client.Client) error {
	answer, err := llmClient.SenRandomContextMessage("Extract JSON from the following message: " + message)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(answer), cmd)
	return err
}

func replaceFileWithContent(filePath string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filePath, data, 0644)
	return err
}
