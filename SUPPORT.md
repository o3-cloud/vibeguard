# Getting Support for VibeGuard

This document describes how to get help with VibeGuard, the available support channels, and expected response times.

## Support Channels

### GitHub Issues

For bugs, feature requests, and technical problems, open an issue on GitHub:

- [Bug Reports](https://github.com/vibeguard/vibeguard/issues/new?template=bug_report.md) — Report unexpected behavior or crashes
- [Feature Requests](https://github.com/vibeguard/vibeguard/issues/new?template=feature_request.md) — Suggest new features or improvements

**Best for:**
- Reproducible bugs with clear steps
- Feature suggestions with use cases
- Documentation errors or gaps

### GitHub Discussions

For questions, ideas, and community conversations, use GitHub Discussions:

- [Q&A](https://github.com/vibeguard/vibeguard/discussions/categories/q-a) — Ask questions and get help from the community
- [Ideas](https://github.com/vibeguard/vibeguard/discussions/categories/ideas) — Share ideas for new features or improvements
- [Show and Tell](https://github.com/vibeguard/vibeguard/discussions/categories/show-and-tell) — Share your VibeGuard configurations and use cases
- [General](https://github.com/vibeguard/vibeguard/discussions/categories/general) — General discussions about VibeGuard

**Best for:**
- Questions about usage or best practices
- Discussing implementation approaches
- Sharing configurations and patterns
- Connecting with other users

### Documentation

Before opening an issue or discussion, check the existing documentation:

- [README.md](README.md) — Quick start and CLI reference
- [CONTRIBUTING.md](CONTRIBUTING.md) — Development workflow and code standards
- [SECURITY.md](SECURITY.md) — Security model and vulnerability reporting
- [Architecture Decision Records](docs/adr/) — Design rationale and constraints
- [AI-Assisted Setup Guide](docs/ai-assisted-setup.md) — Using `vibeguard init --assist`

## Response Time SLAs

We aim to respond to community requests within the following timeframes:

| Channel | Priority | Initial Response | Resolution Target |
|---------|----------|------------------|-------------------|
| **Security Issues** | Critical | 24 hours | 7 days |
| **Bug Reports (P0/P1)** | High | 48 hours | 14 days |
| **Bug Reports (P2/P3)** | Medium | 5 business days | Best effort |
| **Feature Requests** | Normal | 7 business days | Roadmap review |
| **Discussions/Q&A** | Normal | 7 business days | Community-driven |

### Priority Definitions

- **P0 (Critical):** Complete system failure, data loss, security vulnerability
- **P1 (High):** Major feature broken, no workaround available
- **P2 (Medium):** Feature partially broken, workaround available
- **P3 (Low):** Minor issue, cosmetic problem, documentation improvement

### Notes on Response Times

- Response times are targets, not guarantees. VibeGuard is maintained by volunteers.
- "Initial response" means acknowledgment and triage, not resolution.
- Complex issues may take longer to resolve regardless of priority.
- Community contributions to fixes are always welcome.

## Before Asking for Help

To help us help you faster, please:

1. **Check existing resources:**
   - Search [existing issues](https://github.com/vibeguard/vibeguard/issues) and [discussions](https://github.com/vibeguard/vibeguard/discussions)
   - Review the [documentation](#documentation)
   - Check the [CHANGELOG](CHANGELOG.md) for recent changes

2. **Gather information:**
   - VibeGuard version: `vibeguard --version`
   - Operating system and version
   - Go version (if building from source)
   - Relevant configuration (sanitize sensitive data)
   - Steps to reproduce the issue

3. **Create a minimal example:**
   - Simplify your configuration to isolate the problem
   - Include only what's necessary to reproduce the issue

## Asking Good Questions

Good questions get faster answers. Include:

- **What you're trying to do** — The goal, not just the error
- **What you've tried** — Steps already taken
- **What happened** — Actual behavior vs. expected behavior
- **Environment details** — Version, OS, configuration

Example:

```markdown
## Summary
Coverage check fails with assertion error even though coverage exceeds threshold.

## Environment
- VibeGuard version: 0.1.0
- OS: Ubuntu 22.04
- Go version: 1.21

## Configuration (vibeguard.yaml)
```yaml
version: "1"
checks:
  - id: coverage
    run: go test ./... -coverprofile=cover.out && go tool cover -func=cover.out
    grok:
      - 'total:.*\(statements\)\s+%{NUMBER:coverage}%'
    assert: "coverage >= 80"
```

## Steps to Reproduce
1. Run `vibeguard check coverage`
2. Observe assertion failure

## Expected Behavior
Check should pass (actual coverage is 85%)

## Actual Behavior
Assertion fails with "coverage (85.50) < 80"

## Additional Context
The grok pattern matches correctly (verified with verbose output).
```

## Security Issues

For security vulnerabilities, **do not** open a public issue. Instead:

1. Review our [Security Policy](SECURITY.md)
2. Email findings privately to maintainers
3. Include steps to reproduce and potential impact

See [SECURITY.md](SECURITY.md) for our complete security policy and responsible disclosure process.

## Contributing

If you find a bug and know how to fix it, consider contributing:

1. Read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines
2. Fork the repository and create a branch
3. Submit a pull request with your fix

We welcome contributions from everyone, regardless of experience level.

## Commercial Support

VibeGuard is an open-source project without commercial support options at this time. All support is provided by the community on a best-effort basis.

## Stay Updated

- **Watch** the repository for release notifications
- **Star** the repository to show support
- Check the [CHANGELOG](CHANGELOG.md) for version history

## Code of Conduct

Please be respectful and inclusive in all interactions. We're committed to providing a welcoming and supportive environment for all community members.
