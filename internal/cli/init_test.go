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
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

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

func TestRunInit_CreateDefault(t *testing.T) {
	// Create a temp directory for the test
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Change to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	// Save and restore flag state
	oldForce := initForce
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initForce = oldForce
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initForce = false
	initTemplate = ""
	initAssist = false

	// Run init command
	err = runInit(initCmd, []string{})
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify the file was created
	content, err := os.ReadFile("vibeguard.yaml")
	if err != nil {
		t.Fatalf("failed to read created config: %v", err)
	}

	// Verify it contains expected content
	if !strings.Contains(string(content), "version:") {
		t.Error("created config missing version field")
	}
	if !strings.Contains(string(content), "checks:") {
		t.Error("created config missing checks field")
	}
}

func TestRunInit_WithTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	oldForce := initForce
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initForce = oldForce
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initForce = false
	initTemplate = "go-minimal"
	initAssist = false

	err = runInit(initCmd, []string{})
	if err != nil {
		t.Fatalf("runInit with template failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat("vibeguard.yaml"); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}

func TestRunInit_UnknownTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	oldForce := initForce
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initForce = oldForce
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initForce = false
	initTemplate = "nonexistent-template"
	initAssist = false

	err = runInit(initCmd, []string{})
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
	if !strings.Contains(err.Error(), "unknown template") {
		t.Errorf("expected 'unknown template' in error message, got: %v", err)
	}
}

func TestRunInit_AlreadyExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create an existing config
	existingConfig := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(existingConfig, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create existing config: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	oldForce := initForce
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initForce = oldForce
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initForce = false
	initTemplate = ""
	initAssist = false

	err = runInit(initCmd, []string{})
	if err == nil {
		t.Fatal("expected error when config already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' in error message, got: %v", err)
	}
}

func TestRunInit_ForceOverwrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create an existing config
	existingConfig := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(existingConfig, []byte("old content"), 0644); err != nil {
		t.Fatalf("failed to create existing config: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	oldForce := initForce
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initForce = oldForce
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initForce = true
	initTemplate = ""
	initAssist = false

	err = runInit(initCmd, []string{})
	if err != nil {
		t.Fatalf("runInit with --force failed: %v", err)
	}

	// Verify content was overwritten
	content, err := os.ReadFile(existingConfig)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if string(content) == "old content" {
		t.Error("config was not overwritten")
	}
}

func TestListTemplates(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := listTemplates()

	_ = w.Close()
	os.Stdout = oldStdout
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Errorf("listTemplates failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Available templates:") {
		t.Errorf("output missing 'Available templates:' header")
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("output missing usage information")
	}
}

func TestRunInit_ListTemplates(t *testing.T) {
	oldTemplate := initTemplate
	oldAssist := initAssist
	defer func() {
		initTemplate = oldTemplate
		initAssist = oldAssist
	}()

	initTemplate = "list"
	initAssist = false

	// This should call listTemplates() and return nil
	err := runInit(initCmd, []string{})
	if err != nil {
		t.Errorf("runInit with --template list failed: %v", err)
	}
}

func TestRunInit_ListTemplatesFlag(t *testing.T) {
	oldListTemplates := initListTemplates
	oldAssist := initAssist
	defer func() {
		initListTemplates = oldListTemplates
		initAssist = oldAssist
	}()

	initListTemplates = true
	initAssist = false

	// This should call listTemplates() and return nil
	err := runInit(initCmd, []string{})
	if err != nil {
		t.Errorf("runInit with --list-templates flag failed: %v", err)
	}
}

func TestListTemplates_OutputFormat(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := listTemplates()

	_ = w.Close()
	os.Stdout = oldStdout
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("listTemplates failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Verify output structure
	if len(lines) < 3 {
		t.Errorf("output too short, expected at least 3 lines")
	}

	// First line should be header
	if !strings.Contains(lines[0], "Available templates:") {
		t.Errorf("first line should contain 'Available templates:', got: %s", lines[0])
	}

	// Should have template entries with name and description
	hasTemplateEntry := false
	for _, line := range lines {
		if strings.Contains(line, "go-") {
			hasTemplateEntry = true
			break
		}
	}
	if !hasTemplateEntry {
		t.Errorf("output should contain at least one template entry starting with 'go-'")
	}

	// Should have usage line
	usageFound := false
	for _, line := range lines {
		if strings.Contains(line, "Usage:") {
			usageFound = true
			break
		}
	}
	if !usageFound {
		t.Errorf("output missing usage line")
	}
}

func TestRunAssist_NotADirectory(t *testing.T) {
	// Create a temp file (not a directory)
	tmpFile, err := os.CreateTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpFile.Name()) }()

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
