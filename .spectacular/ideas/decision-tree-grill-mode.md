---
type: idea
status: parked
priority: medium
owner: alex
origin: (captured)
updated: 2026-05-29
promoted_to: null
related: []
---

# Idea — decision-tree-grill-mode

## Hypothesis

The grill engine should offer a **decision-tree style** for tasks that carry several independent decisions: divide the work into themes, then resolve **one question at a time** — each with a clear, self-contained explanation and 2-4 concrete options that spell out their trade-offs — closing one topic before opening the next. This is distinct from the existing slot-filling grill (which walks a *known* slot list); here the questions are *emergent design decisions*, not pre-declared slots.

## Context

Surfaced 2026-05-29 while expanding `specs/doc-engine/SPEC.md`. The doc-engine work threw off ~6 open decisions at once (mode count, families, freeform handling, registry count, planned-work framing). Batching them was unfollowable; switching to one-question-at-a-time, topic-by-topic made the decision tree legible and auditable.

Today's grill modes are all **slot-driven** — the agent walks a fixed `slots:` list. This idea is a new interaction *shape* where the agent surfaces decisions as they emerge from analysis, recaps each topic's resolution, and queues remaining themes. It mirrors the existing "one question at a time" core principle of `grill.md` but applies it to open-ended design decisions rather than form-filling.

Captured in parallel as a personal working-preference memory (Alex's global store) — this idea is the *product* candidate: making the interaction style a first-class engine capability others can invoke.

## Open questions

- Is this a new **mode** (`grill-tree`?) or a *flag* on existing grill (`--decide`) that switches from slot-walk to emergent-decision-walk?
- Where do the decisions come from — does the agent pre-enumerate them (a "decision manifest") or surface them lazily as analysis proceeds?
- How are resolved decisions persisted? (DECISIONS.md ADR entries? inline in the doc being grilled? a transient session log?)
- Does it need topic/theme grouping as a first-class structure, or is linear-with-recap enough?
- Overlap with the verify-walk request (v1.11) — both are "skill walks the user through a sequence." Shared engine or separate?

## Promoted to

—
