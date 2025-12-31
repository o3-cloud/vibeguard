# ADR-004: Establish Code Quality Standards and Tooling

## Status
Proposed

## Context and Problem Statement

As the VibeGuard project matures and scales with multiple contributors, maintaining consistent code quality, readability, and reliability becomes critical. Without established standards, we risk:

- Inconsistent code style making reviews harder
- Undetected bugs due to insufficient static analysis
- Unpredictable test coverage
- Unclear or missing documentation
- Developer friction from manual quality checks

We need a comprehensive, automated approach to code quality that enforces standards locally during development while keeping the review experience frictionless.

## Considered Options

### Option A: Minimal tooling (gofmt only)
- Use only Go's built-in `gofmt` for formatting
- Rely on manual code review for quality
- **Pros**: Zero setup, no dependencies
- **Cons**: Limited bug detection, inconsistent patterns, manual enforcement

### Option B: Comprehensive automated tooling with pre-commit hooks (SELECTED)
- Use golangci-lint for linting, goimports for formatting, go vet for analysis
- Enforce standards via pre-commit hooks locally
- Establish test coverage and documentation expectations
- Provide IDE integration guides
- **Pros**: Catches issues early, consistent enforcement, developer experience focused
- **Cons**: Setup overhead, requires pre-commit hook discipline

### Option C: CI/CD-only enforcement
- Enforce all standards only in CI/CD pipeline
- No local tooling requirements
- **Pros**: No local setup needed
- **Cons**: Slower feedback loop, rejected PRs due to formatting, frustrating developer experience

## Decision Outcome

**Chosen option:** Option B - Comprehensive automated tooling with pre-commit hooks

**Rationale:**
- Shift-left approach catches issues early, before code review
- Pre-commit hooks provide immediate feedback during development
- golangci-lint offers comprehensive linting with no configuration needed (uses sensible defaults)
- goimports ensures consistent import ordering and formatting
- Local enforcement reduces CI/CD failures and improves developer velocity
- Standards are explicit and discoverable in the codebase

**Tradeoffs:**
- Requires developers to configure pre-commit hooks
- Setup cost upfront, but saved through faster reviews and fewer corrections
- Slightly slower commits due to checks, but prevents rework

## Consequences

### Positive Outcomes
1. **Consistent code style** - All code follows the same formatting conventions
2. **Early bug detection** - Static analysis catches issues before review
3. **Reduced review friction** - Reviewers focus on logic, not style
4. **Test quality** - Coverage requirements ensure tests are comprehensive
5. **Maintainability** - Clear documentation standards make code easier to understand
6. **Developer confidence** - Automated checks catch mistakes before submission

### Negative Outcomes
1. **Setup friction** - Developers must configure pre-commit hooks initially
2. **Commit time** - Hooks add a few seconds to each commit
3. **False positives** - Occasionally linter rules may require suppression with comments
4. **Dependency on tools** - Must maintain linter configurations as tools evolve

### Neutral Impacts
1. **Documentation** - Will need CONTRIBUTING.md with setup instructions
2. **Configuration files** - Adds `.golangci.yaml`, `.pre-commit-config.yaml`, `.editorconfig`

## Implementation Details

### 1. Go Code Formatting
- **Tool**: `goimports` (via golangci-lint)
- **Enforced**: Auto-format via pre-commit hook
- **Standard**: Follows Go conventions, sorts imports into stdlib, external, internal groups
- **Configuration**: None needed - uses Go defaults

### 2. Code Linting
- **Tool**: `golangci-lint` with strict configuration
- **Enforced**: Pre-commit hook (non-blocking), optional in IDE
- **Included linters**:
  - `vet` - Go's built-in static analyzer
  - `revive` - Go's linter alternative to golint
  - `gosimple` - Suggest code simplifications
  - `goerr113` - Proper error handling patterns
  - `noctx` - Detect missing context in function calls
  - `errorlint` - Better error wrapping patterns
  - Plus 20+ other critical checks
- **Configuration**: `.golangci.yaml` with sensible defaults
- **Suppression**: Via inline `//nolint` comments with justification

### 3. Static Analysis
- **Tool**: `go vet` (included in golangci-lint)
- **Enforced**: Pre-commit hook
- **Checks**: Unreachable code, shadowed variables, incorrect format strings, etc.

### 4. Testing Standards
- **Coverage requirement**: Minimum 70% coverage for all packages
- **Test naming**: Follow convention `Test{FunctionName}_{Scenario}` (e.g., `TestEvaluate_WithNilPolicy`)
- **Table-driven tests**: Preferred pattern for multiple test cases
- **Test execution**: `go test ./...` with coverage reporting
- **Enforced**: Pre-commit hook runs tests

### 5. Documentation Standards
- **Exported symbols**: All exported functions, types, and constants require doc comments
- **Doc comment format**: First line is a complete sentence starting with symbol name
  ```go
  // Runner evaluates policies against resources.
  type Runner interface { ... }

  // Evaluate runs a policy against a resource and returns the result.
  func (r *Runner) Evaluate(ctx context.Context, ...) { ... }
  ```
- **Package documentation**: Each package should have a package-level comment
- **Complex logic**: Inline comments explain "why", not "what"
- **Enforced**: Code review, linter checks for undocumented exports

### 6. Pre-Commit Hooks
- **Tool**: `pre-commit` framework
- **Hooks configured**:
  - Run goimports to auto-format imports
  - Run golangci-lint for linting
  - Run go test for tests
  - Prevent commits with lint errors or failed tests
- **Configuration**: `.pre-commit-config.yaml` in repo root
- **Setup**: `pre-commit install` after cloning
- **Bypass**: `git commit --no-verify` only for emergencies (discouraged)

### 7. IDE Integration
- **VS Code**: Install `golangci-lint` extension (uses project's config)
- **GoLand/IntelliJ**: Configure runner with project's `.golangci.yaml`
- **Vim/Neovim**: Use ALE or nvim-lint with golangci-lint
- **All IDEs**: Configure gofmt on save

### 8. Editor Configuration
- **Tool**: `.editorconfig` for consistent editor settings across IDEs
- **Settings**: Indentation, line endings, trailing whitespace, final newlines
- **Support**: VS Code, GoLand, Vim, and all major editors

## Quality Gates

- ✅ All tests must pass
- ✅ No golangci-lint violations
- ✅ No vet warnings
- ✅ All exported symbols documented
- ✅ Imports formatted with goimports
- ✅ Code follows Go conventions

## Related Documentation

- Contributing guide (to be created)
- `.golangci.yaml` configuration file
- `.pre-commit-config.yaml` configuration file
- `.editorconfig` configuration file

## Alternatives Evaluated

See "Considered Options" section above for detailed analysis of alternatives.
