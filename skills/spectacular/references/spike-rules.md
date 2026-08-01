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

Spikes use canonical `SPK-NNN` IDs and answer “does this approach work?” Execution requires human authorization and should use an isolated `spike/prototype-*` branch through [[afk-git-hygiene]]. Completion requires evidence plus `result: supported | refuted | inconclusive`; inconclusive stays fog. Spike code is disposable, while evidence and outcomes remain durable. `PRT` is reserved but has no standalone lifecycle yet. See [[lifecycle-contract]].

Verbs: `spectacular spike new|list|resolve`. See [[canonical-ids]].
