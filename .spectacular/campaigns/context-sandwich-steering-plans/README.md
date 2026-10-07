---
type: Reference
version: "0.1"
created: "2026-08-23T00:06:46+02:00"
updated: "2026-08-23T00:06:46+02:00"
---

# Context-sandwich Mission plan inputs

These files are `mission start` inputs, not canonical started Missions.

- `M15-governance-contract-baseline.md` is the next activation-ready plan.
- M16-M21 are preserved future sketches. Re-prepare only the next sketch after
  its predecessor closes and current Evidence is available.
- Create and switch to the Mission branch/worktree with native Git before running
  `spectacular mission start`.
- Do not start any plan directly on `main`.

The Campaign remains the planning map. Only `mission start` creates a canonical
record under `.spectacular/missions/`.
