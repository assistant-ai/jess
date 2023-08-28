package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/urfave/cli/v2"
)

type UserStoryCommand struct{}

func (c *UserStoryCommand) Name() string {
	return "user_story"
}

func (c *UserStoryCommand) Usage() string {
	return "generate description of user story based on the topic"
}

func (c *UserStoryCommand) Flags() []cli.Flag {
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

func (c *UserStoryCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := "user story or topic: " + cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := text.UserStoryPrompt
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
