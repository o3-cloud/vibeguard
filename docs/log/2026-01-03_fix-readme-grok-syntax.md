---
summary: Fixed incorrect grok syntax documentation in README schema
event_type: code
sources:
  - README.md line 164
  - examples/go-project.yaml
  - examples/advanced.yaml
  - vibeguard.yaml
tags:
  - documentation
  - grok
  - bug-fix
  - readme
  - syntax
---

# Fix README Grok Syntax Example

## Issue
Task `vibeguard-iro`: README grok syntax example was incorrect in the configuration schema documentation.

## Problem
The README.md schema section (line 164) showed:
```yaml
grok:
  - pattern_name: pattern
```

This suggested grok patterns were specified as a key-value map with `pattern_name: pattern` format, which was incorrect.

## Solution
Changed the documentation to show the correct syntax:
```yaml
grok:
  - pattern_string
```

The `grok` field is a list of pattern strings, not a map. Each pattern is a string containing the grok pattern expression.

## Verification
Verified the fix against existing examples in the codebase:
- `examples/go-project.yaml` line 50: `- total:.*\(statements\)\s+%{NUMBER:coverage}%`
- `examples/advanced.yaml`: Multiple grok pattern examples all use string syntax
- `vibeguard.yaml` line 38: `- total:.*\(statements\)\s+%{NUMBER:coverage}%`

All real-world examples in the codebase use the correct list-of-strings syntax.

## Changes Made
- Line 164 in README.md: Changed `pattern_name: pattern` to `pattern_string`
- This is a documentation-only fix with no code changes

## Impact
This fix corrects the user-facing documentation to accurately reflect the actual API, preventing confusion when users try to configure grok patterns.
