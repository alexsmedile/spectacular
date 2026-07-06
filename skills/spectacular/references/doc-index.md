---
description: Human catalog of every doc type, its mode and location. Dispatch lives in each <doc-id>-rules.md frontmatter, not here.
when_to_use: Looking up which doc types exist or where one lives.
---

# Doc Index — what Spectacular knows how to write

Human catalog of every document type in a Spectacular workspace.

**Dispatch lives in each doc's rules file** (`references/<doc-id>-rules.md` → frontmatter). This index is for browsing, not parsed for behavior.

## Project-wide canonical docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `prd` | grill | `.spectacular/PRD.md` | [prd-rules](prd-rules.md) |
| `principles` | stub | `.spectacular/PRINCIPLES.md` | [principles-rules](principles-rules.md) |
| `policy` | structured (always-set) | `.spectacular/POLICY.md` | [policy-rules](policy-rules.md) |
| `architecture` | stub | `.spectacular/ARCHITECTURE.md` | [architecture-rules](architecture-rules.md) |
| `spec` | stub | `.spectacular/SPEC.md` | [spec-rules](spec-rules.md) |
| `roadmap` | grill-each | `.spectacular/ROADMAP.md` | [roadmap-rules](roadmap-rules.md) |
| `stack` | stub | `.spectacular/STACK.md` | [stack-rules](stack-rules.md) |
| `agents` | stub | `.spectacular/AGENTS.md` | [agents-rules](agents-rules.md) |
| `decisions` | append | `.spectacular/DECISIONS.md` | [decisions-rules](decisions-rules.md) — **ADR / architecture-decision log**; write with `spectacular decide` |
| `memory` | index | `.spectacular/MEMORY.md` + `memory/` | [memory-rules](memory-rules.md) |
| `sessions` | index | `.spectacular/SESSIONS.md` + `sessions/` | [sessions-rules](sessions-rules.md) |
| `personas` | grill-each | `.spectacular/PERSONAS.md` | [personas-rules](personas-rules.md) |
| `feedback` | index | `.spectacular/feedback/` (+ `requests/<slug>/feedback/`) | [feedback-rules](feedback-rules.md) |
| `idea` | index | `.spectacular/ideas/` | [idea-rules](idea-rules.md) |
| `audit` | index | `.spectacular/audit/` | [audit-rules](audit-rules.md) — **bug investigation** before a fix is planned; write with `spectacular audit new\|list\|resolve` (v1.25.0) |
| `fixes` | index | `.spectacular/fixes/` | [fixes-rules](fixes-rules.md) — **verified-fix log**; write only once resolved, with `spectacular fix new\|list` (v1.25.0) |

## Per-request docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `plan` | grill | `.spectacular/requests/<slug>/PLAN.md` | [plan-rules](plan-rules.md) |
| `tasks` | stub | `.spectacular/requests/<slug>/TASKS.md` | [tasks-rules](tasks-rules.md) |
| `vision` | index (`imagine`) | `.spectacular/requests/<slug>/vision/` | [vision-rules](vision-rules.md) |

## User-scope docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `pack` | grill | `~/.spectacular/packs/<name>/pack.md` | [pack-rules](pack-rules.md) |

> `pack` was registered as `convention-pack` before v1.19.0; the old id is still accepted as an alias (`doc-id-aliases:` in pack-rules.md).

## Stub default behavior

`mode: stub` docs are scaffolded from a template and **edited directly** by the user — there is no grill interview. Unless a doc's own `<doc-id>-rules.md` body overrides it, every stub doc behaves identically across the three verbs:

| Verb | Behavior |
|---|---|
| `grill` | Polite no-op + hint: "this is a stub doc — open it in your editor, or pass `--wide` to grill it ad-hoc." |
| `refine` | Whole-doc rewrite pass (vibe → spec); does not interview slot-by-slot. |
| `review` | Structural check only: frontmatter present, no template placeholders left unfilled. |

`snapshot-on-edit: true` for the project-wide canonical stubs (PRINCIPLES, ARCHITECTURE, SPEC, STACK, AGENTS) — the skill snapshots before any edit (e.g. `PRINCIPLES@v2.md`). Per-request stubs (TASKS) do not snapshot.

A stub's rules file therefore only needs its **frontmatter** (the engine's dispatch) plus any genuinely doc-specific note. Files that restate the table above are thinned to a single pointer back here (see [[rules-files-audit]], decision D8). Two "stubs" carry real bodies and are *not* thinned: `spec` (index role + archive-time sync) and `tasks` (review-gate + refine patterns).

