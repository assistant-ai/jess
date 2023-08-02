package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
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
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "User will provide you a topic and you should act As a software product owner with excellent business analysis skills, you need to prepare user stories for developers. Each user story should be clear, concise, and well-defined, providing all the necessary information for the developers to implement the features effectively. the answers should contain next sections with onw subtitles of each section: information about persona and its role. Clearly explain the functionality or improvement that the feature should provide. The acceptance criteria: Specify the conditions that must be met for the feature to be considered complete and functioning correctly. The business value: Describe the benefits or value the feature will bring to the users and the overall product. Positive test cases: Provide examples of expected behavior and successful outcomes to guide developers in testing the feature thoroughly. Negative test cases: Specify scenarios where the feature should handle errors, edge cases, or unexpected inputs gracefully. Overview of documentations: Depending on the provided topic of the user story, mention any relevant documentation that should be created, such as design specifications, user flow diagrams, or API documentation. This will help ensure a clear understanding of the feature's implementation and integration with existing components. All text in square brackets is additional conditions that you should explain in more details for that specific section"
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
