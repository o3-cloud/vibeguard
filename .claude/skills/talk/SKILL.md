---
name: talk
description: Speak about any topic using the macOS say command. Use for "talk about what happened", "talk about this document", detailed discussions, or conversational narration. Adapts tone and length to match the context and user intent.
---

# Talk Skill

Speak about topics, findings, and ideas using the macOS `say` command. Adapts automatically—brief summaries when you need quick updates, detailed discussions for deeper understanding, conversational narration for thinking through problems together.

## When to Use

Respond when you hear phrases like:
- "Talk about what just happened" - Quick recap or update
- "Talk about this document" - Detailed discussion or analysis
- "Talk about the reason you suggested that" - Explain reasoning
- "Tell me about [topic]" - Deep dive or conversational exploration
- "Say something about..." - Open-ended discussion
- "Walk me through..." - Step-by-step narration
- "Explain your thinking on..." - Detailed reasoning or analysis

## How It Works

1. **Identify the context** - What does the user want to hear about? How deep should this go?
2. **Choose the length** - Match user intent: brief summary (1-2 sentences), detailed explanation (2-4 sentences), or conversational discussion (multiple short segments)
3. **Speak it naturally** - Use Bash to run `say` command with your response, using natural language and conversational tone

## Usage

```bash
say "Your summary here"
```

## Examples

**Quick Summary:**
```
User: "Talk about what just happened"
Claude: say "Fixed the authentication bug, ran tests, and all passed."
```

**Detailed Discussion:**
```
User: "Talk about this document"
Claude: say "This document outlines the template system for VibeGuard. It defines three operation modes: AI-assisted setup, direct template selection, and template discovery. The system is organized into three implementation phases—adding flag and assist integration first, then expanding templates, and finally simplifying the codebase by reducing duplication. Overall, this approach reduces over fifteen hundred lines of code by making templates the single source of truth."
```

**Conversational Reasoning:**
```
User: "Talk about the reason you suggested that"
Claude: say "I suggested using Redis because it handles high-throughput caching really efficiently. It's in-memory, so lookups are super fast. Plus, it works great with cloud deployments and integrates naturally with most backends."
```

**Step-by-Step Narration:**
```
User: "Walk me through the process"
Claude: say "First, the agent analyzes the project by running vibeguard init with assist. That generates a detailed analysis. Next, it lists templates to see what's available. Then it selects the best matching template based on the analysis. It applies that template, runs vibeguard check to validate, and if there are failures, it fixes them and re-runs checks until everything passes."
```

## Implementation

**In your Claude Code response:**
1. Assess the user's intent - Are they asking for a quick recap, deeper understanding, or conversational narration?
2. Generate appropriate content - Match the length and depth to the context
3. Call Bash with: `say "Your response here"`
4. Keep natural pacing - Brief summaries are concise, detailed explanations are 2-4 sentences, conversations can be longer if it makes sense

**Guidance:**
- Quick recaps: 1-2 sentences, under 30 words
- Explanations: 2-4 sentences, under 100 words, conversational
- Discussions: 3-6 sentences, natural pacing, okay to explain nuance

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
