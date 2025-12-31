---
summary: Integrated go-grok library for pattern matching in check output extraction
event_type: code
sources:
  - https://github.com/elastic/go-grok
  - internal/grok/grok.go
  - internal/orchestrator/orchestrator.go
tags:
  - grok
  - pattern-matching
  - phase-2
  - output-extraction
  - vibeguard-c9m
---

# Grok Pattern Matching Integration

Completed task **vibeguard-c9m.1**: Integrate go-grok library for pattern matching in check output.

## Implementation Summary

### Dependencies Added
- `github.com/elastic/go-grok v0.3.1` - Elastic's grok pattern matching library
- `github.com/magefile/mage v1.15.0` - Transitive dependency

### New Package: `internal/grok`

Implemented `grok.Matcher` with the following features:

1. **Pattern Compilation** - Patterns are compiled eagerly at matcher creation to catch errors early
2. **Single and Multiple Patterns** - Supports both single pattern and array of patterns via `GrokSpec` type
3. **Value Extraction** - Returns `map[string]string` with named captures from patterns
4. **Pattern Merging** - When multiple patterns are provided, all are applied and results merged (later patterns can override)
5. **No-Match Handling** - Unmatched patterns result in missing keys (not empty strings)

### Orchestrator Integration

Updated both `Run()` and `RunCheck()` methods in `internal/orchestrator/orchestrator.go`:

```go
// Apply grok patterns to extract values from output
extracted := make(map[string]string)
if len(check.Grok) > 0 {
    matcher, matcherErr := grok.New(check.Grok)
    if matcherErr != nil {
        return matcherErr
    }
    extracted, matcherErr = matcher.Match(execResult.Combined)
    if matcherErr != nil {
        return matcherErr
    }
}
```

Extracted values are:
- Stored in `CheckResult.Extracted`
- Propagated to `Violation.Extracted` for failed checks
- Available for interpolation in suggestions via `{{.key}}` syntax

## Test Coverage

### Grok Package Tests (19 tests)
- Empty/nil patterns handling
- Valid/invalid pattern compilation
- Single and multiple pattern matching
- IP address, port, UUID, coverage percentage extraction
- Multiline input handling
- Pattern override behavior
- Built-in grok patterns (%{IP}, %{NUMBER}, %{UUID}, etc.)

### Orchestrator Integration Tests (8 new tests)
- `TestRun_GrokExtractsValues` - Basic extraction
- `TestRun_GrokMultiplePatterns` - Multiple pattern merging
- `TestRun_GrokNoPatterns_EmptyExtracted` - No patterns case
- `TestRun_GrokNoMatch_EmptyValues` - Pattern doesn't match
- `TestRun_GrokInvalidPattern_ReturnsError` - Error handling
- `TestRunCheck_GrokExtractsValues` - Single check extraction
- `TestRun_GrokExtractedInViolation` - Values in violations

## Usage Example

```yaml
checks:
  - id: test-coverage
    run: go test -cover ./...
    grok:
      - "coverage: (?P<coverage>[0-9.]+)%"
    assert: "float(coverage) >= 80.0"  # Phase 2
    suggestion: "Coverage is {{.coverage}}%, needs 80%"
```

## Key Decisions

1. **Eager Compilation** - Patterns compile at matcher creation rather than lazily to fail fast on invalid patterns
2. **Combined Output** - Grok patterns match against `execResult.Combined` (stdout + stderr interleaved)
3. **Error Propagation** - Invalid patterns return errors that stop orchestration (not silent failures)
4. **Named Captures** - Using Go regex named captures `(?P<name>...)` alongside built-in grok patterns

## Next Steps

Remaining Phase 2 tasks:
- **vibeguard-c9m.2**: Assertion expression parser and evaluator
- **vibeguard-c9m.3**: Templated suggestions
- **vibeguard-c9m.4**: JSON output mode
