# VibeGuard Governance

This document describes the governance model for the VibeGuard project.

## Overview

VibeGuard is an open source project that welcomes contributions from the community. This document outlines how decisions are made, how contributors can become maintainers, and how the project is managed.

## Roles and Responsibilities

### Users

Users are community members who use VibeGuard. Anyone can be a user; there are no special requirements. Users are encouraged to:

- Use VibeGuard in their projects
- Report bugs and request features through GitHub Issues
- Participate in discussions on GitHub Discussions
- Help other users in community channels

### Contributors

Contributors are community members who contribute to VibeGuard. Contributions can include code, documentation, bug reports, feature requests, and community support. Anyone can become a contributor by:

- Submitting pull requests
- Filing or commenting on issues
- Improving documentation
- Helping users in discussions
- Reviewing pull requests

Contributors are expected to follow the [Code of Conduct](CODE_OF_CONDUCT.md) and [contribution guidelines](CONTRIBUTING.md).

### Maintainers

Maintainers are contributors who have demonstrated a sustained commitment to the project. Maintainers have write access to the repository and are responsible for:

- Reviewing and merging pull requests
- Triaging issues and discussions
- Making decisions about project direction
- Ensuring code quality and consistency
- Mentoring new contributors
- Representing the project in the community

The current maintainers are listed in [MAINTAINERS.md](MAINTAINERS.md).

#### Becoming a Maintainer

Contributors may be invited to become maintainers after demonstrating:

- A track record of high-quality contributions
- Understanding of the project's goals and architecture
- Ability to review code and provide constructive feedback
- Commitment to the project's long-term success
- Adherence to the Code of Conduct

New maintainers are nominated by existing maintainers and approved through consensus.

#### Maintainer Responsibilities

Maintainers are expected to:

- Respond to issues and pull requests in a timely manner
- Review code thoroughly and provide constructive feedback
- Participate in project planning and decision-making
- Help onboard new contributors
- Follow the Code of Conduct in all interactions

Maintainers who become inactive may be removed from the maintainer list after consultation.

## Decision Making

### Lazy Consensus

Most decisions are made through "lazy consensus." This means that:

1. A proposal is made (via issue, PR, or discussion)
2. The proposal is discussed and refined
3. If no objections are raised within a reasonable time, the proposal is accepted
4. Silence is interpreted as agreement

For routine changes (bug fixes, documentation updates, minor improvements), lazy consensus with a single maintainer approval is sufficient.

### Voting

For significant decisions that cannot reach lazy consensus, maintainers may call for a vote:

- Each maintainer has one vote
- A simple majority (>50%) is required for approval
- The vote must be open for at least 72 hours
- The vote and result should be documented in the relevant issue or discussion

Significant decisions include:

- Major architectural changes
- Breaking changes to public APIs
- Changes to governance policies
- Adding or removing maintainers
- Changing the license

### Technical Decisions

Technical decisions follow these principles:

1. **Prefer simplicity** — Simple solutions are preferred over complex ones
2. **Follow conventions** — Adhere to established project patterns and Go idioms
3. **Document decisions** — Major decisions are documented as Architecture Decision Records (ADRs) in `docs/adr/`
4. **Consider users** — Prioritize user experience and backwards compatibility
5. **Test thoroughly** — All changes should be tested appropriately

### Conflict Resolution

If contributors cannot reach agreement through discussion:

1. Escalate to maintainers for mediation
2. If maintainers cannot reach consensus, call for a vote
3. The project lead (if designated) may make a final decision in deadlock situations

## Releases

### Release Process

1. Maintainers decide when to cut a release based on accumulated changes
2. Release candidates may be published for testing
3. Releases follow [Semantic Versioning](https://semver.org/)
4. Release notes are generated from commit history (Conventional Commits)
5. Releases are tagged and published to GitHub Releases

### Version Support

- The latest major version receives full support (bugs and features)
- Previous major versions receive security fixes only
- Support periods are announced with major releases

## Communication

### Official Channels

- **GitHub Issues** — Bug reports and feature requests
- **GitHub Discussions** — Questions, ideas, and general discussion
- **GitHub Pull Requests** — Code contributions and reviews

### Response Times

Maintainers aim to:

- Acknowledge new issues within 1 week
- Provide initial review of PRs within 2 weeks
- These are goals, not guarantees; maintainers are volunteers

## Changes to Governance

Changes to this governance document require:

1. A pull request with the proposed changes
2. Discussion period of at least 2 weeks
3. Approval by a majority of maintainers

## License

VibeGuard is released under the Apache License 2.0. All contributions must be compatible with this license.

## Attribution

This governance model is inspired by open source projects including:

- [Node.js](https://github.com/nodejs/node/blob/main/GOVERNANCE.md)
- [Rust](https://www.rust-lang.org/governance)
- [Apache Software Foundation](https://www.apache.org/foundation/governance/)
