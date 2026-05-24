---
status: locked
locked: 2026-05-24
updated: 2026-05-24
related:
  - PLAN.md
---

# Discovery — substrate-clarity

**Status: LOCKED 2026-05-24.** Decisions 1-5b are agreed and binding. Decisions 6-7 deferred to M2 (spec-refine) or a future grill. PLAN.md M2 onward translates this into concrete file changes.

## Decision 1+2 — Grill is a verb and a mode family

**`grill` is the verb.** The interactive-form action. The skill always means "the agent interviews the user."

**The mode value names the grill style directly.** No separate sub-mode field. The `grill-*` family clusters visually:

| Mode value | Style | Use case |
|---|---|---|
| `grill` | Alias for `grill-wide` (default for the family) | Most docs |
| `grill-wide` | One broad session covering all slots at once | Fast first pass; user has clarity upfront |
| `grill-each` | Checkpoint per block/chapter — granular, one block at a time | Each block deserves focused attention (ROADMAP per-version, PERSONAS per-person) |
| `grill-loop` | Wide pass first → then narrow on remaining open questions | Hybrid: get the shape fast, then deepen where it's vague |

**Registry usage:**

```yaml
prd:
  mode: grill          # = grill-wide

roadmap:
  mode: grill-each

personas:
  mode: grill-each
```

**CLI/skill flag override** (user wins over declared mode):
- `spectacular roadmap grill` — uses doc's declared mode (here, `grill-each`)
- `spectacular roadmap grill --wide` — force `grill-wide` for this session
- `spectacular roadmap grill --each` — force `grill-each`
- `spectacular roadmap grill --loop` — force `grill-loop`

**Top-level mode taxonomy** (5 conceptual modes, 8 values because grill has 3 named variants + the shorthand):

- `grill` / `grill-wide` / `grill-each` / `grill-loop` — interactive-form family
- `append` — capture one entry, append
- `stub` — scaffold + exit
- `freeform` — agent improvises (reserved)
- `reference` — engine-internal, not user-facing

The previous `reps` mode collapses into `grill-each`. No separate `sub-mode` / `grill-style` field — the mode value carries it.

**v1.4.0 ships all three grill variants.** `grill-loop` is a genuinely useful pattern (wide-then-deep matches real thinking) so it lands now, not in a follow-up.

## Decision 3 — Three orthogonal axes

Never collapse:

1. **Phase** — lifecycle state (discover / spec-refine / mvp / iterate / test / release-prep / release). Lives in PLAN.md frontmatter and ROADMAP version blocks.
2. **Verb** — what the agent does right now (grill / refine / review / new / archive). Invoked via CLI or skill.
3. **Mode** — doc interaction style (grill / append / stub / freeform / reference). Property of the doc, declared in registry.

They correlate but never map 1:1. You can `refine` a PRD during the `mvp` phase if you realize the spec drifted.

## Constraint — Grill (and refine) are skill-only, never CLI

**Surfaced during Decision 5 grilling.** Locks the verb taxonomy at a structural level:

| Verb | Runs in |
|---|---|
| `grill` | **Skill only.** Requires LLM to reason about answers, ask follow-ups, mini-refine, gate quality. |
| `refine` | **Skill only.** Vibe→spec rewrite is LLM work. |
| `review` | **Mixed.** Structural checks (frontmatter, sections present, no placeholders) run in CLI; semantic checks (vague words, slot quality) need LLM. |
| `new` / `archive` / `snapshot` / `init` / `doctor` | **CLI primarily.** Mechanical scaffolding, file moves, integrity checks. |

**Honest CLI behavior:** if a user types `spectacular roadmap grill` at terminal, the CLI prints:

> "Grill requires an agent. Run this in Claude Code or Codex via `/spectacular roadmap grill`."

Not silent failure. Not partial scaffolding. Just a clean redirect.

**Implication for G1 (registry):** the CLI never needed full dispatch — agentic fields (slots, sub-modes, rules body) were never CLI-readable concerns anyway. The CLI legitimately needs only: catalog (doc-ids, descriptions, locations) + mechanical dispatch (mode, template, scope, snapshot-on-edit). Agentic dispatch (slots, sub-modes, rules) belongs to the agent.

