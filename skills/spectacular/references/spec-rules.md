---
doc-id: spec
mode: stub
location: .spectacular/SPEC.md
scope: project-wide
template: templates/spec/base.md
snapshot-on-edit: true
summary: "System spec — index of what the system actually is and how it behaves right now"
status: active
---

# SPEC Rules

Stub doc — `grill` / `refine` follow the [[doc-index]] § Stub default behavior; snapshot-on-edit is true (snapshotted aggressively — SPEC evolves with every shipped request). Always scaffolded at `spectacular init`.

**Doc-specific role — the index.** SPEC.md is the **index** of system capabilities: dense bullets pointing to per-capability spec files under `specs/<capability>/SPEC.md` when needed (see request `spec-refactor`).

**`review` override.** Beyond the default structural check, `spec review` also runs a sync-with-archive check at archive-time (see [[spec-sync]]).

**Sync flow:** archive of a verified request triggers the `spec-sync` skill flow — proposes SPEC.md edits to reflect what was just shipped. See [[spec-sync]].
