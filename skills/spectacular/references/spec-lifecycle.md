---
description: Wayfinding lifecycle for feature specifications from collaborative or AFK drafting through historically verified implementation.
when_to_use: Creating, merging, revising, implementing, deprecating, archiving, or acting on a specification contract.
---

# Specification Lifecycle

Ideas capture possibilities. Optional Vision workspaces shape materially unsettled direction through grounded alternatives and explicit human approval. A specification is the convergent implementation contract that follows either an already-settled destination or an approved Vision.

The authoritative state machine and gates live in [[lifecycle-contract]]. In short:

```text
draft|unconfirmed → approved → implemented → archived
       └──────────────────────────────→ archived (rejected/abandoned)
implemented → superseded|deprecated → archived
```

`draft` means collaboratively unfinished; `unconfirmed` means AFK-authored. A spec becomes an approved shared contract when its spec-branch commit is merged into the configured shared base (normally `main`). `status: approved` is retained only as a compatibility signal, never the authority for new work. `implemented` is a historical claim that requires a verified request, closed docs impact, `implemented_at`, and `verified_against`; it never claims continuous agreement with code. After verified integration/merge, the detailed spec should leave active context through explicit dry-run-first `spec archive`; code/tests and the capability index become the live truth.

New specs use UUIDv7 `id` values and `specs/<slug>.md` paths. The ID is immutable and global; the slug is the human-facing locator and may change without changing links. Author a spec on a `spec/` branch and merge its PR into the configured shared base to approve it. `spectacular request new --from <slug|UUIDv7>` mechanically creates a planned PLAN/TASKS bundle only when that merged commit is in the execution branch ancestry. `approve`/`confirm` are compatibility commands for older records. See [[request-workflow]]. Behavior-changing revisions use a new spec with `supersedes`; an active implemented predecessor transitions atomically to superseded, while an already archived implemented predecessor remains immutable and is linked from the replacement.

`/spectacular vision derive <slug>` requires `status: approved` and drafts an SPC from the Vision's Chosen direction, Boundaries, approved fragments, and linked evidence. It records the Vision as provenance and updates `derived_to` after the draft exists. Derivation never approves the SPC, creates a request, or copies rejected/pending fragments into requirements.

For a natural-language request to draft a specification, first run
[[intent-routing]] and show the user an intent receipt. A new SPC is appropriate
only when the user confirms a new durable implementation boundary; nearby
repository state, a related request, or `docs_impact` can support that decision
but cannot make it. An explicit `spectacular spec new` command remains the
mechanical escape hatch for a user who has already made that routing decision.

Draft, unconfirmed, or approved specs that are rejected/abandoned are archived with a reason rather than purged. This removes them from active context without losing the failed path. See [[artifact-retention]].

Release planning may distinguish `vX.Y.Z-discovery` from `vX.Y.Z-execution`. Dependencies override those targets: execution cannot precede a required discovery node.
