---
type: Mission
id: 01a00d6e-df08-75d3-a6c3-5975bb421630
title: Remove dead v1 surface and close the deferred cleanup ask
status: active
created: "2026-08-17T10:13:17Z"
updated: "2026-08-17T10:24:45Z"
activation:
    at: "2026-08-17T10:13:17Z"
    by: Alex
    fingerprint: sha256:a7ae29ba1e6416f466ef73acf40b9688d89e7ab56e71265e8608c99a7324311c
authority:
    operator:
        - inspect
        - edit-in-scope
        - choose-reversible-implementation
        - run-checks
        - generate-derived-files
        - bounded-repair
        - commit-local
    requires_owner:
        - activate-mission
        - change-outcome-or-completion
        - expand-scope
        - push
        - merge
        - release
        - irreversible-change
        - destructive-data
        - secret-change
baseline:
    branch: main
    commit: 661fccab9511cf701e34d4e263ceff80c1c2d31b
completion:
    - claim: dead-surface-removed
      pass_boundary: The unreachable v1 context-compiler chain — `internal/context`, `internal/projection`, and `internal/guardrails` — is deleted as one unit rather than piecemeal, the build and the full verification suite pass without it, and no remaining package imports a deleted one. Partially-live packages lose only their unreachable members; a symbol that any live caller reaches is retained even when it looks like v1, so `internal/governance` is pruned rather than removed. Before deletion, the capabilities that die with the chain are written down as a recorded follow-up, so the cleanup discards code without discarding the ideas the code was the only record of.
      proof_requirement: A dependency walk from the `cmd/spectacular` main package proves the three deleted packages are absent from its transitive closure before removal, and that every retained governance symbol has a named live caller; `bash test/verify.sh all` passes after deletion; a test asserts the deleted import paths appear nowhere in the module; the follow-up record names the conflict reporting, omission reporting, and loaded-versus-available record counts that no v2 surface currently provides.
    - claim: next-action-reads-reviews
      pass_boundary: The derived next action consults the reviews the bundle carries. Once every Objective is implemented and a review bound to the current activation fingerprint records a passing verdict, the next action names owner completion; a review bound to a stale fingerprint, or carrying a failing verdict, continues to ask for a review. No projection writes to the canonical tree.
      proof_requirement: Fixtures assert the next action across four states — no review, passing review on the current fingerprint, passing review on a stale fingerprint, and failing review — and a mutation that ignores the reviews list is caught by the suite; the M7 golden fixtures still render byte-identically where no review exists.
    - claim: working-tree-accounted
      pass_boundary: Every path in the working tree is tracked, ignored by a rule that states why, or removed. `_snapshots/` has a recorded decision on whether it is local-only recovery or travels with the repository, and the decision is enforced by a rule rather than by memory. No file is deleted without the owner naming it.
      proof_requirement: '`git status --porcelain` is empty on a clean tree; a test or check asserts no untracked, non-ignored path exists; the `_snapshots/` decision is written into the record and the corresponding `.gitignore` rule carries an explanatory comment.'
    - claim: gaps-swept
      pass_boundary: Every open Gap is either closed with the change that closes it named, restated with a current reason to stay open, or converted into a recorded follow-up. A Gap is never closed by deleting it. `dead-v1-governance-code` reaches a terminal state consistent with what this Mission actually removed.
      proof_requirement: Each open Gap at Mission start is listed with its resolution and the record that carries it; a check asserts no Gap references a deleted package or a path that no longer exists.
contract:
    fingerprint: sha256:1ffd39b498b44dce4e77cdf902f5f827bdf40eb2b317573c38a405f9b9ae9a0b
    ref: Contract:01a00aae-8921-7b27-96a9-1a4c175e7dc6
dependencies:
    - M8 completed with a recorded review and owner acceptance.
