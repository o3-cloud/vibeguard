# Pre-Commit Integration Examples

This directory contains example `.pre-commit-config.yaml` configurations for integrating VibeGuard with the [pre-commit](https://pre-commit.com) framework.

## Quick Start

1. Install pre-commit:
   ```bash
   pip install pre-commit
   # or
   brew install pre-commit
   ```

2. Copy an example to your project root:
   ```bash
   cp examples/pre-commit/go-project.pre-commit-config.yaml .pre-commit-config.yaml
   ```

3. Install the hooks:
   ```bash
   pre-commit install
   ```

4. Run manually (optional):
   ```bash
   pre-commit run --all-files
   ```

## Examples

### basic.pre-commit-config.yaml

Minimal configuration with just VibeGuard. Use this when VibeGuard handles all your checks.

### go-project.pre-commit-config.yaml

Go project configuration combining VibeGuard with Go-specific pre-commit hooks.

### node-project.pre-commit-config.yaml

Node.js/TypeScript project configuration combining VibeGuard with JavaScript ecosystem hooks.

### multi-language.pre-commit-config.yaml

Polyglot project configuration showing how to combine VibeGuard with multiple language-specific hooks.

## Using VibeGuard from GitHub

Once VibeGuard is published, you can reference it directly:

```yaml
repos:
  - repo: https://github.com/vibeguard/vibeguard
    rev: v1.0.0
    hooks:
      - id: vibeguard
```

## Using VibeGuard Locally

For local development or when VibeGuard isn't published yet:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard
        name: VibeGuard Policy Checks
        entry: vibeguard check
        language: system
        pass_filenames: false
        always_run: true
```

## Available Hook IDs

VibeGuard provides multiple hook variants (defined in `.pre-commit-hooks.yaml`):

| Hook ID | Description |
|---------|-------------|
| `vibeguard` | Run all VibeGuard checks |
| `vibeguard-fail-fast` | Run checks, stop on first failure |
| `vibeguard-verbose` | Run checks with verbose output |
| `vibeguard-validate` | Validate config without running checks |

## Configuration Tips

### Running Specific Checks

Override the entry to run specific checks:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard-quick
        name: Quick VibeGuard Checks
        entry: vibeguard check fmt vet
        language: system
        pass_filenames: false
```

### Fail-Fast Mode

Use fail-fast for faster feedback:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard
        name: VibeGuard
        entry: vibeguard check --fail-fast
        language: system
        pass_filenames: false
```

### Verbose Output

Get detailed output during commits:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard
        name: VibeGuard
        entry: vibeguard check --verbose
        language: system
        pass_filenames: false
        verbose: true
```

### Conditional Execution

Only run when certain files change:

```yaml
repos:
  - repo: local
    hooks:
      - id: vibeguard-go
        name: VibeGuard Go Checks
        entry: vibeguard check fmt vet
        language: system
        pass_filenames: false
        files: '\.go$'
```

## Troubleshooting

### "vibeguard: command not found"

Ensure VibeGuard is installed and in your PATH:

```bash
# Install globally
go install github.com/vibeguard/vibeguard/cmd/vibeguard@latest

# Or add to PATH if built locally
export PATH=$PATH:/path/to/vibeguard
```

### Hook is slow

1. Use `--fail-fast` to stop on first failure
2. Run only necessary checks: `vibeguard check fmt vet`
3. Increase parallelism: `vibeguard check --parallel 8`

### Skipping hooks temporarily

```bash
git commit --no-verify -m "WIP: work in progress"
```

## See Also

- [Pre-commit Documentation](https://pre-commit.com)
- [VibeGuard CI/CD Integrations](../../docs/INTEGRATIONS.md)
- [VibeGuard Configuration Schema](../../README.md#configuration-schema)
