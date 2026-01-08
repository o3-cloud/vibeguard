# Getting Started with VibeGuard

This guide will help you quickly get up and running with VibeGuard. We'll cover installation, basic configuration, and running your first checks.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Your First Check](#your-first-check)
4. [Configuration Basics](#configuration-basics)
5. [Running Checks](#running-checks)
6. [Common Use Cases](#common-use-cases)
7. [Troubleshooting](#troubleshooting)

## Installation

### Prerequisites

- **Go 1.21 or later** (if building from source)
- **A Unix-like shell** (bash, sh, zsh, etc.)

### Download Pre-built Binary

Pre-built binaries are available for major platforms:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-darwin-arm64 -o vibeguard
chmod +x vibeguard
sudo mv vibeguard /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-darwin-amd64 -o vibeguard
chmod +x vibeguard
sudo mv vibeguard /usr/local/bin/

# Linux (x86_64)
curl -L https://github.com/yourusername/vibeguard/releases/latest/download/vibeguard-linux-amd64 -o vibeguard
chmod +x vibeguard
sudo mv vibeguard /usr/local/bin/

# Verify installation
vibeguard --version
```

### Build from Source

```bash
git clone https://github.com/yourusername/vibeguard.git
cd vibeguard
go build -o vibeguard ./cmd/vibeguard
sudo mv vibeguard /usr/local/bin/
```

### Docker Installation

Create a Dockerfile for your project:

```dockerfile
FROM golang:1.24-alpine as builder
WORKDIR /build
COPY . .
RUN go build -o vibeguard ./cmd/vibeguard

FROM alpine:latest
RUN apk add --no-cache bash
COPY --from=builder /build/vibeguard /usr/local/bin/
ENTRYPOINT ["vibeguard"]
```

## Quick Start

### 1. Initialize a New Configuration

The fastest way to get started is with guided initialization:

```bash
# First, see what templates are available
vibeguard init --list-templates

# Interactive setup with AI assistance
cd your-project
vibeguard init --assist

# Or use a template for your language
vibeguard init -t go-standard    # Go projects
vibeguard init -t node-typescript # TypeScript projects
vibeguard init -t python-poetry  # Python projects
```

This creates a `vibeguard.yaml` file in your project root with sensible defaults for your language.

### 2. Review the Configuration

```bash
cat vibeguard.yaml
```

The generated config includes common checks for your language. Feel free to customize it:

```yaml
version: "1"

vars:
  go_packages: "./..."

checks:
  - id: fmt
    run: gofmt -l .
    severity: error
    suggestion: "Run 'gofmt -w .'"
    timeout: 5s

  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    timeout: 10s

  - id: test
    run: go test {{.go_packages}}
    severity: error
    timeout: 30s
```

### 3. Run Your First Check

```bash
# Run all checks
vibeguard check

# Run a specific check
vibeguard check fmt

# Verbose output (show all results)
vibeguard check -v

# Stop on first failure
vibeguard check --fail-fast
```

**Expected output** (if all checks pass):

```
Exit code: 0
```

**Expected output** (if a check fails):

```
✗ fmt: Failed to apply gofmt
  Suggestion: Run 'gofmt -w .'

Exit code: 3
```

## Your First Check

Let's create a simple configuration with one check:

```bash
cat > vibeguard.yaml << 'EOF'
version: "1"

checks:
  - id: hello
    run: echo "Hello from VibeGuard"
    severity: error
    timeout: 5s
EOF
```

Run it:

```bash
vibeguard check
```

You should see:

```
✓ hello (45ms)
Exit code: 0
```

Now let's add a check that validates output:

```bash
cat > vibeguard.yaml << 'EOF'
version: "1"

checks:
  - id: hello
    run: echo "Hello from VibeGuard"
    grok:
      - 'Hello from (?P<source>\w+)'
    assert: 'source == "VibeGuard"'
    severity: error
    timeout: 5s
EOF
```

Run it again:

```bash
vibeguard check -v
```

## Configuration Basics

### Basic Check Structure

Every check needs an `id` and `run` command:

```yaml
checks:
  - id: my-check
    run: some-command --flag
    severity: error
```

### Adding Severity Levels

Control how failures are treated:

```yaml
checks:
  - id: lint
    run: eslint .
    severity: error        # Failure causes exit code 3

  - id: coverage
    run: npm test -- --coverage
    severity: warning      # Failure logged but doesn't fail
```

### Using Variables

Define reusable values:

```yaml
version: "1"

vars:
  go_packages: "./..."
  min_coverage: "70"
  timeout_quick: "10s"
  timeout_long: "60s"

checks:
  - id: test
    run: go test {{.go_packages}}
    timeout: {{.timeout_long}}

  - id: coverage
    run: go test {{.go_packages}} -coverprofile=cover.out
    assert: "coverage >= {{.min_coverage}}"
    timeout: {{.timeout_long}}
```

### Pattern Extraction (Grok)

Extract structured data from command output:

```yaml
checks:
  - id: coverage
    run: go tool cover -func=cover.out
    grok:
      - 'total:\s+%{NUMBER:coverage}%'
    assert: "coverage >= 80"
```

### Check Dependencies

Make checks run in order:

```yaml
checks:
  - id: fmt
    run: gofmt -l .

  - id: vet
    run: go vet ./...
    requires: [fmt]        # Only runs if fmt passes

  - id: test
    run: go test ./...
    requires: [fmt, vet]   # Only runs if both pass
```

## Running Checks

### Basic Command

```bash
# Run all checks
vibeguard check

# Run specific check
vibeguard check fmt

# Run multiple specific checks
vibeguard check fmt vet test
```

### Common Flags

```bash
# Show detailed output (all results, not just failures)
vibeguard check -v

# Output as JSON (useful for CI/CD)
vibeguard check --json

# Stop after first failure
vibeguard check --fail-fast

# Use custom config file
vibeguard check -c ./path/to/vibeguard.yaml

# Limit parallel execution
vibeguard check -p 2   # Run max 2 checks in parallel (default: 4)

# Save check logs to custom directory
vibeguard check --log-dir ./custom-logs
```

### List Available Checks

```bash
vibeguard list
```

Shows all configured checks and their dependencies:

```
Check: fmt (error severity)
  Command: gofmt -l .
  Timeout: 5s

Check: test (error severity)
  Command: go test ./...
  Timeout: 30s
  Requires: fmt
```

### Validate Configuration

```bash
vibeguard validate
```

Checks for YAML errors and invalid configuration without running checks.

## Common Use Cases

### Go Project Quality Gates

```yaml
version: "1"

vars:
  go_packages: "./..."
  min_coverage: "70"

checks:
  - id: fmt
    run: test -z "$(gofmt -l .)"
    severity: error
    suggestion: "Run 'gofmt -w .'"
    timeout: 5s

  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    suggestion: "Run 'go vet {{.go_packages}}' and fix issues"
    timeout: 10s

  - id: test
    run: go test {{.go_packages}} -coverprofile=cover.out
    severity: error
    suggestion: "Run 'go test {{.go_packages}}' to see failures"
    timeout: 60s

  - id: coverage
    run: go tool cover -func=cover.out | tail -1
    grok:
      - 'total:\s+%{NUMBER:coverage}%'
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Coverage is {{.coverage}}%, target is {{.min_coverage}}%"
    timeout: 5s
    requires: [test]
```

### Node.js/TypeScript Project

```yaml
version: "1"

vars:
  min_coverage: "75"

checks:
  - id: lint
    run: npx eslint . --max-warnings 0
    severity: error
    suggestion: "Run 'npm run lint:fix' to auto-fix"
    timeout: 30s

  - id: type-check
    run: npx tsc --noEmit
    severity: error
    suggestion: "Fix TypeScript errors and rerun"
    timeout: 30s

  - id: test
    run: npm test -- --coverage --json --outputFile=coverage.json
    severity: error
    suggestion: "Fix failing tests"
    timeout: 60s
    requires: [lint, type-check]

  - id: coverage
    run: cat coverage.json
    grok:
      - 'coveragePercent[^:]*:\s*%{NUMBER:coverage}'
    assert: "coverage >= {{.min_coverage}}"
    severity: warning
    suggestion: "Coverage is {{.coverage}}%, target is {{.min_coverage}}%"
    timeout: 5s
    requires: [test]
```

### Python Project

```yaml
version: "1"

checks:
  - id: lint
    run: pylint src/ --fail-under=8.0
    severity: error
    suggestion: "Run 'pylint src/' to see issues"
    timeout: 30s

  - id: format
    run: black --check src/
    severity: error
    suggestion: "Run 'black src/' to auto-format"
    timeout: 10s

  - id: test
    run: pytest --cov=src --cov-report=xml
    severity: error
    suggestion: "Run 'pytest' to see failures"
    timeout: 60s
    requires: [lint, format]
```

### Custom Multi-Language Setup

```yaml
version: "1"

checks:
  # JavaScript tests
  - id: js-test
    run: npm test
    severity: error
    timeout: 60s

  # Go tests
  - id: go-test
    run: go test ./...
    severity: error
    timeout: 60s

  # Run final integration test after both languages pass
  - id: integration
    run: ./scripts/integration-test.sh
    severity: error
    timeout: 120s
    requires: [js-test, go-test]
```

## Troubleshooting

### Configuration Not Found

**Problem:** `error: config file not found`

**Solution:** VibeGuard searches for config files in this order:
1. Path specified with `-c` flag
2. `vibeguard.yaml`
3. `vibeguard.yml`
4. `.vibeguard.yaml`
5. `.vibeguard.yml`

Place your config file in the project root or specify the path:

```bash
vibeguard check -c ./path/to/config.yaml
```

### Checks Timing Out

**Problem:** `exit code: 4 (timeout)`

**Solution:** Increase the timeout in your config:

```yaml
checks:
  - id: slow-check
    run: some-slow-command
    timeout: 120s    # Increased from default 30s
```

### Variables Not Substituting

**Problem:** Check command shows `{{.variable_name}}` instead of the value

**Solution:** Make sure your variable is defined in `vars`:

```yaml
vars:
  my_var: "value"  # Define it here

checks:
  - id: test
    run: command {{.my_var}}  # Use it here
```

### Pattern Extraction Not Working

**Problem:** Assertion always fails even though output looks correct

**Solution:** Check your Grok pattern:

1. Run the command manually and verify output
2. Test your pattern with correct syntax:

```bash
# View actual output
your-command | cat -A

# Test grok pattern
vibeguard check -v
```

Common issues:
- **Whitespace** - Use `\s+` for flexible whitespace
- **Regex escaping** - Escape special characters: `\.`, `\(`, `\)`
- **Named captures** - Use `%{NUMBER:name}` or `(?P<name>\d+)`

### Circular Dependency Error

**Problem:** `error: circular dependency detected`

**Solution:** Check your `requires` declarations for cycles:

```yaml
# BAD - circular dependency
checks:
  - id: a
    requires: [b]
  - id: b
    requires: [a]   # This creates a cycle!

# GOOD - linear dependency
checks:
  - id: a
  - id: b
    requires: [a]
```

### Exit Code Confusion

VibeGuard uses specific exit codes:

| Code | Meaning | Action |
|------|---------|--------|
| 0 | All checks passed | Continue |
| 2 | Configuration error | Fix YAML/config |
| 3 | Check failed (error severity) | Fix the issue |
| 4 | Timeout or command not found | Increase timeout or install tool |

### Parallel Execution Issues

**Problem:** Checks fail intermittently

**Solution:** Try reducing parallelism:

```bash
vibeguard check -p 1    # Run sequentially
```

If this fixes it, your checks may have race conditions or shared state issues.

### Output Not in Logs

**Problem:** Check output is missing from `.vibeguard/log/`

**Solution:** Check logs are written to:

```
.vibeguard/log/{check-id}.log
```

Verify the `.vibeguard/` directory exists and is writable:

```bash
ls -la .vibeguard/
cat .vibeguard/log/my-check.log
```

## Next Steps

- **[Read ARCHITECTURE.md](./ARCHITECTURE.md)** - Understand how VibeGuard works
- **[Read CLI Reference](./CLI-REFERENCE.md)** - Detailed flag and command documentation
- **[See Integration Guides](./INTEGRATIONS.md)** - CI/CD integration examples
- **[Check Examples](../examples/)** - Real-world configuration examples
- **[Read ADRs](./adr/)** - Architecture decisions and rationale

## Getting Help

- **Check [TROUBLESHOOTING.md](../TROUBLESHOOTING.md)** for common issues
- **Review [examples/](../examples/)** for working configurations
- **Check [docs/INTEGRATIONS.md](./INTEGRATIONS.md)** for CI/CD help
- **See [CONTRIBUTING.md](../CONTRIBUTING.md)** for development setup
