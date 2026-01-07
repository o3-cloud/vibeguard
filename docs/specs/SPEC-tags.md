# Specification: Tags for vibeguard.yml

**Status:** Ready
**Author:** Claude Code
**Created:** 2026-01-07
**Related:** [ADR-005](../adr/ADR-005-adopt-vibeguard.md)

## Overview

Add support for tags in vibeguard.yml to enable flexible categorization and filtering of checks. Tags allow users to run subsets of checks by category (e.g., `--tags security,lint`) without modifying configuration files.

## Problem Statement

Currently, vibeguard only supports two modes of check execution:
1. Run all checks: `vibeguard check`
2. Run a single check by ID: `vibeguard check fmt`

Users cannot easily:
- Run all "fast" checks during pre-commit
- Run only "security" checks in CI
- Skip "slow" LLM checks during local development
- Group related checks for selective execution

The `examples/advanced.yaml` file uses comment-based "phases" as a workaround, demonstrating clear user need for this feature.

## Goals

1. Allow checks to be tagged with arbitrary labels
2. Enable filtering checks by tag via CLI flags (`--tags` with OR logic)
3. Maintain backwards compatibility (tags are optional)
4. Support exclusion via `--exclude-tags`
5. Display tags in list output and JSON results
6. Provide `vibeguard tags` command to discover all tags in config

## Non-Goals

1. Predefined/enforced tag vocabulary (tags are freeform)
2. Tag-based execution ordering (use `requires` for ordering)
3. Tag inheritance or hierarchies
4. Tag aliases or grouping
5. AND logic for tag filtering (OR is sufficient)
6. Validation of unused or typo'd tags (silently ignore non-matching filters)

## Design

### YAML Schema

Add a `tags` field to the Check struct:

```yaml
checks:
  - id: fmt
    run: gofmt -l .
    tags: [format, fast, pre-commit]
    severity: error
    suggestion: "Code formatting issues found"

  - id: security-scan
    run: gosec ./...
    tags: [security, slow, ci-only]
    severity: error
    suggestion: "Security vulnerabilities detected"
```

Tags are:
- A list of strings
- Optional (empty list if not specified)
- Case-sensitive
- Must match pattern: `^[a-z][a-z0-9-]*$` (lowercase, alphanumeric, hyphens)

### CLI Interface

#### Running checks with tags

```bash
# Run checks with ANY of the specified tags (OR logic)
vibeguard check --tags format,lint

# Exclude checks with ANY of the specified tags
vibeguard check --exclude-tags slow,llm

# Combine inclusion and exclusion
vibeguard check --tags ci --exclude-tags slow

# Run specific check (existing behavior unchanged)
vibeguard check fmt
```

#### Listing checks and tags

```bash
# List all checks with their tags
vibeguard list -v

# Filter list by tags
vibeguard list --tags security

# List all unique tags in the config
vibeguard tags
```

### Tag Filtering Logic

1. **No flags**: Run all checks (current behavior)
2. **`--tags`**: Run checks matching ANY specified tag (OR)
3. **`--exclude-tags`**: Exclude checks matching ANY specified tag
4. **Combined**: Apply inclusion first, then exclusion
5. **Non-matching tags**: Silently ignored (no error for typos or unused tags)

Example with checks:
```yaml
checks:
  - id: fmt        # tags: [format, fast]
  - id: lint       # tags: [lint, fast]
  - id: test       # tags: [test]
  - id: security   # tags: [security, slow]
```

| Command | Checks Run |
|---------|-----------|
| `vibeguard check` | fmt, lint, test, security |
| `vibeguard check --tags fast` | fmt, lint |
| `vibeguard check --tags security,test` | test, security |
| `vibeguard check --exclude-tags slow` | fmt, lint, test |
| `vibeguard check --tags fast --exclude-tags lint` | fmt |
| `vibeguard check --tags nonexistent` | (none - no error) |

### Dependency Handling

When filtering by tags, dependency resolution follows strict filtering:

- Only run checks matching the tag filter
- If a filtered check depends on an excluded check, **skip the check** with a warning
- Dependencies are NOT automatically included
- Rationale: Explicit is better than implicit; users should tag dependencies appropriately

Example:
```yaml
checks:
  - id: fmt
    tags: [format]
  - id: test
    tags: [test]
    requires: [fmt]  # depends on fmt
```

Running `vibeguard check --tags test` will:
1. Select only `test` (matches tag)
2. Skip `test` with warning because `fmt` (its dependency) is not included
3. Output: `Skipping check 'test': required dependency 'fmt' not in filtered set`

### Standard Tag Conventions

Document recommended (not enforced) tag conventions:

| Tag | Description |
|-----|-------------|
| `format` | Code formatting checks |
| `lint` | Static analysis / linting |
| `test` | Unit/integration tests |
| `security` | Security scanning |
| `build` | Compilation / build checks |
| `fast` | Quick checks (<5s) |
| `slow` | Long-running checks (>30s) |
| `pre-commit` | Suitable for pre-commit hooks |
| `ci` | CI/CD pipeline checks |
| `llm` | LLM-powered checks |

## Implementation

### Phase 1: Schema & Validation

**Files:** `internal/config/schema.go`, `internal/config/config.go`

1. Add `Tags` field to Check struct:
```go
type Check struct {
    // ... existing fields ...
    Tags []string `yaml:"tags,omitempty"`
}
```

2. Add tag validation regex:
```go
var validTag = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
```

3. Add validation in `Validate()`:
```go
for _, tag := range check.Tags {
    if !validTag.MatchString(tag) {
        return &ConfigError{
            Message: fmt.Sprintf("check %q has invalid tag %q: must be lowercase alphanumeric with hyphens", check.ID, tag),
            LineNum: c.FindCheckNodeLine(check.ID, i),
        }
    }
}
```