fallbacks:
    - approach: Delete `internal/projection` alone and leave `internal/context` and `internal/guardrails` in place.
      invalidated_if: A live caller for `internal/context` appears, in which case the chain is load-bearing and none of the three may be deleted.
      rejected_because: '`internal/context/compiler.go` imports both `projection` and `guardrails`, so removing only `projection` breaks the build. The three form one context-compiler chain and are unreachable only as a unit.'
    - approach: Keep the context compiler unwired rather than deleting it, on the chance v2 wants conflict and omission reporting later.
      invalidated_if: A near-term Mission is planned that would rewire the compiler as-is rather than reimplementing the idea against the v2 model.
      rejected_because: Unwired code is not a design record. It compiles, so it looks maintained, and it drifts silently against every schema change without a test noticing. Git history preserves it, and a written follow-up preserves the idea in a form a reader will actually find.
    - approach: Fix the next action by comparing the review count to zero rather than checking fingerprint binding.
      invalidated_if: Reviews are made to re-bind automatically on re-activation, which would make a stale binding unreachable.
      rejected_because: A review bound to a stale fingerprint would retire the instruction to review, which is exactly the drift the fingerprint exists to catch.
    - approach: Track `_snapshots/` in the repository rather than ignoring it.
      invalidated_if: A review or release is found to require a snapshot as evidence, which would make local-only storage a loss of record.
      rejected_because: Snapshots duplicate what git history already holds, and TODO records an open question about whether the mechanism earns its keep at all.
gaps: []
objectives:
    - claims:
        - dead-surface-removed
      id: 01a00d6e-df08-7f59-81fe-aeaec3fd2bc9
      outcome: Prove reachability, record what the context compiler could do that v2 cannot, then delete the three-package chain as a unit and prune only the unreachable governance members.
      ref: O1
      status: implemented
    - claims:
        - next-action-reads-reviews
      id: 01a00d6e-df08-70ab-9b64-939b8a11f63e
      outcome: Make the derived next action read the bundle's reviews, distinguishing a review bound to the current fingerprint from a stale one.
      ref: O2
      status: implemented
    - claims:
        - working-tree-accounted
      id: 01a00d6e-df08-7465-8ec6-b67c41f151da
      outcome: Account for every untracked working-tree path and record the `_snapshots/` retention decision as an enforced rule.
      ref: O3
      status: implemented
    - after:
        - O1
        - O3
      claims:
        - gaps-swept
      id: 01a00d6e-df08-7ecf-bec1-cc525345bca6
      outcome: Sweep every open Gap to a terminal state consistent with what was removed.
      ref: O4
      status: implemented
outcome: The repository carries only reachable code, the derived next action agrees with the reviews the record holds, and every untracked working-tree path is either tracked with a stated reason or deliberately ignored, closing the cleanup ask M8 recorded as deferred.
owner: Alex
ref: M9
repair_budget: 3
request:
    asks:
        - ask: Review the open gaps and dead weight and perform a repository cleanup
          claims:
            - dead-surface-removed
            - working-tree-accounted
            - gaps-swept
          disposition: covered
        - ask: Fix the NEXT line that asks for a review that already exists
          claims:
            - next-action-reads-reviews
          disposition: covered
        - ask: Move the matrix-proposal-loop material out of spectacular for review
          disposition: declined
          reason: Already done before this Mission was planned; the files are staged in skills_db for owner review and are no longer in this repository.
        - ask: Decide whether TODO.md remains local-only and gitignored
          disposition: deferred
          reason: A durable-follow-up policy question that outlives a cleanup Mission; it needs its own decision rather than a cleanup side effect.
        - ask: Reconsider whether the snapshot mechanism earns its keep versus git history
          disposition: deferred
          reason: This Mission records the retention decision but does not relitigate the mechanism; removing it is a schema change, not cleanup.
    captured_at: "2026-08-17T00:00:00Z"
    source: chat, session continuing from M8 completion
review: independent
run:
    current_objective: O1
    id: 01a00d6e-df08-72cb-9aa0-37920de5735c
    operator: Alex
    ref: R1
    repairs: 0
    started_at: "2026-08-17T10:13:17Z"
    status: active
scope:
    mechanical:
        - cmd/spectacular/
        - internal/
        - test/
        - skills/spectacular/generated/
        - .spectacular/
        - .gitignore
    semantic:
        - Removal of the unreachable v1 context-compiler chain (`internal/context`, `internal/projection`, `internal/guardrails`) and of unreachable `internal/governance` members.
        - A recorded follow-up naming the reporting capabilities lost with the deleted chain.
        - Derivation of the next action from recorded reviews and their fingerprint binding.
        - Working-tree accounting and the `_snapshots/` retention decision.
        - Terminal resolution of open Gaps.
