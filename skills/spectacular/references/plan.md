# Durable planning

Use this when: the user invokes `/spectacular plan`, `$spectacular plan`, or asks
Spectacular to plan authorized work before implementation.

## Use the host’s planning behavior

This is a Skill route, not a `spectacular plan` binary command. The host keeps
its own plan-mode controls, tool restrictions, approval behavior, and native plan
cache. Use its native planning facility when available; do not claim to toggle a
host mode through prose. If the user must switch the mode in the UI, explain that
once and continue permitted planning. Otherwise plan directly in the conversation.

Planning authorization produces a plan, not implementation or Mission activation.
Follow host restrictions even when they prevent writing a plan file. In that case,
return the draft and intended durable path, state that it is unsaved, and save only
when writes become permitted. Do not emulate unavailable planning APIs.

## Persist one useful plan

1. Resolve the target project from its anchors and requested working directory.
   Read the relevant linked context and inspect the affected implementation.
   Finish when the outcome, known constraints, and real implementation boundaries
   are understood. Ask only for decisions that materially change the approach.
2. Reuse the plan that already owns this outcome. Otherwise choose
   `.spectacular/plans/<descriptive-slug>.md`; create only that needed directory.
   Check the destination before writing; preserve existing content and unknown
   metadata. Nested domain folders are useful only when several plans justify them.
3. Record the intended outcome, constraints, approach, implementation slices,
   verification, and unresolved choices. Link existing specifications and domain
   objects rather than copying them. Use the maintained metadata agreement from
   [knowledge.md](knowledge.md); no UUID, governed schema, or lifecycle is required.
4. Persist when the host permits writes, then read back the file and check that
   its links resolve. Return its path and any consequential unresolved choice.
   A saved-plan claim requires the file to exist with the reported content.

Keep one coherent outcome per plan, usually 50–120 lines. Split by domain, code
module, phase, or delivered tranche when each part stands alone; keep the parent
as a short linking plan. Large reusable behavior belongs in optional `specs/`;
lasting explanations belong in `atlas/`. See [workspace.md](workspace.md).

## Resume and evolve

Reuse the same path through planning, implementation, and delivery. Record useful
progress, checks, and outcomes there; planning again updates the existing plan.
Approval and execution follow the host and user instructions, not a new Spectacular
gate. Do not create a Proposal, Contract, Mission, or audit record automatically.

A native host plan file can remain a cache or pointer. `.spectacular/plans/` is the
durable project copy; choose one content owner and reconcile changes when importing
from a native cache. If the copies conflict, resolve the difference instead of
silently overwriting either. A plan does not become a spec merely because it grew,
or a Mission because implementation began.
