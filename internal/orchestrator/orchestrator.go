// Package orchestrator coordinates check execution with dependency management
// and parallel execution.
package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vibeguard/vibeguard/internal/assert"
	"github.com/vibeguard/vibeguard/internal/config"
	"github.com/vibeguard/vibeguard/internal/executor"
	"github.com/vibeguard/vibeguard/internal/grok"
)

// CheckResult represents the result of evaluating a single check.
type CheckResult struct {
	Check            *config.Check
	Execution        *executor.Result
	Passed           bool
	Extracted        map[string]string // Values extracted via grok patterns
	TriggeredPrompts []*TriggeredPrompt
}

// RunResult contains the complete results of running all checks.
type RunResult struct {
	Results           []*CheckResult
	Violations        []*Violation
	Duration          time.Duration
	ExitCode          int
	FailFastTriggered bool // True if execution was stopped early due to fail-fast
}

// TriggeredPrompt represents a prompt that was triggered by a check result.
type TriggeredPrompt struct {
	Event   string // "success", "failure", or "timeout"
	Source  string // "init", "code-review", etc. (prompt ID) or "inline" for inline content
	Content string // The actual prompt content
}

// Violation represents a check failure.
type Violation struct {
	CheckID          string
	Severity         config.Severity
	Command          string
	Suggestion       string
	Fix              string
	Extracted        map[string]string
	Timedout         bool
	LogFile          string // Path to log file containing check output
	TriggeredPrompts []*TriggeredPrompt
}

// TagFilter specifies which checks to include/exclude based on tags.
type TagFilter struct {
	Include []string // Run checks matching ANY of these tags (OR logic)
	Exclude []string // Exclude checks matching ANY of these tags (OR logic)
}

// Orchestrator coordinates check execution.
type Orchestrator struct {
	executor      *executor.Executor
	config        *config.Config
	maxParallel   int
	failFast      bool
	verbose       bool
	logDir        string // Directory for check output logs
	errorExitCode int    // Configurable exit code for failures (default: 1)
	tagFilter     *TagFilter
}

// DefaultLogDir is the default directory for check output logs.
const DefaultLogDir = ".vibeguard/log"

// Option is a functional option for configuring an Orchestrator.
type Option func(*Orchestrator)

// WithErrorExitCode sets the exit code to use for failures (FAIL and TIMEOUT).
func WithErrorExitCode(code int) Option {
	return func(o *Orchestrator) {
		o.errorExitCode = code
	}
}

// SetTagFilter sets the tag filter for selective check execution.
func (o *Orchestrator) SetTagFilter(filter TagFilter) {
	o.tagFilter = &filter
}

// New creates a new Orchestrator.
func New(cfg *config.Config, exec *executor.Executor, maxParallel int, failFast, verbose bool, logDir string, errorExitCode int) *Orchestrator {
	if maxParallel <= 0 {
		maxParallel = config.DefaultParallel
	}
	if logDir == "" {
		logDir = DefaultLogDir
	}
	if errorExitCode <= 0 {
		errorExitCode = 1
	}
	return &Orchestrator{
		executor:      exec,
		config:        cfg,
		maxParallel:   maxParallel,
		failFast:      failFast,
		verbose:       verbose,
		logDir:        logDir,
		errorExitCode: errorExitCode,
	}
}

// getAnalysisOutput returns the content to analyze for grok patterns and assertions.
// If the check specifies a file field, it reads from that file.
// Otherwise, it returns the command output.
func (o *Orchestrator) getAnalysisOutput(check *config.Check, execResult *executor.Result) (string, error) {
	if check.File != "" {
		// Interpolate variables in the file path
		filePath := o.interpolatePath(check.File)

		// Validate that the file path doesn't escape the current directory (prevent directory traversal)
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve file path %q: %w", filePath, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		absWd, err := filepath.Abs(wd)
		if err != nil {
			return "", fmt.Errorf("failed to resolve working directory: %w", err)
		}

		// Ensure the resolved path is within the working directory
		relPath, err := filepath.Rel(absWd, absPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			return "", fmt.Errorf("file path %q is outside the working directory", filePath)
		}

		content, err := os.ReadFile(absPath) // #nosec G304 - path is validated to be within working directory
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", filePath, err)
		}
		return string(content), nil
	}
	return execResult.Combined, nil
}

