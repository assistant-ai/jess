package commands

import "github.com/urfave/cli/v2"

func PromptFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "prompt",
		Aliases:  []string{"p"},
		Usage:    "prompt to suppy with file",
		Required: true,
	}
}

func InputFilesFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:     "input",
		Aliases:  []string{"i"},
		Usage:    "input files",
		Required: true,
	}
}

func InputFileFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "input",
		Aliases:  []string{"i"},
		Usage:    "input file",
		Required: true,
	}
}

func ContextFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "context",
		Aliases: []string{"c"},
		Usage:   "context id to store this to",
	}
}

func OutputFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "output file path, if not specificed stdout will be used",
	}
}
