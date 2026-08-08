---
type: synthesis-checkpoint
checkpoint: 015
status: current
date: 2026-08-08
authority: central-orchestration
accepted_contracts: [S01, S03A]
next_session: S02
---

# Synthesis 015 — S03A truth/provenance floor accepted

## Central disposition

H06/S03A is accepted. Five explicit owner dispositions now establish the
minimum truth and provenance conditions needed before evidence can be scored.

## Accepted floor

```text
question-specific authority
  + pointer-first provenance
  + claim-scoped freshness
  + non-authoritative projections
  + explicit unknowns / stop-and-expose conflicts
  → safe basis for S02 evidence design
```

The authoritative result is
[`../TRUTH-PROVENANCE-FLOOR.md`](../TRUTH-PROVENANCE-FLOOR.md). Its hash sidecar
binds the exact accepted version.

## Important boundary

This acceptance does not choose storage, a graph model, file taxonomy, public
vocabulary, or a command surface. The `live-decision-graph` idea is parked for
later evaluation as a generated projection only.

## Repaired sequence

```text
S01 Product Constitution accepted
→ S03A truth/provenance floor accepted
→ S02 success/evidence constitution ready
→ S03B full truth/contract model
```

S02 is now authorized as the next owner-decision session. No later session or
implementation work is authorized by this checkpoint.
