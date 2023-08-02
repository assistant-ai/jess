package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

type SolveProblem struct{}

func (c *SolveProblem) Name() string {
	return "solve"
}

func (c *SolveProblem) Usage() string {
	return "Suggest problem solving"
}

func (c *SolveProblem) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[optional] Add information about your problem",
			Value:    "",
			Required: false,
		},
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}

}

func (c *SolveProblem) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "User going to provide you with short test of description of his problem. Proved him a solution with next sections: 1) What is the problem and it roots? Try to explain that it should be understandable even for idiots. 2) What is the solution? 3) Steps of implementing solution 4) What is the benefit of the solution?  5) List of the risks, that could happened if problem won't disappear. 6) Suggest list of steps how to track solving of the problem. 7) Describe success criteria for provided problem and when User could count that problem is fully solved. Describe how to measure success criteria. 7) Describe in details what is call to action for provided problem? Additional instruction: 1) if problem contains multiple participants or side of problem, describe each section for all sides of the problem. User might provide additional requirements."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
