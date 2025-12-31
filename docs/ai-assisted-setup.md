# AI-Assisted Setup Guide

VibeGuard's AI-assisted setup feature enables AI coding agents to automatically generate project-specific configurations. This guide explains how the feature works and how to use it effectively.

## Overview

The `--assist` flag for `vibeguard init` analyzes your project and generates a comprehensive setup guide that AI agents can understand. Instead of manually writing a `vibeguard.yaml` configuration, you can:

1. Run `vibeguard init --assist` to generate a setup prompt
2. Provide the prompt to an AI coding agent (Claude Code, Cursor, etc.)
3. The agent generates a customized configuration based on your project

This reduces VibeGuard onboarding from manual configuration to a few minutes of AI-assisted setup.

## Quick Start

### Basic Usage

```bash
# Generate setup guide to stdout
vibeguard init --assist

# Save to a file
vibeguard init --assist --output setup-guide.txt

# Show detection details
vibeguard init --assist --verbose
```

### Using with Claude Code

1. Navigate to your project directory
2. Run `vibeguard init --assist`
3. Copy the output (or pipe it directly to your clipboard)
4. Open Claude Code and paste the setup guide
5. Ask Claude to generate the configuration
6. Review the generated `vibeguard.yaml`
7. Run `vibeguard validate` to verify the configuration

## What Gets Detected

### Project Types

The inspector detects the following project types:

| Type | Detection Indicators | Confidence Factors |
|------|---------------------|-------------------|
| Go | `go.mod`, `go.sum`, `*.go` files | go.mod: 0.6, go.sum: 0.2, *.go: 0.2 |
| Node.js | `package.json`, `yarn.lock`, `pnpm-lock.yaml` | package.json: 0.6, lock files: 0.2 |
| Python | `pyproject.toml`, `setup.py`, `requirements.txt` | pyproject.toml: 0.6, setup.py: 0.3 |
| Rust | `Cargo.toml`, `Cargo.lock` | Cargo.toml: 0.7, Cargo.lock: 0.2 |
| Ruby | `Gemfile`, `*.gemspec` | Gemfile: 0.6, gemspec: 0.3 |
| Java | `pom.xml`, `build.gradle` | pom.xml: 0.7, build.gradle: 0.7 |

### Tools Detected

The inspector scans for common development tools:

**Go Tools:**
- golangci-lint (config: `.golangci.yml`, `.golangci.yaml`)
- gofmt (built-in)
- go vet (built-in)
- go test (built-in)
- goimports

**Node.js Tools:**
- ESLint (config: `.eslintrc.*`, `eslint.config.js`)
- Prettier (config: `.prettierrc.*`)
- Jest, Mocha, Vitest (test frameworks)
- TypeScript (config: `tsconfig.json`)
- npm audit

**Python Tools:**
- Black (config in `pyproject.toml`)
- Pylint (config: `.pylintrc`, `pyproject.toml`)
- pytest (config: `pytest.ini`, `pyproject.toml`)
- mypy (config: `mypy.ini`, `pyproject.toml`)
- Ruff, Flake8, isort, pip-audit

**CI/CD:**
- GitHub Actions (`.github/workflows/`)
- GitLab CI (`.gitlab-ci.yml`)
- CircleCI, Jenkins, Travis CI

**Git Hooks:**
- pre-commit (`.pre-commit-config.yaml`)
- husky (`.husky/`)
- lefthook (`lefthook.yml`)

### Project Structure

The inspector analyzes:
- **Source directories**: `src/`, `pkg/`, `lib/`, `internal/`
- **Test directories**: `tests/`, `__tests__/`, `test/`
- **Entry points**: `main.go`, `src/index.js`, `main.py`
- **Build output**: `bin/`, `dist/`, `build/`
- **Monorepo patterns**: npm workspaces, Cargo workspaces, lerna

## Generated Recommendations

Based on detected tools, the inspector generates check recommendations:

### Example: Go Project with golangci-lint

