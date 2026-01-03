package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Create a temporary config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
vars:
  packages: "./..."
checks:
  - id: test
    run: go test {{.packages}}
    severity: error
    timeout: 60s
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Version != "1" {
		t.Errorf("expected version '1', got: %s", cfg.Version)
	}

	if len(cfg.Checks) != 1 {
		t.Fatalf("expected 1 check, got: %d", len(cfg.Checks))
	}

	check := cfg.Checks[0]
	if check.ID != "test" {
		t.Errorf("expected check id 'test', got: %s", check.ID)
	}

	// Verify interpolation happened
	if check.Run != "go test ./..." {
		t.Errorf("expected interpolated run command, got: %s", check.Run)
	}

	if check.Timeout != Duration(60*time.Second) {
		t.Errorf("expected timeout 60s, got: %v", check.Timeout)
	}
}

func TestLoad_MinimalConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// Minimal valid config
	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	check := cfg.Checks[0]
	// Check defaults were applied
	if check.Severity != SeverityError {
		t.Errorf("expected default severity 'error', got: %s", check.Severity)
	}

	if check.Timeout != Duration(DefaultTimeout) {
		t.Errorf("expected default timeout %v, got: %v", DefaultTimeout, check.Timeout)
	}
}

func TestLoad_DefaultVersion(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// Config without explicit version
	content := `
checks:
  - id: test
    run: go test ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Version != "1" {
		t.Errorf("expected default version '1', got: %s", cfg.Version)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/vibeguard.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: [invalid yaml structure here
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_UnsupportedVersion(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "2"
checks:
  - id: test
    run: go test ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestLoad_NoChecks(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks: []
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for no checks")
	}
}

func TestLoad_CheckMissingID(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - run: go test ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for check without id")
	}
}

func TestLoad_CheckMissingRun(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for check without run")
	}
}

func TestLoad_DuplicateCheckID(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
  - id: test
    run: go vet ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for duplicate check id")
	}
}

func TestLoad_InvalidSeverity(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    severity: critical
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestLoad_UnknownRequires(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    requires:
      - nonexistent
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for unknown requires reference")
	}
}

func TestLoad_ForwardRequires(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// test requires build which is defined after - this should be valid
	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    requires:
      - build
  - id: build
    run: go build ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected forward reference to be valid, got: %v", err)
	}

	if len(cfg.Checks) != 2 {
		t.Errorf("expected 2 checks, got: %d", len(cfg.Checks))
	}
}

func TestLoad_SelfRequires(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    requires:
      - test
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for self-referencing requires")
	}
}

func TestLoad_CyclicDependency_TwoNodes(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// A requires B, B requires A
	content := `
version: "1"
checks:
  - id: a
    run: echo a
    requires:
      - b
  - id: b
    run: echo b
    requires:
      - a
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}

	// Check that the error message mentions the cycle
	errMsg := err.Error()
	if !strings.Contains(errMsg, "cyclic dependency") {
		t.Errorf("expected error to mention 'cyclic dependency', got: %s", errMsg)
	}
}

func TestLoad_CyclicDependency_ThreeNodes(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// A requires B, B requires C, C requires A
	content := `
version: "1"
checks:
  - id: a
    run: echo a
    requires:
      - b
  - id: b
    run: echo b
    requires:
      - c
  - id: c
    run: echo c
    requires:
      - a
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "cyclic dependency") {
		t.Errorf("expected error to mention 'cyclic dependency', got: %s", errMsg)
	}
}

func TestLoad_CyclicDependency_PartialCycle(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// A requires B, B requires C, C requires B (cycle is B->C->B, not involving A)
	content := `
version: "1"
checks:
  - id: a
    run: echo a
    requires:
      - b
  - id: b
    run: echo b
    requires:
      - c
  - id: c
    run: echo c
    requires:
      - b
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "cyclic dependency") {
		t.Errorf("expected error to mention 'cyclic dependency', got: %s", errMsg)
	}
}

func TestLoad_NoCycle_ValidDAG(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// Valid DAG: D requires B and C, B requires A, C requires A
	content := `
version: "1"
checks:
  - id: a
    run: echo a
  - id: b
    run: echo b
    requires:
      - a
  - id: c
    run: echo c
    requires:
      - a
  - id: d
    run: echo d
    requires:
      - b
      - c
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error for valid DAG, got: %v", err)
	}

	if len(cfg.Checks) != 4 {
		t.Errorf("expected 4 checks, got: %d", len(cfg.Checks))
	}
}

func TestLoad_NoCycle_DiamondDependency(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// Diamond: A <- B, A <- C, B <- D, C <- D (D depends on both B and C which both depend on A)
	content := `
version: "1"
checks:
  - id: a
    run: echo a
  - id: b
    run: echo b
    requires:
      - a
  - id: c
    run: echo c
    requires:
      - a
  - id: d
    run: echo d
    requires:
      - b
      - c
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error for diamond dependency, got: %v", err)
	}

	if len(cfg.Checks) != 4 {
		t.Errorf("expected 4 checks, got: %d", len(cfg.Checks))
	}
}

