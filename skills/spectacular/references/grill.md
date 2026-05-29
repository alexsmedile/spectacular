# Grill — interactive slot-filling skill

Loaded when the user runs `spectacular <doc> grill` (or `spectacular <doc>` when the doc is empty and the doc's `mode:` is grill-family).

This is the **doc-agnostic** skill. PRD-specific behavior lives in `prd-rules.md`; PLAN-specific behavior in `plan-rules.md`; etc. The skill reads the requested doc's **rules file frontmatter** to know which template, slots, location, sub-mode, and per-doc rules to use.

> **Skill-only.** The `grill` verb requires an LLM to reason about answers, ask follow-ups, and run mini-refine. The CLI (`cli/spectacular`) does not dispatch grill — when called at terminal, it prints a friendly redirect to run inside Claude Code or Codex.

## Core principle

**One question at a time. Slot order respects the sub-mode. Stop when ready.**

No upfront interview. No multi-step research. The skill walks the slots declared in the rules file, runs an inline mini-refine on each answer using the doc's per-doc patterns (if any), writes to disk immediately, and exits to a review gate at the end.

## Behavior

### 0. Mode resolution

Each doc declares a `mode:` in its rules-file frontmatter. The `grill-*` family has four values:

| `mode:` | Behavior |
|---|---|
| `grill` | Alias for `grill-wide` (default style). |
| `grill-wide` | Walk all slots once, in order. One pass. (PRD, PLAN, convention-pack.) |
| `grill-each` | Per-block walk — same slots repeated for each block. Agent asks "add another?" after each completed block. (ROADMAP per-version, PERSONAS per-person.) |
| `grill-loop` | Wide pass first (fast, short answers ok), then deep pass over slots flagged as vague/incomplete. |

**Flag override.** If the user passes `--wide` / `--each` / `--loop`, the flag wins over the declared mode for this invocation. The rules-file mode is the default; the flag is a per-session override.

If `mode:` is **not** grill-family (i.e. `append` / `index` / `stub` / `freeform` / `reference`), route accordingly:
- `append` → `refine` skill in append mode (`refine.md` § append)
- `index` → soft-folder DB; grill doesn't fill a single file. Redirect: "`<doc>` is a soft-DB collection — add entries with its CLI mutator (e.g. `spectacular remember`, `spectacular decide`, `spectacular session start`, `spectacular idea new`, `spectacular feedback-loop new`). `grill`/`refine`/`review` operate on the collection, not one file." See [[doc-index]] mode taxonomy.
- `stub` → polite no-op: "`<doc>` is a stub doc. Open in editor, or pass `--wide` to grill it ad-hoc for this session."
- `freeform` → open-ended prompt: "What do you want to capture here?" Skill infers a slot list from the answer and walks it.
- `reference` → error: "`<doc>` is skill-internal, not user-facing."

### 1. Load rules file

1. Read `references/<doc-id>-rules.md`, parse frontmatter.
2. Load: `template`, `slots`, `location`, `scope`, `snapshot-on-edit`, `kit-support`, `mode`.
3. Body of the rules file contains slot prompts, vague-word lists, mini-refine patterns, gate checks.

**Substrate check (auto-invoked on failure):** if the rules file won't parse or is missing, or the active `kit:` file fails to load — auto-run `spectacular doctor kits frontmatter` and surface findings before refusing to grill. See [[doctor-substrate]].

### 2. Pre-flight

Before asking anything:

1. Check if the file at `location` already exists.
   - If yes and not empty: ask "There's an existing <doc>. Refine it (`<doc> refine`) or start over (creates `<DOC>@vN.md` snapshot first if `snapshot-on-edit: true`)?"
   - If yes and effectively empty (only template placeholders): proceed to grill directly.
   - If no: scaffold from `template`. For per-request docs, supply `<slug>` from context (either the user provided it, or this was invoked via `spectacular new`).
2. If the doc supports kits (`kit-support: true`) and no kit is set in frontmatter, run kit selection per the doc's rules file (see `prd-rules.md` § Kit selection for PRD's flow; contract is documented in [[kits-contract]]).
3. Confirm or infer project name + per-doc context (e.g. PLAN.md needs a request slug).

### 2a. Kit application (only when kit-support: true and a kit was selected)

After scaffolding the base template, apply the kit's deltas:

1. **Read kit file** — parse frontmatter (`adds-slots`, `modifies-slots`, `triggers-docs`). Project-local override wins over bundled.
2. **Insert added slots** — for each `adds-slots` entry, find the `after:` base slot in the scaffolded file and insert the new slot heading right after it. Renumber slots accordingly.
3. **Layer modify-slot notes** — for each `modifies-slots` entry, append the note to that slot's prompt (used during slot loop).
4. **Set frontmatter** — write `kit: <kit-id>` to the file's frontmatter so the review gate knows which kit checks to apply.
5. **(Smart-init only)** — `triggers-docs.always` is consumed downstream by the init flow; the grill itself does not scaffold sibling docs.

### 3. The slot loop

For each slot in the **resolved slot list** (rules file's `slots:` + active kit's `adds-slots` inserted at declared positions), the slot loop respects the resolved sub-mode:

**Inner loop (always the same):**

