package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAssist_Success(t *testing.T) {
	// Get the project root (this test is in internal/cli)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	projectRoot := filepath.Join(cwd, "..", "..")

	// Save and restore flag state
	oldAssist := initAssist
	oldOutput := initOutput
	defer func() {
		initAssist = oldAssist
		initOutput = oldOutput
	}()

	initAssist = true
	initOutput = ""

	// Create a buffer to capture stdout
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Run the assist command on the project root
	err = runAssist(initCmd, []string{projectRoot})
	if err != nil {
		t.Fatalf("runAssist failed: %v", err)
	}

	// Verify the output contains expected sections
	// Note: runAssist prints to stdout, not to cmd.Out() in our implementation
	// So we just verify it doesn't error
}

func TestRunAssist_NonexistentDir(t *testing.T) {
	oldAssist := initAssist
	oldOutput := initOutput
	defer func() {
		initAssist = oldAssist
		initOutput = oldOutput
	}()

	initAssist = true
	initOutput = ""

	err := runAssist(initCmd, []string{"/nonexistent/path/that/does/not/exist"})
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}

	exitErr, ok := err.(*ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Message, "does not exist") {
		t.Errorf("expected error message to contain 'does not exist', got: %s", exitErr.Message)
	}
}

func TestRunAssist_UndetectableProject(t *testing.T) {
	// Create a temp directory with no project files
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldAssist := initAssist
	oldOutput := initOutput
	defer func() {
		initAssist = oldAssist
		initOutput = oldOutput
	}()

	initAssist = true
	initOutput = ""

	err = runAssist(initCmd, []string{tmpDir})
	if err == nil {
		t.Fatal("expected error for undetectable project")
	}

	exitErr, ok := err.(*ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Message, "unable to detect project type") {
		t.Errorf("expected error message about project type detection, got: %s", exitErr.Message)
	}
}

func TestRunAssist_OutputToFile(t *testing.T) {
	// Get the project root
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	projectRoot := filepath.Join(cwd, "..", "..")

	// Create temp output file
	tmpFile, err := os.CreateTemp("", "vibeguard-prompt-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	oldAssist := initAssist
	oldOutput := initOutput
	defer func() {
		initAssist = oldAssist
		initOutput = oldOutput
	}()

	initAssist = true
	initOutput = tmpFile.Name()

	err = runAssist(initCmd, []string{projectRoot})
	if err != nil {
		t.Fatalf("runAssist failed: %v", err)
	}

	// Verify file was created and contains content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("output file is empty")
	}

	// Verify expected sections
	// After integration with validator_guide.go, the prompt now includes
	// comprehensive validation sections from the assist package
	contentStr := string(content)
	expectedSections := []string{
		"# VibeGuard AI Agent Setup Guide",
		"## Project Analysis",
		"## Recommended Checks",
		"## YAML Syntax Requirements",
		"## Check Structure Requirements",
		"## Your Task",
	}

	for _, section := range expectedSections {
		if !strings.Contains(contentStr, section) {
			t.Errorf("output missing expected section: %s", section)
		}
	}
}

func TestRunAssist_NotADirectory(t *testing.T) {
	// Create a temp file (not a directory)
	tmpFile, err := os.CreateTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	oldAssist := initAssist
	oldOutput := initOutput
	defer func() {
		initAssist = oldAssist
		initOutput = oldOutput
	}()

	initAssist = true
	initOutput = ""

	err = runAssist(initCmd, []string{tmpFile.Name()})
	if err == nil {
		t.Fatal("expected error for non-directory")
	}

	exitErr, ok := err.(*ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Message, "not a directory") {
		t.Errorf("expected error message about not being a directory, got: %s", exitErr.Message)
	}
}
