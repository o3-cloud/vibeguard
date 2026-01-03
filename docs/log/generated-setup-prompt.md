# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.

---

## Project Analysis

**Project Name:** github.com/vibeguard/vibeguard
**Project Type:** go
**Detection Confidence:** 100%

### Main Tools Detected:
- golangci-lint
- gofmt
- go vet
- go test
- goimports

### Project Structure:
- Source directories: internal
- Entry points: cmd/vibeguard/main.go

### Build System:
- Build output: default


---

## Tooling Inspection Instructions

Before generating the configuration, **inspect the existing tool configurations** in this project to understand how they're set up:

### What to Look For

When inspecting each configuration file:

1. **Enabled rules/checks**: What specific rules or checks are enabled?
2. **Disabled rules**: Are any rules explicitly disabled? (They may conflict with project needs)
3. **Custom settings**: Are there custom paths, thresholds, or exclusions?
4. **Integration points**: Does the tool integrate with other tools (e.g., editor configs, CI)?
5. **Command variations**: Are there specific flags or options being used?

Use this understanding to craft VibeGuard checks that:
- Run tools with the same flags/options as configured
- Respect existing exclusions and thresholds
- Maintain consistency with the project's existing standards


---

## Additional Tooling Research

Based on the project type, consider recommending additional quality and security tools that aren't currently configured.

### Suggested Tools to Consider

#### staticcheck
**Category:** Linter
**Purpose:** Advanced static analysis for Go beyond go vet
**Value:** Catches bugs, suggests simplifications, and identifies deprecated code
**Example command:** `staticcheck ./...`

#### gosec
**Category:** Security
**Purpose:** Security-focused static analysis for Go code
**Value:** Identifies potential security issues like SQL injection, hardcoded credentials
**Example command:** `gosec ./...`

#### errcheck
**Category:** Linter
**Purpose:** Checks for unchecked errors in Go code
**Value:** Ensures error handling is not accidentally ignored
**Example command:** `errcheck ./...`

#### ineffassign
**Category:** Linter
**Purpose:** Detects ineffectual assignments in Go code
**Value:** Finds variables assigned but never used
**Example command:** `ineffassign ./...`

#### govulncheck
**Category:** Security
**Purpose:** Checks Go dependencies for known vulnerabilities
**Value:** Official Go team tool for CVE scanning in dependencies
**Example command:** `govulncheck ./...`

### How to Present Suggestions

When suggesting additional tools to the user:

1. **Explain the benefit**: Why would this tool help the project?
2. **Assess compatibility**: Will it work well with existing tools?
3. **Provide options**: Let the user decide which suggestions to include
4. **Include installation notes**: If a tool requires installation, mention it in the suggestion field

Ask the user: "Would you like me to include checks for any of these additional tools?"


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
  packages: "./..."
  coverage_threshold: "70"
```

### 3. Checks (required)
Array of check definitions. Each check must have:
- **id:** Unique identifier (alphanumeric + underscore + hyphen)
- **run:** Shell command to execute

Optional fields:
- **grok:** Array of patterns to extract data from output
- **assert:** Condition that must be true
- **requires:** Array of check IDs that must pass first
- **severity:** "error" or "warning" (default: error)
- **suggestion:** Message shown on failure
- **timeout:** Duration string (e.g., "30s", "5m")
- **file:** Path to read output from instead of command stdout

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

### Lint Check with golangci-lint
```yaml
- id: lint
  run: golangci-lint run {{.go_packages}}
  severity: error
  suggestion: "Fix linting issues. Run 'golangci-lint run --fix' for auto-fixes."
  timeout: 60s
  requires:
    - fmt
    - vet
```

### Go Vet Check
```yaml
- id: vet
  run: go vet {{.go_packages}}
  severity: error
  suggestion: "Fix go vet issues. These often indicate real bugs."
  timeout: 30s
```

### Test with Coverage
```yaml
- id: test
  run: go test -race {{.go_packages}}
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s

- id: coverage
  run: go test ./... -coverprofile cover.out && go tool cover -func cover.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%. Add tests."
  requires:
    - test
  timeout: 300s
```

### Build Check
```yaml
- id: build
  run: go build -o /dev/null {{.go_packages}}
  severity: error
  suggestion: "Fix compilation errors before committing"
  timeout: 60s
  requires:
    - vet
