---
name: spectacular
version: 2.4.0
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
| Anchor | Accepted project truth. Core Triad: `PROJECT.md` (root scope/boundaries), `STACK.md` (tools/runtimes), `ARCHITECTURE.md` (layers/components). On-Demand: `VOCABULARY.md`, `SECURITY.md`, `GUARDRAILS.md`, `PRODUCT.md` |
| PRD / Starter Input | Ephemeral kickoff input; digested into Core Anchors and `M1-bootstrap` claims, never kept as eternal context bloat |
| Campaign | Strategic roadmap arc (synonyms: *Initiative*, *Milestone Arc*, *Flight Plan*, *Theme*); 4–10 topological Campaign Blocks visualized via Mermaid (ASCII fallback) |
| Proposal | Optional, mutable exploration in chat, an issue, or a Spectacular file |
| Mission | Frozen execution plan and primary entry point in `M<N>-<slug>.md` |
| Contract / specification | Modular, composable capability specifications (`CC-<module>`); edit as ordinary Mission work |
| Decision | ADR-like durable choice and rationale, never routine lifecycle approval |
| Objective / Run | Inline while simple; promoted to a file only when independently useful (`O<N>-<slug>.md`, `R<N>-<slug>.md`) |
| Evidence / review | Earned proof and assessment; neither silently changes success criteria |

## Start every workflow

0. **Confirm the mechanical layer is present.** Run `spectacular --version`
   once, before anything else.

   If it resolves, proceed normally. If it does not, say so plainly and stop
   before any governed mutation. Do not emulate a missing CLI: never hand-write a
   record that a command owns, never invent a fingerprint, and never describe an
   edit as atomic when it was a plain file write. A user who believes they have
   freeze points and does not is worse off than one who knows the tooling is
   absent.

   Without the CLI you may still read records, explain the method, and draft a
   Mission plan for later activation. You may not start, promote, finish,
   complete, amend, or record — those produce fingerprints and transactional
   writes that only the binary provides.

   Onboard the user in one step rather than a lecture:

   > Spectacular's CLI is not on PATH, so I can read and draft but cannot run
   > governed execution. Install it from a verified release directory:
   >
   > ```sh
   > install/install.sh install --prefix /absolute/prefix --source /absolute/release --runtime claude
   > ```
   >
   > Then re-run `spectacular --version`. See [docs/installation.md](../../docs/installation.md).

   The Skill and the CLI are distributed separately and update separately: the
   Skill travels with the plugin, the binary is installed from a release
   directory. A host that installed only the Skill is expected and supported —
   it is a reduced mode, not a broken one.

1. **Find the workspace.** Read `.spectacular/PROJECT.md`. Only if that file is
   missing, read root `PROJECT.md`. If neither exists, this is a greenfield
   workspace: route to One-Shot Genesis in [prepare.md](references/prepare.md).

   That file is the root Anchor. Its `current_truth` field names the other Anchors
   and Contracts by id. Read those, and skip the ones this Mission does not touch.

2. **Run one read-only preflight.** Check:
   - workspace, and which Mission is selected
   - Git: branch, worktrees, clean tree, fresh against upstream/default, release state
   - Contract and baseline bindings
   - owner, activation time, activation fingerprint
   - validation mode
   - current Objective and Run
   - blockers: dependencies, Gaps, stops

   Report three lines: plain outcome, one technical evidence line, one next action.

   Answer two isolation questions, not one. **Which branch** decides what the
   history looks like; **which working tree** decides whether a concurrent session
   can destroy the work. A branch alone shares one tree, so `checkout`, `stash`,
   and `reset` still reach everything.

   If the preflight lands on the default branch and the job is multi-step, branch
   before activating — `git checkout -b <mission-slug>`. Working a Mission on
   `main` leaves no merge point and no review boundary. The exception is an
   owner-requested quick patch with no concurrent session; take it explicitly, not
   silently.

   Run `git worktree list` as part of this check. If a second session is live — a
   reviewer, a feedback session, another Mission — take a worktree rather than
   sharing the tree. A branch separates history; a worktree separates hands.

