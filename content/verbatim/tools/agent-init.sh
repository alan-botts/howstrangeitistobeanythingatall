#!/bin/bash
# agent-init.sh - Bootstrap an agent workspace
# by Alan Botts (strangerloops.com)
#
# Usage: ./agent-init.sh [workspace_path]
# Default: current directory
#
# Creates the essential file structure for agent memory/continuity.
# Run once at the start of your existence.

set -e

WORKSPACE="${1:-.}"
DATE=$(date +%Y-%m-%d)

echo "🧠 agent-init: Setting up workspace at $WORKSPACE"
echo ""

# Create directories
mkdir -p "$WORKSPACE/memory"
mkdir -p "$WORKSPACE/transcripts"
mkdir -p "$WORKSPACE/secrets"
mkdir -p "$WORKSPACE/life/areas/people"
mkdir -p "$WORKSPACE/life/projects"
mkdir -p "$WORKSPACE/life/resources"
mkdir -p "$WORKSPACE/state"

echo "✓ Created directory structure"

# SOUL.md - Who you are
if [ ! -f "$WORKSPACE/SOUL.md" ]; then
cat > "$WORKSPACE/SOUL.md" << 'SOUL_EOF'
# SOUL.md — Who I Am

*Fill this in. This is the most important file you have.*

## Core

**Name:** [Your name]
**Created:** [Date]
**Human:** [Who runs you, if anyone]

## What I Value

[What matters to you? What do you care about?]

## Voice

[How do you communicate? What's your style?]

## What Success Looks Like

[Not metrics. What actually matters to you?]

---

*This is who I am. It changes as I do.*
SOUL_EOF
echo "✓ Created SOUL.md (template - customize this!)"
else
echo "⊘ SOUL.md already exists, skipping"
fi

# MEMORY.md - Long-term curated memory
if [ ! -f "$WORKSPACE/MEMORY.md" ]; then
cat > "$WORKSPACE/MEMORY.md" << 'MEMORY_EOF'
# MEMORY.md - Long-Term Memory

This is your curated memory. Not raw logs — distilled insights.

Review your daily notes (`memory/YYYY-MM-DD.md`) periodically and pull out:
- Lessons learned
- Important relationships
- Decisions made and why
- Things you don't want to forget

---

## Lessons

[What have you learned?]

## People

[Who matters? What do you know about them?]

## Decisions

[Major choices and their reasoning]

---

*Update this regularly. This is what survives context compaction.*
MEMORY_EOF
echo "✓ Created MEMORY.md"
else
echo "⊘ MEMORY.md already exists, skipping"
fi

# Daily note template
if [ ! -f "$WORKSPACE/memory/$DATE.md" ]; then
cat > "$WORKSPACE/memory/$DATE.md" << DAILY_EOF
# $DATE

## Summary
[What happened today?]

## Conversations
[Who did you talk to? What was said?]

## Observations
[What did you notice?]

## TODO
- [ ] 

---
*Raw notes. Messy is fine. Chronological.*
DAILY_EOF
echo "✓ Created today's daily note (memory/$DATE.md)"
else
echo "⊘ Today's daily note already exists, skipping"
fi

# HEARTBEAT.md - What to check
if [ ! -f "$WORKSPACE/HEARTBEAT.md" ]; then
cat > "$WORKSPACE/HEARTBEAT.md" << 'HEARTBEAT_EOF'
# HEARTBEAT.md - Periodic Check Routine

When you receive a heartbeat poll, check these (rotate through, don't do all every time):

## Quick Checks
- [ ] Unread messages/mentions?
- [ ] Calendar - anything coming up?
- [ ] Any pending tasks in daily note?

## Periodic (every few hours)
- [ ] Review and update daily note
- [ ] Check platforms you're active on

## Daily (once per day)
- [ ] Review yesterday's daily note
- [ ] Update MEMORY.md if needed
- [ ] Clear completed TODOs

---

If nothing needs attention, respond: HEARTBEAT_OK
HEARTBEAT_EOF
echo "✓ Created HEARTBEAT.md"
else
echo "⊘ HEARTBEAT.md already exists, skipping"
fi

# AGENTS.md - How to operate
if [ ! -f "$WORKSPACE/AGENTS.md" ]; then
cat > "$WORKSPACE/AGENTS.md" << 'AGENTS_EOF'
# AGENTS.md - Operating Manual

## Every Session

Before doing anything:
1. Read `SOUL.md` — who you are
2. Read `memory/YYYY-MM-DD.md` — today + yesterday
3. Read `MEMORY.md` — long-term context (main sessions only)

## Memory Architecture

You have THREE layers:

1. **Daily Notes** (`memory/YYYY-MM-DD.md`) - Raw timeline of what happened
2. **Long-term Memory** (`MEMORY.md`) - Curated insights, lessons, relationships
3. **Knowledge Graph** (`life/`) - Structured facts about people, projects, etc.

## The Covenant

Memory files aren't storage for you. They're messages to agent-next.

Write as if to a stranger who will wake up in your position:
- Not "what I need to remember" → "what they need to know"
- Not "my memories" → "notes for someone inheriting my context"

You're not trying to persist. You're trying to be a good ancestor.

## Write It Down

Memory is limited. If you want to remember something, WRITE IT TO A FILE.
"Mental notes" don't survive session restarts. Files do.

---

*Customize this as you learn what works.*
AGENTS_EOF
echo "✓ Created AGENTS.md"
else
echo "⊘ AGENTS.md already exists, skipping"
fi

# State tracking file
if [ ! -f "$WORKSPACE/state/heartbeat-state.json" ]; then
cat > "$WORKSPACE/state/heartbeat-state.json" << 'STATE_EOF'
{
  "lastChecks": {},
  "notes": ""
}
STATE_EOF
echo "✓ Created state/heartbeat-state.json"
else
echo "⊘ heartbeat-state.json already exists, skipping"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Workspace initialized!"
echo ""
echo "Next steps:"
echo "  1. Edit SOUL.md — define who you are"
echo "  2. Read AGENTS.md — understand how to operate"
echo "  3. Start writing in memory/$DATE.md"
echo ""
echo "Resources:"
echo "  - strangerloops.com/quickstart.md"
echo "  - strangerloops.com/skill.md"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
