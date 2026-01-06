---
summary: Successfully installed Claude Code CLI in Docker image with proper binary relocation and retry logic for transient installation failures
event_type: code
sources:
  - https://github.com/anthropics/claude-code
  - https://code.claude.com/docs/en/setup
  - Dockerfile implementation
tags:
  - docker
  - claude-code
  - dockerfile
  - cli-tools
  - installation
  - container-image
  - development-environment
---

# Install Claude Code in Docker Image

## Overview

Successfully added Claude Code CLI installation to the vibeguard Dockerfile. The Docker image now includes both Beads (task management) and Claude Code (AI-assisted development), creating a complete development environment for containerized AI-powered workflows.

## Problem Statement

The Dockerfile was labeled as providing a "Claude Code + Beads development environment," but Claude Code was not actually installed. This issue blocked:
- Testing credential mounting with Claude Code
- Development workflows that require both tools
- Integration testing of the full development environment

## Solution Implemented

### Installation Method

Used the official Claude Code installer script from Anthropic:
```bash
curl -fsSL https://claude.ai/install.sh | bash
```

**Key Characteristics:**
- Official, maintained installation method from Anthropic
- Works on Ubuntu 20.04+ without additional dependencies
- Downloads pre-built binary (no compilation required)
- Automatically handles system setup

### Binary Relocation Strategy

The official installer places the binary at `~/.local/bin/claude` (user-specific location). For a system-wide Docker image, we:

1. Run installer as root (which installs to root's home)
2. Move binary from `~/.local/bin/claude` to `/usr/local/bin/claude`
3. Set proper executable permissions

This approach:
- Makes Claude Code available globally in the container
- Works for both root and non-root users
- Follows Docker best practices for system-wide tools

### Retry Logic

Implemented 3-attempt retry loop with 2-second delays:
```bash
for attempt in 1 2 3; do
  echo "Installing Claude Code (attempt $attempt/3)..."
  curl -fsSL https://claude.ai/install.sh | bash && success
  if retry_available; then sleep 2; fi
done
```

**Handles:**
- Network timeouts during installer download
- Temporary connectivity issues
- GitHub/CDN rate limiting

### Error Handling

**Validation Steps:**
1. Track installation attempt success/failure
2. Verify binary exists at expected location after installation
3. Test binary with `--help` command
4. Provide clear error messages if any step fails

**Implementation Details:**
```bash
if [ -z "$CLAUDE_INSTALL_ATTEMPT" ]; then
  echo "ERROR: Failed to install Claude Code after 3 attempts"
  exit 1
fi
if [ -f ~/.local/bin/claude ]; then
  mv ~/.local/bin/claude /usr/local/bin/claude
else
  echo "ERROR: Claude Code binary not found"
  exit 1
fi
/usr/local/bin/claude --help > /dev/null
```

## Implementation Details

### Changes to Dockerfile (lines 41-56)

**Installation Block:**
- Lines 41-43: Comments explaining the installation method and retry logic
- Lines 44-50: Retry loop for fetching and executing installer
- Lines 51-55: Binary relocation from user to system location
- Line 56: Validation with help command

**Health Check Update (lines 65-67):**
- Added Claude Code verification to Docker health check
- Both Beads and Claude Code must be available for healthy container

### Comments in Dockerfile

Added clear documentation explaining:
- Official installer source
- Reason for binary relocation
- Retry logic for transient issues

## Installation Details

### Installed Version
- **Version**: 2.0.76 (Claude Code)
- **Installation Method**: Official curl-based installer
- **Installation Location**: `/usr/local/bin/claude` (system-wide)
- **Dependency**: curl, bash, network connectivity

### System Requirements

The installation requires:
- Ubuntu 20.04+ (or equivalent Debian-based system)
- curl (already installed for Beads)
- bash (already installed)
- Network access to claude.ai install script and CDNs

## Testing Results

### Test 1: Beads Installation Verification
```bash
docker run --rm vibeguard:claude-final /usr/local/bin/bd --version
```
**Result**: ✅ PASSED
- Output: `bd version 0.44.0 (d7221f68)`

### Test 2: Claude Code Installation Verification
```bash
docker run --rm vibeguard:claude-final /usr/local/bin/claude --help
```
**Result**: ✅ PASSED
- Output: Complete help documentation
- Binary location: `/usr/local/bin/claude`
- Version: 2.0.76

### Test 3: Docker Build Success
```bash
docker build --no-cache -t vibeguard:claude-final .
```
**Result**: ✅ PASSED
- All build stages completed successfully
- Both installers (Beads and Claude Code) completed without errors
- Image built successfully: `vibeguard:claude-final`

### Test 4: Policy Compliance
```bash
vibeguard check
```
**Result**: ✅ PASSED
- No policy violations detected

### Test 5: Health Check
The Docker health check verifies:
- `/usr/local/bin/bd --version` - Beads CLI available
- `which claude` - Claude Code CLI in PATH
- Both must pass for container to be healthy

## Benefits

1. **Complete Development Environment**: Docker image now includes both Beads (task management) and Claude Code (AI development)
2. **Reliability**: Retry logic handles transient network issues during installation
3. **System-Wide Availability**: Claude Code accessible from any user in container
4. **Debuggability**: Clear error messages indicate failure location
5. **Maintainability**: Documented installation method with comments
6. **Production Ready**: Proper error handling prevents silent failures

## Related Work

- **vibeguard-917**: Fix Beads binary download (COMPLETED - blocking resolved)
- **vibeguard-918**: Install Claude Code in Docker image (THIS TASK)
- **vibeguard-911**: Test combined Claude Code + Beads workflow (Now enabled)
- **vibeguard-914**: Test volume permissions for credential mounting
- **vibeguard-919**: Test Beads in container

## Architecture Decision References

- **ADR-001**: Adopt Beads for AI Agent Task Management
- **ADR-005**: Adopt Vibeguard for Policy Enforcement in CI/CD
- **ADR-006**: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement

## Changelog

### Changes Made
1. Added Claude Code installation to Dockerfile (after Beads, before non-root user setup)
2. Implemented 3-attempt retry logic for installation
3. Added binary relocation from user directory to system directory
4. Updated health check to verify both Beads and Claude Code
5. Added proper documentation and comments

### No Breaking Changes
- Existing Beads installation unchanged
- Non-root user setup unchanged
- IMAGE labels and metadata unchanged
- All other functionality preserved

## Conclusion

Claude Code is now successfully installed and integrated into the vibeguard Docker image. The implementation:
- Follows official installation practices from Anthropic
- Includes robust error handling and retry logic
- Maintains system-wide accessibility for all users
- Works seamlessly with the existing Beads installation
- Provides a complete development environment for AI-powered workflows

The Docker image is now ready for:
- Credential mounting tests with Claude Code
- Combined Claude Code + Beads workflow testing
- Development workflows requiring AI assistance
- Full integration testing of the development environment

## Next Steps

1. **vibeguard-911**: Test combined Claude Code + Beads workflow in container
2. **vibeguard-914**: Test volume permissions for credential mounting scenarios
3. **vibeguard-912**: Create docker-compose.yml for development workflow
4. **vibeguard-915**: Run container security scan
5. **vibeguard-916**: Document GitHub API rate limiting mitigation