```yaml
checks:
  - id: fmt
    run: test -z "$(gofmt -l .)"
    severity: error
    suggestion: "Run 'gofmt -w .' to format your Go code."

  - id: lint
    run: golangci-lint run ./...
    severity: error
    suggestion: "Fix linting issues. Run 'golangci-lint run --fix' for auto-fixes."

  - id: test
    run: go test ./...
    severity: error
    suggestion: "Fix failing tests before committing."

  - id: coverage
    run: go test -cover ./... 2>&1 | tail -1
    grok:
      - coverage: %{NUMBER:coverage}%
    assert: "coverage >= 70"
    severity: warning
    suggestion: "Coverage is {{.coverage}}%, target is 70%. Add tests."
    requires:
      - test
```

### Example: Node.js Project with ESLint and Jest

```yaml
checks:
  - id: lint
    run: npx eslint .
    severity: error
    suggestion: "Fix ESLint errors. Run 'npx eslint . --fix' for auto-fixes."

  - id: format
    run: npx prettier --check .
    severity: error
    suggestion: "Run 'npx prettier --write .' to format code."

  - id: test
    run: npm test
    severity: error
    suggestion: "Fix failing tests."
```

### Example: Python Project with pytest and Black

```yaml
checks:
  - id: format
    run: black --check .
    severity: error
    suggestion: "Run 'black .' to format Python code."

  - id: lint
    run: pylint src/
    severity: warning
    suggestion: "Address pylint warnings to improve code quality."

  - id: test
    run: pytest
    severity: error
    suggestion: "Fix failing tests before committing."
```

## Setup Guide Structure

The generated setup guide includes these sections:

### 1. Project Analysis
- Project name and type
- Detection confidence score
- Detected tools with config file locations
- Project structure overview

### 2. Recommended Checks
For each detected tool:
- Check ID and category
- Description and rationale
- Shell command to run
- Severity level
- Optional grok patterns and assertions
- Failure suggestions
- Dependencies

### 3. Configuration Requirements
- YAML syntax rules
- Required and optional fields
- Variable interpolation syntax
- Grok pattern format
- Validation constraints

### 4. Language-Specific Examples
- Format check example
- Lint check example
- Test with coverage example
- Build check with dependencies

### 5. Validation Rules
- Must-have requirements
- DO NOT guidelines
- Common mistakes to avoid

## Customizing the Output

### Verbose Mode

Use `--verbose` to see detailed detection information:

```bash
vibeguard init --assist --verbose
```

This shows:
- All detection indicators found
- Confidence scores for each detection
- Config file locations
- Why certain tools were or weren't detected

### Output to File

Save the guide for later use or sharing:

```bash
vibeguard init --assist --output my-project-guide.txt
```

## Best Practices

### Review Before Committing

Always review the AI-generated configuration:
- Verify check commands are correct for your project
- Adjust timeouts based on your test suite size
- Customize suggestions for your team
- Remove checks for tools you don't use

### Start Simple

Begin with a minimal configuration:
1. Format check
2. Lint check
3. Test check

Add more checks (coverage, security, build) as needed.

### Test the Configuration

After generating:

```bash
# Validate syntax
vibeguard validate

# Run all checks
vibeguard check

# Run with verbose output
vibeguard check -v
```

### Iterate with AI

If the initial configuration needs adjustments:
1. Run `vibeguard check` and note failures
2. Share the error output with the AI agent
3. Ask for specific fixes or improvements

## Troubleshooting

### "Unknown project type"

The inspector couldn't confidently detect your project type. This happens when:
- No standard manifest files exist (`go.mod`, `package.json`, etc.)
- Multiple project types are present with equal confidence

**Solution:** Specify the project type manually when prompting the AI, or add the appropriate manifest file.

### Missing Tool Detection

A tool you use isn't being detected. Common reasons:
- Non-standard config file location
- Tool configured via CLI flags only
- Tool is a transitive dependency

**Solution:** Mention the tool explicitly when prompting the AI.

### Incorrect Recommendations

The generated checks don't match your workflow:
- Custom test commands (e.g., `make test` instead of `go test`)
- Project-specific conventions
- Monorepo structure

**Solution:** Use the recommendations as a starting point and customize as needed.

## Check Templates

The recommendation engine uses predefined check templates for each supported tool. These templates provide sensible defaults that AI agents can use or customize.

### Template Structure

Each template includes:

