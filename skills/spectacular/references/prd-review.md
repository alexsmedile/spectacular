# PRD Review — quality gate

Loaded when the user runs:
- `prd review` — explicit review pass on `.spectacular/PRD.md`
- Implicitly at the end of `prd-grill` (final gate before exit)

Review is a **pass/fail quality gate**. It produces a punch list, not a rewrite. The user decides what to fix.

## Core principle

**No score. No rewrite. Just a checklist + actionable findings.**

Review is the safety net at the end of every grill and every refine. It reads `.spectacular/PRD.md` and reports exactly what blocks it from being "done".

## The gate (all must pass)

| # | Check | How |
|---|---|---|
| 1 | All 6 required slots non-empty | No `<PLACEHOLDER>` markers remain in slots 1-6 |
| 2 | No unresolved clarifications | Zero `[NEEDS CLARIFICATION: ...]` markers anywhere |
| 3 | Success criteria is measurable | Slot 3 contains ≥1 number AND ≥1 verb AND ≥1 date/timeframe |
| 4 | No vague-word hits in critical slots | Slots 1, 3, 6 contain none of the vague-word list (see `prd-refine.md`) |
| 5 | Non-goals are specific | Slot 4 has ≥2 items, none of which are vague meta-terms ("scope creep", "feature bloat") |
| 6 | Single primary user | Slot 2 reads as a singular role + situation + constraint (heuristic: no plural noun phrases like "users", "customers", "developers" without qualifiers) |

If all 6 pass → PRD is usable. Exit with a green confirmation.
If any fail → produce a punch list.

## Detection rules

### Check 1: Placeholders

Look for these unresolved markers anywhere in slots 1-6:
- `<PROBLEM>`, `<PRIMARY USER>`, `<SUCCESS CRITERIA>`, `<NON-GOAL 1>`, `<CONSTRAINT 1>`, `<MILESTONE>`
- Any uppercase angle-bracket placeholder (`<ANYTHING>`)

### Check 2: Clarifications

Regex: `\[NEEDS CLARIFICATION:.*?\]` — any match anywhere fails the gate.

### Check 3: Measurable success

Slot 3 content must contain:
- **A number** — digits (`50`, `1000`) OR quantified comparison (`half of`, `3x`, `most`)
- **A verb** — action or observable outcome (`open`, `complete`, `return`, `reach`, `ship`, `launch`)
- **A timeframe** — explicit date (`2026-09-01`), relative window (`within 30 days`, `by Q3`), or named event (`at launch`, `by demo day`)

Missing any → fails. Punch list says exactly which signal is absent.

### Check 4: Vague-word scan

Use the vague-word list from `prd-refine.md`. Scan slots 1, 3, and 6 only (not 2/4/5 — those have legitimate uses).

Vague words: `fast`, `slow`, `intuitive`, `simple`, `easy`, `scalable`, `seamless`, `great`, `good`, `nice`, `clean`, `flexible`, `robust`, `powerful`, `elegant`, `modern`, `smart`, `lightweight`, `solid`, `polished`.

Each hit produces one punch-list entry: `Slot N: "<word>" — quantify or replace`.

### Check 5: Specific non-goals

Slot 4 fails if:
- Fewer than 2 list items
- Any item is a vague meta-term: `scope creep`, `feature bloat`, `over-engineering`, `complexity`, `tech debt` (without context)

These meta-terms describe what you *don't want*, not what you're *excluding*. Push for concrete: "Not a multi-agent pipeline" beats "scope creep".

### Check 6: Single primary user

Heuristic — slot 2 fails if:
- The first sentence uses a bare plural noun ("Users", "Customers", "Developers", "Teams") without a qualifier like "Solo devs who..." or "5-15 person teams that..."
- The slot lists multiple distinct user types (e.g. "Both PMs and engineers")

This is heuristic, not strict — the user can override by explicitly answering "yes, this is a singular primary user" during grill.

## Output format

### When the gate passes

```
✓ PRD review passed

  .spectacular/PRD.md is ready.
  All 6 required slots filled, measurable success, no clarification gaps.

  Next:
    - spectacular new <slug>      scaffold the first request
    - spectacular snapshot PRD    create a versioned snapshot
```

### When the gate fails (punch list)

```
✗ PRD review found 4 issues

  Slot 1 — Problem
    L12  vague: "easier" — quantify or replace

  Slot 2 — Who it's for
    L18  plural primary user: "developers"
         → pick one specific role + situation + constraint

  Slot 3 — Success
    L24  missing timeframe — add a date or window

  Slot 6 — First milestone
    L40  [NEEDS CLARIFICATION: define MVP scope]

  Fix path:
    spectacular prd refine        # propose rewrites for all findings
    spectacular prd grill         # re-run the slot loop on flagged slots
    (or edit .spectacular/PRD.md directly and re-run review)
```

## What review does NOT do

- It does not score. Pass or fail.
- It does not rewrite. The user (or `prd refine`) does that.
- It does not check optional sections (stakeholders, risks, etc.). They're optional — only the 6 required slots gate.
- It does not check PLAN.md or TASKS.md. Those have separate gates (future).
- It does not warn about long PRDs. Length isn't a quality signal here.

## Dogfood test

This skill's own PRD must pass:
1. `.spectacular/PRD.md` (the spectacular product PRD) — passes
2. `.spectacular/requests/prd-craft/PLAN.md` — the request-level PLAN should also pass (different artifact, but follows the same measurable-success rule)

## Karpathy alignment

- **Think before coding:** the gate makes assumptions checkable.
- **Simplicity first:** 6 binary checks, no scoring, no NLP.
- **Surgical changes:** review never touches the file, only reads.
- **Goal-driven:** every check has a verifiable signal.

## Related

- [[prd-grill]] — runs review as its stop condition
- [[prd-refine]] — produces the rewrites that fix review findings
- [[versioning]] — snapshot before any significant rewrite
