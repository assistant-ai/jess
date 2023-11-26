package commands_qa

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineCVCommand(llmClient *client.Client) *cli.Command {
	qaRecommendation := commands_common.JessCommand{
		Command: &QaSwaggerInitCommand{},
	}
	qaCheckList := commands_common.JessCommand{
		Command: &QaCheckListCommand{},
	}
	qaExtensiveList := commands_common.JessCommand{
		Command: &QaExtensiveListCommand{},
	}
	qaSwaggerCheckList := commands_common.JessCommand{
		Command: &QaSwaggerInitCommand{},
	}
	qaSwaggerPython := commands_common.JessCommand{
		Command: &QaSwaggerPythonCommand{},
	}
	qaSwaggerCurl := commands_common.JessCommand{
		Command: &QASwaggerCurlCommand{},
	}
	qaListToCases := commands_common.JessCommand{
		Command: &QaListToCasesCommand{},
	}
	qaDetailedTestCases := commands_common.JessActionCommand{
		Command: &QADetailedTestVasesCommand{},
	}

	return &cli.Command{
		Name:  "qa",
		Usage: "Actions that would help qa engineers in their work",
		Subcommands: []*cli.Command{
			qaRecommendation.DefineCommand(llmClient),
			qaCheckList.DefineCommand(llmClient),
			qaExtensiveList.DefineCommand(llmClient),
			qaSwaggerCheckList.DefineCommand(llmClient),
			qaSwaggerPython.DefineCommand(llmClient),
			qaSwaggerCurl.DefineCommand(llmClient),
			qaListToCases.DefineCommand(llmClient),
			qaDetailedTestCases.DefineCommand(llmClient),
		},
	}
}
