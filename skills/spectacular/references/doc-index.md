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
| `spec` | stub | `.spectacular/specs/index.md` | [spec-rules](spec-rules.md) |
| `roadmap` | grill-each | `.spectacular/roadmaps/index.md` | [roadmap-rules](roadmap-rules.md) |
| `stack` | stub | `.spectacular/STACK.md` | [stack-rules](stack-rules.md) |
| `agents` | stub | `.spectacular/AGENTS.md` | [agents-rules](agents-rules.md) |
| `decisions` | append | `.spectacular/decisions/index.md` | [decisions-rules](decisions-rules.md) — **ADR / architecture-decision log**; write with `spectacular decide` |
| `memory` | index | `.spectacular/memories/index.md` + `memories/` | [memory-rules](memory-rules.md) |
| `sessions` | index | `.spectacular/sessions/index.md` + `sessions/` | [sessions-rules](sessions-rules.md) |
| `personas` | grill-each | `.spectacular/PERSONAS.md` | [personas-rules](personas-rules.md) |
| `feedback` | index | `.spectacular/feedbacks/index.md` + `feedbacks/` | [feedback-rules](feedback-rules.md) |
| `idea` | index | `.spectacular/ideas/index.md` + `ideas/` | [idea-rules](idea-rules.md) |
| `audit` | index | `.spectacular/audits/index.md` + `audits/` | [audit-rules](audit-rules.md) — **bug investigation** before a fix is planned; write with `spectacular audit new\|list\|resolve` (v1.25.0) |
| `fixes` | index | `.spectacular/fixes/index.md` + `fixes/` | [fixes-rules](fixes-rules.md) — **verified-fix log**; write only once resolved, with `spectacular fix new\|list` (v1.25.0) |

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
| `bug-workflow` | reference | `skills/spectacular/references/bug-workflow.md` | **runtime core** — ties audit/ + fixes/ into the self-learning loop; loaded on any bug report |
| `bug-workflow-doctrine` | reference | `skills/spectacular/references/bug-workflow-doctrine.md` | the *why* behind bug-workflow's gates — load only when a routing call is uncertain or when editing the workflow |
| `build-workflow` | reference | `skills/spectacular/references/build-workflow.md` | **build-direction orchestrator arc, runtime core** (mirror of bug-workflow) — assemble a closed milestone brief from the request chain, decide build-inline vs dispatch `spec-builder`, confirm + tick the ledger; routes the optional fleet (`repo-explorer` map-before-plan, `code-reviewer` + `test-verifier` arms-length gates); loaded when implementing a milestone |
| `build-workflow-doctrine` | reference | `skills/spectacular/references/build-workflow-doctrine.md` | the *why* behind build-workflow's gates + the relation-to-bug-workflow table — load only when a routing call is uncertain or when editing the workflow |
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

## Naming conventions under OKF

The rule: **Plural folder name = a category directory containing either an index.md index file and sequential entry files (soft-DB collections), or sub-directories representing execution tasks/traces.** 

1. **Two `SESSION`s, opposite scope:** Per-request `SESSION.md` (singular — one request's working state, created on `active`) is **unrelated** to the top-level `sessions/` category folder and its `sessions/index.md` (the work-session time-log).
2. **Consolidated Specs:** Capability specifications are stored as flat files inside the `specs/` directory (e.g., `specs/cli.md`, `specs/skill.md`), indexed by `specs/index.md`. There are no nested spec folders.
3. **Plural Folders:** All collection and execution folders are strictly plural (`memories/`, `roadmaps/`, `decisions/`, `sessions/`, `audits/`, `fixes/`, `feedbacks/`, `ideas/`, `requests/`, `debugs/`).
4. **Soft-DB Index-Logged Collections vs. Execution Trees**:
   - **Soft-DB Collections** (`memories/`, `decisions/`, `sessions/`, `audits/`, `fixes/`, `feedbacks/`, `ideas/`, `roadmaps/`) contain a central `index.md` and flat, sequential/date-logged `.md` entries.
   - **Execution Trees** (`requests/`, `debugs/`) contain sub-directories (`requests/<slug>/`, `debugs/<slug>/`) holding active state files or run logs rather than simple flat document entries.
5. **Sequential Prefix Identifiers:** Collection entry files under `decisions/` and `memories/` are prefixed with their ID sequence: `decisions/D<N>-<slug>.md` and `memories/M<N>-<slug>.md`.

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
