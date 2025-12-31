---
summary: Discovered Claude Code hook exit code behavior - exit 2 blocks actions, exit 1 is non-blocking
event_type: research
sources:
  - https://docs.anthropic.com/en/docs/claude-code/hooks
tags:
  - claude-code
  - hooks
  - vibeguard
  - exit-codes
  - policy-enforcement
---

# Claude Code Hook Exit Code Behavior

## Key Finding

Claude Code hooks use exit codes to communicate action decisions:

| Exit Code | Behavior |
|-----------|----------|
| **0** | Success - action continues normally |
| **2** | Block - action is prevented, stderr shown to user/Claude |
| **Other (1, etc.)** | Non-blocking error - logged but execution continues |

## Implications for Vibeguard

Currently, `vibeguard check` exits with code **1** on policy violations. This means:
- Policy violations are logged as errors
- But the action (prompt submission, file edit) **continues anyway**

To actually enforce policies and block actions, vibeguard should exit with code **2** instead of **1** when checks fail.

## Hook-Specific Output Handling

| Hook Event | stdout | stderr (exit 2) |
|------------|--------|-----------------|
| `UserPromptSubmit` | Added to conversation context | Shown to user, blocks prompt |
| `PostToolUse` | Hidden (verbose only) | Shown to Claude |

## Next Steps

- Update vibeguard to exit with code 2 on policy violations to enable blocking behavior
- Consider adding a `--strict` flag to toggle between blocking (exit 2) and warning (exit 1) modes
