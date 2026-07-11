---
type: idea
status: parked
priority: medium
owner: alex
updated: 2026-07-11
origin: user pain (2026-07-11) — spec-builder runs for hours with no orchestrator visibility
promoted_to: null
related:
  - ../../skills/spectacular/references/build-workflow.md
  - ../AGENTS.md
  - ../../skills/spectacular/references/tasks-rules.md
  - ./builder-trace-and-fanout.md
---

# Idea — size-and-decompose gate (visible sub-steps inside a milestone)

## The pain (user, 2026-07-11)

> "Sometimes I see [a spec-builder] running for hours, with no knowledge of what it's doing
> from the main orchestrator. We should have the orchestrator think how to distribute the task
> into one or more agents or sub-steps, like internal checkpoints it reviews."

## Diagnosis — two problems, one chosen fix

- **A — decomposition stops at the milestone boundary.** Every dispatch path in `build-workflow.md`
  assumes the unit is a *whole milestone*. A single milestone can be a "small PR" (the doc's words)
  that runs for hours as one opaque block. There is no step that says "this milestone is too big —
  break it into sub-steps and checkpoint between them."
- **B — dispatch is fire-and-forget by design.** Step 2b spawns a `spec-builder` and blocks until it
  returns one atomic result. A 5-min and a 3-hour build look identical to the orchestrator. True
  in-flight streaming bumps into the platform's fire-and-forget subagent model (this is what
  [[builder-trace-and-fanout]]'s deferred M3 trace-folder was gated on).

**Chosen fix (grill, 2026-07-11): "smaller units + checkpoints" (A).** Decomposing a fat milestone
into sequential sub-steps fixes BOTH: nothing runs for hours as one block (A), and the checkpoints
between sub-steps *are* the visibility (B) — without needing to crack open a running agent. It
sidesteps the platform constraint entirely and is mostly a doc + policy change. This is the
smaller-slice version of M3 and may make M3 permanently unnecessary.

## Key realization — the substrate already exists

Nothing new needs building. Two checkpoint surfaces already exist; `build-workflow.md` just never
tells the orchestrator to use them when a milestone is fat:

1. **Nested `- [ ]` bullets in TASKS.md** — already documented (`tasks-rules.md`: "indented sub-bullets
   allowed as a nested acceptance checklist, not counted"). Reused here as **durable decomposition
   checkpoints** the orchestrator fills when it decides to break a milestone down.
2. **Harness `TaskCreate`/`TaskUpdate`** — already documented in `AGENTS.md` § Task tracking as the
   finer layer that "decompose[s] a single TASKS.md milestone into the concrete edits/commits/tests."
   This is the **live in-session progress signal** (drives the CLI UI). The visibility channel.

The gap is purely the **trigger to invoke this inside the build arc** — a missing rung, one level
below the milestone-level gates the workflow already has.

## Decisions locked (grill, 2026-07-11)

- **Fix = decomposition, not live trace** — turns B into A; no platform fight.
- **Sub-steps live in ephemeral nested TASKS bullets** — zero new substrate, no doctor area, no
  trace folder. Rejected: reviving M3's `build/<m>/` trace folder (heavier, was gated as unearned);
  a pure no-artifact habit (loses the durable checkpoint record).
- **D1 — a new `build-workflow.md` Step 1.5 (size-and-decompose gate)** between Step 1 (worth-it) and
  Step 2 (build/dispatch). Distinct decision from dispatch-vs-inline, so it's its own step.
- **D2 — doc step + a lightweight `@Implementation` warn policy** `decompose-large-milestone`, so the
  rule is injected at runtime, not just a poster (Principle 3 enforcement pattern). Warn, not block —
  sizing is judgment.

## The buildable shape

### Step 1.5 — Size-and-decompose gate (new, in build-workflow.md)
Once Step 0 closes the brief and Step 1 picks inline-vs-dispatch, ask: **is this milestone one
coherent pass, or several?** A milestone is "large" when it spans multiple phases that each have
their own verify point — e.g. schema → CLI wiring → tests, or parser → renderer → doctor check.

If multi-phase:
- **Decompose into sequential sub-steps** as nested `- [ ]` bullets under the milestone's TASKS block
  (and mirror them as harness `TaskCreate` items for the live signal).
- **Build/dispatch one sub-step at a time**, verifying + reporting at each boundary before starting
  the next. Each sub-step is a short, visible unit — not a 3-hour opaque block.
- For a *dispatched* milestone: send **one closed sub-brief per sub-step in sequence**, confirming
  each returned diff before dispatching the next — never one fat brief that runs unwatched. (Ordered
  sub-steps within a milestone are the same-file/ordered case: serialize, don't parallelize.)

If single-phase: proceed to Step 2 as today. Most small milestones skip this gate.

### Policy `decompose-large-milestone` (@Implementation, warn)
```
### decompose-large-milestone
- principle: 10
- severity: warn
- check: a milestone that spans multiple verify-points is built/dispatched as sequential
  sub-steps (nested TASKS bullets + harness tasks) with a checkpoint between each — not as one
  unbounded pass

A milestone that spans several phases (each with its own verify point) is built one visible
sub-step at a time, confirming at each boundary — not dispatched as a single opaque brief that
runs for hours. Decompose into nested `- [ ]` checkpoints, build/dispatch sequentially, report
between. **Override:** a genuinely single-phase milestone needs no sub-steps — don't manufacture
ceremony for a coherent one-pass change.
```
*(Authored legibly per the v1.30.1 patch: warn ⇒ no ⛔ marker; carries an **Override:** clause.)*

## Scope / cost
Small. One new build-workflow step + one mirror line in bug-workflow's Step 3 (a fix can also be
multi-phase, though rarer) + one `@Implementation` policy + a note in `AGENTS.md` § Task tracking
tying the nested-bullet checkpoint use to the two-layer model + a `tasks-rules.md` line blessing
nested bullets as decomposition checkpoints (not only acceptance criteria). No CLI change, no new
substrate, no doctor area.

## Explicitly out of scope (deferred)
- **Live in-flight trace from a running agent (problem B's hard form)** — stays parked in
  [[builder-trace-and-fanout]] M3. Only revisit if decomposition proves insufficient (a *single*
  sub-step still runs too long to be visible).
- **Any git worktree / parallel sub-step machinery** — ordered sub-steps serialize by definition.

## Suggested next step
Cut a narrow request "milestone-decomposition" shipping the Step 1.5 gate + the policy + the two
doc notes. Independent of stance-layer (b26). Verify by walkthrough: a synthetic multi-phase
milestone produces nested checkpoints + per-step reports; a single-phase one skips the gate.
