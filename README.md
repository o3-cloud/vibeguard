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
| `id` | Yes (per check) | string | Unique check identifier | — |
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

VibeGuard is released under [LICENSE_NAME] (see LICENSE file for details).

## Further Reading

- [Conventional Commits](https://www.conventionalcommits.org/) — Commit message specification
- [Go Best Practices](https://go.dev/doc/effective_go) — Go coding guidelines
- [Project Conventions](./CONVENTIONS.md) — VibeGuard-specific standards

## Support

For issues, questions, or suggestions:

- Open an issue on GitHub
- Check existing ADRs and documentation
- Review the implementation patterns in `docs/patterns/`
