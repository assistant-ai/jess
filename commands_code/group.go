package commands_code

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/gpt"
	"github.com/urfave/cli/v2"
)

func DefineCodeCommand(gpt *gpt.GptClient) *cli.Command {
	questionCommand := commands_common.JessCommand{
		Command: &QuestionCommand{},
	}
	explainCommand := commands_common.JessCommand{
		Command: &ExplainCommand{},
	}
	refactorCommand := commands_common.JessCommand{
		Command: &RefactorCommand{},
	}
	return &cli.Command{
		Name:  "code",
		Usage: "Actions to take with code",
		Subcommands: []*cli.Command{
			questionCommand.DefineCommand(gpt),
			explainCommand.DefineCommand(gpt),
			refactorCommand.DefineCommand(gpt),
		},
	}
}
