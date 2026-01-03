# Contributing to VibeGuard

Thank you for your interest in contributing to VibeGuard! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please be respectful and inclusive in all interactions. We're committed to providing a welcoming and supportive environment for all contributors.

## Getting Started

### Prerequisites

- **Go 1.21+** — Latest stable Go version
- **git** — For version control
- **golangci-lint** — For code linting
- **gofmt** — For code formatting (included with Go)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/vibeguard.git
   cd vibeguard
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/vibeguard/vibeguard.git
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Verify setup by running tests:
   ```bash
   go test -v ./...
   ```

## Development Workflow

### Creating a Feature Branch

Create a feature branch from `main`:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-timeout-support` — Adding a new feature
- `fix/panic-on-empty-config` — Fixing a bug
- `docs/update-readme` — Documentation changes
- `refactor/simplify-executor` — Code refactoring

### Making Changes

1. Write code following the conventions in `CONVENTIONS.md`
2. Write tests for new functionality:
   ```bash
   go test -v ./...
   ```
3. Run linting and formatting:
   ```bash
   gofmt -w .
   go vet ./...
   golangci-lint run ./...
   ```
4. Run VibeGuard on itself to verify your changes:
   ```bash
   vibeguard check
   ```

### Code Quality Standards

- **Test Coverage:** Aim for 70%+ code coverage. Run `go test -cover ./...` to check coverage.
- **Code Style:** Follow the conventions in `CONVENTIONS.md`
- **Comments:** Add comments for non-obvious logic. Document public APIs.
- **Error Handling:** Always handle errors explicitly. Use meaningful error messages.
- **Performance:** Consider performance implications for frequently-used code paths.

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): short description (imperative mood)

Optional longer explanation. Reference issues if applicable.

Fixes #123
```

**Types:**
- `feat` — New feature
- `fix` — Bug fix
- `docs` — Documentation changes
- `refactor` — Code refactoring (no functional changes)
- `test` — Test additions or changes
- `chore` — Build, dependency, or tooling changes

**Examples:**
```
feat(config): add variable interpolation support

Implement {{.var}} syntax for injecting variables into check commands.
Variables can be defined in the top-level vars: section.

Fixes #42
```

```
fix(executor): handle timeout correctly for long-running checks

Previously, timeouts were being applied incorrectly when context was
cancelled. Now properly waiting for process termination.
```

### Submitting a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a pull request on GitHub with:
   - Clear description of what changed and why
   - Reference to any related issues (e.g., "Fixes #123")
   - Screenshots/examples if the change is user-visible

3. Your PR will be reviewed. Address feedback by pushing additional commits.

4. Once approved, your PR will be merged by a maintainer.

## Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v ./... -cover

# Run a specific test
go test -v ./... -run TestCheckExecutor

# Run tests in a specific package
go test -v ./internal/config
```

### Writing Tests

1. Create test files with `_test.go` suffix in the same package
2. Use table-driven tests for multiple scenarios:
   ```go
   func TestConfigParsing(t *testing.T) {
       tests := []struct {
           name    string
           input   string
           want    *Config
           wantErr bool
       }{
           {
               name:  "valid config",
               input: validYAML,
               want:  expectedConfig,
           },
           {
               name:    "invalid version",
               input:   invalidVersionYAML,
               wantErr: true,
           },
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test implementation
           })
       }
   }
   ```

3. Aim for meaningful test names that describe what's being tested
4. Include both happy path and error cases

## Architecture Decisions

Major architectural decisions are documented in `docs/adr/` using the MADR format. Review these documents to understand design rationale and constraints:

- ADR-001 — Adopt Beads for AI Agent Task Management
- ADR-002 — Adopt Conventional Commits
- ADR-003 — Adopt Go as the Primary Implementation Language
- ADR-004 — Code Quality Standards and Tooling
- ADR-005 — Adopt VibeGuard for Policy Enforcement
- ADR-006 — Integrate VibeGuard as Git Pre-Commit Hook
- ADR-007 — Adopt Gremlins for Mutation Testing

When making architectural changes, consider creating a new ADR to document your decision.

## Reporting Issues

### Bug Reports

Include:
- Clear description of the bug
- Steps to reproduce
- Expected behavior vs. actual behavior
- VibeGuard version (`vibeguard --version`)
- OS and Go version
- Configuration that triggers the issue

Example:
```
**Bug:** Coverage check fails with grok pattern matching error

**Steps to reproduce:**
1. Create vibeguard.yaml with coverage check
2. Run `vibeguard check`

**Expected:** Coverage percentage is extracted and assertion is evaluated

