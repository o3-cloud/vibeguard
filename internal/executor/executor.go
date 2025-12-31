// Package executor provides command execution capabilities for checks.
package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ExitCodeTimeout is the exit code used when a command times out.
const ExitCodeTimeout = 3

// Result contains the execution result of a check command.
type Result struct {
	CheckID  string
	ExitCode int
	Stdout   string
	Stderr   string
	Combined string
	Duration time.Duration
	Success  bool
	Timedout bool
	Error    error
}

// Executor runs check commands and captures their output.
type Executor struct {
	workDir string
	env     []string
}

// New creates a new Executor with optional working directory.
// If workDir is empty, the current working directory is used.
func New(workDir string) *Executor {
	if workDir == "" {
		workDir, _ = os.Getwd()
	}
	return &Executor{
		workDir: workDir,
		env:     os.Environ(),
	}
}

// Execute runs a command and captures its output.
func (e *Executor) Execute(ctx context.Context, checkID, command string) (*Result, error) {
	// Create command with shell
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir = e.workDir
	cmd.Env = e.env

	// Capture stdout and stderr separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute with timing
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	// Determine exit code and timeout status
	exitCode := 0
	timedout := false

	if err != nil {
		// Check if context was cancelled due to timeout
		if ctx.Err() == context.DeadlineExceeded {
			timedout = true
			exitCode = ExitCodeTimeout
			err = nil // Timeout is a recognized condition, not an error
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			err = nil // Non-zero exit is not an error for us
		}
	}

	// Build combined output (stdout + stderr)
	combined := stdout.String() + stderr.String()

	return &Result{
		CheckID:  checkID,
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Combined: combined,
		Duration: duration,
		Success:  exitCode == 0,
		Timedout: timedout,
		Error:    err,
	}, nil
}

// String returns a human-readable representation of the result.
func (r *Result) String() string {
	status := "passed"
	if r.Timedout {
		status = "timeout"
	} else if !r.Success {
		status = "failed"
	}
	return fmt.Sprintf("%s: %s (exit=%d, duration=%v)", r.CheckID, status, r.ExitCode, r.Duration)
}
