---
summary: Exploration of vibeguard project infrastructure to plan Dockerfile implementation for Claude Code and Beads
event_type: deep dive
sources:
  - docs/log/2026-01-05_docker-claude-code-setup.md
  - .goreleaser.yaml
  - .github/workflows/ci.yml
tags:
  - docker
  - infrastructure
  - planning
  - claude-code
  - beads
  - ci-cd
  - goreleaser
  - deployment
---

# Dockerfile Planning: Project Infrastructure Analysis

## Overview
Explored the vibeguard project infrastructure to understand the current deployment model and determine the best approach for implementing a Dockerfile that includes Claude Code and Beads.

## Current State

**Docker Infrastructure:** None exists currently in the repository.

**Build & Distribution Model:**
- Go-based project (v1.24.4)
- Single-binary distribution via GoReleaser
- Cross-platform builds: Linux, macOS, Windows
- Multi-architecture: amd64, arm64
- Homebrew tap publishing configured

**CI/CD Infrastructure:**
- GitHub Actions workflows in `.github/workflows/`
- `ci.yml` - Main CI (lint, test, vibeguard checks)
- `release.yml` - Automated releases with GoReleaser
- `mutation.yml` - Weekly mutation testing
- `commitlint.yml` - Conventional commit validation

## Directory Structure Analysis

**Current Layout:**
```
vibeguard/
├── .github/workflows/    # CI/CD
├── .beads/               # Task management
├── cmd/vibeguard/        # CLI entry
├── internal/             # Core implementation
├── docs/                 # Documentation
│   ├── adr/              # Architecture decisions
│   └── log/              # Research logs
├── examples/             # Config examples
├── dist/                 # Built binaries
└── .goreleaser.yaml      # Release config
```

**Dockerfile Location Options:**
1. **Project root** (`/Dockerfile`) - Simplest, standard convention
2. **`/docker/`** - Organized, allows multiple Dockerfiles
3. **`/build/`** - Conventional for build artifacts
4. **`/deploy/`** - Common for deployment configs

## Prior Research Available

Comprehensive Docker research was completed earlier (2026-01-05) and documented at:
`docs/log/2026-01-05_docker-claude-code-setup.md`

**Key findings from prior research:**
- Two Dockerfile options designed (single-stage and multi-stage)
- Base image: `ubuntu:24.04` (LTS, pinned)
- Dependencies: curl, bash, ca-certificates, git
- Claude Code: curl-based native installer
- Beads: pre-built binary or go install
- Security: non-root user creation
- Credential strategy: mount `~/.claude.json` (not bake into image)

## Implementation Considerations

### Dependencies Required
| Dependency | Purpose |
|------------|---------|
| curl | Download installers |
| bash | Shell environment |
| ca-certificates | HTTPS connections |
| git | Beads git-backed storage |

### Volume Mount Strategy
| Volume | Container Path | Purpose |
|--------|----------------|---------|
| `~/.claude.json` | `/home/claude/.claude.json` | Claude Code credentials |
| `~/.claude/` | `/home/claude/.claude/` | Claude Code state |
| Project dir | `/home/claude/project` | Working directory |

### Security Best Practices
- Non-root user (`claude`)
- Read-only credential mounts (`:ro`)
- No credentials baked into image
- Minimal package installation

## CI/CD Integration Opportunity

The existing CI/CD could be extended:
1. Add Docker build step to `ci.yml`
2. Publish to GitHub Container Registry (GHCR)
3. Multi-platform builds (linux/amd64, linux/arm64)
4. Tag with semantic versions

## Key Decisions Needed

1. **Dockerfile location** - Root vs `/docker/` directory
2. **Build approach** - Single-stage (simple) vs multi-stage (optimized)
3. **Beads installation** - Pre-built binary vs go install
4. **CI integration** - Whether to add Docker build to workflows
5. **Registry publishing** - GHCR, Docker Hub, or none

## Related ADRs

- **ADR-001**: Beads adoption (core to this Docker setup)
- **ADR-003**: Go implementation (relevant for Beads build)
- **ADR-005**: VibeGuard policy enforcement
- **ADR-006**: Git pre-commit hooks

## Next Steps

1. Decide on Dockerfile location and approach
2. Implement Dockerfile based on prior research
3. Test credential mounting and Claude Code authentication
4. Optionally integrate Docker build into CI/CD
