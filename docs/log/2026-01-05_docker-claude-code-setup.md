---
summary: Researched Docker image setup with Ubuntu latest, Claude Code, and Beads (git-backed task management) installation requirements, best practices, and implementation strategy
event_type: research
sources:
  - https://code.claude.com/docs/en/setup
  - https://docs.docker.com/build/building/best-practices/
  - https://hub.docker.com/_/ubuntu
  - https://github.com/steveyegge/beads
tags:
  - docker
  - containerization
  - claude-code
  - beads
  - task-management
  - ubuntu
  - devops
  - installation
  - dockerfile
  - infrastructure
  - go
---

# Docker + Claude Code + Beads Setup Research

## Overview
Researched the feasibility and best practices for creating a Docker image based on Ubuntu (latest) with Claude Code and Beads (distributed task management) installed and ready to use.

## Claude Code System Requirements

**Supported Operating Systems:**
- macOS 10.15+
- Ubuntu 20.04+ / Debian 10+
- Windows 10+ (requires WSL 1, WSL 2, or Git for Windows)

**Hardware & Dependencies:**
- Minimum 4GB RAM
- Internet connection (required for authentication and API calls)
- Shell environment: Bash, Zsh, or Fish
- ripgrep (usually included with Claude Code)
- Node.js 18+ (only required for npm installation method)
- Must be in an Anthropic-supported country

## Installation Methods for Docker

### 1. Native Installation (Recommended for Docker)
```bash
curl -fsSL https://claude.ai/install.sh | bash
```
**Advantages:** No external dependencies, lightweight, official method
**Best for:** Container environments

### 2. NPM Installation
```bash
npm install -g @anthropic-ai/claude-code
```
**Requires:** Node.js 18+
**Less optimal:** Adds Node.js dependency to image

### 3. Homebrew Installation
```bash
brew install --cask claude-code
```
**Not suitable:** Homebrew is macOS-only

## Beads System Requirements & Installation

**What is Beads:**
- Git-backed distributed issue tracker for AI agents (per ADR-001)
- Enables persistent task management with dependency tracking
- Addresses context loss in long-horizon AI agent tasks

**Supported Operating Systems:**
- Linux (glibc 2.32+)
- macOS
- Windows

**Language & Implementation:**
- Primarily written in Go (93.9% of codebase)
- Single binary distribution (ideal for containers)

**Installation Methods for Docker:**

### 1. Go Installation (Recommended for Dockerfile)
```bash
go install github.com/steveyegge/beads/cmd/bd@latest
```
**Advantages:** Builds from source, lightweight, integrates well in multi-stage Docker builds
**Requires:** Go toolchain during build

### 2. NPM Installation
```bash
npm install -g @beads/bd
```
**Requires:** Node.js
**Less suitable:** Adds Node.js dependency

### 3. Homebrew Installation
```bash
brew install steveyegge/beads/bd
```
**Not suitable:** Homebrew is macOS-only

### 4. Pre-built Binary (Simplest for Containers)
Download from GitHub releases and add to PATH

**Post-Installation Setup:**
After installing Beads, initialize configuration:
```bash
bd init
```
This creates `.beads/` directory structure for git-versioned task tracking.

## Docker Best Practices Applied

**Base Image Selection:**
- Use `ubuntu:24.04` (latest LTS) instead of `ubuntu:latest`
- Pin version explicitly for reproducibility and stability
- Ubuntu 24.04 is actively supported with long-term maintenance

**Layer Optimization:**
- Combine `apt-get update` and `apt-get install` in single RUN command
- Prevents stale package cache issues
- Reduces final image layers

**Package Management:**
- Use `--no-install-recommends` flag to exclude unnecessary packages
- Minimize image size (improves security and portability)
- Clean up apt cache: `rm -rf /var/lib/apt/lists/*`

**Security Best Practice:**
- Create non-root user for running Claude Code
- Prevents privilege escalation risks
- Follows container security standards

## Reference Dockerfile Implementations

### Option 1: Single-Stage Build (Simpler)

```dockerfile
FROM ubuntu:24.04

# Install dependencies in one layer
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    bash \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Claude Code
RUN curl -fsSL https://claude.ai/install.sh | bash

# Install Beads from pre-built binary
RUN curl -fsSL https://api.github.com/repos/steveyegge/beads/releases/latest | \
    grep '"browser_download_url"' | \
    grep -E 'linux.*x86_64' | \
    cut -d'"' -f4 | \
    head -1 | \
    xargs curl -fsSL -o /tmp/bd && \
    chmod +x /tmp/bd && \
    mv /tmp/bd /usr/local/bin/bd

# Set up non-root user (security best practice)
RUN useradd -m -s /bin/bash claude
USER claude

WORKDIR /home/claude

# Initialize Beads configuration
RUN bd init

ENTRYPOINT ["/bin/bash"]
```

### Option 2: Multi-Stage Build (Optimized Size)

```dockerfile
# Build stage - compiles Go tools
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

# Build Beads from source
RUN go install github.com/steveyegge/beads/cmd/bd@latest

# Runtime stage - minimal final image
FROM ubuntu:24.04

# Install only runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    bash \
    ca-certificates \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy Beads binary from builder stage
COPY --from=builder /root/go/bin/bd /usr/local/bin/bd

# Install Claude Code
RUN curl -fsSL https://claude.ai/install.sh | bash

# Set up non-root user
RUN useradd -m -s /bin/bash claude
USER claude

WORKDIR /home/claude

# Initialize Beads
RUN bd init

ENTRYPOINT ["/bin/bash"]
```

