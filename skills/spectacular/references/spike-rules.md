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

Spikes use canonical `SPK-NNN` IDs and answer “does this technical approach work?” Execution requires human authorization and should use an isolated `spike/prototype-*` branch through [[afk-git-hygiene]]. Completion requires evidence plus `result: supported | refuted | inconclusive`; inconclusive stays fog. Spike code is disposable, while evidence and outcomes remain durable. A Vision may link the SPK and present a concise demo, transcript, screenshot, or comparison for human reaction, but the feasibility result never approves a product direction. A UX/workflow prototype belongs to its Vision (or feedback/request when a target already exists); `PRT` remains reserved. A production tracer bullet is approved `SPC` execution, never a spike or prototype. See [[discovery-protocol]], [[vision-rules]], and [[lifecycle-contract]].

Verbs: `spectacular spike new|list|resolve`. See [[canonical-ids]].
