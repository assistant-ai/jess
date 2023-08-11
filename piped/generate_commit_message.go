package piped

import (
	"fmt"
	jess_cli "github.com/assistant-ai/jess/cli"
	"github.com/assistant-ai/jess/prompt"
	"github.com/assistant-ai/jess/utils"
	"github.com/assistant-ai/llmchat-client/client"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

type ChangedFile struct {
	Path        string `json:"path"`
	Diff        string `json:"diff"`
	jessComment string `json:""`
}

type GenerateCommitJessCommand struct {
	Command DoubleBaseCommand
}

func (c *GenerateCommitJessCommand) DefineCommand(llmClient *client.Client) *cli.Command {
	return &cli.Command{
		Name:   c.Command.Name(),
		Usage:  c.Command.Usage(),
		Action: c.handleActionForCommit(llmClient),
		Flags:  c.Command.Flags(),
	}
}

func (c *GenerateCommitJessCommand) handleActionForCommit(llmClient *client.Client) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		filePathForFinalAnswer := cliContext.String("output_file")
		inputFolder, err := utils.ExpandTilde(cliContext.String("input_folder"))
		if err != nil {
			return err
		}

		_, err = utils.IsValidPath(inputFolder)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		err = utils.IsGitRepository(inputFolder)
		if err != nil {
			return nil
		}

		changedFiles, err := getChangedFilesWithDiffs(inputFolder)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		ChangedFilesMapWithComments, err := iterateChangeFilesAndaddComments(changedFiles, llmClient)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		resultString := joinCommentsToString(ChangedFilesMapWithComments)

		err = utils.AnswersOutput(filePathForFinalAnswer, resultString)
		if err != nil {
			log.Errorf("Error while saving file: %v", err)
			return err
		}

		return nil
	}

}

type GenerateCommitCommand struct{}

func (c *GenerateCommitCommand) Name() string {
	return "gcm"
}

func (c *GenerateCommitCommand) Usage() string {
	return "generate commit message on provided folder. right now it contains some limitation to work with git only. And on a size of new files. "
}

func (c *GenerateCommitCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "input_folder",
			Aliases:  []string{"i"},
			Usage:    "[Mandatory] Path to folder git folder",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output_file",
			Aliases: []string{"o"},
			Usage:   "[Optional] Output file path, by default output will be printed to terminal",
		},
	}
}

func (c *GenerateCommitCommand) PreparePromptForDoubleAction(cliContext *cli.Context) (string, error) {
	filePaths := cliContext.StringSlice("input")
	userPrompt := cliContext.String("prompt")
	urls := cliContext.StringSlice("url")
	gDriveFiles := cliContext.StringSlice("gdrive")
	prePrompt := "Generate concise commit descriptions from the given JSON data containing changed file paths and their respective diffs. Each description should be in a Markdown-friendly format suitable for GitHub commits. The JSON contains an array of objects, each having a \"path\" field representing the file path and a \"diff\" field with the file's changes. Provide a formatted list of commit descriptions corresponding to each file path and its changes."
	finalPrompt, err := prompt.FilePromptBuilder(prePrompt, filePaths, urls, gDriveFiles, userPrompt)
	if err != nil {
		return "", err
	}
	return finalPrompt, nil
}

func joinCommentsToString(ChangedFilesMapWithComments []ChangedFile) string {
	var result string
	for _, file := range ChangedFilesMapWithComments {
		result += fmt.Sprintf("%s:\n%s\n\n", file.Path, file.jessComment)
	}
	return result
}

func getChangedFilesWithDiffs(inputFolder string) ([]ChangedFile, error) {
	var stagedChangedFiles []ChangedFile
	var unStagedChangedFiles []ChangedFile

	if err := os.Chdir(inputFolder); err != nil {
		return nil, err
	}
	stagedCmd := exec.Command("git", "diff", "--name-status", "--ignore-blank-lines")
	StagedOutput, err := stagedCmd.Output()
	if err != nil {
		return nil, err
	}
	stagedChangedFiles = parseDiffOutput(StagedOutput, stagedChangedFiles)
	unStagedCmd := exec.Command("git", "diff", "--staged", "--name-status")
	unStagedOutput, err := unStagedCmd.Output()
	if err != nil {
		return nil, err
	}
	unStagedChangedFiles = parseDiffOutput(unStagedOutput, unStagedChangedFiles)

	for _, uf := range unStagedChangedFiles {
		for _, sf := range stagedChangedFiles {
			if sf.Path != uf.Path {
				unStagedChangedFiles = append(unStagedChangedFiles, sf)
			}
		}
	}
	if len(unStagedChangedFiles) == 0 {
		unStagedChangedFiles = stagedChangedFiles
	}

	return unStagedChangedFiles, nil
}

func parseDiffOutput(output []byte, changedFiles []ChangedFile) []ChangedFile {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			status := parts[0]
			filePath := parts[1]

			switch status {
			case "M":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "")
			case "A":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "Added new file")
			case "D":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "Deleted file")
			case "R":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "Renamed file")
			case "C":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "Copied file")
			case "T":
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "File changed type")
			default:
				changedFiles = setChangeFiledInitialComment(filePath, changedFiles, "Other changes with file")
			}
		}
	}
	return changedFiles
}

func setChangeFiledInitialComment(filePath string, changedFiles []ChangedFile, initialComment string) []ChangedFile {
	diffCmd := exec.Command("git", "diff", filePath)
	diffOutput, _ := diffCmd.Output() // Ignoring error for simplicity
	changedFiles = append(changedFiles, ChangedFile{
		Path:        filePath,
		Diff:        string(diffOutput),
		jessComment: initialComment,
	})
	return changedFiles
}

func iterateChangeFilesAndaddComments(changedFiles []ChangedFile, llmClient *client.Client) ([]ChangedFile, error) {

	prePrompt := "forget about all previos request of same requests. analyze it from start. Generate concise commit descriptions from the given JSON data containing changed file paths and their respective diffs. Each description should be in a Markdown-friendly format suitable for GitHub commits. The JSON contains an array of objects, each having a \"path\" field representing the file path and a \"diff\" field with the file's changes. analyze changes, by methods or functions, and provided explanation of method/functions, not a line. If there were one or more comment per file just add that 'there was some comments added'. Provide a formatted list of commit descriptions corresponding to each file path and its changes. your comments should be as short as possible, but cover sense and main idea of changes. "

	for i := range changedFiles {
		if changedFiles[i].jessComment == "" {
			var err error
			fmt.Printf("           about commit message for : %s", changedFiles[i].Path)
			generatedPrompt := prePrompt + "{\npath:" + changedFiles[i].Path + ", \ndiff:" + changedFiles[i].Diff + "}"
			changedFiles[i].jessComment, err = jess_cli.ExecutePrompt(llmClient, generatedPrompt, "")
			fmt.Println("")
			if err != nil {
				return nil, err
			}
		}
	}
	return changedFiles, nil
}