```
Ask the slot's question
  ↓
Receive answer
  ↓
Run mini-refine inline (load patterns from rules file if present)
  ↓
Write/update file immediately (replace <PLACEHOLDER> for that slot)
  ↓
Confirm: "Looks good? (y / edit / next)"
  ↓
On "next" → advance to next slot
On "edit" → re-ask with the user's nudge
On "y"   → advance
```

**Outer loop, per sub-mode:**

| Sub-mode | Outer loop |
|---|---|
| `grill-wide` | Walk slots 1..N once. After slot N: run review gate. If it passes, exit. If not, show punch list and revisit flagged slots (out-of-order revisit allowed here). |
| `grill-each` | Walk slots 1..N for one block. After slot N: ask "review this block? (y/skip)". On review pass, ask "add another block?". Repeat until user declines. Each block gets independent slot walks. |
| `grill-loop` | **Pass 1 (wide):** walk slots 1..N fast — accept one-line / short answers, mini-refine relaxed. Mark each slot as needs-deepening if it matches the [grill-loop heuristic](#grill-loop-heuristic). **Pass 2 (deep):** revisit only flagged slots; run full grill-wide quality on each. Exit on review gate pass. |
| `grill` (= `grill-wide`) | Same as grill-wide. |

### 3a. grill-loop heuristic

In pass 1 of `grill-loop`, mark a slot as **needs-deepening** if **any of:**

1. Answer length < 30 characters
2. Answer matches any word from the rules file's vague-word list (scoped to that slot)
3. Answer contains placeholder strings: `<…>`, `TODO`, `tbd`
4. Slot has an explicit gate-check (in the rules file) that fails on this answer

In pass 2, walk only flagged slots with full grill-wide quality (mini-refine, confirmation step, etc.).

If the user explicitly asks to skip pass 2, accept it and exit to the review gate.

### 4. Slot prompts

The skill needs a question per slot. Sources, in order of preference:

1. **Kit-added slot** — if the slot comes from the active kit's `adds-slots`, use the kit's `prompt:` (and `example:` if present)
2. **Rules file** — if `rules:` is set, look for a `## Slot prompts` section listing per-slot prompts for base slots
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

After every answer, the skill scans for vague-language patterns. If hit, *propose* a tighter version and ask the user to accept or override.

Pattern sources:
- **Base patterns** (universal): vague adjectives applied to slots not exempted by the rules file.
- **Per-doc patterns** (rules file): doc-specific rules like PRD's "plural-user → singular" or PLAN's "unbounded milestone → dated".

If the user can't resolve a flag right now, insert `[NEEDS CLARIFICATION: <specific gap>]` inline and continue. The review gate will catch it later.

Slots can be exempted from mini-refine via the rules file's `## Mini-refine exemptions` section. (Example: PRD's Vision slot is exempt because narrative abstraction is expected.)

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

### PRD grill (mode: grill = grill-wide)

`prd-rules.md` frontmatter: `mode: grill`, 8 slots, `kit-support: true`.

Skill: scaffolds from `templates/prd/base.md`, runs kit selection, walks 8 slots once, applies PRD-specific mini-refine (Vision-exempt, plural-user → singular, etc.), exits on review gate pass.

### PLAN grill (mode: grill)

`plan-rules.md` frontmatter: `mode: grill`, 7 slots.

Skill: scaffolds from `templates/plan/base.md`, walks 7 slots once, applies PLAN mini-refine (milestone ordering, dependency-link validation), exits on review gate pass.

### ROADMAP grill-each

`roadmap-rules.md` frontmatter: `mode: grill-each`, 6 slots per version block.

Skill: scaffolds from `templates/roadmap/base.md`, walks 6 slots for the current version block, runs per-block review gate, then asks "add another version block?". Each version is independent.

### PERSONAS grill-each

`personas-rules.md` frontmatter: `mode: grill-each`, 5 slots per persona.

Skill: scaffolds from `templates/personas/base.md`, walks 5 slots for one persona block, asks "add another persona?". User can stop at 1-5 personas (review gate warns if >5).

### PRD grill-loop (flag override)

User runs `spectacular prd grill --loop`. Override forces `grill-loop` regardless of declared mode.

Skill: pass 1 walks all 8 slots fast accepting short answers, flags slots that match the heuristic (e.g. "users love it" hits vague-word list for slot 5). Pass 2 revisits flagged slots with full grill-wide quality.

### DECISIONS append

`decisions-rules.md` frontmatter: `mode: append`. Skill does **not** invoke the slot loop — routes to `refine.md` § append, which asks for title + context + decision + consequences, then appends one entry.

### AGENTS stub (grill called)

User runs `spectacular agents grill`. `agents-rules.md` frontmatter: `mode: stub`.

Skill prints: *"AGENTS.md is a stub doc. Open in editor, or pass `--wide` to grill it ad-hoc for this session."* Exits.

If `--wide` was passed: skill generates a slot list on the fly from the existing AGENTS.md sections, then walks them with grill-wide behavior. The stub mode is unchanged — this is a one-off session.

## Related

- [[doc-index]] — human catalog of doc types
- [[refine]] — vibe→spec rewriter, also handles append mode
- [[review]] — quality gate run at the end
- [[prd-rules]] — per-doc rules (reference example)
- [[plan-rules]] — per-doc rules
- [[roadmap-rules]] — grill-each example
- [[personas-rules]] — grill-each example
- [[scaffold-reference]] — what templates look like
