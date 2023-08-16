package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/prompt_storage/text"
	"github.com/assistant-ai/jess/utils"
	"github.com/urfave/cli/v2"
)

type TLDR struct{}

func (c *TLDR) Name() string {
	return "tldr"
}

func (c *TLDR) Usage() string {
	return "summarise big chunks of texts"
}

func (c *TLDR) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[optional] Add information about your problem",
			Value:    "",
			Required: false,
		},
		// Experimental, not sure if this is a good idea. but during debug of personal prompts it's useful
		&cli.BoolFlag{
			Name:     "showPrompt",
			Aliases:  []string{"sp"},
			Usage:    "[optional] print prompt to console, default is false",
			Value:    false,
			Required: false,
		},
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
		commands_common.UrlsFlag(),
	}

}

func (c *TLDR) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	showPrompt := cliContext.Bool("showPrompt")
	jsonPrompt := text.Tldr
	finalPrompt, err := prompt.FilePromptBuilder(jsonPrompt, filePaths, urls, gDriveFiles, userPrompt)
	utils.PrintPrompt(showPrompt, jsonPrompt, finalPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