This three-way split:

| What | Lives where | Read by |
|---|---|---|
| **Catalog** (doc-id, description, location) | `doc-index.md` | Humans + CLI |
| **Mechanical dispatch** (mode, template, scope, snapshot-on-edit) | Rules file **frontmatter** | Agent + CLI |
| **Agentic dispatch** (slots, sub-modes, rules body) | Rules file **body** | Agent only |

## Decision 4 — Verb × Mode matrix

| | **grill** | **append** | **stub** | **freeform** | **reference** |
|---|---|---|---|---|---|
| **grill** (verb) | Run interactive form (sub-mode resolves wide/each/loop) | Capture one new entry interactively, append | **Polite no-op + hint** — "stub doc, open in editor; or pass `--wide` to grill it ad-hoc for this session" | **Open-ended prompt** — "what do you want to capture?" Agent infers slot list on the fly and walks it | Error: not user-facing |
| **refine** (verb) | Vibe→spec rewrite across slots | **Ask user**: "Refine latest / all / pick entry?" | Whole-doc refine pass | Whole-doc refine pass | N/A |
| **review** (verb) | Slot-based gate check (per-block if sub-mode=each) | Validate entry shape | Structural review (frontmatter + sections + no template placeholders left) | Subjective review ("does this read well?") | Internal validation |

**Notable resolutions:**
- `grill × stub` → friendly hint, never a dead-end. Optional ad-hoc grill via `--wide`.
- `grill × freeform` → onramp prompt; agent generates a slot list from the user's first answer, then walks it. Honors freeform's spirit without leaving the user staring at a blank canvas.
- `refine × append` → user picks scope (latest / all / pick). Don't default to either extreme — append docs are sensitive history.

## Decision 5 — Doc-registry fate

**Rename** `doc-registry.md` → `doc-index.md`.

**Reframe** as human catalog. Drop the dispatch-contract framing ("no code changes required" — false claim, removed).

**Move dispatch fields** into each rules file's frontmatter. Engine reads rules files directly. Index becomes pointer + description.

**What stays in `doc-index.md`:**
- Doc-id → one-line description
- Mode (high-level, for scanning)
- Location
- Link to the rules file

**What moves to rules-file frontmatter:**
- `template`, `mode`, `default-sub-mode`, `slots`, `location`, `scope`, `snapshot-on-edit`, `kit-support`, `doc-id`

### Decision 5b — Docs without rules files

**Every doc gets a rules file**, even minimal stubs (frontmatter only, no body).
Affected docs: PRINCIPLES, ARCHITECTURE, STACK, AGENTS, SPEC, DECISIONS.

Each becomes a 5-10 line file: frontmatter declaring dispatch, optional one-paragraph note explaining when the doc is used. Consistency over brevity.

This produces uniform engine behavior: every doc-id resolves to a rules file. No conditional "if rules file exists" branches. Adding a new doc = create rules file + template (+ optional body content in rules file). The index entry is a downstream artifact (could even be auto-generated from rules-file frontmatter — out of scope for v1.4.0).

## Decision 6 — Rules-file frontmatter schema

Every doc gets a rules file. Frontmatter is the dispatch contract — agent + CLI both read it.

### Required fields (4)

| Field | Purpose |
|---|---|
| `doc-id` | Unique identifier, kebab-case. Engine looks up rules by doc-id. |
| `mode` | Determines engine behavior. One of: `grill / grill-wide / grill-each / grill-loop / append / stub / freeform / reference`. |
| `location` | Where the doc lives. Supports `<slug>` interpolation for per-request. |
| `scope` | `project-wide` / `per-request` / `user` / `skill-internal`. |

### Mode-conditional fields

| Field | Required for | Forbidden for |
|---|---|---|
| `template` | grill-family, append, stub | reference |
| `slots` | grill-family | append, stub, freeform, reference |
| `kit-support` | (only `true` permitted for grill-family) | non-grill |
| `snapshot-on-edit` | — (optional, default `false`; typically `true` for project-wide canonical docs) | — |

