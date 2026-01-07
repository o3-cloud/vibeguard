package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunList_Success(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: check-one
    run: "true"
    severity: error
    timeout: 10s
  - id: check-two
    run: "false"
    severity: warning
    timeout: 30s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runList(listCmd, []string{})
	if err != nil {
		t.Errorf("runList failed: %v", err)
	}
}

func TestRunList_WithVerbose(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: verbose-check
    run: "echo hello"
    severity: error
    timeout: 10s
    suggestion: "Run echo hello"
    requires:
      - dep-check
  - id: dep-check
    run: "true"
    severity: warning
    timeout: 5s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runList(listCmd, []string{})
	if err != nil {
		t.Errorf("runList with verbose failed: %v", err)
	}
}

func TestRunList_EmptyChecks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// This is an invalid config (no checks), but Load should handle it
	configContent := `version: "1"
checks:
  - id: minimal
    run: "true"
    severity: error
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runList(listCmd, []string{})
	if err != nil {
		t.Errorf("runList failed: %v", err)
	}
}

func TestRunList_ConfigNotFound(t *testing.T) {
	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = "/nonexistent/path/vibeguard.yaml"

	err := runList(listCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunList_WithDependencies(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: first
    run: "true"
    severity: error
    timeout: 10s
  - id: second
    run: "true"
    severity: error
    timeout: 10s
    requires:
      - first
  - id: third
    run: "true"
    severity: warning
    timeout: 10s
    requires:
      - first
      - second
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runList(listCmd, []string{})
	if err != nil {
		t.Errorf("runList with dependencies failed: %v", err)
	}
}

func TestRunList_WithTags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: fmt
    run: "gofmt -l ."
    severity: error
    timeout: 10s
    tags: [format, fast, pre-commit]
  - id: security
    run: "gosec ./..."
    severity: error
    timeout: 30s
    tags: [security, slow]
  - id: test
    run: "go test ./..."
    severity: error
    timeout: 60s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runList(listCmd, []string{})
	if err != nil {
		t.Errorf("runList with tags failed: %v", err)
	}

	output := buf.String()
	// Verify tags are in output
	if !bytes.Contains(buf.Bytes(), []byte("Tags:")) {
		t.Errorf("expected 'Tags:' in output, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("format")) {
		t.Errorf("expected 'format' tag in output, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("security")) {
		t.Errorf("expected 'security' tag in output, got: %s", output)
	}
}
