# PRD Grill — interactive slot-filling loop

Loaded when the user runs:
- `prd` or `prd grill` (no existing PRD.md, or explicit request to start fresh)
- `/spectacular prd <kit>` — kit-specific grill

The grill is the **primary interactive mode** for crafting a PRD from scratch.

## Core principle

**One question at a time. Strict slot order. Stop when ready.**

No upfront interview. No multi-agent research pipeline. The grill walks the 6 required slots in order, runs a mini-refine on each answer inline, writes to disk immediately, and exits to a review gate at the end.

## Behavior

### 1. Pre-flight

Before asking anything:

1. Check if `.spectacular/PRD.md` already exists.
   - If yes and not empty: ask user "There's an existing PRD. Refine it (`prd refine`) or start over (creates `PRD@vN.md` snapshot first)?"
   - If yes and effectively empty (only template placeholders): proceed to grill directly.
   - If no: create from template (see Kit selection).
2. Confirm or ask for the project name (from `config.yaml` if available).

### 2. Kit selection

If the user didn't specify a kit, ask:

> What kind of project is this?
> 1. **Coding** — CLI, library, app, service
> 2. **Product** — consumer or B2B with user flows
> 3. **Content** — course, newsletter, book, video
> 4. **Research** — investigation feeding a decision
> 5. **Blank** — none of the above / maximum flexibility

Copy the chosen kit from `templates/prd/kits/<kit>.md` to `.spectacular/PRD.md`. Fill `<DATE>` and `summary` frontmatter from context.

**Override path:** if `.spectacular/templates/prd/kits/<kit>.md` exists locally, use that instead of the skill's bundled template.

### 3. The slot loop (strict order)

Walk slots **1 → 6 in order**. For each slot:

```
Ask the slot's question
  ↓
Receive answer
  ↓
Run mini-refine inline (see Mini-refine below)
  ↓
Write/update PRD.md immediately (replace <PLACEHOLDER> for that slot)
  ↓
Confirm: "Looks good? (y / edit / next)"
  ↓
On "next" → advance to next slot
On "edit" → re-ask with the user's nudge
On "y"   → advance
```

After slot 6: run full review gate. If it passes, exit. If not, show the punch list and loop on flagged items (out-of-order revisit is allowed here — strict order applies only during initial grill).

### 4. Slot prompts

Each prompt is short. Show one example of good vs bad to anchor expectations.

**Slot 1 — Problem**
> What concrete pain does this solve? One sentence. Who is hurting, in what specific situation, how often.
>
> *Avoid:* "make X better", "improve Y experience"
> *Example:* "Solo devs writing PRDs from scratch waste 2+ hours and end up with vague documents nobody references."

**Slot 2 — Who it's for**
> Describe **one** primary user. Not a list, not "everyone". A specific role, situation, and constraint.
>
> *Example:* "Solo devs on side projects who use Claude Code and don't have a PM to write specs for them."

**Slot 3 — What success looks like**
> Measurable. Time-boxed. At least one number, one verb, and one date or timeframe.
>
> *Avoid:* "users love it", "ship fast"
> *Example:* "30 days after launch, 50% of users who run `/spectacular prd` open their PRD.md again within 7 days."

**Slot 4 — Non-goals**
> What are you **not** doing? List 3-5 explicit exclusions you'd push back on if asked to expand.

**Slot 5 — Constraints**
> What's fixed before you start? Budget, time, tech, policy, team.
>
> *Example:* "Markdown-only, no new binaries", "ships before 2026-07-01".

**Slot 6 — First milestone**
> One concrete, demoable outcome that proves this is real. Date-bound.

### 5. Optional sections

After slot 6 passes the review gate, ask:

> The PRD has all 6 required slots. Want to add any of these now?
> - Stakeholders
> - Risks
> - Open questions
> - Prior art
> (or skip — you can add them later)

Only prompt for sections the kit includes. Skip silently if the user declines.

## Mini-refine (inline)

After every answer, scan for these patterns. If hit, *propose* a tighter version and ask the user to accept or override:

| Pattern | Trigger | Proposed action |
|---|---|---|
| Vague adjective | `fast`, `intuitive`, `scalable`, `seamless`, `great`, `simple`, `flexible` | "What does '<word>' mean concretely? A number or comparison?" |
| Plural user | `users`, `customers`, `developers`, `people` (in slot 2) | "Pick **one** primary user. Who's the most important?" |
| Unbounded success | No number AND no date in slot 3 | "Add a number and a timeframe. Example: 'X by date Y'." |
| Tech jargon in problem | (slot 1 only) words like `microservices`, `embeddings`, `framework` | "Restate in plain language — what's the user-visible pain?" |
| Empty exclusion | Slot 4 has fewer than 2 items | "What would you push back on if someone tried to expand scope?" |

If the user can't resolve a flag right now, insert `[NEEDS CLARIFICATION: <specific gap>]` inline and continue. The review gate will catch it later.

## Stop condition

The grill ends when:

1. All 6 required slots have non-placeholder content, AND
2. The review gate passes (see `prd-review.md`).

If the user wants to bail mid-grill, accept it — save what's filled, leave `<PLACEHOLDER>` markers for the rest, and tell them to run `prd review` later to see what's missing.

## What the grill does NOT do

- It does not write the PLAN.md or TASKS.md for any request — that happens via `spectacular new <slug>` later.
- It does not research the domain. No web searches, no NotebookLM, no source ingestion. The user supplies the content; the grill structures it.
- It does not propose user personas, success metrics, or non-goals on its own. It can suggest sharper *phrasing*, but the substance comes from the user.
- It does not loop indefinitely. If the user keeps answering vaguely after 2 nudges per slot, accept the answer and move on. The review gate is the safety net, not the grill.

## Karpathy alignment

- **Think before coding:** the grill makes assumptions explicit by forcing measurable success.
- **Simplicity first:** 6 slots, no abstractions, no upfront research pipeline.
- **Surgical changes:** writes only the slot being filled. Never reformats the rest of the file.
- **Goal-driven:** the stop condition is the review gate, not "feels done".

## Related

- [[prd-refine]] — vibe→spec patterns used by mini-refine
- [[prd-review]] — quality gate run at the end
- [[scaffold-reference]] — full file templates
