# Sample Setup Prompt: Python Project with Black, pytest, and mypy

This is an example of the setup guide generated for a Python project.

---

# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.

---

## Project Analysis

**Project Name:** example-package
**Project Type:** python
**Detection Confidence:** 100%
**Language Version:** 3.11

### Main Tools Detected:
- Black (config: pyproject.toml [tool.black])
- isort (config: pyproject.toml [tool.isort])
- Ruff (config: pyproject.toml [tool.ruff])
- mypy (config: pyproject.toml [tool.mypy])
- pytest (config: pyproject.toml [tool.pytest])
- pip-audit (security scanner)

### Project Structure:
- Source directories: src/example_package
- Test directories: tests
- Entry points: src/example_package/__main__.py

### Build System:
- Build tool: setuptools (pyproject.toml)

---

## Recommended Checks

Based on the detected tools, here are the recommended checks:

### format (format)
**Description:** Check Python code formatting with Black
**Rationale:** Consistent formatting improves readability and reduces diffs
**Command:** `black --check .`
**Severity:** error
**Suggestion on failure:** Run 'black .' to format your Python code.

### isort (format)
**Description:** Check import sorting with isort
**Rationale:** Consistent import ordering improves readability
**Command:** `isort --check-only .`
**Severity:** error
**Suggestion on failure:** Run 'isort .' to sort imports.
**Requires:** format

### lint (lint)
**Description:** Run Ruff linter for code quality
**Rationale:** Ruff is a fast Python linter that catches common issues
**Command:** `ruff check .`
**Severity:** error
**Suggestion on failure:** Fix Ruff errors. Run 'ruff check . --fix' for auto-fixes.
**Requires:** isort

### typecheck (typecheck)
**Description:** Run mypy for static type checking
**Rationale:** Type checking catches bugs before runtime
**Command:** `mypy src/`
**Severity:** error
**Suggestion on failure:** Fix type errors before committing.

### test (test)
**Description:** Run pytest tests
**Rationale:** Tests verify that code behaves as expected
**Command:** `pytest`
**Severity:** error
**Suggestion on failure:** Fix failing tests before committing.
**Requires:** typecheck

### coverage (test)
**Description:** Check test coverage meets minimum threshold
**Rationale:** Code coverage helps identify untested code paths
**Command:** `pytest --cov=src --cov-report=term-missing`
**Severity:** warning
**Grok Patterns:** `TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%`
**Assertion:** `coverage >= 80`
**Suggestion on failure:** Coverage is {{.coverage}}%, target is 80%. Add tests.
**Requires:** test

### audit (security)
**Description:** Check for known vulnerabilities in dependencies
**Rationale:** Security vulnerabilities in dependencies can compromise your application
**Command:** `pip-audit`
**Severity:** warning
**Suggestion on failure:** Update vulnerable dependencies. Run 'pip-audit --fix' for auto-fixes.

---

## Configuration Requirements

A valid vibeguard.yaml must contain:

### 1. Version (required)
```yaml
version: "1"
```

### 2. Variables (optional)
Global variables for interpolation in commands and assertions.

```yaml
vars:
  src_dir: "src"
  coverage_threshold: "80"
```

### 3. Checks (required)
Array of check definitions. Each check must have:
- **id:** Unique identifier (alphanumeric + underscore + hyphen)
- **run:** Shell command to execute

Optional fields:
- **grok:** Array of patterns to extract data from output
- **assert:** Condition that must be true
- **requires:** Array of check IDs that must pass first
- **severity:** "error" or "warning" (default: error)
- **suggestion:** Message shown on failure
- **timeout:** Duration string (e.g., "30s", "5m")

---

## Python-Specific Examples

### Format Check with Black
```yaml
- id: format
  run: black --check .
  severity: error
  suggestion: "Run 'black .' to format code"
  timeout: 30s
```

### Import Sorting with isort
```yaml
- id: isort
  run: isort --check-only .
  severity: error
  suggestion: "Run 'isort .' to sort imports"
  timeout: 30s
  requires:
    - format
```

### Lint Check with Ruff
```yaml
- id: lint
  run: ruff check .
  severity: error
  suggestion: "Fix Ruff errors. Run 'ruff check . --fix' for auto-fixes."
  timeout: 60s
  requires:
    - isort
```

### Type Checking with mypy
```yaml
- id: typecheck
  run: mypy {{.src_dir}}/
  severity: error
  suggestion: "Fix mypy type errors before committing"
  timeout: 120s
```

### pytest with Coverage
```yaml
- id: test
  run: pytest
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - typecheck

- id: coverage
  run: pytest --cov={{.src_dir}} --cov-report=term-missing
  grok:
    - TOTAL\s+\d+\s+\d+\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%."
  requires:
    - test
  timeout: 300s
```

### Security Audit with pip-audit
```yaml
- id: audit
  run: pip-audit
  severity: warning
  suggestion: "Update vulnerable packages. Run 'pip-audit --fix' for auto-fixes."
  timeout: 120s
```

---

## Validation Rules

Your generated configuration must:

1. Be valid YAML syntax
2. Have `version: "1"` at the top level
3. Include at least one check in the `checks` array
4. Each check must have a unique `id`
5. Each check must have a non-empty `run` command
6. All `requires` references must point to existing check IDs
7. No circular dependencies in `requires`
8. All variables used must be defined in `vars`
9. `severity` must be "error" or "warning"
10. Grok patterns must be valid Grok syntax

**DO NOT:**
- Include YAML comments in the generated config
- Add extra top-level keys beyond version, vars, checks
- Use undefined variables
- Create checks for tools not detected in the project

---

## Your Task

Based on the project analysis above, generate a vibeguard.yaml configuration that:

1. Includes version: "1"
2. Defines appropriate variables for this project
3. Creates checks for the detected tools
4. Follows the syntax rules described above
5. Includes helpful suggestions for each check
6. Uses appropriate timeouts for each check type

Output the configuration in a YAML code block.
