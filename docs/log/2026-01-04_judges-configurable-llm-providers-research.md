---
summary: Research on implementing judges as AI agents in vibeguard with configurable provider support (Claude, Gemini, Codex, Ollama). Judges are specialized checks that review criteria, provide binary pass/fail verdicts with reasoning.
event_type: research
sources:
  - examples/advanced.yaml (LLM judge patterns)
  - docs/log/2025-12-30_llm-as-judge-and-templated-suggestions.md
  - docs/log/2025-12-30_vibeguard-implementation-patterns.md
  - docs/adr/ADR-003-adopt-golang.md (single-binary deployment)
tags:
  - vibeguard
  - judges
  - llm-integration
  - ai-agents
  - provider-abstraction
  - claude
  - gemini
  - codex
  - ollama
  - policy-enforcement
---

# Judges with Configurable LLM Providers

## Executive Summary

Judges are specialized checks within vibeguard that invoke AI agents to perform semantic analysis and return binary pass/fail verdicts with reasoning. The current implementation treats judges as CLI commands (e.g., `claude`, `gemini`, `ollama`), but expanding this to support configurable providers with structured input/output would improve reliability, testability, and user experience.

## Current State

### How Judges Work Today

From `examples/advanced.yaml`, judges are already functional via CLI invocation:

```yaml
- id: llm-architecture-review
  run: |
    claude -p "Review this Go code diff for architectural issues.
    Output ONLY one line in this exact format:
    VERDICT: PASS or FAIL | REASON: <brief explanation>

    $(git diff HEAD~1 --name-only -- '*.go' | head -5 | xargs cat 2>/dev/null || echo 'No Go files changed')"
  grok:
    - 'VERDICT: %{WORD:verdict} \| REASON: %{GREEDYDATA:reason}'
  assert: 'verdict == "PASS"'
  severity: warning
  suggestion: "{{.reason}}"
  timeout: 60s
  requires:
    - fmt
    - vet
    - lint
    - test
```

**Key characteristics:**
- Judges are just checks with `run`, `grok`, `assert`, and `suggestion` fields
- No special abstraction—they use shell commands to invoke LLM CLIs
- Grok patterns extract the verdict and reasoning
- Assert evaluates the verdict as a boolean
- Suggestion templates surface the reasoning to users

### Supported CLI Tools (Implicit)

| Provider | CLI Command | Status | Notes |
|----------|-------------|--------|-------|
| Claude (Anthropic) | `claude -p "prompt"` | Documented | Works with Claude Code CLI |
| Gemini (Google) | `gemini -p "prompt"` | Documented | Google CLI tool |
| Ollama (Local) | `ollama run <model> "prompt"` | Documented | Offline, open-source models |
| OpenAI | `openai api chat.completions.create` | Possible | CLI interface exists |
| LLM (Simon Willison) | `llm "prompt"` | Possible | Multi-provider abstraction CLI |

## Proposed Enhancement: Judge Type with Provider Abstraction

### Problem Statement

While the CLI-based approach works, expanding judges to be first-class citizens with structured provider configuration would provide:

1. **Provider Flexibility** — Switch providers via config without rewriting shell commands
2. **Structured Input/Output** — Use JSON for cleaner context bundling and parsing
3. **Testability** — Mock providers in unit tests without shell execution
4. **Error Handling** — Consistent error handling across provider implementations
5. **Authentication** — Centralized API key/credential management
6. **Caching** — Optional response caching to reduce API calls
7. **Rate Limiting** — Graceful handling of provider rate limits
8. **Observability** — Detailed logging of judge invocations and reasoning

### Proposed Schema Extension

Add an optional `judge` field to checks schema:

```yaml
checks:
  - id: architecture-review
    judge:
      provider: claude          # claude | gemini | codex | ollama | openai | custom
      model: claude-opus-4.5    # Model identifier for the provider
      context:
        - git_diff              # Built-in context: git diff
        - file_glob: "cmd/**"   # Include files matching glob
        - custom: |             # Custom context
            Architecture rules:
            1. No circular dependencies
            2. Clear ownership boundaries
      rubric:
        criteria:
          - name: "Separation of Concerns"
            weight: 30
            description: "Code is organized into clear, focused modules"
          - name: "Dependency Management"
            weight: 40
            description: "Dependencies flow in one direction; no cycles"
          - name: "Interface Design"
            weight: 30
            description: "Interfaces are narrow and purpose-specific"
      output:
        format: json            # json | text (with expected structure)
        schema:                 # Expected response schema
          verdict: boolean      # true = pass, false = fail
          reasoning: string     # Explanation of verdict
          issues: array         # [{ line, severity, suggestion }]
      timeout: 90s
      retry:
        max_attempts: 2
        backoff: exponential
    severity: warning
    suggestion: "{{reasoning}}"  # Template from judge output
    requires: [fmt, vet, lint]
```

