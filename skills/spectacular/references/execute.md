# Start and resume

Use this when: Orchestrator or primary operator activating, resuming, or executing a Mission; detailed Git isolation.

## Read-Only Preflight Checklist

Run one read-only preflight before governed execution:
- **Workspace & Mission**: Confirm `.spectacular/PROJECT.md` and which Mission is selected.
- **Git State**: Branch, worktrees (`git worktree list`), clean tree, fresh against upstream/default, release state.
- **Bindings**: Contract fingerprint and baseline commit binding.
- **Identity**: Owner, activation time, activation fingerprint.
- **Validation Mode**: Supported CLI mechanics vs. declared `manual-bootstrap`.
- **Execution State**: Current active Objective and Run.
- **Blockers**: Upstream dependencies, unresolved Gaps, active stops.

### Report in Three Lines:
1. **Plain outcome**: Current project direction, selected Mission, and lifecycle status.
2. **Technical evidence**: Git branch/worktree, commit SHA, Contract fingerprint, validation mode.
3. **Next action**: Exactly one safe next action, or one owner gate.

## Start

Start only an owner-confirmed Mission. Before activating, confirm all of it:

- the design is sufficient, with no blocking Gap
- exact Contract and Git baseline
- complete claim coverage
- coherent Objectives
- explicit authority and scope
- review level
- budgets and stops

Then record the owner and the activation time, plus a fingerprint. All three land
in the Mission's `activation:` block:

```yaml
activation:
    at: "2026-08-16T14:19:42Z"
    by: Alex
    fingerprint: sha256:ef7e695fc1f496fabb92c05a0369ca148fc6404e49b728521db3a2cdf120f389
```

**The fingerprint covers the frozen semantic envelope** — outcome, review,
completion, authority, scope, budgets, dependencies, Gaps, stops. It must not
cover mutable state: status, Objective progress, Run state, repair count.

The tooling computes it. Do not hand-roll a hash to make a check pass.

A normal start creates one file:
`.spectacular/missions/<slug>/<ref>-<slug>.md`, with inline Objectives and R1.

### Settle the execution mode in the same breath as activation

Activation authorizes the plan. It does not say how often the owner wants to be
interrupted while it runs — and asking that question later, one Objective at a
time, is how a Mission turns into homework.

Ask it once, with activation, as a single question offering three modes:

| Mode | The owner is asked | Fits |
|---|---|---|
| **Autopilot** | only at owner gates and stops | mechanical work, reversible steps, a plan the owner already read closely |
| **Checkpoints** | at named points, briefly | most Missions; the default when unsure |
| **HITL** | at a specific named human activity | work that genuinely needs the owner's hands or eyes at a known moment |

Two rules make the answer usable:

- **Restate the checkpoints when offering the mode.** "Checkpoints" means nothing
  until the owner sees which ones — name them, one line each, from the Objective
  graph: end of each Objective, end of each Run, or a specific boundary. Three to
  five is a checkpoint list; every Objective is not.
- **For HITL, name the activity and the moment.** Not "I'll check in during O3" but
  "before deleting `internal/context`, you confirm the dependency walk". If you
  cannot name what the human does, it is a checkpoint, not HITL.

Record the answer in the Run body and honor it for the whole Mission. Owner gates
and stop conditions still fire in every mode — autopilot means fewer interruptions,
never fewer gates. An owner who chose autopilot has not authorized scope expansion,
an irreversible effect, or working past a stop.

If the mode was never settled, default to checkpoints at Objective boundaries and
say that you are doing so.

### Plan checkpoints as flexible Run-body gates

A checkpoint is a planned place to review progress, run a check, collect a
decision, or decide whether to resume. It is optional, lives in the Run body,
and is not automatically a human-review gate or a durable authority record.

Plan only the checkpoints that alter how the Run proceeds. Use this routing
table when one produces a durable result:

| Checkpoint result | Record to create or update |
| --- | --- |
| owner choice or changed direction | Decision |
| observation, command output, or test result | Evidence |
| verdict on implementation or claims | Review or Assessment |
| transfer to another operator or runtime | Handoff |
| ordinary progress check with none of the above | Run-body note only |

Use a compact body shape:

```md
### Checkpoint: after O2 integration
- Trigger: focused integration checks complete.
- Reviewed: <what changed and which checks ran>.
- Result: <continue, repair, stop, or owner gate>.
- Next: <one action or linked durable record>.
```

Checkpoint records under `runs/.../checkpoints/` remain available for historical
or advanced workspace layouts. Normal v2 execution does not create one for an
ordinary check-in.

Two things that are not activation authority:

- A Proposal is optional input. Mission start never creates one.
- A Decision is not activation authority. It records a choice; it does not
  authorize a start.

Use the typed tooling for atomic generation and validation:

```bash
spectacular mission start <plan.md|-> --json   # mutating; generates the canonical path
spectacular mission check <ref> --json         # read-only; verify after
```

`mission start` takes the approved plan as a Markdown file, or `-` for stdin. It
generates identities, bindings, and activation atomically — you do not author the
UUID or the fingerprint yourself.

Under `manual-bootstrap`, build the same canonical shape directly and verify it
with focused scripts — see [bootstrap.md](bootstrap.md). Never route the work
through an incompatible legacy command sequence to make it look validated.

## Resume

Read three things: the Mission card, the current Objective section, and the exact
pointers they name.

