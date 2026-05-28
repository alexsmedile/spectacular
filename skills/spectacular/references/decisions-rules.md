---
doc-id: decisions
mode: append
location: .spectacular/DECISIONS.md
scope: project-wide
template: templates/decisions/entry.md
snapshot-on-edit: false
summary: "ADR-style decision log — append one entry per decision"
status: active
---

# DECISIONS Rules

Append-only ADR log. Each new decision is one entry, appended to the bottom. Entries are immutable by convention — corrections go in a follow-up entry, not by rewriting history.

**Verbs:**
- `grill` → capture one new entry interactively, append it
  - Prompts for: title, context, decision, consequences
  - Generates ADR-style block from `templates/decisions/entry.md`
- `refine` → asks the user: refine latest entry / all entries / pick one (default: latest)
  - Append-only history is sensitive; refining all is opt-in for a reason
- `review` → validate entry shape (each entry has title + context + decision + consequences)

**Snapshot-on-edit: false** — the file itself is the append log; per-entry snapshots add no value. If a wholesale rewrite is ever needed, the user can snapshot manually first.

**Mutator verb (CLI, v1.5.0+):** `spectacular decide "<text>"` appends a new entry. Auto-derives a title slug from the first ~6 words. If a session is open, the entry includes a `Session:` link to the active session (see [[sessions-rules]] D9). **If `DECISIONS.md` does not exist, the verb bootstraps it** (frontmatter + `# Decisions` heading) before appending — `decide` never fails on a missing file inside a valid workspace.

**Dry run:** `spectacular decide "<text>" --dry-run` previews the entry and writes nothing to disk (v1.8.3+). On a workspace with no `DECISIONS.md` yet, it prints `would create` + `would append` but does **not** bootstrap the file — the bootstrap only happens on a real write.

**Entry format:**

```markdown
## YYYY-MM-DD — Short title

**Context:** What's the situation that forced a decision?
**Decision:** What did we choose?
**Consequences:** What does this enable, foreclose, or imply?
**Session:** [[sessions/2026-05-24-foo]]   <!-- optional, set if session open -->
```

> v1.6.x will introduce optional migration to a soft-folder shape (`decisions/<slug>.md` + `DECISIONS.md` index), unlocking query verbs (`spectacular decisions --7d`). v1.5.0 leaves the flat format intact.
