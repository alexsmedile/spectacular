---
type: feedback
target: feedback-loop CLI ergonomics after M0-M4
scope: request
status: resolved
opened: 2026-05-25
resolved: 2026-05-25
proposal_summary: "Walk through the new/list/resolve/archive verbs as if shipping a real feedback session, surface friction"
next_action: ship-as-is
request: feedback-loop
spawned_request: null
promoted_to: null
related: []
---

# Feedback — feedback-loop CLI ergonomics after M0-M4

## Target
feedback-loop CLI ergonomics after M0-M4

## Hypothesis / hunch

After implementing M0-M4 in one go, the verbs work but there will be ergonomic friction visible only when actually using them. Surface it before M6 ships.

## Proposal

Run every verb in sequence on the live workspace and observe what breaks:
- `feedback-loop --help`
- `feedback-loop new` (request-scoped)
- top-level `--help` (confirm feedback-loop shows, aliases don't)
- `feedback-loop list` with one entry present

## Question asked

Of the 5 friction points surfaced, which to fix before M6 release vs ship and revisit?

## Friction observed

1. **Column overflow in `list`** — slug `2026-05-25-feedback-loop-cli-ergonomics-after-m0-m4` is 51 chars; column padded to 32 → STATUS/SCOPE/TARGET cells misalign for long slugs. Cosmetic.
2. **Date prefix redundancy** — slugs are `YYYY-MM-DD-<target>` so `list` shows the date twice (DATE column + slug prefix). Either drop the prefix from the on-disk slug or drop the DATE column.
3. **`--promote` is advisory, not active** — the flag prints what to run instead of running it. Safe but the name oversells. Either rename (`--promote-hint`) or make it interactive (y/n prompt then runs `remember`).
4. **No mid-loop notes verb** — entry body is all placeholders; user hand-edits to fill sections. Acceptable since the loop happens in /spectacular, but worth noting if friction compounds.
5. **`resolve` without `--next-action`** — leaves `next_action: tbd` silently. Should warn or require the flag.

## User response

Via AskUserQuestion (4 single-select):
1. Column overflow → **truncate slug to 32 chars with ellipsis**
2. Date prefix duplication → **keep date in slug, drop DATE column**
3. `--promote` is advisory → **rename to `--promote-hint`**
4. `resolve` without `--next-action` → **require `--next-action` (error if missing)**

Friction #4 (no mid-loop notes verb) was not raised — implicitly ship-as-is since loops happen in /spectacular.

## Insight

The dogfood pass surfaced 5 issues; 4 were small mechanical fixes worth shipping in the same release as M0-M4 rather than carrying them to a follow-up request. All four resolved one way:

- Make the CLI output narrower and predictable (truncate, drop redundant column)
- Make flag names honest (`--promote-hint` advertises what it actually does)
- Make resolve a deliberate act (require `--next-action`)

Pattern: when a verb's name implies a stronger contract than the implementation delivers, rename the flag rather than weaken the contract. Memory promotion is a judgment call that needs an LLM in the loop — the CLI's job is to print the right command, not pretend to make the judgment.

This is itself a durable preference worth a memory entry: **"prefer flag names that match actual behavior; CLI should hint, not pretend, when the action requires human judgment."**

## Decision

`next_action: ship-as-is` — all four fixes applied in this same session (commits: pending). Insight captured here; consider promoting to memory after one more session validates the pattern holds.

The dogfood proved the proactive-surfacing principle: this loop ran *because* we hit the M0-M4 checkpoint, and it caught real friction that would have shipped otherwise. The mode works as designed.
