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

Statuses are `open | deferred | resolved`. Fog/frontier readiness is computed by the Wayfinding sequencer and is never stored as status. Product and business trade-offs cannot be resolved autonomously.

Mechanical verbs:

- `spectacular question new <slug>`
- `spectacular question list [--status ...]`
- `spectacular question resolve <id|alias> --answer <text>`

At session start, surface unresolved user-input questions concisely: highest priority first, then dependency impact. Deferred questions remain durable but are not presented as immediate blockers.

See [[canonical-ids]] for identity and alias rules.
