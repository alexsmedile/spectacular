---
id: SPC-001
type: specification
status: implemented
approved_at: 2026-08-01
approved_by: user
target_version: v1.36.0-execution
updated: 2026-08-03
summary: "Canonical discovery records, fog/frontier sequencing, and approved-spec execution handoff"
related: []
version: 1.5
implemented_at: 2026-08-03
verified_against: commit bf0756f
---

# SPC-001 — Canonical discovery records, fog/frontier sequencing, and approved-spec execution handoff

## Intent

Turn a user's broad dream and unresolved alternatives into durable, low-entropy discovery state that survives sessions and converges into an approved executable specification.

## Requirements

- Maintain distinct `decisions/`, `questions/`, and `ideas/` Markdown databases.
- Represent evidence gathering as `RES-NNN` research records and feasibility work as human-authorized `SPK-NNN` spikes.
- Use canonical cross-store IDs while accepting compact conversational aliases.
- Compute fog/frontier from canonical dependencies; do not persist derived readiness.
- Surface unresolved `requires_user_input` questions at session start unless deferred.
- Reject dangling/cyclic graphs and present discovery nodes in deterministic, dependency-first topological order.
- Rank frontier nodes by explicit priority, then prefer high-uncertainty work before deterministic specification/execution work.
- Keep collaborative specs `draft` and AFK-authored specs `unconfirmed` until explicit approval; only `approved` specs may seed implementation requests.
- Use `vX.Y.Z-discovery` and `vX.Y.Z-execution` target conventions while dependency order remains primary.
- Park unexpected execution scope in ideas or a later target instead of expanding the current milestone.
- Map “park this idea,” “put it on ice,” “find your way to,” and “act on goal” through the same gated CLI verbs used by explicit commands.
- Detect strong dependency signals across PRD, roadmap, plans, and specs; warn and propose explicit edges without mutating or reslotting source documents.
- Keep AFK Git behavior opt-in and dry-run-first; isolate draft specs, spikes, forks, and approved execution using host-prefix-compatible branch classes.
- Require clean-tree/primary-branch preflight, durable branch provenance, and archive-before-local-delete; never delete remote branches autonomously.
- Gate PR handoff on a verified request, passing verification evidence, approved source spec, fresh tests, and separate breaking-change approval; stop before merge.

## Evidence and decisions

- User decisions R1–R4 supplied 2026-08-01 and captured in `requests/wayfinding-contract/PLAN.md § Decisions`.
- `references/canonical-ids.md`, `question-rules.md`, `research-rules.md`, `spike-rules.md`, and `spec-lifecycle.md` define the working contract.
- `references/wayfinding-sequencer.md` defines strict ordering, metaphor routing, and execution-scope boundaries.
- `references/afk-git-hygiene.md` defines the AFK authorization and branch lifecycle contract.

## Confirmation

**Implemented 2026-08-03** — the user's R1–R4 decisions establish the product contract. The archived `wayfinding-contract` request records 53 focused Wayfinding assertions and the full 18-file regression suite passing; implementation is verified against commit `bf0756f`.