### Provider-Specific Implementation

#### Claude Provider

```yaml
judge:
  provider: claude
  model: claude-opus-4.5
  api_key: ${ANTHROPIC_API_KEY}  # From environment or secrets
  config:
    max_tokens: 1024
    temperature: 0.2  # Low for deterministic verdicts
```

**Implementation approach:**
- Use Anthropic SDK directly (no CLI shell invocation)
- Messages API for structured requests/responses
- Tool use for JSON schema enforcement

#### Gemini Provider

```yaml
judge:
  provider: gemini
  model: gemini-2.0-flash
  api_key: ${GOOGLE_API_KEY}
  config:
    temperature: 0.2
    max_output_tokens: 1024
```

#### Ollama Provider (Local)

```yaml
judge:
  provider: ollama
  model: neural-chat:7b
  endpoint: http://localhost:11434  # Local or remote endpoint
  config:
    temperature: 0.2
```

**Key advantage:** Completely offline, no API keys, useful for CI/CD without external dependencies

#### OpenAI/Codex Provider

```yaml
judge:
  provider: openai
  model: gpt-4
  api_key: ${OPENAI_API_KEY}
  config:
    temperature: 0.2
    max_tokens: 1024
```

#### Custom/Webhook Provider

```yaml
judge:
  provider: custom
  endpoint: https://internal-llm.company.com/judge
  method: POST
  headers:
    Authorization: "Bearer ${CUSTOM_API_KEY}"
  request_template: |
    {
      "instruction": "Review code for architecture",
      "context": "{{.context}}",
      "rubric": "{{.rubric}}"
    }
  response_parser: |
    # jq-like syntax to extract verdict and reasoning
    {
      "verdict": .pass,
      "reasoning": .explanation
    }
```

## Implementation Architecture

### Judge Package Structure

```
internal/judge/
├── judge.go              # Core Judge struct and interface
├── provider.go           # Provider interface and factory
├── providers/
│   ├── claude.go         # Claude provider implementation
│   ├── gemini.go         # Gemini provider implementation
│   ├── ollama.go         # Ollama provider implementation
│   ├── openai.go         # OpenAI provider implementation
│   └── custom.go         # Custom HTTP endpoint provider
├── context.go            # Context bundler (git diff, files, etc.)
├── response/
│   └── parser.go         # Response parsing and validation
└── cache.go              # Optional response caching
```

### Core Interfaces

```go
// judge/provider.go
type Provider interface {
    // Invoke sends a request to the provider and returns verdict + reasoning
    Invoke(ctx context.Context, req *InvokeRequest) (*Response, error)

    // Name returns the provider identifier
    Name() string

    // Config returns the provider's configuration
    Config() ProviderConfig
}

type InvokeRequest struct {
    Instruction string            // The rubric/criteria
    Context     map[string]string // Bundled context (diff, files, etc.)
    Rubric      []Criterion       // Evaluation criteria
    Timeout     time.Duration
}

type Response struct {
    Verdict    bool              // true = pass, false = fail
    Reasoning  string            // Explanation
    Issues     []Issue           // Optional detailed issues
    Extracted  map[string]string // For {{var}} templating
    RawResponse string           // Original response for debugging
}

// judge/judge.go
type Judge struct {
    ID       string
    Provider Provider
    Config   JudgeConfig
}

func (j *Judge) Evaluate(ctx context.Context, bundledContext map[string]string) (*Response, error) {
    req := &InvokeRequest{
        Instruction: j.Config.Instruction,
        Context:     bundledContext,
        Rubric:      j.Config.Criteria,
        Timeout:     j.Config.Timeout,
    }
    return j.Provider.Invoke(ctx, req)
}
```

### Context Bundling

Judges need curated input. The context bundler should support:

```go
type ContextBundler struct {
    GitDiff    bool              // Include staged/unstaged changes
    Files      []string          // File globs to include
    Custom     []string          // Custom context strings
    MaxSize    int               // Max context size (prevent token limits)
}

func (cb *ContextBundler) Bundle(ctx context.Context) (map[string]string, error) {
    result := make(map[string]string)

    if cb.GitDiff {
        diff, err := executeGitDiff()
        result["git_diff"] = diff
    }

    for _, glob := range cb.Files {
        contents, err := bundleFiles(glob)
        result[fmt.Sprintf("files_%s", glob)] = contents
    }

    return result, nil
}
```

### Response Parsing

Structured parsing ensures consistent evaluation:

