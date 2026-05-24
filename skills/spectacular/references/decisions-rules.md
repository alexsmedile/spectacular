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

**Entry format:**

```markdown
## YYYY-MM-DD — Short title

**Context:** What's the situation that forced a decision?
**Decision:** What did we choose?
**Consequences:** What does this enable, foreclose, or imply?
```
