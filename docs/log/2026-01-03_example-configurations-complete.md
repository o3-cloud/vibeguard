---
summary: Completed example configurations task (vibeguard-o6g.2) by adding fix fields to existing examples and creating advanced.yaml with LLM-as-judge patterns.
event_type: code
sources:
  - docs/log/2025-12-30_llm-as-judge-and-templated-suggestions.md
  - examples/advanced.yaml
  - examples/go-project.yaml
  - examples/node-project.yaml
  - examples/simple.yaml
tags:
  - vibeguard
  - examples
  - llm-as-judge
  - documentation
  - phase-4
---

# Example Configurations Complete

Completed task vibeguard-o6g.2: Example configurations.

## Work Performed

### 1. Added `fix` Field to Existing Examples

All three existing examples were missing the recently-added `fix` field. Updated:

- **simple.yaml** - Added fix commands for fmt, vet, and test checks
- **go-project.yaml** - Added fix commands for all 6 checks, including coverage with HTML report generation
- **node-project.yaml** - Added fix commands for all 6 checks (Prettier, ESLint, TypeScript, tests, coverage, build)

### 2. Created advanced.yaml

New example demonstrating advanced VibeGuard features:

- **LLM-as-judge checks** using Claude CLI for architecture review, security audit, and PR quality assessment
- **Complex grok patterns** with multiple extractions (coverage, complexity, LLM verdicts)
- **Dependency chains** organizing checks into phases (fast deterministic -> testing -> LLM-powered)
- **Templated suggestions** using extracted values like `{{.coverage}}`, `{{.reason}}`, `{{.score}}`
- **Cyclomatic complexity check** using gocyclo with configurable threshold
- **Commented alternatives** for Ollama (local) and Gemini CLI

### 3. Validation

All 4 example files pass `vibeguard validate`:
- simple.yaml: 3 checks
- go-project.yaml: 6 checks
- node-project.yaml: 6 checks
- advanced.yaml: 10 checks

## Key Patterns Demonstrated

| Pattern | Example Location |
|---------|------------------|
| Basic checks | simple.yaml |
| Variable interpolation | go-project.yaml (go_packages, min_coverage) |
| Grok extraction | go-project.yaml (coverage), advanced.yaml (complexity, LLM verdicts) |
| Assertions | go-project.yaml (coverage >= min), advanced.yaml (verdict == "PASS") |
| Dependencies | All files use `requires` for execution ordering |
| LLM integration | advanced.yaml (Claude, with Ollama/Gemini alternatives) |
| Fix field | All files now include actionable fix commands |

## Task Closure

Closed vibeguard-o6g.2 in beads tracking system.