```

---

## YAML Syntax Requirements

The vibeguard.yaml file must be valid YAML. Follow these rules:

### Basic Structure
- The file must start with valid YAML syntax
- Use 2-space indentation consistently
- Strings with special characters must be quoted
- Arrays can use either block style (- item) or flow style ([item1, item2])

### Required Top-Level Keys
The configuration must have these top-level keys:
- **version** (required): Must be "1" (as a string)
- **vars** (optional): Map of variable names to string values
- **checks** (required): Array of check definitions

### String Quoting Rules
Quote strings that contain:
- Colons followed by space (: )
- Special YAML characters: {, }, [, ], ,, &, *, #, ?, |, -, <, >, =, !, %, @, \
- Leading/trailing spaces
- Numbers that should be treated as strings
- Boolean-like values: yes, no, true, false, on, off

### Example Valid Structure
```yaml
version: "1"

vars:
  packages: "./..."
  coverage_threshold: "70"

checks:
  - id: fmt
    run: gofmt -l .
    severity: error
```

### Common Syntax Errors to Avoid
1. Missing quotes around version number: version: 1 (wrong) vs version: "1" (correct)
2. Tabs instead of spaces for indentation
3. Inconsistent indentation levels
4. Missing space after colon in key-value pairs
5. Unquoted strings with special characters

## Check Structure Requirements

Each check in the **checks** array must follow this structure:

### Required Fields
- **id** (string, required): Unique identifier for the check
  - Must start with a letter or underscore
  - Can contain letters, numbers, underscores, and hyphens
  - Must be unique across all checks
  - Examples: "fmt", "lint", "go-test", "npm_audit", "_private"

- **run** (string, required): Shell command to execute
  - Must be non-empty
  - Will be executed via shell (sh -c)
  - Can reference variables using {{.varname}} syntax

### Optional Fields
- **grok** (string or array of strings): Patterns to extract data from command output
  - Uses Grok syntax (similar to Logstash)
  - Common patterns: %{NUMBER:varname}, %{WORD:varname}, %{GREEDYDATA:varname}
  - Extracted values can be used in assertions

- **assert** (string): Condition that must be true for the check to pass
  - Supports comparisons: ==, !=, <, <=, >, >=
  - Supports boolean operators: &&, ||, !
  - References extracted grok values: coverage >= 70
  - References special variables: exit_code == 0, stdout == ""

- **severity** (string): "error" or "warning"
  - Default: "error"
  - "error": Check failure fails the overall run
  - "warning": Check failure is reported but doesn't fail the run

- **suggestion** (string): Message shown when check fails
  - Supports variable interpolation: {{.varname}}
  - Can reference grok-extracted values

- **requires** (array of strings): IDs of checks that must pass first
  - Creates dependency ordering
  - Referenced checks must exist
  - No circular dependencies allowed

- **timeout** (string): Maximum time for check execution
  - Format: Go duration string (e.g., "30s", "5m", "1h")
  - Default: "30s"

- **file** (string): File to read output from instead of command stdout
  - Useful for reading generated reports

### Example Complete Check
```yaml
- id: coverage
  run: go test -cover ./... 2>&1 | tail -1
  grok:
    - "coverage: %{NUMBER:coverage}%"
  assert: "coverage >= 70"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is 70%. Add more tests."
  requires:
    - test
  timeout: 60s
```

## Dependency Validation Rules

The **requires** field creates dependencies between checks. These rules must be followed:

### Basic Rules
1. All IDs in **requires** must reference existing checks
2. A check cannot require itself (no self-reference)
3. No circular dependencies allowed

### Circular Dependency Detection
A circular dependency exists when check A requires B, and B (directly or indirectly) requires A.

Examples of INVALID circular dependencies:
```yaml
# Direct cycle: A -> B -> A
checks:
  - id: A
    run: echo A
    requires: [B]
  - id: B
    run: echo B
    requires: [A]

# Indirect cycle: A -> B -> C -> A
checks:
  - id: A
    run: echo A
    requires: [C]
  - id: B
    run: echo B
    requires: [A]
  - id: C
    run: echo C
    requires: [B]
```

### Valid Dependency Patterns
```yaml
# Linear chain: fmt -> vet -> lint -> test
checks:
  - id: fmt
    run: gofmt -l .
  - id: vet
    run: go vet ./...
    requires: [fmt]
  - id: lint
    run: golangci-lint run
    requires: [vet]
  - id: test
    run: go test ./...
    requires: [lint]