// interpolatePath performs variable substitution on a file path.
func (o *Orchestrator) interpolatePath(path string) string {
	result := path
	for key, value := range o.config.Vars {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// filterChecksByTags applies tag-based filtering to the checks.
// Returns a filtered list of checks and a set of excluded check IDs.
func (o *Orchestrator) filterChecksByTags(checks []config.Check) ([]config.Check, map[string]bool) {
	if o.tagFilter == nil || (len(o.tagFilter.Include) == 0 && len(o.tagFilter.Exclude) == 0) {
		// No filter, return all checks
		return checks, make(map[string]bool)
	}

	filtered := []config.Check{}
	excluded := make(map[string]bool)

	for _, check := range checks {
		// Check inclusion: if Include is specified, check must match at least one tag
		if len(o.tagFilter.Include) > 0 {
			matchesInclude := false
			for _, includeTag := range o.tagFilter.Include {
				for _, checkTag := range check.Tags {
					if includeTag == checkTag {
						matchesInclude = true
						break
					}
				}
				if matchesInclude {
					break
				}
			}
			if !matchesInclude {
				excluded[check.ID] = true
				continue
			}
		}

		// Check exclusion: skip if matches any exclude tag
		if len(o.tagFilter.Exclude) > 0 {
			matchesExclude := false
			for _, excludeTag := range o.tagFilter.Exclude {
				for _, checkTag := range check.Tags {
					if excludeTag == checkTag {
						matchesExclude = true
						break
					}
				}
				if matchesExclude {
					break
				}
			}
			if matchesExclude {
				excluded[check.ID] = true
				continue
			}
		}

		// Check passed both inclusion and exclusion filters
		filtered = append(filtered, check)
	}

	return filtered, excluded
}

// Run executes all checks and returns the results.
func (o *Orchestrator) Run(ctx context.Context) (*RunResult, error) {
	start := time.Now()

	// Apply tag filtering
	filteredChecks, excludedByTag := o.filterChecksByTags(o.config.Checks)

	// Pre-process filtered checks to identify those with missing dependencies
	// (dependencies excluded by tag filter, not genuinely unknown)
	// These checks will be skipped during execution with a warning
	checkIDs := make(map[string]bool)
	for _, check := range filteredChecks {
		checkIDs[check.ID] = true
	}

	// Build a set of all check IDs in the original config to detect tag-filtered dependencies
	allCheckIDs := make(map[string]bool)
	for _, check := range o.config.Checks {
		allCheckIDs[check.ID] = true
	}

	// Build a new list of checks that have all their dependencies available
	// and track checks with missing dependencies separately
	var validChecks []config.Check
	var skippedChecks []*config.Check // Checks with missing deps (tag-filtered) will be skipped
	for i := range filteredChecks {
		check := &filteredChecks[i]
		hasMissingDep := false
		for _, dep := range check.Requires {
			if !checkIDs[dep] {
				// Check if this dependency exists in the original config
				// If it does, it's excluded by tag filter, so we skip this check
				// If it doesn't, it's an unknown dependency - let BuildGraph handle the error
				if allCheckIDs[dep] {
					// This dependency is missing due to tag filtering
					// Mark this check as excluded so it gets skipped with a warning
					excludedByTag[check.ID] = true
					hasMissingDep = true
					skippedChecks = append(skippedChecks, check)
					break
				}
				// Otherwise, it's an unknown dependency - let it through to BuildGraph
				// which will return a proper error
			}
		}
		if !hasMissingDep {
			validChecks = append(validChecks, *check)
		}
	}
	filteredChecks = validChecks

	// Build dependency graph to determine execution order
	graph, err := BuildGraph(filteredChecks)
	if err != nil {
		return nil, err
	}

	// Build lookup maps for checks by ID and index by ID
	checkByID := make(map[string]*config.Check)
	checkIndexByID := make(map[string]int)
	for i := range filteredChecks {
		checkByID[filteredChecks[i].ID] = &filteredChecks[i]
		checkIndexByID[filteredChecks[i].ID] = i
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
			checkIndex := checkIndexByID[checkID]

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
				// Verify all dependencies passed and are not excluded by tag filter
				allDepsPassed := true
				missingDep := ""
				for _, depID := range check.Requires {
					if excludedByTag[depID] {
						allDepsPassed = false
						missingDep = depID
						break
					}
					if !passedChecks[depID] {
						allDepsPassed = false
						missingDep = depID
						break
					}
				}
				mu.Unlock()

				// Skip this check if a required dependency failed or is excluded by tag filter
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

					var suggestion string
					if excludedByTag[missingDep] {
						suggestion = fmt.Sprintf("Skipped: required dependency %q not in filtered set", missingDep)
					} else {
						suggestion = "Skipped: required dependency failed"
					}

					violation := &Violation{
						CheckID:    check.ID,
						Severity:   check.Severity,
						Command:    check.Run,
						Suggestion: suggestion,
						Fix:        check.Fix,
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

				// Write check output to log file (best-effort, don't fail if this fails)
				_ = o.writeCheckLog(check.ID, execResult.Combined)

				// Get the content to analyze (either from file or command output)
				analysisOutput, analysisErr := o.getAnalysisOutput(check, execResult)
				if analysisErr != nil {
					// Wrap file reading error with check context
					lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
					return &config.ExecutionError{
						Message:   analysisErr.Error(),
						Cause:     analysisErr,
						CheckID:   check.ID,
						LineNum:   lineNum,
						ErrorType: "file",
					}
				}

				// Apply grok patterns to extract values from output
				extracted := make(map[string]string)
				if len(check.Grok) > 0 {
					matcher, matcherErr := grok.New(check.Grok)
					if matcherErr != nil {
						// Wrap grok error with check context
						lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
						return &config.ExecutionError{
							Message:   "failed to compile grok pattern",
							Cause:     matcherErr,
							CheckID:   check.ID,
							LineNum:   lineNum,
							ErrorType: "grok",
						}
					}
					extracted, matcherErr = matcher.Match(analysisOutput)
					if matcherErr != nil {
						// Wrap grok error with check context
						lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
						return &config.ExecutionError{
							Message:   "failed to parse grok pattern",
							Cause:     matcherErr,
							CheckID:   check.ID,
							LineNum:   lineNum,
							ErrorType: "grok",
						}
					}
				}

				// Determine pass/fail based on exit code and assertion (if specified)
				passed := execResult.Success
				if passed && check.Assert != "" {
					evaluator := assert.New()
					assertPassed, assertErr := evaluator.Eval(check.Assert, extracted)
					if assertErr != nil {
						// Wrap assert error with check context
						lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
						return &config.ExecutionError{
							Message:   "failed to evaluate assertion",
							Cause:     assertErr,
							CheckID:   check.ID,
							LineNum:   lineNum,
							ErrorType: "assert",
						}
					}
					passed = assertPassed
				}

				result := &CheckResult{
					Check:            check,
					Execution:        execResult,
					Passed:           passed,
					Extracted:        extracted,
					TriggeredPrompts: o.evaluateTriggeredPrompts(check, passed, execResult.Timedout),
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
						CheckID:          check.ID,
						Severity:         check.Severity,
						Command:          check.Run,
						Suggestion:       suggestion,
						Fix:              check.Fix,
						Extracted:        result.Extracted,
						Timedout:         execResult.Timedout,
						LogFile:          filepath.Join(o.logDir, check.ID+".log"),
						TriggeredPrompts: o.evaluateTriggeredPrompts(check, passed, execResult.Timedout),
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

	// Add skipped checks (with missing dependencies) to results and violations
	for _, check := range skippedChecks {
		result := &CheckResult{
			Check:  check,
			Passed: false,
			Execution: &executor.Result{
				CheckID:  check.ID,
				ExitCode: -1,
				Success:  false,
			},
			Extracted: make(map[string]string),
		}

		// Find the missing dependency
		var missingDep string
		for _, dep := range check.Requires {
			if !checkIDs[dep] {
				missingDep = dep
				break
			}
		}

		violation := &Violation{
			CheckID:    check.ID,
			Severity:   check.Severity,
			Command:    check.Run,
			Suggestion: fmt.Sprintf("Skipped: required dependency %q not in filtered set", missingDep),
			Fix:        check.Fix,
			Extracted:  result.Extracted,
		}

		results = append(results, result)
		violations = append(violations, violation)
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
// Exit codes: 0 = success, errorExitCode = failure (both FAIL and TIMEOUT)
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

	// Both timeout and error use the same configurable exit code
	if hasTimeout || hasError {
		return o.errorExitCode
	}
	return executor.ExitCodeSuccess
}

// RunCheck executes a single check by ID.
func (o *Orchestrator) RunCheck(ctx context.Context, checkID string) (*RunResult, error) {
	start := time.Now()

	// Find the check and its index
	var check *config.Check
	var checkIndex int
	for i := range o.config.Checks {
		if o.config.Checks[i].ID == checkID {
			check = &o.config.Checks[i]
			checkIndex = i
			break
		}
	}
	if check == nil {
		return nil, &config.ConfigError{
			Message: fmt.Sprintf("check with ID %q not found", checkID),
		}
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

	// Write check output to log file (best-effort, don't fail if this fails)
	_ = o.writeCheckLog(check.ID, execResult.Combined)

	// Get the content to analyze (either from file or command output)
	analysisOutput, analysisErr := o.getAnalysisOutput(check, execResult)
	if analysisErr != nil {
		// Wrap file reading error with check context
		lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
		return nil, &config.ExecutionError{
			Message:   analysisErr.Error(),
			Cause:     analysisErr,
			CheckID:   check.ID,
			LineNum:   lineNum,
			ErrorType: "file",
		}
	}

	// Apply grok patterns to extract values from output
	extracted := make(map[string]string)
	if len(check.Grok) > 0 {
		matcher, matcherErr := grok.New(check.Grok)
		if matcherErr != nil {
			// Wrap grok error with check context
			lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
			return nil, &config.ExecutionError{
				Message:   "failed to compile grok pattern",
				Cause:     matcherErr,
				CheckID:   check.ID,
				LineNum:   lineNum,
				ErrorType: "grok",
			}
		}
		extracted, matcherErr = matcher.Match(analysisOutput)
		if matcherErr != nil {
			// Wrap grok error with check context
			lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
			return nil, &config.ExecutionError{
				Message:   "failed to parse grok pattern",
				Cause:     matcherErr,
				CheckID:   check.ID,
				LineNum:   lineNum,
				ErrorType: "grok",
			}
		}
	}

	// Determine pass/fail based on exit code and assertion (if specified)
	passed := execResult.Success
	if passed && check.Assert != "" {
		evaluator := assert.New()
		assertPassed, assertErr := evaluator.Eval(check.Assert, extracted)
		if assertErr != nil {
			// Wrap assert error with check context
			lineNum := o.config.FindCheckNodeLine(check.ID, checkIndex)
			return nil, &config.ExecutionError{
				Message:   "failed to evaluate assertion",
				Cause:     assertErr,
				CheckID:   check.ID,
				LineNum:   lineNum,
				ErrorType: "assert",
			}
		}
		passed = assertPassed
	}

	result := &CheckResult{
		Check:            check,
		Execution:        execResult,
		Passed:           passed,
		Extracted:        extracted,
		TriggeredPrompts: o.evaluateTriggeredPrompts(check, passed, execResult.Timedout),
	}

	var violations []*Violation
	exitCode := executor.ExitCodeSuccess
	if !passed {
		suggestion := check.Suggestion
		if execResult.Timedout {
			suggestion = "Check timed out. Consider increasing the timeout value or optimizing the command."
		}
		violation := &Violation{
			CheckID:          check.ID,
			Severity:         check.Severity,
			Command:          check.Run,
			Suggestion:       suggestion,
			Fix:              check.Fix,
			Extracted:        result.Extracted,
			Timedout:         execResult.Timedout,
			LogFile:          filepath.Join(o.logDir, check.ID+".log"),
			TriggeredPrompts: result.TriggeredPrompts,
		}
		violations = append(violations, violation)
		if execResult.Timedout || check.Severity == config.SeverityError {
			exitCode = o.errorExitCode
		}
	}

	return &RunResult{
		Results:    []*CheckResult{result},
		Violations: violations,
		Duration:   time.Since(start),
		ExitCode:   exitCode,
	}, nil
}

// evaluateTriggeredPrompts evaluates which event is triggered and returns the prompts to display.
// Event precedence: timeout > failure > success
func (o *Orchestrator) evaluateTriggeredPrompts(check *config.Check, passed bool, timedout bool) []*TriggeredPrompt {
	var eventName string
	var eventValue config.EventValue

	// Determine which event was triggered (precedence: timeout > failure > success)
	if timedout {
		eventName = "timeout"
		eventValue = check.On.Timeout
	} else if !passed {
		eventName = "failure"
		eventValue = check.On.Failure
	} else {
		eventName = "success"
		eventValue = check.On.Success
	}

	// Return early if no event handler is defined
	if eventName == "success" && !eventValue.IsInline && len(eventValue.IDs) == 0 {
		return nil
	}
	if eventName == "failure" && !eventValue.IsInline && len(eventValue.IDs) == 0 {
		return nil
	}
	if eventName == "timeout" && !eventValue.IsInline && len(eventValue.IDs) == 0 {
		return nil
	}

	// Handle inline content
	if eventValue.IsInline && eventValue.Content != "" {
		return []*TriggeredPrompt{
			{
				Event:   eventName,
				Source:  "inline",
				Content: eventValue.Content,
			},
		}
	}

	// Handle prompt ID references
	var triggered []*TriggeredPrompt
	for _, promptID := range eventValue.IDs {
		// Find the prompt in config
		for _, prompt := range o.config.Prompts {
			if prompt.ID == promptID {
				triggered = append(triggered, &TriggeredPrompt{
					Event:   eventName,
					Source:  promptID,
					Content: prompt.Content,
				})
				break
			}
		}
	}

	return triggered
}

// writeCheckLog writes check output to <logDir>/<check-id>.log
func (o *Orchestrator) writeCheckLog(checkID, output string) error {
	if err := os.MkdirAll(o.logDir, 0750); err != nil {
		return err
	}

	logPath := filepath.Join(o.logDir, checkID+".log")
	return os.WriteFile(logPath, []byte(output), 0600)
}