**Usage with Credential Mounting:**
```bash
# Run Claude Code with host credentials
docker run -it \
  -v ~/.claude.json:/home/claude/.claude.json \
  -v ~/.claude:/home/claude/.claude \
  image-name claude

# Or run Beads commands
docker run -it \
  -v ~/.claude.json:/home/claude/.claude.json \
  -v /path/to/project:/home/claude/project \
  image-name bd ready
```

**Key Design Decisions:**

**Common Elements:**
- `ca-certificates`: Required for HTTPS connections during authentication
- `git`: Required for Beads (git-backed task storage)
- Non-root user: Improves security posture
- `bd init`: Initializes .beads/ directory structure

**Option 1 (Single-Stage):**
- Simpler, easier to debug and modify
- Uses pre-built Beads binary (faster build)
- Slightly larger final image

**Option 2 (Multi-Stage):**
- Smaller final image (Go build tools not included)
- Builds Beads from source (control over version)
- More complex, but follows Docker best practices
- Better for production deployments

**ENTRYPOINT Decision:**
- Changed from `claude` to `/bin/bash` to allow flexibility for both Claude Code and Beads usage
- Users can launch as: `docker run -it image-name claude` or `docker run -it image-name bd ready`

## Important Considerations

### Authentication & Runtime

**Credential Management Strategy:**
- Claude Code stores credentials in `~/.claude.json` (user's home directory)
- Mount this file into container for persistent authentication
- Avoids re-authentication on every container run
- Credentials remain on host machine (not in image)

**Recommended Docker Run Command:**
```bash
docker run -it \
  -v ~/.claude.json:/home/claude/.claude.json \
  -v ~/.claude:/home/claude/.claude \
  image-name
```

**Alternative with Project Mounting:**
```bash
docker run -it \
  -v ~/.claude.json:/home/claude/.claude.json:ro \
  -v ~/.claude:/home/claude/.claude:ro \
  -v /path/to/project:/home/claude/project \
  image-name
```

**Important Notes:**
- Mount as read-only (`:ro`) on non-interactive hosts for security
- Ensure container user (claude) has permissions to read mounted files
- `.claude/` directory contains session history and other state
- OAuth flow typically already completed on host, so container can use cached credentials

**Other Considerations:**
- Still requires interactive TTY (keep `-it` flags)
- Network connectivity required for API calls
- Requires valid Anthropic account with active billing

### Network Requirements
- Internet connectivity required for all Claude Code operations
- API calls to Anthropic servers (can't work offline)
- Consider network policies in containerized environments

### Image Size
- Minimal dependencies = smaller image footprint
- Faster pulls and deployments
- Reduced attack surface

### Signal Handling
- Long-running Claude Code sessions need proper TTY signal handling
- Docker default ENTRYPOINT structure handles this correctly

### Beads-Specific Considerations

**Git Integration:**
- Beads requires functional git in the container
- `.beads/` directory is initialized with `bd init`
- Works best when container directory is a git repository
- Task tracking depends on git commits for versioning

**Volume Mounting:**
- `.beads/` directory should persist (use volume mount or bind mount)
- Project root must be git repository for Beads to function
- Recommended: `docker run -v /path/to/project:/home/claude/work image-name`

**Workflow Integration:**
- Beads and Claude Code work together for AI agent task management
- Claude Code can execute `bd` commands within container
- Persistent task state requires proper volume management

## Potential Challenges

1. **Interactive Authentication**: OAuth flow in containerized environment may be complex
2. **TTY Requirements**: Claude Code expects proper terminal connection
3. **Billing Verification**: Need active Anthropic account setup before use
4. **Network Dependencies**: Can't function in air-gapped environments
5. **Git Repository Requirement**: Beads requires container to be inside a git repository
6. **Volume Persistence**: `.beads/` directory must persist across container restarts
7. **Beads Binary Download**: Pre-built binary fetch depends on GitHub API availability

## Next Steps

1. Build proof-of-concept Dockerfile with both Options (single-stage and multi-stage)
2. Test credential mounting strategy: `docker run -it -v ~/.claude.json:/home/claude/.claude.json image-name claude`
3. Verify Claude Code authentication works with mounted `.claude.json`
4. Test Beads initialization and basic commands with mounted credentials (`bd ready`, `bd list`)
5. Evaluate mounting local projects/codebase volumes for persistent task tracking
6. Verify git integration with Beads task persistence
7. Test combined workflow: Claude Code + Beads in same container
8. Document volume mount best practices and permissions handling
9. Create docker-compose.yml template for typical development workflow

## Related Decisions

- **ADR-001**: Beads adoption for AI agent task management (core to this Docker setup)
- **ADR-003**: Go as primary implementation language (relevant for Beads, written in Go)
- **ADR-005**: VibeGuard for policy enforcement (Docker images can enforce via Claude Code + Beads)
- **ADR-006**: Git pre-commit hooks for policy enforcement (may interact with Claude Code in containers)

---

## Implementation Status

**Implemented: 2026-01-05**

Based on this research, a Dockerfile was created at the project root using the **single-stage build approach**.

**Files Created:**
- `/Dockerfile` - Single-stage Ubuntu 24.04 image with Claude Code and Beads
- `/.dockerignore` - Excludes unnecessary files from build context

**Build & Run:**
```bash
# Build
docker build -t claude-code-beads .

# Run with credential mounting
docker run -it \
  -v ~/.claude.json:/home/claude/.claude.json \
  -v ~/.claude:/home/claude/.claude \
  -v $(pwd):/home/claude/project \
  claude-code-beads
```
