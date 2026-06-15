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

**Mutator verb (CLI, v1.5.0+):** `spectacular decide "<decision>"` appends a new entry. The positional argument fills `**Decision:**`. Auto-derives a title slug from the first ~6 words of the decision. If a session is open, the entry includes a `Session:` link to the active session (see [[sessions-rules]] D9). **If `DECISIONS.md` does not exist, the verb bootstraps it** (frontmatter + `# Decisions` heading) before appending — `decide` never fails on a missing file inside a valid workspace.

**Filling Context + Consequences (v1.8.4+):** pass `--context "..."` and `--consequences "..."` to populate those sections at write time:

```
spectacular decide "Use bash for the CLI" \
  --context "team wants zero-install distribution; targets vary" \
  --consequences "ships everywhere with no runtime; harder to unit-test, no static types"
```

Omitted sections are written as **empty headers** (not dropped) so the ADR shape stays visible and fillable by hand or by the skill's `grill`/`refine` flow. The verb never invents Context/Consequences from the decision text — blank means "not yet supplied", not "none".

**Dry run:** `spectacular decide "<text>" --dry-run` previews the entry and writes nothing to disk (v1.8.3+). On a workspace with no `DECISIONS.md` yet, it prints `would create` + `would append` but does **not** bootstrap the file — the bootstrap only happens on a real write.

**Entry format:**

```markdown
## YYYY-MM-DD — Short title

**Context:** What's the situation that forced a decision?
**Decision:** What did we choose?
**Consequences:** What does this enable, foreclose, or imply?
**Session:** [[sessions/2026-05-24-foo]]   <!-- optional, set if session open -->
```

---

## Index mode (large projects — 50+ decisions)

When `DECISIONS.md` grows past ~50 entries the flat file becomes a context-tax. Use **index mode** instead:

```
.spectacular/
├── DECISIONS.md          ← index only: one line per decision
└── decisions/
    ├── D1.md
    ├── D2.md
    └── ...
```

**Index line format:**
```markdown
- **D42** — Reject field-mode storage for v1 — folders-only until v2 ships
```

**Per-decision file format** (`decisions/D42.md`):
```markdown
# D42 — Reject field-mode storage for v1

**Context:** ...
**Decision:** ...
**Consequences:** ...
```

**Agent read pattern:** always load `DECISIONS.md` (index, cheap). Load `decisions/Dxxx.md` on demand when that decision is directly relevant to current work.

**Migration:** read flat DECISIONS.md, extract each `### Dxxx` block into `decisions/Dxxx.md`, rewrite DECISIONS.md as one-liner index. CLI support (`spectacular decisions migrate`) is planned for v1.6.x. Until then, do it manually or write a project-local script.

**Detected by:** presence of a `decisions/` subfolder next to `DECISIONS.md`. Absence = flat mode (backwards compat).
