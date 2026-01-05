---
summary: Added integration tests for main binary error handling and exit codes
event_type: code
sources:
  - cmd/vibeguard/main.go
  - cmd/vibeguard/main_integration_test.go
  - internal/cli/check.go
  - internal/config/config.go
tags:
  - testing
  - integration-tests
  - error-handling
  - exit-codes
  - coverage
---

# Main Binary Integration Tests

Implemented integration tests for the main binary (`cmd/vibeguard/main.go`) to cover error handling paths that were previously untested (0% coverage).

## Problem

The main binary error handling logic was untested:
- `ExitError` detection and exit code propagation
- `ConfigError` detection resulting in exit code 2
- Default exit code 1 for generic errors
- Custom exit code via `--error-exit-code` flag

## Solution

Created `cmd/vibeguard/main_integration_test.go` with comprehensive tests that:

1. Build the actual binary in a temp directory
2. Run the binary with various configurations
3. Verify correct exit codes are returned

## Test Coverage

### Exit Code Tests
- **Exit 0**: Successful check passes
- **Exit 1**: Failing check with default error code
- **Exit 2**: Config errors (missing file, invalid YAML, no checks, unsupported version, duplicate IDs, invalid ID format, cyclic dependencies)
- **Custom exit code**: Via `--error-exit-code` flag
- **Warning severity**: Does not cause non-zero exit

### Output Tests
- Config errors write messages to stderr
- Successful checks produce no output (silence is success)

## Key Implementation Details

- Tests use `exec.Command` to run the actual compiled binary
- Binary is built fresh per test run using `go build`
- Temp directories are used for isolation
- Helper functions: `buildTestBinary`, `writeConfig`, `assertExitCode`

## Related Issue

Closes vibeguard-pm5: "Add main binary integration tests for error handling paths"
