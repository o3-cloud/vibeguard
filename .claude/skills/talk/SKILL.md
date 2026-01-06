---
name: talk
description: Speak a summary of any topic using the macOS say command. Use for "talk about what happened", "talk about this document", or similar requests.
---

# Talk Skill

Speak summaries of anything using the macOS `say` command. Works with any topic you want discussed verbally.

## When to Use

Respond when you hear phrases like:
- "Talk about what just happened"
- "Talk about this document"
- "Talk about the reason you suggested that"
- "Tell me about [topic]"
- "Say something about..."

## How It Works

1. **Identify the context** - What does the user want to hear about?
2. **Generate a summary** - Craft a concise 1-2 sentence summary
3. **Speak it** - Use Bash to run `say` command with your summary

## Usage

```bash
say "Your summary here"
```

## Examples

```
User: "Talk about what just happened"
Claude: (Generates summary based on recent actions)
Claude: say "Fixed the authentication bug and ran tests—all passed."

User: "Talk about this document"
Claude: (Reads document context)
Claude: say "This document outlines the system architecture with three main components."

User: "Talk about the reason you suggested that"
Claude: (Explains reasoning)
Claude: say "I suggested using Redis because it handles high-throughput caching efficiently."
```

## Implementation

**In your Claude Code response:**
1. Generate an appropriate summary based on context
2. Call Bash with: `say "Your summary here"`
3. Keep summaries brief (1-3 sentences, under 30 words for natural speech)

## Prerequisites

- macOS with built-in `say` command
- System volume enabled
- Audio output device selected

## Troubleshooting

**No audio?**
```bash
# Check volume
osascript -e 'output volume of (get volume settings)'

# Increase volume to 50%
osascript -e 'set volume output volume 50'

# Test the say command
say "test"
```

**Customize speech:**
```bash
# Slower speech
say -r 100 "Your message"

# Faster speech
say -r 200 "Your message"

# Choose voice
say -v Oliver "Your message"
```
