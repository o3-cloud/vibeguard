---
summary: Research on integrating OPA-inspired policy concepts into VibeGuard to prevent specific actions through declarative policy definitions
event_type: research
sources:
  - https://www.openpolicyagent.org/docs
  - https://www.openpolicyagent.org/docs/policy-language
  - https://spacelift.io/blog/what-is-open-policy-agent-and-how-it-works
  - https://www.wiz.io/academy/application-security/open-policy-agent-opa
  - https://docs.internationaldataspaces.org/ids-ram-4/layers-of-the-reference-architecture-model/3-layers-of-the-reference-architecture-model/3_4_process_layer/3_4_6_policy_enforcement
tags:
  - policy-enforcement
  - opa
  - rego
  - architecture
  - vibeguard-enhancement
  - rule-engine
  - declarative-policies
  - action-blocking
  - design-patterns
---

# Policy System Research: OPA Integration Possibilities for VibeGuard

## Overview

This research explores integrating Open Policy Agent (OPA) concepts into VibeGuard to add a policy enforcement layer that prevents specific actions from being performed when they violate defined policies. The goal is to move beyond passive checks to active policy-driven action blocking.

## Current VibeGuard Architecture

VibeGuard is a lightweight, composable policy enforcement system for CI/CD pipelines that:
- Runs shell commands and captures output
- Extracts structured data using grok patterns
- Evaluates assertions on extracted data
- Reports violations with actionable suggestions
- Executes checks in parallel with dependency management

**Current Focus:** Validation and reporting (checking if something is compliant)

## Open Policy Agent (OPA) Overview

OPA is a general-purpose policy engine that decouples policy decision-making from enforcement:

### Key Characteristics
- **Declarative Language:** Policies written in Rego (a purpose-built policy language)
- **Domain-Agnostic:** Can describe almost any kind of invariant
- **Decision-Making:** Accepts policy + input data, generates structured decisions
- **Distributed:** Can be deployed as daemon, sidecar, or embedded library
- **Fast Feedback:** In-memory policy and data for low-latency decisions
- **Rich Output:** Not limited to yes/no answers; can generate arbitrary structured data

### Rego Language Basics
```rego
package payments

default allow := false

allow if {
    input.account.state == "open"
    input.user.risk_score in ["low", "medium"]
    input.transaction.amount <= 1000
}
```

**Module Structure:**
- Package declaration
- Import statements (optional)
- Rule definitions (head + body)
- Declarative logic: "rule head is true if rule body is true"

## Policy Enforcement Architecture Patterns

### Standard PEP/PDP Model
The industry standard for policy enforcement uses two main components:

1. **Policy Decision Point (PDP):** Evaluates input against policies and makes decisions
2. **Policy Enforcement Point (PEP):** Executes/prevents actions based on PDP decisions

### Key Design Patterns
- **Interceptor Pattern:** Intercepts actions before execution for complete control over data flows
- **Policy Decoupling:** Separates policy logic from resource code for flexibility
- **Distributed PEPs:** Embed policy checks at enforcement points to reduce latency
- **Preventive vs Detective:**
  - **Preventive:** Block actions before they occur (what VibeGuard could become)
  - **Detective:** Monitor and verify compliance after actions (current VibeGuard)

## Integration Possibilities for VibeGuard

### Option 1: Preventive Policy Layer (Action Blocking)
Add a new policy type that prevents actions from being executed when they violate rules:

```yaml
version: "1"
policies:
  - id: prevent-secrets-in-commits
    type: prevent
    description: "Block commits containing secrets"
    rules:
      - match: file_content
        contains: "aws_secret_key|private_key|password="
        action: deny
        message: "Secrets detected in files"

  - id: prevent-large-files
    type: prevent
    description: "Block commits with files > 100MB"
    rules:
      - match: file_size
        greater_than: 104857600  # 100MB
        action: deny
        message: "File exceeds size limit"

checks:
  - id: validate-changes
    run: git diff HEAD --name-only --diff-filter=A
    requires_policy: prevent-secrets-in-commits
    requires_policy: prevent-large-files
```

**Pros:**
- Blocks unwanted actions before they happen
- Reuses existing VibeGuard check infrastructure
- Minimal language addition needed

**Cons:**
- Requires explicit policy evaluation in CI/CD
- Limited to declarative rules (less expressive than Rego)
- Action blocking logic must be implemented in checks themselves

### Option 2: Embedded Rego Policy Engine
Fully embed OPA/Rego into VibeGuard for rich, expressive policies:

```yaml
version: "1"
policies:
  - id: commit-policy
    type: rego
    rules: |
      package vibeguard.commit

      deny["commit violates branch protection"] {
          input.branch == "main"
          input.author not in ["release-bot", "admin"]
      }

      deny["too many files changed"] {
          count(input.changed_files) > 50
      }

      deny["risky patterns detected"] {
          input.diff_content[_] contains pattern
          pattern in ["DROP TABLE", "exec(", "eval("]
      }

checks:
  - id: check-commit-safety
    run: git show --stat
    policy: commit-policy
    action: block_commit  # If policy denies, prevent commit
```

