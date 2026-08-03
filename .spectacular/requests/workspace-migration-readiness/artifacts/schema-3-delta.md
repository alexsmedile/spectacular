---
updated: 2026-08-03
status: proposed
source_schema: "2.0"
target_schema: "3.0"
---

# Proposed schema-3 contract delta

This is a readiness proposal, not an approved specification and not migration authority.

## Confirmed invariants

- Product versions and workspace-schema versions are independent.
- `.spectacular/` is committed shared knowledge; `.spectacular.local/` is ignored private operational state.
- Local state supplements machine/user operation only. It cannot override intent, specifications, decisions, lifecycle, dependencies, authority, policy, security gates, or verification.
- Local feature paths are created lazily. Init guarantees ignore protection only.
- Unknown newer schemas are read-only until the CLI is upgraded.
- Migration is dry-run-first, branch-isolated, manifest-driven, and separate from private-data conversion.

## Candidate schema-3 requirements

### Shared workspace

1. Freeze an explicit required/optional/forbidden path contract for schema 3 instead of deriving validity from a version number alone.
2. Define the root-anchor allowlist and the status of migration provenance such as `migrations.log`.
3. Define canonical plural collection paths and a judgment path for legacy singular content.
4. Add machine-readable request relationship fields only after SPC-003 fixes their exact schema and reciprocity rules; do not smuggle them into migration code.
5. Require every mutating command to reject a workspace whose declared schema is newer than the CLI understands.

### Private local workspace

1. Define an allowlisted local-purpose registry: private ideas, GitHub machine/account settings, protected security state, and feature-owned caches/preferences.
2. Define file permissions and protected-output behavior for security-bearing paths.
3. Detect tracked/history-visible local paths without reading bodies and fail closed.
4. Do not copy, archive, or back up private material into `.spectacular/`.
5. Decide whether local state needs its own compatibility marker or whether each feature owns versioning independently.

### Compatibility window

1. Schema 2.0 remains accepted while optional schema-3 fields/paths soak.
2. Doctor validates new fields when present but does not require them during soak.
3. Dry-run and fixtures ship before `CURRENT_SCHEMA` changes.
4. `status --against-latest` and all mutators distinguish older, equal, and newer schemas.
5. Only the breaking release makes schema 3.0 required; older CLIs stop with an upgrade message.

## Not automatically part of schema 3

- GitHub Issues, Discussions, Actions, labels, or PR behavior: owned by SPC-003.
- `.DS_Store` and tracked undo breadcrumb cleanup: repository hygiene.
- Resolution of the existing legacy debug trace: judgment cleanup after classification.
- Snapshot pruning and legacy memory-status warnings: ordinary maintenance.
- Public documentation rewrite: downstream of the approved contract.

## Open design questions

1. Which concrete shared-layout or required-field change is sufficiently breaking to earn schema 3.0? If none exists, keep schema 2.0 and ship the boundary protections additively.
2. Should `.spectacular.local/` have one root compatibility marker, or should each feature-owned local file/version evolve independently?
3. Is `migrations.log` durable shared evidence that should receive a canonical home, or disposable/local execution output?
4. Should the two-file legacy `debug/` trace be moved to `debugs/`, archived, or deleted after its outcome is classified?
