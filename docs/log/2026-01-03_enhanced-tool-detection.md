---
summary: Implemented tool detection from Makefile, CI configs, and scripts directories for improved coverage when config files are absent
event_type: code
sources:
  - internal/cli/inspector/tools.go
  - internal/cli/inspector/tools_test.go
  - docs/adr/ADR-007-adopt-mutation-testing.md
tags:
  - ai-assisted-setup
  - tool-detection
  - inspector
  - makefile
  - ci-cd
  - github-actions
  - gitlab-ci
  - enhancement
---

# Enhanced Tool Detection from Makefile and CI Configs

Completed task vibeguard-ytq to detect tools from Makefile and CI configs when traditional config files are absent.

## Problem Statement

Previously, tools like golangci-lint were only detected if they had a config file (e.g., `.golangci.yml`). Many projects use tools with default configs and only reference them in Makefiles or CI workflows, leading to missed detection opportunities.

## Solution

Added four new helper methods to `ToolScanner`:

1. **`scanMakefileForTool(toolName string)`** - Scans `Makefile`, `makefile`, and `GNUmakefile` for tool references
2. **`scanCIWorkflowsForTool(toolName string)`** - Scans CI workflow files:
   - GitHub Actions (`.github/workflows/*.yml`, `*.yaml`)
   - GitLab CI (`.gitlab-ci.yml`)
   - CircleCI (`.circleci/config.yml`, `config.yaml`)
   - Travis CI (`.travis.yml`)
   - Jenkins (`Jenkinsfile`, `jenkins/Jenkinsfile`)
3. **`scanScriptsForTool(toolName string)`** - Scans scripts directories (`scripts/`, `script/`, `bin/`, `tools/`) for shell scripts containing tool references
4. **`enhanceToolDetection(toolName string)`** - Combines all three methods and returns confidence + indicators

## Tools Enhanced

Applied enhanced detection to all major tools:

**Go Tools:**
- golangci-lint
- goimports

**Node.js Tools:**
- eslint
- prettier
- jest
- mocha
- vitest
- typescript (via `tsc` command)

**Python Tools:**
- black
- pylint
- pytest
- mypy
- ruff
- flake8
- isort

## Confidence Levels

- **0.75** - CI workflow reference (strong signal)
- **0.70** - Makefile reference (good signal)
- **0.65** - Scripts directory reference (moderate signal)

These are lower than config file detection (0.9-0.95) but still provide useful detection when traditional indicators are absent.

## Test Coverage

Added 20+ new tests covering:
- Makefile detection (standard, lowercase, GNUmakefile)
- CI workflow detection (GitHub Actions, GitLab CI, CircleCI, Travis CI, Jenkins)
- Scripts directory detection
- Multiple indicators from different sources
- Tool-specific detection across Go, Node.js, and Python ecosystems

## Example Detection

A Makefile with:
```makefile
lint:
    golangci-lint run ./...
```

Now triggers golangci-lint detection with confidence 0.7 and indicator `golangci-lint in Makefile`, even without `.golangci.yml`.

## Files Changed

- `internal/cli/inspector/tools.go` - Added 4 new methods and integrated with existing tool scanners
- `internal/cli/inspector/tools_test.go` - Added comprehensive test coverage

## Next Steps

- Consider adding detection from `package.json` scripts field for Node.js tools
- Could expand to detect tools from Docker/Containerfile
- Monitor for false positives in real-world usage
