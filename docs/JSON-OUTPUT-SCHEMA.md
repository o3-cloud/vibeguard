# VibeGuard JSON Output Schema

VibeGuard can output check results in JSON format, which is useful for programmatic parsing, integration with other tools, and detailed analysis. This document describes the complete JSON output schema.

## Quick Start

Run VibeGuard with the `--json` flag to get JSON output:

```bash
vibeguard check --json
```

## Output Structure

The JSON output is a single object with the following top-level structure:

```json
{
  "checks": [...],
  "violations": [...],
  "exit_code": 0,
  "fail_fast_triggered": false
}
```

## Top-Level Fields

| Field | Type | Description |
|-------|------|-------------|
| `checks` | array | Array of check execution results |
| `violations` | array | Array of policy violations detected |
| `exit_code` | integer | Exit code indicating overall result (0=success, 1=failure/timeout by default, 2=config error) |
| `fail_fast_triggered` | boolean | Whether execution stopped early due to `--fail-fast` flag (omitted if false) |

## Check Object

Each object in the `checks` array represents the execution result of a single check:

```json
{
  "id": "fmt",
  "status": "passed",
  "duration_ms": 150
}
```

### Check Fields

| Field | Type | Description | Values |
|-------|------|-------------|--------|
| `id` | string | The check's unique identifier (from config) | any string |
| `status` | string | The execution status of the check | `"passed"`, `"failed"`, `"cancelled"` |
| `duration_ms` | integer | How long the check took to execute in milliseconds | >= 0 |

### Status Values

- **`passed`** — Check executed successfully and passed all assertions
- **`failed`** — Check executed but failed its assertions or produced errors
- **`cancelled`** — Check execution was cancelled (typically due to timeout or `--fail-fast`)

## Violation Object

Each object in the `violations` array represents a policy violation detected during a failed check:

```json
{
  "id": "coverage",
  "severity": "error",
  "command": "go test -cover ./...",
  "suggestion": "Coverage is 72%, need 80%. Run 'go test ./...' with coverage analysis.",
  "fix": "Add unit tests to improve coverage",
  "extracted": {
    "coverage": "72"
  }
}
```

### Violation Fields

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `id` | string | The check ID that produced this violation | Yes |
| `severity` | string | Severity level of the violation | Yes |
| `command` | string | The command that was executed | Yes |
| `suggestion` | string | Actionable suggestion for fixing the issue | No |
| `fix` | string | Interpolated fix instructions from the config | No |
| `extracted` | object | Data extracted from command output via grok patterns | No |

### Severity Values

- **`"error"`** — Critical violation; indicates a failed check that must be addressed
- **`"warning"`** — Non-critical issue; check failed but doesn't block execution in non-strict mode

### Extracted Data

The `extracted` field is an object where:
- **Keys** are grok pattern group names from the configuration
- **Values** are strings containing the matched values

Example:
```yaml
# Config
checks:
  - id: test-coverage
    run: go test ./... -cover
    grok:
      - total:.*\(statements\)\s+%{NUMBER:coverage}%
```

Would produce:
```json
"extracted": {
  "coverage": "82"
}
```

## Complete Examples

### All Checks Passing

```json
{
  "checks": [
    {
      "id": "fmt",
      "status": "passed",
      "duration_ms": 150
    },
    {
      "id": "vet",
      "status": "passed",
      "duration_ms": 320
    }
  ],
  "violations": [],
  "exit_code": 0
}
```

### Checks with Violations

```json
{
  "checks": [
    {
      "id": "fmt",
      "status": "passed",
      "duration_ms": 150
    },
    {
      "id": "coverage",
      "status": "failed",
      "duration_ms": 900
    }
  ],
  "violations": [
    {
      "id": "coverage",
      "severity": "error",
      "command": "go test -cover ./...",
      "suggestion": "Coverage is 72%, need 80%. Run 'go test ./...' with comprehensive tests.",
      "fix": "Add unit tests to improve coverage to 80%",
      "extracted": {
        "coverage": "72"
      }
    }
  ],
  "exit_code": 1
}
```

### Timeout Scenario

```json
{
  "checks": [
    {
      "id": "fmt",
      "status": "passed",
      "duration_ms": 150
    },
    {
      "id": "slow-test",
      "status": "cancelled",
      "duration_ms": 30000
    }
  ],
  "violations": [],
  "exit_code": 1,
  "fail_fast_triggered": true
}
```

