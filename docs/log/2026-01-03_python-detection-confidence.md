---
summary: Increased pyproject.toml confidence weight to 0.6 for parity with other modern manifest files
event_type: code
sources:
  - internal/cli/inspector/detector.go
  - internal/cli/inspector/detector_test.go
tags:
  - python
  - detection
  - confidence
  - ai-assisted-setup
  - inspector
---

# Python Project Detection Confidence Fix

## Context

Task: `vibeguard-9pl` - Increase Python project detection confidence

Python projects with `pyproject.toml` were getting 0.5 base confidence, while Go projects with `go.mod` get 0.6 base confidence. This caused Python to be ranked lower than expected in mixed-language projects.

## Analysis

Examined the confidence weights across all supported project types:

| Language | Manifest File     | Weight |
|----------|-------------------|--------|
| Go       | go.mod            | 0.6    |
| Node     | package.json      | 0.6    |
| Ruby     | Gemfile           | 0.6    |
| Rust     | Cargo.toml        | 0.7    |
| Java     | pom.xml           | 0.6    |
| Python   | pyproject.toml    | 0.5 (was) |

The `pyproject.toml` file is the modern standard manifest for Python, defined by PEP 517/518/621. It serves the same purpose as `go.mod`, `package.json`, or `Cargo.toml` in their respective ecosystems and should have equal weight.

## Changes

1. **detector.go:185-191**: Changed `pyproject.toml` confidence weight from 0.5 to 0.6
   - Added comment explaining the rationale (PEP 517/518/621 modern standard)

2. **detector_test.go:237-238**: Updated test expectations for "Python with pyproject.toml" test case
   - Changed `minConfidence` from 0.7 to 0.8
   - Changed `maxConfidence` from 0.75 to 0.85
   - This accounts for pyproject.toml (0.6) + *.py files (0.2) = 0.8

## Verification

- All inspector tests pass
- All project tests pass
- No regressions in other language detection
