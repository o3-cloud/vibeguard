FROM ubuntu:24.04

# Metadata labels
LABEL org.opencontainers.image.source="https://github.com/anthropics/vibeguard"
LABEL org.opencontainers.image.description="Claude Code + Beads development environment"

# Install dependencies in one layer
# Include libc6 explicitly to ensure runtime libraries are present
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    bash \
    ca-certificates \
    git \
    jq \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

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

# Set up non-root user (security best practice)
RUN useradd -m -s /bin/bash claude

# Switch to non-root user for default execution
USER claude
WORKDIR /home/claude

# Set up PATH for installed tools (add to user's PATH)
ENV PATH="/usr/local/bin:${PATH}"

# Health check to ensure Beads CLI is available
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD /usr/local/bin/bd --version > /dev/null || exit 1