```go
type ResponseParser interface {
    Parse(rawResponse string) (*ParsedResponse, error)
}

// Built-in JSON parser
type JSONParser struct {
    Schema *jsonschema.Schema
}

func (jp *JSONParser) Parse(raw string) (*ParsedResponse, error) {
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(raw), &data); err != nil {
        return nil, fmt.Errorf("invalid JSON: %w", err)
    }

    // Validate against schema
    if err := jp.Schema.Validate(data); err != nil {
        return nil, fmt.Errorf("response doesn't match schema: %w", err)
    }

    return &ParsedResponse{
        Verdict:   data["verdict"].(bool),
        Reasoning: data["reasoning"].(string),
        Issues:    parseIssues(data["issues"]),
    }, nil
}
```

## Integration with Orchestrator

### Changes to ExecutionFlow

1. **Check Type Detection**
   - In `orchestrator.Run()`, detect if check has `judge` field
   - Route to `judge.Evaluate()` instead of shell execution

2. **Result Mapping**
   - Judge response maps to existing `Result` struct
   - `Response.Verdict` → Check pass/fail
   - `Response.Reasoning` → Suggestion content
   - `Response.Issues` → Additional violation details

3. **Dependency Ordering**
   - Judges can depend on other checks via `requires`
   - Judges are treated as normal checks in DAG

4. **Timeout Handling**
   - Judge timeout is independent of check timeout
   - Context cancellation propagates to provider

### Updated Check Schema

```yaml
checks:
  - id: architecture-review

    # Option 1: Traditional shell command
    run: |
      claude -p "Review this diff..."
      $(git diff)
    grok: 'VERDICT: %{WORD:verdict} \| REASON: %{GREEDYDATA:reason}'

    # Option 2: Structured judge (new)
    judge:
      provider: claude
      model: claude-opus-4.5
      context:
        - git_diff
        - file_glob: "cmd/**/*.go"
      rubric:
        - "No circular dependencies"
        - "Clear ownership boundaries"
      output_format: json

    assert: 'verdict == "PASS"'
    severity: warning
    suggestion: "{{reasoning}}"
```

Both options work; the structured judge offers better ergonomics and testability.

## Provider Selection Criteria

### Claude (Recommended for vibeguard dogfooding)

| Criterion | Rating | Notes |
|-----------|--------|-------|
| Quality | ⭐⭐⭐⭐⭐ | Latest model (Opus 4.5) very capable for code review |
| Cost | ⭐⭐⭐ | Higher cost but excellent quality |
| Latency | ⭐⭐⭐⭐ | Fast API, good for CI/CD |
| Availability | ⭐⭐⭐⭐⭐ | Reliable SLA |
| Local Option | ✗ | API-only |
| **Best For** | Architecture/security reviews, high-stakes decisions |

**Why Claude is good for vibeguard:**
- Aligns with ADR-005 (dogfood vibeguard on vibeguard)
- Can leverage Claude Code integration
- Excellent reasoning transparency (good for explaining violations)

### Gemini

| Criterion | Rating | Notes |
|-----------|--------|-------|
| Quality | ⭐⭐⭐⭐ | Strong model, good for code |
| Cost | ⭐⭐⭐⭐ | Competitive pricing |
| Latency | ⭐⭐⭐⭐ | Fast |
| Availability | ⭐⭐⭐⭐⭐ | Reliable |
| Local Option | ✗ | API-only |
| **Best For** | Cost-conscious projects, multi-modal (if using images) |

### Ollama (Local)

| Criterion | Rating | Notes |
|-----------|--------|-------|
| Quality | ⭐⭐⭐ | Depends on model (llama, neural-chat, codellama) |
| Cost | ⭐⭐⭐⭐⭐ | Free (runs locally) |
| Latency | ⭐⭐⭐ | Depends on hardware; typically slower than cloud |
| Availability | ⭐⭐⭐⭐⭐ | No API dependency |
| Local Option | ✓ | Fully offline |
| **Best For** | Privacy-sensitive projects, no internet CI, learning |

**Limitation:** Smaller models (7B, 13B) often struggle with complex code reasoning; 70B+ recommended for production.

### OpenAI (Codex/GPT-4)

| Criterion | Rating | Notes |
|-----------|--------|-------|
| Quality | ⭐⭐⭐⭐ | GPT-4 is excellent; Codex (legacy) discontinued |
| Cost | ⭐⭐⭐ | Higher cost than alternatives |
| Latency | ⭐⭐⭐⭐ | Good |
| Availability | ⭐⭐⭐⭐ | Reliable |
| Local Option | ✗ | API-only |
| **Best For** | Teams already invested in OpenAI ecosystem |

## Configuration Management

### Environment-Based Authentication

```bash
# Claude
export ANTHROPIC_API_KEY=sk-ant-...

# Gemini
export GOOGLE_API_KEY=AIz...

# OpenAI
export OPENAI_API_KEY=sk-...

# Custom endpoint
export CUSTOM_API_KEY=...
```

### Vibeguard Config

