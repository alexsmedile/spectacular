---
description: Strict dependency sequencing, uncertainty ranking, metaphor routing, and execution-scope boundaries for Wayfinding records.
when_to_use: Choosing the next discovery action, interpreting Wayfinder metaphors, or handling an unexpected discovery during implementation.
---

# Wayfinding Sequencer

The sequencer derives readiness from canonical `blocked_by` edges. Those edges are authoritative; fog and frontier are views, never persisted state.

## Ordering contract

1. Reject dangling edges and cycles before sequencing.
2. Topologically order nodes so every prerequisite appears before its dependents.
3. Consider only frontier nodes for `wayfind next`.
4. Rank priority first (`high → medium → low`), then uncertainty: user-input question → spike → research → other question → specification.
5. Use canonical-ID order as the deterministic final tie-break.

Dependencies override target versions. `doctor wayfinding` warns when a node targets an earlier SemVer/phase than one of its prerequisites. It also finds strong dependency language in PRD, roadmap, PLAN, and specification text and proposes explicit edges; it never writes or reslots them.

## Metaphor routing

| User language | Safe route | Boundary |
|---|---|---|
| Park this idea | `idea new` | Creates a parked `IDEA-NNN`; never edits current PLAN/TASKS. |
| Put it on ice / Icebox | `wayfind defer` | Requires a reason; preserves the record. |
| Find your way to X | `wayfind path` | Shows prerequisites before recommending action. |
| Act on goal X | `/spectacular act SPC-NNN` | Requires an approved specification and runs the request/HITL/activation gates in [[request-workflow]]. |

The CLI exposes these mappings through `spectacular wayfind route <phrase> ...`. The skill may translate natural language directly, but it must call the same underlying gated verb.

## Execution boundary

During active implementation, an unexpected requirement, tangent, or optimization is not permission to expand the milestone. Create an idea with `origin:` pointing to the active request, or assign the discovery to a later target version. Do not add, reorder, or silently broaden current PLAN/TASKS scope. Resume the current milestone after parking the loop.

The sequencer may recommend read-only research. A spike still requires authorization. It cannot merge code, make breaking API/schema changes, or resolve product/business questions for the user.
