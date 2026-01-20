---
summary: Research on Claude Code hooks system and identified 6 integration points for vibeguard prompts
event_type: deep dive
sources:
  - https://code.claude.com/docs/en/hooks
  - https://docs.claude.com/en/docs/claude-code/hooks
  - https://www.eesel.ai/blog/hooks-in-claude-code
  - https://github.com/disler/claude-code-hooks-mastery
  - internal/cli/prompt.go
tags:
  - hooks
  - claude-code
  - integration
  - automation
  - prompts
  - context-injection
  - policy-enforcement
  - workflow-enhancement
---

# Claude Code Hooks and Vibeguard Prompts Integration

## Overview

Claude Code hooks are automated event handlers that execute at specific points during Claude Code workflows. They enable policy enforcement, context injection, and workflow automation. Vibeguard prompts can be integrated with hooks to provide intelligent guidance, validation, and context at key workflow moments.

## Claude Code Hooks System Basics

### Available Hook Events

| Hook Event | Trigger | Use Case |
|-----------|---------|----------|
| **SessionStart** | Session begins/resumes | Load context, setup environment |
| **SessionEnd** | Session terminates | Cleanup, logging, final validation |
| **UserPromptSubmit** | User submits message | Validate input, add context |
| **PreToolUse** | Before tool execution | Validate/approve tool calls |
| **PostToolUse** | After tool completes | Validate results, enforce policy |
| **PermissionRequest** | Permission dialog shown | Auto-approve/deny permissions |
| **Notification** | System notification sent | React to alerts |
| **Stop/SubagentStop** | Agent finishes | Block/allow session stop |
| **PreCompact** | Context compaction | Trigger on `/compact` |

### Hook Configuration

**Location**: `.claude/settings.json` (project-level) or `~/.claude/settings.json` (user-level)

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "ToolPattern",  // regex to match tool names
        "hooks": [
          {
            "type": "command",
            "command": "your-command-here",
            "timeout": 60
          }
        ]
      }
    ]
  }
}
```

### Hook Input/Output

**Input**: JSON via stdin containing:
- `session_id`, `transcript_path`, `cwd`
- `hook_event_name`, `tool_name`
- `tool_input` (parameters for the tool)
- `tool_response` (for PostToolUse hooks)

**Output**: JSON response controlling behavior:
```json
{
  "continue": true,
  "systemMessage": "Message for Claude",
  "hookSpecificOutput": {
    "additionalContext": "Context to include"
  }
}
```

**Exit Codes**:
- 0: Success, continue normally
- 2: Blocking error, operation denied

## Integration Opportunity 1: SessionStart Context Injection

### Purpose
Automatically inject vibeguard prompts when Claude Code sessions start, priming Claude with project-specific guidelines.

### Implementation

**Configuration** (`.claude/settings.json`):
```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "vibeguard prompt code-review 2>/dev/null",
            "timeout": 10
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- Inject `code-review` prompt to guide code analysis
- Inject `security-audit` for security-focused sessions
- Inject `init` prompt for onboarding new contributors

### Benefits
- Claude automatically has project guidelines in context
- Consistent coding standards across all sessions
- No manual prompt copying needed

---

## Integration Opportunity 2: UserPromptSubmit Dynamic Context

### Purpose
Add relevant vibeguard prompts based on what the user is asking about.

### Implementation

**Script** (`hooks/context-enhancer.sh`):
```bash
#!/bin/bash
set -e

INPUT=$(cat)
USER_PROMPT=$(echo "$INPUT" | jq -r '.tool_input.prompt // empty')

# Enhance context based on user input keywords
if echo "$USER_PROMPT" | grep -qi "review\|examine\|analyze"; then
  CONTEXT=$(vibeguard prompt code-review 2>/dev/null || echo "")
  if [ -n "$CONTEXT" ]; then
    echo "{\"hookSpecificOutput\":{\"additionalContext\":\"Context for code review:\n\n$CONTEXT\"}}"
    exit 0
  fi
fi

if echo "$USER_PROMPT" | grep -qi "security\|vulnerability\|audit"; then
  CONTEXT=$(vibeguard prompt security-audit 2>/dev/null || echo "")
  if [ -n "$CONTEXT" ]; then
    echo "{\"hookSpecificOutput\":{\"additionalContext\":\"Context for security analysis:\n\n$CONTEXT\"}}"
    exit 0
  fi
fi

if echo "$USER_PROMPT" | grep -qi "test\|unit"; then
  CONTEXT=$(vibeguard prompt test-generator 2>/dev/null || echo "")
  if [ -n "$CONTEXT" ]; then
    echo "{\"hookSpecificOutput\":{\"additionalContext\":\"Context for test generation:\n\n$CONTEXT\"}}"
    exit 0
  fi
fi

# No matching context
exit 0
```

