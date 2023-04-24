package cli

import (
	"github.com/assistant-ai/jess/model"
	"github.com/urfave/cli/v2"
)

func DefineQuestionCommand(ctx *model.AppContext) *cli.Command {
	return &cli.Command{
		Name:   "question",
		Usage:  "Ask questions about the code",
		Action: HandleProcessAction(ctx),
		Flags:  ProcessFlags(),
	}
}
