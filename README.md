# VibeGuard

**VibeGuard** is a lightweight, composable policy enforcement system designed for seamless integration with CI/CD pipelines, agent loops, and Cloud Code workflows.

## Overview

VibeGuard enforces policies at scale with minimal overhead. It combines declarative policy definition, flexible runner patterns, and LLM judge integration to provide intelligent policy evaluation for modern development workflows.

### Key Features

- **Single Binary Deployment** — Compiled Go binary with zero runtime dependencies
- **Fast & Lightweight** — Minimal resource usage and quick startup time for frequent invocation
- **Declarative Policies** — Define policies in YAML for clarity and maintainability
- **Multiple Runner Patterns** — Supports various policy evaluation approaches (see patterns below)
- **Judge Integration** — Leverage LLMs for nuanced policy evaluation
- **Cross-Platform** — Runs seamlessly on Linux, macOS, and Windows

## Quick Start

### Building

```bash
# Build the binary
go build -o vibeguard ./cmd/vibeguard

# Run the CLI
./vibeguard --help
```

### Running Tests

```bash
go test -v ./...
```

## Implementation Patterns

VibeGuard supports five core implementation patterns:

1. **Declarative Policy Runner** — Evaluate policies defined in structured formats (YAML/JSON)
2. **Judge-as-a-Policy** — Use LLMs as policy evaluators
3. **Cloud Code Native** — Native integration with Anthropic Cloud Code
4. **Event-Based Policy Graph** — React to events with graph-based policy evaluation
5. **Git-Aware Guardrails** — Policies based on git history and code changes

(See `docs/patterns/` for detailed documentation on each pattern)

## Project Structure

```
vibeguard/
├── cmd/
│   └── vibeguard/              # Main CLI application
├── internal/
│   ├── cli/                    # Command-line interface (Cobra-based)
│   ├── config/                 # Configuration loading
│   ├── judge/                  # Judge API integration
│   ├── policy/                 # Policy evaluation engine
│   └── runner/                 # Runner implementations
├── pkg/
│   └── models/                 # Shared data models
├── docs/
│   ├── adr/                    # Architecture Decision Records
│   └── patterns/               # Implementation patterns
└── README.md, CONVENTIONS.md   # Documentation
```

## Development

### Prerequisites

- **Go 1.21+** — Latest stable Go version
- **git** — For version control

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/vibeguard/vibeguard.git
cd vibeguard

# Install dependencies
go mod tidy

# Run tests
go test -v ./...

# Run linting
go fmt ./...
go vet ./...
```

### Code Style and Conventions

See `CONVENTIONS.md` for detailed code style guidelines, naming conventions, and development standards.

### Making Changes

1. Create a feature branch from `main`
2. Make your changes following the conventions in `CONVENTIONS.md`
3. Write or update tests as needed
4. Commit using Conventional Commits format (see ADR-002)
5. Push and create a pull request

## Architecture Decisions

Major architectural decisions are documented as Architecture Decision Records (ADRs) in `docs/adr/`:

- **ADR-001** — Adopt Beads for AI Agent Task Management
- **ADR-002** — Adopt Conventional Commits
- **ADR-003** — Adopt Go as the Primary Implementation Language

Review these documents to understand the project's design rationale and constraints.

## Contributing

Contributions are welcome! Please:

1. Read `CONVENTIONS.md` for code style requirements
2. Review existing ADRs to understand architectural constraints
3. Write tests for new functionality
4. Follow Conventional Commits for commit messages
5. Keep PRs focused and well-documented

## License

VibeGuard is released under [LICENSE_NAME] (see LICENSE file for details).

## Further Reading

- [Conventional Commits](https://www.conventionalcommits.org/) — Commit message specification
- [Go Best Practices](https://go.dev/doc/effective_go) — Go coding guidelines
- [Project Conventions](./CONVENTIONS.md) — VibeGuard-specific standards

## Support

For issues, questions, or suggestions:

- Open an issue on GitHub
- Check existing ADRs and documentation
- Review the implementation patterns in `docs/patterns/`
