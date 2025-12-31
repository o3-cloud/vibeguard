---
summary: LLM judges are just CLI commands (claude, gemini, ollama) using the same checks schema. Suggestions support {{var}} templating from grok-extracted values for dynamic output.
event_type: deep dive
sources:
  - docs/log/2025-12-30_simplified-checks-schema.md
tags:
  - vibeguard
  - llm-as-judge
  - templating
  - grok
  - claude-cli
  - gemini-cli
  - suggestions
---

# LLM as Judge & Templated Suggestions

## Core Insight

LLM judges don't need a special abstraction. They're just CLI commands that output verdicts. Use the same `checks` schema with `run`, `grok`, and `assert`.

## LLM as a Check

```yaml
checks:
  - id: architecture-review
    run: |
      claude -p "Review this diff for architectural issues.
      Output exactly: PASS or FAIL: <reason>
      $(git diff main --staged)"
    grok: '%{WORD:verdict}: %{GREEDYDATA:reason}'
    assert: verdict == "PASS"
    suggestion: "{{reason}}"
```

Same schema. LLM is just another tool.

## Available CLI Tools

| Tool | Command | Notes |
|------|---------|-------|
| Claude Code | `claude -p "prompt"` | Anthropic CLI |
| Gemini | `gemini -p "prompt"` | Google CLI |
| OpenAI | `openai api chat.completions.create` | OpenAI CLI |
| Ollama | `ollama run llama3 "prompt"` | Local models |
| LLM | `llm "prompt"` | Simon Willison's multi-provider CLI |

Any CLI that takes a prompt and outputs text works.

## Templated Suggestions

Suggestions support `{{var}}` templating with grok-extracted values:

```yaml
- id: coverage
  run: go test -cover ./...
  grok: 'coverage: %{NUMBER:coverage}%'
  assert: coverage >= 80
  suggestion: "Coverage is {{coverage}}%, need 80%."
```

Output:
```
FAIL  coverage (error)
      → go test -cover ./...

      Tip: Coverage is 72.5%, need 80%.
```

## LLM Examples

### Code Review with Reasoning

```yaml
- id: llm-review
  run: |
    claude -p "Review. Output: VERDICT: PASS/FAIL | REASON: <why>
    $(git diff main)"
  grok: 'VERDICT: %{WORD:verdict} | REASON: %{GREEDYDATA:reason}'
  assert: verdict == "PASS"
  suggestion: "{{reason}}"
```

Output:
```
FAIL  llm-review (warning)
      → claude -p "Review..."

      Tip: Missing error handling in auth.go:45, potential nil pointer.
```

### Security Review with JSON

```yaml
- id: security-review
  run: |
    gemini -p "Security audit. Output JSON: {\"safe\": bool, \"reason\": string}
    $(cat src/auth/*.go)"
  assert: json.safe == true
  suggestion: "{{json.reason}}"
```

### PR Description Quality

```yaml
- id: pr-quality
  run: |
    claude -p "Rate this PR description 1-10 for clarity.
    Output: SCORE: <number> | FEEDBACK: <why>
    $(gh pr view --json body -q .body)"
  grok: 'SCORE: %{INT:score} | FEEDBACK: %{GREEDYDATA:feedback}'
  assert: score >= 7
  suggestion: "{{feedback}}"
```

## Multiple Extractions

```yaml
- id: complexity
  run: gocyclo -top 5 .
  grok: '%{INT:score} %{WORD:func} %{PATH:file}:%{INT:line}'
  assert: score <= 15
  suggestion: "{{func}} in {{file}}:{{line}} has complexity {{score}}."
```

Output:
```
FAIL  complexity (error)
      → gocyclo -top 5 .

      Tip: calculateTotal in src/billing.go:142 has complexity 23.
```

## Full Config Example

```yaml
version: "1"

checks:
  # Fast deterministic checks first
  - id: fmt
    run: "! gofmt -l . | grep ."

  - id: vet
    run: go vet ./...

  - id: lint
    run: golangci-lint run

  - id: test
    run: go test ./...

  - id: coverage
    run: go test -cover ./...
    grok: 'coverage: %{NUMBER:coverage}%'
    assert: coverage >= 80
    suggestion: "Coverage is {{coverage}}%, need 80%."

  # LLM checks (slower, run after deterministic)
  - id: llm-review
    run: |
      claude -p "Review this diff. Output: VERDICT: PASS/FAIL | REASON: <why>
      $(git diff main --staged)"
    grok: 'VERDICT: %{WORD:verdict} | REASON: %{GREEDYDATA:reason}'
    assert: verdict == "PASS"
    requires: [fmt, vet, lint, test]
    timeout: 60s
    severity: warning
    suggestion: "{{reason}}"
```

## JSON Output

```json
{
  "violations": [
    {
      "id": "coverage",
      "severity": "error",
      "command": "go test -cover ./...",
      "suggestion": "Coverage is 72.5%, need 80%.",
      "extracted": {
        "coverage": "72.5"
      }
    }
  ]
}
```

Include `extracted` for tooling that wants raw values.

## Prompt Engineering Tips

For reliable LLM output:

1. **Constrain format strictly**
   ```
   Output ONLY one line.
   Format: VERDICT: PASS or FAIL | REASON: <brief explanation>
   ```

2. **JSON is more reliable than free-form**
   ```
   Output JSON only: {"pass": boolean, "reason": string}
   ```

3. **Give examples**
   ```
   Examples:
   VERDICT: PASS | REASON: Code follows best practices
   VERDICT: FAIL | REASON: Missing input validation
   ```

## Fallback Behavior

| Scenario | Suggestion Output |
|----------|-------------------|
| `{{var}}` exists | Replaced with extracted value |
| `{{var}}` missing | Left as empty or literal |
| No suggestion defined | No "Tip:" line shown |

## Design Benefits

| Benefit | How |
|---------|-----|
| No special LLM integration | CLI tools handle auth, API, rate limits |
| Same schema | `run` + `grok` + `assert` + `suggestion` |
| Dynamic output | `{{var}}` templating surfaces extracted values |
| Composable | Chain LLM checks after fast checks via `requires` |
| Swappable | Replace `claude` with `gemini` or local `ollama` |

---

**Status:** Design complete
