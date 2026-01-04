# VibeGuard Architecture

This document describes the high-level architecture of VibeGuard, its core components, and the design principles that guide development.

## Table of Contents

1. [Overview](#overview)
2. [Core Components](#core-components)
3. [Execution Model](#execution-model)
4. [Data Flow](#data-flow)
5. [Configuration System](#configuration-system)
6. [Pattern Matching and Assertions](#pattern-matching-and-assertions)
7. [Output Formats](#output-formats)
8. [Design Principles](#design-principles)

## Overview

VibeGuard is a lightweight, declarative quality gate automation tool built in Go. It enables teams to define and enforce policy checks through YAML configuration files, with minimal overhead and zero external dependencies at runtime.

**Key characteristics:**
- **Single binary deployment** - No runtime dependencies, works in any CI/CD environment
- **Declarative configuration** - Policy checks defined in simple YAML files
- **Composable and modular** - Checks can depend on each other with transparent DAG execution
- **Fast and efficient** - Built in Go for performance; runs in seconds
- **Tool-agnostic** - Works with any shell command, language, or tool
- **Rich output** - Supports both human-readable and machine-readable formats

## Core Components

### 1. Command-Line Interface (`internal/cli/`)

The CLI layer provides user-facing commands built with [Cobra](https://github.com/spf13/cobra).

**Available commands:**

| Command | Purpose |
|---------|---------|
| `vibeguard check [id]` | Run all checks or a specific check by ID |
| `vibeguard init [--assist]` | Initialize a new VibeGuard config with optional AI assistance |
| `vibeguard list` | Display all configured checks and their dependencies |
| `vibeguard validate` | Validate the config file without running checks |
| `vibeguard --version` | Show version information |

**Global flags:**

| Flag | Short | Default | Purpose |
|------|-------|---------|---------|
| `--config` | `-c` | Auto-detect | Path to config file |
| `--verbose` | `-v` | false | Show all check results (not just failures) |
| `--json` | — | false | Output in JSON format (to stderr) |
| `--parallel` | `-p` | 4 | Maximum concurrent checks |
| `--fail-fast` | — | false | Stop on first error-severity violation |
| `--log-dir` | — | `.vibeguard/log` | Directory for check execution logs |

**Exit codes:**

| Code | Meaning | Usage |
|------|---------|-------|
| 0 | Success | All checks passed or no violations detected |
| 2 | Config error | Invalid YAML, validation failure, or missing file |
| 3 | Violation | Error-severity check failed |
| 4 | Timeout | Check execution timeout or command not found |

### 2. Configuration System (`internal/config/`)

Manages loading, parsing, and validation of VibeGuard configuration files.

**Configuration structure:**

```yaml
version: "1"              # Config version (currently "1")
vars:                     # Optional variables for templating
  key: value
checks:                   # List of checks to execute
  - id: check-id          # Unique identifier
    run: shell command    # Command to execute
    grok: [patterns]      # Optional pattern extraction
    assert: expression    # Optional assertion (requires grok)
    severity: error       # error | warning
    suggestion: text      # Help text on failure
    fix: command          # Suggested fix command
    requires: [ids]       # Depend on other checks
    timeout: 30s          # Execution timeout
    file: path            # Read output from file instead of stdout
```

**Key features:**

- **Variable interpolation** - Reference variables with `{{.var_name}}` syntax
- **Auto-discovery** - Searches for `vibeguard.yaml`, `vibeguard.yml`, `.vibeguard.yaml`, `.vibeguard.yml`
- **Type validation** - Strict parsing with helpful error messages
- **Duration parsing** - Human-readable timeouts (e.g., "5s", "1m", "30s")

### 3. Executor (`internal/executor/`)

Handles low-level command execution with timeout and context management.

**Responsibilities:**
- Execute shell commands in subprocess
- Capture stdout, stderr, and exit codes
- Enforce per-check timeout limits
- Context-aware cancellation for parallel execution
- Return structured execution results

**Key types:**

```go
type Result struct {
    ExitCode   int
    Stdout     string
    Stderr     string
    Duration   time.Duration
    TimedOut   bool
}
```

### 4. Orchestrator (`internal/orchestrator/`)

Manages check dependency resolution, DAG construction, and parallel execution.

**Responsibilities:**
- Build directed acyclic graph (DAG) from check dependencies
- Topological sorting using Kahn's algorithm
- Level-based execution (checks at same level run in parallel)
- Concurrency control via semaphore
- Fail-fast mode to stop on first error
- Logging of check outputs to `.vibeguard/log/`

**Execution flow:**

```
Parse config
    ↓
Build DAG from "requires" declarations
    ↓
Validate no circular dependencies
    ↓
Topological sort → execution levels
    ↓
For each level (parallel with semaphore):
  - Execute all checks in level
  - Collect results
  - If fail-fast and error found: stop
    ↓
Aggregate results
    ↓
Format and output results
```

### 5. Pattern Matching (`internal/grok/`)

Implements pattern extraction from command output using Grok syntax.

**Features:**
- **Grok patterns** - Pre-built patterns like `%{NUMBER:name}`, `%{WORD:name}`
- **Custom regex** - Named capture groups with `(?P<name>regex)`
- **Named captures** - Extract structured data from unstructured output
- **Elastic grok library** - Uses [elastic/go-grok](https://github.com/elastic/go-grok)

**Examples:**

```yaml
# Extract test coverage from output
grok:
  - 'total:\s+%{NUMBER:coverage}%'

# Extract version number
grok:
  - 'version:\s+%{DATA:version}'

# Multiple patterns to extract different values
grok:
  - 'passed:\s+%{NUMBER:passed}'
  - 'failed:\s+%{NUMBER:failed}'
```

### 6. Assertions (`internal/assert/`)

Evaluates boolean expressions against extracted data.

**Supported operators:**

| Category | Operators | Example |
|----------|-----------|---------|
| Comparison | `>=`, `>`, `<=`, `<`, `==`, `!=` | `coverage >= 75` |
| Logical | `&&`, `\|\|`, `!` | `passed > 0 && failed == 0` |
| Arithmetic | `+`, `-`, `*`, `/` | `(passed + failed) >= 10` |
| Literals | Numbers, strings, booleans | `true`, `false`, `"text"`, `42` |

**Precedence** (highest to lowest):
1. Unary operators (`!`, `-`)
2. Arithmetic (`*`, `/`)
3. Arithmetic (`+`, `-`)
4. Comparison (`>=`, `>`, `<=`, `<`, `==`, `!=`)
5. Logical AND (`&&`)
6. Logical OR (`||`)

**Examples:**

```yaml
assert: "coverage >= 75"
assert: "tests > 0 && failures == 0"
assert: "version != '0.0.0'"
```

### 7. Output Formatting (`internal/output/`)

Produces human-readable and machine-readable output.

**Text output modes:**

- **Quiet mode (default)** - Only violations shown ("silence is success")
  ```
  ✗ lint: ESLint found style violations
    Suggestion: Run 'npm run lint:fix'

  Exit code: 3
  ```

- **Verbose mode** - All check results with status and timing
  ```
  ✓ fmt (5ms)
  ✓ vet (156ms)
  ✗ test (2.3s)
    Suggestion: Fix failing tests
  ⊘ coverage (cancelled - test failed)

  Exit code: 3
  ```

**JSON output format:**

Outputs structured results to stderr for CI/CD integration:

```json
{
  "checks": [
    {
      "id": "test",
      "status": "failed",
      "exit_code": 1,
      "duration_ms": 2300,
      "suggestion": "Fix failing tests"
    }
  ],
  "violations": ["test"],
  "exit_code": 3,
  "fail_fast_triggered": false
}
```

## Execution Model

### Phase 1: Configuration Loading

1. Search for config file (if not specified via `-c`)
2. Parse YAML and validate schema
3. Perform variable interpolation
4. Create Check objects

### Phase 2: DAG Construction and Validation

1. Build dependency graph from `requires` declarations
2. Check for circular dependencies
3. Topologically sort checks using Kahn's algorithm
4. Group checks into execution levels

### Phase 3: Check Execution

1. **Level-based execution**: Each level runs in parallel
2. **Concurrency control**: Limited by `--parallel` flag (default: 4)
3. **Timeout enforcement**: Each check has individual timeout
4. **Output logging**: Results written to `.vibeguard/log/`
5. **Fail-fast mode**: If error found and `--fail-fast` enabled, remaining levels skipped

### Phase 4: Result Processing

1. Execute Grok patterns on command output
2. Evaluate assertions against extracted values
3. Determine pass/fail based on:
   - Exit code (if no grok/assert)
   - Assertion result (if grok/assert present)
4. Filter results based on severity
5. Aggregate violations

### Phase 5: Output and Exit

1. Format results (text or JSON)
2. Output to stdout/stderr
3. Exit with appropriate code (0, 2, 3, or 4)

## Data Flow

```
┌──────────────────────────────┐
│ Config File (vibeguard.yaml) │
└──────────────┬───────────────┘
               │
               ↓
        ┌─────────────┐
        │Config Parser│
        └──────┬──────┘
               │
               ↓
        ┌─────────────────┐
        │Check Objects    │
        │+ Dependencies   │
        └────────┬────────┘
                 │
                 ↓
           ┌──────────────┐
           │Orchestrator  │
           │- DAG Build   │
           │- Topsorting  │
           └────┬─────────┘
                │
                ↓
       ┌─────────────────┐
       │Level Executor   │
       │(Parallel Run)   │
       └────┬────────────┘
            │
    ┌───────┴───────┐
    ↓               ↓
 ┌────────┐    ┌──────────┐
 │Executor│ -> │Check Log │
 └────┬───┘    └──────────┘
      │
      ↓
 ┌─────────────┐
 │Result Data  │
 │- Output     │
 │- Exit Code  │
 │- Duration   │
 └──────┬──────┘
        │
        ↓
 ┌──────────────────┐
 │Pattern Matching  │
 │(Grok)            │
 └────────┬─────────┘
          │
          ↓
 ┌──────────────────┐
 │Assertion Eval    │
 │(Assert)          │
 └────────┬─────────┘
          │
          ↓
 ┌──────────────────┐
 │Result Processor  │
 │- Severity Filter │
 │- Violation Check │
 └────────┬─────────┘
          │
          ↓
 ┌──────────────────┐
 │Output Formatter  │
 │(Text/JSON)       │
 └────────┬─────────┘
          │
          ↓
      stdout/stderr
```

## Configuration System

### Variable Interpolation

Variables provide reusable values across checks:

```yaml
version: "1"

vars:
  go_packages: "./..."
  min_coverage: "70"
  lint_tool: "golangci-lint"

checks:
  - id: test
    run: go test {{.go_packages}}
    timeout: {{.timeout | default "30s"}}

  - id: coverage
    run: go test {{.go_packages}} -coverprofile=cover.out
    assert: "coverage >= {{.min_coverage}}"
```

### Check Inheritance Pattern

While VibeGuard doesn't have inheritance, you can reuse patterns:

```yaml
vars:
  timeout_fast: "5s"
  timeout_slow: "60s"

checks:
  - id: fmt
    run: gofmt -l .
    severity: error
    timeout: {{.timeout_fast}}

  - id: vet
    run: go vet ./...
    severity: error
    timeout: {{.timeout_fast}}

  - id: test
    run: go test ./...
    severity: error
    timeout: {{.timeout_slow}}
    requires: [fmt, vet]
```

## Pattern Matching and Assertions

### Grok Pattern Examples

Extract coverage percentage from test output:

```yaml
checks:
  - id: coverage
    run: go test ./... -coverprofile=cover.out && go tool cover -func=cover.out
    grok:
      - 'total:\s+coverage:\s+%{NUMBER:coverage}%'
    assert: "coverage >= 70"
```

Extract multiple values:

```yaml
checks:
  - id: test-summary
    run: npm test -- --json --outputFile=test-results.json && cat test-results.json
    grok:
      - 'numPassedTests[^:]*:\s*%{NUMBER:passed}'
      - 'numFailedTests[^:]*:\s*%{NUMBER:failed}'
    assert: "failed == 0 && passed > 0"
```

### Custom Regex Patterns

Use named capture groups for custom regex:

```yaml
grok:
  - 'version:\s+(?P<version>\d+\.\d+\.\d+)'
  - 'memory:\s+(?P<memory>\d+(?:\.\d+)?)\s*(?P<unit>MB|GB)'
```

## Output Formats

### Text Output (Default)

```
$ vibeguard check

✗ fmt: Code formatting check failed
  Suggestion: Run 'gofmt -w .' to format your code
  Output: cmd/main.go:12:5

Exit code: 3
```

### Text Output (Verbose)

```
$ vibeguard check -v

✓ fmt (12ms)
✓ vet (234ms)
✗ test (3.5s)
  Suggestion: Fix failing tests
  Exit code: 1
⊘ coverage (skipped - test failed)

Exit code: 3
```

### JSON Output

```
$ vibeguard check --json

{
  "checks": [
    {
      "id": "fmt",
      "status": "passed",
      "exit_code": 0,
      "duration_ms": 12
    },
    {
      "id": "test",
      "status": "failed",
      "exit_code": 1,
      "duration_ms": 3500,
      "suggestion": "Fix failing tests"
    }
  ],
  "violations": ["test"],
  "exit_code": 3,
  "fail_fast_triggered": false
}
```

## Design Principles

### 1. Single Responsibility

Each component has one clear purpose:
- **Config** - Parse and validate
- **Executor** - Execute commands
- **Orchestrator** - Manage dependencies
- **Output** - Format results

### 2. Separation of Concerns

- Command execution logic is isolated from result processing
- Pattern matching is separate from assertion evaluation
- Configuration is decoupled from execution

### 3. Declarative Configuration

Users declare "what" they want to check, not "how" to check it. Implementation details are hidden.

### 4. Minimal Overhead

- Single Go binary with no external dependencies at runtime
- Fast startup and check execution
- Efficient memory usage for CI/CD pipelines

### 5. Tool Agnostic

VibeGuard doesn't care what tools or languages are used:
- Works with any shell command
- Supports any language or tool that can produce output
- Easy integration with existing workflows

### 6. Clear Error Messages

- Helpful suggestions when checks fail
- Specific exit codes for different failure modes
- Logged output available in `.vibeguard/log/` for debugging

### 7. Composability

- Checks can depend on each other
- Dependencies are transparent and validated
- Parallel execution when possible
- Clean failure handling with fail-fast option

## Deployment Model

VibeGuard is designed for CI/CD deployment:

1. **No installation needed** - Single binary
2. **No configuration management** - YAML files in repository
3. **CI-friendly** - Works with GitHub Actions, GitLab CI, Jenkins, etc.
4. **Container-ready** - Can be installed in Docker images
5. **Pre-commit hook compatible** - Integrates with git hooks

## Future Architecture Considerations

See Architecture Decision Records in `docs/adr/` for detailed discussion of:
- Code quality standards (ADR-004)
- Mutation testing integration (ADR-007)
- Pre-commit hook integration (ADR-006)
- Policy enforcement dogfooding (ADR-005)