| Field | Description | Example |
|-------|-------------|---------|
| `ID` | Unique check identifier | `lint`, `fmt`, `test` |
| `Description` | Human-readable description | "Run golangci-lint for code quality" |
| `Rationale` | Why this check is recommended | "Catches bugs the compiler misses" |
| `Command` | Shell command to execute | `golangci-lint run ./...` |
| `Grok` | Output extraction patterns | `coverage: %{NUMBER:coverage}%` |
| `Assert` | Condition for pass/fail | `coverage >= 70` |
| `Severity` | `error` or `warning` | `error` |
| `Suggestion` | Failure guidance (supports `{{.var}}`) | "Coverage is {{.coverage}}%" |
| `Requires` | Dependency check IDs | `["test"]` |
| `Category` | Check category | `lint`, `format`, `test`, `security` |
| `Priority` | Execution order (lower = earlier) | `10` (format), `30` (test) |

### Available Templates by Language

**Go:**
| Tool | Check ID | Category | Priority |
|------|----------|----------|----------|
| gofmt | `fmt` | format | 10 |
| goimports | `imports` | format | 11 |
| go vet | `vet` | lint | 15 |
| golangci-lint | `lint` | lint | 20 |
| go test | `test` | test | 30 |
| go test -cover | `coverage` | test | 35 |
| go build | `build` | build | 5 |

**Node.js:**
| Tool | Check ID | Category | Priority |
|------|----------|----------|----------|
| Prettier | `fmt` | format | 10 |
| ESLint | `lint` | lint | 20 |
| TypeScript | `typecheck` | typecheck | 25 |
| Jest | `test`, `coverage` | test | 30, 35 |
| Mocha | `test` | test | 30 |
| Vitest | `test` | test | 30 |
| npm audit | `security` | security | 50 |

**Python:**
| Tool | Check ID | Category | Priority |
|------|----------|----------|----------|
| Black | `fmt` | format | 10 |
| isort | `imports` | format | 11 |
| Pylint | `lint` | lint | 20 |
| Ruff | `lint` | lint | 20 |
| Flake8 | `lint` | lint | 20 |
| mypy | `typecheck` | typecheck | 25 |
| pytest | `test`, `coverage` | test | 30, 35 |
| pip-audit | `security` | security | 50 |

### Priority Ordering

Checks are ordered by priority (lower values run first):

1. **Build** (5) - Verify compilation
2. **Format** (10-11) - Code formatting
3. **Lint** (15-20) - Static analysis
4. **Typecheck** (25) - Type checking
5. **Test** (30) - Run tests
6. **Coverage** (35) - Check coverage
7. **Security** (50) - Vulnerability scanning

### Customizing Templates

AI agents can customize the generated checks based on project needs:

- Adjust coverage thresholds (default: 70%)
- Add project-specific linter rules
- Modify timeout values
- Change severity levels
- Add custom dependencies

Example customization in the generated config:

```yaml
# Default template
- id: coverage
  run: go test -cover ./...
  assert: "coverage >= 70"

# Customized for stricter requirements
- id: coverage
  run: go test -cover ./... -coverprofile=coverage.out
  grok:
    - total:.*\(statements\)\s+%{NUMBER:coverage}%
  assert: "coverage >= 85"
  severity: error  # Changed from warning
  timeout: 300s    # Extended for large test suite
```

## Sample Prompts

Example setup prompts for different project types are available in [`docs/sample-prompts/`](sample-prompts/):

- [Go project with golangci-lint](sample-prompts/go-project.md)
- [Node.js project with ESLint and Jest](sample-prompts/node-project.md)
- [Python project with Black and pytest](sample-prompts/python-project.md)

## Architecture

The AI-assisted setup feature consists of four components:

```
┌─────────────────┐
│  Project Type   │
│   Detector      │ ──► Identifies Go, Node, Python, etc.
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Tool Scanner   │ ──► Finds linters, formatters, test frameworks
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Metadata      │
│   Extractor     │ ──► Project name, structure, versions
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Recommendation  │
│    Engine       │ ──► Generates check suggestions
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Prompt       │
│   Generator     │ ──► Creates AI-readable setup guide
└─────────────────┘
```

For implementation details, see `internal/cli/inspector/`.

## Related Documentation

- [README.md](../README.md) - Main project documentation
- [Configuration Schema](../README.md#configuration-schema) - Full YAML reference
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Development guidelines