# Diamond pattern: A <- B, A <- C, B <- D, C <- D
checks:
  - id: A
    run: echo A
  - id: B
    run: echo B
    requires: [A]
  - id: C
    run: echo C
    requires: [A]
  - id: D
    run: echo D
    requires: [B, C]
```

### Execution Order
- Checks with no dependencies run first (potentially in parallel)
- A check only runs after all its required checks pass
- If a required check fails, dependent checks are skipped

## Variable Interpolation Rules

Variables defined in **vars** can be referenced throughout the configuration.

### Syntax
- Variables are referenced using Go template syntax: {{.varname}}
- The variable name must match exactly (case-sensitive)
- Variables must be defined in the **vars** section before use

### Where Variables Can Be Used
Variables can be interpolated in these fields:
- **run**: Command to execute
- **assert**: Assertion expression
- **suggestion**: Failure message
- **file**: Output file path
- **grok**: Pattern strings (array elements)

### Example Usage
```yaml
version: "1"

vars:
  packages: "./cmd/... ./internal/... ./pkg/..."
  coverage_min: "70"
  test_timeout: "5m"

checks:
  - id: test
    run: go test {{.packages}} -timeout {{.test_timeout}}
    severity: error
    suggestion: "Tests failed in {{.packages}}"

  - id: coverage
    run: go test {{.packages}} -cover 2>&1 | tail -1
    grok:
      - "coverage: %{NUMBER:coverage}%"
    assert: "coverage >= {{.coverage_min}}"
    suggestion: "Coverage {{.coverage}}% is below {{.coverage_min}}%"
```

### Variable Naming Rules
- Use alphanumeric characters and underscores
- Start with a letter or underscore
- Case-sensitive: {{.Packages}} != {{.packages}}
- Good names: packages, coverage_min, test_timeout
- Avoid: kebab-case ({{.test-timeout}}) or spaces

### Grok-Extracted Values
Values extracted by grok patterns are available in:
- **assert**: Reference extracted values directly (coverage >= 70)
- **suggestion**: Use template syntax ({{.coverage}})

Config vars take precedence over grok-extracted values if names conflict.

### Common Mistakes
1. Using undefined variables: {{.undefined_var}} will remain as literal text
2. Wrong syntax: {.var} or {{ .var }} (must be {{.var}})
3. Referencing grok values in **run** (grok runs after the command)

## Explicit DO NOT List

When generating a vibeguard.yaml configuration, DO NOT:

### YAML Structure
- DO NOT include YAML comments (# comment) in the generated config
- DO NOT add extra top-level keys beyond version, vars, checks
- DO NOT use YAML anchors and aliases (&anchor, *alias)
- DO NOT use multi-document YAML (--- separator)

### Check Definitions
- DO NOT create checks for tools not detected in the project
- DO NOT use empty strings for required fields (id, run)
- DO NOT use duplicate check IDs
- DO NOT reference undefined variables
- DO NOT create circular dependencies in requires

### Assertions and Grok
- DO NOT use invalid assertion syntax (missing operators, invalid comparisons)
- DO NOT reference grok-extracted values that don't exist
- DO NOT use unsupported assertion operators
- DO NOT write grok patterns that don't match Grok syntax

### Commands
- DO NOT assume tools are installed without detection evidence
- DO NOT use interactive commands (require user input)
- DO NOT use commands that modify system state destructively
- DO NOT hardcode paths that are project-specific without using variables

### Severity and Timeouts
- DO NOT use severity values other than "error" or "warning"
- DO NOT use invalid timeout formats (use Go duration: "30s", "5m", not "30 seconds")
- DO NOT set unreasonably short timeouts that would cause false failures

### Variable Interpolation
- DO NOT use variables before defining them in vars
- DO NOT mix up variable syntax (use {{.var}}, not $var or ${var})
- DO NOT name variables with special characters or spaces

### General Guidelines
- DO NOT generate configs that would fail vibeguard validate
- DO NOT assume the execution environment has specific tools without verification
- DO NOT create overly complex configurations when simple ones suffice
- DO NOT add checks that provide no value for the project type

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

After generating the configuration:

1. **Validate the YAML syntax and schema:**
   ```bash
   vibeguard validate
   ```
   This verifies the configuration file has correct YAML syntax and adheres to the vibeguard schema.

2. **Run the checks to verify they execute properly:**
   ```bash
   vibeguard check
   ```
   This runs all defined checks and ensures they execute successfully. Fix any failing checks before considering the task complete.

Only consider this task complete when both commands pass without errors.