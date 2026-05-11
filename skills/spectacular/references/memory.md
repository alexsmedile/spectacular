# Memory — Writing Operational Learning

Triggered by: `spectacular remember this`, or as part of the archive sequence.

---

## Memory location

`.spectacular/memory/` — git-committed, team-visible. This is NOT `.claude/` personal memory.

Default files:
```
memory/
├── failures.md
├── lessons.md
├── architecture-traps.md
└── recurring-bugs.md
```

Create new files if the content doesn't fit existing categories. Keep files focused.

---

## Anti-collision rule

**IMPORTANT:** Avoid phrasing that triggers Claude Code's own auto-memory system.

- Do NOT use phrases like "remember that", "note that", "I should remember"
- Write memory entries as factual operational records, not instructions to an agent
- Target audience is the team reading `.spectacular/memory/`, not the AI assistant

---

## Write triggers

### On demand: `spectacular remember this`

1. Ask what specifically to capture if not clear from context
2. Propose which memory file it belongs in
3. Draft the entry
4. Write immediately on user confirmation (no second confirmation needed)

### On archive

As part of the archive sequence, propose memory entries for:
- Blockers that weren't obvious
- Architecture decisions made mid-implementation
- Bugs or failures discovered
- Implementation patterns worth reusing
- Risks that materialized

---

## Entry format

Keep entries concise and specific. Include context so the entry is useful months later.

```md
## <Short title> — <date>

<What happened or was learned>

**Context:** <Which request / what system / what conditions>
**Impact:** <Why this matters>
```

Example — `failures.md`:
```md
## Stripe webhook deduplication — 2026-05-11

Webhooks fired twice on network retry; our handler wasn't idempotent.
Added `stripe_event_id` deduplication table.

**Context:** add-team-billing request
**Impact:** Without idempotency, billing events double-charged on Stripe retry
```

---

## Reading memory

During normal operation, skill reads `.spectacular/memory/` file list and counts for the status briefing.
Load full memory content only when directly relevant to the current task.
