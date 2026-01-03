# VibeGuard Coding Conventions

This document outlines the coding standards and conventions for the VibeGuard project.

## Language Standards

VibeGuard is implemented in **Go 1.21+** to ensure compatibility with recent language features and security updates.

## Code Organization

### Directory Structure

```
vibeguard/
├── cmd/
│   └── vibeguard/          # Main application entry point
├── internal/
│   ├── assert/             # Assertion expression evaluation
│   ├── cli/                # Command-line interface and Cobra commands
│   ├── config/             # Configuration loading and parsing
│   ├── executor/           # Policy execution logic
│   ├── grok/               # Grok pattern matching and parsing
│   ├── orchestrator/       # Policy orchestration and coordination
│   ├── output/             # Output formatting and rendering
│   └── version/            # Version information
├── docs/                   # Documentation and ADRs
└── go.mod, go.sum          # Go module files
```

### Package Layout

- **`cmd/`** — Application entry points (may have multiple binaries in future)
- **`internal/`** — Private packages used only within VibeGuard; not importable by external packages

## Go Code Style

### Naming Conventions

- **Interfaces**: End with `-er` (e.g., `Runner`, `Evaluator`, `PolicyReader`)
- **Errors**: Use descriptive error variables prefixed with `Err` (e.g., `ErrPolicyNotFound`)
- **Constants**: Use PascalCase for unexported constants, UPPER_CASE for public constants where appropriate
- **Functions**: Use CamelCase; exported functions start with uppercase letter
- **Receiver names**: Use single or two-letter abbreviations (e.g., `p` for `*Policy`, `pr` for `*PolicyRunner`)

### Error Handling

1. **Always check errors explicitly** — Do not ignore error returns
2. **Wrap errors with context** — Use `fmt.Errorf("context: %w", err)` for error wrapping
3. **Custom error types** — Define custom error types for domain-specific errors
4. **Error messages** — Start with lowercase letter, no trailing period (e.g., `"policy file not found"`)

### Comments

- **Exported items** — All exported functions, types, and constants must have doc comments
- **Complex logic** — Explain the "why", not the "what" (the code shows what it does)
- **Line comments** — Use `//` for single-line comments within functions
- **Doc format** — Follow Go conventions: start with the name of the function/type

Example:
```go
// PolicyRunner evaluates a policy against a resource.
type PolicyRunner interface {
	Run(ctx context.Context, policy *Policy) error
}
```

### Testing

- **Unit tests** — Colocate with the code they test (e.g., `policy.go` has `policy_test.go`)
- **Table-driven tests** — Use for parameterized test cases
- **Test naming** — Use `TestFunctionName` or `TestFunctionName_Scenario` for clarity
- **Integration tests** — Place in `tests/` directory for broader integration scenarios

Example:
```go
func TestPolicyEval_ValidPolicy(t *testing.T) {
	tests := []struct {
		name    string
		policy  *Policy
		input   interface{}
		want    bool
		wantErr bool
	}{
		// test cases
	}
	// test logic
}
```

## Formatting and Tooling

### Automated Tools

- **Format** — Run `go fmt ./...` before committing
- **Vet** — Run `go vet ./...` to catch common mistakes
- **Lint** — Use `golangci-lint run ./...` for advanced linting (optional but recommended)
- **Test** — Run `go test -v ./...` to verify tests pass

### Line Length

- Aim for lines under 100 characters
- Up to 120 characters is acceptable for readability in complex expressions
- Long function signatures or import blocks may exceed this when necessary

## Dependencies

### Adding Dependencies

1. Use `go get github.com/owner/repo@version` to add new dependencies
2. Minimize external dependencies; prefer standard library where reasonable
3. Keep dependencies up-to-date with `go get -u` and regular maintenance
4. Document why a dependency is needed in code comments or ADRs if non-obvious

### Dependency Guidelines

- **Stability** — Prefer mature, well-maintained packages
- **Size** — Minimize dependency bloat; prefer smaller, focused libraries
- **License** — Ensure compatibility with VibeGuard's license

## Security Considerations

- **Input validation** — Always validate external inputs (CLI args, file contents, API responses)
- **Error messages** — Avoid leaking sensitive information in error output
- **Credentials** — Never hardcode credentials; use environment variables or secure vaults
- **Unsafe operations** — Document and review any use of `unsafe` package

## Documentation

### Code Comments

- Document complex algorithms or non-obvious design decisions
- Use comments to explain the intent, not repeat the code
- Keep comments synchronized with code changes

### Module Documentation

Create a `README.md` in each major package directory explaining its purpose and usage.

## Performance Considerations

- **Benchmarking** — Use `go test -bench` for performance-critical code
- **Profiling** — Use `pprof` to identify bottlenecks before optimizing
- **Goroutines** — Be cautious with unbounded goroutine creation; use worker pools for scalability

## Version Control and Commits

- Follow **Conventional Commits** specification (see ADR-002)
- Example: `feat(policy): add YAML policy parser` or `fix(judge): handle timeout errors`
- Commits should be atomic and focused on a single concern

## Continuous Integration

- All code must pass `go test`, `go vet`, and formatting checks before merging
- Maintain test coverage above 70% for critical packages
- Review Go module updates regularly for security patches
