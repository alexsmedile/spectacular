---
name: spectacular
version: 2.1.0
description: Guide work in a canonical Spectacular v2 workspace through optional Proposal exploration, compact Mission preparation and activation, governed execution, earned Objective/Run expansion, Evidence and review, owner completion, audit, and cold recovery. Use for `/spectacular` jobs such as orient, explore, propose, plan, start, resume, handoff, review, complete, or audit; for compiling bounded runtime context or an Autopilot charter; and for safely continuing after session or runtime replacement.
---

# Spectacular

Run one bounded Mission at a time, from truth the owner already accepted.

## Working model

Every surface has one job. Judgment stays in this Skill and the host session;
fixed invariants stay in the tooling; execution stays in the host runtime; effects
stay in their own providers.

| Surface | Responsibility |
|---|---|
| Anchor | Accepted project truth: direction, boundaries, constraints. `PROJECT.md` is the root one |
| Proposal | Optional, mutable exploration in chat, an issue, or a Spectacular file |
| Mission | Frozen execution plan and primary entry point in `MISSION.md` |
| Contract/specification | Accepted product behavior; edit as ordinary Mission work |
| Decision | ADR-like durable choice and rationale, never routine lifecycle approval |
| Objective / Run | Inline while simple; promoted to a file only when independently useful |
| Evidence / review | Earned proof and assessment; neither silently changes success criteria |

## Start every workflow

1. **Find the workspace.** Read `.spectacular/PROJECT.md`. Only if that file is
   missing, read root `PROJECT.md`.

   That file is the root Anchor. Its `current_truth` field names the other Anchors
   and Contracts by id. Read those, and skip the ones this Mission does not touch.

2. **Run one read-only preflight.** Check:
   - workspace, and which Mission is selected
   - Git: branch, clean tree, fresh against upstream/default, release state
   - Contract and baseline bindings
   - owner, activation time, activation fingerprint
   - validation mode
   - current Objective and Run
   - blockers: dependencies, Gaps, stops

   Report three lines: plain outcome, one technical evidence line, one next action.

3. **Enter the Mission through `MISSION.md`.** Read its frontmatter control card
   and its body. Follow pointers to Objective, Run, Evidence, or review files only
   when the current work needs them. Do not preload project history, every record,
   or generated catalogs.

4. **Match the validation mode.**
   - Supported mechanical validation: use the typed `show` / `check` command.
   - `manual-bootstrap`: the CLI is out of service. Edit canonical Markdown
     directly and verify by hand — see [bootstrap.md](references/bootstrap.md) for
     the checklist. Never cite the old CLI as proof.

5. **Route to one reference** using the table below. Before any consequential
   effect, and again when resuming, revalidate bindings, authority, budgets, and
   stops.

## Route the guided job

Pick by what the session actually needs. The user rarely names the job.

| Load | When the session needs it | Guided job |
|---|---|---|
| [orient.md](references/orient.md) | You do not yet know the project direction, which Mission is active, or where it stands. Cold start, vague opener, or several active Missions. | `orient` |
| [prepare.md](references/prepare.md) | The work is not frozen yet. The problem is open, approaches compete, or a Mission must be drafted and previewed. | `explore`, `propose`, `plan` |
| [execute.md](references/execute.md) | A Mission is ready to activate, or an already-active one has work left to do. | `start`, `resume` |
| [runtime.md](references/runtime.md) | The work leaves this session — delegation, a subagent, another operator, or unattended Autopilot. | `handoff`, Autopilot |
| [close.md](references/close.md) | The work looks done and its claims now need assessment, Evidence, review, or owner acceptance. | `review`, `complete` |
| [audit.md](references/audit.md) | A claim is doubted, or someone asks whether a stated result is genuinely proven. Retrospective, not lifecycle. | `audit` |

Two rules for the ambiguous case:

- When the state is unclear, load `orient.md` first. It is read-only and cheap.
- Load exactly one. If the job changes mid-session, finish the current one and route again.

## Divide meaning from mechanics

A simple active Mission is one file: `<mission>/MISSION.md`. Its frontmatter holds
the frozen parts; its Markdown body holds the explanation. Split into
`objectives/` or `runs/` only when detail, delegation, ownership, or an
independent boundary earns it.

**The plan carries meaning** — outcome, criteria, scope, authority, rationale.
You may draft or edit these canonical files directly when that is the fastest
clear path. That is normal during exploration and during a declared bootstrap.

**The tooling carries repeatability** — schema validation, identity and
fingerprints, dependency integrity, atomic transitions. A Mission never picks
which mandatory validators run; the active schema registry decides that.

Which one to use:

| Use tooling when | Use judgment when |
|---|---|
| failure is expensive | meaning depends on context |
| the rule is exact and repeated | the prose is the value |
| the transition must be atomic | several answers are valid |
| | encoding it mechanically costs more than checking the result |

Full field lists, and the rules for splitting a Mission, live in
[mission-anatomy.md](references/mission-anatomy.md).

## Ask owner questions

Ask only when one of these is still open: a semantic fork, a Mission boundary, an
authority question, a risk, an irreversible effect, or a conflict with the current
Contract.

When you ask, use four short parts:

1. Plain outcome.
2. Technical basis.
3. Concrete options as `action -> consequence`.
4. Recommended default, and why.

Do not run an interview when the Mission already answers the question. Turn silent
product assumptions into explicit choices or Gaps. Carry accepted defaults forward.

## Authority and execution

**Owner only.** Outcome, completion criteria, semantic scope, review
independence, forbidden-effect ceiling. Record the owner and the activation time,
plus a fingerprint over that frozen semantic envelope. Objective and Run progress
stays mutable and outside the fingerprint.

**Operator freely.** Reversible implementation attempts, tools, checks, bounded
repairs — all inside the Mission.

**Return to the owner.** Scope expansion, irreversible or provider effects,
exhausted repair, any stop condition.

Execute in this order:

```
Mission card -> current Objective -> exact sources -> work -> focused checks
```

Batch cohesive work, then run one full tree-bound gate before independent review or
completion.

Fan out only when the slices are outcome-sized and disjoint. Give each one exact
inputs, dependencies, authority, stops, and a return contract. No recursive critic
loops. No tiny sessions.

Keep these five separate: Evidence, deterministic checks, independent review, owner
acceptance, completion. A green check proves its own observation and nothing more.

## Continuity

Always return the state a cold session would need:

- Mission, owner, current Objective and Run
- Contract and Git baseline
- validation mode
- review and Evidence state
- remaining dependencies and Gaps
- repair use
- recovery point
- exactly one safe next action, or one owner gate

When Spectacular develops itself, an active Mission keeps the schema and completion
boundary frozen at its activation. Later product changes apply to later Missions.
