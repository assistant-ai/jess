package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

type GrammarCommand struct{}

func (c *GrammarCommand) Name() string {
	return "grammar"
}

func (c *GrammarCommand) Usage() string {
	return "Improve spelling and grammar"
}

func (c *GrammarCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFileFlag(),
		commands_common.ContextFlag(),
		commands_common.PromptFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *GrammarCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePath := cliContext.String("input")
	userPrompt := cliContext.String("prompt")
	filePaths := []string{filePath}
	urls := []string{}
	gDriveFiles := []string{}
	prePrompt := "User going to provide you with text as well as some context to it. Figure out which text user wants to update yourself. You should fix all misspelling and fix all grammar issues in this text, if it requires you could change some phrasal verbs and phrases that text should sound more clear."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
