---
doc-id: decisions
mode: index | flat
location: .spectacular/decisions/index.md (flat) | .spectacular/decisions/index.md + .spectacular/decisions/D<N>-<slug>.md (index)
scope: project-wide
template: templates/decisions/entry.md
snapshot-on-edit: false
summary: "ADR-style decision log — flat append or index mode (one-liner root + per-entry files)"
status: active
---

# DECISIONS Rules

Append-only ADR log. Each new decision is one entry, appended to the bottom. Entries are immutable by convention — corrections go in a follow-up entry, not by rewriting history.

> **This is where ADRs (Architecture Decision Records) live in Spectacular.** When asked to "record an ADR" / "record an architecture decision," write here with `spectacular decide` — do **not** create a separate `docs/adr/` folder.

## Is it a DECISIONS.md entry? (store-worthy decisions)

Record a decision here when it is **architectural or product-shaping, has lasting consequences, and a future reader would ask "why was this done this way?"**

**Litmus test (one line):** *would reversing this later need an explanation?* Yes → ADR. No → not store-worthy here.

| ✅ DECISIONS.md (write an ADR) | ❌ Not an ADR — record elsewhere |
|---|---|
| Project-wide choice with lasting trade-offs (e.g. "CLI mutates, skill orchestrates") | **Request-scoped design call** → that request's `PLAN.md` § Decisions |
| A rejected alternative worth remembering ("why not X") | Implementation detail / how-to → code comment or `SPEC.md` |
| A reversal that supersedes an earlier decision | A task or next step → `TASKS.md` |
| A cross-cutting convention other work must respect | A transient session note → `SESSION.md` / memory |
| 1 ADR per breaking change (e.g. the v2 file-contract) | A routine fix or change → `CHANGELOG.md` |
| Something an outside contributor needs to understand the system | Personal/ephemeral preference → don't record |

The key distinction agents get wrong: **project-wide architecture → DECISIONS.md; one-request design decision → that request's PLAN.md § Decisions.** When in doubt, if the decision outlives the request that prompted it, it's an ADR.

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

## Mode detection

**Flat mode** (default): `DECISIONS.md` contains full ADR prose. No `decisions/` folder. Backwards-compatible — all existing workspaces stay in flat mode with zero changes.

**Index mode**: detected by presence of a `decisions/` subfolder next to `DECISIONS.md`. The root file is a cheap one-liner index; full ADR prose lives in per-entry files. The CLI switches to index-mode behavior automatically when it detects the folder.

---

## Index mode (large projects — 50+ decisions)

When `DECISIONS.md` grows past ~50 entries the flat file becomes a context-tax. Use **index mode** instead:

```
.spectacular/
└── decisions/
    ├── index.md          ← index only: one line per decision
    ├── D1-<slug>.md
    ├── D2-<slug>.md
    └── ...
```

**Index line format** (canonical):
```markdown
- **D42** — Reject field-mode storage for v1 — folders-only until v2 ships
```

Three parts separated by ` — `: D-number (bold), short title, one-sentence rationale. No trailing period.

**Per-entry file format** (`decisions/D42-<slug>.md`):
```markdown
# D42 — Reject field-mode storage for v1

**Context:** ...
**Decision:** ...
**Consequences:** ...
**Session:** [[sessions/2026-05-24-foo]]   <!-- optional -->
```

Heading format: `# D<N> — Title` (same title as the index line's short title). Sections follow the ADR schema (Context / Decision / Consequences). Session link is optional.

**Agent read pattern:** always load `DECISIONS.md` (index, cheap — ~1 line per decision). Load `decisions/D<N>-<slug>.md` on demand when that specific decision is directly relevant to current work. Never load all per-entry files at once.

**Migration:** `spectacular decisions migrate` reads flat `DECISIONS.md`, splits each `## YYYY-MM-DD —` block into `decisions/D<N>-<slug>.md`, then rewrites `DECISIONS.md` as the one-liner index. `--dry-run` previews without writing. Idempotent if `decisions/` already exists.

**Detected by:** presence of `decisions/` subfolder next to `DECISIONS.md`. Absence = flat mode (backwards compat).