start_key: sha256:ee854ade160eed19fa89d6128db97d0d22854201029e7ee26111fd0d43c7d939
stops:
    - The accepted command surface grows beyond the ten commands CC-missioncli enumerates, or a new noun is introduced.
    - A projection writes to the canonical tree or caches a rendered view.
    - A package is deleted before a dependency walk proves it is absent from the main package's transitive closure.
    - A capability is deleted without a recorded follow-up naming what was lost.
    - A Gap is closed by deletion rather than by a named resolution.
    - A completed Mission's frontmatter is rewritten to satisfy a new field.
    - A file the owner has not named is deleted from the working tree.
validation:
    mode: cli
    schema: mission.v2
---
# Remove dead v1 surface and close the deferred cleanup ask

M8 recorded "review the open gaps and dead weight and perform a repository
cleanup" as `deferred` — the first ask Spectacular has ever tracked as dropped.
This Mission closes it. Leaving the first recorded deferral unresolved would make
the request record a place asks go to be forgotten politely.

## What the dead code is, and why it is one unit

`TODO.md` records `internal/projection` as 917 dead lines "referenced by neither
`internal/command` nor `internal/missionbundle`". True, and incomplete in a way
that matters: `internal/context/compiler.go` imports both `projection` and
`guardrails`, and nothing outside `internal/context` imports `internal/context`.
A dependency walk from `cmd/spectacular` confirms all three are absent from the
main package's transitive closure.

The three are v1's **context compiler**: `guardrails` supplied declared guidance,
`projection` built cards and pointers over the workspace, and `context` assembled
them into a bounded, fingerprinted `Bundle` answering "what should be loaded right
now". Its own package comment states the discipline this repository still keeps:
the output "is a disposable projection and never owns Mission or Contract truth."

It is unused because v2 answered the same question more cheaply. `Bundle.Derive()`
computes state on read, directly from the record, and it was placed in
`missionbundle` beside the Bundle it reads rather than in the package named for
projection. Once that choice was made the chain lost its only consumer, and
nothing ever unwired it. This is supersession, not abandonment.

`internal/governance` is a different case and the TODO overstates it. It *is*
reachable from main. The proposal and `candidate_*` machinery is unreachable, but
`ApplyTransaction`, `FileChange`, and `RecoverTransactions` are live in
`internal/command` and `internal/missionbundle/service.go`. The package stays;
only the unreachable members go.

This is why the pass boundary demands the dependency walk before deletion rather
than after. A cleanup Mission that breaks the build has done the opposite of
cleanup.

## What dies with it, and why that is written down first

Two capabilities have no v2 equivalent, and deleting the code deletes the only
record that they were ever considered:

- **Conflicts and omissions.** The Bundle reported what it could not reconcile
  and what it deliberately left out. Nothing in v2 states its own limits this way.
- **Loaded versus available record counts.** The Bundle reported that it loaded
  twelve of forty records, making a bounded-context claim checkable rather than
  asserted.

Neither justifies keeping unwired code — code that compiles looks maintained,
drifts silently against every schema change, and no test notices. Git history
preserves the implementation. O1 writes the ideas into a recorded follow-up
before the deletion, so the salvage happens in a form a reader will find.

## Why the next action is a design question, not a patch

`nextAction` selects "record a review" on `state.Done == len(b.Objectives)`
alone and never reads `b.Reviews`, so a recorded review cannot retire the
instruction. That much is a bug.

The fix is not. A review binds to an activation fingerprint. If the boundary is
later amended, the fingerprint changes and the review no longer describes what
the record now claims. Whether a passing review survives a boundary amendment
decides whether the next action trusts it — and that is the difference between a
gate and a formality. This Mission answers it explicitly: a stale binding keeps
asking.

M7 froze `state-line` on the boundary that every field is derived from data the
bundle already carries, and listed as a stop that a state line may not disagree
with the record it summarizes. The stop did not fire because no fixture covered a
bundle holding both implemented Objectives and a recorded review. This Mission
adds that fixture.

## Why the working tree counts as cleanup

Four untracked paths sit in the tree with no Mission explaining them:
`_research/`, `_snapshots/`, `articles/`, and `.claude/settings.json`. Untracked
is not the same as ignored — an untracked path is invisible to the record and
silently lost on a fresh clone. Each one is either tracked with a reason or
ignored by a rule that says why.

