package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunTags_WithTags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: fmt
    run: "true"
    severity: error
    timeout: 10s
    tags: [format, fast, pre-commit]
  - id: lint
    run: "true"
    severity: error
    timeout: 10s
    tags: [lint, fast]
  - id: security
    run: "true"
    severity: error
    timeout: 10s
    tags: [security, slow]
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

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runTags(tagsCmd, []string{})
	if err != nil {
		t.Errorf("runTags failed: %v", err)
	}

	output := buf.String()
	expectedTags := []string{"fast", "format", "lint", "pre-commit", "security", "slow"}
	outputLines := strings.TrimSpace(output)
	outputTags := strings.Split(outputLines, "\n")

	if len(outputTags) != len(expectedTags) {
		t.Errorf("expected %d tags, got %d", len(expectedTags), len(outputTags))
	}

	for i, expected := range expectedTags {
		if i < len(outputTags) && outputTags[i] != expected {
			t.Errorf("tag at index %d: expected %q, got %q", i, expected, outputTags[i])
		}
	}
}

func TestRunTags_EmptyConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: fmt
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

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runTags(tagsCmd, []string{})
	if err != nil {
		t.Errorf("runTags with no tags failed: %v", err)
	}

	output := buf.String()
	if strings.TrimSpace(output) != "" {
		t.Errorf("expected no output for config with no tags, got: %q", output)
	}
}

func TestRunTags_DuplicateTags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: check1
    run: "true"
    severity: error
    timeout: 10s
    tags: [fast, lint]
  - id: check2
    run: "true"
    severity: error
    timeout: 10s
    tags: [fast, lint, security]
  - id: check3
    run: "true"
    severity: error
    timeout: 10s
    tags: [fast]
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

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runTags(tagsCmd, []string{})
	if err != nil {
		t.Errorf("runTags with duplicate tags failed: %v", err)
	}

	output := buf.String()
	outputLines := strings.TrimSpace(output)
	outputTags := strings.Split(outputLines, "\n")

	// Should have 3 unique tags: fast, lint, security (in sorted order)
	expectedTags := []string{"fast", "lint", "security"}
	if len(outputTags) != len(expectedTags) {
		t.Errorf("expected %d unique tags, got %d", len(expectedTags), len(outputTags))
	}

	for i, expected := range expectedTags {
		if i < len(outputTags) && outputTags[i] != expected {
			t.Errorf("tag at index %d: expected %q, got %q", i, expected, outputTags[i])
		}
	}
}

func TestRunTags_Sorted(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: check1
    run: "true"
    severity: error
    timeout: 10s
    tags: [zebra, alpha, beta]
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

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err = runTags(tagsCmd, []string{})
	if err != nil {
		t.Errorf("runTags with unsorted tags failed: %v", err)
	}

	output := buf.String()
	outputLines := strings.TrimSpace(output)
	outputTags := strings.Split(outputLines, "\n")

	// Should be sorted alphabetically
	expectedTags := []string{"alpha", "beta", "zebra"}
	if len(outputTags) != len(expectedTags) {
		t.Errorf("expected %d tags, got %d", len(expectedTags), len(outputTags))
	}

	for i, expected := range expectedTags {
		if i < len(outputTags) && outputTags[i] != expected {
			t.Errorf("tag at index %d: expected %q, got %q", i, expected, outputTags[i])
		}
	}
}

func TestRunTags_ConfigNotFound(t *testing.T) {
	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = "/nonexistent/path/vibeguard.yaml"

	err := runTags(tagsCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}
