# VibeGuard Configuration Examples

This directory contains example VibeGuard configurations for different project types and use cases.

## Getting Started

Choose an example that matches your project type, then customize it for your needs:

```bash
# Copy an example to your project root
cp examples/go-project.yaml vibeguard.yaml

# Review and edit the configuration as needed
vim vibeguard.yaml

# Validate the configuration
vibeguard validate

# Run checks
vibeguard check
```

## Examples

### simple.yaml

**Use case:** Getting started with VibeGuard

A minimal configuration demonstrating the basic checks every project should have:
- Code formatting
- Static analysis (go vet)
- Unit tests

**Good for:** First-time VibeGuard users or simple projects with few dependencies.

### go-project.yaml

**Use case:** Comprehensive Go project quality gates

A complete quality gate for Go projects including:
- Code formatting (gofmt)
- Static analysis (go vet)
- Linting (golangci-lint)
- Unit tests with coverage validation
- Build verification

**Features demonstrated:**
- Variable interpolation (`{{.var_name}}`)
- Check dependencies (`requires`)
- Grok pattern extraction and assertions
- Coverage threshold validation

**Good for:** Production Go projects that need comprehensive quality checks.

### node-project.yaml

**Use case:** Comprehensive Node.js/TypeScript project quality gates

A complete quality gate for JavaScript/TypeScript projects including:
- Code formatting (Prettier)
- Linting (ESLint)
- Type checking (TypeScript)
- Unit tests with coverage validation
- Build verification

**Features demonstrated:**
- Integration with npm scripts
- JSON parsing with grok patterns
- Coverage percentage assertions
- Multi-stage check dependencies

**Good for:** Production JavaScript/TypeScript projects.

## Configuration Tips

### Adjusting Timeouts

If checks are timing out, increase the `timeout` value:

```yaml
checks:
  - id: build
    run: npm run build
    timeout: 120s  # Increase from 60s to 120s
```

### Variable Interpolation

Use variables to avoid duplication:

```yaml
vars:
  src_dir: "src"
  test_dir: "test"

checks:
  - id: lint
    run: eslint {{.src_dir}}
  - id: test
    run: jest {{.test_dir}}
```

### Conditional Checks

Skip optional checks by using `severity: warning`:

```yaml
checks:
  - id: lint
    severity: warning  # Warns but doesn't fail the pipeline
```

### Check Ordering

Control execution order with `requires`:

```yaml
checks:
  - id: test
    requires:
      - lint        # lint must pass before test runs
      - build       # build must also pass before test runs
```

## Running Examples

To test an example configuration:

```bash
# Validate the configuration
vibeguard validate -c examples/go-project.yaml

# Run checks from an example
vibeguard check -c examples/go-project.yaml
```

## Adding New Examples

If you create a new configuration that might be useful for others:

1. Save it to `examples/` with a descriptive name
2. Add a README section describing the use case and features
3. Include inline comments explaining key configuration options
4. Submit a pull request to share it with the community

### pre-commit/

**Use case:** Integration with the pre-commit framework

Example `.pre-commit-config.yaml` files showing how to integrate VibeGuard with pre-commit hooks:

- `basic.pre-commit-config.yaml` — Minimal configuration with just VibeGuard
- `go-project.pre-commit-config.yaml` — Go project with VibeGuard and Go-specific hooks
- `node-project.pre-commit-config.yaml` — Node.js/TypeScript project configuration
- `multi-language.pre-commit-config.yaml` — Polyglot project with multiple language hooks

See `pre-commit/README.md` for detailed setup instructions.

## See Also

- [Configuration Schema](../README.md#configuration-schema)
- [CLI Reference](../README.md#cli-reference)
- [CONVENTIONS.md](../CONVENTIONS.md)
- [Pre-commit Framework](https://pre-commit.com)
- [Ecosystem Comparison](../docs/ECOSYSTEM-COMPARISON.md)
