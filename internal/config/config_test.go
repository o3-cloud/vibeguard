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
