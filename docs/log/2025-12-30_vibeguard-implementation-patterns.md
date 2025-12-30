---
summary: Five concrete implementation patterns for VibeGuard policy enforcement, ranging from minimal deterministic (YAML-driven) to advanced extensible (event-based), all preserving the core design principle of binary signals and actionable next-step prompts.
event_type: deep dive
sources:
  - CLAUDE.md
tags:
  - vibeguard
  - policy-enforcement
  - architecture
  - implementation-patterns
  - llm-judges
  - yaml-driven
  - event-based
  - git-aware
---

# VibeGuard Implementation Patterns

This deep dive explores five concrete implementation patterns for VibeGuard, a policy enforcement system designed to emit actionable signals only when violations occur.

## Core Design Principle

All patterns share a common principle:
- **Silence is success** — No output when policies pass
- **Failures are structured** — Violations include machine-readable context and next-step prompts
- **Clean LLM context** — LLMs never see raw tool spam; judges receive curated context bundles
- **Actionable guidance** — Every violation answers: "What should I do next?"

## Pattern 1: Declarative Policy Runner (YAML-Driven, Deterministic First)

**Best for:** Initial implementation, tool-agnostic workflows, predictable evaluation

Policies are defined declaratively in `vibeguard.yaml`. VibeGuard is a thin orchestrator that:
1. Executes shell commands or scripts
2. Extracts and evaluates output against thresholds
3. Emits signals (non-zero exit code + JSON + optional prompt) only on violation

**Example:**
```yaml
policies:
  - id: coverage
    type: threshold
    command: "go test ./... -coverprofile=coverage.out"
    extract: "total"
    min: 80
    on_violation:
      prompt: "Code coverage dropped below 80%. Inspect uncovered files and add targeted unit tests."
```

**Strengths:** Simple, tool-agnostic, ideal for first implementation
**Limitations:** No historical or semantic judgment; deterministic only

---

## Pattern 2: Judge-as-a-Policy (LLMs as First-Class Evaluators)

**Best for:** Qualitative policy checks, semantic analysis, architectural reasoning

LLMs are treated exactly like other policies. Each judge:
1. Receives a constrained context bundle (diff, relevant files, rubric)
2. Returns a strict verdict: `{ verdict: pass | fail, reason, suggested_action }`
3. Only failures propagate upward

**Example:**
```yaml
judges:
  - id: architecture-review
    type: llm
    model: gpt-4.1
    context:
      - git_diff
      - src/**/*
      - ARCHITECTURE.md
    rubric:
      - No circular dependencies
      - Clear ownership boundaries
      - Consistent abstraction levels
    on_violation:
      prompt: "Architecture concerns detected. Review dependency graph and refactor to restore layering."
```

**Strengths:** Powerful qualitative checks, composable judges, clean LLM context
**Limitations:** Requires careful rubric design, non-deterministic by nature

---

## Pattern 3: Cloud Code–Native Wrapper (Opinionated, Integrated)

**Best for:** AI agent integration, autonomous workflows, seamless Cloud Code loops

VibeGuard acts as a companion command designed to be invoked inside agent loops:

```bash
cloud-code vibe check || cloud-code agent fix --from=vibeguard
```

Violation payload includes policy, severity, next-step prompt, and suggested tools.

**Strengths:** Seamless agent integration, zero noise in success cases, ideal for autonomous loops
**Limitations:** Initially Cloud Code–centric; less generic for other ecosystems

---

## Pattern 4: Event-Based Policy Graph (Advanced, Extensible)

**Best for:** Complex workflows, policy cascades, advanced orchestration

Policies emit events, not just results. Other policies or judges subscribe to those events, enabling cascading checks.

**Example:**
```yaml
policies:
  - id: coverage
    emits: coverage.failed

  - id: test-gap-judge
    listens: coverage.failed
    type: llm
```

When coverage fails, it triggers dependent judges (gap analysis, regression detection, etc.).

**Strengths:** Enables policy cascades, expressive and powerful, encourages separation of concerns
**Limitations:** More complex mental model; overkill for small projects

---

## Pattern 5: Git-Aware Guardrails (Change-Scope Enforcement)

**Best for:** High signal-to-noise ratio, change-scoped validation, intent-aligned policies

Policies reason about *what changed*, not just current state.

**Examples:**
- "If files under `/core` changed, require architecture judge"
- "If API schema changed, require backward-compat check"
- "If no test files touched, raise warning"

**Example:**
```yaml
policies:
  - id: api-change-requires-test
    when:
      files_changed: "api/**"
    require:
      files_changed: "tests/**"
```

**Strengths:** Very high signal-to-noise ratio, aligns with real engineering intent
**Limitations:** Requires good diff classification, slightly more logic-heavy

---

## Recommendation for Next Steps

1. **Pattern 1** (Declarative Policy Runner) is the strongest starting point — simple, tool-agnostic, immediately useful
2. **Pattern 5** (Git-Aware Guardrails) should be layered in early — dramatically improves signal-to-noise
3. **Pattern 2** (Judge-as-a-Policy) can be adopted incrementally for policies that benefit from semantic analysis
4. **Pattern 4** (Event-Based Policy Graph) can be deferred until complexity demands it
5. **Pattern 3** (Cloud Code Native) is a natural integration point once core patterns stabilize

All patterns preserve the core principle: **binary signal + actionable prompt, only on violation**.

---

## Decision Status

These patterns are **documented options**, not yet committed to a specific implementation choice. Recommend creating an ADR once the decision on which pattern(s) to prioritize is made.
