---
type: idea
status: parked
priority: low
owner: alex
updated: 2026-07-11
origin: builder-agent request (b21, archived 2026-07-11 on M1+M2 — M3 gated, never earned)
promoted_to: null
related:
  - ../../skills/spectacular/references/build-workflow.md
  - ../archive/builder-agent/PLAN.md
---

# Idea — builder trace + fan-out (deferred from builder-agent M3)

The `builder-agent` request (b21) shipped M1 (`spec-builder` agent) + M2 (`build-workflow.md`
arc) and archived on those. Two things were deliberately **not** built, gated by Principle 11
("don't build speculatively — earn each step"). They live here until a real need proves them.

## 1. M3 — durable trace + CLI signal (gated, never earned)

The debug fleet has a trace substrate (`debugs/<job>/` JSON spine). M3 asked whether the *build*
fleet needs an analog:
- **A builder trace folder** (analog of `debugs/<slug>/`) — one artifact per dispatched milestone.
- **A `--delegable` / `--brief` CLI emit** — surface a milestone's closed brief for dispatch.

**Why parked:** M1/M2 never produced the fan-out volume that would justify either. A trace with
no fan-out to trace is ceremony; a brief-emit verb with one consumer is premature surface. Build
only if/when fan-out becomes routine and hand-assembling briefs becomes the bottleneck.

**Trigger to promote:** the *first time* a real multi-milestone fan-out is painful to orchestrate
by hand — that pain is the evidence M3 was waiting for.

## 2. M2's fan-out walkthrough — verify-on-first-real-use

`build-workflow.md`'s worth-it/fan-out gate is shipped and parity-checked against `bug-workflow.md`,
but the walkthrough on a **real multi-milestone fan-out** never ran (no such fan-out occurred before
archive). The mechanism is present; its end-to-end exercise is owed on first real use. **Not a
defect — a verification awaiting a natural trigger.** When the first genuine fan-out happens,
observe whether the gate produces the right dispatch/inline call and whether the context-assembly
chain (task row → milestone block → PLAN §2/§3/§6/§7 → brief) holds.

## 3. v2 fleet candidates (from builder-agent TASKS v2)

- **Test Agent** — independently verify builder output. *(Note: `test-verifier` shipped since,
  in fleet-arc-wiring (b25) — this candidate is largely satisfied. Re-check before reviving.)*
- **Review fleet** — correctness / inefficiency / dead-code. *(Note: `code-reviewer` shipped since,
  also b25. Likewise largely satisfied.)*

Both v2 candidates were superseded by the b25 fleet-arc-wiring work — kept here only for the
audit trail; likely nothing to build.
