---
summary: Adopt Beads (bd) as a git-backed, distributed issue tracker designed specifically for AI agents to maintain persistent, structured task management with dependency tracking across sessions.
event_type: code
sources:
  - https://github.com/steveyegge/beads
  - https://github.com/steveyegge/beads/tree/main/docs
  - CLAUDE.md
tags:
  - architecture
  - beads
  - task-management
  - ai-agents
  - decision
  - git-native
  - persistence
---

## Context and Problem Statement

As we work with AI coding agents (Claude Code and other agents) on complex, long-horizon tasks, we face a critical challenge: agents lose context across sessions and struggle to maintain coherent progress on multi-step work. Currently, we rely on:

- Markdown-based todo lists that become stale
- Manual plan updates in the CLI
- Fragmented task tracking across multiple formats
- No structured dependency management between tasks
- Difficulty recovering context when agents resume work

This fragmentation increases cognitive overhead, causes agents to re-discover information, and makes it hard to track what's been done versus what remains. We need a persistent, structured memory system that works naturally with how agents operate.

## Considered Options

### Option A: Adopt Beads (bd)
A git-backed, distributed issue tracker designed specifically for AI agents. Features include:
- JSONL-based storage in `.beads/` directory (versioned with git)
- Hash-based task IDs preventing merge conflicts
- Dependency graphs for blocking relationships
- JSON output optimized for agent consumption
- Semantic summarization for context compression
- SQLite caching for performance

### Option B: Use GitHub Issues + Automation
Leverage GitHub's native issue tracking with:
- Structured labels and projects
- Built-in API and automation
- Team collaboration features
- However: requires external platform, slower for agent consumption, not git-native

### Option C: Continue Current Approach
Stick with markdown-based todo tracking in CLI with manual updates:
- Low friction to start
- However: scales poorly, loses context, agents don't naturally maintain it, no dependency tracking

## Decision Outcome

**Chosen option: Option A - Adopt Beads**

**Rationale:**
1. **Agent-Optimized**: Beads is purpose-built for AI agents with JSON output, dependency awareness, and ready-task identification
2. **Git-Native**: Tasks are stored as JSONL files, making them version-controlled alongside code—enabling branching, merging, and history
3. **Scalability**: As tasks grow complex, dependency tracking and context compression prevent context window bloat
4. **Conflict-Free**: Hash-based IDs solve merge conflicts that would occur with other approaches
5. **Persistence**: Unlike in-memory plans, beads maintains context across agent sessions automatically

**Tradeoffs:**
- New tool to learn and integrate into workflows
- Requires `.beads/` directory in repository
- Developers need to adopt `bd` commands and mindset

## Consequences

### Positive Outcomes
- AI agents can maintain persistent, structured task lists across sessions
- Dependency management prevents "what do I do next?" ambiguity
- Context compression helps manage token budgets for long projects
- Git-native approach aligns with developer workflows
- Reduced manual context-switching between sessions

### Negative Outcomes
- Learning curve for team members unfamiliar with beads
- Additional `.beads/` directory in repository (small footprint)
- Requires initial setup and adoption of new command patterns

### Implementation Path
1. Initialize beads: `bd init` ✓ Completed
2. Migrate existing plans to beads tasks
3. Document beads workflows for the team
4. Integrate with existing Claude Code processes
5. Monitor adoption and refine as needed

## Related Decisions
- Will inform how Claude Code agents structure long-running projects
- May influence how we document complex features requiring multiple-agent coordination
