---
summary: Comprehensive analysis of vibeguard error handling and identified 5 issues for task o6g.1
event_type: deep dive
sources:
  - internal/config/config.go
  - internal/executor/executor.go
  - internal/orchestrator/orchestrator.go
  - internal/grok/grok.go
  - internal/assert/eval.go
tags:
  - error-handling
  - code-quality
  - user-experience
  - vibeguard-o6g.1
  - exit-codes
  - configuration
---

# Vibeguard Error Message Analysis

## Overview
Conducted comprehensive analysis of error handling across the vibeguard codebase for task **vibeguard-o6g.1 (Comprehensive error messages)**. The current implementation has a solid foundation but lacks important context and consistency in error reporting.

## Current Strengths

- ✓ Well-defined exit codes (0=success, 2=violation, 3=config error, 4=timeout)
- ✓ ConfigError type properly wraps configuration errors with context
- ✓ Error cause chains using %w for proper error wrapping
- ✓ Violation formatting includes severity, command, and suggestions
- ✓ Timeout-specific helpful suggestions provided

## Issues Identified

### 1. vibeguard-cce (HIGH) - RunCheck Silent Failure for Unknown Checks
**Location:** internal/orchestrator/orchestrator.go:323-327

When `RunCheck()` is called with an unknown check ID, it silently returns `ExitCodeConfigError` (3) without any error message. This makes debugging harder and breaks normal CLI workflows.

```go
// Current: returns no error message
if check == nil {
    return &RunResult{
        Duration: time.Since(start),
        ExitCode: executor.ExitCodeConfigError,
    }, nil
}
```

**Impact:** Users see exit code 3 but no explanation of what went wrong.

### 2. vibeguard-bd9 (HIGH) - Config Validation Missing YAML Line Numbers
**Location:** internal/config/config.go:100-157 (Validate method)

Configuration validation errors show array indices instead of YAML file line numbers. This makes it difficult for users to locate the problematic line in their config file.

**Current example:** "check at index 5 has no id"
**Better example:** "vibeguard.yaml:42: check at index 5 has no id"

**Technical Note:** YAML unmarshaling needs to preserve `yaml.Node.Line` information during config loading.

### 3. vibeguard-trb (HIGH) - Grok/Assert Errors Lack Check Context
**Locations:**
- internal/grok/grok.go:31-33 (pattern compilation)
- internal/assert/eval.go:76-77 (parse errors)
- internal/assert/parser.go (multiple error sites)

Error messages from grok patterns and assertion expressions don't indicate which check owns them, making diagnosis difficult.

**Current examples:**
- "failed to compile grok pattern %d" - shows index, not pattern or check
- "parse error: unexpected token 'x' at position 15" - no check context

**Better examples:**
- "Check 'parse-log' has invalid grok pattern '...(pattern)...': failed to compile - invalid pattern syntax"
- "Check 'fmt' assertion 'count(x) > 100' has parse error: unexpected token 'x' at position 15"

### 4. vibeguard-8p1 (MEDIUM) - Exit Code Inconsistency
**Locations:**
- cmd/vibeguard/main.go:14-16 (exit logic)
- internal/orchestrator/orchestrator.go:312-327 (RunCheck)
- internal/orchestrator/graph.go:23-29 (dependency validation)

Inconsistent handling of config vs execution errors:
- `RunCheck()` returns ExitCodeConfigError (3) for missing check (should be error)
- Dependency validation errors from `BuildGraph()` aren't wrapped as ConfigError
- Main handler doesn't distinguish between config-time and execution-time errors
- Need clear separation: config-time errors → exit 2, execution errors → exit 3+

### 5. vibeguard-eb4 (MEDIUM) - Error Messages Need Actionable Suggestions
**Scattered across multiple files**

Many error messages could be improved with actionable guidance:

**Example improvements needed:**
1. Unknown check: Show available check IDs
   - "check 'build-all' not found. Available checks: lint, test, build, deploy"

2. Invalid severity: Show valid options
   - "check 'fmt' has invalid severity: 'invalid'. Must be 'error' or 'warning'"

3. Grok pattern errors: Provide debugging resources
   - "Check 'parse-log' has invalid grok pattern: ... See https://grokdebug.herokuapp.com"

4. Timeout errors: Show how to fix
   - "Check 'benchmark' timed out after 30s. To increase: add 'timeout: 60s' to check config"

5. Circular dependencies: Provide fix instructions
   - "circular dependency detected: a -> b -> c -> a. Check the 'requires' fields and break this cycle."

## Recommended Priority

1. **vibeguard-cce** (Blocks workflows) - Fix RunCheck to return proper error
2. **vibeguard-bd9** (Major usability) - Add YAML line numbers to validation errors
3. **vibeguard-trb** (Major usability) - Add check context to grok/assert errors
4. **vibeguard-8p1** (CI/CD integration) - Fix exit code handling consistency
5. **vibeguard-eb4** (UX improvement) - Make error messages actionable

## Related ADRs
- ADR-004: Code Quality Standards - error messages are part of quality
- ADR-006: Git Pre-Commit Hook Integration - clear errors critical for hook compatibility

## Next Steps
Address the identified issues in priority order. RunCheck fix and YAML line numbers are blocking UX issues that should be tackled first.