**Pros:**
- Extremely expressive and flexible
- Rich boolean logic, aggregations, pattern matching
- Industry-standard policy language
- Reusable across teams and projects

**Cons:**
- Adds significant complexity
- New language for teams to learn (Rego)
- Larger binary size (needs Rego runtime)
- Steeper learning curve vs simple rules

### Option 3: Hybrid Approach (Recommended)
Start with simple declarative blocking rules, with hooks for advanced Rego policies:

```yaml
version: "1"
policies:
  simple_rules:
    - id: no-secrets
      type: pattern_deny
      patterns:
        - "aws_secret_key"
        - "private_key"
      message: "Secrets detected"

    - id: max-file-size
      type: constraint
      max_size: 100MB
      message: "File exceeds 100MB limit"

  advanced_rules:
    - id: commit-policy
      type: rego
      source: policies/commit.rego

checks:
  - id: validate-changes
    run: git diff HEAD
    policies: [no-secrets, max-file-size, commit-policy]
    on_policy_violation: block  # or warn, skip
```

**Pros:**
- Gradual complexity: start simple, add Rego later
- Backward compatible with existing VibeGuard checks
- Lower barrier to entry for simple use cases
- Extensible architecture

**Cons:**
- More implementation complexity upfront
- Need to maintain two rule systems

## Action Blocking Mechanisms

### Git Pre-Commit Hook Integration
VibeGuard already integrates as a pre-commit hook (per ADR-006). Policies could:
1. Evaluate staged changes against policy rules
2. Return non-zero exit code if policies violated
3. Print actionable denial messages
4. Optionally run remediation checks with `--auto-fix`

### CI/CD Pipeline Integration
```yaml
jobs:
  validate:
    steps:
      - run: vibeguard check --policy-only
        # Fails pipeline if any preventive policies violated
```

### Git Hook Advanced Flow
```
git commit --message "..."
  → vibeguard check --hooks
    → Load policies
    → Evaluate staged changes
    → If deny: print message, exit 1 → commit blocked
    → If allow: exit 0 → commit proceeds
```

## Comparison: VibeGuard Checks vs Policies

| Aspect | Current Checks | Proposed Policies |
|--------|---|---|
| **Purpose** | Validate state (testing, coverage, linting) | Prevent actions (blocking commits, pushes) |
| **Timing** | After action | Before action |
| **Input** | Command output, files | Action metadata (commit, push, PR) |
| **Output** | Pass/Fail + suggestions | Allow/Deny + reason |
| **Rules** | Shell commands + assertions | Declarative patterns or Rego rules |
| **Action** | Report results | Block or allow action |

## Implementation Roadmap

### Phase 1: Simple Pattern-Based Blocking
- Add `policy` config section for simple deny patterns
- Implement pattern matching against input data
- Integrate with existing git hook infrastructure
- Output: `allow` or `deny` decision

### Phase 2: Constraint-Based Rules
- Add size, count, and threshold constraints
- Boolean logic: `AND`, `OR`, `NOT`
- Named rules for reusability

### Phase 3: Rego Integration (Optional)
- Embed Rego runtime via `github.com/open-policy-agent/opa/rego` package
- Load `.rego` files in config
- Rich boolean logic, aggregations, pattern matching

## Key Findings

1. **OPA's Architecture is Solid:** PDP/PEP separation is industry standard and applicable to VibeGuard
2. **Rego is Powerful but Complex:** Excellent for advanced policies, but adds learning curve
3. **Hybrid Approach is Pragmatic:** Start with simple pattern-based rules, add Rego later
4. **VibeGuard is Already Well-Positioned:**
   - Already has check execution framework
   - Already integrates as pre-commit hook (ADR-006)
   - Can reuse variable interpolation and expression evaluation
5. **Action Blocking is Distinct from Validation:**
   - Current VibeGuard checks answer: "Is the system compliant?"
   - Policies would answer: "Should this action be allowed?"
   - Both valuable, complementary capabilities

## Recommended Next Steps

1. **Create ADR for Policy System:** Document the decision to add policies, architecture approach
2. **Design Policy Schema:** Define YAML structure for simple pattern-based policies
3. **Prototype Phase 1:** Implement basic pattern matching and blocking logic
4. **Get Feedback:** Test with a simple use case (e.g., block commits with secrets)
5. **Plan Rego Integration:** Evaluate if/when Rego would add value

## Related ADRs

- [ADR-005: Adopt Vibeguard for Policy Enforcement in CI/CD](docs/adr/ADR-005-adopt-vibeguard.md)
- [ADR-006: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement](docs/adr/ADR-006-integrate-vibeguard-as-claude-code-hook.md)
