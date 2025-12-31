// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
)

// CheckResult represents the result of evaluating a single check.
type CheckResult struct {
	Check     *config.Check
	Execution *executor.Result
	Passed    bool
	Extracted map[string]string // Values extracted via grok patterns
}

// RunResult contains the complete results of running all checks.
type RunResult struct {
	Results           []*CheckResult
	Violations        []*Violation
	Duration          time.Duration
	ExitCode          int
	FailFastTriggered bool // True if execution was stopped early due to fail-fast
}

// Violation represents a check failure.
type Violation struct {
	CheckID    string
	Severity   config.Severity
	Command    string
	Suggestion string
	Extracted  map[string]string
	Timedout   bool
}

// Orchestrator coordinates check execution.
type Orchestrator struct {
	executor    *executor.Executor
	config      *config.Config
	maxParallel int
	failFast    bool
	verbose     bool
}

// New creates a new Orchestrator.
func New(cfg *config.Config, exec *executor.Executor, maxParallel int, failFast, verbose bool) *Orchestrator {
	if maxParallel <= 0 {
		maxParallel = config.DefaultParallel
	}
	return &Orchestrator{
		executor:    exec,
		config:      cfg,
		maxParallel: maxParallel,
		failFast:    failFast,
		verbose:     verbose,
	}
}

// Run executes all checks and returns the results.
func (o *Orchestrator) Run(ctx context.Context) (*RunResult, error) {
	start := time.Now()

	// Build dependency graph to determine execution order
	graph, err := BuildGraph(o.config.Checks)
	if err != nil {
		return nil, err
	}

	// Build lookup map for checks by ID
	checkByID := make(map[string]*config.Check)
	for i := range o.config.Checks {
		checkByID[o.config.Checks[i].ID] = &o.config.Checks[i]
	}

	results := make([]*CheckResult, 0, len(o.config.Checks))
	violations := make([]*Violation, 0)

	// Track which checks have passed (for dependency validation)
	passedChecks := make(map[string]bool)

	// Mutex for thread-safe access to shared state
	var mu sync.Mutex
	// Flag to signal fail-fast termination
	failFastTriggered := false

	// Create a cancellable context for fail-fast cancellation
	failFastCtx, cancelFailFast := context.WithCancel(ctx)
	defer cancelFailFast()

	// Execute checks level by level (topological order)
	// Within each level, checks run in parallel (limited by maxParallel)
	for _, level := range graph.Levels() {
		// Check if fail-fast was triggered in a previous level
		if failFastTriggered {
			break
		}

		// Create results slice for this level to maintain order within level
		levelResults := make([]*CheckResult, len(level))
		levelViolations := make([]*Violation, 0)

		// Use errgroup for parallel execution with context cancellation
		g, gctx := errgroup.WithContext(failFastCtx)

		// Semaphore to limit concurrency
		sem := make(chan struct{}, o.maxParallel)

		for i, checkID := range level {
			i, checkID := i, checkID // capture for goroutine
			check := checkByID[checkID]

			g.Go(func() error {
				// Acquire semaphore
				select {
				case sem <- struct{}{}:
				case <-gctx.Done():
					return gctx.Err()
				}
				defer func() { <-sem }()

				// Check if fail-fast was triggered by another goroutine
				mu.Lock()
				if failFastTriggered {
					mu.Unlock()
					return nil
				}
				// Verify all dependencies passed
				allDepsPassed := true
				for _, depID := range check.Requires {
					if !passedChecks[depID] {
						allDepsPassed = false
						break
					}
				}
				mu.Unlock()

				// Skip this check if a required dependency failed
				if !allDepsPassed {
					result := &CheckResult{
						Check:  check,
						Passed: false,
						Execution: &executor.Result{
							CheckID:  checkID,
							ExitCode: -1,
							Success:  false,
						},
						Extracted: make(map[string]string),
					}

					violation := &Violation{
						CheckID:    check.ID,
						Severity:   check.Severity,
						Command:    check.Run,
						Suggestion: "Skipped: required dependency failed",
						Extracted:  result.Extracted,
					}

					mu.Lock()
					levelResults[i] = result
					levelViolations = append(levelViolations, violation)
					mu.Unlock()
					return nil
				}

				// Apply timeout
				checkCtx := gctx
				var cancel context.CancelFunc
				if check.Timeout > 0 {
					checkCtx, cancel = context.WithTimeout(gctx, check.Timeout.AsDuration())
				}

				// Execute the check
				execResult, execErr := o.executor.Execute(checkCtx, check.ID, check.Run)
				if cancel != nil {
					cancel()
				}
				if execErr != nil {
					// Execution error (not just non-zero exit)
					return execErr
				}

				// For Phase 1, pass/fail is based on exit code only
				passed := execResult.Success

				result := &CheckResult{
					Check:     check,
					Execution: execResult,
					Passed:    passed,
					Extracted: make(map[string]string),
				}

				mu.Lock()
				levelResults[i] = result
				passedChecks[checkID] = passed

				if !passed {
					suggestion := check.Suggestion
					if execResult.Timedout {
						suggestion = "Check timed out. Consider increasing the timeout value or optimizing the command."
					}
					violation := &Violation{
						CheckID:    check.ID,
						Severity:   check.Severity,
						Command:    check.Run,
						Suggestion: suggestion,
						Extracted:  result.Extracted,
						Timedout:   execResult.Timedout,
					}
					levelViolations = append(levelViolations, violation)

					if o.failFast && check.Severity == config.SeverityError {
						failFastTriggered = true
						cancelFailFast() // Cancel in-flight checks
					}
				}
				mu.Unlock()

				return nil
			})
		}

		// Wait for all goroutines in this level to complete
		if err := g.Wait(); err != nil {
			// If fail-fast was triggered, context.Canceled is expected
			if failFastTriggered && err == context.Canceled {
				// Continue to collect results from this level
			} else {
				return nil, err
			}
		}

		// Append level results in order
		for _, r := range levelResults {
			if r != nil {
				results = append(results, r)
			}
		}
		violations = append(violations, levelViolations...)

		// If fail-fast was triggered, stop processing further levels
		if failFastTriggered {
			break
		}
	}

	return &RunResult{
		Results:           results,
		Violations:        violations,
		Duration:          time.Since(start),
		ExitCode:          o.calculateExitCode(violations),
		FailFastTriggered: failFastTriggered,
	}, nil
}

