# Working with Claude Code on Vibeguard

This document describes how Claude Code agents work with the vibeguard project, including architectural decisions and project conventions.

## Architecture Decisions

Major architectural decisions are documented as Architecture Decision Records (ADRs) in the `docs/adr/` directory. These records explain not just what we decided, but *why* and *what tradeoffs* we accepted.

### Current ADRs

- **[ADR-001: Adopt Beads for AI Agent Task Management](docs/adr/ADR-001-adopt-beads.md)**
  - Decision to use Beads (bd) as a git-backed, distributed issue tracker for AI agents
  - Addresses the problem of agents losing context on long-horizon tasks
  - Provides persistent, structured task management with dependency tracking

- **[ADR-002: Adopt Conventional Commits](docs/adr/ADR-002-adopt-conventional-commits.md)**
  - Decision to structure commit messages using the Conventional Commits specification
  - Enables automated changelog generation and semantic versioning
  - Improves code review experience and git history clarity

- **[ADR-003: Adopt Go as the Primary Implementation Language](docs/adr/ADR-003-adopt-golang.md)**
  - Decision to implement VibeGuard in Go for single-binary deployment and performance
  - Ensures frictionless integration into CI/CD pipelines and agent loops with minimal overhead
  - Aligns with cloud-native DevOps tooling ecosystem and enables strong CLI integration

## Project Skills

The `.claude/skills/` directory contains Claude Code skills that automate common workflows:

- **`adr`** - Create Architecture Decision Records following the MADR template
  - Use when: documenting significant architectural decisions, design choices, or major technical direction changes
  - Location: `.claude/skills/adr/SKILL.md`

- **`log`** - Log related skill
  - Location: `.claude/skills/log/SKILL.md`

## Key Tools & Setup

### Beads Configuration
Once beads is adopted, initialize with:
```bash
bd init
```

This creates a `.beads/` directory for storing tasks as git-versioned JSONL files.

### Creating ADRs
To create new Architecture Decision Records:
1. Ask Claude Code: "Create an ADR for [decision]"
2. Claude will guide you through the MADR template structure
3. Save to `docs/adr/ADR-NNN-decision-title.md` with sequential numbering

## Conventions

- **ADR Numbering**: Sequential (ADR-001, ADR-002, etc.)
- **ADR Location**: `docs/adr/`
- **ADR Template**: MADR format (reference: `docs/adr/TEMPLATE.md`)
- **Skills Location**: `.claude/skills/{skill-name}/SKILL.md`

## References

- [MADR (Markdown Architecture Decision Records)](https://adr.github.io/madr/)
- [Beads Repository](https://github.com/steveyegge/beads)
- [Claude Code Documentation](https://docs.claude.com/claude-code)