### Phase 2: CLI Flags

**Files:** `internal/cli/check.go`, `internal/cli/list.go`, `internal/cli/tags.go`

1. Add flags to check command:
```go
var (
    tags        []string
    excludeTags []string
)

func init() {
    checkCmd.Flags().StringSliceVar(&tags, "tags", nil, "Run checks matching ANY of these tags")
    checkCmd.Flags().StringSliceVar(&excludeTags, "exclude-tags", nil, "Exclude checks matching ANY of these tags")
}
```

2. Add flags to list command for filtering display

3. Add new `tags` command:
```go
var tagsCmd = &cobra.Command{
    Use:   "tags",
    Short: "List all unique tags in the configuration",
    RunE:  runTags,
}

func runTags(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(configFile)
    if err != nil {
        return err
    }

    tags := collectUniqueTags(cfg.Checks)
    for _, tag := range tags {
        fmt.Println(tag)
    }
    return nil
}
```

### Phase 3: Orchestrator Filtering

**Files:** `internal/orchestrator/orchestrator.go`

1. Add filter configuration to Orchestrator:
```go
type TagFilter struct {
    Include []string // OR match
    Exclude []string // OR exclusion
}
```

2. Add filter method:
```go
func (o *Orchestrator) filterChecksByTags(checks []config.Check, filter TagFilter) []config.Check
```

3. Add dependency validation that skips checks with missing dependencies and logs warning

### Phase 4: Output Updates

**Files:** `internal/cli/list.go`, `internal/output/formatter.go`, `internal/output/json.go`

1. Update list command to display tags:
```
Checks (4):

  fmt
    Tags:     format, fast, pre-commit
    Command:  gofmt -l .
    Severity: error

  security-scan
    Tags:     security, slow, ci-only
    Command:  gosec ./...
    Severity: error
```

2. Update JSON output to include tags:
```json
{
  "results": [
    {
      "check_id": "fmt",
      "tags": ["format", "fast", "pre-commit"],
      "passed": true,
      "duration_ms": 150
    }
  ]
}
```

### Phase 5: Documentation

**Files:** `README.md`, `docs/ai-assisted-setup.md`, `examples/*.yaml`

1. Update README with tag examples
2. Update examples to use tags instead of comment-based phases
3. Document standard tag conventions
4. Add troubleshooting for tag filtering issues

## Testing Strategy

### Unit Tests

1. **Schema tests** (`internal/config/schema_test.go`)
   - Parse config with tags
   - Validate tag format
   - Handle empty/missing tags

2. **Validation tests** (`internal/config/config_test.go`)
   - Invalid tag format detection
   - Duplicate tag handling (allowed)

3. **Filter tests** (`internal/orchestrator/orchestrator_test.go`)
   - OR filtering (`--tags`)
   - Exclusion (`--exclude-tags`)
   - Combined filtering
   - Empty filter (all checks)
   - No matches (empty result, no error)
   - Non-existent tags (silently ignored)

4. **Dependency skip tests**
   - Skip check when dependency excluded
   - Warning logged for skipped checks
   - Multiple checks skipped in chain

### Integration Tests

1. CLI flag parsing
2. End-to-end tag filtering with real config
3. JSON output includes tags

### Example Test Cases

```go
func TestFilterChecksByTags(t *testing.T) {
    checks := []config.Check{
        {ID: "fmt", Tags: []string{"format", "fast"}},
        {ID: "lint", Tags: []string{"lint", "fast"}},
        {ID: "test", Tags: []string{"test"}},
        {ID: "security", Tags: []string{"security", "slow"}},
    }

    tests := []struct {
        name     string
        filter   TagFilter
        expected []string
    }{
        {"no filter", TagFilter{}, []string{"fmt", "lint", "test", "security"}},
        {"include fast", TagFilter{Include: []string{"fast"}}, []string{"fmt", "lint"}},
        {"include security or test", TagFilter{Include: []string{"security", "test"}}, []string{"test", "security"}},
        {"exclude slow", TagFilter{Exclude: []string{"slow"}}, []string{"fmt", "lint", "test"}},
        {"include fast exclude lint", TagFilter{Include: []string{"fast"}, Exclude: []string{"lint"}}, []string{"fmt"}},
        {"nonexistent tag", TagFilter{Include: []string{"nonexistent"}}, []string{}},
    }
    // ...
}
```

## Migration

### Backwards Compatibility

- Tags field is optional (`yaml:"tags,omitempty"`)
- Existing configs work without modification
- Default behavior (no flags) runs all checks

### Upgrade Path

1. Update vibeguard binary
2. Optionally add tags to existing configs
3. Update CI scripts to use tag filtering

## Resolved Questions

1. **Tag case sensitivity**: Enforce lowercase only for consistency

2. **Maximum tags per check**: No limit, document best practices (3-5 tags)

3. **Reserved tags**: No reserved tags initially

4. **Tag discovery command**: Yes, `vibeguard tags` included in Phase 2

5. **AND logic (`--tags-all`)**: Not needed, OR logic is sufficient

6. **Validation of unused tags**: No validation, silently ignore non-matching filters

7. **Dependency handling**: Skip checks with missing dependencies, log warning

## Success Metrics

1. Users can filter checks by tag without config modification
2. No breaking changes to existing configs
3. Clear error messages for invalid tags
4. Documentation updated with examples

## References

- [golangci-lint tag filtering](https://golangci-lint.run/usage/configuration/#run-configuration)
- [pre-commit hook stages](https://pre-commit.com/#confining-hooks-to-run-at-certain-stages)
- [pytest markers](https://docs.pytest.org/en/stable/example/markers.html)
