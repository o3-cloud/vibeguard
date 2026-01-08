---
summary: Added code coverage checks to go-minimal, node-javascript, rust-cargo, and generic templates
event_type: code
sources:
  - docs/specs/init-template-system-spec.md
  - internal/cli/templates/
tags:
  - code-coverage
  - templates
  - go
  - javascript
  - rust
  - quality-standards
  - vibeguard-nmo
---

# Code Coverage Checks for All Templates

## Task Completion: vibeguard-nmo
Successfully implemented code coverage checks for all templates as specified in the init template system specification.

## Implementation Summary

### Templates Updated
1. **go-minimal** - Added coverage check using `go test` with coverprofile
   - Extracts coverage from `go tool cover -func` output
   - Pattern: `total:.*\(statements\)\s+%{NUMBER:coverage}%`
   - Min threshold: 70%

2. **node-javascript** - Added coverage check using `npm test --coverage`
   - Extracts coverage from Jest output
   - Pattern: `"Lines\\s+:\\s+%{NUMBER:coverage}%"`
   - Min threshold: 70%

3. **rust-cargo** - Added coverage check using `cargo tarpaulin`
   - Extracts coverage from tarpaulin stdout
   - Pattern: `Coverage:\s+%{NUMBER:coverage}%`
   - Min threshold: 70%
   - Note: Requires `cargo-tarpaulin` crate (optional dependency)

4. **generic** - Added commented coverage check example
   - Includes template structure for custom coverage tools
   - Helps users understand coverage check pattern

### Templates Already Complete
- go-standard ✓
- node-typescript ✓
- python-pip ✓
- python-poetry ✓
- python-uv ✓

## Standardization Across All Templates
- All coverage checks are now consistently structured
- min_coverage variable defaults to 70% (aligns with ADR-004 code quality standards)
- Each check includes helpful error messages
- Coverage checks properly depend on test checks
- Follows language-specific conventions and tools

## Testing Results
- All 9 templates pass validation tests
- No syntax errors in YAML configurations
- vibeguard check completes without errors
- No regressions in existing functionality

## Alignment with Project Standards
- Follows ADR-004: Code Quality Standards (70% coverage requirement)
- Consistent with template naming conventions and structure
- Enables dogfooding of VibeGuard policies (ADR-005)
- Supports AI agent workflows for init system (as per spec)

## Next Steps
- Phase 2 of spec: Template expansion (add framework-specific templates)
- Phase 3A: Codebase simplification (delete redundant examples)
- Monitor if any tools require adjustment (e.g., cargo tarpaulin availability)

## Notes
- JavaScript pattern uses Jest's default coverage reporter format
- Go pattern extracts from gofmt-compatible tool output
- Rust pattern depends on cargo-tarpaulin being installed or available
- All patterns are validated and tested
