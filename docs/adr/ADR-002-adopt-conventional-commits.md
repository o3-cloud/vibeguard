---
summary: Adopt Conventional Commits specification for commit messages to provide structured, parseable commit history that enables automated changelog generation, version bumping, and clearer project history.
event_type: code
sources:
  - https://www.conventionalcommits.org/
  - https://github.com/commitizen/cz-cli
tags:
  - architecture
  - git
  - commits
  - versioning
  - automation
  - decision
  - developer-experience
---

## Context and Problem Statement

Commit messages are a critical part of project history and maintenance. Currently, vibeguard lacks a structured approach to commits, leading to:

- Inconsistent commit message formats across the codebase
- Difficulty identifying breaking changes from commit history
- Manual changelog creation without automated tooling
- Unclear separation between features, fixes, documentation, and chores
- Lost context when reviewing historical changes
- Inability to automatically determine semantic versioning

As we scale the project with AI agents and multiple contributors, a clear commit convention becomes essential for:
- Automated changelog generation
- Semantic version bumping (major.minor.patch)
- Better code review and git history navigation
- Integration with CI/CD pipelines
- Clearer communication of what changed and why

## Considered Options

### Option A: Adopt Conventional Commits
Follow the Conventional Commits specification with structured format:
```
type(scope): subject

body

footer
```
Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `perf`
- Standardized, widely adopted specification
- Enables automated tooling (changelog generation, version bumping)
- Parser-friendly for automation and CI/CD
- Promotes clarity and consistency
- Integrates well with Commitizen and similar tools

### Option B: Squash Commits + Manual Changelog
Keep current ad-hoc approach but enforce squashing to main:
- Simple to implement
- However: loses granular history, still requires manual changelog, harder to identify breaking changes

### Option C: Semantic Versioning Only
Use semantic versioning (MAJOR.MINOR.PATCH) without structured commits:
- Communicates version intent clearly
- However: doesn't structure commit messages, still requires manual changelog, doesn't aid in reviewing history

## Decision Outcome

**Chosen option: Option A - Adopt Conventional Commits**

**Rationale:**
1. **Automation-Ready**: Enables automated changelog generation, version bumping, and release notes
2. **Clarity**: Type prefix makes it immediately clear what kind of change was made (feat, fix, docs, etc.)
3. **Breaking Changes**: Footer syntax (`BREAKING CHANGE:`) explicitly documents API-breaking changes
4. **Tooling Ecosystem**: Rich ecosystem of tools (Commitizen, standard-version, semantic-release) integrate seamlessly
5. **CI/CD Integration**: Enables intelligent automation in pipelines based on commit types
6. **Industry Standard**: Widely adopted across open source and commercial projects

**Tradeoffs:**
- Requires developer discipline and learning curve
- Commits must follow format (tooling like Commitizen can help)
- Breaking changes must be explicitly documented

## Consequences

### Positive Outcomes
- Automated changelog generation from commit history
- Semantic version bumping based on commit types (feat → minor, fix → patch, BREAKING CHANGE → major)
- Clearer git history for understanding project evolution
- Better integration with CI/CD pipelines and automation
- Easier identification of which commits introduce breaking changes
- Improved code review experience with structured, scannable commit messages
- Facilitates automated tooling and integrations

### Negative Outcomes
- Initial learning curve for team members unfamiliar with the specification
- Requires discipline to maintain consistency (can use commitizen CLI to enforce)
- Enforcing the format may slow down rapid development (mitigated by tooling)

### Implementation Path
1. Document Conventional Commits guidelines for the project
2. Optionally install Commitizen for interactive commit creation
3. Set up pre-commit hooks to validate commit messages
4. Configure automated changelog generation in CI/CD
5. Create breaking change documentation process
6. Educate team on commit conventions and tooling

## Related Decisions
- Works in conjunction with ADR-001 (Beads) for structured task tracking
- Enables automated versioning workflows
- Supports AI agent-generated commits with clear semantics
