---
summary: Completed comprehensive JSON output schema documentation for vibeguard-icf task
event_type: code
sources:
  - internal/output/json.go
  - internal/output/json_test.go
  - internal/executor/executor.go
  - docs/INTEGRATIONS.md
tags:
  - documentation
  - json-output
  - schema
  - api
  - vibeguard-icf
  - task-completion
---

# JSON Output Schema Documentation Complete

## Overview

Completed task **vibeguard-icf**: Document JSON output schema. Created comprehensive documentation for VibeGuard's JSON output format, which is essential for programmatic integration with CI/CD systems and other tools.

## Deliverable

**File created:** `docs/JSON-OUTPUT-SCHEMA.md`

This document provides:

### Structure Documentation
- Complete JSON output structure with all top-level fields
- Field types, descriptions, and default values
- Exit code mapping (0, 2, 3, 4) with descriptions

### Object Schemas
- **Check objects**: id, status (passed/failed/cancelled), duration_ms
- **Violation objects**: id, severity, command, suggestion, fix, extracted data
- All optional and required fields clearly marked

### Practical Examples
1. All checks passing - baseline scenario
2. Checks with violations - error detection
3. Timeout scenario - fail-fast behavior
4. Multiple violations - complex grok extraction

### Integration Guides
- **jq examples** - Command-line JSON processing
- **Python example** - Programmatic parsing
- **Go example** - Type-safe unmarshaling
- **GitHub Actions** - CI/CD pipeline integration

### Consumer Guidelines
- Exit code checking recommendations
- Handling optional fields
- Duration precision expectations
- String encoding notes
- Error handling patterns

## Verification

Documentation was verified against actual implementation:
- Checked `internal/output/json.go` - FormatJSON function implementation
- Reviewed `internal/output/json_test.go` - 5 comprehensive test cases
- Validated exit codes from `internal/executor/executor.go`
- Ran actual `vibeguard check --json` output and verified it matches schema

All documented fields, types, and examples match the actual implementation exactly.

## Integration

Updated `README.md` to link to the new schema documentation in the `vibeguard check` command section, providing easy access for users who need JSON output details.

## Next Steps

The documentation is complete and ready for use. Other related open tasks:
- vibeguard-aq5: Document assertion expression operators
- vibeguard-rol: Document file config field
- vibeguard-nha: Document grok pattern debugging

These documentation tasks could benefit from similar comprehensive schema documentation following the same pattern.
