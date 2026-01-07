---
summary: Completed comprehensive TAGS.md documentation for tag conventions and usage patterns
event_type: code
sources:
  - docs/specs/SPEC-tags.md
  - README.md
  - docs/TAGS.md
tags:
  - tags-feature
  - documentation
  - vibeguard-935
  - user-guide
  - conventions
  - tag-filtering
---

# Tag Conventions Documentation Completion

## Overview

Completed the task of documenting standard tag conventions for VibeGuard tag filtering feature (vibeguard-935). Created comprehensive `docs/TAGS.md` guide covering all aspects of tag usage, conventions, and best practices.

## What Was Created

### New File: `docs/TAGS.md`

A comprehensive guide document (~450 lines) covering:

#### 1. Standard Tag Conventions (Documented in Table)
- **Execution Performance:** `fast` (<5s), `slow` (>30s)
- **Check Categories:** `format`, `lint`, `test`, `build`, `security`
- **Execution Context:** `pre-commit`, `ci`, `ci-only`
- **Special Capabilities:** `llm` (LLM-powered checks)

#### 2. Detailed Tag Descriptions
Each tag includes:
- Clear description of what it means
- Use cases and typical scenarios
- Real-world examples
- Configuration snippets
- Related tags and patterns

**Example structure for `fast` tag:**
```
- Completion time < 5 seconds
- Use cases: pre-commit, local development, CI fast-path
- Examples: gofmt, go vet, eslint
- Often combined with: pre-commit or category tags
```

#### 3. Usage Patterns
Practical examples for:
- **Local Development:** Fast iteration with `--tags fast`
- **Pre-Commit Hooks:** Running `--tags pre-commit` before commits
- **CI/CD Pipelines:** Comprehensive validation with `--tags ci`
- **Tag Discovery:** Using `vibeguard tags` command

#### 4. Best Practices
Guidelines covering:
- Multi-tag assignment strategies
- Realistic timing classification (when to use fast vs. slow)
- Dependency handling with tags
- Custom tag creation for project-specific needs
- Why AND logic isn't needed (OR is sufficient)

#### 5. Dependency Handling Deep Dive
Explained the strict filtering model where:
- Dependencies are NOT automatically included
- Filtered checks skip if their dependency is excluded
- How to tag dependencies to work together
- Solutions for breaking dependency chains

#### 6. Troubleshooting Guide
Common questions addressed:
- "Why did my slow check run on --tags pre-commit?"
- "Why was my check skipped with 'dependency not in filtered set'?"
- "How do I handle different workflows?"
- "Can I use AND logic for tags?"

## Implementation Details

### Documentation Structure
- **Overview section** - Quick intro to tag benefits
- **Standard conventions table** - Alphabetical reference
- **Detailed sections** - Each tag with deep explanations
- **Usage patterns** - Copy-paste ready examples
- **Best practices** - Prescriptive guidance
- **Troubleshooting** - Common problems solved
- **Reference section** - Links to related docs

### Quality Assurance
✓ All `vibeguard check` validations pass (10/10 checks)
✓ Markdown is properly formatted
✓ All code examples are syntactically valid
✓ Internal cross-references work (spec, README, CLI-REFERENCE)
✓ Examples use actual tag conventions from spec

### Integration Points
- Linked to SPEC-tags.md for implementation details
- Referenced from README.md guidelines
- Supports vibeguard tags discovery command
- Consistent with existing documentation style

## What This Solves

Before: Tag conventions were scattered across README.md and SPEC-tags.md
After: Centralized, comprehensive guide with:
- Quick reference for all tags
- Real-world usage examples
- Best practices guidance
- Troubleshooting help
- Links to related documentation

Users can now:
1. Discover what tags are available: `vibeguard tags`
2. Understand each tag: Read TAGS.md
3. See practical examples: Check "Usage Patterns" section
4. Debug issues: Check "Troubleshooting" section

## Related Work

This documentation completes the Tags Feature Phase 5 deliverables:
- ✓ README.md updates (already done in previous session)
- ✓ Standard tag conventions documented (NEW - this work)
- ✓ examples/basic.yaml with tags (already done)
- ✓ examples/advanced.yaml refactored (already done)
- ✓ docs/ai-assisted-setup.md updated (already done)

## Git Commit

```
docs: Document standard tag conventions and use cases

Add comprehensive TAGS.md documentation covering:
- Standard tag conventions (format, lint, test, security, build, fast, slow, pre-commit, ci, llm)
- Execution context guidelines (pre-commit, ci, ci-only)
- Usage patterns for local development, CI/CD, and pre-commit hooks
- Best practices for tag assignment and dependency handling
- Troubleshooting guide for common tag filtering issues
- Real-world examples for each tag category
```

Commit hash: 379f3ef
All vibeguard checks passed.

## Next Steps

Task vibeguard-935 is now complete. Available work items:
- vibeguard-936: Add tag filtering troubleshooting guide (but coverage here!)
- vibeguard-937: Add explicit test for tag filtering with no matches

Both remaining items are addressed or could leverage this documentation.
