---
summary: Added language-specific static analysis checks to all templates
event_type: code
sources:
  - docs/specs/init-template-system-spec.md
  - internal/cli/templates/
tags:
  - static-analysis
  - templates
  - code-quality
  - go
  - javascript
  - python
  - rust
  - vibeguard-693
---

# Static Analysis Checks for All Templates

## Task Completion: vibeguard-693
Successfully implemented language-specific static analysis checks across all templates for deeper code quality insights beyond standard linting.

## Implementation Summary

### Go Templates
1. **go-minimal** - Added staticcheck
   - Command: `staticcheck {{.go_packages}}`
   - Severity: warning
   - Provides deeper static analysis than go vet
   - Depends on: vet check

2. **go-standard** - Added staticcheck
   - Command: `staticcheck {{.go_packages}}`
   - Severity: warning
   - Complements golangci-lint for comprehensive coverage
   - Depends on: lint check

### Node.js Templates
3. **node-javascript** - Added analyze check
   - Command: `npm run analyze`
   - Severity: warning
   - Allows projects to configure their own static analysis tools
   - Suggests: sonarjs, code-inspector
   - Depends on: lint check

4. **node-typescript** - Added analyze check
   - Command: `npm run analyze`
   - Severity: warning
   - Same pattern as JavaScript for consistency
   - Depends on: lint check

### Python Templates
5. **python-pip** - Added pylint
   - Command: `pylint {{.source_dir}} --disable=all --enable=E,F`
   - Severity: warning
   - Focuses on errors and fatals, not style warnings
   - Complements ruff for comprehensive analysis
   - Depends on: lint check

6. **python-poetry** - Added pylint
   - Command: `poetry run pylint {{.source_dir}} --disable=all --enable=E,F`
   - Uses Poetry for dependency management
   - Depends on: lint check

7. **python-uv** - Added pylint
   - Command: `uv run pylint {{.source_dir}} --disable=all --enable=E,F`
   - Uses uv for fast dependency resolution
   - Depends on: lint check

### Rust Template
8. **rust-cargo** - Added cargo-deny
   - Command: `cargo deny check`
   - Severity: warning
   - Analyzes dependencies for security/license issues
   - Complements clippy for comprehensive coverage
   - Depends on: clippy check

## Design Decisions

### Graceful Degradation
All analyze checks use the `|| echo "... not installed (optional)"` pattern to:
- Fail gracefully when tools aren't installed
- Not block the build pipeline
- Encourage optional tool adoption
- Provide helpful installation instructions

### Tool Selection

**Go**: `staticcheck` provides detailed static analysis beyond `go vet`:
- Detects unused code
- Identifies likely bugs
- Finds performance issues
- Available as single binary

**Node.js**: Delegates to project-defined `npm run analyze`:
- Respects project preferences
- Supports multiple tools (sonarjs, code-inspector, etc.)
- Allows gradual adoption

**Python**: `pylint` with focused scope (`--enable=E,F`):
- Focuses on errors and fatals
- Avoids conflicts with ruff
- More comprehensive than ruff for certain issues
- Requires explicit setup per package manager

**Rust**: `cargo-deny` for dependency analysis:
- Checks for security vulnerabilities
- Validates licensing
- Complementary to clippy's code analysis

### Dependency Structure
All analyze checks come after core linting checks to:
- Ensure code quality before deeper analysis
- Fail fast on preventable issues
- Keep check ordering logical

## Testing Results
- All 9 templates pass validation tests
- No syntax errors in YAML configurations
- vibeguard check completes without errors
- No regressions in existing functionality

## Check Ordering (Full Pipeline)

Example from go-minimal:
1. fmt (formatting)
2. vet (built-in static analysis)
3. **analyze** (staticcheck - deeper analysis) ← NEW
4. test (unit tests)
5. coverage (code coverage threshold)
6. build (compilation)

## Installation Guidance
Each analyze check includes helpful installation instructions:
- Go: `go install honnef.co/go/tools/cmd/staticcheck@latest`
- Python (pip): `pip install pylint`
- Python (poetry): `poetry add --group dev pylint`
- Python (uv): `uv pip install pylint`
- Rust: `cargo install cargo-deny`

## Future Considerations
- Monitor adoption of analyze checks in real projects
- Consider adding more specialized tools (e.g., Go race detector enhancement)
- Evaluate optional tool packaging with templates
- Consider analyze tooling for generic template

## Alignment with Project Standards
- Follows ADR-004: Code Quality Standards
- Enables comprehensive quality validation
- Supports multi-language development
- Consistent with template design principles

## Notes
- All analyze checks are warnings, not errors (optional enforcement)
- Installation instructions provided in suggestions
- Tool compatibility checked with `2>/dev/null || ...` pattern
- Consistent 120s timeout for all analyze checks