func TestFormatCycle(t *testing.T) {
	tests := []struct {
		path     []string
		expected string
	}{
		{[]string{"a", "a"}, "a -> a"},
		{[]string{"a", "b", "a"}, "a -> b -> a"},
		{[]string{"a", "b", "c", "a"}, "a -> b -> c -> a"},
	}

	for _, tt := range tests {
		result := formatCycle(tt.path)
		if result != tt.expected {
			t.Errorf("formatCycle(%v) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}

func TestLoad_GrokSingleString(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: coverage
    run: go test -cover ./...
    grok: 'coverage: %{NUMBER:coverage}%'
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(cfg.Checks[0].Grok) != 1 {
		t.Errorf("expected 1 grok pattern, got: %d", len(cfg.Checks[0].Grok))
	}

	if cfg.Checks[0].Grok[0] != "coverage: %{NUMBER:coverage}%" {
		t.Errorf("expected grok pattern, got: %s", cfg.Checks[0].Grok[0])
	}
}

func TestLoad_GrokMultiplePatterns(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: lint
    run: golangci-lint run
    grok:
      - 'High: %{INT:high}'
      - 'Medium: %{INT:medium}'
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(cfg.Checks[0].Grok) != 2 {
		t.Errorf("expected 2 grok patterns, got: %d", len(cfg.Checks[0].Grok))
	}
}

func TestLoad_InvalidDuration(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    timeout: invalid
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestLoad_ValidDurations(t *testing.T) {
	tests := []struct {
		duration string
		expected time.Duration
	}{
		{"30s", 30 * time.Second},
		{"5m", 5 * time.Minute},
		{"1h", 1 * time.Hour},
		{"1m30s", 90 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.duration, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "vibeguard.yaml")

			content := `
version: "1"
checks:
  - id: test
    run: go test ./...
    timeout: ` + tt.duration

			if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}

			cfg, err := Load(configPath)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if cfg.Checks[0].Timeout != Duration(tt.expected) {
				t.Errorf("expected timeout %v, got: %v", tt.expected, cfg.Checks[0].Timeout)
			}
		})
	}
}

func TestLoad_WareSeverity(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"
checks:
  - id: lint
    run: golangci-lint run
    severity: warning
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Checks[0].Severity != SeverityWarning {
		t.Errorf("expected severity 'warning', got: %s", cfg.Checks[0].Severity)
	}
}

func TestFindConfigFile(t *testing.T) {
	// Create a temp directory and change to it
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	// No config file initially
	_, err := findConfigFile()
	if err == nil {
		t.Fatal("expected error when no config file exists")
	}

	// Create vibeguard.yaml
	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
`
	if err := os.WriteFile("vibeguard.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	path, err := findConfigFile()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if path != "vibeguard.yaml" {
		t.Errorf("expected 'vibeguard.yaml', got: %s", path)
	}
}

func TestFindConfigFile_Priority(t *testing.T) {
	// Test that config files are searched in priority order
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	content := `
version: "1"
checks:
  - id: test
    run: go test ./...
`
	// Create .vibeguard.yaml first
	if err := os.WriteFile(".vibeguard.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	path, err := findConfigFile()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if path != ".vibeguard.yaml" {
		t.Errorf("expected '.vibeguard.yaml', got: %s", path)
	}

	// Create vibeguard.yaml - should take priority
	if err := os.WriteFile("vibeguard.yaml", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	path, err = findConfigFile()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if path != "vibeguard.yaml" {
		t.Errorf("expected 'vibeguard.yaml' to take priority, got: %s", path)
	}
}

func TestDuration_AsDuration(t *testing.T) {
	d := Duration(5 * time.Minute)
	if d.AsDuration() != 5*time.Minute {
		t.Errorf("expected 5m, got: %v", d.AsDuration())
	}
}

func TestLoad_ComplexConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `
version: "1"

vars:
  MIN_COVERAGE: "80"
  packages: "./..."

checks:
  - id: fmt
    run: "! gofmt -l . | grep ."
    severity: error
    suggestion: "Run 'gofmt -w .' to format code"

  - id: vet
    run: go vet {{.packages}}
    severity: error

  - id: test
    run: go test -cover {{.packages}}
    requires:
      - vet
      - fmt
    timeout: 5m
    severity: error

  - id: lint
    run: golangci-lint run
    severity: warning
    requires:
      - vet
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(cfg.Checks) != 4 {
		t.Errorf("expected 4 checks, got: %d", len(cfg.Checks))
	}

	// Verify vars were stored
	if cfg.Vars["MIN_COVERAGE"] != "80" {
		t.Errorf("expected MIN_COVERAGE='80', got: %s", cfg.Vars["MIN_COVERAGE"])
	}

	// Verify interpolation happened
	vetCheck := cfg.Checks[1]
	if vetCheck.Run != "go vet ./..." {
		t.Errorf("expected interpolated vet command, got: %s", vetCheck.Run)
	}

	// Verify test check has correct requires
	testCheck := cfg.Checks[2]
	if len(testCheck.Requires) != 2 {
		t.Errorf("expected 2 requires for test, got: %d", len(testCheck.Requires))
	}
}

func TestConfigError(t *testing.T) {
	// Test ConfigError with cause
	cause := errors.New("underlying error")
	err := &ConfigError{Message: "test error", Cause: cause}

	expected := "test error: underlying error"
	if err.Error() != expected {
		t.Errorf("expected %q, got: %q", expected, err.Error())
	}

	if err.Unwrap() != cause {
		t.Error("expected Unwrap to return the cause")
	}

	// Test ConfigError without cause
	errNoCause := &ConfigError{Message: "standalone error"}
	if errNoCause.Error() != "standalone error" {
		t.Errorf("expected 'standalone error', got: %q", errNoCause.Error())
	}

	if errNoCause.Unwrap() != nil {
		t.Error("expected Unwrap to return nil when no cause")
	}
}

func TestIsConfigError(t *testing.T) {
	// Direct ConfigError
	configErr := &ConfigError{Message: "config error"}
	if !IsConfigError(configErr) {
		t.Error("expected IsConfigError to return true for ConfigError")
	}

	// Non-ConfigError
	regularErr := errors.New("regular error")
	if IsConfigError(regularErr) {
		t.Error("expected IsConfigError to return false for regular error")
	}

	// nil error
	if IsConfigError(nil) {
		t.Error("expected IsConfigError to return false for nil")
	}
}

func TestLoad_ReturnsConfigError(t *testing.T) {
	// Test that Load returns ConfigError for various failures

	t.Run("file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/vibeguard.yaml")
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsConfigError(err) {
			t.Errorf("expected ConfigError, got: %T", err)
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "vibeguard.yaml")
		if err := os.WriteFile(configPath, []byte("invalid: [yaml"), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(configPath)
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsConfigError(err) {
			t.Errorf("expected ConfigError, got: %T", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "vibeguard.yaml")
		content := `version: "1"
checks: []
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(configPath)
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsConfigError(err) {
			t.Errorf("expected ConfigError, got: %T", err)
		}
	})
}

func TestLoad_ConfigErrorWithLineNumbers(t *testing.T) {
	// Test that validation errors include line numbers
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	// Create a config with a duplicate check ID at line 7
	content := `version: "1"
checks:
  - id: test
    run: go test ./...
  - id: test
    run: go vet ./...
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for duplicate check id")
	}

	configErr, ok := err.(*ConfigError)
	if !ok {
		t.Fatalf("expected ConfigError, got: %T", err)
	}

	// Check that line number is present
	if configErr.LineNum == 0 {
		t.Error("expected non-zero line number for duplicate check error")
	}

	// Check that error message includes line number
	errMsg := configErr.Error()
	if !strings.Contains(errMsg, "line") {
		t.Errorf("expected error message to include 'line', got: %s", errMsg)
	}
}

func TestLoad_CyclicDependencyWithLineNumbers(t *testing.T) {
	// Test that cyclic dependency errors include line numbers
	dir := t.TempDir()
	configPath := filepath.Join(dir, "vibeguard.yaml")

	content := `version: "1"
checks:
  - id: a
    run: echo a
    requires:
      - b
  - id: b
    run: echo b
    requires:
      - a
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for cyclic dependency")
	}

	configErr, ok := err.(*ConfigError)
	if !ok {
		t.Fatalf("expected ConfigError, got: %T", err)
	}

	// Check that line number is present
	if configErr.LineNum == 0 {
		t.Error("expected non-zero line number for cyclic dependency error")
	}

	// Check that error message includes line number
	errMsg := configErr.Error()
	if !strings.Contains(errMsg, "line") {
		t.Errorf("expected error message to include 'line', got: %s", errMsg)
	}
}

// ExecutionError tests
func TestExecutionError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ExecutionError
		expected string
	}{
		{
			name: "with cause and line number",
			err: &ExecutionError{
				Message: "grok pattern failed",
				CheckID: "test-check",
				Cause:   errors.New("invalid pattern"),
				LineNum: 42,
			},
			expected: `grok pattern failed in check "test-check": invalid pattern (line 42)`,
		},
		{
			name: "with cause, no line number",
			err: &ExecutionError{
				Message: "assertion failed",
				CheckID: "build",
				Cause:   errors.New("coverage < 80%"),
			},
			expected: `assertion failed in check "build": coverage < 80%`,
		},
		{
			name: "no cause, with line number",
			err: &ExecutionError{
				Message: "check timed out",
				CheckID: "lint",
				LineNum: 15,
			},
			expected: `check timed out in check "lint" (line 15)`,
		},
		{
			name: "no cause, no line number",
			err: &ExecutionError{
				Message: "command failed",
				CheckID: "test",
			},
			expected: `command failed in check "test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExecutionError_Unwrap(t *testing.T) {
	cause := errors.New("underlying cause")
	err := &ExecutionError{
		Message: "wrapper",
		CheckID: "test",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}

	// Test with nil cause
	errNoCause := &ExecutionError{
		Message: "no cause",
		CheckID: "test",
	}
	if unwrapped := errNoCause.Unwrap(); unwrapped != nil {
		t.Errorf("Unwrap() with nil cause = %v, want nil", unwrapped)
	}
}

func TestIsExecutionError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "direct execution error",
			err:      &ExecutionError{Message: "test", CheckID: "check"},
			expected: true,
		},
		{
			name:     "wrapped execution error",
			err:      errors.Join(errors.New("wrapper"), &ExecutionError{Message: "test", CheckID: "check"}),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("not an execution error"),
			expected: false,
		},
		{
			name:     "config error",
			err:      &ConfigError{Message: "config issue"},
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsExecutionError(tt.err)
			if got != tt.expected {
				t.Errorf("IsExecutionError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoad_InvalidCheckID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		// Valid IDs
		{name: "simple lowercase", id: "test", wantErr: false},
		{name: "simple uppercase", id: "TEST", wantErr: false},
		{name: "mixed case", id: "myTest", wantErr: false},
		{name: "with underscore", id: "go_test", wantErr: false},
		{name: "with hyphen", id: "go-test", wantErr: false},
		{name: "with numbers", id: "test123", wantErr: false},
		{name: "underscore prefix", id: "_private", wantErr: false},
		{name: "complex valid", id: "Go_Test-123", wantErr: false},
		{name: "single letter", id: "a", wantErr: false},
		{name: "single underscore prefix", id: "_", wantErr: false},

		// Invalid IDs
		{name: "starts with number", id: "123test", wantErr: true},
		{name: "starts with hyphen", id: "-test", wantErr: true},
		{name: "contains space", id: "go test", wantErr: true},
		{name: "contains dot", id: "go.test", wantErr: true},
		{name: "contains colon", id: "go:test", wantErr: true},
		{name: "contains slash", id: "go/test", wantErr: true},
		{name: "contains special char", id: "test@check", wantErr: true},
		{name: "unicode characters", id: "tëst", wantErr: true},
		{name: "emoji", id: "test🚀", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "vibeguard.yaml")

			content := `
version: "1"
checks:
  - id: ` + tt.id + `
    run: echo hello
`
			if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}

			_, err := Load(configPath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for invalid check ID %q, got nil", tt.id)
				} else if !strings.Contains(err.Error(), "invalid id format") {
					t.Errorf("expected 'invalid id format' error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for valid check ID %q: %v", tt.id, err)
				}
			}
		})
	}
}

func TestValidCheckIDRegex(t *testing.T) {
	// Direct regex tests for edge cases
	tests := []struct {
		id    string
		valid bool
	}{
		{"a", true},
		{"A", true},
		{"_", true},
		{"_a", true},
		{"a1", true},
		{"a-b", true},
		{"a_b", true},
		{"A-B_c123", true},
		{"", false},
		{"1", false},
		{"1a", false},
		{"-a", false},
		{"a b", false},
		{"a.b", false},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := validCheckID.MatchString(tt.id)
			if got != tt.valid {
				t.Errorf("validCheckID.MatchString(%q) = %v, want %v", tt.id, got, tt.valid)
			}
		})
	}
}
