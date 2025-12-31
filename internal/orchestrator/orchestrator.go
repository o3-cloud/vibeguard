// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"context"
	"time"

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
	Results    []*CheckResult
	Violations []*Violation
	Duration   time.Duration
	ExitCode   int
}

// Violation represents a check failure.
type Violation struct {
	CheckID    string
	Severity   config.Severity
	Command    string
	Suggestion string
	Extracted  map[string]string
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

	// Execute checks level by level (topological order)
	// Within each level, checks can be run in parallel (future: vibeguard-v3m.3)
	// For now, execute sequentially within each level
	for _, level := range graph.Levels() {
		for _, checkID := range level {
			check := checkByID[checkID]

			// Verify all dependencies passed (fail-fast may have stopped a dependency)
			allDepsPassed := true
			for _, depID := range check.Requires {
				if !passedChecks[depID] {
					allDepsPassed = false
					break
				}
			}

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
				results = append(results, result)

				violation := &Violation{
					CheckID:    check.ID,
					Severity:   check.Severity,
					Command:    check.Run,
					Suggestion: "Skipped: required dependency failed",
					Extracted:  result.Extracted,
				}
				violations = append(violations, violation)
				continue
			}

			// Apply timeout
			checkCtx := ctx
			var cancel context.CancelFunc
			if check.Timeout > 0 {
				checkCtx, cancel = context.WithTimeout(ctx, check.Timeout.AsDuration())
			}

			// Execute the check
			execResult, execErr := o.executor.Execute(checkCtx, check.ID, check.Run)
			if cancel != nil {
				cancel()
			}
			if execErr != nil {
				// Execution error (not just non-zero exit)
				return nil, execErr
			}

			// For Phase 1, pass/fail is based on exit code only
			// Grok/assert will be added in Phase 2
			passed := execResult.Success
			passedChecks[checkID] = passed

			result := &CheckResult{
				Check:     check,
				Execution: execResult,
				Passed:    passed,
				Extracted: make(map[string]string),
			}
			results = append(results, result)

			if !passed {
				violation := &Violation{
					CheckID:    check.ID,
					Severity:   check.Severity,
					Command:    check.Run,
					Suggestion: check.Suggestion,
					Extracted:  result.Extracted,
				}
				violations = append(violations, violation)

				if o.failFast {
					// Return early, remaining checks won't be executed
					return &RunResult{
						Results:    results,
						Violations: violations,
						Duration:   time.Since(start),
						ExitCode:   o.calculateExitCode(violations),
					}, nil
				}
			}
		}
	}

	return &RunResult{
		Results:    results,
		Violations: violations,
		Duration:   time.Since(start),
		ExitCode:   o.calculateExitCode(violations),
	}, nil
}

// calculateExitCode determines the exit code based on violations.
func (o *Orchestrator) calculateExitCode(violations []*Violation) int {
	for _, v := range violations {
		if v.Severity == config.SeverityError {
			return 1
		}
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
		violation := &Violation{
			CheckID:    check.ID,
			Severity:   check.Severity,
			Command:    check.Run,
			Suggestion: check.Suggestion,
			Extracted:  result.Extracted,
		}
		violations = append(violations, violation)
		if check.Severity == config.SeverityError {
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
