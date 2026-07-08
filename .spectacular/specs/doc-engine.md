---
status: draft
updated: 2026-05-29
summary: "Registry-driven engine that grills/refines/reviews any registered canonical doc through one set of skill flows"
related:
  - index.md
  - ../ARCHITECTURE.md
---

# Doc-writing engine

> Promoted from `SPEC.md` in v1.10.0 (was a single dense bullet). The index keeps a one-line pointer; the full contract lives here.

## Purpose

One engine operates *every* canonical doc. Instead of a bespoke flow per document, three generic skill flows — `grill.md`, `refine.md`, `review.md` — handle any registered doc verb (`spectacular <doc> grill|refine|review`). Per-doc behavior is data, not code: each doc carries a `references/<doc-id>-rules.md` file whose **frontmatter** declares dispatch and whose **body** declares prompts and gate checks. This mirrors Spectacular's own philosophy — small files, progressive context loading, behavior driven by the registry rather than hardcoded.

## Dispatch contract

Each `references/<doc-id>-rules.md` frontmatter declares:

| Field | Meaning |
|---|---|
| `mode` | Which engine behavior to run (see Modes below) |
| `slots` | The structured fields the doc is decomposed into |
| `template` | Path to the scaffold template |
| `location` | Where the doc lives in `.spectacular/` |
| `scope` | Project-wide vs per-request |
| `snapshot-on-edit` | Whether to snapshot before a substantive edit |

The rules-file **body** holds the per-doc grill prompts, vague-word lists, and gate checks. `references/doc-index.md` is a human-readable catalog only — it is **not** parsed for behavior (demoted from `doc-registry.md` in v1.4.0).

## Modes

A doc's `mode:` declares how the engine interacts with it. Nine modes exist:

| Mode | Behavior |
|---|---|
| `grill` | Interactive slot-fill. Alias for `grill-wide` (the default style). |
| `grill-wide` | One broad session — all slots filled in a single pass. |
| `grill-each` | Per-block walk — same slots repeated, one block at a time; agent asks "add another?" after each. (roadmap per-version, personas per-person.) |
| `grill-loop` | Wide pass first (short answers ok), then a deep pass over slots flagged vague/incomplete. |
| `append` | Capture one entry, append to the file. No slot loop. (decisions.) |
| `index` | **Soft-folder DB** — index file regenerated from `entries/`. CLI mutators write entries; agentic verbs operate on the *collection*, not a single file. Distinct from the slot-driven modes. (memory, sessions, feedback, idea.) |
| `stub` | Scaffold + exit; user edits the file directly thereafter. |
| `freeform` | Agent improvises the shape. **Reserved — no docs use this as of v1.x.** |
| `reference` | Skill-internal / read-only; the doc is informational, not interrogated. |

A `--wide` / `--each` / `--loop` flag overrides the declared `grill-*` mode for a single invocation.

## Registered docs

The live registry **is** the set of `references/*-rules.md` files — one rules file per dispatch contract (some serve more than one doc-id). `doc-index.md` is the human catalog. Rather than pin a headline count (it drifts every time a doc ships), here is the registry by **scope** — where each doc lives:

| Scope | Doc-ids | Rules file |
|---|---|---|
| **Project-wide** | `prd`, `spec`, `principles`, `architecture`, `roadmap`, `stack`, `agents`, `decisions`, `personas` | one each |
| | `memory`, `sessions`, `feedback`, `idea` | one each (soft-DB / `index` mode) |
| **Per-request** | `plan`, `tasks` | one each |
| **External** | `convention-pack` (`~/.spectacular/packs/`) | `pack-rules.md` |
| | `docs-manifest`, `docs-page` (`docs/`) — *deprecated v1.2.0 → pageworks* | shared `docs-rules.md` |
| **Skill-internal (reference)** | `migrations-contract`, `migration` | *none* — catalog markers only, no dispatch |

So: most ids map 1:1 to a rules file; `docs-rules.md` serves two; and the two reference-only ids have no rules file at all (they're skill-internal, documented in `migrations-contract.md`). The number of *live dispatch contracts* is the count of `*-rules.md` files — read it off disk rather than trusting a figure here.

## CLI vs skill split

Agentic verbs (`grill`, `refine`) are **skill-only** — running them through the CLI prints a friendly redirect, because interrogation is judgment work. Mechanical verbs (`new`, `archive`, `snapshot`, `init`, `doctor`, `pack`, `migrate`) run in the CLI. This is the mutation principle applied to docs: **skill orchestrates the conversation; CLI mutates the files.**

## Related references

- `references/grill.md` / `references/refine.md` / `references/review.md` — the three engine flows
- `references/doc-index.md` — human catalog of every doc type
- `references/prd-overrides.md` / `plan-overrides.md` / `tasks-overrides.md` — per-doc rule examples
- `references/kits-contract.md` — how kits add/modify slots and trigger docs

## Decisions

> Running design-log captured while writing this spec (2026-05-29), one decision at a time. The *why* behind the spec text above.

### Topic: grill modes — CLOSED

| # | Decision | Resolution | Why |
|---|---|---|---|
| 1 | Mode count | Document all **9** (include `index`) | `index` is real + shipped (v1.5–v1.7); the spec should be the accurate source of truth, not match the stale 8-mode docs |
| 2 | `index` treatment | Listed in the flat table, **annotated** as the soft-DB mode | It behaves differently (CLI mutators write entries; agentic verbs operate on the collection) — flag it, don't hide it |
| 3 | Table structure | **Flat 9-row list**, not family-grouped | Mirrors `doc-index.md`; least restructuring |
| 4 | `freeform` | **Kept** in the table, marked **(reserved — unused)** | Honest about latent capability; matches doc-index.md's existing label |

**Follow-up flagged:** `grill.md` + the root `SPEC.md` bullet document only 8 modes (omit `index`). They are now stale-on-modes — reconcile in a separate edit.

### Topic: registry shape — CLOSED

| # | Decision | Resolution | Why |
|---|---|---|---|
| 5 | Headline count | **No number** — describe, don't pin | "18 docs" matched nothing on disk (17 rules-files / 19 ids / 2 reference-only). A number drifts every release; the rules-file set is the real registry — read it off disk |
| 6 | Breakdown axis | **By scope** — project-wide / per-request / external | Shows where each doc lives; more durable than a mode-grouping that duplicates the modes table |
| 7 | Reference-only ids | **4th row: skill-internal** | `migrations-contract` / `migration` have no rules file — include them honestly rather than imply every id is engine-driven |

### Topic: planned work — CLOSED

| # | Decision | Resolution | Why |
|---|---|---|---|
| 8 | A "what's planned" section | **Omitted** | SPEC.md is present-tense by Spectacular convention — describes what *is*, not what's next. Forward work belongs in ROADMAP |
| 9 | Where adjacent items live | **Already homed — no new recording** | The defined-but-unbuilt `roadmap grill --icebox` (ROADMAP line 383) + `workflows/` doc-type (line 334) are in ROADMAP; the `decision-tree-grill-mode` idea is parked in `.spectacular/ideas/`. Nothing falls through |
