package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunValidate_ValidConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: valid-check
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

	err = runValidate(validateCmd, []string{})
	if err != nil {
		t.Errorf("runValidate failed for valid config: %v", err)
	}
}

func TestRunValidate_WithVerbose(t *testing.T) {
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
    severity: warning
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
	defer func() {
		configFile = oldConfig
		verbose = oldVerbose
	}()

	configFile = configPath
	verbose = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runValidate(validateCmd, []string{})
	if err != nil {
		t.Errorf("runValidate with verbose failed: %v", err)
	}
}

func TestRunValidate_ConfigNotFound(t *testing.T) {
	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = "/nonexistent/path/vibeguard.yaml"

	err := runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunValidate_InvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Invalid YAML syntax
	configContent := `version: "1"
checks:
  - id: bad-check
    run: "true"
    severity: error
    timeout: 10s
  - this is not valid yaml: [
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected 'validation failed' in error, got: %v", err)
	}
}

func TestRunValidate_MissingCheckID(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Check without ID
	configContent := `version: "1"
checks:
  - run: "true"
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing check ID")
	}
}

func TestRunValidate_InvalidSeverity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: bad-severity
    run: "true"
    severity: critical
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestRunValidate_CyclicDependency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Cyclic dependency: a -> b -> a
	configContent := `version: "1"
checks:
  - id: check-a
    run: "true"
    severity: error
    timeout: 10s
    requires:
      - check-b
  - id: check-b
    run: "true"
    severity: error
    timeout: 10s
    requires:
      - check-a
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}
}

func TestRunValidate_UnknownRequires(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: check-a
    run: "true"
    severity: error
    timeout: 10s
    requires:
      - nonexistent-check
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for unknown requires")
	}
}

func TestRunValidate_DuplicateCheckID(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: duplicate
    run: "true"
    severity: error
    timeout: 10s
  - id: duplicate
    run: "false"
    severity: warning
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

	err = runValidate(validateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for duplicate check ID")
	}
}