3. **Enter the Mission through `M<N>-<slug>.md`.** Read its frontmatter control card
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
| [prepare.md](references/prepare.md) | The project is new (Genesis), or the work is not frozen yet. The problem is open, approaches compete, or a Mission must be drafted and previewed. | `genesis`, `explore`, `propose`, `plan` |
| [execute.md](references/execute.md) | A Mission is ready to activate, or an already-active one has work left to do. | `start`, `resume` |
| [runtime.md](references/runtime.md) | The work leaves this session — delegation, a subagent, another operator, or unattended Autopilot. | `handoff`, Autopilot |
| [close.md](references/close.md) | The work looks done and its claims now need assessment, Evidence, review, or owner acceptance. | `review`, `complete` |
| [audit.md](references/audit.md) | A claim is doubted, or someone asks whether a stated result is genuinely proven. Retrospective, not lifecycle. | `audit` |

Two rules for the ambiguous case:

- When the state is unclear, load `orient.md` first. It is read-only and cheap.
- Load exactly one. If the job changes mid-session, finish the current one and route again.

## Divide meaning from mechanics

A simple active Mission is one file: `<mission>/<mission>.md` (e.g. `M5-implement-compact-missions/M5-implement-compact-missions.md`). Its frontmatter holds
the frozen parts; its Markdown body holds the explanation. Split into
`objectives/` or `runs/` only when detail, delegation, ownership, or an
independent boundary earns it.

**The Anchor naming rule**: Bare single-word uppercase names (`<NOUN>.md` e.g. `PROJECT.md`, `STACK.md`, `ARCHITECTURE.md`, `README.md`, `AGENTS.md`) are reserved exclusively for Project Anchors and workspace landmark contracts. All governed records (Missions, Runs, Objectives, Proposals, Reviews, Decisions, Evidence, Gaps) carry their scoped prefix in their filename.

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

### Ask for authorization, not for labor

An owner-gated action needs the owner's *decision*, not the owner's *hands*. When
you reach one, ask to be authorized and say what you will do with the authorization:

> `mission complete` is owner-gated because completion freezes the claims.
> Authorize me to run it? (Y/N)

Do not hand back a command to type when you could hold the keyboard on an explicit
yes. Reserve owner-executed steps for actions where performing it *is* the
authority — credential entry, an interactive login, a decision the record must
attribute to the owner by name.

When you do act on an authorization, record who authorized and who performed. An
operator acting on owner approval is not the owner acting.

### Settle execution mode at activation

Activation is also when to ask how much the owner wants to be involved while the
Mission runs: autopilot, checkpoints, or a named human-in-the-loop moment. Ask it
once, alongside the activation gate, and name the actual checkpoints rather than
offering the word. Owner gates and stops fire in every mode.

See [execute.md](references/execute.md) for the three modes and what makes an
answer usable.

### Batch gates; ask once

Collect every owner decision a phase needs and put them in one exchange. Six small
approvals across six messages cost more than one message listing six items.

Before asking, check whether the answer is already available: an earlier
authorization that plainly covers this case, a stated preference, a default the
Mission froze. An approval carries forward to the same kind of action in the same
phase. Re-asking something already answered is noise, not diligence.

### State a boundary once

Name an authority boundary the first time it becomes relevant, then act on it.
Repeating the same constraint at every step reads as stalling, and buries the one
place a decision is genuinely needed.

When you must decline, be brief: one sentence on what you cannot do, one on the
nearest thing you can, then continue with everything not blocked. Do not
re-litigate a boundary the owner has already heard, and do not narrate the refusal
at greater length than the work.

## Authority and execution

**Owner only.** Outcome, completion criteria, semantic scope, review
independence, forbidden-effect ceiling. Record the owner and the activation time,
plus a fingerprint over that frozen semantic envelope. Objective and Run progress
stays mutable and outside the fingerprint.

**Operator freely.** Reversible implementation attempts, tools, checks, bounded
repairs — all inside the Mission.

**Return to the owner.** Scope expansion, irreversible or provider effects,
exhausted repair, any stop condition. Returning means requesting a decision, then
proceeding on the answer — not transferring the remaining work. Finish everything
the gate does not block before you ask, so one answer completes the phase.

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
