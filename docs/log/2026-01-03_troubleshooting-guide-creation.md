---
summary: Created comprehensive TROUBLESHOOTING.md guide covering common VibeGuard issues and solutions
event_type: code
sources:
  - docs/adr/ADR-003-adopt-golang.md
  - README.md
  - internal/executor/executor.go
  - internal/grok/grok.go
tags:
  - documentation
  - troubleshooting
  - user-support
  - common-errors
  - configuration
  - grok-patterns
  - debugging
  - error-handling
---

# VibeGuard Troubleshooting Guide Creation

## Summary

Created a comprehensive TROUBLESHOOTING.md document that addresses the 80% of support burden identified in task vibeguard-zfq. The guide covers the most common issues users encounter when setting up and running VibeGuard.

## Coverage

The troubleshooting guide includes:

1. **Configuration Errors**
   - Configuration file not found
   - YAML validation failures
   - Invalid field types and required fields

2. **Cyclic Dependency Errors**
   - Detection and explanation of circular dependencies
   - Examples of cyclic dependencies
   - Solutions for breaking cycles

3. **Grok Pattern Errors**
   - Failed pattern compilation
   - Pattern matching failures
   - Undefined pattern references
   - Debugging strategies and working examples

4. **Assertion Errors**
   - Invalid assertion syntax
   - Type mismatches
   - Common assertion patterns
   - Field reference issues

5. **Timeout Errors (Exit Code 4)**
   - Timeout cause analysis
   - Solutions for increasing timeouts
   - Command optimization strategies

6. **Check Not Found Errors**
   - Identification of missing checks
   - Check ID case sensitivity
   - Configuration validation

7. **Exit Codes Reference**
   - Complete mapping of exit codes 0-4
   - What each code means
   - Actions to take for each code

8. **Debugging Strategies**
   - Verbose output usage
   - JSON output for parsing
   - Single check isolation
   - Manual command testing

## Key Implementation Details

- Organized by error type with clear, actionable solutions
- Includes code examples and YAML snippets
- References to relevant documentation sections
- Step-by-step troubleshooting procedures
- Common patterns and anti-patterns
- Links to configuration schema and references

## Design Decisions

- Structured for scanability with clear headers and bullet points
- Placed at project root level (TROUBLESHOOTING.md) for easy discovery
- Cross-references README.md for schema details
- Includes both theory (why errors occur) and practice (how to fix them)
- Provides debugging tools and strategies rather than just solutions

## Impact

This documentation should significantly reduce support burden by:
- Providing self-service solutions for common issues
- Enabling faster troubleshooting for both users and developers
- Reducing repeated questions about grok patterns and timeouts
- Improving overall user experience during initial setup

## Next Steps

- Monitor for additional common issues not covered in current version
- Consider adding video tutorials or visual diagrams for complex concepts
- Link from main README.md to this guide
- Collect user feedback on coverage and clarity
