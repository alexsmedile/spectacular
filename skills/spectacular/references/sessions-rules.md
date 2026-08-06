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
- `spectacular session end` → before closing, print a read-only Git commit review when inside a repository; then flip `status: closed`, set `end_date`, recompute linked-entry counts

**Lifecycle invariant:** at most one session has `status: open` at any time. `spectacular session start` errors if one is already open and suggests `end` first.

**Session boundaries:** before the working briefing, surface `spectacular wayfind status --blockers-only`. At session end, if the session contained a major update, heavy request, release/roadmap change, or architecture decision, re-evaluate the live indexes (`roadmaps/index.md`, `specs/index.md`, and affected collection indexes) and record any deliberate deferral. Routine small edits do not trigger index churn. See [[artifact-retention]].

**Git commit review:** `session end` inspects only the current branch, Git porcelain
status, and staged/unstaged diff statistics before it writes its own session
metadata. It separately reports staged, unstaged, and untracked files; direct
`.spectacular/requests/<slug>/` path matches are ownership hints only. It labels
any follow-up as **"Suggested, human must verify"** because paths and diffs do
not prove whether one change or several commits are appropriate. It may print
ordinary `git diff`, `git status`, `git add -p`, and `git commit` commands, but
never runs Git mutation: no staging, commit, amend, push, merge, reset, or stash.
Outside a Git repository, session closure still succeeds and reports that the
review is unavailable.

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
