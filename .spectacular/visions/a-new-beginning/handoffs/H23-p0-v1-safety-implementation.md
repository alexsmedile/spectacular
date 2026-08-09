---
type: implementation-mission-handoff
schema_version: spectacular.handoff.v2
handoff_id: H23
mission: P0-v1-safety-stabilization
mode: isolated-local-implementation
authority: central-orchestration
status: authorized-for-dispatch
implementation_baseline_commit: ea0aba4eeceba008066aabd1d672235284aa9cd0
implementation_baseline_tree: f50202c8c700ae4778c23f053981f8c33bfc5a80
required_branch: codex/fix/v1-safety-stabilization
push: forbidden
draft_pr: forbidden
date: 2026-08-09
---

# H23 — implement P0 v1 safety stabilization

## Authority

Implement the accepted `P0-V1-SAFETY-MISSION-CHARTER.md@1.0` on the exact named local branch and
baseline. The charter is available from the launch prompt's program commit through `git show`; it
is not part of the older implementation baseline. Revalidate baseline, branch, clean state, and
charter before mutation.

Allowed: bounded local edits, disposable Git test repositories, required checks, and coherent
local commits. Forbidden: push, PR, merge, release/tag, deployment, migration, W0 work, real remote
deletion/provider effects, public-doc edits, and unrelated refactors.

## Work

1. Inspect the exact direct reader and cleanup callers; record the final bounded inventory.
2. Implement the private `kind`-first / legacy-`type` fallback reader and route the approved
   semantic readers through it; remove the unused direct read.
3. Remove remote deletion, `--keep-remote`, and ignored `--remote` behavior from Workspace/AFK
   cleanup while preserving local recovery and honest remote-state reporting.
4. Correct CLI help, operative internal references, and focused tests. Report exact stale
   `docs/commands.md` paragraphs for Pageworks; do not edit them.
5. Run focused checks, then the required baseline suite. Use at most the chartered R1 repair cycle.
6. Inspect the final diff for scope and secret/sensitive output; create coherent local commit(s).
7. Return implementation evidence for an independent reviewer. Do not claim P0 accepted/closed.

## Required evidence

Run every command named in the charter. Record command, exit status, final-head relevance, and any
limitations. Primary tests must prove remote preservation, local cleanup/recovery, kind precedence,
legacy fallback, deterministic ranking, Doctor alignment, and help execution/accuracy.

## Stop conditions

Use the charter's stop conditions. In particular, stop on any real remote/provider effect,
baseline conflict, v2 adoption, scope expansion, unresolved required check after R1, or need to edit
Pageworks-owned public documentation.

## Return

Return `spectacular.handoff-return.v2` with exact baseline/final head, changed files, coherent
commit refs, complete command results, implementation decisions within discretion, assumptions,
limitations, Pageworks correction report, scope deviations, and one next action:
`central dispatch independent P0 review; do not accept P0 or unblock W0`.
