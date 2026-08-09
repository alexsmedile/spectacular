---
type: mission-preparation-handoff
schema_version: spectacular.handoff.v2
handoff_id: H22
session: P0-preparation
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
provenance_parent_commit: 4900b1bc7e18b2cf1eeb62f4445cb6b884e0a4ab
provenance_parent_tree: 6e20fff0b97ec8d87ba53a222c46037a62f60f39
date: 2026-08-09
---

# H22 — P0 v1 safety-stabilization preparation

## Mission

Prepare—but do not activate—the narrow P0 repair Mission for the two proven v1 contract defects:
canonical record-kind reading and separately authorized remote deletion. Produce an owner-approved
Mission charter with exact behavior, code/test/reference scope, authority envelope, evidence,
independent-review gate, branch posture, repair budget, recovery point, and stop conditions.

This is read-only Mission preparation. It may inspect code/tests and grill Alex. It may not edit
files, create a branch/commit/PR, run destructive or provider mutations, activate P0, change v2
semantics, release/freeze v1, migrate a project, or authorize W0.

## Baseline gate

The launch prompt supplies the immutable dispatch commit/tree. Verify that pair, detached clean
state, this handoff and sidecar, and the current program contract/sidecar before analysis.
`provenance_parent_*` is pre-dispatch history only. Stop on any mismatch.

## Required inputs

Read only:

- `AGENTS.md`, `.spectacular/AGENTS.md`, this handoff and sidecar;
- `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.1` and sidecar;
- `EXECUTION-AUTHORITY-CONTRACT.md@1.0`, `EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0`, and
  `MISSION-PREPARATION-CONTRACT.md@1.0` with sidecars;
- PZL-047, PZL-048, `source-005-cli-contract-audit.md`, the H21 boundary bounce, H21-R1 return, and
  H21-R1 integration correction;
- the exact current CLI readers/cleanup paths, focused tests, help/reference/spec assertions, and
  nearby helpers needed to bound the repair.

Do not preload unrelated v2 evidence, archives, or the entire test suite.

## Locked P0 outcome

P0 must:

1. introduce one shared v1 record-kind reader that prefers canonical `kind`, falls back to legacy
   `type`, and replace the affected direct reads proven by the audit;
2. make local cleanup retain remote branches by default and require separate explicit authority
   before any remote deletion;
3. reconcile focused CLI help, tests, and operative skill/reference promises with that behavior;
4. prove both repairs on the final Mission head and receive an independent safety/boundary review.

P0 must not redesign Wayfinding, AFK, workspace coordination, v2 vocabulary, providers, or the
future Go core. It must not release v1 or perform an actual remote mutation outside disposable test
repositories.

## Interactive protocol — mandatory

For each cluster: inspect primary evidence; present 2–3 coherent options with consequences;
recommend one; ask exactly one owner-decision question; stop and wait; record Alex's verbatim and
normalized answer. Do not infer approval from “start P0” or earlier program acceptance.

### A — Remote-deletion consent contract

Compare at least:

- explicit `--delete-remote` in addition to `--apply --yes`, with remote retained by default;
- removing remote deletion from Spectacular and requiring a separate native Git/provider action;
- any materially better narrow alternative supported by current code.

Address workspace and AFK consistency, dry-run output, conflicts with existing `--keep-remote`, and
how tests prove no implicit remote deletion. Ask Alex to choose the P0 behavior. Stop.

### B — Record-kind and repair scope

Map every affected direct `type` reader and the smallest shared helper boundary. Present the exact
code/test/help/reference surfaces that must change and explicit exclusions. Apply Design
Sufficiency and Slice Quality: `sufficient | needs-evidence | needs-decision`, and `coherent |
too-broad | fragmented | dependency-bound`.

Ask Alex to approve or change the repair scope. Stop.

### C — Authority, evidence, and delivery envelope

Recommend a P0 charter. Default proposal:

- one isolated `codex/fix/v1-safety-stabilization` worktree/branch from the accepted baseline;
- local edits, focused tests, full required baseline suite, and coherent commits allowed;
- no push or draft PR unless Alex explicitly grants it here;
- no real remote deletion, release, tag, merge, deployment, migration, W0 work, or unrelated refactor;
- R1 repair budget: one hypothesis-changing repair cycle after a failed required check, then bounce;
- independent reviewer must not implement the changes and must inspect the final diff plus primary
  test evidence;
- terminal return includes baseline/final head, files, exact commands/results, reviewer findings,
  scope deviations, remaining limitations, and one next action.

Ask Alex to approve or change the envelope. Stop.

### D — Activation recommendation

Skeptically review the complete proposed charter for hidden v2 adoption, provider/destructive
authority, incomplete affected-reader inventory, test gaps, reference drift, branch conflict, or
scope creep. Present only material findings and corrections.

Ask whether Alex approves the corrected P0 charter for central activation, bounces a named part, or
escalates a decision. Stop.

## Required return

Only after all four clusters have explicit owner dispositions, return a complete
`spectacular.handoff-return.v2` packet with:

- exact baseline/tree, handoff/program hashes, clean-state result, and read set;
- owner dispositions verbatim and normalized by cluster;
- final P0 outcome, non-goals, affected paths/surfaces, branch/effect permissions, evidence and
  independent-review contract, R1 budget, recovery point, and stop conditions;
- Design Sufficiency and Slice Quality verdicts with rationale;
- skeptical findings and resolutions;
- remaining Type-2 implementation details that the executor may decide within the charter; and
- `next_action: central accept|bounce|escalate; do not activate P0 or W0`.

Central orchestration alone records the charter and dispatches an implementation Mission.
