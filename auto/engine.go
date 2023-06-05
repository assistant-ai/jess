package auto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/assistant-ai/llmchat-client/client"
	"github.com/assistant-ai/llmchat-client/db"
	"github.com/sirupsen/logrus"
	"github.com/b0noi/go-utils/v2/fs"
)

func StartProcess(userAsk string, rootPath string, contextId string, llmClient *client.Client, logger *logrus.Logger) error {
	memory := ""
	exist, err := db.CheckIfContextExists(contextId)
	if err != nil {
		return err
	}
	if !exist {
		err = db.CreateContext(contextId, "")
		if err != nil {
			return err
		}
	} else {
		err = db.UpdateContext(contextId, "")
		if err != nil {
			return err
		}
	}
	prevCommnads := ""
	messageToSend, err := GenerateLsPrompt(userAsk, memory, rootPath, prevCommnads)
	fmt.Println("Jess is asking about ls")
	prevCommnads = "ls"
	if err != nil {
		return err
	}
	counter := 0
	err_counter := 0
	for {
		if counter > 10 {
			fmt.Println("stop right there")
			break
		}
		counter++
		if !strings.Contains(messageToSend, SystemContext) {
			messageToSend = SystemContext + "\n" + messageToSend
		}
		logger.Debugln(messageToSend)
		response, err := llmClient.SendMessage(messageToSend, contextId)
		if err != nil {
			return err
		}
		logger.Debugln(response)
		var cmd Command
		if cmd.Memory != "" {
			memory = cmd.Memory
		}
		if memory != "" {
			fmt.Println("Memory: " + memory)
		}
		err = json.Unmarshal([]byte(response), &cmd)
		if err != nil {
			err_counter++
			if err_counter > 3 {
				return err
			}
			fmt.Println("Jess have not responded with proper JSON, but I will give her one more chancse.")
			continue
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

func replaceFileWithContent(filePath string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filePath, data, 0644)
	return err
}
