package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrompt_ListAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.
    tags: [setup, initialization, guidance]
  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer.
    tags: [review, quality]
  - id: security-audit
    description: "Security-focused code analysis"
    content: |
      You are a security auditor.
    tags: [security, audit]
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("runPrompt failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Prompts (3):") {
		t.Errorf("expected 'Prompts (3):' in output, got: %s", output)
	}
	if !strings.Contains(output, "init") {
		t.Errorf("expected 'init' in output, got: %s", output)
	}
	if !strings.Contains(output, "code-review") {
		t.Errorf("expected 'code-review' in output, got: %s", output)
	}
	if !strings.Contains(output, "security-audit") {
		t.Errorf("expected 'security-audit' in output, got: %s", output)
	}
}

func TestRunPrompt_ListAllVerbose(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.
    tags: [setup, initialization, guidance]
  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer.
    tags: [review, quality, go]
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("runPrompt with verbose failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Prompts (2):") {
		t.Errorf("expected 'Prompts (2):' in output, got: %s", output)
	}
	if !strings.Contains(output, "Description:") {
		t.Errorf("expected 'Description:' in output, got: %s", output)
	}
	if !strings.Contains(output, "Guidance for initializing vibeguard configuration") {
		t.Errorf("expected description text in output, got: %s", output)
	}
	if !strings.Contains(output, "Tags:") {
		t.Errorf("expected 'Tags:' in output, got: %s", output)
	}
	if !strings.Contains(output, "setup") {
		t.Errorf("expected 'setup' tag in output, got: %s", output)
	}
}

func TestRunPrompt_ReadSpecific(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.

      Guide them through:
      1. Detecting their project type
      2. Recommending checks
      3. Creating configuration
    tags: [setup, initialization, guidance]
  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer.
    tags: [review, quality]
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{"init"})
	if err != nil {
		t.Errorf("runPrompt for specific prompt failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "You are an expert in helping users set up VibeGuard.") {
		t.Errorf("expected prompt content in output, got: %s", output)
	}
}

func TestRunPrompt_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.
    tags: [setup, initialization, guidance]
checks:
  - id: test
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

	err = runPrompt(promptCmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for non-existent prompt")
	}
	if !strings.Contains(err.Error(), "prompt not found") {
		t.Errorf("expected 'prompt not found' in error, got: %v", err)
	}
}

func TestRunPrompt_NoPrompts(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
checks:
  - id: test
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
	promptCmd.SetOut(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Fatalf("expected no error when listing prompts (built-in init prompt available), got: %v", err)
	}

	output := buf.String()
	// Should show built-in init prompt
	if !strings.Contains(output, "init") {
		t.Errorf("expected 'init' in output, got: %q", output)
	}
	if !strings.Contains(output, "built-in") {
		t.Errorf("expected 'built-in' in output, got: %q", output)
	}
}

func TestRunPrompt_JSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: init
    description: "Guidance for initializing vibeguard configuration"
    content: |
      You are an expert in helping users set up VibeGuard.
    tags: [setup, initialization, guidance]
  - id: code-review
    description: "System prompt for code review assistance"
    content: |
      You are an expert code reviewer.
    tags: [review, quality]
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("runPrompt with JSON output failed: %v", err)
	}

	output := buf.String()
	var result []map[string]interface{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		t.Errorf("failed to parse JSON output: %v, output: %s", err, output)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 prompts in JSON (2 from config + 1 built-in), got %d", len(result))
	}

	foundInit := false
	for _, prompt := range result {
		if id, ok := prompt["id"].(string); ok && id == "init" {
			foundInit = true
			if desc, ok := prompt["description"].(string); !ok || desc == "" {
				t.Errorf("expected description in init prompt JSON")
			}
			if tags, ok := prompt["tags"].([]interface{}); !ok || len(tags) == 0 {
				t.Errorf("expected tags in init prompt JSON")
			}
			break
		}
	}
	if !foundInit {
		t.Errorf("expected 'init' prompt in JSON output")
	}
}

func TestRunPrompt_ConfigNotFound(t *testing.T) {
	oldConfig := configFile
	defer func() {
		configFile = oldConfig
	}()

	configFile = "/nonexistent/path/vibeguard.yaml"

	err := runPrompt(promptCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunPrompt_PromptWithoutDescription(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: simple
    content: |
      Simple prompt content
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("runPrompt with prompt lacking description failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "simple") {
		t.Errorf("expected 'simple' in output, got: %s", output)
	}
}

func TestRunPrompt_MultilinePromptContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configContent := `version: "1"
prompts:
  - id: multiline
    description: "A prompt with multiple lines"
    content: |
      Line 1
      Line 2
      Line 3
      Multiple lines of content
      With formatting preserved
checks:
  - id: test
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
	promptCmd.SetOut(&buf)
	promptCmd.SetErr(&buf)

	err = runPrompt(promptCmd, []string{"multiline"})
	if err != nil {
		t.Errorf("runPrompt for multiline prompt failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Line 1") {
		t.Errorf("expected 'Line 1' in output")
	}
	if !strings.Contains(output, "Multiple lines of content") {
		t.Errorf("expected 'Multiple lines of content' in output")
	}
}

func TestRunPrompt_BuiltinInitPrompt_Read(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a config without an init prompt
	configContent := `version: "1"
prompts:
  - id: code-review
    content: "You are a code reviewer"
checks:
  - id: test
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
	promptCmd.SetOut(&buf)

	// Request built-in init prompt
	err = runPrompt(promptCmd, []string{"init"})
	if err != nil {
		t.Errorf("reading built-in init prompt failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "VibeGuard") {
		t.Errorf("expected 'VibeGuard' in built-in init prompt, got: %q", output)
	}
	if !strings.Contains(output, "policy enforcement") {
		t.Errorf("expected 'policy enforcement' in built-in init prompt, got: %q", output)
	}
}

func TestRunPrompt_BuiltinInitPrompt_InList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a config with some prompts
	configContent := `version: "1"
prompts:
  - id: code-review
    description: "Code review prompt"
    content: "You are a code reviewer"
checks:
  - id: test
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
	promptCmd.SetOut(&buf)

	// List prompts - should include built-in init
	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("listing prompts with built-in init failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "code-review") {
		t.Errorf("expected 'code-review' in output")
	}
	if !strings.Contains(output, "init") {
		t.Errorf("expected 'init' (built-in) in output")
	}
	if !strings.Contains(output, "built-in") {
		t.Errorf("expected '(built-in)' indicator in output")
	}
}

func TestRunPrompt_BuiltinInitPrompt_JSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibeguard-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Config with no prompts - only built-in init will be available
	configContent := `version: "1"
checks:
  - id: test
    run: "true"
    severity: error
    timeout: 10s
`
	configPath := filepath.Join(tmpDir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfig := configFile
	oldJSON := jsonOutput
	defer func() {
		configFile = oldConfig
		jsonOutput = oldJSON
	}()

	configFile = configPath
	jsonOutput = true

	var buf bytes.Buffer
	promptCmd.SetOut(&buf)

	err = runPrompt(promptCmd, []string{})
	if err != nil {
		t.Errorf("listing prompts in JSON failed: %v", err)
	}

	output := buf.String()
	var result []map[string]interface{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		t.Errorf("failed to parse JSON output: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 built-in prompt, got %d", len(result))
	}

	if result[0]["id"] != "init" {
		t.Errorf("expected 'init' prompt ID, got %q", result[0]["id"])
	}

	if builtIn, ok := result[0]["built_in"].(bool); !ok || !builtIn {
		t.Errorf("expected built_in: true in JSON output")
	}
}
