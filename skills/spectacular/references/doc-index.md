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
| `decisions` | append | `.spectacular/DECISIONS.md` | [decisions-rules](decisions-rules.md) |
| `memory` | index | `.spectacular/MEMORY.md` + `memory/` | [memory-rules](memory-rules.md) |
| `sessions` | index | `.spectacular/SESSIONS.md` + `sessions/` | [sessions-rules](sessions-rules.md) |
| `personas` | grill-each | `.spectacular/PERSONAS.md` | [personas-rules](personas-rules.md) |
| `feedback` | index | `.spectacular/feedback/` (+ `requests/<slug>/feedback/`) | [feedback-rules](feedback-rules.md) |
| `idea` | index | `.spectacular/ideas/` | [idea-rules](idea-rules.md) |

## Per-request docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `plan` | grill | `.spectacular/requests/<slug>/PLAN.md` | [plan-rules](plan-rules.md) |
| `tasks` | stub | `.spectacular/requests/<slug>/TASKS.md` | [tasks-rules](tasks-rules.md) |

## User-scope docs

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `convention-pack` | grill | `~/.spectacular/packs/<name>/pack.md` | [pack-rules](pack-rules.md) |

## Public-facing docs (deprecated v1.2.0)

> See [pageworks](https://github.com/alexsmedile/pageworks) for current public-docs authoring. These entries remain for backwards compatibility; removal target v2.0.0.

| Doc | Mode | Location | Rules | Status |
|---|---|---|---|---|
| `docs-manifest` | stub | `docs/docs.yaml` | [docs-rules](docs-rules.md) | deprecated v1.2.0 |
| `docs-page` | stub | `docs/<section>/<slug>.md` | [docs-rules](docs-rules.md) | deprecated v1.2.0 |

## Skill-internal references

| Doc | Mode | Location | Rules |
|---|---|---|---|
| `migrations-contract` | reference | `skills/spectacular/references/migrations-contract.md` | — |
| `migration` | reference | `skills/spectacular/references/migrations/v<from>-to-v<to>.md` | — |

## Mode taxonomy

| Mode | Behavior |
|---|---|
| `grill` | Interactive form. Alias for `grill-wide` (the default style). |
| `grill-wide` | One broad session — all slots filled in one pass. |
| `grill-each` | Per-block walk — same slots, one block at a time. Agent asks "add another?" after each block. |
| `grill-loop` | Wide pass first, then narrows to deepen slots that look vague. |
| `append` | Capture one entry, append to file. No slot loop. |
| `index` | Soft-folder DB. Index file regenerated from entries in `entries-dir/`. CLI mutators write entries; agentic verbs operate on the collection. |
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

## Adding a new doc type

1. Create rules file at `skills/spectacular/references/<doc-id>-rules.md` with frontmatter (schema in [scaffold-reference](scaffold-reference.md))
2. Create template at the path declared in the rules file's `template:` field
3. Add a row to this index (cosmetic — no behavior depends on the index)

CLI auto-discovers via filesystem walk. No edits to `cli/spectacular` needed for new docs in v1.4.0+.

## Related

- [[scaffold-reference]] — frontmatter schema + template stubs for every file type
- [[grill]] / [[refine]] / [[review]] — verb behaviors
- [[kits-contract]] — kit extension schema (PRD-only in v1)
- [[packs-contract]] — convention-pack schema
