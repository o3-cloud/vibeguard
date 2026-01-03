---
summary: Security review of shell command injection in variable interpolation - concluded no vulnerability exists due to intentional design where config author controls both variables and commands
event_type: deep dive
sources:
  - internal/config/interpolate.go
  - internal/config/interpolate_test.go
  - internal/executor/executor.go
tags:
  - security
  - shell-injection
  - variable-interpolation
  - threat-model
  - security-review
---

# Shell Command Injection Security Review

Completed security review for beads issue `vibeguard-8hs`: "Review shell command injection prevention in variable interpolation"

## Issue Background

Commands in vibeguard are executed via `sh -c` with variable interpolation. The concern was that variables could enable command injection if they contain shell metacharacters like `;`, `|`, `$()`, or backticks.

## Security Analysis

### Variable Sources (Trust Boundaries)

1. **Config file variables (`vars`)**: Defined by config author in `vibeguard.yaml`
   - Same person controls both variable values AND the commands that use them
   - No injection risk - they can already put malicious commands directly in `run:`

2. **Grok-extracted values**: Come from command output
   - ONLY used for display purposes (suggestions, fix messages)
   - NEVER used in command execution
   - Verified in `internal/output/formatter.go` - `InterpolateWithExtracted` is only called for suggestion/fix rendering

### Key Finding: No Vulnerability

The "vulnerability" is actually by design:

1. **Config author == trust boundary**: The person writing `vibeguard.yaml` is the same person defining what commands to run. If malicious, they can already execute arbitrary commands.

2. **No external input**: Variables are defined statically in the YAML file, not from environment variables, user input, or command line arguments at runtime.

3. **Escaping would break legitimate use cases**:
   ```yaml
   vars:
     packages: "./cmd/... ./internal/..."
   checks:
     - run: go test {{.packages}}
   ```
   Escaping dots, slashes, or other shell characters would break this primary use case.

## Code Flow Analysis

```
vibeguard.yaml (vars section)
        ↓
config.Load() - parses YAML
        ↓
config.Interpolate() - replaces {{.var}} with values from vars map
        ↓
orchestrator.Run() - schedules checks
        ↓
executor.Execute() - runs via sh -c <interpolated_command>
```

The only source of variable values is the static YAML config file. No external/untrusted input enters this flow.

## Changes Made

1. **Added security model documentation** to `interpolate.go` explaining the design decisions

2. **Added comprehensive tests** to `interpolate_test.go`:
   - `TestInterpolateShellMetacharacters` - documents that metacharacters pass through by design
   - `TestInterpolateWithExtractedSecurityNote` - confirms grok values are display-only

## Conclusion

**No vulnerability exists.** The security model is intentional:
- Trust boundary is at the config file level
- Config author controls both variables and commands
- Shell metacharacter escaping would be security theater that breaks legitimate use cases

## Related

- Closed beads issue: `vibeguard-8hs`
- Files modified: `internal/config/interpolate.go`, `internal/config/interpolate_test.go`