**Configuration**:
```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/hooks/context-enhancer.sh",
            "timeout": 15
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- User says "review this code" → inject code-review prompt
- User says "find security issues" → inject security-audit prompt
- User says "write tests" → inject test-generator prompt

### Benefits
- Intelligent context injection based on user intent
- Reduces need for explicit prompt inclusion
- Seamless workflow enhancement

---

## Integration Opportunity 3: PostToolUse Code Review Guidance

### Purpose
After Claude edits or writes files, inject review prompts to guide self-assessment.

### Implementation

**Script** (`hooks/post-edit-review.sh`):
```bash
#!/bin/bash
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name')
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# After file edits/writes, provide review context
if [[ "$TOOL_NAME" =~ ^(Edit|Write)$ ]] && [[ "$FILE_PATH" == *.go ]]; then
  REVIEW_PROMPT=$(vibeguard prompt code-review 2>/dev/null || echo "")
  if [ -n "$REVIEW_PROMPT" ]; then
    cat <<EOF
{
  "hookSpecificOutput": {
    "additionalContext": "After modifying $FILE_PATH, please review your changes against these guidelines:\n\n$REVIEW_PROMPT"
  }
}
EOF
  fi
fi

exit 0
```

**Configuration**:
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/hooks/post-edit-review.sh",
            "timeout": 15
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- After editing Go files, inject code-review guidance
- After writing test files, inject test-generator prompt
- Self-review reminder with project standards

### Benefits
- Claude self-corrects based on project guidelines
- Consistent code quality in generated changes
- Catches issues before user review

---

## Integration Opportunity 4: PreToolUse Security Validation

### Purpose
Inject security guidance before executing potentially dangerous commands.

### Implementation

**Script** (`hooks/security-gate.sh`):
```bash
#!/bin/bash
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name')
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

# Security validation for shell commands
if [[ "$TOOL_NAME" == "Bash" ]]; then
  # Check for potentially dangerous patterns
  if echo "$COMMAND" | grep -qE "(rm|chmod|chown|curl.*\|.*bash|wget.*\|.*sh|dd if=)"; then
    SECURITY_PROMPT=$(vibeguard prompt security-audit 2>/dev/null || echo "")
    if [ -n "$SECURITY_PROMPT" ]; then
      cat <<EOF
{
  "hookSpecificOutput": {
    "additionalContext": "This command may have security implications. Verify using security guidelines:\n\n$SECURITY_PROMPT"
  }
}
EOF
    fi
  fi
fi

exit 0
```

**Configuration**:
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/hooks/security-gate.sh",
            "timeout": 10
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- Validate shell commands for security risks
- Prevent pipe-to-bash vulnerabilities
- Remind about dangerous file operations

### Benefits
- Security review before command execution
- Prevents accidental destructive operations
- Consistent security awareness

---

## Integration Opportunity 5: File Type Aware Prompt Selection

### Purpose
Automatically select relevant vibeguard prompts based on file type and context.

### Implementation

**Script** (`hooks/smart-prompts.sh`):
```bash
#!/bin/bash
INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# Select prompt based on file type and operation
case "$FILE_PATH" in
  *_test.go|*_test.ts|*_test.py)
    # Test files get test-generator guidance
    vibeguard prompt test-generator 2>/dev/null || true
    ;;
  *.go|*.ts|*.js|*.py)
    # Source files get code-review guidance
    vibeguard prompt code-review 2>/dev/null || true
    ;;
  *.yaml|*.yml)
    # Config files may benefit from init guidance
    if grep -q "prompts:\|checks:" "$FILE_PATH" 2>/dev/null; then
      vibeguard prompt init 2>/dev/null || true
    fi
    ;;
esac