### Multiple Violations from One Check

A single check can produce multiple violations if grok patterns extract multiple values:

```json
{
  "checks": [
    {
      "id": "quality-metrics",
      "status": "failed",
      "duration_ms": 1200
    }
  ],
  "violations": [
    {
      "id": "quality-metrics",
      "severity": "error",
      "command": "custom-quality-tool",
      "suggestion": "Code coverage below threshold",
      "extracted": {
        "coverage": "65",
        "cyclomatic_complexity": "12"
      }
    }
  ],
  "exit_code": 1
}
```

## Exit Code Mapping

The `exit_code` field corresponds to VibeGuard's exit codes:

| Value | Name | Description |
|-------|------|-------------|
| 0 | Success | All checks passed, no violations |
| 1 | Error (default) | One or more violations detected, or timeout (configurable via `--error-exit-code`) |
| 2 | ConfigError | Configuration parsing or validation failed |

### Configurable Error Exit Code

By default, VibeGuard uses exit code 1 for both check failures (error-severity violations) and timeouts. This can be customized using the `--error-exit-code` flag:

```bash
# Use exit code 3 for failures (matches legacy behavior)
vibeguard check --error-exit-code=3
```

**Note:** Exit codes 3 and 4 were previously used to distinguish violations from timeouts. This behavior is now deprecated in favor of a unified configurable exit code. The JSON output's `checks` array still indicates whether a check timed out via its status field.

## Using JSON Output Programmatically

### With jq (Command-line JSON Processor)

Extract all violation IDs:
```bash
vibeguard check --json | jq '.violations[].id'
```

Get the exit code:
```bash
vibeguard check --json | jq '.exit_code'
```

Check if any violations have "error" severity:
```bash
vibeguard check --json | jq '.violations[] | select(.severity == "error")'
```

### With Python

```python
import json
import subprocess

result = subprocess.run(['vibeguard', 'check', '--json'], capture_output=True, text=True)
data = json.loads(result.stdout)

if data['exit_code'] != 0:
    for violation in data['violations']:
        print(f"Check {violation['id']}: {violation['suggestion']}")
```

### With Go

```go
var output map[string]interface{}
if err := json.Unmarshal([]byte(jsonOutput), &output); err != nil {
    log.Fatal(err)
}

exitCode := int(output["exit_code"].(float64))
violations := output["violations"].([]interface{})
```

## Integration with CI/CD

JSON output is particularly useful in CI/CD pipelines where you want to:

1. **Parse detailed results** — Extract specific violation information
2. **Generate reports** — Convert to HTML, markdown, or custom formats
3. **Conditional actions** — Take different actions based on violation severity
4. **Metrics collection** — Track check duration and failure rates over time

### Example: GitHub Actions with JSON Parsing

```yaml
- name: Run VibeGuard checks
  id: vibeguard
  run: vibeguard check --json > results.json

- name: Parse and report results
  if: always()
  run: |
    VIOLATIONS=$(jq '.violations | length' results.json)
    if [ "$VIOLATIONS" -gt 0 ]; then
      echo "Found $VIOLATIONS violations:"
      jq '.violations[] | "\(.id): \(.suggestion)"' results.json
    fi

- name: Fail if critical violations
  if: failure()
  run: exit 1
```

## Compatibility Notes

- JSON output is produced with 2-space indentation for readability
- All string fields use UTF-8 encoding
- The output is a single JSON object (not an array)
- Empty arrays (e.g., no violations) are included in the output
- The `fail_fast_triggered` field is only included when `true`
- Field ordering within objects is not guaranteed; rely on field names

## Notes for Consumers

- **Always check `exit_code`** — Don't rely solely on the presence/absence of violations
- **Handle missing optional fields** — `suggestion`, `fix`, and `extracted` may be absent
- **Duration precision** — Millisecond precision may vary slightly depending on system load
- **Extracted data encoding** — Values in `extracted` are always strings (no type inference)
- **Error handling** — If JSON output fails to serialize, VibeGuard will exit with code 4 and output error text

## See Also

- [CLI Reference](README.md#cli-reference) — General command-line options
- [Configuration Schema](README.md#configuration-schema) — How to define checks with grok patterns and assertions
- [Assertion Expression Operators](ASSERTION-OPERATORS.md) — How to write assertion conditions
