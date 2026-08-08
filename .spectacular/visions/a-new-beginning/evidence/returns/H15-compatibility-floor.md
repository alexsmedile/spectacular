---
type: handoff-return
schema_version: spectacular.handoff-return.v2
handoff_id: H15
session: compatibility-floor
status: accepted
central_disposition: accept
baseline_commit: 1dbeb40bbb611667de76e6b632811451211f72c1
baseline_tree: 23db697521fdb79cbf3f4c3a43ac63adac53fe09
baseline_dirty: declared-pre-existing-untracked
date: 2026-08-09
---

# H15 return — clean-break cutover and recovery

## Central disposition

**Accept.** H15 verified all ten immutable inputs and obtained explicit owner
dispositions for all four remaining clusters. Its clean-break boundary preserves
the accepted truth, authority, evidence, retrieval, and public-interface
contracts. It does not authorize S10 deletions or choose S11 mechanics.

## Accepted decisions

1. v1 freezes at a final immutable release without routine maintenance.
2. v2 supports only v2 workspaces and contains no legacy compatibility logic.
3. Migration is a whole-project atomic transaction through a recoverable source
   snapshot and separately validated v2 candidate.
4. Legacy interpretation lives in an isolated, removable migration capsule.
5. S10 retirement is itemized and evidence-gated; unique truth and recovery must
   be preserved before removal.

## Owner dispositions

> “freeze v1”

> “a” — whole-project atomic migration

> “ok” — sealed, dependency-isolated disposable migration capsule

> “a” — evidence-gated retirement with unique-truth preservation

## Boundaries retained

- Failed validation or withheld owner approval leaves v1 authoritative.
- Ambiguous semantic mappings stop for owner disposition.
- The migration capsule cannot become a dependency of v2 core.
- A Git tag alone is insufficient deletion evidence.
- S10 decides survival; S11 decides implementation and migration mechanics.

## Result

The reconciled contract is
[`../../CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md`](../../CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md).
S10 becomes next-ready after validation and commit; it is not independently
authorized by this record.