## Skill-internal references

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `migrations-contract` | reference | `skills/spectacular/references/migrations-contract.md` | — |
| `migration` | reference | `skills/spectacular/references/migrations/v<from>-to-v<to>.md` | — |
| `bug-workflow` | reference | `skills/spectacular/references/bug-workflow.md` | ties audit/ + fixes/ into the self-learning loop; loaded on any bug report |
| `build-workflow` | reference | `skills/spectacular/references/build-workflow.md` | **build-direction orchestrator arc** (mirror of bug-workflow) — assemble a closed milestone brief from the request chain, decide build-inline vs dispatch `spec-builder`, confirm + tick the ledger; loaded when implementing a milestone |
| `soft-db-index` | reference | `skills/spectacular/references/soft-db-index.md` | **canonical routing index** for the 7 soft-DB collections — role, purpose, structure, boundary rules |

## Mode taxonomy

| Mode | Behavior |
|---|---|
| `grill` | Interactive form. Alias for `grill-wide` (the default style). |
| `grill-wide` | One broad session — all slots filled in one pass. |
| `grill-each` | Per-block walk — same slots, one block at a time. Agent asks "add another?" after each block. |
| `grill-loop` | Wide pass first, then narrows to deepen slots that look vague. |
| `append` | Capture one entry, append to file. No slot loop. |
| `index` | Soft-folder DB. Index file regenerated from entries in `entries-dir/`. CLI mutators write entries; agentic verbs operate on the collection. |
| `imagine` | Generative-first variant of `index`. Agent **renders** ASCII fragments (stories/ui/arch) into a soft-folder + spine, human reacts per-fragment (`approved:`), then it **derives a draft PLAN** from the approved vision. Distinct from `grill` (which interrogates). Only the `vision` doc uses it. See [vision-rules](vision-rules.md). |
| `stub` | Scaffold + exit. User edits directly thereafter. |
| `freeform` | Agent improvises shape (reserved — no docs use this in v1.4). |
| `reference` | Skill-internal doc, not user-facing. |

### Verb × mode behavior

| | grill family | append | index | stub | freeform | reference |
|---|---|---|---|---|---|---|
| **grill** | Run interactive form (style per mode value; `--wide`/`--each`/`--loop` flags override) | Capture one new entry interactively | Capture one entry interactively, hand off to CLI mutator (`spectacular remember` / `session start`) | Polite no-op + hint ("stub doc — open in editor, or `--wide` to grill ad-hoc") | Open prompt; agent infers slot list on the fly | Error: not user-facing |
| **refine** | Vibe→spec rewrite across slots | Ask user: refine latest / all / pick | Refine specific entry by slug, or tagged subset | Whole-doc refine pass | Whole-doc refine pass | N/A |
| **review** | Slot gate-check (per-block for `grill-each`) | Validate entry shape | Validate entries + index-vs-entries drift + (sessions) stale-open check | Structural review (frontmatter + sections + no template placeholders) | Subjective review | Internal validation |

**Verbs `grill` and `refine` are skill-only** — they require an LLM. The CLI redirects with a friendly message when called at terminal. `review` is mixed (structural checks run in CLI; semantic in skill).

## Naming conventions & traps

The rule: **plural filename = a top-level index; singular = one file per request.** Two lexical exceptions and one hard trap — internalize these before routing:

1. **Two `SESSION`s, opposite scope.** Per-request `SESSION.md` (singular — one request's working state, created on `active`) is **unrelated** to the top-level `SESSIONS.md` + `sessions/` collection (the work-session time-log). Same word, different system. Do not conflate — the biggest confusion trap.
2. **`SPEC.md` is overloaded.** Top-level `.spectacular/SPEC.md` is a lightweight *index*; `specs/<cap>/SPEC.md` are the per-capability *truth* docs.
3. **`MEMORY.md` is singular-form but plays the plural index role** (indexes `memory/`, exactly as `SESSIONS.md` indexes `sessions/`).
4. **`VISION.md` is singular but acts as a spine/index** of the `vision/` soft-folder.
5. `FEEDBACKS.md` at the repo root is a Spectacular-repo dev artifact, **not** a canonical file type — the canonical store is the `feedback/` folder (no `FEEDBACK.md` index is emitted).

Full catalog with roles + usage: [[scaffold-reference]] § File-type catalog (also mirrored in `docs/scaffold.md` for external readers).

## Adding a new doc type

1. Create rules file at `skills/spectacular/references/<doc-id>-rules.md` with frontmatter (schema in [scaffold-reference](scaffold-reference.md))
2. Create template at the path declared in the rules file's `template:` field
3. Add a row to this index (cosmetic — no behavior depends on the index)

CLI auto-discovers via filesystem walk. No edits to `cli/spectacular` needed for new docs in v1.4.0+.

## Related

- [[scaffold-reference]] — frontmatter schema + template stubs for every file type
- [[grill]] / [[refine]] / [[review]] — verb behaviors
- [[kits-contract]] — kit extension schema (PRD-only in v1)
- [[packs-contract]] — pack schema
