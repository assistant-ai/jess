package auto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/b0noi/go-utils/v2/fs"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	fmt.Println("Jess is asking about ls")
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
		response, err := llmClient.SendMessageWithContextDepth(messageToSend, contextId, 0, false)
		if err != nil {
			return err
		}
		logger.Debugln(response)
		var cmd Command
		if memory != "" {
			fmt.Println("Memory: " + memory)
		}
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
			exists, err := fs.PathExists(cmd.Path)
			if err != nil {
				return err
			}
			if !exists {
				prevCommnads = prevCommnads + "\ncat " + cmd.Path + " - ERROR, no such file"
				fmt.Println("but there is NO such file, so continuing")
				continue
			}
			messageToSend, err = GenerateCatPrompt(userAsk, memory, cmd.Path, prevCommnads)
			prevCommnads = prevCommnads + "\ncat " + cmd.Path
			if err != nil {
				return err
			}
		} else if cmd.Action == "memory" {
			fmt.Println("Jess is asking update memory with: " + cmd.Context)
			if memory == "" {
				memory = cmd.Context
			} else {
				memoryUpdate, err := joinMemoryOldAndNew(memory, cmd.Context, llmClient)
				if err != nil {
					return err
				}
				memory = memoryUpdate
			}
			fmt.Println("Memory after the update: " + memory)
			messageToSend, err = GenerateMemoryPrompt(userAsk, memory, prevCommnads)
			prevCommnads = prevCommnads + "\nmemory"
			if err != nil {
				return err
			}
		} else if cmd.Action == "update" || cmd.Action == "new" {
			fmt.Println("Jess is updating file: " + cmd.Path)
			if cmd.Action == "update" {
				prevCommnads = prevCommnads + "\nupdate " + cmd.Path
			} else if cmd.Action == "new" {
				prevCommnads = prevCommnads + "\nnew " + cmd.Path
			}

			err = replaceFileWithContent(cmd.Path, cmd.Context)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("WTF?!?!: " + cmd.Action)
			break
		}
		time.Sleep(10000 * time.Millisecond)
	}
	return nil
}

func joinMemoryOldAndNew(oldMemory string, newMemory string, llmClient *client.Client) (string, error) {
	answer, err := llmClient.SendRandomContextMessage("There are old memory and new memory, you have to join them in the way that will capture all the same information with all the details but result message in size should not be more than 500 words. Old memory: " + oldMemory + "\n\n new memory: " + newMemory)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func tryToExtractJsonMessage(message string, cmd *Command, llmClient *client.Client) error {
	answer, err := llmClient.SendRandomContextMessage("Extract JSON from the following message: " + message)
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
