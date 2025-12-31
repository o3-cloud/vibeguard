---
summary: Simplified VibeGuard schema from tools+policies to unified checks with grok pattern support for unstructured output parsing. Reduces config from 30+ lines to 5 lines for common cases.
event_type: deep dive
sources:
  - docs/log/2025-12-30_vibeguard-architecture-spike.md
  - docs/log/2025-12-30_vibeguard-cli-output-design.md
  - https://github.com/elastic/go-grok
tags:
  - vibeguard
  - schema-design
  - simplification
  - grok
  - output-parsing
  - checks
  - developer-experience
---

# Simplified Checks Schema with Grok Support

## Problem

The original spike design required separate `tools` and `policies` sections with OPA/Rego for even simple checks. A basic coverage threshold required 30+ lines of YAML.

## Solution

Unified `checks` schema that:
1. Combines tool execution and assertion in one block
2. Uses exit code as default pass/fail (no config needed)
3. Supports grok patterns for parsing unstructured output
4. Automatically includes command in violation output
5. Optional `suggestion` for additional guidance

## Final Schema

```yaml
checks:
  - id: string              # Required: unique identifier
    run: string             # Required: command to execute

    # Extraction (optional)
    grok: string | [string] # Grok pattern(s) for stdout
    file: path              # Read from file instead of stdout

    # Evaluation
    assert: string          # Optional (default: exit 0)
    severity: error|warning # Optional (default: error)
    suggestion: string      # Optional: extra guidance on failure
    requires: [id, ...]     # Optional: dependencies
    timeout: duration       # Optional (default: 30s)
```

## Examples

### Simple (Exit Code Only)

```yaml
checks:
  - id: vet
    run: go vet ./...

  - id: lint
    run: golangci-lint run
```

### With Grok Extraction

```yaml
checks:
  - id: coverage
    run: go test -cover ./...
    grok: 'coverage: %{NUMBER:coverage}%'
    assert: coverage >= 80
    suggestion: "Run 'go tool cover -func=coverage.out' to find uncovered functions."

  - id: complexity
    run: gocyclo -avg ./...
    grok: 'Average: %{NUMBER:avg}'
    assert: avg <= 15
```

### Multiple Grok Patterns

```yaml
checks:
  - id: security
    run: gosec ./...
    grok:
      - 'High: %{INT:high}'
      - 'Medium: %{INT:medium}'
    assert: high == 0
    severity: error
```

## Output Format

### Failure Without Suggestion

```
FAIL  lint (error)
      > golangci-lint run
```

### Failure With Suggestion

```
FAIL  coverage (error)
      > go test -cover ./...

      Tip: Run 'go tool cover -func=coverage.out' to find uncovered functions.
```

### JSON Output

```json
{
  "violations": [
    {
      "id": "coverage",
      "severity": "error",
      "command": "go test -cover ./...",
      "suggestion": "Run 'go tool cover -func=coverage.out' to find uncovered functions."
    }
  ]
}
```

## Key Design Decisions

### 1. Command Is Always Shown

The violation output always includes the command that failed. This lets agents rerun the command to investigate details. No need to parse/format/truncate tool output.

### 2. Grok Over Raw Regex

Using [elastic/go-grok](https://github.com/elastic/go-grok) for pattern matching because:
- More readable than raw regex
- Built-in patterns: `%{NUMBER}`, `%{INT}`, `%{WORD}`, `%{PATH}`, etc.
- Battle-tested from Logstash ecosystem
- Named captures for clean extraction

### 3. Layered Complexity

| Level | Config | Use Case |
|-------|--------|----------|
| 1 | `run` only | Exit code pass/fail |
| 2 | `run` + `grok` + `assert` | Extract and validate values |
| 3 | Full `tools` + `policies` + Rego | Complex multi-tool analysis |

Simple things stay simple. Power available when needed.

### 4. Suggestion Is Optional

- Command shown automatically (the "what to run")
- Suggestion adds context (the "what else to know")
- Not required — many checks are self-explanatory

## Comparison: Before vs After

### Before (31 lines)

```yaml
tools:
  - id: coverage
    command: ["go", "test", "-coverprofile=coverage.out", "./..."]
    outputs:
      coverage_data: coverage.out

policies:
  - id: coverage-check
    requires: [coverage]
    severity: error
    rego: |
      package vibeguard.coverage
      import rego.v1
      default allow := false
      allow if { input.coverage >= 80 }
      violation contains msg if {
        input.coverage < 80
        msg := sprintf("Coverage %v%% below 80%%", [input.coverage])
      }
      suggestion contains s if {
        count(violation) > 0
        s := "Add unit tests to improve coverage"
      }
```

### After (5 lines)

```yaml
checks:
  - id: coverage
    run: go test -cover ./...
    grok: 'coverage: %{NUMBER:coverage}%'
    assert: coverage >= 80
```

## Grok Patterns Reference

| Pattern | Matches | Example |
|---------|---------|---------|
| `%{NUMBER:n}` | Float/int | `85.5`, `42` |
| `%{INT:n}` | Integer only | `42` |
| `%{WORD:w}` | Single word | `main` |
| `%{DATA:d}` | Non-greedy any | `foo bar` |
| `%{GREEDYDATA:d}` | Greedy any | `everything...` |
| `%{PATH:p}` | File path | `/src/main.go` |
| `%{IP:ip}` | IP address | `192.168.1.1` |

## Fallback Behavior

| Scenario | Behavior |
|----------|----------|
| No `grok`, no `assert` | Pass/fail from exit code |
| `grok` but no match | Warning, fall back to exit code |
| `grok` match + `assert` | Evaluate assertion |
| JSON stdout detected | Auto-parse, query with `json.path` |

## Next Steps

1. Implement `checks` parser that expands to internal tool+policy model
2. Integrate go-grok for pattern extraction
3. Build assertion evaluator for extracted values
4. Keep `tools` + `policies` for advanced users (backwards compatible)

---

**Status:** Design complete, ready for implementation
