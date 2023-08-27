package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/urfave/cli/v2"
)

type BugHunterCommand struct{}

func (c *BugHunterCommand) Name() string {
	return "bug_hunter"
}

func (c *BugHunterCommand) Usage() string {
	return "generate of typical bugs and testcases for provided user story / task / topic"
}

func (c *BugHunterCommand) Flags() []cli.Flag {
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

func (c *BugHunterCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := "topic or user story for searching bugs: " + cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := text.BUG_HUNTER_PROMPT
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
