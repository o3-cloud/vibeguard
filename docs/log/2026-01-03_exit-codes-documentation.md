---
summary: Documented exit codes for CI/CD integration in README
event_type: code
sources:
  - internal/executor/executor.go
  - README.md
tags:
  - documentation
  - exit-codes
  - ci-cd
  - integration
  - readme
---

# Exit Codes Documentation

## Completed Task: vibeguard-x1l

Added comprehensive documentation of VibeGuard exit codes to the README for CI/CD integration users.

## Changes Made

### Exit Codes Section Added to README
- **Location**: Between "CLI Reference" and "Configuration Schema" sections
- **Content**:
  - Table of all exit codes with descriptions
  - CI/CD integration guidance
  - GitHub Actions example

### Documented Exit Codes
| Exit Code | Name | Description |
|-----------|------|-------------|
| 0 | Success | All checks passed successfully |
| 2 | ConfigError | Configuration file error (invalid YAML, validation failure, etc.) |
| 3 | Violation | One or more error-severity violations detected during execution |
| 4 | Timeout | Check execution error (timeout exceeded, command not found, etc.) |

## Key Findings

1. **Well-Designed Exit Code Schema**: The exit codes in `internal/executor/executor.go` are thoughtfully designed for CI/CD compatibility where exit codes ≥2 are blocking, aligning with standards like Claude Code hooks.

2. **Documentation Gap**: Despite exit codes being clearly defined in the codebase (lines 14-20), no user-facing documentation existed. This is a critical gap for CI/CD integration use cases.

3. **CI/CD Integration Context**: Added context explaining how exit codes control pipeline behavior:
   - Exit code 0 allows pipeline to proceed
   - Exit codes ≥2 block the pipeline (suitable for pre-commit hooks and CI checks)

## Commit

- **Hash**: 7cf5bf8
- **Message**: `docs: add Exit Codes section to README for CI/CD integration`
- **Files Changed**: README.md (+26 lines)

## Next Steps

- Exit codes are now discoverable and actionable for CI/CD users
- Consider adding exit code documentation to additional integration guides (INTEGRATIONS.md, if created per vibeguard-wo8)
- Related task: vibeguard-wo8 (Create INTEGRATIONS.md with CI/CD examples) could reference this section

## Related ADRs

- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement (uses exit code design)
- ADR-005: Adopt Vibeguard for Policy Enforcement in CI/CD (context for this documentation)