`_snapshots/` carries an open question from `TODO.md`: whether ignored canonical
snapshots are intentionally local recovery artifacts, or whether any snapshot
must travel with a review or release. The decision gets recorded and enforced,
not remembered.

## What this Mission does not do

It does not remove the snapshot mechanism, decide `TODO.md`'s durability policy,
or touch the matrix-proposal-loop material already staged elsewhere for review.
Each is recorded in the request block with its disposition, so a reader can tell
what was set aside and why.

## Decisions taken (O3)

Four untracked paths were resolved. Each is now ignored by a rule in `.gitignore`
that states why, so the reason survives the session that decided it.

- **`_snapshots/` — local recovery only.** Snapshots duplicate what git history
  already holds. Nothing here travels with a review or a release; if a snapshot is
  ever needed as evidence, it is exported deliberately rather than by force-adding
  the folder. This answers the open question `TODO.md` recorded and closes it.
- **`articles/` — local-only, as it already declared.** `articles/AGENTS.md` says
  "Gitignored — never committed to the spectacular repo". The intent was written
  down but never enforced by a rule, which is precisely how it stayed untracked
  and unexplained. The rule now matches the declaration.
- **`_research/` — regenerable input, not product surface.** Saved NotebookLM
  queries, agent-memory reports, and source dumps. Inputs to thinking.
- **`.claude/settings.json` and `.claude/worktrees/` — per-machine state.** Shared
  project settings would belong in a tracked file; these are not that. The
  worktrees directory holds separate checkouts and is not part of this module,
  which is also why the deleted-import check skips it.

`TestWorkingTreeHasNoUnexplainedUntrackedPaths` enforces the general rule: an
untracked path is invisible to the record and silently lost on a fresh clone, so
"not committed yet" and "deliberately excluded" must never look the same.

## Gap sweep (O4)

Every Gap open at Mission start is listed below with its resolution. The Gaps live
in `CC-projsurf`, and this Mission is frozen against that Contract: editing it
refuses with `stale_fingerprint field contract.fingerprint: bound Contract content
changed`. That refusal is correct — it is the binding doing its job — so the sweep
is recorded here and the Contract is amended after M9 completes, when nothing is
bound to it. `TODO.md` records the same collision on 2026-08-16, when a note
written into one of these Gaps had to be reverted.

- **`dead-v1-governance-code` — closed by O1.** The Gap asked for "a decision on
  whether removing them belongs with the Proposal schema work or with a separate
  cleanup". Answer: a separate cleanup, this one. The Gap's framing was also too
  narrow. It named only `ProposalInput`, `CreateProposal`, and the `candidate_*`
  machinery; a dependency walk showed every governance *service* export
  unreachable — 3,432 lines across `service.go`, `closure.go`, `progress.go`,
  `typed.go`, and `models.go`, plus a 1,077-line test file. `transaction.go` is
  self-contained and stayed. The package was pruned, not removed, because
  `ApplyTransaction`, `FileChange`, and `RecoverTransactions` have live callers.
- **`lifecycle-diagram-ungenerated` — stays open, reason unchanged.** Still blocked
  on extracting a declarative transition model; transitions remain implicit in
  `service.go` control flow. Nothing in M9 touched it.
- **`concurrent-run-timelines` — stays open, reason unchanged.** Still blocked on a
  Run-model change spanning run start, fingerprints, atomicity, and review
  boundaries. Out of scope here.
- **`mission-ref-frontmatter-drift` — already closed** by this Contract before M9
  began, via M7's normalized decoder and M8's `mission-order-integrity`. Listed for
  completeness; no action taken.

No Gap was closed by deletion, and no Gap references a package or path this
Mission removed — `dead-v1-governance-code` names `internal/governance`, which
still exists.

### Found during the sweep, not closed

`internal/index` is unreachable: zero importers, not even from a test outside
itself. It is the v1 predecessor of `discovery.Workspace.Lookup`, carrying 8 tests
and no callers. M9's frozen scope names three packages, and adding a fourth is
`expand-scope`, so it was left in place and recorded in `TODO.md` instead. It
needs an owner decision or a follow-up Mission.
