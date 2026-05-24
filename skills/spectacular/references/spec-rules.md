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

Stub doc. Always scaffolded at `spectacular init`. Acts as the **index** of system capabilities — dense bullets pointing to per-capability spec files under `specs/<capability>/SPEC.md` when needed (see request `spec-refactor`).

**Verbs:**
- `grill` → polite no-op + hint
- `refine` → whole-doc rewrite pass
- `review` → structural check + sync-with-archive check at archive-time (see [[spec-sync]])

**Snapshot-on-edit: true** — canonical project-wide doc. Snapshotted aggressively because SPEC evolves with every shipped request.

**Sync flow:** archive of a verified request triggers `spec-sync` skill flow — proposes SPEC.md edits to reflect what was just shipped. See [[spec-sync]].
