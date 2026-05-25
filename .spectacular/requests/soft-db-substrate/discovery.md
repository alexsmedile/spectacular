---
status: planned
updated: 2026-05-24
related:
  - PLAN.md
  - TASKS.md
---

# Discovery — soft-db-substrate

Locked decisions from the grill session on 2026-05-24. All seven PLAN-level open questions resolved, plus one follow-on (session linkage).

---

## D1 — Session shape: SESSIONS.md (plural)

**Decision.** Use plural `SESSIONS.md` as the index, all sessions are entries in `sessions/<date>-<slug>.md`. Matches `SPEC.md`/`specs/` and (forthcoming) `DECISIONS.md`/`decisions/`, `MEMORY.md`/`memory/` precedent.

**Implication.** No special "current session" file at root. The active session is identified by frontmatter (`status: open`) in its entry, surfaced by `SESSIONS.md` index and `doctor sessions`.

```
.spectacular/
├── SESSIONS.md              ← index (newest first)
└── sessions/
    ├── 2026-05-20-foo.md    status: closed
    ├── 2026-05-22-bar.md    status: closed
    └── 2026-05-24-baz.md    status: open   ← current
```

**Why this matters.** Perfect symmetry across all three new doc-types means one rules-file template, one doctor area shape, one mental model.

---

## D2 — Slug generation: auto + warn on collision

**Decision.** `spectacular decide "use bash for cli"` → `decisions/use-bash-for-cli.md`. Auto-derived from first ~6 words, slugified. On collision, append `-2`, `-3`, etc. and print a warning. Never silently overwrite. Never block.

**Implication.** Shared slug util reused from request scaffolder (`cli/spectacular` already has this). `spectacular decide` is single-arg.

> Note: this lands closer to the "Auto + warn on collision" preview than strict "Auto only" — collision handling is implicit in any auto-slug system, so calling it out explicitly here.

---

## D3 — Migration timing: wait for v1.6.x

**Decision.** v1.5.0 ships with:

- **MEMORY.md + memory/** — new, folder-shape from day one
- **SESSIONS.md + sessions/** — new, folder-shape from day one
- **DECISIONS.md** — stays flat (existing format), unchanged

`spectacular doctor decisions --fix` (migration to folder shape) is deferred to v1.6.x, where the query verb `spectacular decisions --7d` makes the folder shape's value obvious.

**Implication.** M6 (migration helper) is dropped from this request. Move it to v1.6.x scope.

---

## D4 — Init kit defaults: coding kit only

**Decision.**

- **Always-set:** unchanged (6 files — PRD, PLAN, TASKS, ROADMAP, AGENTS, README)
- **Coding kit:** add `MEMORY.md` + `SESSIONS.md` to existing `DECISIONS.md`, `STACK.md`, `ARCHITECTURE.md`
- **Minimal/blank kits:** stay lean, no new files

**Implication.** Kit registry update in `init-workflow.md` references. The three capture surfaces ship together — coding projects get the full substrate as one bundle.

---

## D5 — Auto-tagging on `remember`: explicit --tag only

**Decision.** `spectacular remember "<text>" [--tag a,b,c]`. No auto-derivation from cwd, active request, or time. User passes tags or skips them.

**Implication.** Frontmatter `tags` is `[]` by default. RAG layer (future) will derive contextual signal from `related` and `summary` fields instead of tag soup.

```yaml
---
type: memory
date: 2026-05-24
tags: [cli, regression]      # explicit, can be empty
related: []                  # explicit, can be empty
summary: "1-line summary"    # filled by CLI from text if not provided
session: null                # auto-set if session open (see D8)
---
```

---

## D6 — Session lifecycle: explicit start/end only

**Decision.** `spectacular session start [--tag ...]` opens. `spectacular session end` closes. No auto-open on other verbs. No idle timeout.

**Implication.** Users will sometimes forget. That's why D7 lands a 4h doctor warning — recovery surface, not enforcement.

```bash
$ spectacular session start --tag substrate-work
  → sessions/2026-05-24-substrate-work.md (status: open)
  → SESSIONS.md updated

$ spectacular session end
  → status: open → closed
  → end_date set
  → SESSIONS.md updated with linked decisions/memories
```

---

## D7 — Stale-session warning: 4h

**Decision.** `spectacular doctor sessions` warns on any `status: open` session older than 4h.

**Implication.** Tight enough to catch "forgot to end" same-day. Users who want longer sessions can `--ignore-stale` (TBD) or just close + reopen.

---

## D8 — RAG-ready frontmatter: type + summary

**Decision.** Entry frontmatter schema across all three doc-types:

```yaml
---
type: decision | memory | session-note
date: YYYY-MM-DD
tags: [...]                 # explicit only (per D5)
related: [...]              # wikilinks/paths
summary: "1-line summary"   # cheap RAG embedding target
session: <session-slug>|null # auto-set during open session (per D9)
---
```

**Dropped** from the candidate set: `author/agent` (deferred to v1.7.x multi-agent work), `context_ref` (folded into `related` — no need for a separate field today).

**Why `summary`.** A future embedding layer can cheaply scan summaries first, then load full body only for top-k matches. The CLI auto-fills from the first sentence if user doesn't pass `--summary`.

**Why `type`.** Discriminator for cross-doc search (`recall --type decision`).

---

## D9 — Session linkage: auto-link when session open

**Decision.** When `spectacular decide` or `spectacular remember` runs while a session is open, the entry's frontmatter gets `session: <slug>` automatically. When the session is closed, `SESSIONS.md` entry summary auto-lists the count of decisions and memories captured during it.

**Implication.**

- Free retrospective: each closed session shows what was decided/remembered during it
- `recall --session <slug>` (v1.6.x) becomes a trivial query
- One inverse-link computation (doctor sessions reads decisions/* and memory/* for matching `session:` field)

**Edge case.** Decisions/memories made outside any session have `session: null` — that's fine, they're just "ambient" captures.

---

## Schema reference — index files

### `SESSIONS.md`

```markdown
---
type: index
doc: sessions
updated: 2026-05-24
---

# Sessions

| Date       | Slug            | Status | Decisions | Memories |
|------------|-----------------|--------|-----------|----------|
| 2026-05-24 | substrate-work  | open   | 0         | 0        |
| 2026-05-22 | bar             | closed | 3         | 1        |
```

### `MEMORY.md`

```markdown
---
type: index
doc: memory
updated: 2026-05-24
---

# Memory

| Date       | Slug         | Tags          | Summary                              |
|------------|--------------|---------------|--------------------------------------|
| 2026-05-24 | use-haiku    | [cli, perf]   | Haiku is fast enough for slug gen    |
```

### `DECISIONS.md` (unchanged in v1.5.0 per D3)

Keep flat ADR-style format. Migration in v1.6.x.

---

## Out of scope (confirmed deferred)

- DECISIONS.md → decisions/ migration → **v1.6.x**
- Query verbs (`decisions --7d`, `recall`, `sessions`) → **v1.6.x**
- `author/agent` provenance → **v1.7.x** (multi-agent advisory)
- Embedding / RAG implementation → **future** (schema is ready)
- UI/browser surface → **never in CLI**

## Updated milestone impact

- **M6 (migration helper)** dropped from this request → moves to v1.6.x
- **M3 CLI writers** gain: auto session linkage logic (small)
- **M5 doctor** gains: 4h stale-session check + inverse-link computation (small)
- **M4 index files** simplifies: DECISIONS.md untouched, only MEMORY.md + SESSIONS.md bootstrapped
