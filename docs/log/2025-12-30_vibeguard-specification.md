---
summary: Comprehensive specification for VibeGuard v1.0 consolidating architecture spike findings, simplified checks schema, CLI output design, and LLM-as-judge patterns into an authoritative implementation reference.
event_type: deep dive
sources:
  - docs/log/2025-12-30_vibeguard-architecture-spike.md
  - docs/log/2025-12-30_vibeguard-implementation-patterns.md
  - docs/log/2025-12-30_vibeguard-cli-output-design.md
  - docs/log/2025-12-30_simplified-checks-schema.md
  - docs/log/2025-12-30_llm-as-judge-and-templated-suggestions.md
tags:
  - vibeguard
  - specification
  - v1.0
  - checks
  - cli
  - grok
  - llm-as-judge
---

# VibeGuard v1.0 Specification

## 1. Overview

VibeGuard is a unified code quality and policy enforcement tool designed for CI/CD pipelines and AI agent workflows. It orchestrates external tools, evaluates assertions against their output, and emits actionable signals only when violations occur.

### Core Principles

| Principle | Description |
|-----------|-------------|
| **Silence is Success** | No output when all checks pass; violations produce structured, actionable output |
| **Tools are Commands** | Any CLI tool (linters, test runners, LLMs) integrates via shell commands |
| **Simple by Default** | Exit code pass/fail requires no config; complexity opt-in via grok/assert |
| **Actionable Output** | Every violation answers: "What failed?" and "What should I do next?" |

---

## 2. Configuration Schema

### 2.1 File Location

VibeGuard looks for configuration in this order:
1. `vibeguard.yaml`
2. `vibeguard.yml`
3. `.vibeguard.yaml`
4. `.vibeguard.yml`

### 2.2 Full Schema

```yaml
version: "1"                    # Schema version (required)

vars:                           # Optional: variable definitions
  KEY: "value"

checks:                         # Required: list of checks
  - id: string                  # Required: unique identifier
    run: string                 # Required: shell command to execute

    # Output extraction (optional)
    grok: string | [string]     # Grok pattern(s) to extract values from stdout
    file: string                # Read file contents instead of stdout (treated identically to stdout)

    # Evaluation (optional)
    assert: string              # Expression to evaluate (default: exit code == 0)
    severity: error | warning   # Violation severity (default: error)
    suggestion: string          # Guidance shown on failure (supports {{.var}} templating)

    # Execution control (optional)
    requires: [id, ...]         # Dependencies (run after these checks complete)
    timeout: duration           # Max execution time (default: 30s)
```

### 2.3 Minimal Example

```yaml
version: "1"

checks:
  - id: vet
    run: go vet ./...

  - id: test
    run: go test ./...
```

### 2.4 Full Example

```yaml
version: "1"

vars:
  MIN_COVERAGE: "80"

checks:
  # Fast deterministic checks
  - id: fmt
    run: "! gofmt -l . | grep ."

  - id: vet
    run: go vet ./...

  - id: lint
    run: golangci-lint run

  - id: test
    run: go test ./...

  - id: coverage
    run: go test -cover ./...
    grok: 'coverage: %{NUMBER:coverage}%'
    assert: coverage >= {{.MIN_COVERAGE}}
    suggestion: "Coverage is {{.coverage}}%, need {{.MIN_COVERAGE}}%."

  # LLM check (runs after fast checks)
  - id: llm-review
    run: |
      claude -p "Review this diff. Output: VERDICT: PASS/FAIL | REASON: <why>
      $(git diff main --staged)"
    grok: 'VERDICT: %{WORD:verdict} | REASON: %{GREEDYDATA:reason}'
    assert: verdict == "PASS"
    requires: [fmt, vet, lint, test]
    timeout: 60s
    severity: warning
    suggestion: "{{.reason}}"
```

---

## 3. Check Evaluation

### 3.1 Evaluation Modes

| Mode | Config Required | Pass Condition |
|------|-----------------|----------------|
| **Exit Code** | `run` only | Command exits with code 0 |
| **Grok + Assert** | `run` + `grok` + `assert` | Assertion evaluates to true |
| **JSON + Assert** | `run` + `assert` (JSON stdout) | JSON path assertion evaluates to true |

### 3.2 Grok Patterns

