package commands_code

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

type RefactorCommand struct{}

func (c *RefactorCommand) Name() string {
	return "refactor"
}

func (c *RefactorCommand) Usage() string {
	return "Refactor file for me"
}

func (c *RefactorCommand) Flags() []cli.Flag {
	return []cli.Flag{
		commands_common.InputFileFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
	}
}

func (c *RefactorCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePath := cliContext.String("input")
	filePaths := []string{filePath}
	urls := []string{}
	gDriveFiles := []string{}
	prePrompt := "Let me show you code file."
	userPrompt := "Please refactor it using the best practices of the language that is used. Ideally final version should be as readable as possible. You can refactor any way you want as long as public APIs are no changed. Feel free to extract anything that is needed to a helper function/methods entities. Add comments with describing major chunks of code. Add explanation as comments why it was optimized. Your final output should be full content of the file end-to-end without any text before/after, I will use your output and override my original file."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
