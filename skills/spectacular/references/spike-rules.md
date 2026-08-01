---
doc-id: spike
mode: index
location: .spectacular/spikes/
entries-dir: .spectacular/spikes/
scope: project-wide
snapshot-on-edit: false
summary: "Human-authorized feasibility experiments whose durable output is evidence."
status: active
---

# Spike Rules

Spikes use canonical `SPK-NNN` IDs and answer “does this approach work?” Execution requires human authorization and should use an isolated `spike/prototype-*` branch through [[afk-git-hygiene]]. Spike code is disposable; its sources, evidence, outcome, and optional `PRT-NNN` artifact references are durable.

Verbs: `spectacular spike new|list|resolve`. See [[canonical-ids]].
