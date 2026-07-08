---
doc-id: sessions
mode: index
location: .spectacular/sessions/index.md
entries-dir: .spectacular/sessions/
scope: project-wide
template: templates/sessions/entry.md
snapshot-on-edit: false
summary: "Working sessions — index of `spectacular session start|end` boundaries with auto-linked decisions and memories"
status: active
---

# SESSIONS Rules

> **@SessionEnd policy gate.** On `spectacular session end` / handoff, run `spectacular policy @SessionEnd` and follow every active policy. The default (`summarize-before-handoff`) is `warn`: summarize what changed, what's left, and what's next, then continue. See [policy-injection.md](policy-injection.md).

Soft-folder database. The index file (`SESSIONS.md`) is regenerated from individual entry files (`sessions/<date>-<slug>.md`).

**Mode: `index`** — same pattern as [[memory-rules]]. Canonical content in `entries-dir`; root file is the catalog.

**Verbs:**
- `grill` → polite no-op + hint: "Start a session with `spectacular session start --tag <tag>`"
- `refine` → asks the user: refine session retrospective notes for a closed session, or for the open session
- `review` → validate entry frontmatter shape; flag stale-open sessions (>4h); recompute decision/memory link counts

**Mutator verbs (CLI, not skill):**
- `spectacular session start [--tag a,b]` → create entry with `status: open`, append index row
- `spectacular session end` → flip `status: closed`, set `end_date`, recompute linked-entry counts

**Lifecycle invariant:** at most one session has `status: open` at any time. `spectacular session start` errors if one is already open and suggests `end` first.

**Snapshot-on-edit: false** — entries are factual records of when work happened; immutable by convention.

**Entry frontmatter:**

```yaml
---
type: session
status: open|closed
start_date: YYYY-MM-DDTHH:MM:SS
end_date: YYYY-MM-DDTHH:MM:SS|null
tags: [...]
related: [...]
summary: "1-line summary of session purpose"
decisions_count: 0   # recomputed at session end
memories_count: 0    # recomputed at session end
---

<body — session notes, free-form>

## Linked decisions
- [[decisions/<slug>]]

## Linked memories
- [[memory/<slug>]]
```

**Auto-link mechanic (D9):** when a session is open, `spectacular decide` and `spectacular remember` set `session: <session-slug>` in the new entry's frontmatter. At `session end`, the writer scans `decisions/*` and `memory/*` for matching `session:` fields and:
1. Updates `decisions_count` + `memories_count` in the session entry frontmatter
2. Appends "Linked decisions" / "Linked memories" sections to the session body

**Doctor area:** `spectacular doctor sessions` checks:
- At most one `status: open` session
- Open session age — warn at >4h (D7)
- Every entry file has valid frontmatter
- Every index row points to a real entry file
- For each closed session, `decisions_count` matches actual scan of `decisions/*` with `session: <slug>`
- For each closed session, `memories_count` matches actual scan of `memory/*` with `session: <slug>`
