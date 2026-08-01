---
description: Wayfinding lifecycle for feature specifications from idea-loop output to current execution truth.
when_to_use: Creating, confirming, updating, deprecating, or acting on an SPC-NNN specification.
---

# Specification Lifecycle

The idea loop explores possibilities. A specification is the convergent artifact that follows.

```text
idea/discovery → unconfirmed spec → human confirmation → current spec → request → implementation → verification
```

- `unconfirmed` — drafted from vision, questions, research, spikes, and decisions; not executable.
- `current` — explicitly confirmed and updated; eligible to seed an implementation request.
- `deprecated` — retained for history but no longer executable.

Canonical specs use `SPC-NNN-<slug>.md`. `spectacular spec confirm` snapshots before changing state. `spectacular spec act` refuses anything not current, scaffolds a normal request, and stores `source_spec: SPC-NNN` provenance.

Release planning may distinguish `vX.Y.Z-discovery` from `vX.Y.Z-execution`. Dependencies override those targets: execution cannot precede a required discovery node.
