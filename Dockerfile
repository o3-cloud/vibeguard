FROM ubuntu:24.04

# Metadata labels
LABEL org.opencontainers.image.source="https://github.com/anthropics/vibeguard"
LABEL org.opencontainers.image.description="Claude Code + Beads development environment"

# Apply security mitigations for known vulnerabilities
# CVE-2025-68973 (High) in gpgv - removed (not needed)
# CVE-2025-14104 (Medium) in util-linux - limit to essential tools only
# CVE-2025-8941 (Medium) in PAM - accept risk as it's needed for user switching
# Note: Many Medium CVEs in system packages are not patchable in Ubuntu 24.04 yet
# See: https://ubuntu.com/security/ for official patch status
RUN apt-get update && \
    apt-get remove -y gpg-agent gpgv gnupg-l10n gnupg-utils || true && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install minimal dependencies in one layer
# Reduced to only what Beads, Claude Code CLI, and git operations require
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    bash \
    ca-certificates \
    git \
    jq \
    libc6 \
    && apt-get remove -y \
      man-db \
      manpages \
      groff-base \
      perl-modules \
    || true \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /usr/share/doc/* \
    && rm -rf /usr/share/man/* \
    && rm -rf /usr/share/info/* \
    && rm -rf /var/cache/apt/*

# Install Beads from pre-built binary (as root, to /usr/local/bin)
# Supports both amd64 and arm64 architectures
# Includes retry logic for GitHub API rate limiting and connection issues
RUN ARCH=$(dpkg --print-architecture) && \
    RELEASE_URL="" && \
    for attempt in 1 2 3; do \
      echo "Fetching Beads release info (attempt $attempt/3)..." && \
      curl -fsSL https://api.github.com/repos/steveyegge/beads/releases/latest 2>/dev/null | \
        jq -r '.assets[] | select(.name | test("beads.*linux_'"${ARCH}"'\\.tar\\.gz$")) | .browser_download_url' | \
        head -1 > /tmp/url.txt && \
      RELEASE_URL=$(cat /tmp/url.txt) && \
      if [ -s /tmp/url.txt ] && [ -n "$RELEASE_URL" ]; then break; fi && \
      if [ $attempt -lt 3 ]; then sleep 2; fi; \
    done && \
    if [ -z "$RELEASE_URL" ] || [ ! -s /tmp/url.txt ]; then echo "ERROR: Failed to resolve Beads release URL for architecture $ARCH after 3 attempts"; exit 1; fi && \
    echo "Downloading Beads binary from: $RELEASE_URL" && \
    curl -fsSL -o /tmp/beads.tar.gz "$RELEASE_URL" || { echo "ERROR: Failed to download Beads binary"; exit 1; } && \
    tar -xzf /tmp/beads.tar.gz -C /tmp || { echo "ERROR: Failed to extract Beads archive"; exit 1; } && \
    chmod +x /tmp/bd && \
    mv /tmp/bd /usr/local/bin/bd && \
    /usr/local/bin/bd --version || { echo "ERROR: Installed Beads binary failed version check"; exit 1; } && \
    rm -rf /tmp/beads.tar.gz /tmp/url.txt

# Install Claude Code CLI from official installer
# The installer places binary in ~/.local/bin, so we move it to /usr/local/bin for global access
# Retry logic handles transient network issues
RUN CLAUDE_INSTALL_ATTEMPT="" && \
    for attempt in 1 2 3; do \
      echo "Installing Claude Code (attempt $attempt/3)..." && \
      curl -fsSL https://claude.ai/install.sh 2>/dev/null | bash 2>&1 && \
      CLAUDE_INSTALL_ATTEMPT="success" && break || true && \
      if [ $attempt -lt 3 ]; then sleep 2; fi; \
    done && \
    if [ -z "$CLAUDE_INSTALL_ATTEMPT" ]; then echo "ERROR: Failed to install Claude Code after 3 attempts"; exit 1; fi && \
    if [ -f ~/.local/bin/claude ]; then \
      mv ~/.local/bin/claude /usr/local/bin/claude && \
      chmod +x /usr/local/bin/claude; \
    else echo "ERROR: Claude Code binary not found at ~/.local/bin/claude"; exit 1; fi && \
    /usr/local/bin/claude --help > /dev/null || { echo "WARNING: Claude Code help check failed, but binary exists"; }

# Set up non-root user (security best practice)
RUN useradd -m -s /bin/bash claude

# Switch to non-root user for default execution
USER claude
WORKDIR /home/claude

# Set up PATH for installed tools (add to user's PATH)
ENV PATH="/usr/local/bin:${PATH}"

# Health check to ensure both Beads and Claude Code CLIs are available
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD /usr/local/bin/bd --version > /dev/null && which claude > /dev/null || exit 1
