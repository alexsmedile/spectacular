# Grill — generic interactive slot-filling engine

Loaded when the user runs `spectacular <doc> grill` (or `spectacular <doc>` when the doc is empty and `mode: grill`).

This is the **doc-agnostic** engine. PRD-specific behavior lives in `prd-overrides.md`; PLAN-specific behavior in `plan-overrides.md`; etc. The engine reads the registry entry for the requested doc to know which template, slots, location, and override file to use.

## Core principle

**One question at a time. Strict slot order. Stop when ready.**

No upfront interview. No multi-step research. The grill walks the slots declared in the registry, in order, runs an inline mini-refine on each answer using the doc's override patterns (if any), writes to disk immediately, and exits to a review gate at the end.

## Behavior

### 1. Resolve the registry entry

1. Read `references/doc-registry.md`, look up the requested `<doc>`.
2. If `mode != grill`, route accordingly: `append` → `refine` engine in append mode; `freeform` → just scaffold the template and exit.
3. Load: `template`, `slots`, `location`, `snapshot-on-edit`, `overrides`.

### 2. Pre-flight

Before asking anything:

1. Check if the file at `location` already exists.
   - If yes and not empty: ask "There's an existing <doc>. Refine it (`<doc> refine`) or start over (creates `<DOC>@vN.md` snapshot first if `snapshot-on-edit: true`)?"
   - If yes and effectively empty (only template placeholders): proceed to grill directly.
   - If no: scaffold from `template`. For per-request docs, supply `<slug>` from context (either the user provided it, or this was invoked via `spectacular new`).
2. If the doc supports kits (`kit-support: true`) and no kit is set in frontmatter, run kit selection per the doc's override file (see `prd-overrides.md` § Kit selection for PRD's flow; contract is documented in [[kits-contract]]).
3. Confirm or infer project name + per-doc context (e.g. PLAN.md needs a request slug).

### 2a. Kit application (only when kit-support: true and a kit was selected)

After scaffolding the base template, apply the kit's deltas:

1. **Read kit file** — parse frontmatter (`adds-slots`, `modifies-slots`, `triggers-docs`). Project-local override wins over bundled.
2. **Insert added slots** — for each `adds-slots` entry, find the `after:` base slot in the scaffolded file and insert the new slot heading right after it. Renumber slots accordingly.
3. **Layer modify-slot notes** — for each `modifies-slots` entry, append the note to that slot's prompt (used during slot loop).
4. **Set frontmatter** — write `kit: <kit-id>` to the file's frontmatter so the review gate knows which kit checks to apply.
5. **(Smart-init only)** — `triggers-docs.always` is consumed downstream by the init flow; the grill itself does not scaffold sibling docs.

### 3. The slot loop (strict order)

For each slot in the **resolved slot list** (registry's `slots:` list + active kit's `adds-slots` inserted at declared positions), in order:

```
Ask the slot's question
  ↓
Receive answer
  ↓
Run mini-refine inline (load patterns from overrides file if present)
  ↓
Write/update file immediately (replace <PLACEHOLDER> for that slot)
  ↓
Confirm: "Looks good? (y / edit / next)"
  ↓
On "next" → advance to next slot
On "edit" → re-ask with the user's nudge
On "y"   → advance
```

After the last slot: run the review gate (`review.md`). If it passes, exit. If not, show the punch list and loop on flagged items (out-of-order revisit is allowed here — strict order applies only during initial grill).

### 4. Slot prompts

The engine needs a question per slot. Sources, in order of preference:

1. **Kit-added slot** — if the slot comes from the active kit's `adds-slots`, use the kit's `prompt:` (and `example:` if present)
2. **Override file** — if `overrides:` is set, look for a `## Slot prompts` section listing per-slot prompts for base slots
3. **Template inline comments** — `<!-- ... -->` comments at the top of each slot section in the template
4. **Generic fallback** — `"Fill in the <Slot Name> section."`

For base slots that the active kit has in `modifies-slots`, **append** the kit's `note:` to the resolved prompt — never replace.

The user always sees the slot's section heading + the prompt. Examples (good vs bad) are optional but encouraged.

### 5. Optional sections

After all required slots pass the review gate, if the template includes optional sections (delimited by `<!-- ──── OPTIONAL SECTIONS ──── -->`), ask:

> The <doc> has all required slots filled. Want to add any optional sections now?
> - <list of optional section names>
> (or skip — you can add them later)

Skip silently if the user declines.

## Mini-refine (inline)

After every answer, the engine scans for vague-language patterns. If hit, *propose* a tighter version and ask the user to accept or override.

Pattern sources:
- **Base patterns** (universal): vague adjectives applied to slots not exempted by the override file.
- **Per-doc patterns** (override file): doc-specific rules like PRD's "plural-user → singular" or PLAN's "unbounded milestone → dated".

If the user can't resolve a flag right now, insert `[NEEDS CLARIFICATION: <specific gap>]` inline and continue. The review gate will catch it later.

Slots can be exempted from mini-refine via the override file's `## Mini-refine exemptions` section. (Example: PRD's Vision slot is exempt because narrative abstraction is expected.)

## Stop condition

The grill ends when:

1. All required slots have non-placeholder content, AND
2. The review gate passes (see `review.md`).

If the user wants to bail mid-grill, accept it — save what's filled, leave `<PLACEHOLDER>` markers for the rest, and tell them to run `<doc> review` later to see what's missing.

## What the grill does NOT do

- It does not research the domain. No web searches, no NotebookLM, no source ingestion. The user supplies the content; the grill structures it.
- It does not propose substance on its own — only sharper *phrasing*. The grill never invents user personas, success metrics, milestones, etc.
- It does not loop indefinitely. If the user keeps answering vaguely after 2 nudges per slot, accept the answer and move on. The review gate is the safety net.
- It does not write other docs. Grilling a PLAN never edits PRD; grilling a PRD never edits PLAN. Cross-doc generation is deferred to v2.

## Karpathy alignment

- **Think before coding:** the grill makes assumptions explicit by forcing measurable signals.
- **Simplicity first:** slot loop + inline mini-refine + review gate. No multi-agent pipelines.
- **Surgical changes:** writes only the slot being filled. Never reformats the rest.
- **Goal-driven:** the stop condition is the review gate, not "feels done".

## Examples

### PRD grill

Registry says: `mode: grill`, `slots: [Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone]`, `overrides: references/prd-overrides.md`.

Engine: scaffolds from `templates/prd/base.md`, walks 8 slots, runs PRD-specific mini-refine (kit-aware, Vision-exempt), exits on review gate pass.

### PLAN grill

Registry says: `mode: grill`, `slots: [Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables]`, `overrides: references/plan-overrides.md`.

Engine: scaffolds from `templates/plan/base.md`, walks 7 slots, runs PLAN-specific mini-refine (milestone ordering, dependency-link validation), exits on review gate pass.

### DECISIONS append

Registry says: `mode: append`. Engine does **not** invoke grill — routes to `refine` engine's append mode, which asks for a one-line title + decision + reasoning + tradeoffs, then appends a single entry to DECISIONS.md.

### ROADMAP freeform

Registry says: `mode: freeform`. Engine does **not** invoke grill — just scaffolds the template and exits. The user edits the file directly.

## Related

- [[doc-registry]] — the registry the engine consumes
- [[refine]] — vibe→spec rewriter, also handles append mode
- [[review]] — quality gate run at the end
- [[prd-overrides]] — per-doc rules (reference example)
- [[plan-overrides]] — per-doc rules
- [[scaffold-reference]] — what templates look like (separate concern)
