---
summary: Implemented --version flag for vibeguard CLI to display semantic version
event_type: code
sources:
  - internal/version/version.go
  - internal/cli/root.go
  - cmd/vibeguard/main.go
tags:
  - cli
  - feature
  - version-management
  - semantic-versioning
---

# Vibeguard --version Flag Implementation

## Task
Implemented the `--version` flag for the vibeguard CLI to allow users to verify the installed version of vibeguard. This addresses requirement vibeguard-qo4.

## Implementation Details

### New Package: `internal/version`
Created a new `version` package (`internal/version/version.go`) that:
- Exports a `Version` variable with default value `v0.1.0-dev`
- Supports override at build time using `-ldflags`
- Provides `String()` function for accessing the version

This pattern allows semantic versioning and CI/CD integration where version is injected at build time:
```bash
go build -ldflags="-X github.com/vibeguard/vibeguard/internal/version.Version=v1.0.0"
```

### Modified: `internal/cli/root.go`
Updated the root command to:
1. Import the new `version` package and `fmt`
2. Add `showVersion` boolean flag variable
3. Implement `Run` handler on `rootCmd` that:
   - Checks if `--version` flag is set
   - Prints version string if true
   - Shows help text if no flags provided
4. Register `--version` flag in init function using `rootCmd.Flags()` (not PersistentFlags to avoid inheritance to subcommands)

## Testing
- Built binary successfully with `go build`
- Tested `--version` flag outputs `v0.1.0-dev`
- Tested version override at build time produces correct output
- All existing tests pass (98+ tests across packages)
- Help text still displays when called without arguments

## Key Design Decisions
1. **Flag Location**: Used `Flags()` instead of `PersistentFlags()` so `--version` is only available on root command, not inherited by subcommands
2. **Version Storage**: Dedicated package `internal/version` keeps versioning concern separated from CLI logic
3. **Build-time Injection**: Uses standard Go ldflags pattern for CI/CD flexibility
4. **Exit Behavior**: Returns cleanly after printing version without error

## Files Modified
- `internal/cli/root.go` - Added version flag and handler
- `internal/version/version.go` - New file with version constant

## Verification
```
$ ./bin/vibeguard --version
v0.1.0-dev

$ go build -ldflags="-X github.com/vibeguard/vibeguard/internal/version.Version=v1.0.0" -o bin/vibeguard-test ./cmd/vibeguard
$ ./bin/vibeguard-test --version
v1.0.0
```

## Task Status
✅ vibeguard-qo4 closed successfully