exit 0
```

**Configuration**:
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/hooks/smart-prompts.sh",
            "timeout": 15
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- Test files → test-generator prompt
- Source code → code-review prompt
- Config files → init/setup prompt

### Benefits
- Context-aware guidance without user intervention
- Different prompts for different file types
- Improved workflow efficiency

---

## Integration Opportunity 6: Vibeguard Policy Enforcement Hook

### Purpose
Combine vibeguard policy checks with prompt context for comprehensive validation.

### Implementation

**Script** (`hooks/policy-enforcer.sh`):
```bash
#!/bin/bash
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name')
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# Run vibeguard checks after file modifications
if [[ "$TOOL_NAME" =~ ^(Edit|Write)$ ]]; then
  # Run quick pre-commit checks
  if vibeguard check --tags pre-commit 2>&1 | grep -q "VIOLATION"; then
    # Policy violation - add context
    REVIEW_PROMPT=$(vibeguard prompt code-review 2>/dev/null || echo "")
    echo "{\"hookSpecificOutput\":{\"additionalContext\":\"Policy violation detected. Review guidelines:\n\n$REVIEW_PROMPT\"}}"
    exit 0
  fi
fi

exit 0
```

**Configuration**:
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/hooks/policy-enforcer.sh",
            "timeout": 120
          },
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/vibeguard check --tags pre-commit",
            "timeout": 120
          }
        ]
      }
    ]
  }
}
```

### Use Cases
- Enforce code quality policies before finalization
- Combine policy checks with remediation prompts
- Ensure compliance in AI-generated code

### Benefits
- Automated policy enforcement
- Remediation guidance automatically injected
- Prevents non-compliant code from being committed

---

## Recommended Hook Configuration

Complete hook setup for vibeguard prompt integration:

**File**: `.claude/settings.json`

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "echo 'Project review guidelines:' && vibeguard prompt code-review 2>/dev/null || true",
            "timeout": 10
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/vibeguard check --tags pre-commit",
            "timeout": 120
          }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/hooks/context-enhancer.sh",
            "timeout": 15
          }
        ]
      }
    ]
  }
}
```

---

## Security Considerations

### Best Practices
- Use absolute paths with `$CLAUDE_PROJECT_DIR`
- Quote all variables to prevent shell injection
- Set reasonable timeouts to avoid blocking
- Handle errors gracefully with fallbacks
- Test hooks in isolation before deployment
- Use `2>/dev/null || true` to handle missing commands gracefully

### Security Warnings
- Hooks execute arbitrary shell commands with your user permissions
- Malicious hooks can access or modify any files
- Review all hook scripts before enabling
- Never trust hook configurations from untrusted sources
- Be cautious with dynamic prompt selection based on user input

---

## Implementation Roadmap

### Phase 1: Foundation (Complete)
- ✅ Implement `vibeguard prompt` CLI command
- ✅ Add built-in prompts to vibeguard.yaml
- ✅ Document prompt schema and validation

### Phase 2: Hook Integration
- ⬜ Create hook scripts in `.claude/hooks/`
- ⬜ Enable SessionStart context injection
- ⬜ Enable UserPromptSubmit dynamic context
- ⬜ Test with real Claude Code workflows

### Phase 3: Refinement
- ⬜ Monitor hook execution performance
- ⬜ Gather user feedback on prompt effectiveness
- ⬜ Add additional built-in prompts based on feedback
- ⬜ Optimize hook configurations

### Phase 4: Advanced Integration
- ⬜ Hook-aware prompt customization
- ⬜ Dynamic prompt composition
- ⬜ Integration with policy enforcement
- ⬜ Metrics and analytics on hook usage

---

## Summary

Claude Code hooks provide powerful integration points for vibeguard prompts:

| Integration | Hook Event | Value |
|-------------|-----------|-------|
| **Context Injection** | SessionStart | Prime Claude with guidelines |
| **Dynamic Context** | UserPromptSubmit | Smart context based on intent |
| **Review Guidance** | PostToolUse | Self-review after edits |
| **Security Validation** | PreToolUse | Validate dangerous operations |
| **File-Aware Prompts** | PostToolUse | Context-specific guidance |
| **Policy Enforcement** | PostToolUse | Validate policy compliance |

The `vibeguard prompt` command outputs raw text that can be easily piped into hook scripts, making seamless integration possible.

## Next Steps

1. **Create hook scripts** in `.claude/hooks/` directory
2. **Configure `.claude/settings.json`** with recommended hooks
3. **Test hook execution** with actual workflows
4. **Document hook behavior** for team usage
5. **Gather feedback** on prompt effectiveness
6. **Iterate and refine** based on real-world usage

## Related ADRs

- ADR-005: Adopt VibeGuard for policy enforcement
- ADR-006: Integrate VibeGuard as Git Pre-Commit Hook
- Future: ADR-010: Adopt Claude Code Hooks for Workflow Enhancement
