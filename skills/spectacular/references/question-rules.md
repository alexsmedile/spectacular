---
doc-id: question
mode: index
location: .spectacular/questions/
entries-dir: .spectacular/questions/
scope: project-wide
template: templates/question/entry.md
snapshot-on-edit: false
summary: "Active ambiguities and blockers surfaced for human resolution."
status: active
---

# Question Rules

Questions are explicit open loops, not ideas and not ADRs. Each entry has a canonical `QUE-NNN` ID, states the ambiguity, carries options and dependency links, and defaults to `requires_user_input: true`.

Active statuses are `open | deferred`. Resolution records answer, actor, provenance, and optional resulting decision, then moves the file to `archive/questions/` with `status: archived` and `archived_from: resolved`. Fog/frontier readiness is computed by the Wayfinding sequencer and is never stored as status. Product and business trade-offs cannot be resolved autonomously.

Mechanical verbs:

- `spectacular question new <slug>`
- `spectacular question list [--status ...]`
- `spectacular question resolve <id|alias> --answer <text>`

At every human-agent session start—including `/spectacular` status and onboarding—surface unresolved user-input questions before other work: highest priority first, then dependency impact. Deferred questions remain durable but are not presented as immediate blockers.

An answer creates or links a DEC only when it represents a real choice between alternatives. Factual clarification, evidence, rejection, or “not applicable” resolution remains answer provenance without manufacturing a decision. Archived questions remain valid satisfied canonical dependencies but are excluded from normal fog/frontier iteration.

See [[canonical-ids]] for identity and alias rules and [[artifact-retention]] for the active/history boundary.
