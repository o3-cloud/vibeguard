---
name: log
description: Create log entries with timestamps. Use this when documenting discoveries, research findings, meetings, code reviews, articles, videos, or significant events during development.
---

# Log Entry Skill

## Instructions

When creating a log entry, follow these steps:

1. **Determine the event type** - Choose from: deep dive, meeting, research, code, code review, article, video
2. **Write a summary** - Create a brief, clear summary of the entry
3. **Gather sources** - List relevant references, links, or sources (up to 5)
4. **Add tags** - Tag the entry for easy discovery and organization (up to 10)
5. **Create the file** - The file will be named with a date and reference slug: `YYYY-MM-DD_reference-slug.md`

## Log Entry Structure

Log entries follow this format:

```yaml
---
summary: A brief summary of the log entry
event_type: deep dive | meeting | research | code | code review | article | video
sources:
  - source1
  - source2
  - source3
  - source4
  - source5
tags:
  - tag1
  - tag2
  - tag3
  - tag4
  - tag5
  - tag6
  - tag7
  - tag8
  - tag9
  - tag10
---

[Your detailed entry content here]
```

## Event Types

- **deep dive** - In-depth exploration or analysis of a topic
- **meeting** - Notes from meetings or discussions
- **research** - Research findings or investigations
- **code** - Code changes, implementations, or technical work
- **code review** - Review notes and feedback on code changes
- **article** - Summary of an article or blog post
- **video** - Notes from a video or presentation

## File Naming

Log entries are saved in `docs/log/` with date and reference slug:

**Format:** `YYYY-MM-DD_reference-slug.md`

**Examples:**
- `2025-12-30_authentication-research.md`
- `2025-12-30_meeting-notes.md`
- `2025-12-30_performance-review.md`

The slug should be:
- Lowercase
- Hyphen-separated words
- Descriptive and meaningful (e.g., `auth-framework`, `performance-review`, `bug-investigation`)

This format ensures chronological ordering while keeping entries discoverable by their descriptive slug.

## Content Guidelines

After the YAML frontmatter, include:
- Clear, concise description of the log entry
- Key findings or decisions
- Next steps if applicable
- Related decisions or references to ADRs

## Example

```yaml
---
summary: Evaluated authentication frameworks for the project
event_type: research
sources:
  - https://example.com/auth-comparison
  - https://example.com/jwt-guide
  - https://example.com/oauth2-tutorial
tags:
  - authentication
  - security
  - framework-evaluation
  - oauth2
  - jwt
---

# Authentication Framework Research

Researched three main authentication approaches...

## Key Findings
- Option A provides the best developer experience
- Security model aligns with our requirements
- ...

## Next Steps
- Prototype with Option A
- Get team review
```

## Template Reference

The project uses a standardized template located at `docs/log/TEMPLATE.md` for reference.
