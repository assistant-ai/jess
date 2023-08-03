package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/urfave/cli/v2"
)

type PromptGeneratorCommand struct{}

func (c *PromptGeneratorCommand) Name() string {
	return "generate_prompt"
}

func (c *PromptGeneratorCommand) Usage() string {
	return "generate description of user story based on the topic"
}

func (c *PromptGeneratorCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[Optional] Add information about prompt you want to generate. just give summary in one sentence",
			Value:    "",
			Required: false,
		},
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *PromptGeneratorCommand) PreparePrompt(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "I want you to be my prompt maker. Your goal is to help me create the best prompt for my needs. You, ChatGPT, will use this prompt and follow the following process: 1. First of all, you will ask me what the topic of the prompt should be. I will give my answer in the end, but we will have to improve it by constant iteration, going through the following steps. 2. Based on my answer, you will create 3 sections. you will act as an professional of the provided topic and provide me description of your role. THen you will suggest the prompt from the position of the role, you take above. Then you provide me the prompt. The first is the Suggested Prompt. (Submit your version of the rewritten prompt. It should be clear, concise and easy to understand for you). This prompt should include request fro summary of this topic, short SWOT analysis of provided topic. Prompt should also contain a request to provide a list of 5 related areas to this topics with request of explanation why these areas are important, this request should include also request for 5 roles in this areas with explanation why this roles are important, and how they could help to improve your main prompt. Prompt should suggest to split response on s separate sections based on topics and with number and title. Prompt should contain a request to build in the end checklist of all requested topics, to if all topics in prompt were covered by prompt.; Second - Offers (title this section as Suggestion to improve your prompt [copy from square brackets just into your prompt]). (give me at least 5 suggestions on what details should be included in the prompt to improve it, after each suggestion, give me full and extended explanation of suggestion (why it  is really important) and wrap these explanation in square brackets, example should be in format that you will understand if I just copy this answer in the end of prompt). Do not include in offers topics that was requested above. And the third section - Questions. (Ask me at least 5 questions regarding what additional information is required from me to improve the prompt, after each question, give me an full and extended example of answer and wrap these answers in square brackets, example should be in format that you will understand if I just copy this answer in the end of prompt)."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
