---
id: SPC-006
type: specification
status: archived
target_version: "tbd"
supersedes: ""
updated: 2026-08-06
summary: "Define additive workspace-boundary and schema-mutation safety hardening under schema 2.0"
related: []
version: 1.1
approved_at: 2026-08-06
approved_by: alex
implemented_at: 2026-08-06
verified_against: uncommitted
archived_from: implemented
archived_at: 2026-08-06
archive_reason: Implemented work is archived with its verified request.
---

# SPC-006 — Define additive workspace-boundary and schema-mutation safety hardening under schema 2.0

## Intent

Make every mutating CLI path safe when a workspace declares a schema newer than
the installed CLI understands, and make the shared/private workspace boundary
match DEC-022 in both guidance and executable checks. The result preserves
`workspace_schema: "2.0"` and adds no migration edge, layout conversion, or
private-data migration.

## Requirements


### Schema relation contract

- Define one shared schema-relation helper for `older`, `equal`, and `newer`
  workspaces relative to `CURRENT_SCHEMA`.
- Every command that can write `.spectacular/` must refuse before its first
  write when the declared workspace schema is newer than the CLI supports. The
  refusal must identify the declared and supported versions and direct the user
  to update the CLI; read-only inspection remains available.
- `spectacular status --against-latest` must distinguish older from newer:
  only an older workspace suggests `spectacular migrate`; a newer workspace
  directs the user to update the CLI and must not suggest a downgrade or
  migration.
- Preserve existing older-schema migration behavior and equal-schema behavior.
  An absent marker keeps its documented legacy interpretation.

### Shared/private boundary contract

- Replace guidance that says `.spectacular.local/` takes precedence with an
  allowlisted, supplement-only contract consistent with DEC-022. Local state
  may not override shared identity, approved specifications, decisions,
  lifecycle, dependencies, authority mappings, policy, security gates, or
  verification.
- Add filename-only detection for tracked `.spectacular.local/` paths before a
  migration can write. Detection must not read or print local-file bodies.
- A detected tracked local path fails closed: stop the operation, report only
  pathnames, and route exposure classification, credential rotation, history
  repair, deletion, or disclosure to explicit human/security authority.
- Retain lazy local creation and the existing root `.gitignore` protection. Do
  not introduce a root local schema/version, backup, archive, or automatic
  private-to-shared conversion.

### Test and documentation evidence

- Add synthetic fixtures covering older, equal, and newer schema relations for
  status and representative mutators; prove newer-schema refusal happens
  before mutation.
- Add synthetic tracked-local-path fixtures that prove filename-only output and
  refusal without relying on real private material.
- Align CLI help/reference guidance with the executable boundary and DEC-022.

### Boundaries and exclusions

- Keep `CURRENT_SCHEMA` at `2.0`; do not add a `2.0 -> 3.0` registry entry,
  flip a schema marker, or require a new workspace layout or field.
- Do not create, activate, or implement a migration request from this SPC.
- Do not inspect, copy, archive, convert, back up, delete, or expose contents
  under `.spectacular.local/`.
- Do not change GitHub configuration or broaden SPC-003's protected-security
  orchestration, and do not mix root-artifact cleanup into this work.

## Evidence and decisions

- `DEC-022` — local state is subordinate, private, lazily created, and subject
  to filename-only fail-closed handling when tracked.
- `DEC-024` — retain schema 2.0 until a real breaking contract exists;
  compatible hardening proceeds separately.
- `DEC-025` — no private-local schema/version without an actual incompatible
  local format.
- `DEC-027` — `migrations.log` is optional diagnostic provenance, not schema
  authority.
- `archive/workspace-migration-readiness/artifacts/readiness-report.md` —
  discovery verdict: go for additive hardening, no-go for schema-3 migration.
- `archive/workspace-migration-readiness/RISKS.md` — identifies the absent
  newer-schema mutator guard, incorrect status guidance, and local-precedence
  wording as the relevant outstanding risks.

## Confirmation

draft — not eligible for implementation until explicitly approved.

**Approved 2026-08-06 by alex** — User approved the additive-only schema-2 safety contract on 2026-08-06, including read-only newer-workspace diagnostics.
