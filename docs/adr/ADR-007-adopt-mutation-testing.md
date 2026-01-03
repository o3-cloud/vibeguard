# ADR-007: Adopt Gremlins for Mutation Testing

## Status
Accepted

## Context and Problem Statement

While ADR-004 establishes code coverage requirements (70% minimum), code coverage alone is an incomplete measure of test suite quality. A test suite can achieve high coverage while still missing critical bugs because:

- Tests may execute code without actually verifying behavior
- Assertions may be too weak or missing entirely
- Edge cases and boundary conditions may be untested
- Tests may pass regardless of whether the code is correct

We need a way to measure and improve test suite effectiveness beyond simple line coverage metrics.

## Considered Options

### Option A: Rely solely on code coverage (Current approach)
- Continue using `go test -cover` for coverage metrics
- Rely on code review to assess test quality
- **Pros**: No additional tooling, simple to understand
- **Cons**: Coverage can be gamed, doesn't measure assertion quality, false confidence

### Option B: Adopt Gremlins for mutation testing (SELECTED)
- Use [Gremlins](https://gremlins.dev) to introduce small mutations and verify tests catch them
- Run periodically to identify weak tests
- Use mutation score as a complementary quality metric
- **Pros**: Modern tooling, YAML configuration, PR-diff support, good documentation, active development
- **Cons**: Longer test execution time, 0.x versioning (API may change), not suited for very large codebases

### Option C: Adopt go-mutesting
- Use go-mutesting as the mutation testing framework
- **Pros**: Established tool, simple usage
- **Cons**: Less actively maintained, no YAML configuration, limited documentation

### Option D: Manual test review process
- Establish formal test review guidelines
- Require explicit approval of test quality during code review
- **Pros**: Human judgment on test intent
- **Cons**: Subjective, time-consuming, inconsistent, doesn't scale

## Decision Outcome

**Chosen option:** Option B - Adopt Gremlins for mutation testing

**Rationale:**
- Mutation testing provides an objective measure of test suite effectiveness
- Complements existing coverage requirements from ADR-004
- Gremlins offers modern tooling with YAML-based configuration (consistent with VibeGuard's approach)
- Supports testing only PR changes, reducing CI time for incremental work
- Well-documented with official website and guides
- Docker and GitHub Actions support for easy CI integration

**Tradeoffs:**
- Significant increase in test execution time (runs tests multiple times per mutation)
- Not suitable for pre-commit hooks; better suited for CI or periodic analysis
- Still in 0.x versioning - configuration may change between minor releases
- Works best on smallish Go modules (microservices scale)

## Consequences

### Positive Outcomes
1. **Improved test quality** - Tests that survive mutation testing are genuinely effective
2. **Identified weak spots** - Pinpoints specific code locations where tests are insufficient
3. **Objective metric** - Mutation score provides quantifiable test effectiveness measure
4. **PR-focused testing** - Can run mutations only on changed code for faster feedback
5. **Developer education** - Teaches developers to write more effective assertions

### Negative Outcomes
1. **Execution time** - Mutation testing is significantly slower than regular tests
2. **Version instability** - 0.x versioning means potential breaking changes
3. **Scale limitations** - May perform poorly on very large codebases
4. **CI resource usage** - Increases compute time for mutation testing runs

### Neutral Impacts
1. **Complementary to coverage** - Does not replace coverage metrics, adds another dimension
2. **Optional integration** - Can be run on-demand rather than on every commit

## Implementation Details

### 1. Installation
```bash
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
```

Or via Docker:
```bash
docker pull gogremlins/gremlins
```

### 2. Configuration
Create `.gremlins.yaml` in the repository root:
```yaml
# .gremlins.yaml
unleash:
  # Only test mutations in code covered by tests
  dry-run: false
  # Output format
  output: colored
  # Timeout for each mutation test
  timeout-coefficient: 1

mutants:
  # Enable all mutation operators
  arithmetic-base: true
  conditionals-boundary: true
  conditionals-negation: true
  increment-decrement: true
  invert-assignments: true
  invert-bitwise: true
  invert-booleans: true
  invert-logical: true
  invert-loopctrl: true
  invert-negatives: true
  remove-self-assignments: true
```

### 3. Basic Usage
```bash
# Run mutation testing on entire module
gremlins unleash

# Run on specific package
gremlins unleash ./internal/policy/...

# Run only on files changed in PR (diff mode)
gremlins unleash --diff
```

### 4. Mutation Operators
Gremlins applies the following mutation types:

| Operator | Description | Example |
|----------|-------------|---------|
| `arithmetic-base` | Swap arithmetic operators | `+` to `-` |
| `conditionals-boundary` | Change boundary conditions | `<` to `<=` |
| `conditionals-negation` | Negate conditions | `==` to `!=` |
| `increment-decrement` | Swap inc/dec | `++` to `--` |
| `invert-assignments` | Swap assignment operators | `+=` to `-=` |
| `invert-bitwise` | Swap bitwise operators | `&` to `\|` |
| `invert-booleans` | Flip boolean values | `true` to `false` |
| `invert-logical` | Swap logical operators | `&&` to `\|\|` |
| `invert-loopctrl` | Swap loop control | `break` to `continue` |
| `invert-negatives` | Remove negation | `-x` to `x` |
| `remove-self-assignments` | Remove self-assignments | Remove `x += 0` |

### 5. Interpreting Results
```
Mutation testing completed
Total mutants: 60
Killed: 45 (75.00%)
Survived: 12 (20.00%)
Not covered: 3 (5.00%)
```

- **Killed**: Test suite detected and failed on the mutation (good)
- **Survived**: Test suite passed despite the mutation (weak test)
- **Not covered**: Mutation in code not covered by tests
- **Score**: Ratio of killed to total mutations (higher is better)

### 6. CI Integration
Add as a GitHub Actions workflow:
```yaml
# .github/workflows/mutation.yml
name: Mutation Testing

on:
  schedule:
    - cron: '0 2 * * 0'  # Weekly on Sunday at 2 AM
  workflow_dispatch:  # Manual trigger

jobs:
  mutation:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run Gremlins
        uses: go-gremlins/gremlins-action@v1
        with:
          version: latest
          args: --tags="" --timeout-coefficient=2

      - name: Upload mutation report
        uses: actions/upload-artifact@v4
        with:
          name: mutation-report
          path: gremlins-report.*
```

For PR-only diff testing:
```yaml
# In PR workflow
- name: Run Gremlins on PR diff
  uses: go-gremlins/gremlins-action@v1
  with:
    args: --diff
```

### 7. Quality Targets
- **Initial target**: 50% mutation score
- **Long-term target**: 70% mutation score
- **Focus areas**: Core policy evaluation, security-critical paths

## Quality Gates

Mutation testing is advisory, not blocking:
- Regular CI: Not included (too slow)
- Weekly CI job: Generate mutation report
- PR CI (optional): Run diff-only mutation testing
- Pre-release: Review mutation score for critical packages
- Code review: Consider mutation results when improving test coverage

## Related Documentation

- [ADR-004: Code Quality Standards](ADR-004-code-quality-standards.md) - Testing standards this complements
- [Gremlins Documentation](https://gremlins.dev/latest/) - Official documentation
- [Gremlins GitHub](https://github.com/go-gremlins/gremlins) - Source repository

## References

- Gremlins documentation: https://gremlins.dev/latest/
- Mutation testing concept: https://en.wikipedia.org/wiki/Mutation_testing
