---
updated: 2026-08-03
status: reserved
source_schema: "2.0"
target_schema: "3.0"
---

# Proposed schema-3 contract delta

This is a reserved future-breaking-contract proposal, not an approved specification or current migration authority. D24 keeps the live workspace on schema 2.0 because the presently identified protections are additive.

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
2. Define the root-anchor allowlist and D27's lightweight `migrations.log` operational exception.
3. Define canonical plural collection paths and a judgment path for legacy singular content.
4. Add machine-readable request relationship fields only after SPC-003 fixes their exact schema and reciprocity rules; do not smuggle them into migration code.
5. Require every mutating command to reject a workspace whose declared schema is newer than the CLI understands.

### Private local workspace

1. Define an allowlisted local-purpose registry: private ideas, GitHub machine/account settings, protected security state, and feature-owned caches/preferences.
2. Define file permissions and protected-output behavior for security-bearing paths.
3. Detect tracked/history-visible local paths without reading bodies and fail closed.
4. Do not copy, archive, or back up private material into `.spectacular/`.
5. Per D25, add no local schema/version by default. A feature may introduce a narrow format marker only after a real incompatible change cannot be safely detected or migrated.

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

## Interview resolutions

1. **Resolved — D24:** no current change is sufficiently breaking; keep schema 2.0 and ship boundary protections additively. Reopen schema 3.0 only when a real breaking delta is proposed.
2. **Resolved — D25:** neither by default; versioning is introduced only when a real incompatible local format earns it.
3. **Resolved — D27:** keep `migrations.log` as an optional one-line-per-success user receipt, not authority or normal context; add no subsystem around it.
4. **Resolved:** remove the generated legacy trace; the synthetic test fixture and Git history preserve everything useful without polluting live or archived Spectacular context.
