# ADR-005: Adopt Vibeguard for Policy Enforcement in CI/CD

## Context and Problem Statement

The VibeGuard project requires a reliable, maintainable way to enforce code quality and policy checks in CI/CD pipelines. Historically, projects rely on ad-hoc shell scripts, GitHub Actions workflows, or custom CI/CD logic to enforce checks like:

- Code formatting and linting
- Testing requirements
- Build validation
- Deployment readiness

These approaches suffer from:
- **Fragmentation**: Different projects use different check mechanisms, making them inconsistent
- **Maintenance Burden**: Shell scripts and CI-specific syntax are difficult to version control, understand, and update
- **Limited Composability**: Hard to create reusable, shareable policy definitions
- **Unclear Semantics**: Implicit dependencies between checks, unclear failure modes, and difficulty in prioritizing failures

We need a declarative, composable policy enforcement system that can:
1. Define policies clearly in a standard format
2. Run checks with transparent dependencies and timeouts
3. Provide actionable feedback for remediation
4. Integrate seamlessly into CI/CD, agent loops, and development workflows
5. Support flexible evaluation (e.g., LLM-based policy judges)

## Considered Options

### Option A: Continue with Ad-Hoc Shell Scripts and GitHub Actions
- Use `make` targets, bash scripts, and GitHub Actions workflows
- **Pros**: Simple to start, no new dependencies, immediate familiarity
- **Cons**: Fragmented across projects, difficult to maintain, no standard format, inconsistent error handling

### Option B: Adopt a Third-Party CI/CD Policy Tool (e.g., OPA/Rego)
- Use existing open-source policy engines like OPA (Open Policy Agent)
- **Pros**: Battle-tested, large community, mature ecosystem
- **Cons**: Designed for infrastructure policies, not application-level checks; overkill for simple use cases; additional language to learn

### Option C: Adopt VibeGuard for Policy Enforcement
- Use VibeGuard as the unified policy enforcement system across the project
- **Pros**: Lightweight Go binary, single-binary deployment, YAML-based policy definition, designed for CI/CD integration, supports flexible runners including LLM judges, aligns with Go-first architecture decision
- **Cons**: Relatively new tool (requires community validation), adds responsibility to maintain as project matures

## Decision Outcome

**Chosen option:** Adopt VibeGuard for Policy Enforcement

### Rationale

VibeGuard is the optimal choice because:

1. **Alignment with ADR-003**: Go binary aligns with the project's decision to use Go for single-binary deployment and minimal overhead
2. **Designed for This Purpose**: VibeGuard was built specifically for lightweight CI/CD policy enforcement, avoiding unnecessary complexity
3. **Declarative and Maintainable**: YAML-based policy definition is easy to version control, review, and understand
4. **Composable Patterns**: Supports multiple runner patterns (declarative, judge-based, event-driven, git-aware) enabling evolution without breaking changes
5. **Low Barrier to Entry**: No new languages to learn, integrates cleanly with standard CI/CD tools and agent loops
6. **Self-Validating**: Using VibeGuard in its own CI/CD pipeline provides immediate real-world validation and drives product improvement
7. **Future Extensibility**: LLM judge integration enables intelligent policy evaluation as workflows evolve

## Consequences

### Positive Outcomes
- **Single Source of Truth**: All policies defined in `vibeguard.yaml`, reducing cognitive load and maintenance burden
- **Transparent Dependencies**: Policy execution order and dependencies are explicit, preventing silent failures
- **Better Error Messages**: Structured suggestions guide developers to remediation steps
- **Framework for Innovation**: Judge integration enables AI-driven policy evaluation in future iterations
- **Dogfooding**: Using VibeGuard in its own CI/CD pipeline validates the tool and drives improvements
- **Consistency**: All projects can adopt the same policy enforcement approach

### Negative Outcomes
- **New Tool Adoption**: Teams must learn VibeGuard concepts and configuration format
- **Maintenance Responsibility**: As an in-house tool, the project bears maintenance and evolution responsibility
- **Community Maturity**: Unlike established tools, VibeGuard must prove stability and community viability over time

### Neutral Impacts
- **Configuration Overhead**: Initial setup of `vibeguard.yaml` requires understanding the policy format
- **Check Timeout Management**: Must balance timeout values with CI/CD SLA requirements

## Implementation Steps

1. Integrate `vibeguard check` into the primary CI/CD pipeline
2. Define comprehensive policies in `vibeguard.yaml` covering:
   - Go code formatting (`fmt`)
   - Static analysis (`vet`, `lint`)
   - Testing (`test`)
   - Build validation (`build`)
3. Document policy enforcement in `CONVENTIONS.md`
4. Add `vibeguard check` to pre-commit hooks for local validation
5. Monitor policy compliance and iterate on policies based on team feedback
