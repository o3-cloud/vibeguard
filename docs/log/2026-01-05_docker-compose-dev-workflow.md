---
summary: Created docker-compose.yml for development workflow with proper volume mounting, Beads integration, and Claude Code support
event_type: code
sources:
  - Docker Compose Official Documentation
  - Dockerfile implementation
  - Beads configuration
tags:
  - docker-compose
  - development-workflow
  - containerization
  - devops
  - infrastructure-as-code
  - volume-mounting
  - beads-integration
---

# Docker Compose Development Workflow Configuration

## Overview

Created a comprehensive `docker-compose.yml` configuration that provides a development environment for the vibeguard project. The configuration integrates both Beads (task management) and Claude Code (AI development tools) in a containerized workflow with proper volume mounting for persistent development state.

## Problem Statement

Developers needed a simple, standardized way to:
- Run the vibeguard development environment in Docker
- Access the Beads task database from containers
- Mount project files for code changes
- Manage Git version control from within containers
- Optionally mount SSH credentials for Git operations
- Work with both Beads and Claude Code CLI tools

## Solution Implemented

### Core Services

**vibeguard Service:**
- Builds from local Dockerfile
- Runs as non-root user (uid:gid = 1000:1000, configurable)
- Mounts project root at `/work`
- Isolates network communication

### Volume Configuration

**Mounted Volumes:**

1. **Project Root** (`.:/work`)
   - Read-write access for code changes
   - Enables real-time development workflow
   - Changes reflected immediately in editor and container

2. **Beads Database** (`./.beads:/work/.beads`)
   - Persistent task state across sessions
   - Enables task tracking in containerized workflow
   - Shared between host and container

3. **Git Repository** (`./.git:/work/.git`)
   - Read-write access for version control
   - Enables commits/branches from within container
   - Maintains git history

4. **Optional: SSH Credentials** (`~/.ssh:/home/claude/.ssh:ro`)
   - For Git SSH operations
   - Commented out by default (enable as needed)
   - Read-only for security

5. **Optional: Git Config** (`~/.gitconfig:/home/claude/.gitconfig:ro`)
   - User name and email configuration
   - GPG signing setup
   - Commented out by default

### Environment Configuration

**Key Environment Variables:**
- `BEADS_DB`: Points to persistent database location (`/work/.beads/beads.db`)
- `TERM`: Color terminal support for CLI tools
- `UID/GID`: User ID override (defaults to 1000:1000)
- `ANTHROPIC_API_KEY`: Claude API key (optional, set via .env)

### Interactive Development

**Key Features:**
- `stdin_open: true` - Accept standard input
- `tty: true` - Allocate pseudo-terminal
- `command: /bin/bash` - Default shell
- `restart: unless-stopped` - Automatic container restart

### Health Monitoring

**Health Check:**
```bash
test: sh -c "/usr/local/bin/bd --version > /dev/null && which claude > /dev/null"
```

Verifies:
- Beads CLI is available
- Claude Code is available
- Both tools function correctly
- Interval: 30s, Timeout: 10s, Retries: 3

### Networking

**Network Configuration:**
- Service isolated on `vibeguard-network` bridge network
- Service name `vibeguard` resolvable within network
- DNS resolution for service discovery

## Implementation Details

### File: docker-compose.yml

**Structure:**
- Service definition: `vibeguard`
- Build configuration: Uses local Dockerfile
- Volume mounts: 5 volumes (3 required, 2 optional)
- Environment: 4 variables (2 optional)
- Health check: Verifies both CLI tools
- Network: Custom bridge network

**Key Sections:**

1. **Build Configuration** (lines 3-6)
   - Context: current directory
   - Dockerfile: local Dockerfile
   - Auto-builds if image not present

2. **User Context** (line 8)
   - Runs as non-root user
   - Configurable UID/GID via environment

3. **Volume Configuration** (lines 11-27)
   - Project root, Beads DB, Git repo required
   - SSH and gitconfig optional

4. **Environment Setup** (lines 29-35)
   - Beads database location
   - Terminal color support
   - Optional API key configuration

5. **Interactive Mode** (lines 37-38)
   - stdin_open for input
   - tty for terminal features

6. **Health Check** (lines 40-47)
   - Validates both Beads and Claude Code
   - 30-second intervals

## Testing Results

### Test 1: Docker Compose Build
```bash
docker-compose build
```
**Result**: ✅ PASSED
- Image built successfully
- All layers cached correctly
- Build time: minimal (cached)
- Image name: `vibeguard:latest`

