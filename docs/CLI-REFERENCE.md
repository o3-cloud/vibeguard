# CLI Reference

Complete reference for all VibeGuard commands and flags.

## Table of Contents

1. [Global Flags](#global-flags)
2. [Commands](#commands)
   - [check](#vibeguard-check)
   - [init](#vibeguard-init)
   - [list](#vibeguard-list)
   - [validate](#vibeguard-validate)
3. [Exit Codes](#exit-codes)
4. [Environment Variables](#environment-variables)
5. [Configuration File Discovery](#configuration-file-discovery)

## Global Flags

Global flags are available on all commands:

### `-c, --config` (string)

Path to VibeGuard configuration file. If not specified, VibeGuard searches for configuration files in the current directory.

**Default:** Auto-discovery (searches for `vibeguard.yaml`, `vibeguard.yml`, `.vibeguard.yaml`, `.vibeguard.yml`)

**Examples:**
```bash
vibeguard -c ./config/vibeguard.yaml check
vibeguard --config /etc/vibeguard/config.yaml check
```

### `-v, --verbose` (boolean)

Show all check results, not just failures. In verbose mode, all checks are displayed with their status (pass, fail, or cancelled) and execution time.

**Default:** `false`

**Examples:**
```bash
vibeguard check -v
vibeguard check --verbose
```

**Output comparison:**

Without `--verbose`:
```
✗ test: Tests failed
  Suggestion: Run tests locally
Exit code: 3
```

With `--verbose`:
```
✓ fmt (12ms)
✓ vet (156ms)
✗ test (2.3s)
  Suggestion: Run tests locally
⊘ coverage (cancelled - dependency failed)

Exit code: 3
```

### `--json` (boolean)

Output results in JSON format to stderr. Useful for CI/CD integration and automation.

**Default:** `false`

**Examples:**
```bash
vibeguard check --json
vibeguard check --json 2>/tmp/results.json
```

**JSON Output Format:**
```json
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
      "duration_ms": 2300,
      "suggestion": "Run 'go test ./...' to see failures"
    }
  ],
  "violations": ["test"],
  "exit_code": 3,
  "fail_fast_triggered": false
}
```

### `-p, --parallel` (int)

Maximum number of checks to run in parallel. VibeGuard respects check dependencies and runs checks at the same dependency level in parallel.

**Default:** `4`

**Valid range:** `1` to `256`

**Examples:**
```bash
vibeguard check -p 1    # Run checks sequentially
vibeguard check --parallel 8    # Run up to 8 in parallel
```

**Notes:**
- Setting to `1` effectively runs checks sequentially
- Checks respect dependencies regardless of this setting
- Useful for debugging race conditions or limiting system load

### `--fail-fast` (boolean)

Stop execution on the first error-severity check failure. Remaining checks are cancelled and not executed.

**Default:** `false`

**Examples:**
```bash
vibeguard check --fail-fast
```

**Behavior:**
- If an error-severity check fails, immediately stop
- Warning-severity checks do not trigger fail-fast
- Remaining checks in the queue are cancelled
- Cancelled checks show status as `⊘` in output
- Exit code is still `3` (violation)

### `--log-dir` (string)

Directory where individual check logs are written. Each check creates a separate log file.

**Default:** `.vibeguard/log`

**Examples:**
```bash
vibeguard check --log-dir ./logs
vibeguard check --log-dir /tmp/vibeguard-logs
```

**Output locations:**
```
.vibeguard/log/
├── fmt.log
├── vet.log
├── test.log
└── coverage.log
```

**Note:** Directory is created if it doesn't exist.

## Commands

### `vibeguard check` [id...]

Run configured checks. Optionally specify one or more check IDs to run only those checks.

**Syntax:**
```bash
vibeguard check [id...]
```

**Examples:**
```bash
# Run all checks
vibeguard check

# Run specific checks
vibeguard check fmt vet

# Run all checks with verbose output
vibeguard check -v

# Run specific check with verbose output
vibeguard check fmt -v

# Run all checks, stop on first failure
vibeguard check --fail-fast

# Run with custom config and parallel limit
vibeguard -c custom.yaml check -p 2
```

**Behavior:**
1. Loads configuration from disk
2. Builds dependency graph
3. Executes checks respecting dependencies
4. Extracts patterns from output (if grok specified)
5. Evaluates assertions (if specified)
6. Formats and outputs results
7. Exits with appropriate code

**Exit codes:**
- `0` - All checks passed
- `2` - Configuration error
- `3` - Error-severity check failed
- `4` - Timeout or command not found

### `vibeguard init` [--assist]

Initialize a new VibeGuard configuration file.

**Syntax:**
```bash
vibeguard init [--assist]
```

**Flags:**

#### `--assist` (boolean)

Use AI-assisted setup for guided configuration creation.

**Examples:**
```bash
vibeguard init --assist
```

#### `-t, --template` (string)

Use a predefined template for your project type.

**Available templates:**
- `go-standard` - Comprehensive Go project setup
- `go-minimal` - Minimal Go setup
- `node-typescript` - TypeScript/Node.js
- `node-javascript` - JavaScript/Node.js
- `python-poetry` - Python with Poetry
- `python-pip` - Python with pip
- `rust-cargo` - Rust/Cargo
- `generic` - Generic/minimal

**Examples:**
```bash
vibeguard init -t go-standard
vibeguard init --template node-typescript
```

#### `-f, --force` (boolean)

Overwrite existing configuration file without prompting.

**Examples:**
```bash
vibeguard init -t go-standard -f
```

#### `-o, --output` (string)

Output file path for configuration (mainly for `--assist` mode).

**Default:** `vibeguard.yaml`

**Examples:**
```bash
vibeguard init --assist -o ./config/vibeguard.yaml
```

**Behavior:**
1. If `--assist` specified: Analyzes project and guides creation via prompts
2. If `-t` specified: Uses template for the language
3. If neither: Creates default Go-based configuration
4. If file exists: Prompts for confirmation (unless `-f` specified)
5. Writes configuration to specified output file

**Example output:**
```yaml
version: "1"

vars:
  go_packages: "./..."
  min_coverage: "70"

checks:
  - id: fmt
    run: gofmt -l .
    severity: error
    suggestion: "Run 'gofmt -w .'"
    timeout: 5s

  - id: vet
    run: go vet {{.go_packages}}
    severity: error
    suggestion: "Run 'go vet {{.go_packages}}'"
    timeout: 10s

  - id: test
    run: go test {{.go_packages}}
    severity: error
    suggestion: "Run 'go test {{.go_packages}}'"
    timeout: 30s
    requires: [fmt, vet]
```

### `vibeguard list`

Display all configured checks and their metadata.

**Syntax:**
```bash
vibeguard list
```

**Examples:**
```bash
vibeguard list
vibeguard list -c ./custom.yaml
```

**Output format:**
```
Check: fmt (error severity)
  Command: gofmt -l .
  Timeout: 5s

Check: vet (error severity)
  Command: go vet {{.go_packages}}
  Timeout: 10s
  Requires: fmt

Check: test (error severity)
  Command: go test {{.go_packages}}
  Timeout: 30s
  Requires: fmt, vet

Check: coverage (warning severity)
  Command: go tool cover -func=cover.out
  Grok: total:\s+coverage:\s+%{NUMBER:coverage}%
  Assert: coverage >= {{.min_coverage}}
  Timeout: 5s
  Requires: test
```

### `vibeguard validate`

Validate configuration file without running checks.

**Syntax:**
```bash
vibeguard validate
```

**Examples:**
```bash
vibeguard validate
vibeguard validate -c ./config/vibeguard.yaml
```

**Checks performed:**
1. YAML syntax validation
2. Required fields present
3. Field type validation
4. Circular dependency detection
5. Timeout format validation
6. Severity value validation
7. Check ID uniqueness

**Output:**
- Success: No output, exit code `0`
- Error: Error message with details, exit code `2`

**Example error:**
```
error: validation failed: check 'test' requires non-existent check 'build'
```

### `vibeguard --version`

Display version information.

**Syntax:**
```bash
vibeguard --version
vibeguard -V
```

**Output format:**
```
VibeGuard version 1.0.0 (build: abc123def, go version: go1.24.4)
```

### `vibeguard --help`

Display help message for commands.

**Syntax:**
```bash
vibeguard --help
vibeguard check --help
vibeguard init --help
```

**Examples:**
```bash
vibeguard -h              # Main help
vibeguard check -h        # Help for check command
vibeguard init -h         # Help for init command
```

## Exit Codes

VibeGuard uses specific exit codes to indicate different types of results:

| Code | Name | Meaning | Action |
|------|------|---------|--------|
| 0 | SUCCESS | All checks passed | Continue normally |
| 2 | CONFIG_ERROR | Configuration error | Fix YAML/config file |
| 3 | VIOLATION | Error-severity check failed | Fix the issue |
| 4 | TIMEOUT | Check timeout or command not found | Increase timeout or install tool |

**Exit code selection logic:**
1. If configuration is invalid → exit code `2`
2. If error-severity check fails → exit code `3`
3. If check times out or command not found → exit code `4`
4. If all checks pass → exit code `0`
5. If warning-severity checks fail → exit code `0` (but message shown)

## Environment Variables

VibeGuard respects the following environment variables:

### `VIBEGUARD_CONFIG` (string)

Default configuration file path. Overridden by `-c` flag if specified.

**Example:**
```bash
export VIBEGUARD_CONFIG=./config/vibeguard.yaml
vibeguard check    # Uses VIBEGUARD_CONFIG
```

### `VIBEGUARD_LOG_DIR` (string)

Default log directory. Overridden by `--log-dir` flag if specified.

**Example:**
```bash
export VIBEGUARD_LOG_DIR=/tmp/vibeguard-logs
vibeguard check
```

### `VIBEGUARD_PARALLEL` (int)

Default parallel execution limit. Overridden by `-p` flag if specified.

**Example:**
```bash
export VIBEGUARD_PARALLEL=2
vibeguard check    # Runs max 2 checks in parallel
```

## Configuration File Discovery

When no `-c` flag is specified, VibeGuard searches for configuration files in this order:

1. `vibeguard.yaml` - Current directory
2. `vibeguard.yml` - Current directory
3. `.vibeguard.yaml` - Current directory (hidden file)
4. `.vibeguard.yml` - Current directory (hidden file)

The first file found is used. If none exist, an error is displayed.

**To use a specific config file:**
```bash
vibeguard -c /path/to/config.yaml check
vibeguard -c ./config/vibeguard.yaml check
vibeguard --config ~/projects/myapp/vibeguard.yaml check
```

## Common Command Combinations

### Run checks locally before committing
```bash
vibeguard check --fail-fast
```

### Run checks in CI with detailed output
```bash
vibeguard check -v --json
```

### Run only specific checks
```bash
vibeguard check fmt vet test
```

### Run with limited parallelism for debugging
```bash
vibeguard check -p 1 -v
```

### Validate config without running checks
```bash
vibeguard validate -c ./custom.yaml
```

### Initialize and run Go project checks
```bash
vibeguard init -t go-standard && vibeguard check
```

### Run with AI-assisted setup
```bash
vibeguard init --assist
```

### Save results to file
```bash
vibeguard check --json > results.json 2>&1
```

### See all available checks
```bash
vibeguard list
```

## Flag Combinations

### Debugging Checks
```bash
# Verbose output with sequential execution
vibeguard check -v -p 1

# Show results as JSON
vibeguard check --json -v

# Check specific checks sequentially
vibeguard check fmt vet -p 1 -v
```

### CI/CD Pipelines
```bash
# Fast feedback, stop on first error
vibeguard check --fail-fast

# JSON output for tool integration
vibeguard check --json

# Parallel with custom config
vibeguard -c ci/vibeguard.yaml check -p 8
```

### Local Development
```bash
# See all results
vibeguard check -v

# Run specific checks frequently
vibeguard check fmt

# Stop on first failure
vibeguard check --fail-fast
```
