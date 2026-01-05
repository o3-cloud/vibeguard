---
summary: Implemented support infrastructure including SUPPORT.md and GitHub Discussions templates
event_type: code
sources:
  - SUPPORT.md
  - .github/DISCUSSION_TEMPLATE/q_and_a.yml
  - .github/DISCUSSION_TEMPLATE/ideas.yml
  - .github/DISCUSSION_TEMPLATE/show_and_tell.yml
  - .github/DISCUSSION_TEMPLATE/general.yml
tags:
  - support
  - documentation
  - github-discussions
  - community
  - response-sla
---

# Support Infrastructure Implementation

Completed vibeguard-89o: Create support infrastructure (SUPPORT.md, GitHub Discussions).

## Files Created

### SUPPORT.md

Created a comprehensive support document that includes:

- **Support channels overview** — GitHub Issues for bugs, GitHub Discussions for Q&A
- **Response time SLAs** — Defined expectations for different priority levels:
  - Security issues: 24-hour initial response, 7-day resolution target
  - P0/P1 bugs: 48-hour response, 14-day resolution target
  - P2/P3 bugs: 5 business day response, best effort resolution
  - Feature requests: 7 business day response, roadmap review
  - Discussions/Q&A: 7 business day response, community-driven
- **Priority definitions** — P0 (critical) through P3 (low) with clear criteria
- **Guidance for asking good questions** — Template and checklist
- **Security issue reporting** — Reference to SECURITY.md

### GitHub Discussion Templates

Created four discussion category templates in `.github/DISCUSSION_TEMPLATE/`:

1. **q_and_a.yml** — For questions with context, configuration, and version fields
2. **ideas.yml** — For feature ideas with problem statement and proposed solution
3. **show_and_tell.yml** — For sharing configurations and use cases
4. **general.yml** — For open-ended discussions

### README.md Update

Updated the Support section to reference:
- GitHub Discussions for questions
- GitHub Issues for bug reports
- SECURITY.md for security issues
- SUPPORT.md for detailed support information

## Design Decisions

1. **SLA targets are realistic for volunteer-maintained project** — Chose reasonable timeframes that can be met consistently rather than aggressive targets that might be missed
2. **Structured discussion templates** — Using GitHub's YAML-based discussion templates for better organization and required fields
3. **Priority-based response times** — Different SLAs for different priority levels recognizes that not all issues are equally urgent
4. **Pre-submission checklist** — All templates include checklist items to encourage users to search existing resources first

## Verification

Ran `vibeguard check -v` — all 8 checks passed:
- vet, fmt, actionlint, lint, test, test-coverage, build, mutation

## Next Steps

To enable GitHub Discussions on the repository:
1. Go to repository Settings > General > Features
2. Enable "Discussions"
3. The discussion templates will automatically be available

The templates are ready and will work as soon as Discussions is enabled in repository settings.
