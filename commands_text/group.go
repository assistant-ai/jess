package commands_text

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineTextCommand(llmClient *client.Client) *cli.Command {
	mailCommand := commands_common.JessCommand{
		Command: &MailCommand{},
	}
	grammarCommand := commands_common.JessCommand{
		Command: &GrammarCommand{},
	}
	solveProblemCommand := commands_common.JessCommand{
		Command: &SolveProblem{},
	}
	userStoryCommand := commands_common.JessCommand{
		Command: &UserStoryCommand{},
	}
	generatePromptCommand := commands_common.JessCommand{
		Command: &PromptGeneratorCommand{},
	}
	generatePromptCommandjson := commands_common.JessCommand{
		Command: &PromptGeneratorCommandJson{},
	}
	tldrCommand := commands_common.JessCommand{
		Command: &TLDR{},
	}
	return &cli.Command{
		Name:  "text",
		Usage: "Actions to take with text",
		Subcommands: []*cli.Command{
			mailCommand.DefineCommand(llmClient),
			grammarCommand.DefineCommand(llmClient),
			solveProblemCommand.DefineCommand(llmClient),
			userStoryCommand.DefineCommand(llmClient),
			generatePromptCommand.DefineCommand(llmClient),
			generatePromptCommandjson.DefineCommand(llmClient),
			tldrCommand.DefineCommand(llmClient),
		},
	}
}
