---
summary: Implemented YAML line number reporting for config validation errors
event_type: code
sources:
  - internal/config/config.go
  - internal/config/schema.go
  - internal/config/config_test.go
tags:
  - config-validation
  - error-messages
  - yaml
  - line-numbers
  - vibeguard-bd9
  - developer-experience
---

# Config Validation with YAML Line Numbers

## Overview

Successfully implemented YAML line number reporting for config validation errors in vibeguard-bd9. Users can now easily identify the exact line in their config file where validation errors occur, significantly improving the developer experience when debugging configuration issues.

## Problem Statement

Previously, validation errors would report issues like "duplicate check id: test" without any indication of where in the config file the problem occurred. For large configuration files with many checks, users had to manually search through the file to find the problematic check.

## Solution

Implemented a multi-part solution that preserves YAML node information during parsing and uses it to annotate validation errors with line numbers.

### 1. ConfigError Enhancement

Modified the `ConfigError` struct to include location information:

```go
type ConfigError struct {
	Message  string
	Cause    error
	LineNum  int    // Line number in the config file (0 if not available)
	FileName string // File name for reference (optional)
}
```

The Error() method now formats line numbers in error messages when available:
```
Error: duplicate check id: test (line 5)
```

### 2. YAML Node Preservation

Updated the `Load()` function to preserve the YAML node tree during parsing:

```go
var root yaml.Node
if err := yaml.Unmarshal(data, &root); err != nil {
	return nil, &ConfigError{Message: "failed to parse config file", Cause: err}
}

var cfg Config
if err := root.Decode(&cfg); err != nil {
	return nil, &ConfigError{Message: "failed to parse config file", Cause: err}
}

cfg.yamlRoot = &root  // Store for line number lookups
```

### 3. YAML Node Tree Navigation

Added `findCheckNodeLine()` helper that:
- Handles the DocumentNode structure returned by yaml.v3
- Traverses to the checks sequence
- Returns the line number of the check at a given index

Key insight: yaml.v3's Unmarshal to a Node returns a DocumentNode as the root, with the actual mapping in Content[0].

### 4. Validation Error Updates

Updated all validation errors to include line numbers:
- Duplicate check IDs
- Missing required fields
- Invalid severity values
- Self-referencing requirements
- Unknown check references
- Cyclic dependencies

Each error now calls `findCheckNodeLine()` to include the line number.

## Testing

Added comprehensive test coverage:

1. **TestLoad_ConfigErrorWithLineNumbers** - Verifies duplicate check errors include line numbers
2. **TestLoad_CyclicDependencyWithLineNumbers** - Verifies cyclic dependency errors include line numbers

Both tests confirm:
- LineNum field is non-zero
- Error message includes "line" text
- Line numbers are accurate

All 50 existing config tests continue to pass, confirming backward compatibility.

## Example Usage

Before:
```
Error: configuration validation failed: duplicate check id: test
```

After:
```
Error: configuration validation failed: duplicate check id: test (line 5)
```

## Technical Details

### YAML Node Structure in yaml.v3

When unmarshaling into a `yaml.Node`, the structure is:
- Root: DocumentNode (Kind 1)
- Content[0]: The actual mapping (Kind 4 = MappingNode)
- The mapping contains alternating key/value pairs
- For "checks" key, the value is a SequenceNode with check items

### Line Number Semantics

- Line numbers are 1-based (matching editor conventions)
- Each element in the sequence has a Line property
- The check item's Line is the line of the first property in that check's mapping

## Impact

- **Developer Experience**: Users can immediately locate config errors
- **Debugging Time**: Significant reduction in time spent searching config files
- **Error Quality**: Errors are now actionable with clear location context

## Related Issues

- Closes: vibeguard-bd9 "Config validation errors don't show YAML line numbers"
- Related: vibeguard-trb "Grok and assert errors lack file/line context" (future work)
- Related: vibeguard-8p1 "Distinguish config errors from execution errors" (already implemented)
