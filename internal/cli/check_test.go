package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunCheck_Success(t *testing.T) {
	// Create a temp directory with a valid config
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a simple config with a passing check
	configContent := `version: "1"
checks:
  - id: pass
    run: "true"
    severity: error
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Save and restore flag state
	oldConfig := configFile
	oldVerbose := verbose
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = false

	// Create a buffer to capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Run check command
	err = runCheck(checkCmd, []string{})
	if err != nil {
		t.Errorf("runCheck failed: %v", err)
	}
}

func TestRunCheck_SingleCheck(t *testing.T) {
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
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Run single check by ID
	err = runCheck(checkCmd, []string{"check-one"})
	if err != nil {
		t.Errorf("runCheck with single check failed: %v", err)
	}
}

func TestRunCheck_Failing(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a config with a failing check
	configContent := `version: "1"
checks:
  - id: fail
    run: "false"
    severity: error
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runCheck(checkCmd, []string{})
	if err == nil {
		t.Fatal("expected error for failing check")
	}

	exitErr, ok := err.(*ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	// Exit code 3 is ExitCodeViolation for error-severity check failures
	if exitErr.Code != 3 {
		t.Errorf("expected exit code 3 (violation), got %d", exitErr.Code)
	}
}

func TestRunCheck_WithVerbose(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: verbose-check
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
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = true
	jsonOutput = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runCheck(checkCmd, []string{})
	if err != nil {
		t.Errorf("runCheck with verbose failed: %v", err)
	}
}

func TestRunCheck_WithJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: json-check
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
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runCheck(checkCmd, []string{})
	if err != nil {
		t.Errorf("runCheck with JSON output failed: %v", err)
	}
}

func TestRunCheck_ConfigNotFound(t *testing.T) {
	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = "/nonexistent/path/vibeguard.yaml"

	err := runCheck(checkCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunCheck_UnknownCheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: real-check
    run: "true"
    severity: error
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = configPath

	// Run with a non-existent check ID
	err = runCheck(checkCmd, []string{"nonexistent-check"})
	if err == nil {
		t.Fatal("expected error for unknown check")
	}
}

func TestRunCheck_WithDependencies(t *testing.T) {
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
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runCheck(checkCmd, []string{})
	if err != nil {
		t.Errorf("runCheck with dependencies failed: %v", err)
	}
}

func TestRunCheck_Warning(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// A warning severity check that fails returns exit code 0
	// (warnings don't cause overall failure)
	configContent := `version: "1"
checks:
  - id: warn
    run: "false"
    severity: warning
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldVerbose := verbose
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
		jsonOutput = oldJSON
	}()

	configFile = configPath
	verbose = false
	jsonOutput = false

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runCheck(checkCmd, []string{})
	// Warnings should not cause an error
	if err != nil {
		t.Errorf("runCheck with warning severity should not fail: %v", err)
	}
}
