---
name: rapid-prototyping
description: >-
  Rapidly explore and de-risk ambiguous product, UI, UX, or architectural decisions by generating three tracer fragments (A, B, C)
  across a 5-tier fidelity ladder. Use when a requirement has multiple viable implementation paths or when comparing distinct
  design choices before committing to code.
  Triggers on "prototype", "explore options", "tracer fragment", "matrix prototype", "compare designs", "fidelity ladder",
  "rapid prototyping", or "grow fragment".
  Do not invoke for clearly specified single-path tasks, simple bugs, finished code reviews, or open-ended brainstorming without decision intent.
---

# Rapid Prototyping

Explore competing options quickly by generating three concrete tracer fragments (A, B, C) and growing the selected direction along a strict fidelity ladder.

## Fidelity Ladder

| Level | Fragment | Deliverable / Evidence |
|---|---|---|
| **1. Atom** | Token, label, control, state cue, data snippet | Focused specimen or constraint trace |
| **2. Component** | Card, form section, toolbar, isolated widget | Standalone component or focused mock |
| **3. Composition** | Dashboard region, page layout, panel grouping | Responsive layout and state variations |
| **4. Screen / Flow** | Full page, modal flow, multi-step sequence | Connected states and interactive paths |
| **5. Integration** | Production path in the live repository | Verified codebase changes, tests, and builds |

Call the candidate artifact a **tracer fragment**: the smallest representative slice that lets the user judge the direction before committing full implementation effort.

## Decision Ledger Invariant

Always maintain the 4-field decision ledger at the start of every exploration round:
- `Locked`: Dimensions agreed upon and held constant across all options.
- `Open`: The single specific design or architectural axis being settled this round.
- `Level`: Current fidelity level (1 to 5).
- `Lineage`: History of choices (e.g. `root → B → A`).

## Workflow

1. **Frame one decision:** Identify the single open axis, the fidelity level, locked constraints, and 2–3 success signals.
2. **Generate 3 tracer fragments (A, B, C):** Produce three distinct, concrete options spanning the trade-off space.
3. **Present the matrix:** Display the options side-by-side with clear differentiators.
4. **Transition upon feedback:**
   - `A / B / C` : Select an option and advance to the next fidelity level (read [references/transitions.md](references/transitions.md)).
   - `M` : Merge specific attributes from multiple options.
   - `R` : Replace/refocus the open axis.
   - `S` : Step down to a simpler fidelity level to settle a prerequisite.
   - `F` : Finalize the selected fragment into the target codebase.

## Expansion Handoffs

| Out-of-Scope Need | Action / Delegate |
|---|---|
| Bounded contexts, C4 architecture, ADR documentation | Invoke `system-architecture` companion skill |
| Database schema design, Crow's Foot ERD, DDL migration | Invoke `data-modeling` companion skill |
| Mission governance, flight plans, verification receipts | Invoke `spectacular` mission governance |

## Core Invariants & Negative Constraints

- **DO NOT mutate production code before Level 5 integration.** Keep candidate fragments isolated or disposable until final lineage is uniquely approved.
- **DO NOT present options without explicit decision axes.** Every matrix round must isolate exactly one open axis.