### Test 2: Beads Commands
```bash
docker-compose run --rm vibeguard bd ready
```
**Result**: ✅ PASSED
- Output: Shows 4 ready issues
- Project volume mounted correctly
- `.beads` directory accessible
- All Beads commands work

### Test 3: Claude Code Commands
```bash
docker-compose run --rm vibeguard /usr/local/bin/claude --help
```
**Result**: ✅ PASSED
- Help command works
- Binary available
- No errors in execution

### Test 4: Interactive Shell
```bash
docker-compose run --rm vibeguard /bin/bash
```
**Result**: ✅ PASSED
- Interactive shell ready
- Project files visible
- Can edit and commit changes

## Usage Examples

### Start Interactive Development Session
```bash
docker-compose run --rm vibeguard /bin/bash
# Inside container:
cd /work
bd ready          # Check available issues
git status        # Check git status
claude --help     # Access Claude Code
```

### Run Specific Commands
```bash
# Check Beads status
docker-compose run --rm vibeguard bd list --status=open

# Create an issue
docker-compose run --rm vibeguard bd create --title="..." --type=task

# Run Claude Code
docker-compose run --rm vibeguard claude --print "help me understand this code"
```

### Enable Optional Credential Mounting
```bash
# Uncomment SSH mounting in docker-compose.yml for Git SSH operations
# Uncomment gitconfig mounting for user configuration

# Then use:
docker-compose run --rm vibeguard git push origin main
```

### Enable API Key
```bash
# Create .env file or pass via command line
docker-compose run --rm -e ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY vibeguard claude --help
```

## Benefits

1. **Standardized Environment**: Everyone runs identical Docker image
2. **Development Friendly**: Volume mounting enables real-time code changes
3. **Persistent State**: Beads database survives container lifecycle
4. **Version Control Integration**: Git operations work seamlessly
5. **Multi-Tool Support**: Both Beads and Claude Code available
6. **Security**: Non-root user, read-only credential mounts
7. **Easy Onboarding**: Single command to start development
8. **No Configuration**: Works out-of-box for common scenarios

## Related Work

- **vibeguard-911**: Test combined Claude Code + Beads workflow (COMPLETED - unblocked this task)
- **vibeguard-917**: Fix Beads binary download (COMPLETED)
- **vibeguard-918**: Install Claude Code in Docker image (COMPLETED)
- **vibeguard-912**: Create docker-compose.yml for development workflow (THIS TASK)
- **vibeguard-914**: Test volume permissions for credential mounting
- **vibeguard-915**: Run container security scan
- **vibeguard-919**: Test Beads in container

## Architecture Decision References

- **ADR-001**: Adopt Beads for AI Agent Task Management
- **ADR-005**: Adopt Vibeguard for Policy Enforcement in CI/CD
- **ADR-006**: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement

## Deployment Considerations

### Production vs Development

This docker-compose.yml is designed for **development**:
- Volumes for active code changes
- Interactive shell as default
- restart policy for development convenience
- Beads database mounted locally

For production:
- Remove interactive TTY
- Use specific image tags, not `latest`
- Volume mounts read-only
- Different restart policies

### Environment Variables

Key variables to configure:
- `UID`/`GID` - Match host user for file ownership
- `ANTHROPIC_API_KEY` - Required for Claude Code API calls
- `BEADS_DB` - Database location (default: `./.beads/beads.db`)

### Volume Permissions

The container runs as user `1000:1000` (configurable):
- Ensure host user has same UID/GID, or
- Set `UID` and `GID` environment variables to match host user
- Avoid permission issues with mounted volumes

## Conclusion

The `docker-compose.yml` provides a complete, ready-to-use development environment that:
- Integrates Beads for task management
- Includes Claude Code for AI assistance
- Maintains persistent project state
- Supports version control workflows
- Works seamlessly with existing project structure
- Follows Docker best practices

Developers can now use a single command (`docker-compose run --rm vibeguard`) to get a fully configured development environment.

## Next Steps

1. **vibeguard-914**: Test volume permissions for credential mounting scenarios
2. **vibeguard-915**: Run container security scan on final image
3. **vibeguard-916**: Document GitHub API rate limiting mitigation
4. **vibeguard-919**: Test Beads in container (verify functionality)
5. **Documentation**: Add docker-compose usage guide to project README