```yaml
version: "1"

# Global judge configuration
judges:
  default_provider: claude
  default_model: claude-opus-4.5
  providers:
    claude:
      api_key: ${ANTHROPIC_API_KEY}
      timeout: 90s
    gemini:
      api_key: ${GOOGLE_API_KEY}
      timeout: 60s
    ollama:
      endpoint: http://localhost:11434
      timeout: 120s

checks:
  - id: architecture-review
    judge:
      provider: claude  # Uses default model from config
      context: [git_diff]
      rubric: [...]
```

## Error Handling

### Judge Failure Scenarios

| Scenario | Behavior | Example |
|----------|----------|---------|
| Provider unavailable | Timeout, retry, then fail check | API down, no network |
| Invalid response | Parse error, log raw response, fail check | JSON malformed |
| Rate limited | Exponential backoff, respect Retry-After | Too many requests |
| Context too large | Truncate context, warn, proceed | Large diff |
| Schema validation | Response doesn't match expected format | Missing `verdict` field |
| Provider error | Return error message in suggestion | "API key invalid" |

### Error Messages to Users

```
FAIL  architecture-review (error)
      → Judge invocation failed: Claude provider timeout after 90s

      Tip: Check that ANTHROPIC_API_KEY is set and API is accessible.
      Suggestion: Re-run after checking network connectivity.
```

## Testing Strategy

### Unit Tests

```go
// judge/claude_test.go
func TestClaudeProvider_ValidResponse(t *testing.T) {
    // Mock HTTP responses
    // Test parsing of valid verdict + reasoning
    // Test template substitution
}

func TestClaudeProvider_InvalidResponse(t *testing.T) {
    // Test handling of malformed JSON
    // Test fallback behaviors
}

// judge/context_test.go
func TestContextBundler_GitDiff(t *testing.T) {
    // Git diff extraction
}

func TestContextBundler_FilesGlob(t *testing.T) {
    // File inclusion by glob
}

func TestContextBundler_SizeLimit(t *testing.T) {
    // Context truncation when too large
}
```

### Integration Tests

```bash
# End-to-end judge evaluation
vibeguard check --judge-provider=ollama --judge-model=neural-chat:7b

# Mock provider (for CI/CD without API keys)
vibeguard check --judge-provider=mock
```

### Prompt Testing (Iteration)

```go
// Test framework for judge prompts
type JudgePromptTest struct {
    Name     string
    Code     string
    Expected bool
}

tests := []JudgePromptTest{
    {
        Name: "Clean code",
        Code: "func Add(a, b int) int { return a + b }",
        Expected: true,
    },
    {
        Name: "Circular dependency",
        Code: "a -> b -> a",
        Expected: false,
    },
}
```

## Migration Path

### Phase 1: CLI-Based Judges (Current)

- Judges work via shell commands (`claude`, `gemini`, `ollama` CLIs)
- No code changes required
- Fully functional today

### Phase 2: Structured Judge Type (Proposed)

1. Add `judge` field to check schema (backward compatible)
2. Implement provider interface
3. Implement 3-4 core providers (Claude, Gemini, Ollama, custom)
4. Update orchestrator to detect and route judge checks
5. Add integration tests

### Phase 3: Enhanced Features

1. Judge response caching
2. Rate limiting
3. Streaming responses
4. Custom parsers/validators
5. Multi-judge aggregation

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| Judges as checks, not separate | Consistent with vibeguard philosophy: everything is a check |
| Provider abstraction | Flexibility to swap providers without config rewrite |
| Structured context bundling | Prevents token limit issues, cleaner LLM input |
| JSON response format | More reliable than free-form parsing |
| Timeout + retry | Handle transient API failures gracefully |
| Environment-based auth | Standard practice, works with secrets systems |
| Mock provider in tests | No external dependencies in unit tests |

## Next Steps

1. **Design Review**
   - Get team feedback on schema and provider abstraction
   - Discuss error handling edge cases
   - Confirm provider priorities

2. **Prototype Phase**
   - Implement Claude provider (best bang for buck)
   - Add integration test with mock API
   - Test with real vibeguard config

3. **Documentation**
   - Create user guide for judge configuration
   - Provider setup guides (API keys, endpoints)
   - Prompt engineering tips
   - Troubleshooting guide

4. **Example Configs**
   - Architecture review judge
   - Security review judge
   - PR quality judge
   - Test coverage analysis judge

## References

- **Current LLM judge patterns**: `examples/advanced.yaml`
- **LLM as judge deep dive**: `docs/log/2025-12-30_llm-as-judge-and-templated-suggestions.md`
- **Implementation patterns**: `docs/log/2025-12-30_vibeguard-implementation-patterns.md`
- **Architecture Decision Records**: `docs/adr/`
  - ADR-003: Go as primary language
  - ADR-005: Adopt vibeguard for CI/CD
  - ADR-006: Git pre-commit hook integration

---

**Status:** Research complete, ready for design review and prototype phase

**Created:** 2026-01-04

