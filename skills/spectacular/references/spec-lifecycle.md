---
description: Wayfinding lifecycle for feature specifications from collaborative or AFK drafting through historically verified implementation.
when_to_use: Creating, approving, revising, implementing, deprecating, archiving, or acting on an SPC-NNN specification.
---

# Specification Lifecycle

The idea loop explores possibilities. A specification is the convergent artifact that follows.

The authoritative state machine and gates live in [[lifecycle-contract]]. In short:

```text
draft|unconfirmed → approved → implemented → archived
       └──────────────────────────────→ archived (rejected/abandoned)
implemented → superseded|deprecated → archived
```

`draft` means collaboratively unfinished; `unconfirmed` means AFK-authored. Only `approved` authorizes implementation. `implemented` is a historical claim that requires a verified request, closed docs impact, `implemented_at`, and `verified_against`; it never claims continuous agreement with code. After verified integration/merge, the detailed SPC should leave active context through explicit dry-run-first `spec archive`; code/tests and the capability index become the live truth.

Canonical specs use `SPC-NNN-<slug>.md`. `spectacular spec approve` snapshots before approval (`confirm` is a compatibility alias). `spectacular request new --from SPC-NNN` mechanically creates a planned PLAN/TASKS bundle; `/spectacular act SPC-NNN` (or the unambiguous `/spectacular SPC-NNN`) owns authorization, activation provenance, compiled context, native planning, and implementation. Terminal `spec act` redirects instead of partially executing that agentic flow. See [[request-workflow]]. Behavior-changing revisions use a new SPC with `supersedes`; an active implemented predecessor transitions atomically to superseded, while an already archived implemented predecessor remains immutable and is linked from the replacement.

For a natural-language request to draft a specification, first run
[[intent-routing]] and show the user an intent receipt. A new SPC is appropriate
only when the user confirms a new durable implementation boundary; nearby
repository state, a related request, or `docs_impact` can support that decision
but cannot make it. An explicit `spectacular spec new` command remains the
mechanical escape hatch for a user who has already made that routing decision.

Draft, unconfirmed, or approved specs that are rejected/abandoned are archived with a reason rather than purged. This removes them from active context without losing the failed path. See [[artifact-retention]].

Release planning may distinguish `vX.Y.Z-discovery` from `vX.Y.Z-execution`. Dependencies override those targets: execution cannot precede a required discovery node.
