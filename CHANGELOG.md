# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `--version` flag to display semantic version (6c6fec6)
- Gremlins mutation testing integration per ADR-007 (b7c9c80)
- `fix` field support to CheckRecommendation structs for AI agent-assisted setup (af7d1f3)
- Comprehensive check output format improvements for AI agents (7acd0c7)
- `fix` field to Check struct for suggested fixes (2f8c638)
- File field support to check recommendations (af7d1f3)
- Tool detection from Makefile and CI configs in inspector (9652049)
- AI-assisted setup feature via `init --assist` command (1739a55)
- Predefined templates for init command (d66c680)
- Tooling discovery instructions for AI agent-assisted setup (8ebc66e)
- Project type detection for AI-assisted setup (91ceab6)
- Tool detection system for AI-assisted setup (3079243)
- Metadata extraction for AI-assisted setup (982c016)
- Recommendation engine for diverse project types (3954585)
- isort and pip-audit check templates for Python (6f0e6b)
- Go template syntax for suggestions with grok extraction (8a8da96)
- YAML line numbers to config validation errors (d28bd69)
- File/line context to grok and assert error messages (cd40283)

### Fixed
- Exit code handling: config errors (2), execution errors (3) distinction (5399279)
- Return error for unknown check instead of silent exit (d30a014)
- Error message actionability with context and visual pointers (960d8b6)
- Unknown check exit handling (d30a014)
- Help text suppression when checks fail (f5de2c2)
- Check ID format validation (3472549)
- Inspector depth calculation and performance optimization (9106a31)
- Co-located Go test directory detection (c03ff68)
- Python pyproject.toml confidence detection (24004b4)
- stderr and exit code 2 for Claude Code hook compatibility (4ef4346)
- Coverage in JSON output for cancelled status (0fa8142)

### Documentation
- Security model documentation for shell injection in interpolate.go (a37d962)
- Epic closure and issue triage session log (0e31724)
- Generated setup prompt with improved detection (25ede72)
- Architecture Decision Record (ADR-007) for mutation testing (6044725)
- AI-assisted setup feature comprehensive documentation (0f2d086)
- Complete specification (SPEC.md) (4ef4346)

### Testing
- Mutation testing with Gremlins achieving 98.70% efficacy (acbe741)
- Comprehensive unit tests for CLI package (165c2ad)
- Test coverage improvement from 89.0% to 90.5% (ba730e3)
- Comprehensive integration tests for Phase 4 Polish (be8fa81)
- Comprehensive integration tests for diverse project types (89172a5)
- Comprehensive unit tests for AI-assisted setup (65513b4)
- Increased grok package test coverage to 90%+ (ADR requirement)
- Main binary integration tests for error handling paths

## Key Features

### Core Functionality
- Single-binary deployment with zero runtime dependencies
- Fast and lightweight policy enforcement for CI/CD pipelines
- Declarative policy definitions in YAML format
- Multiple runner patterns for flexible policy evaluation
- LLM judge integration for nuanced policy assessment
- Cross-platform support (Linux, macOS, Windows)

### CLI Commands
- `vibeguard check [id]` - Run all or specific checks
- `vibeguard init [--assist]` - Initialize configuration with optional AI assistance
- `vibeguard list` - List configured checks with dependencies
- `vibeguard validate` - Validate configuration without execution

### Configuration Features
- Global variable interpolation with `{{.variable}}` syntax
- Grok pattern extraction for structured data parsing
- Assertion expressions for conditional validation
- Check dependencies for ordered execution
- Configurable severity levels (error, warning)
- Timeout support for long-running checks
- Go template syntax for dynamic suggestions

### AI-Assisted Setup
- Automatic project type detection (Go, Node.js, Python, Rust, Ruby, Java)
- Existing tool discovery and configuration analysis
- Context-aware check recommendations
- Generated setup guides for AI coding agents (Claude Code, Cursor, etc.)

## Development & Quality Standards

### Code Quality
- golangci-lint configuration for style enforcement
- goimports for automatic import formatting
- Pre-commit hooks for shift-left quality assurance
- 70% minimum test coverage requirement
- Mutation testing with Gremlins (ADR-007)

### Version Control
- Conventional Commits specification (ADR-002)
- Beads task management for AI agents (ADR-001)
- Architecture Decision Records (ADRs) for major decisions

### Architecture
- Go language for single-binary deployment (ADR-003)
- Declarative policy system with composable checks (ADR-005)
- Git pre-commit hooks for policy enforcement (ADR-006)

## Known Limitations

- Configuration file auto-discovery searches for `vibeguard.yaml`, `vibeguard.yml`, `.vibeguard.yaml`, `.vibeguard.yml` (in that order)
- Default parallel check limit is 4 (configurable via `--parallel` flag)
- Default check timeout is 30 seconds (configurable per-check)
- JSON output currently limited to 30000 characters (will truncate if exceeded)

## Getting Started

### Installation
```bash
git clone https://github.com/vibeguard/vibeguard.git
cd vibeguard
go build -o vibeguard ./cmd/vibeguard
```

### Basic Usage
```bash
# Initialize configuration
vibeguard init

# Run all checks
vibeguard check

# Run with AI assistance
vibeguard init --assist
```

### For Contributors
See [CONVENTIONS.md](./CONVENTIONS.md) for code style requirements and [CLAUDE.md](./CLAUDE.md) for AI agent guidelines.

---

**Note**: Early development phase. API and configuration format may change. See ADRs in `docs/adr/` for architectural details and rationale.
