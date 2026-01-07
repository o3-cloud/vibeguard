# VibeGuard Tags - Standard Conventions and Use Cases

**Last Updated:** 2026-01-07

This document describes the standard tag conventions for VibeGuard checks. Tags enable flexible categorization and filtering of checks without modifying configuration files.

## Overview

Tags are optional labels you can assign to checks for categorization and filtering. Each check can have multiple tags, allowing you to run subsets of checks by category using CLI flags.

**Key Benefits:**
- Run fast checks locally during development (`--tags fast`)
- Configure pre-commit hooks to run only suitable checks (`--tags pre-commit`)
- Run comprehensive CI validation with both fast and slow checks
- Skip expensive LLM checks during local iterations (`--exclude-tags llm`)

## Standard Tag Conventions

The following tags are recommended (not enforced). You can create custom tags following the naming pattern: lowercase alphanumeric with optional hyphens (`^[a-z][a-z0-9-]*$`).

### Execution Performance

#### `fast`
**Description:** Quick checks that complete in less than 5 seconds
**Use Cases:**
- Pre-commit hooks (run before every commit)
- Local development iteration (fast feedback loop)
- CI fast-path validation
- Initial verification before pushing

**Examples:**
- `gofmt -l .` (formatting check)
- `go vet ./...` (static analysis)
- `eslint .` (JavaScript linting)
- `black --check .` (Python formatting)

**Configuration:**
```yaml
- id: fmt
  run: gofmt -l .
  tags: [format, fast, pre-commit]
```

#### `slow`
**Description:** Long-running checks that take more than 30 seconds
**Use Cases:**
- CI/CD pipelines (run comprehensive validation)
- End-of-shift validation (when user can wait)
- Scheduled nightly builds
- Post-merge validation

**Examples:**
- `go test ./...` (full test suite)
- `gosec ./...` (security scanning)
- Integration test suites
- Performance profiling

**Configuration:**
```yaml
- id: test-integration
  run: go test -tags=integration ./...
  tags: [test, slow, ci]
```

### Check Categories

#### `format`
**Description:** Code formatting and style enforcement
**Use Cases:**
- Ensuring consistent code style
- Automated formatting checks
- Pre-commit hooks
- CI gating

**Related Tags:** Usually combined with `fast` or `pre-commit`

**Examples:**
- `gofmt -l .`
- `black --check .`
- `prettier --check .`

#### `lint`
**Description:** Static analysis and linting
**Use Cases:**
- Code quality checks
- Detecting common mistakes
- Enforcing best practices
- Pre-commit hooks

**Related Tags:** Usually combined with `fast`

**Examples:**
- `golangci-lint run`
- `eslint .`
- `pylint src/`

#### `test`
**Description:** Unit, integration, and functional tests
**Use Cases:**
- Validating functionality
- Regression prevention
- CI/CD validation
- Coverage tracking

**Related Tags:** Can be both `fast` and `slow` depending on test scope

**Examples:**
- `go test ./...` (fast unit tests)
- `go test -tags=integration ./...` (slow integration tests)
- `npm test` (JavaScript tests)
- `pytest` (Python tests)

#### `build`
**Description:** Compilation, packaging, and artifact creation
**Use Cases:**
- Ensuring code compiles
- Building Docker images
- Creating distribution packages
- Verifying deployment artifacts

**Related Tags:** Usually `slow`, often `ci`

**Examples:**
- `go build ./...`
- `docker build .`
- `npm run build`
- Maven build

#### `security`
**Description:** Security scanning and vulnerability detection
**Use Cases:**
- Detecting known vulnerabilities
- Checking for security issues
- Compliance validation
- CI gating

**Related Tags:** Usually `slow` and `ci`

**Examples:**
- `gosec ./...`
- `npm audit`
- `safety check` (Python)
- `bandit` (Python security)

### Execution Context

#### `pre-commit`
**Description:** Suitable for pre-commit git hooks
**Use Cases:**
- Automatic validation before committing
- Local development workflow
- Fast feedback loop
- Prevents committing broken code

**Characteristics:**
- Should be fast (< 5 seconds)
- Should be deterministic
- Should focus on local issues
- Should not require external services

**Configuration Example:**
```yaml
- id: fmt
  run: gofmt -l .
  tags: [format, fast, pre-commit]

- id: lint
  run: golangci-lint run
  tags: [lint, fast, pre-commit]
```

**Git Hook Integration:**
```bash
#!/bin/bash
# .git/hooks/pre-commit
vibeguard check --tags pre-commit --fail-fast
```

#### `ci`
**Description:** Checks suitable for CI/CD pipelines
**Use Cases:**
- Comprehensive validation in CI/CD
- Running all checks (fast and slow)
- Production readiness verification
- Automated pull request validation

**Characteristics:**
- Can include both fast and slow checks
- Often includes security scanning
- May use external services
- Runs on every push or pull request

**Configuration Example:**
```yaml
checks:
  - id: fmt
    run: gofmt -l .
    tags: [format, fast, ci]

  - id: test
    run: go test ./...
    tags: [test, slow, ci]

  - id: security
    run: gosec ./...
    tags: [security, slow, ci]
```

**CI Pipeline Integration:**
```bash
# In GitHub Actions or similar
- name: Run VibeGuard checks
  run: vibeguard check --tags ci
```

#### `ci-only`
**Description:** Checks that should only run in CI/CD (not locally)
**Use Cases:**
- Expensive external API calls
- LLM-powered analysis
- Resource-intensive operations
- Checks requiring specific environment setup