```bash
spectacular mission show <ref> --json    # card, current Objective, current Run
spectacular run show <ref>/<run> --json  # e.g. M7/R1, when the Run detail matters
```

Then recheck:

- Contract fingerprint
- Git baseline
- activation fingerprint
- validation mode
- authority and scope
- budgets
- dependencies and Gaps
- stops

A material semantic change goes back to the owner. Reversible implementation
changes stay with the operator.

## Pick the isolation the job needs

Two independent questions, and answering only the first is how work gets lost.

**Which branch does this commit belong on?** — decides what the history looks like.
**Which working tree do I edit in?** — decides whether a concurrent session can
destroy it.

A branch alone gives no protection from a second session: branches share one
working tree, and `checkout`, `stash`, and `reset` operate on all of it. A worktree
gives a second directory with its own checked-out branch, so two sessions cannot
reach each other's files.

| Situation | Branch | Worktree |
|---|---|---|
| One session, one Mission | Mission branch | no — the main tree is yours |
| Owner-requested quick patch, nothing else running | `main` directly | no |
| A second session runs while a Mission is in flight | separate branches | **yes, one per session** |
| Reviewing a Mission while its author keeps working | reviewer takes a worktree at the reviewed commit | **yes** |
| A subagent editing files you also edit | same branch, same tree | no — a worktree would fork the work |

The rule of thumb: **a branch separates history; a worktree separates hands.** Add
the worktree when two sets of hands are live at once, not merely when two topics
are.

### Two sessions at once

This is the case that needs setting up before it is needed, not during. A Mission
session and a feedback session sharing one tree will collide the first time either
switches branches.

```bash
git worktree add ../<repo>-<purpose> <branch>   # second directory, own branch
git worktree add ../<repo>-review <commit>      # detached, for reviewing a fixed tree
git worktree list                               # what exists now
git worktree remove ../<repo>-<purpose>         # when done
```

Each session works in its own directory and never `cd`s into the other. Both share
one `.git`, so history, branches, and merges behave normally.

Split the work by surface, and the branches almost never touch the same files:

- **Mission code** — `internal/`, `cmd/`, `test/` — on the Mission branch.
- **Skill, feedback, and docs** — `skills/`, `FEEDBACKS.md`, `TODO.md`, `docs/` —
  on `main`. These record what the Mission is teaching you while it runs, and
  holding them hostage to the Mission's review is what makes them never get written.

When a worktree is live, say which tree you are in before any Git operation. "On
branch X" is not enough information when three trees exist.

### Worktree hazards

- **A branch can only be checked out in one tree.** While `main` is checked out in
  a worktree, the primary tree cannot switch to it. This is a feature — it is what
  makes the isolation real — but it surprises the first time.
- **Stale worktrees keep old code alive.** Agent worktrees under `.claude/` hold
  full copies at whatever commit they were created. After deleting a package, those
  copies still contain it, and a repository-wide grep will find it. Exclude
  non-module directories from structural checks, and remove worktrees the moment
  they have served their purpose.
- **`git worktree list` before assuming.** A tree you did not create may already
  hold the branch you want.

## Branch before you activate

**A Mission gets its own branch. Create it before `mission start`, not after.**

Multi-step work on `main` has no merge point, no review boundary, and no cheap
way back. The Mission's own authority block gates `merge` for the owner, which
presupposes there is something to merge *from*; activating on `main` quietly
removes that gate by removing its subject.

```bash
git checkout -b <mission-slug>      # before mission start
```

`mission start` records the current branch into `baseline:`, so the branch must
exist first. A Mission whose `baseline: branch: main` is a Mission that skipped
this step — check for it when resuming, and say so rather than continuing
silently.

The one exception: the owner asks for immediate action or a quick patch, **and**
no concurrent session is running. Then patching `main` directly is allowed. State
that you are taking the exception and why. Anything multi-step is not a quick
patch, whatever it was called when it started.

Retrofitting a branch after commits have landed on `main` is possible but not
free — `git branch <name> <sha>` then `git reset --hard <base>` moves the commits,
and the reset **discards every uncommitted change in the working tree**. Commit or
stash first.

## Choose the branch before you edit

A Mission's work belongs on the branch its `baseline:` names. Before starting a
side task — a doc pass, a rename, a cleanup — check whether an active branch
already modifies those files:

```bash
git status --short <paths>                       # is this work already in flight here?
git branch --format='%(refname:short)' | while read -r b; do
  git diff --stat main "$b" -- <paths>           # does another branch already touch them?
done
```

If an active branch already edits those files, **work on that branch.** Splitting
the same files across two branches guarantees a manual conflict resolution later,
and the resolution is where work gets silently reverted.

Start a separate branch only when the file sets are genuinely disjoint. Isolation
buys a clean review only when nothing else is editing the same lines.

If both branches have already diverged on the same file, reconcile on the feature
branch — `git merge main` there, resolve, verify — and only then move it onto
`main`. Never resolve a content conflict directly on `main`.

## Work in outcome-sized clusters

```
[claims + dependencies] -> [work] -> [focused checks] -> [boundary integration] -> [local commit]
```

- Run one full repository and release gate after integration (e.g. host project test suite, or `bash test/verify.sh all` in this repo).
  Do not repeat it per small edit.
- Read detailed logs only when something fails.
- Keep Git, secret, and distribution checks at the boundary where they apply:
  commit, push, PR, or release.