**Actual:** Error: "Invalid grok pattern syntax"

**VibeGuard version:** 0.1.0
**OS:** macOS 14.1
**Go version:** 1.21
```

### Feature Requests

Include:
- Use case: Why do you need this feature?
- Proposed solution: How should it work?
- Alternatives: Any alternative approaches?

Example:
```
**Feature:** Support for parallel check execution with failure tolerance

**Use case:** Large projects with many checks want faster feedback while still reporting all failures

**Proposed solution:** Add `fail-mode` option (fail-fast vs fail-all) to control whether to stop on first failure
```

## Project Philosophy

- **Minimal Overhead** — Keep VibeGuard lightweight and fast
- **Composable** — Support flexible policy definitions
- **Simple by Default** — Basic usage should be straightforward
- **Actionable Output** — Errors should clearly communicate what failed and how to fix it
- **Zero Dependencies** — Binary should have no runtime dependencies

## Documentation

- Keep README.md current with new features
- Document CLI changes in `README.md#cli-reference`
- Add examples to `examples/` for new patterns
- Document configuration changes in `README.md#configuration-schema`
- Update CONVENTIONS.md if changing code style guidelines

## AI-Assisted Setup Development

The `--assist` feature uses the inspector package to analyze projects and generate AI-friendly setup guides.

### Inspector Package Structure

```
internal/cli/inspector/
├── detector.go         # Project type detection
├── tools.go            # Tool scanning (linters, formatters, etc.)
├── metadata.go         # Metadata and structure extraction
├── recommendations.go  # Check recommendation generation
└── *_test.go          # Comprehensive tests for each component
```

### Adding Support for New Tools

To add detection for a new tool:

1. **Add to tools.go:**
   ```go
   // In the appropriate scan function (e.g., scanGoTools)
   if exists(".newtool.yml") || exists(".newtool.yaml") {
       tools = append(tools, ToolInfo{
           Name:       "newtool",
           Category:   CategoryLinter,
           Detected:   true,
           ConfigFile: ".newtool.yml",
           Confidence: 0.9,
           Indicators: []string{".newtool.yml config found"},
       })
   }
   ```

2. **Add check template in recommendations.go:**
   ```go
   var newtoolCheck = CheckRecommendation{
       ID:          "newtool",
       Description: "Run newtool analysis",
       Rationale:   "Newtool catches common issues...",
       Command:     "newtool check ./...",
       Severity:    "error",
       Suggestion:  "Fix newtool issues",
       Category:    "lint",
       Tool:        "newtool",
       Priority:    20,
   }
   ```

3. **Add tests in tools_test.go and recommendations_test.go**

### Adding Support for New Project Types

To add detection for a new language/framework:

1. **Add type constant in detector.go:**
   ```go
   const (
       ProjectTypeNewLang ProjectType = "newlang"
   )
   ```

2. **Add detection logic in detector.go:**
   ```go
   func (d *Detector) detectNewLang() DetectionResult {
       confidence := 0.0
       indicators := []string{}

       if d.exists("newlang.config") {
           confidence += 0.6
           indicators = append(indicators, "newlang.config found")
       }
       // ... more detection logic

       return DetectionResult{
           Type:       ProjectTypeNewLang,
           Confidence: min(confidence, 1.0),
           Indicators: indicators,
       }
   }
   ```

3. **Add tool scanning in tools.go**

4. **Add check templates in recommendations.go**

5. **Add comprehensive tests**

### Testing AI-Assisted Setup

Run the inspector tests:

```bash
# Unit tests
go test -v ./internal/cli/inspector/...

# Integration tests with real projects
go test -v ./internal/cli/inspector/... -run Integration

# Test prompt generation
go test -v ./internal/cli/inspector/... -run TestGenerateSetupPrompt
```

### Example: Testing with a Real Project

```bash
# Clone a test project
git clone https://github.com/example/project /tmp/test-project

# Run inspector analysis
cd /tmp/test-project
vibeguard init --assist --verbose

# Review the generated guide
```

## Performance Considerations

- Minimize allocations in hot paths
- Use goroutines judiciously (already parallelizing checks)
- Profile before optimizing: `go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...`
- Consider impact on single-check execution time (frequent in pre-commit hooks)

## Questions?

- Review existing issues and PRs to see if your question has been answered
- Check the documentation in `docs/` and `README.md`
- Open a GitHub discussion for questions
- Reference relevant ADRs for architectural context

## Recognition

Contributors will be recognized in:
- Release notes for significant contributions
- CONTRIBUTORS.md file (to be created)
- GitHub's contributor list

Thank you for contributing to VibeGuard!