**Characteristics:**
- Not suitable for pre-commit hooks
- Usually slow
- May require credentials or configuration
- Deterministic but expensive

**Configuration Example:**
```yaml
- id: llm-security-review
  run: ./scripts/llm-security-review.sh
  tags: [llm, security, slow, ci-only]
```

**Local Development (skip these):**
```bash
# Run fast checks, skip ci-only
vibeguard check --exclude-tags ci-only,slow
```

### Special Capabilities

#### `llm`
**Description:** Checks powered by Large Language Models
**Use Cases:**
- Architecture review
- Security analysis
- Code quality assessment
- PR quality evaluation

**Characteristics:**
- Usually slow (requires API calls)
- Usually ci-only
- May have rate limiting
- More subjective results

**Configuration Example:**
```yaml
- id: architecture-review
  run: llm-check analyze-architecture --file /dev/stdin
  tags: [llm, slow, ci-only]
```

**Local Development Workaround:**
```bash
# Run everything except LLM checks for faster iteration
vibeguard check --exclude-tags llm
```

## Usage Patterns

### Local Development

Fast iteration with only essential checks:
```bash
# Run only fast pre-commit checks
vibeguard check --tags fast

# Or be more specific
vibeguard check --tags format,lint --exclude-tags slow
```

### Pre-Commit Hooks

Automatic validation before committing:
```bash
# Install as git hook
vibeguard check --tags pre-commit --fail-fast
```

### CI/CD Pipeline

Comprehensive validation:
```bash
# Run all CI checks (includes slow tests and security scanning)
vibeguard check --tags ci

# Or run everything except ci-only (e.g., skip LLM reviews)
vibeguard check --exclude-tags ci-only
```

### Tag Discovery

Find all available tags in your configuration:
```bash
vibeguard tags
```

Output:
```
build
ci
fast
format
lint
llm
pre-commit
security
slow
test
```

## Best Practices

### 1. Tag Assignment

**DO:**
- Assign multiple tags per check for maximum flexibility
- Include both category (`format`, `lint`, etc.) and timing (`fast`, `slow`)
- Include execution context (`pre-commit`, `ci`)

**Example:**
```yaml
- id: fmt
  run: gofmt -l .
  tags: [format, fast, pre-commit]  # Category, timing, context
```

**DON'T:**
- Use uppercase letters (not enforced, but inconsistent)
- Create tags you won't use for filtering
- Make tags too specific (e.g., `fmt-check-spaces-only`)

### 2. Realistic Timing

**`fast`** = < 5 seconds on typical hardware
**`slow`** = > 30 seconds on typical hardware

If a check takes 10-30 seconds, choose based on your typical use case:
- Local development cycle? → `fast`
- Only runs in CI? → `slow`
- Borderline case? → Include both or skip timing tag

### 3. Dependency Handling

When a check depends on another check:

```yaml
- id: test
  run: go test ./...
  requires: [build]
  tags: [test, slow, ci]

- id: build
  run: go build ./...
  requires: [fmt]
  tags: [build, slow, ci]
```

If you filter by tags and exclude a dependency:
```bash
# This skips 'test' because 'build' is excluded
vibeguard check --tags test --exclude-tags build
# Output: Skipping check 'test': required dependency 'build' not in filtered set
```

**Solution:** Tag dependencies with overlapping tags:
```yaml
- id: test
  run: go test ./...
  requires: [build]
  tags: [test, slow, ci]

- id: build
  run: go build ./...
  requires: [fmt, vet]
  tags: [build, slow, ci]  # Same tags as 'test'
```

Now both work together:
```bash
vibeguard check --tags ci  # Both test and build run
```

### 4. Custom Tags

Use custom tags for project-specific organization:

```yaml
- id: db-migration
  run: ./migrate.sh
  tags: [build, slow, ci, optional]  # 'optional' is custom

- id: generate-docs
  run: cargo doc
  tags: [build, slow, ci, documentation]  # 'documentation' is custom
```

Skip optional checks:
```bash
vibeguard check --exclude-tags optional
```

## Troubleshooting

### Q: I ran `--tags pre-commit` but a slow check ran anyway

**A:** The check probably didn't have the `pre-commit` tag. Use `vibeguard tags` to see all available tags and `vibeguard list -v` to see which tags each check has.

### Q: My filtered check was skipped with a "dependency not in filtered set" error

**A:** The check depends on another check that doesn't match your filter. Either:
1. Add the dependency's tag to your filter
2. Tag the dependency with the same tags
3. Remove the dependency if not needed

### Q: I want different checks for different workflows

**A:** Use multiple tags per check:
```yaml
- id: lint
  run: golangci-lint run
  tags: [lint, fast, pre-commit, ci]
```

Then filter by workflow:
```bash
vibeguard check --tags pre-commit  # Pre-commit hook
vibeguard check --tags ci          # Full CI
vibeguard check --tags fast        # Quick local check
```

### Q: Can I use AND logic for tags? (e.g., checks that are both `fast` AND `lint`)

**A:** No, VibeGuard uses OR logic. Instead, filter for what you want:
```bash
# Run all lint checks (fast or slow)
vibeguard check --tags lint

# Run all fast checks (any category)
vibeguard check --tags fast

# Run fast checks, but skip lint
vibeguard check --tags fast --exclude-tags lint
```

## Reference

- **Spec:** [SPEC-tags.md](specs/SPEC-tags.md)
- **README:** [Check Tags section](../README.md#check-tags)
- **CLI Reference:** [vibeguard check --tags](CLI-REFERENCE.md#tags)
