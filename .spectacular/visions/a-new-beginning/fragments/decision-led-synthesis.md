---
kind: strategy
caption: How drafted sources become explicit choices, specs, and an executable refactor without bloating Spectacular
reaction: pending
reaction_note: ""
updated: 2026-08-07
related: []
---

# Strategy — decision-led-synthesis


## Proposal

Use a repeatable, provenance-preserving funnel for every supplied plan or
suggestion.

### 1. Register the source

Capture a stable label or link, its date/context when known, and what authority
it actually has. Summarize rather than duplicate the full source.

### 2. Produce a short source card

Each card contains only:

- the problem the source claims to solve;
- its core proposals and intended outcomes;
- assumptions that must be true;
- overlap or conflict with current behavior and other sources;
- evidence needed to validate material claims;
- valuable ideas that can stand independently of the source's full design.

### 3. Cluster by decision, not by document

Move claims into bounded product themes such as core model, artifact taxonomy,
lifecycle, CLI/skill experience, Git/provider boundaries, migration, or
distribution. One decision packet may cite several sources; one source may feed
several packets.

### 4. Resolve one decision packet at a time

Present the owner with:

| Slot | Required content |
|---|---|
| Question | The exact choice that changes the design |
| Options | Genuinely distinct paths, including keep/remove when relevant |
| Evidence | Current behavior, tests, source claims, and known unknowns |
| Trade-offs | User value, complexity, compatibility, migration, maintenance |
| Recommendation | One clear recommendation and why |
| Decision | `pending` until the owner explicitly chooses |

If a fact can settle the packet, route bounded research. If feasibility—not
preference—is uncertain, propose a spike. If neither would change the choice,
avoid discovery ceremony.

### 5. Record outcomes without laundering uncertainty

Classify each proposal as `adopt`, `adapt`, `reject`, `defer`, or
`needs-evidence`. Record architectural choices as decisions only after explicit
confirmation. Park valuable out-of-scope possibilities as ideas; do not smuggle
them into the active design.

### 6. Converge through small specifications

Only approved Vision fragments become requirements. Derive cohesive specs around
real behavior or contract boundaries, not one specification per source and not
one mega-spec for the entire refactor. Resolve cross-spec joins before planning
execution.

### 7. Plan and execute in earned order

After spec approval, form requests whose milestones deliver the smallest
verified slices in dependency order. Every code addition must map to an approved
requirement and validation evidence; unused abstractions, duplicate ledgers,
speculative extensibility, and compatibility aliases without demonstrated need
are rejected as scope rather than "future-proofing."

## Decision impact

If approved, this becomes the working method for the refactor session. It does
not decide Spectacular's product architecture; it decides how competing inputs
will be evaluated and how explicit owner choices become specs and execution
without losing provenance or manufacturing consensus.
