---
summary: Fixed CLI to suppress help text when checks fail, improving output clarity for CI/CD and AI agent workflows
event_type: code
sources:
  - internal/cli/root.go
  - internal/cli/check.go
tags:
  - cli
  - cobra
  - bug-fix
  - vibeguard-qn2
  - ux
---

# Fix for vibeguard-qn2: Silence Usage Text on Check Failures

## Problem

When running `vibeguard check` and a check failed, the CLI incorrectly displayed the usage/help text after the error output:

```
FAIL  lint (error)
      > golangci-lint run ./...

      Tip: golangci-lint found lint issues

Error: exit code 3
Usage:
  vibeguard check [id] [flags]

Flags:
  -h, --help   help for check
...
```

This cluttered the output and was confusing for users and AI agents parsing the output.

## Root Cause

Cobra's default behavior is to show usage text whenever a command returns an error. The `check` command returns an `ExitError` when checks fail, which Cobra treated as any other error.

## Solution

Added `SilenceUsage: true` to the root command in `internal/cli/root.go:34`. This prevents Cobra from showing help text on any error while still allowing:

- Explicit `--help` flag to show usage
- Error messages to be displayed normally

## Testing

1. Built and tested with a failing check - confirmed no usage text shown
2. Verified `--help` flag still works correctly
3. All existing tests pass
4. Ran `vibeguard check` on the project itself - passes

## References

- Beads issue: vibeguard-qn2
- Standard Cobra pattern for production CLIs
