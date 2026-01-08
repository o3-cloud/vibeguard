# Meeting Skill

Conducts structured meetings with spoken commentary using the talk skill. Reads meeting agendas from the `docs/meetings/` folder and guides the discussion with audio narration throughout.

## Usage

```bash
/meeting                          # Auto-detect latest meeting agenda
/meeting vibeguard-spec-review    # Run specific meeting by name (without .md)
```

## How It Works

1. **Agenda Discovery** - Scans `docs/meetings/` for meeting files
2. **Agenda Reading** - Loads the meeting agenda (YAML frontmatter + content)
3. **Meeting Flow** - Talks through agenda sections, facilitates discussion
4. **Section Narration** - Uses talk skill to:
   - Announce each agenda section
   - Speak key discussion points
   - Summarize decisions made
   - Introduce next section

## Meeting Agenda Format

Create meeting files in `docs/meetings/` with this structure:

```markdown
---
title: Meeting Title
date: 2026-01-08
participants: [Person1, Person2]
status: scheduled|completed|cancelled
---

# Meeting Title

## Section 1: Overview
[Content about section 1]

## Section 2: Key Discussion Points
[Content about section 2]

## Section 3: Decisions & Actions
[Content about section 3]
```

## Features

- **Auto-discovery** - Finds latest meeting if not specified
- **Spoken narration** - Uses macOS `say` command for audio
- **Agenda-driven** - Follows the meeting structure you define
- **Summary support** - Can create meeting notes or summaries
- **Integration** - Works with other tools (TodoWrite, Skill, etc.)

## Examples

### Run latest meeting with spoken commentary:
```bash
/meeting
```

### Run specific meeting:
```bash
/meeting project-kickoff
```

### Create and run meeting in sequence:
```bash
# First create the agenda file
# Then run it with spoken narration
/meeting my-meeting-name
```

## Technical Details

The meeting skill:
1. Reads meeting agenda from `docs/meetings/<name>.md`
2. Parses YAML frontmatter (title, date, participants, status)
3. Extracts markdown sections from the meeting file
4. For each section:
   - Announces section heading aloud
   - Speaks key discussion points
   - Pauses for live discussion
   - Summarizes outcomes
5. Uses the talk skill (`/talk`) for all audio narration

## Integration with Other Skills

- **`/talk`** - Provides the audio narration capability
- **`/log`** - Create log entries to document meeting findings
- **`/adr`** - Create ADRs from meeting decisions
- **`TodoWrite`** - Track action items from meeting

## Tips

- Write meeting agendas in `docs/meetings/` before running the skill
- Use clear section headings for better narration
- Include action items with owners and due dates
- Create follow-up meeting files for recurring meetings
- Use the talk skill's tone control for emphasis on important points
