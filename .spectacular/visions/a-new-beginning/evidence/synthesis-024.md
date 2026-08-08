---
type: synthesis-checkpoint
checkpoint: 024
status: current
date: 2026-08-09
authority: central-orchestration
accepted_contracts: [S01, S03A, S02, S03B, S04, S05, S06, S07, S08, S09, compatibility-floor]
next_session: S10
---

# Synthesis 024 — clean-break cutover accepted

## Central disposition

H15/compatibility-floor is accepted. Spectacular v2 supports v2 workspaces only;
v1 freezes as an immutable recoverable release, and no legacy compatibility
logic enters v2 core.

## Accepted cutover spine

```text
frozen v1 + verified snapshot
  → isolated migration capsule
  → separate v2 candidate
  → explicit ambiguity/unsupported-state report
  → clean-v2 validation
  → owner acceptance
  → cutover with rollback evidence
```

The authoritative result is
[`../CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md`](../CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md).

## Next gate

S10 is next-ready. It may classify each contested subsystem as keep, simplify,
extract, merge, or retire, but it may remove nothing without replacement,
unique-truth preservation, and verified recovery evidence.
