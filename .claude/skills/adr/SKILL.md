---
name: adr
description: Create Architecture Decision Records (ADRs) following the MADR template. Use this when documenting significant architectural decisions, design choices, or major technical direction changes.
---

# Architecture Decision Record (ADR)

## Instructions

When creating an ADR, follow these steps:

1. **Get the context** - Understand the problem that requires an architectural decision
2. **Identify alternatives** - Gather options that were considered
3. **Document the decision** - Record what was chosen and why
4. **Capture consequences** - Note the impacts and tradeoffs
5. **Update CLAUDE.md** - Add a reference to the new ADR in the "Current ADRs" section of `CLAUDE.md` with a brief summary and link

## ADR Structure

Use the MADR (Markdown Architecture Decision Records) template with these sections:

### Title
A concise heading like: "ADR-001: {Decision Title}" (e.g., "ADR-001: Use React for Frontend Framework")

### Context and Problem Statement
- Explain the situation and constraints
- State the specific issue being addressed
- Include relevant background information

### Considered Options
List alternatives that were evaluated:
- Option A: Description and key characteristics
- Option B: Description and key characteristics
- Option C: Description and key characteristics

### Decision Outcome
- **Chosen option:** State which option was selected
- **Rationale:** Explain why this option was chosen
- **Tradeoffs:** Document what was gained and sacrificed

### Consequences
- **Positive outcomes:** Benefits of the decision
- **Negative outcomes:** Drawbacks or risks
- **Neutral impacts:** Side effects or changes in approach

## File Naming

Save ADRs in `docs/adr/` with sequential numbering:
- `ADR-001-framework-choice.md`
- `ADR-002-database-architecture.md`
- `ADR-003-authentication-strategy.md`

## Syncing with CLAUDE.md

Every new ADR must be registered in `CLAUDE.md` to maintain a discoverable index of all architectural decisions. When you create an ADR:

1. Save the ADR file to `docs/adr/ADR-NNN-title.md`
2. Update the "Current ADRs" section in `CLAUDE.md` with an entry like:

```markdown
- **[ADR-NNN: {Title}](docs/adr/ADR-NNN-title.md)**
  - Brief description of what was decided
  - Why this matters to the project
  - Key tradeoff or benefit
```

This keeps `CLAUDE.md` as the single source of truth for discovering what architectural decisions have been made.

## Template Reference

The project uses the [MADR bare minimal template](https://raw.githubusercontent.com/adr/madr/refs/heads/develop/template/adr-template-bare-minimal.md) located at `docs/adr/TEMPLATE.md`.

## Example

A complete ADR typically spans 200-500 words and addresses one significant decision. When you need to create an ADR, provide me with:
- The decision title
- The context/problem
- The alternatives considered
- Your chosen option
- Key reasoning and consequences

I'll help you structure it following the MADR format.