**Schema is strict** — `spectacular doctor frontmatter` flags violations (e.g. `slots` on a stub doc).

### Optional fields (3)

| Field | Purpose |
|---|---|
| `summary` | One-line description. Used by `doc-index.md` and `--help`. |
| `version` | Schema version of the rules file itself (for future migrations). |
| `status` | Lifecycle: `active` (default) / `experimental` / `deprecated`. Replaces a separate deprecation boolean. |

### Dropped from earlier proposals
- `related`, `since`, `deprecated` (separate field), `tags`, `default-sub-mode`, `grill-style`. None earn their slot.

### Examples

**A — grill (default = wide), PRD:**
```yaml
---
doc-id: prd
mode: grill
location: .spectacular/PRD.md
scope: project-wide
template: templates/prd/base.md
slots: [Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone]
snapshot-on-edit: true
kit-support: true
summary: "Product Requirements Document — what & why & for whom (8 slots)"
status: active
---
```

**B — grill-each, ROADMAP:**
```yaml
---
doc-id: roadmap
mode: grill-each
location: .spectacular/ROADMAP.md
scope: project-wide
template: templates/roadmap/base.md
slots: [Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests]
snapshot-on-edit: true
summary: "Per-version scope + phase + exit criteria"
status: active
---
```

**C — append, DECISIONS:**
```yaml
---
doc-id: decisions
mode: append
location: .spectacular/DECISIONS.md
scope: project-wide
template: templates/decisions/entry.md
summary: "ADR-style decision log"
status: active
---
```

**D — stub, AGENTS:**
```yaml
---
doc-id: agents
mode: stub
location: .spectacular/AGENTS.md
scope: project-wide
template: templates/agents/base.md
snapshot-on-edit: true
summary: "Onboarding doc for agents working in .spectacular/"
status: active
---
```

**E — reference, migrations-contract:**
```yaml
---
doc-id: migrations-contract
mode: reference
location: skills/spectacular/references/migrations-contract.md
scope: skill-internal
snapshot-on-edit: true
summary: "Schema contract for workspace-schema migration files"
status: active
---
```

**F — deprecated, docs-page:**
```yaml
---
doc-id: docs-page
mode: stub
location: docs/<section>/<slug>.md
scope: per-request
template: templates/docs/page.md.tmpl
summary: "Single user-facing docs page"
status: deprecated
---
```

## Decision 7 — Drop "engine"

**Drop the word entirely.** No replacement metaphor. The verb × mode matrix made "the engine" unnecessary as a concept — there isn't really one unified thing, just verb-specific behaviors that dispatch on mode.

**Replacement guide** (M7 cleanup sweep, ~78 occurrences across `skills/spectacular/`):

| Old | New |
|---|---|
| "the generic engine" | "grill / refine / review" or "the skill" |
| "the engine reads" | "the verb reads" or "the skill reads" |
| "engine rules" | "rules" |
| "consumed by the engine" | "consumed by grill/refine/review" — or drop the phrase |
| "engine behavior" | "skill behavior" or "verb behavior" |

When a collective noun is genuinely needed, **"the skill"** carries it. When specificity matters, name the verb directly. This produces sentences that are tighter and more honest about what's happening.

## All decisions locked

| # | Decision |
|---|---|
| 1+2 | Grill is a verb + mode family. Values: `grill` / `grill-wide` / `grill-each` / `grill-loop` |
| 3 | Three orthogonal axes — Phase, Verb, Mode — never collapse |
| Constraint | Grill + refine are skill-only; CLI redirects |
| 4 | Verb × Mode matrix defined for every cell |
| 5 | `doc-registry.md` → `doc-index.md` (human catalog) |
| 5b | Every doc gets a rules file (frontmatter-only is fine for stubs) |
| 6 | Rules-file frontmatter schema: 4 required + mode-conditional + 3 optional |
| 7 | Drop "engine" — use "skill" or name the verb directly |

M1 + M2 (discovery + decision-locking) complete. Ready for M3 onwards.
