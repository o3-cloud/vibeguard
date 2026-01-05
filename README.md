# VibeGuard

**VibeGuard** is a lightweight, composable policy enforcement system designed for seamless integration with CI/CD pipelines, agent loops, and Cloud Code workflows.

## Overview

VibeGuard enforces policies at scale with minimal overhead. It combines declarative policy definition, flexible runner patterns, and LLM judge integration to provide intelligent policy evaluation for modern development workflows.

### Key Features

- **Single Binary Deployment** — Compiled Go binary with zero runtime dependencies
- **Fast & Lightweight** — Minimal resource usage and quick startup time for frequent invocation
- **Declarative Policies** — Define policies in YAML for clarity and maintainability
- **Multiple Runner Patterns** — Supports various policy evaluation approaches
- **Judge Integration** — Leverage LLMs for nuanced policy evaluation
- **Cross-Platform** — Runs seamlessly on Linux, macOS, and Windows

## Quick Start

### Installation

Download a pre-built binary from [releases](https://github.com/vibeguard/vibeguard/releases), or build from source:

```bash
git clone https://github.com/vibeguard/vibeguard.git
cd vibeguard
go build -o vibeguard ./cmd/vibeguard
```

### Basic Usage

Initialize a configuration file:

```bash
vibeguard init
```

Run all checks:

```bash
vibeguard check
```

Run a specific check:

```bash
vibeguard check fmt
```

### Running Tests

```bash
go test -v ./...
```

### AI-Assisted Setup

VibeGuard includes an AI-assisted setup feature that helps AI coding agents (Claude Code, Cursor, etc.) generate project-specific configurations automatically.

```bash
vibeguard init --assist
```

This command analyzes your project and generates a comprehensive setup guide that AI agents can use to create a valid `vibeguard.yaml` configuration. The guide includes:

- Detected project type (Go, Node.js, Python, Rust, Ruby, Java)
- Existing tools and their configuration files
- Recommended checks based on detected tools
- Project structure analysis
- Configuration syntax and validation rules

**Usage with Claude Code:**

1. Run `vibeguard init --assist` in your project directory
2. Copy the generated prompt to Claude Code
3. The AI agent will generate a customized `vibeguard.yaml`
4. Review and commit the configuration

For detailed documentation, see [AI-Assisted Setup Guide](docs/ai-assisted-setup.md).

## CLI Reference

### Global Flags

All commands support the following flags:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to config file | Searches for `vibeguard.yaml`, `vibeguard.yml`, `.vibeguard.yaml`, `.vibeguard.yml` |
| `--fail-fast` | | Stop on first failure | false |
| `--json` | | Output results in JSON format | false |
| `--parallel` | `-p` | Max parallel checks to run | 4 |
| `--verbose` | `-v` | Show all check results, not just failures | false |

### Commands

#### `vibeguard check [id]`

Run all configured checks, or a specific check by ID.

```bash
vibeguard check              # Run all checks
vibeguard check fmt          # Run only the 'fmt' check
vibeguard check -v           # Run all checks with verbose output
vibeguard check --fail-fast  # Stop on first failure
vibeguard check --json       # Output results in JSON format
```

For JSON output format details, see [JSON Output Schema](docs/JSON-OUTPUT-SCHEMA.md).

#### `vibeguard init [flags]`

Create a starter configuration file in the current directory.

```bash
vibeguard init              # Create vibeguard.yaml if it doesn't exist
vibeguard init --force      # Overwrite existing configuration file
vibeguard init --assist     # Generate AI setup guide for this project
vibeguard init --assist -o guide.txt  # Save guide to a file
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--force` | `-f` | Overwrite existing config file | false |
| `--assist` | | Generate AI agent setup guide | false |
| `--output` | `-o` | Output file for --assist (default: stdout) | stdout |

#### `vibeguard list`

List all checks defined in the configuration file, showing IDs, commands, and dependencies.

```bash
vibeguard list              # Show all checks
vibeguard list --json       # Show checks in JSON format
```

#### `vibeguard validate`

Validate the configuration file without running any checks. Useful for catching errors before execution.

```bash
vibeguard validate          # Validate the default config file
vibeguard validate -c prod.yaml  # Validate a specific config file
```

## Exit Codes

VibeGuard uses the following exit codes to indicate the result of check execution. This is particularly useful for CI/CD integration and automated workflows where exit codes determine the success or failure of a step.

| Exit Code | Name | Description |
|-----------|------|-------------|
| 0 | Success | All checks passed successfully |
| 2 | ConfigError | Configuration file error (invalid YAML, validation failure, etc.) |
| 3 | Violation | One or more error-severity violations detected during execution |
| 4 | Timeout | Check execution error (timeout exceeded, command not found, etc.) |

### CI/CD Integration

When integrating VibeGuard into CI/CD pipelines:

- **Exit code 0** — Pipeline can proceed
- **Exit codes ≥ 2** — Pipeline is blocked (suitable for pre-commit hooks and CI checks)

Example with GitHub Actions:

```yaml
- name: Run VibeGuard checks
  run: vibeguard check
  # Automatically fails the step if exit code >= 1
```

## Configuration Schema

VibeGuard configurations are written in YAML. Here's the complete schema:

```yaml
# Configuration version (currently "1")
version: "1"

# Optional: Global variables for interpolation in check commands
vars:
  go_packages: "./..."
  python_version: "3.11"

# List of checks to execute
checks:
  - id: check-name           # Unique check identifier
    run: command to execute  # Shell command to run (variables interpolated with {{.var_name}})

    # Optional: Extract data from command output using grok patterns
    grok:
      - pattern_string
      # Multiple patterns can be specified as a list

    # Optional: Read output from file instead of command stdout
    file: path/to/output.txt

    # Optional: Assert extracted data meets conditions
    assert: "condition"      # e.g., "coverage >= 80" or "result == 'ok'"

    # Optional: Severity level when check fails
    severity: error          # Options: "error", "warning" (default: "error")

    # Optional: Actionable suggestion when check fails
    suggestion: "How to fix this..."

    # Optional: List of check IDs that must pass before this check runs
    requires:
      - other-check-id

    # Optional: Timeout for this check
    timeout: 30s             # Examples: "5s", "1m", "30s" (default: 30s)
```

### Field Details

| Field | Required | Type | Description | Default |
|-------|----------|------|-------------|---------|
| `version` | Yes | string | Config format version | — |
| `vars` | No | map[string]string | Global variables for interpolation | — |
| `checks` | Yes | array | List of checks to run | — |
| `id` | Yes (per check) | string | Unique check identifier. Must start with a letter or underscore, followed by alphanumeric characters, underscores, or hyphens (regex: `^[a-zA-Z_][a-zA-Z0-9_-]*$`) | — |
| `run` | Yes (per check) | string | Shell command with optional `{{.var}}` interpolation | — |
| `grok` | No | array[string] | Grok patterns to extract data from command output | — |
| `file` | No | string | File path to read output from instead of command stdout | — |
| `assert` | No | string | Assertion expression (requires `grok` patterns) | — |
| `severity` | No | string | `error` or `warning` | `error` |
| `suggestion` | No | string | Help text shown when check fails | — |
| `requires` | No | array[string] | Check IDs that must pass first | — |
| `timeout` | No | duration | Max execution time (e.g., `5s`, `1m`) | `30s` |

### Variable Interpolation

Use `{{.variable_name}}` syntax to interpolate global variables into check commands:

```yaml
vars:
  go_packages: "./..."

checks:
  - id: vet
    run: go vet {{.go_packages}}
```

### Grok Pattern Extraction

Extract structured data from command output using grok patterns:

```yaml
checks:
  - id: test-coverage
    run: go test ./... -coverprofile=cover.out && go tool cover -func=cover.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= 80"
    suggestion: "Coverage is below 80%. Run 'go test ./...' with coverage analysis."
```

### Assertion Expression Operators

The `assert` field supports a rich set of operators for flexible condition evaluation:

#### Comparison Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `>=` | Greater than or equal | `coverage >= 80` |
| `>` | Greater than | `score > 90` |
| `<=` | Less than or equal | `latency <= 500` |
| `<` | Less than | `errors < 5` |
| `==` | Equal (numeric or string) | `status == "ok"` or `count == 42` |
| `!=` | Not equal | `result != "fail"` |

#### Logical Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `&&` | Logical AND | `coverage >= 80 && tests_passed == true` |
| `\|\|` | Logical OR | `linting_ok == true \|\| warnings < 10` |
| `!` | Logical NOT | `!failed` |

#### Arithmetic Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `count + 5 >= 10` |
| `-` | Subtraction | `total - errors > 50` |
| `*` | Multiplication | `ratio * 100 >= 80` |
| `/` | Division | `usage / 1024 < 100` |

#### Literals and Values
| Type | Syntax | Example |
|------|--------|---------|
| Numbers | Integer or float | `coverage >= 80` or `ratio > 0.95` |
| Strings | Single or double quoted | `status == "ok"` or `result == 'pass'` |
| Booleans | `true` or `false` | `tests_passed == true` |
| Variables | Grok pattern names | `coverage`, `result`, `errors` |
| Grouping | Parentheses | `(coverage >= 80) && (tests_passed == true)` |

#### Examples

```yaml
checks:
  # Numeric comparison
  - id: coverage-threshold
    run: go test ./... -coverprofile=cover.out && go tool cover -func=cover.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= 80"

  # String comparison
  - id: lint-status
    run: golangci-lint run
    grok:
      - 'status:\s+%{WORD:result}'
    assert: "result == 'pass'"

  # Logical AND
  - id: quality-gates
    run: ./run-checks.sh
    grok:
      - 'coverage:%{NUMBER:coverage}'
      - 'tests:%{WORD:tests_status}'
    assert: "coverage >= 75 && tests_status == 'pass'"

  # Logical OR (fail-safe)
  - id: either-metric
    run: ./check-metrics.sh
    grok:
      - 'metric_a:%{NUMBER:a}'
      - 'metric_b:%{NUMBER:b}'
    assert: "a > 100 || b > 50"

  # Logical NOT
  - id: no-failures
    run: ./test-runner.sh
    grok:
      - 'failures:%{NUMBER:failed}'
    assert: "!(failed > 0)"

  # Arithmetic operations
  - id: normalized-score
    run: ./calculate-score.sh
    grok:
      - 'points:%{NUMBER:points}'
    assert: "(points * 100) / 50 >= 80"

  # Complex expression
  - id: multi-condition
    run: ./full-check.sh
    grok:
      - 'coverage:%{NUMBER:coverage}'
      - 'performance:%{NUMBER:perf}'
      - 'security:%{WORD:sec_status}'
    assert: "(coverage >= 80 && perf < 1000) || sec_status == 'pass'"
```

### Reading Output from Files

The `file` field allows reading check output from a file instead of command stdout. This is useful when tools write results to files (e.g., coverage reports, test result files) rather than printing to stdout:

```yaml
checks:
  - id: coverage
    run: go test ./... -coverprofile=coverage.out
    file: coverage.out
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
    assert: "coverage >= 80"
    suggestion: "Coverage is {{.coverage}}%, target is 80%. Add more tests."
```

When `file` is specified, VibeGuard reads the file contents and applies grok patterns and assertions to that content instead of the command's stdout. The command still runs normally—the `file` field simply changes where the output is read from.

### Grok Pattern Debugging Guide

When a grok pattern fails to match, VibeGuard provides detailed error messages to help you debug. Understanding these messages and common pattern syntax is essential for effective pattern configuration.

#### Error Format and Interpretation

When a grok pattern fails to parse output, you'll see an error like:

```
grok pattern 0 failed to parse
  pattern: 'coverage: %{NUMBER:coverage}%'
  output: 'Total coverage: 85.5%'
  error: <underlying grok error>
```

The error includes:
- **pattern index** (0, 1, 2, etc.) - Which pattern in your `grok` list failed
- **pattern string** - The exact pattern that failed
- **output** - The first 100 characters of the text being matched (truncated for readability)
- **error** - The underlying grok compilation or matching error

#### Common Pattern Syntax

Grok supports built-in patterns for common data types. Here are the most frequently used:

| Pattern | Matches | Example |
|---------|---------|---------|
| `%{NUMBER:name}` | Integer or decimal numbers | `42`, `3.14`, `85.5` |
| `%{INT:name}` | Integers only | `42`, `100`, `-5` |
| `%{WORD:name}` | Single words (letters, digits, underscore) | `PASS`, `coverage`, `test_1` |
| `%{IP:name}` | IPv4 addresses | `192.168.1.1`, `10.0.0.1` |
| `%{IPV6:name}` | IPv6 addresses | `::1`, `2001:db8::1` |
| `%{UUID:name}` | UUID format | `550e8400-e29b-41d4-a716-446655440000` |
| `%{GREEDYDATA:name}` | Any characters (greedy) | Useful for capturing everything to end of line |
| `%{DATA:name}` | Non-greedy data capture | Stops at first match of following pattern |

#### Mixing Built-in and Custom Patterns

You can mix grok built-in patterns with custom regex:

```yaml
checks:
  - id: test-results
    run: ./run-tests.sh
    grok:
      # Built-in pattern
      - '%{NUMBER:tests} tests'
      # Custom regex with named capture group
      - 'passed: (?P<passed>[0-9]+)'
      # Mix both styles (grok built-in + literal text)
      - 'Failures: %{NUMBER:failures}'
```

#### Pattern Matching Behavior

- **Patterns are applied sequentially** - If you specify multiple patterns, each is applied to the output independently
- **Later patterns override earlier values** - If two patterns capture the same field name, the later pattern's value wins
- **Non-matches return empty fields** - If a pattern doesn't match, the fields it would capture simply won't be present (don't generate errors)
- **All patterns can be optional** - You can have patterns that may or may not match; only those that match contribute extracted values

#### Common Debugging Strategies

**1. Test patterns incrementally**

Start with simple patterns and add complexity:

```yaml
grok:
  # Start here - does this basic pattern work?
  - '%{NUMBER:coverage}'
  # Then add context
  - 'coverage: %{NUMBER:coverage}%'
  # Then handle variations
  - 'Total coverage: %{NUMBER:coverage}%'
```

**2. Account for special characters**

Special regex characters need escaping:

```yaml
grok:
  # Literal parentheses must be escaped
  - 'total:.*\(statements\)\s+%{NUMBER:coverage}%'

  # Literal brackets
  - '\[%{WORD:level}\]\s+%{GREEDYDATA:message}'
```

**3. Use capturing groups for flexible matching**

When built-in patterns don't fit, use regex:

```yaml
grok:
  # Capture test count from various formats
  - '(?P<tests>[0-9]+) tests?'  # matches "1 test" or "42 tests"

  # Capture version numbers
  - 'version[: ]+(?P<version>[0-9.]+)'

  # Capture quoted strings
  - 'error: "(?P<message>[^"]+)"'
```

**4. Handle whitespace variations**

Use `\s+` for flexible whitespace:

```yaml
grok:
  # Matches "coverage:85.5%" or "coverage: 85.5 %"
  - 'coverage\s*:\s*%{NUMBER:coverage}\s*%'
```

#### Pattern Examples by Use Case

**Coverage Extraction**
```yaml
grok:
  - 'coverage: %{NUMBER:coverage}%'
  - 'total:.*\(statements\)\s+%{NUMBER:coverage}%'  # Go test format
  - '%{NUMBER:coverage}%\s+coverage'  # Alternative order
```

**Test Count Extraction**
```yaml
grok:
  - '%{NUMBER:tests} tests?'
  - 'ran\s+%{NUMBER:tests}\s+tests?'
  - '(?P<passed>[0-9]+)/(?P<tests>[0-9]+) passed'
```

**Status/Result Extraction**
```yaml
grok:
  - 'status:\s+%{WORD:status}'
  - '%{WORD:result}\s+(PASS|FAIL|OK)'
  - 'result[: =]+(?P<result>\w+)'
```

**Error/Warning Counts**
```yaml
grok:
  - '%{NUMBER:errors} errors?'
  - 'errors:\s*%{INT:errors}\s*warnings:\s*%{INT:warnings}'
  - 'failed:\s+%{NUMBER:failures}'
```

**Duration/Performance**
```yaml
grok:
  - 'took %{NUMBER:duration}s'
  - 'completed in %{NUMBER:duration}ms'
  - 'latency:\s*%{NUMBER:latency}(ms|s)'
```

### Check Dependencies

Specify that a check requires other checks to pass first:

```yaml
checks:
  - id: vet
    run: go vet ./...

  - id: build
    run: go build ./...
    requires:
      - vet  # build only runs if vet passes
```

## Execution Model

VibeGuard uses a sophisticated execution model to efficiently run checks while respecting dependencies and resource constraints.

### Dependency Graph and Topological Ordering

Checks form a directed acyclic graph (DAG) based on their `requires` declarations. VibeGuard builds this graph and uses **Kahn's algorithm** for topological sorting to determine the execution order:

1. **Circular dependency detection** — The system validates that no circular dependencies exist before execution begins
2. **Level-based execution** — Checks are organized into levels where each level contains checks with no unprocessed dependencies
3. **Deterministic ordering** — The execution order is consistent and predictable across runs

For example, with this configuration:

```yaml
checks:
  - id: vet
    run: go vet ./...

  - id: fmt
    run: go fmt ./...

  - id: build
    run: go build ./...
    requires:
      - vet
      - fmt

  - id: test
    run: go test ./...
    requires:
      - build
```

The execution proceeds in **3 levels**:
- **Level 1** (parallel): `vet` and `fmt` run simultaneously
- **Level 2** (parallel): `build` runs after both `vet` and `fmt` complete
- **Level 3** (sequential): `test` runs after `build` completes

### Parallel Execution

Within each level, checks are executed **in parallel** to maximize efficiency:

- **`--parallel` flag** — Controls the maximum number of concurrent checks (default: 4)
  - `--parallel 1` — Run checks sequentially
  - `--parallel 8` — Allow up to 8 concurrent checks per level
  - Higher values increase throughput but consume more resources

Each check acquires a semaphore before execution. When the limit is reached, subsequent checks wait for earlier ones to complete before starting.

Example:
```bash
vibeguard check --parallel 8  # Increase concurrency for faster execution
```

### Fail-Fast Behavior

The `--fail-fast` flag stops execution on the first error-severity violation:

```bash
vibeguard check --fail-fast
```

**Behavior:**
- When an error-severity check fails, no further levels are executed
- In-flight checks (already started) in the current level continue to completion
- The exit code reflects the failure (exit code 3 for violations, 4 for timeouts)
- Useful in CI/CD pipelines where fast feedback on failures is important

**Example:**
```yaml
checks:
  - id: fmt
    run: go fmt ./...  # error severity (default)

  - id: vet
    run: go vet ./...  # error severity (default)

  - id: test
    run: go test ./...
    requires:
      - fmt
      - vet
```

With `--fail-fast`:
- If `fmt` fails, `vet` continues (same level)
- Both `fmt` and `vet` complete
- If either failed with error severity, `test` **does not run** (next level is skipped)

### Dependency Validation

Before a check executes, the orchestrator validates that all required dependencies have **passed**:

- **Passed dependency** — Required check passed all assertions and exited with success
- **Failed dependency** — Check skipped with reason "Skipped: required dependency failed"
- **No re-execution** — Dependencies are not re-run; each check executes exactly once

### Timeout Handling

Each check can have an individual timeout:

```yaml
checks:
  - id: integration-test
    run: ./run-integration-tests.sh
    timeout: 5m  # 5 minutes
```

**Timeout behavior:**
- Check execution is cancelled if it exceeds the timeout
- The check is marked as failed with `timedout: true`
- Timeouts return exit code 4 (takes precedence over error violations)
- Default timeout is 30 seconds if not specified

## Implementation Patterns

VibeGuard supports multiple implementation patterns:

1. **Declarative Policy Runner** — Evaluate policies defined in structured formats (YAML/JSON)
2. **Judge-as-a-Policy** — Use LLMs as policy evaluators
3. **Cloud Code Native** — Native integration with Anthropic Cloud Code
4. **Event-Based Policy Graph** — React to events with graph-based policy evaluation
5. **Git-Aware Guardrails** — Policies based on git history and code changes

(See `docs/patterns/` for detailed documentation on each pattern)

## Project Structure

```
vibeguard/
├── bin/                        # Compiled binary output
├── cmd/
│   └── vibeguard/              # Main CLI application
├── internal/
│   ├── assert/                 # Assertion expression parsing and evaluation
│   ├── cli/                    # Command-line interface (Cobra-based)
│   ├── config/                 # Configuration loading and validation
│   ├── executor/               # Check execution engine
│   ├── grok/                   # Grok pattern extraction and matching
│   ├── orchestrator/           # Check orchestration and dependency management
│   ├── output/                 # Output formatting (text, JSON)
│   └── version/                # Version information and constants
├── docs/
│   ├── adr/                    # Architecture Decision Records
│   ├── log/                    # Work logs and findings
│   └── sample-prompts/         # Sample prompts for AI-assisted setup
├── examples/                   # Example configurations
├── spikes/                     # Research and prototyping work
│   ├── config/                 # Configuration exploration
│   ├── executor/               # Executor prototype experiments
│   ├── opa/                    # OpenPolicyAgent integration exploration
│   └── orchestrator/           # Orchestration pattern research
├── README.md                   # This file
├── CONVENTIONS.md              # Code style and development standards
├── CLAUDE.md                   # Claude Code agent documentation
├── SECURITY.md                 # Security model and responsible disclosure
├── CHANGELOG.md                # Version history and changes
└── vibeguard.yaml              # Default configuration file
```

## Development

### Prerequisites

- **Go 1.21+** — Latest stable Go version
- **git** — For version control

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/vibeguard/vibeguard.git
cd vibeguard

# Install dependencies
go mod tidy

# Run tests
go test -v ./...

# Run linting
go fmt ./...
go vet ./...
```

### Code Style and Conventions

See `CONVENTIONS.md` for detailed code style guidelines, naming conventions, and development standards.

### Making Changes

1. Create a feature branch from `main`
2. Make your changes following the conventions in `CONVENTIONS.md`
3. Write or update tests as needed
4. Commit using Conventional Commits format (see ADR-002)
5. Push and create a pull request

## Architecture Decisions

Major architectural decisions are documented as Architecture Decision Records (ADRs) in `docs/adr/`:

- **ADR-001** — Adopt Beads for AI Agent Task Management
- **ADR-002** — Adopt Conventional Commits
- **ADR-003** — Adopt Go as the Primary Implementation Language
- **ADR-004** — Establish Code Quality Standards and Tooling
- **ADR-005** — Adopt VibeGuard for Policy Enforcement in CI/CD
- **ADR-006** — Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement
- **ADR-007** — Adopt Gremlins for Mutation Testing

Review these documents to understand the project's design rationale and constraints.

## Security

For security information, including our security model, threat boundaries, and responsible disclosure process, see [SECURITY.md](SECURITY.md).

## Contributing

Contributions are welcome! Please:

1. Read `CONVENTIONS.md` for code style requirements
2. Review existing ADRs to understand architectural constraints
3. Write tests for new functionality
4. Follow Conventional Commits for commit messages
5. Keep PRs focused and well-documented

## License

VibeGuard is released under the Apache License 2.0 (see LICENSE file for details).

## Further Reading

- [Conventional Commits](https://www.conventionalcommits.org/) — Commit message specification
- [Go Best Practices](https://go.dev/doc/effective_go) — Go coding guidelines
- [Project Conventions](./CONVENTIONS.md) — VibeGuard-specific standards

## Support

For help with VibeGuard:

- **Questions & Discussions:** [GitHub Discussions](https://github.com/vibeguard/vibeguard/discussions)
- **Bug Reports:** [GitHub Issues](https://github.com/vibeguard/vibeguard/issues)
- **Security Issues:** See [SECURITY.md](SECURITY.md)

For detailed support information including response time SLAs, see [SUPPORT.md](SUPPORT.md).
