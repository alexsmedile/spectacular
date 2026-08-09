---
type: implementation-mission-charter
mission: P0-v1-safety-stabilization
version: 1.0
status: abandoned
source_handoff: H22
accepted_by: owner
accepted_at: 2026-08-09
activated_by: central-orchestration
baseline_commit: ea0aba4eeceba008066aabd1d672235284aa9cd0
branch: codex/fix/v1-safety-stabilization
push: forbidden
draft_pr: forbidden
terminal_reason: owner-deprioritized-v1
result_commit: e6b1bfab5b2bb9e50ec8bdb94944a9ee21f0f054
result_disposition: retained-unmerged
next_action: none
---

# P0 Mission Charter — v1 safety stabilization

> **Terminal disposition:** abandoned by explicit owner direction on 2026-08-09. The implementation
> branch and review evidence remain recoverable but must not be merged. Remaining repair,
> documentation, review, and closure gates are cancelled; this charter is historical.

## Outcome

Restore two proven v1 contracts before any v2 work: canonical record-kind reading and local-only
cleanup wrappers. P0 makes no v2 product or interface commitment.

## Required behavior

1. One private Bash-v1 reader prefers canonical `kind` and falls back to legacy `type` across the
   complete affected semantic-reader inventory.
2. Workspace and AFK cleanup remain local-only. They validate and report remote state but never
   delete a remote branch.
3. Obsolete `--keep-remote` and ignored `--remote` parsing are removed.
4. Local validation, archive-ref verification, restore reporting, merge/base/provider checks, and
   local deletion remain intact.
5. Wayfinding output says `KIND`, and affected help executes without shell-substitution errors or
   destructive promises.

## Scope

Expected implementation surfaces:

- `cli/spectacular`
- `tests/cli/workspace.test.sh`
- `tests/cli/afk-git-hygiene.test.sh`
- `tests/cli/wayfinding-sequencer.test.sh`
- `tests/cli/wayfinding-contract.test.sh`
- `tests/cli/doctor.test.sh`
- `skills/spectacular/references/afk-git-hygiene.md`
- `skills/spectacular/references/wayfinding-sequencer.md`
- one focused test file if it is clearer than extending existing tests

The executor may report—but may not edit—stale Pageworks-owned paragraphs in `docs/commands.md`.
That correction is a P0 acceptance gate owned by the documentation boundary.

## Non-goals

No AFK/Workspace consolidation, command redesign, record migration, writer/schema conversion, v2
implementation/fallback/alias/provider design, release/tag/merge/deployment, project migration, W0,
real remote deletion, unrelated refactor, or public-document edit.

## Authority envelope

- Work only in `codex/fix/v1-safety-stabilization` from the exact baseline.
- Local bounded edits, disposable local/bare Git test repositories, focused/full checks, and
  coherent local commits are allowed.
- Push, PR, real-provider mutation, non-disposable remote access, release, tag, merge, deployment,
  migration, and remote deletion are forbidden.
- One hypothesis-changing repair cycle is allowed after a failed required check or blocking review
  finding. A repeated failure bounces to central orchestration.

## Evidence

Focused checks:

```bash
bash tests/cli/workspace.test.sh
bash tests/cli/afk-git-hygiene.test.sh
bash tests/cli/wayfinding-sequencer.test.sh
bash tests/cli/wayfinding-contract.test.sh
bash tests/cli/doctor.test.sh
```

Required baseline:

```bash
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
./cli/spectacular --help
bash tests/run.sh
```

Evidence must prove no cleanup path deletes a remote; matching/moved remotes remain; valid local
cleanup and recovery work; remote state is reported honestly; `kind` wins over `type`; legacy
`type` works; ranking is deterministic; Doctor uses the shared reader; and affected help executes
and describes reality.

An independent reviewer who did not implement P0 must inspect the final diff/head, test evidence,
reference alignment, reader inventory, provider boundary, and v1/v2 separation. Pageworks-owned
public documentation must be reconciled before P0 can be accepted and closed.

## Stop conditions

Stop on baseline/branch conflict, wider v1 scope, hidden v2 adoption, real provider/destructive
effect, an additional reader that changes the approved boundary, unresolved Pageworks ownership
conflict, exhausted evidence budget, or ambiguous consent/scope/authority/reviewer independence.

## Preparation verdicts

- Design Sufficiency: `sufficient`
- Slice Quality: `coherent`

Executor discretion is limited to the private helper's name/location, remote-report wording,
focused-test organization, and coherent commit grouping; none may alter observable semantics or
authority.