// calculateExitCode determines the exit code based on violations.
// Exit codes: 0 = success, 1 = error severity violation, 3 = timeout
func (o *Orchestrator) calculateExitCode(violations []*Violation) int {
	hasTimeout := false
	hasError := false

	for _, v := range violations {
		if v.Timedout {
			hasTimeout = true
		}
		if v.Severity == config.SeverityError {
			hasError = true
		}
	}

	// Timeout takes precedence over error severity
	if hasTimeout {
		return executor.ExitCodeTimeout
	}
	if hasError {
		return 1
	}
	return 0
}

// RunCheck executes a single check by ID.
func (o *Orchestrator) RunCheck(ctx context.Context, checkID string) (*RunResult, error) {
	start := time.Now()

	// Find the check
	var check *config.Check
	for i := range o.config.Checks {
		if o.config.Checks[i].ID == checkID {
			check = &o.config.Checks[i]
			break
		}
	}
	if check == nil {
		return &RunResult{
			Duration: time.Since(start),
			ExitCode: 2, // Configuration error
		}, nil
	}

	// Apply timeout
	checkCtx := ctx
	if check.Timeout > 0 {
		var cancel context.CancelFunc
		checkCtx, cancel = context.WithTimeout(ctx, check.Timeout.AsDuration())
		defer cancel()
	}

	// Execute the check
	execResult, err := o.executor.Execute(checkCtx, check.ID, check.Run)
	if err != nil {
		return nil, err
	}

	passed := execResult.Success
	result := &CheckResult{
		Check:     check,
		Execution: execResult,
		Passed:    passed,
		Extracted: make(map[string]string),
	}

	var violations []*Violation
	exitCode := 0
	if !passed {
		suggestion := check.Suggestion
		if execResult.Timedout {
			suggestion = "Check timed out. Consider increasing the timeout value or optimizing the command."
		}
		violation := &Violation{
			CheckID:    check.ID,
			Severity:   check.Severity,
			Command:    check.Run,
			Suggestion: suggestion,
			Extracted:  result.Extracted,
			Timedout:   execResult.Timedout,
		}
		violations = append(violations, violation)
		if execResult.Timedout {
			exitCode = executor.ExitCodeTimeout
		} else if check.Severity == config.SeverityError {
			exitCode = 1
		}
	}

	return &RunResult{
		Results:    []*CheckResult{result},
		Violations: violations,
		Duration:   time.Since(start),
		ExitCode:   exitCode,
	}, nil
}
