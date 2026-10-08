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
   objects rather than copying them. For a distributed PRD, read
   [requirements.md](requirements.md) and link the applicable needs and acceptance
   criteria. Use [drafting-methods.md](drafting-methods.md) only when a skeleton or
   anatomy helps the task. Use the compact default in
   [profiles.md](profiles.md); minimal remains valid and no UUID, governed schema,
   or lifecycle is required.
4. Persist when the host permits writes, then read back the file and check that
   its links resolve. Return its path and any consequential unresolved choice.
   A saved-plan claim requires the file to exist with the reported content.

Keep one coherent outcome per plan, usually 50–120 lines. Split by domain, code
module, phase, or delivered tranche when each part stands alone; keep the parent
as a short linking plan. Large reusable behavior belongs in optional `specs/`;
lasting explanations belong in `atlas/`. See [workspace.md](workspace.md).

## Milestones and delivery batches

Define the first usable outcome before decomposing substantial product work. A
milestone states who can complete which workflow, in what environment, and what
evidence closes it. A technical task contributes to that outcome; passing its
unit tests does not close the milestone. Separate a pilot, a broader usable
product, and release only when the project needs those distinctions; use its
own names. Do not infer dates, release acceptance, or additional authority.

For product-stage gates, exposure prerequisites, prioritization, or MVP value
experiments, read [milestone-readiness.md](milestone-readiness.md). Keep its
relevant facts in this plan; technical completion, distribution, and product
validation require different evidence.
Group tasks into a bounded batch with one observable result, dependencies,
checks, and a stopping boundary. Prefer a workflow exercised from input to
result over a collection of disconnected components. State whether evidence is
synthetic, integrated, exercised with a real dependency, or accepted by a person.
Keep correctness, recovery and permission checks alongside each feature.

For each experiment, name the question, smallest representative test, effort
bound, and decision it unlocks. Close with a choice and integration step, a
rejected option, or a reasoned deferral with a revisit condition. If inconclusive,
state the missing evidence and justify a bounded extension; do not append an
indefinite chain of micro-spikes or treat feasibility as working integration.

### Multi-Wave Campaign Plans (`CampaignPlan`)

For multi-milestone initiatives, structure the plan around sequential waves with explicit gates:
- **Wave & Track Topology**: Distinguish between sequential waves, *Blocking Prerequisites* (`blocks: [Wave X]`), and *Orthogonal Parallel Tracks* (disjoint write scopes). Each slice declares a feature branch, physical worktree (`.worktrees/<slug>`), and verification gate.
- **Dispatch Briefs**: Embed compact, ready-to-run handoff prompts for each wave ($\le 1200$ tokens) covering invariants, scope boundaries, test commands, and the structured JSON return receipt schema.
- **Pre-warmed Worktrees**: Ensure the orchestrator bootstraps runtime dependencies (e.g. symlinking `node_modules`, `.env.sandbox`, C-bindings) before dispatching workers.

## Resume and evolve

Reuse the same path through planning, implementation, and delivery. Record useful
progress, checks, and outcomes there; planning again updates the existing plan.
Approval and execution follow the host and user instructions, not a new Spectacular
gate. Do not create a Proposal, Contract, Mission, or audit record automatically.

When implementation is authorized, a repeated “continue” resumes the current
batch through its verified outcome. Keep intermediate updates and checks; a
completed fragment is a progress update, not a default stopping point. Stop for
the agreed boundary, a real blocker, missing effect authorization, or host limits,
and record what remains. Continuation does not authorize scope expansion, live
data access, publication, or a new Mission.

At a batch boundary, briefly reconcile implemented, verified, and open work in
the existing owners. The plan owns delivery detail; anchors summarize current
truth; an existing roadmap owns milestone state; a campaign or status document
may own receipts. Link details rather than copying them. Update only affected
claims, remove superseded next steps, and leave one executable resume point.
Create no extra status file or catalog by default; preserve frozen records.
When Git checkpoints are already authorized, align them with coherent verified
batches and report unpublished changes. Skill edits alone authorize no commit,
push, or tag.

A native host plan file can remain a cache or pointer. `.spectacular/plans/` is the
durable project copy; choose one content owner and reconcile changes when importing
from a native cache. If the copies conflict, resolve the difference instead of
silently overwriting either. A plan does not become a spec merely because it grew,
or a Mission because implementation began.
