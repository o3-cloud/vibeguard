# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.

---

## Project Analysis

**Project Name:** github.com/vibeguard/vibeguard
**Project Type:** go
**Detection Confidence:** 80%
**Language Version:** 

### Main Tools Detected:
- gofmt
- go vet
- go test


### Project Structure:
- Source directories: internal
- Test directories: 
- Entry points: cmd/vibeguard/main.go

### Build System:
- Build output: default

---

## Recommended Checks

Based on the detected tools, here are the recommended checks:


### build (build)
**Description:** Verify Go code compiles successfully
**Rationale:** Catch compilation errors before they reach CI
**Command:** `go build ./...`
**Severity:** error


**Suggestion on failure:** Fix compilation errors before committing.



### fmt (format)
**Description:** Check Go code formatting with gofmt
**Rationale:** Consistent formatting improves readability and reduces diffs
**Command:** `test -z "$(gofmt -l .)"`
**Severity:** error


**Suggestion on failure:** Run 'gofmt -w .' to format your Go code.



### vet (lint)
**Description:** Run go vet to detect suspicious constructs
**Rationale:** go vet finds bugs that the compiler doesn't catch, like incorrect printf format strings
**Command:** `go vet ./...`
**Severity:** error


**Suggestion on failure:** Fix the issues reported by go vet. These are often real bugs.



### test (test)
**Description:** Run Go tests
**Rationale:** Tests verify that code behaves as expected
**Command:** `go test ./...`
**Severity:** error


**Suggestion on failure:** Fix failing tests before committing.



### coverage (test)
**Description:** Check test coverage meets minimum threshold
**Rationale:** Code coverage helps identify untested code paths
**Command:** `go test -cover ./... 2>&1 | tail -1`
**Severity:** warning
**Grok Patterns:** `coverage: %{NUMBER:coverage}%` 
**Assertion:** `coverage >= 70`
**Suggestion on failure:** Coverage is {{.coverage}}%, target is 70%. Add tests to improve coverage.
**Requires:** test



---

## Configuration Requirements

A valid vibeguard.yaml must contain:

### 1. Version (required)
```yaml
version: "1"
```

### 2. Variables (optional)
Global variables for interpolation in commands and assertions.

```yaml
vars:
  go_packages: "./..."
  test_dir: "./..."
```

### 3. Checks (required)
Array of check definitions. Each check must have:
- **id:** Unique identifier (alphanumeric + underscore + hyphen)
- **run:** Shell command to execute

Optional fields:
- **grok:** Array of patterns to extract data from output (uses Grok syntax)
- **assert:** Condition that must be true (e.g., "coverage >= 70")
- **requires:** Array of check IDs that must pass first
- **severity:** "error" or "warning" (default: error)
- **suggestion:** Message shown on failure (supports {{`.variable`}} templating)
- **timeout:** Duration string (e.g., "30s", "5m")

---

## Go-Specific Examples

### Format Check
```yaml
- id: fmt
  run: test -z "$(gofmt -l .)"
  severity: error
  suggestion: "Run 'gofmt -w .' to format code"
  timeout: 5s
```

### Lint Check
```yaml
- id: lint
  run: golangci-lint run {{`.go_packages`}}
  severity: error
  suggestion: "Fix linting issues. Run 'golangci-lint run --fix' for auto-fixes."
  timeout: 30s
```

### Test with Coverage
```yaml
- id: test
  run: go test {{`.go_packages`}}
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 60s

- id: coverage
  run: go test ./... -coverprofile cover.out && go tool cover -func cover.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= 70"
  severity: warning
  suggestion: "Coverage is {{`.coverage`}}%, target is 70%. Add tests to improve."
  requires:
    - test
  timeout: 60s
```

### Build Check with Dependency
```yaml
- id: build
  run: go build {{`.go_packages`}}
  severity: error
  suggestion: "Fix compilation errors"
  timeout: 30s
  requires:
    - vet
```

---

## Validation Rules

Your generated configuration must:

1. Be valid YAML syntax
2. Have `version: "1"` at the top level
3. Include at least one check in the `checks` array
4. Each check must have a unique `id`
5. Each check must have a non-empty `run` command
6. All `requires` references must point to existing check IDs
7. No circular dependencies in `requires`
8. All variables used in double-curly-brace syntax (e.g., .var) must be defined in vars
9. `severity` must be "error" or "warning"
10. Grok patterns must be valid Grok syntax

**DO NOT:**
- Include YAML comments in the generated config
- Add extra top-level keys beyond version, vars, checks
- Use undefined variables
- Create checks for tools not detected in the project

---

## Your Task

Based on the project analysis above, generate a vibeguard.yaml configuration that:

1. Includes version: "1"
2. Defines appropriate variables for this project
3. Creates checks for the detected tools
4. Follows the syntax rules described above
5. Includes helpful suggestions for each check
6. Uses appropriate timeouts for each check type

Output the configuration in a YAML code block:

```yaml
# vibeguard.yaml for github.com/vibeguard/vibeguard
version: "1"

vars:
  # ... your variables ...

checks:
  # ... your checks ...
```

After generating the configuration, verify it would pass the validation rules listed above.
