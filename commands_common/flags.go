package commands_common

import "github.com/urfave/cli/v2"

func PromptFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "prompt",
		Aliases: []string{"p"},
		Usage:   "[optional] prompt to supply with file",
		Value:   "",
	}
}

func InputFilesFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "input",
		Aliases: []string{"i"},
		Usage:   "input files",
	}
}

func UrlsFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "url",
		Aliases: []string{"u"},
		Usage:   "URL to download, extract text and feed to the GPT model",
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
		Usage:   "Context id to store this to",
	}
}

func OutputFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Output file path, by default output will be printed to terminal",
	}
}

func GoogleDriveFilesFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "gdrive",
		Aliases: []string{"g", "gd"},
		Usage:   "Google Drive file ID",
	}
}
