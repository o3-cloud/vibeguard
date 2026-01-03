---
summary: Documented all assertion expression operators in README with comprehensive examples
event_type: code
sources:
  - internal/assert/parser.go
  - internal/assert/lexer.go
  - internal/assert/token.go
  - README.md
tags:
  - documentation
  - assertion-expressions
  - operators
  - readme
  - developer-experience
  - grok-patterns
---

# Assertion Expression Operators Documentation

## Summary

Completed vibeguard-aq5 task by adding comprehensive documentation for assertion expression operators in README.md. Previously, only `coverage >= 80` was shown as an example. Now developers have clear reference documentation for all supported operators, operator types, and practical examples.

## Work Completed

### Documentation Added to README.md

Added a new "Assertion Expression Operators" section that documents:

1. **Comparison Operators** (6 operators)
   - `>=`, `>`, `<=`, `<`, `==`, `!=`
   - Supports both numeric and string comparisons
   - Includes practical examples

2. **Logical Operators** (3 operators)
   - `&&` (AND)
   - `||` (OR)
   - `!` (NOT)
   - Enables complex multi-condition assertions

3. **Arithmetic Operators** (4 operators)
   - `+`, `-`, `*`, `/`
   - Useful for normalization and scaling operations
   - Works in combination with comparisons

4. **Literals and Values**
   - Numeric: integers and floats (e.g., `80`, `0.95`)
   - Strings: single or double quoted (e.g., `"ok"`, `'pass'`)
   - Booleans: `true` and `false`
   - Variables: grok pattern names
   - Grouping: parentheses for explicit precedence

### Practical Examples

Added 7 real-world examples covering:
- Numeric comparison (coverage thresholds)
- String comparison (lint status)
- Logical AND (multi-gate quality checks)
- Logical OR (fail-safe metrics)
- Logical NOT (inverted conditions)
- Arithmetic operations (normalized scoring)
- Complex expressions (combined conditions)

## Key Findings

### Operator Support Analysis

From code review of `internal/assert/`:
- **Parser** (parser.go): Implements recursive descent parser with proper operator precedence
- **Lexer** (lexer.go): Tokenizes all operators correctly
- **Token Types** (token.go): Well-defined token constants for all operators

### Documentation Gaps Identified

The assertion expression system is feature-rich but was severely under-documented:
- **Before**: Only simple `>=` example shown
- **After**: Comprehensive reference with operator categories, syntax, and 7 examples

This was causing knowledge gaps for users trying to write complex policy assertions.

## Quality Assurance

All vibeguard checks pass:
- ✓ vet (0.5s)
- ✓ fmt (0.0s)
- ✓ lint (0.9s)
- ✓ test (4.0s)
- ✓ test-coverage (3.5s)
- ✓ build (0.2s)

## Impact

This documentation improvement:
- **Reduces learning curve** for users writing custom assertion expressions
- **Prevents errors** by showing correct syntax upfront
- **Demonstrates capabilities** that users may not have known existed (arithmetic ops, complex logical expressions)
- **Improves IDE/LLM understanding** as LLMs now have clear reference material for generating assertion expressions

## Related

- Task: vibeguard-aq5
- Related docs: Grok Pattern Extraction section (README.md)
- Related code: internal/assert/ package

## Notes

The assertion expression language is quite powerful but was buried in the codebase without clear documentation. This documentation gap likely led to underutilization of features like arithmetic operators and complex logical expressions. The examples provided should help users construct sophisticated multi-condition checks.
