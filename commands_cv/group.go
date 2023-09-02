package commands_cv

import (
	"github.com/assistant-ai/jess/commands_common"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/urfave/cli/v2"
)

func DefineCVCommand(llmClient *client.Client) *cli.Command {
	cvRequirements := commands_common.JessCommand{
		Command: &CvRequirementsCommand{},
	}
	cvRecommendations := commands_common.JessActionCommand{
		Command: &CvRecommendationCommand{},
	}
	return &cli.Command{
		Name:  "cv",
		Usage: "Actions that should helps to works with CV",
		Subcommands: []*cli.Command{
			cvRequirements.DefineCommand(llmClient),
			cvRecommendations.DefineCommand(llmClient),
		},
	}
}
