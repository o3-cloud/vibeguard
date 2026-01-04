# ADR-008: Adopt actionlint for GitHub Actions Workflow Validation

## Context and Problem Statement

GitHub Actions workflows are critical infrastructure code that orchestrates CI/CD pipelines. However, YAML syntax errors, deprecated action versions, and unused input parameters in workflows are often only discovered at runtime when workflows fail or misbehave.

As VibeGuard promotes policy-driven development and code quality standards (ADR-004), it's inconsistent to enforce quality on Go code but neglect the quality of workflow definitions that run the CI/CD system itself.

The project needs a mechanism to catch workflow issues early in the development process, before they cause failures in production deployments.

## Considered Options

### Option A: Do Nothing
- Rely on GitHub's built-in workflow validation and runtime feedback
- **Characteristics:** No additional tooling, minimal maintenance overhead
- **Trade-off:** Issues only discovered when workflows fail; slower feedback loop during development

### Option B: Use actionlint
- Adopt [actionlint](https://github.com/rhysd/actionlint), a specialized linter for GitHub Actions workflows
- **Characteristics:** Detects deprecated actions, syntax errors, unused inputs, shell injection risks, and more
- **Trade-off:** Adds another dependency and tool to the CI/CD pipeline; requires configuration

### Option C: Use Generic YAML Linting (yamllint)
- Use yamllint for general YAML schema validation
- **Characteristics:** Works for any YAML file, simpler setup
- **Trade-off:** Less specific to GitHub Actions; misses action-specific issues like deprecated versions or input validation

## Decision Outcome

**Chosen option:** Option B - Adopt actionlint

**Rationale:**

1. **Specialized for GitHub Actions:** actionlint is purpose-built for GitHub Actions workflows and catches issues generic YAML linters miss (deprecated actions, unused inputs, dangerous shell interpolation).

2. **Aligns with Code Quality Standards:** ADR-004 establishes comprehensive code quality standards. Workflows are code, and extending quality checks to workflows is consistent with the project's philosophy.

3. **Dogfooding VibeGuard:** As VibeGuard is a policy enforcement tool, using it to validate workflows (through vibeguard.yml configuration) demonstrates its effectiveness and aligns with ADR-005.

4. **Early Feedback:** Integration into CI ensures workflow issues are caught before they cause production failures, reducing debugging time and improving developer experience.

5. **Low Friction:** actionlint is lightweight, fast, and requires minimal configuration to get started.

**Tradeoffs:**

- Gained: Early detection of workflow quality issues; consistency in code quality enforcement
- Sacrificed: No additional tool dependencies are added that wouldn't otherwise be needed; the integration is straightforward

## Consequences

### Positive Outcomes
- Workflow syntax and logic errors are caught in CI before merging
- Deprecated GitHub Actions are identified, prompting timely updates
- Reduces debugging time when workflows misbehave
- Demonstrates VibeGuard's capability to enforce policies across infrastructure code
- Sets precedent for quality checks beyond application code

### Negative Outcomes
- One additional linting tool to maintain and monitor for updates
- Potential for false positives or overly strict rules requiring configuration tuning
- Team members must understand actionlint error messages and fixes

### Neutral Impacts
- actionlint runs alongside existing Go linting (golangci-lint) without conflicts
- Configuration is simple and can be stored in vibeguard.yml, centralizing policy definitions
