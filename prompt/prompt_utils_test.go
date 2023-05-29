package prompt

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFilePromptBuilderSingleFile(t *testing.T) {
	// Creating a temporary file for testing
	tmpFile := "test_file.txt"
	testContent := "This is a test file."
	err := ioutil.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary test file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Testing FilePromptBuilder
	prePrompt := "PrePrompt Test"
	userPrompt := "UserPrompt Test"
	filePaths := []string{tmpFile}
	urls := []string{}
	googleDriveFiles := []string{}

	result, err := FilePromptBuilder(prePrompt, filePaths, urls, googleDriveFiles, userPrompt)
	if err != nil {
		t.Fatalf("Failed to execute FilePromptBuilder: %v", err)
	}

	// Assertions
	if !strings.Contains(result, prePrompt) {
		t.Errorf("FilePromptBuilder output does not contain prePrompt. Result: %v", result)
	}

	if !strings.Contains(result, userPrompt) {
		t.Errorf("FilePromptBuilder output does not contain userPrompt. Result: %v", result)
	}

	if !strings.Contains(result, tmpFile) {
		t.Errorf("FilePromptBuilder output does not contain file path. Result: %v", result)
	}

	if !strings.Contains(result, testContent) {
		t.Errorf("FilePromptBuilder output does not contain file content. Result: %v", result)
	}
}
