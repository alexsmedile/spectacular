---
description: Wayfinding lifecycle for feature specifications from collaborative or AFK drafting through historically verified implementation.
when_to_use: Creating, approving, revising, implementing, deprecating, archiving, or acting on an SPC-NNN specification.
---

# Specification Lifecycle

The idea loop explores possibilities. A specification is the convergent artifact that follows.

The authoritative state machine and gates live in [[lifecycle-contract]]. In short:

```text
draft|unconfirmed → approved → implemented → superseded|deprecated → archived
```

`draft` means collaboratively unfinished; `unconfirmed` means AFK-authored. Only `approved` authorizes implementation. `implemented` is a historical claim that requires a verified request, closed docs impact, `implemented_at`, and `verified_against`; it never claims continuous agreement with code.

Canonical specs use `SPC-NNN-<slug>.md`. `spectacular spec approve` snapshots before approval (`confirm` is a compatibility alias). `spec act` stores `source_spec` on the request. Behavior-changing revisions use a new SPC with `supersedes`; the old SPC stays implemented until the replacement is implemented, when both transition atomically.

Release planning may distinguish `vX.Y.Z-discovery` from `vX.Y.Z-execution`. Dependencies override those targets: execution cannot precede a required discovery node.
