---
doc-id: memory
mode: index
location: .spectacular/MEMORY.md
entries-dir: .spectacular/memory/
scope: project-wide
template: templates/memory/entry.md
snapshot-on-edit: false
summary: "Operational memory — index of long-lived notes captured via `spectacular remember`"
status: active
---

# MEMORY Rules

Soft-folder database. The index file (`MEMORY.md`) is regenerated from individual entry files (`memory/<slug>.md`) — never hand-edited as a flat log.

**Mode: `index`** — the canonical content lives in `entries-dir`. The root file is a table-of-contents derived from entry frontmatter. Editing the index directly is a smell; `spectacular doctor memory` will flag drift between index and entries.

**Verbs:**
- `grill` → capture one new entry interactively (skill prompts for text + tags), then CLI writes the entry via `spectacular remember`
- `refine` → asks the user: refine a specific entry by slug, or rewrite a tagged subset
- `review` → validate entry frontmatter shape (type + date + summary required); flag orphaned entries (no row in index) and orphaned index rows (no entry file)

**Mutator verb (CLI, not skill):** `spectacular remember "<text>" [--tag a,b]` writes one entry. Auto-derives slug + summary. Frontmatter `session:` is set automatically when a session is open (see [[sessions-rules]] D9).

**Snapshot-on-edit: false** — entries are themselves append-only; the index regenerates deterministically from them. No snapshot needed.

**Entry frontmatter:**

```yaml
---
type: memory
date: YYYY-MM-DD
tags: [...]
related: [...]
summary: "1-line summary (auto from text if not provided)"
session: <session-slug>|null
---

<body text>
```

**Index shape:** see [[soft-db-substrate/discovery]] § Schema reference. Table columns: Date, Slug, Tags, Summary.

**Doctor area:** `spectacular doctor memory` checks:
- Every entry file has valid frontmatter (`type: memory`, `date`, `summary` present)
- Every index row points to a real entry file
- Every entry file appears in the index
- No two entries share a slug
