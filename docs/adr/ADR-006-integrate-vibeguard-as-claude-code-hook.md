# ADR-006: Integrate VibeGuard as Git Pre-Commit Hook for Policy Enforcement

## Context and Problem Statement

The VibeGuard project has adopted itself as the policy enforcement system for the repository (ADR-005). Currently, developers must manually run `vibeguard check` to validate code against policies. Git provides a hook system that automatically executes shell commands at key lifecycle points, including before commits are finalized.

The problem: there is no automated enforcement of VibeGuard policies within the development workflow, meaning policy violations can be committed without being caught.

## Considered Options

**Option A: Git pre-commit hooks only**
- Use traditional git pre-commit hooks to run vibeguard before commits
- Pros: Standard approach, works with any editor/IDE, well-understood mechanism
- Cons: Only catches issues at commit time

**Option B: Claude Code hooks in project-level settings**
- Configure vibeguard to run as a Claude Code hook via `.claude/settings.json`
- Pros: Catches issues early in Claude Code workflows
- Cons: Only works with Claude Code, adds latency to all operations, depends on Claude Code-specific mechanism

**Option C: Both git pre-commit hooks AND Claude Code hooks**
- Implement comprehensive coverage with both mechanisms
- Pros: Defense-in-depth, catches issues at multiple points
- Cons: Redundant execution, increased complexity

**Option D: Manual execution only**
- Developers run vibeguard manually when needed
- Pros: No automation overhead
- Cons: Inconsistent enforcement, easy to skip

## Decision Outcome

**Chosen option:** Option A - Git pre-commit hooks only

**Rationale:**
1. **Universal applicability**: Git hooks work regardless of editor, IDE, or development tool
2. **Standard mechanism**: Pre-commit hooks are a well-understood, battle-tested approach
3. **Minimal overhead**: Only runs when committing, not on every operation
4. **No tool lock-in**: Works with any development workflow, not just Claude Code
5. **Simpler maintenance**: Single enforcement point reduces complexity

**Tradeoffs:**
- Gained: Universal compatibility, simpler configuration, no per-operation latency
- Sacrificed: Earlier feedback during interactive development (policies checked at commit time, not during editing)

## Consequences

**Positive outcomes:**
- Policy violations are caught before code is committed
- Works with any development tool or workflow
- Simple, well-understood enforcement mechanism
- No dependency on specific development tools
- Easy to share and maintain via version control

**Negative outcomes:**
- Requires vibeguard binary to be built and available in the environment
- Only catches issues at commit time, not during interactive development
- Developers can bypass with `--no-verify` flag

**Neutral impacts:**
- Hook script lives in `.git/hooks/` (not version-controlled by default) or managed via tooling
- Complements but doesn't replace manual vibeguard execution for pre-commit validation
