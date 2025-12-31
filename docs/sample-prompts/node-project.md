# Sample Setup Prompt: Node.js Project with ESLint, Prettier, and Jest

This is an example of the setup guide generated for a Node.js/TypeScript project.

---

# VibeGuard AI Agent Setup Guide

You are being asked to help set up VibeGuard policy enforcement for a software project.
VibeGuard is a declarative policy tool that runs quality checks and assertions on code.

This guide will help you understand the project structure, existing tools, and how to
generate a valid configuration.

---

## Project Analysis

**Project Name:** @example/web-app
**Project Type:** node
**Detection Confidence:** 100%
**Language Version:** 20.x

### Main Tools Detected:
- ESLint (config: eslint.config.js)
- Prettier (config: .prettierrc.json)
- Jest (config: jest.config.js)
- TypeScript (config: tsconfig.json)
- npm audit (built-in security scanner)

### Project Structure:
- Source directories: src
- Test directories: __tests__, src/**/*.test.ts
- Entry points: src/index.ts

### Build System:
- Build output: dist/
- Package manager: npm

---

## Recommended Checks

Based on the detected tools, here are the recommended checks:

### format (format)
**Description:** Check code formatting with Prettier
**Rationale:** Consistent formatting improves readability and reduces merge conflicts
**Command:** `npx prettier --check .`
**Severity:** error
**Suggestion on failure:** Run 'npx prettier --write .' to format your code.

### lint (lint)
**Description:** Run ESLint to check for code quality issues
**Rationale:** ESLint catches bugs and enforces best practices
**Command:** `npx eslint .`
**Severity:** error
**Suggestion on failure:** Fix ESLint errors. Run 'npx eslint . --fix' for auto-fixes.
**Requires:** format

### typecheck (typecheck)
**Description:** Run TypeScript compiler to check for type errors
**Rationale:** Type checking catches bugs before runtime
**Command:** `npx tsc --noEmit`
**Severity:** error
**Suggestion on failure:** Fix TypeScript errors before committing.

### test (test)
**Description:** Run Jest tests
**Rationale:** Tests verify that code behaves as expected
**Command:** `npm test`
**Severity:** error
**Suggestion on failure:** Fix failing tests before committing.
**Requires:** typecheck

### coverage (test)
**Description:** Check test coverage meets minimum threshold
**Rationale:** Code coverage helps identify untested code paths
**Command:** `npm test -- --coverage --coverageReporters=text-summary`
**Severity:** warning
**Grok Patterns:** `Statements\s+:\s+%{NUMBER:coverage}%`
**Assertion:** `coverage >= 80`
**Suggestion on failure:** Statement coverage is {{.coverage}}%, target is 80%. Add tests.
**Requires:** test

### audit (security)
**Description:** Check for known vulnerabilities in dependencies
**Rationale:** Security vulnerabilities in dependencies can compromise your application
**Command:** `npm audit --audit-level=high`
**Severity:** warning
**Suggestion on failure:** Run 'npm audit fix' to fix vulnerabilities, or review and update dependencies manually.

### build (build)
**Description:** Verify TypeScript compiles successfully
**Rationale:** Catch compilation errors before deployment
**Command:** `npm run build`
**Severity:** error
**Suggestion on failure:** Fix build errors before committing.
**Requires:** lint, typecheck

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
  coverage_threshold: "80"
  audit_level: "high"
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

## Node.js-Specific Examples

### Format Check with Prettier
```yaml
- id: format
  run: npx prettier --check .
  severity: error
  suggestion: "Run 'npx prettier --write .' to format code"
  timeout: 30s
```

### Lint Check with ESLint
```yaml
- id: lint
  run: npx eslint .
  severity: error
  suggestion: "Fix ESLint errors. Run 'npx eslint . --fix' for auto-fixes."
  timeout: 60s
  requires:
    - format
```

### TypeScript Type Checking
```yaml
- id: typecheck
  run: npx tsc --noEmit
  severity: error
  suggestion: "Fix TypeScript type errors before committing"
  timeout: 60s
```

### Jest Tests with Coverage
```yaml
- id: test
  run: npm test
  severity: error
  suggestion: "Fix failing tests before committing"
  timeout: 300s
  requires:
    - typecheck

- id: coverage
  run: npm test -- --coverage --coverageReporters=text-summary
  grok:
    - Statements\s+:\s+%{NUMBER:coverage}%
  assert: "coverage >= {{.coverage_threshold}}"
  severity: warning
  suggestion: "Coverage is {{.coverage}}%, target is {{.coverage_threshold}}%."
  requires:
    - test
  timeout: 300s
```

### Security Audit
```yaml
- id: audit
  run: npm audit --audit-level={{.audit_level}}
  severity: warning
  suggestion: "Run 'npm audit fix' or update vulnerable dependencies"
  timeout: 60s
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
