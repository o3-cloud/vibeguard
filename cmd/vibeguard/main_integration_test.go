package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMain_ExitCodes tests that the main binary returns correct exit codes
// for different error conditions.
func TestMain_ExitCodes(t *testing.T) {
	// Build the binary for testing
	binPath := buildTestBinary(t)

	t.Run("exit code 0 for successful check", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: pass
    run: "true"
    severity: error
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		if err != nil {
			t.Errorf("expected exit code 0, got error: %v", err)
		}
	})

	t.Run("exit code 1 for failing check (default)", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: fail
    run: "false"
    severity: error
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 1)
	})

	t.Run("exit code 2 for config error - missing config file", func(t *testing.T) {
		cmd := exec.Command(binPath, "check", "-c", "/nonexistent/path/vibeguard.yaml")
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Invalid YAML syntax - malformed list item
		configContent := `version: "1"
checks:
  - id: test
    run: "true"
  - this is not valid yaml [
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - no checks defined", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks: []
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - unsupported version", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "99"
checks:
  - id: test
    run: "true"
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - duplicate check ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: duplicate
    run: "true"
    timeout: 10s
  - id: duplicate
    run: "true"
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - invalid check ID format", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: "123-invalid"
    run: "true"
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("exit code 2 for config error - cyclic dependency", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: check_a
    run: "true"
    timeout: 10s
    requires:
      - check_b
  - id: check_b
    run: "true"
    timeout: 10s
    requires:
      - check_a
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		assertExitCode(t, err, 2)
	})

	t.Run("custom exit code via --error-exit-code flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: fail
    run: "false"
    severity: error
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"), "--error-exit-code", "42")
		err := cmd.Run()
		assertExitCode(t, err, 42)
	})

	t.Run("warning severity does not cause non-zero exit", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: warn
    run: "false"
    severity: warning
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		err := cmd.Run()
		if err != nil {
			t.Errorf("expected exit code 0 for warning severity, got error: %v", err)
		}
	})
}

// TestMain_ErrorOutput tests that error messages are written to stderr.
func TestMain_ErrorOutput(t *testing.T) {
	binPath := buildTestBinary(t)

	t.Run("config error message written to stderr", func(t *testing.T) {
		cmd := exec.Command(binPath, "check", "-c", "/nonexistent/path/vibeguard.yaml")
		output, _ := cmd.CombinedOutput()
		if len(output) == 0 {
			t.Error("expected error message in output")
		}
	})

	t.Run("no output for successful check (silence is success)", func(t *testing.T) {
		tmpDir := t.TempDir()
		configContent := `version: "1"
checks:
  - id: pass
    run: "true"
    severity: error
    timeout: 10s
`
		writeConfig(t, tmpDir, configContent)

		cmd := exec.Command(binPath, "check", "-c", filepath.Join(tmpDir, "vibeguard.yaml"))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if len(output) != 0 {
			t.Errorf("expected no output for successful check, got: %s", output)
		}
	})
}

// buildTestBinary builds the vibeguard binary for testing and returns its path.
func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create a temp directory for the binary
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "vibeguard")

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Dir(os.Args[0])
	// Find the actual source directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	// Navigate to the cmd/vibeguard directory if we're not already there
	srcDir := wd
	if filepath.Base(srcDir) != "vibeguard" || filepath.Base(filepath.Dir(srcDir)) != "cmd" {
		srcDir = filepath.Join(wd, "cmd", "vibeguard")
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			// Try from project root
			srcDir = wd
		}
	}
	cmd.Dir = srcDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build test binary: %v\noutput: %s", err, output)
	}

	return binPath
}

// writeConfig writes a config file to the given directory.
func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	configPath := filepath.Join(dir, "vibeguard.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
}

// assertExitCode checks that the error is an exit error with the expected code.
func assertExitCode(t *testing.T, err error, expected int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected exit code %d, got no error", expected)
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != expected {
		t.Errorf("expected exit code %d, got %d", expected, exitErr.ExitCode())
	}
}
