package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/urfave/cli/v2"
)

type TechTaskCommand struct{}

func (c *TechTaskCommand) Name() string {
	return "tech_task"
}

func (c *TechTaskCommand) Usage() string {
	return "generate of technical task for provided topic"
}

func (c *TechTaskCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[Optional] Add information about your user story",
			Value:    "",
			Required: false,
		},
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *TechTaskCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := "topic or user story for generating technical task: " + cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := text.TECH_TASK_PROMPT
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
