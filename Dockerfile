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
RUN ARCH=$(dpkg --print-architecture) && \
    RELEASE_URL=$(curl -fsSL https://api.github.com/repos/steveyegge/beads/releases/latest | jq -r '.assets[] | select(.name | endswith("linux_'"${ARCH}"'.tar.gz")) | .browser_download_url' | head -1) && \
    echo "Downloading from: $RELEASE_URL" && \
    curl -fsSL -o /tmp/beads.tar.gz "$RELEASE_URL" && \
    tar -xzf /tmp/beads.tar.gz -C /tmp && \
    chmod +x /tmp/bd && \
    mv /tmp/bd /usr/local/bin/bd && \
    /usr/local/bin/bd --version && \
    rm -rf /tmp/beads.tar.gz

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