VibeGuard uses [elastic/go-grok](https://github.com/elastic/go-grok) for pattern extraction.

**Output Source:** Grok patterns are applied to combined stdout+stderr output. If `file:` is specified, the file contents are read and treated identically to stdout.

**Match Behavior:** When multiple grok patterns are specified, each pattern is applied independently to the output. If a pattern does not match, its variables are set to empty strings. This allows flexible extraction where some fields may be optional.

**Built-in Patterns:**

| Pattern | Matches | Example |
|---------|---------|---------|
| `%{NUMBER:n}` | Float or integer | `85.5`, `42` |
| `%{INT:n}` | Integer only | `42` |
| `%{WORD:w}` | Single word (no spaces) | `main` |
| `%{DATA:d}` | Non-greedy any text | `foo bar` |
| `%{GREEDYDATA:d}` | Greedy any text | `everything to end...` |
| `%{PATH:p}` | File path | `/src/main.go` |
| `%{IP:ip}` | IP address | `192.168.1.1` |

**Multiple Patterns:**

```yaml
grok:
  - 'High: %{INT:high}'
  - 'Medium: %{INT:medium}'
  - 'Low: %{INT:low}'
```

### 3.3 Assertion Expressions

Assertions support:
- Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
- Logical: `&&`, `||`, `!`
- Arithmetic: `+`, `-`, `*`, `/`
- Variables: extracted values from grok, JSON paths

**Type Coercion:** Grok-extracted values are automatically converted to the appropriate type for comparison:
- Numeric strings (`"72"`, `"85.5"`) are converted to numbers when compared with numeric operators
- Boolean strings (`"true"`, `"false"`) are converted to booleans
- All other values remain strings

**Examples:**
```yaml
assert: coverage >= 80          # coverage="72" → 72 >= 80 → false
assert: high == 0 && medium <= 5
assert: verdict == "PASS"       # string comparison (case-sensitive)
assert: json.safe == true
```

### 3.4 Templated Suggestions

Suggestions use Go template syntax (`{{.var}}`) with extracted values and config variables:

```yaml
suggestion: "Coverage is {{.coverage}}%, need 80%."
suggestion: "{{.reason}}"
suggestion: "{{.func}} in {{.file}}:{{.line}} has complexity {{.score}}."
```

All variables (both grok-extracted and config `vars`) use the same `{{.name}}` syntax.

---

## 4. CLI Interface

### 4.1 Commands

```
vibeguard check [flags]         Run all checks
vibeguard check <id> [flags]    Run specific check
vibeguard init                  Create starter vibeguard.yaml
vibeguard list                  List configured checks
vibeguard validate              Validate configuration
```

### 4.2 Flags

| Flag | Description |
|------|-------------|
| `--config`, `-c` | Path to config file |
| `--verbose`, `-v` | Show all check results, not just failures |
| `--json` | Output in JSON format |
| `--parallel`, `-p` | Max parallel checks (default: 4) |
| `--fail-fast` | Stop on first failure |

### 4.3 Output Modes

#### Default (Quiet) - Success

```
$ vibeguard check
$ echo $?
0
```

No output. Exit code 0 indicates all checks passed.

#### Default (Quiet) - Violation

```
$ vibeguard check
FAIL  coverage (error)
      > go test -cover ./...

      Tip: Coverage is 72%, need 80%.

$ echo $?
2
```

**Note:** All output is written to **stderr**, not stdout. This ensures output is visible to Claude Code PostToolUse hooks, which hide stdout by default.

#### Verbose Mode

```
$ vibeguard check --verbose
✓ fmt              passed (0.1s)
✓ vet              passed (0.3s)
✓ lint             passed (0.8s)
✓ test             passed (1.2s)
✗ coverage         FAIL (0.9s)
  Coverage is 72%, need 80%
```

#### JSON Output

```json
{
  "checks": [
    {"id": "fmt", "status": "passed", "duration_ms": 100},
    {"id": "vet", "status": "passed", "duration_ms": 300},
    {"id": "coverage", "status": "failed", "duration_ms": 900}
  ],
  "violations": [
    {
      "id": "coverage",
      "severity": "error",
      "command": "go test -cover ./...",
      "suggestion": "Coverage is 72%, need 80%.",
      "extracted": {"coverage": "72"}
    }
  ],
  "exit_code": 2
}
```

### 4.4 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed, or only warning-severity violations |
| 2 | One or more error-severity violations |
| 3 | Configuration error |
| 4 | Check execution error (timeout, command not found) |

**Note:** Warning-severity violations are reported in output but do not cause a non-zero exit code. This allows CI pipelines to surface warnings without blocking merges.

**Claude Code Hook Compatibility:** Exit code 2 is used for violations (rather than 1) because Claude Code hooks treat exit codes 0 and 1 as non-blocking, while exit codes ≥2 block the tool call. This ensures VibeGuard policy violations prevent commits when used as a git pre-commit hook.

---

## 5. Execution Model

### 5.1 Working Directory and Environment

**Working Directory:** All checks execute in the current working directory (CWD) where `vibeguard` was invoked. Relative paths in commands and the `file:` field are resolved relative to this directory.

**Environment Inheritance:** Processes executed by VibeGuard inherit the **full shell environment** from the parent process. This means:

- All environment variables (`PATH`, `HOME`, `GOPATH`, etc.) are available
- Tool configurations (`.npmrc`, `.gitconfig`, etc.) work as expected
- Credentials and tokens in the environment are accessible to checks
- No sandboxing or environment isolation is applied

This design ensures checks behave identically whether run directly or via VibeGuard.

### 5.2 Dependency Resolution

Checks execute in topological order based on `requires` dependencies:

```
Level 0 (parallel):  fmt, vet, lint      ─┐
Level 1 (parallel):  test                ←┤ (requires nothing)
Level 2:             coverage            ←┘ (requires test)
Level 3:             llm-review          ← (requires fmt, vet, lint, test)
```

### 5.3 Parallelism

- Checks at the same dependency level run in parallel
- `--parallel` flag controls max concurrent checks (default: 4)
- `--fail-fast` stops execution on first failure

### 5.4 Timeouts

- Default timeout: 30 seconds per check
- Override per-check with `timeout: 60s`
- Supported units: `s`, `m`, `h`

---

## 6. LLM Integration

LLMs are treated as ordinary CLI commands. No special integration required.

### 6.1 Supported CLI Tools

| Tool | Command Pattern |
|------|-----------------|
| Claude Code | `claude -p "prompt"` |
| Gemini | `gemini -p "prompt"` |
| OpenAI | `openai api chat.completions.create -m gpt-4 -p "prompt"` |
| Ollama | `ollama run llama3 "prompt"` |
| LLM (multi-provider) | `llm "prompt"` |

### 6.2 Prompt Engineering Guidelines

For reliable extraction:

1. **Constrain output format strictly**
   ```
   Output ONLY one line.
   Format: VERDICT: PASS or FAIL | REASON: <brief explanation>
   ```

2. **JSON is more reliable than free-form**
   ```
   Output JSON only: {"pass": boolean, "reason": string}
   ```

3. **Provide examples**
   ```
   Examples:
   VERDICT: PASS | REASON: Code follows best practices
   VERDICT: FAIL | REASON: Missing input validation
   ```

### 6.3 Example: Architecture Review

```yaml
- id: architecture-review
  run: |
    claude -p "Review this diff for architectural issues.
    Output exactly: VERDICT: PASS or FAIL | REASON: <explanation>

    $(git diff main --staged)"
  grok: 'VERDICT: %{WORD:verdict} | REASON: %{GREEDYDATA:reason}'
  assert: verdict == "PASS"
  requires: [test, lint]
  timeout: 60s
  severity: warning
  suggestion: "{{.reason}}"
```

---

## 7. Internal Architecture

### 7.1 Package Structure

```
vibeguard/
├── cmd/
│   └── vibeguard/
│       └── main.go              # CLI entrypoint
├── internal/
│   ├── config/
│   │   ├── config.go            # YAML parsing and validation
│   │   ├── interpolate.go       # Variable interpolation
│   │   └── schema.go            # Type definitions
│   ├── executor/
│   │   ├── executor.go          # Command execution
│   │   └── output.go            # Output capture and parsing
│   ├── grok/
│   │   └── grok.go              # Grok pattern matching
│   ├── assert/
│   │   └── eval.go              # Assertion evaluation
│   ├── orchestrator/
│   │   ├── orchestrator.go      # Parallel execution
│   │   └── graph.go             # Dependency graph / toposort
│   └── output/
│       ├── formatter.go         # Output formatting
│       └── json.go              # JSON output
└── examples/
    ├── go-project.yaml
    ├── node-project.yaml
    └── simple.yaml
```

### 7.2 Execution Flow

```
┌─────────────────────────────────────────────────────────────┐
│                      vibeguard check                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. Load Config                                              │
│    - Parse YAML                                             │
│    - Interpolate variables                                  │
│    - Validate schema                                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Build Dependency Graph                                   │
│    - Topological sort checks                                │
│    - Group by execution level                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Execute Checks (parallel within levels)                  │
│    For each check:                                          │
│    - Run command via os/exec                                │
│    - Capture stdout/stderr                                  │
│    - Apply grok patterns (if configured)                    │
│    - Evaluate assertion (if configured)                     │
│    - Collect result                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Collect Violations                                       │
│    - Filter failed checks                                   │
│    - Render suggestions with templates                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Output Results                                           │
│    - Quiet mode: violations only                            │
│    - Verbose mode: all results                              │
│    - JSON mode: structured output                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Exit                                                     │
│    - Code 0: all passed                                     │
│    - Code 2: violations                                     │
│    - Code 3/4: errors                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 8. Implementation Phases

### Phase 1: Core CLI

- [ ] Project scaffolding (go mod, cmd structure)
- [ ] Config parsing and validation
- [ ] Variable interpolation
- [ ] Basic executor (run commands, capture output)
- [ ] Exit code-based pass/fail
- [ ] CLI with `check`, `init`, `validate` commands
- [ ] Quiet and verbose output modes

### Phase 1.5: Dogfooding

- [ ] Create `vibeguard.yaml` for the vibeguard project itself
- [ ] Run `go vet`, `go fmt`, `go test` via vibeguard
- [ ] Use vibeguard in CI from this point forward

**Rationale:** Dogfood as early as possible to validate the tool's usability and catch friction points before adding complexity.

### Phase 2: Grok + Assertions

- [ ] Integrate go-grok library
- [ ] Assertion expression parser and evaluator
- [ ] Templated suggestions
- [ ] JSON output mode

### Phase 3: Orchestration

- [ ] Dependency graph construction
- [ ] Topological sort execution ordering
- [ ] Parallel execution with errgroup
- [ ] Fail-fast mode
- [ ] Timeout handling

### Phase 4: Polish

- [ ] Comprehensive error messages
- [ ] Example configurations
- [ ] Integration tests with real tools
- [ ] Documentation

---

## 9. Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `gopkg.in/yaml.v3` | YAML parsing | latest |
| `github.com/elastic/go-grok` | Grok pattern matching | latest |
| `github.com/spf13/cobra` | CLI framework | v1.8+ |
| `golang.org/x/sync/errgroup` | Parallel execution | latest |

---

## 10. Explicit Non-Goals (v1.0)

The following features are intentionally **not supported** in v1.0 to keep the tool simple and focused:

| Feature | Rationale |
|---------|-----------|
| **Conditional execution** | No `when:`, `skip_if:`, or environment-based conditionals. Use separate config files or wrapper scripts if needed. |
| **Retry mechanism** | No automatic retries for flaky checks. Checks should be deterministic; transient failures should be handled at the tool level (e.g., test runners with built-in retries). |
| **Caching** | No result caching between runs. Each invocation executes all applicable checks fresh. |
| **Remote execution** | All checks run locally. No support for distributed or remote check execution. |

---

## 11. Success Criteria

VibeGuard v1.0 is complete when:

1. **Basic checks work**: Exit code-based pass/fail with 2-line config
2. **Grok extraction works**: Parse coverage percentages, counts, etc.
3. **Assertions work**: `coverage >= 80`, `errors == 0`
4. **Dependencies work**: `requires: [test]` orders execution correctly
5. **Parallel execution works**: Checks at same level run concurrently
6. **Output is clean**: Silence on success, actionable output on failure
7. **LLMs work**: Claude/Gemini/Ollama treated as ordinary commands
8. **CI/CD ready**: Exit codes, JSON output, timeout handling

---

**Specification Version:** 1.0.3
**Status:** Ready for Implementation
**Date:** 2025-12-30
**Last Updated:** 2025-12-31
