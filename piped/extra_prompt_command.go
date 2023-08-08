package piped

import (
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
)

type DoubleBaseCommand interface {
	Flags() []cli.Flag
	PreparePromptForDoubleAction(cliContext *cli.Context) (string, error)
	Name() string
	Usage() string
}

type DoubleJessCommand struct {
	Command DoubleBaseCommand
}

func (c *DoubleJessCommand) DefineCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.handleDoubleAction(llmClient),
		Flags:  c.Command.Flags(),
	}
}

func (c *DoubleJessCommand) handleDoubleAction(llmClient *client.Client) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		context := cliContext.String("context")
		userFinalPrompt := cliContext.String("prompt")
		filePathForFinalAnswer := cliContext.String("output")
		filePathForPromptOutput := cliContext.String("output_prompt")

		initialPrompt, err := c.Command.PreparePromptForDoubleAction(cliContext)

		if err != nil {
			log.Errorf("Error while sending message: %v", err)
			return err
		}

		utils.PrintlnCyan("USER PROMPT:\n" + userFinalPrompt + "\n\n")
		generatedPrompt, err := jess_cli.ExecutePrompt(llmClient, initialPrompt, context)

		err = utils.AnswersOutput(filePathForPromptOutput, generatedPrompt)

		if err != nil {
			log.Errorf("Error while saving file: %v", err)
			return err
		}

		answer, err := jess_cli.ExecutePrompt(llmClient, generatedPrompt, context)

		err = utils.AnswersOutput(filePathForFinalAnswer, answer)
		if err != nil {
			log.Errorf("Error while saving file: %v", err)
			return err
		}

		return nil
	}

}

type DoublePromptCommand struct{}

func (c *DoublePromptCommand) Name() string {
	return "dp"
}

func (c *DoublePromptCommand) Usage() string {
	return "Running prompt in double mode. Means you give a topic. Jess generate basic prompt and generate 2 files - one with generated prompt for future editint and second with the result of the prompt."
}

func (c *DoublePromptCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "prompt",
			Aliases:  []string{"p"},
			Usage:    "[Optional] Add information about prompt you want to generate. just give summary in one sentence",
			Value:    "",
			Required: false,
		},
		&cli.StringFlag{
			Name:    "output_prompt",
			Aliases: []string{"op"},
			Usage:   "[Optional] Output file path for prompt, by default output will be printed to terminal",
		},
		commands_common.InputFilesFlag(),
		commands_common.ContextFlag(),
		commands_common.OutputFlag(),
		commands_common.GoogleDriveFilesFlag(),
	}
}

func (c *DoublePromptCommand) PreparePromptForDoubleAction(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "I want you to be my the professional in making prompts for chatgpt. I will provide you with a topic and you will create a prompt for me.\nyour answer should start from asking chatbot to imagine him and act as person who are most experienced in requested topic.\nThen returned request from the position of the role, you take above with requesting to explain the topic with following additional tasks:\nThis prompt should include request for summary of this topic.\nA request for providing top important definitions in requested topics, as a table.\nA request for short SWOT analysis of provided topic.\nPrompt should also contain a request to provide a list of 5 related areas to this topics with request of explanation why these areas are important.\nPrompt should also contain a request for 5 roles in this areas with explanation why this roles are important, and how they could help to improve your main prompt.For each role, include a short description of the role and how it relates to the topic. For each role, give an example (in one sentence) how would they create prompts for this topic, from their perspective.\nA Request for a typical goals in requested topic. Request of providing such kind of goals. Request for providing these examples as SMART goals. A Request for output this list as a table\nPrompt should suggest to split response on s separate sections based on topics and with number and title.\nAll previous requests should be contain specific knowledge from the perspective of the role you take.\nPrompt should contain a request to build in the end checklist of all requested topics, to if all topics in requested were covered.\nPrompt should also request to provide answer in markdown format.\nYou need to return only suggested request for the topic, no any comments from your side are required.\n"
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}
