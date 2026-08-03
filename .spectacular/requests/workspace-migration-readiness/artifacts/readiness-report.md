---
updated: 2026-08-03
verdict: go-for-spec-no-go-for-migration
against: 70745412435a846089ee5e3019427e6b2af4a3db
traffic: parallel
---

# Workspace migration readiness report

## Verdict

**GO** to write and grill a focused schema-3 contract specification.

**NO-GO** to implement or apply a schema-3 migration now.

The repository is not blocked by corruption or a detected private-data leak. It is blocked by an undefined breaking target: schema 2.0 already exists, while the roadmap's future “v2” plan names no concrete fields/layout and collides with the shipped schema number.

## What is ready

- Baseline and remote main matched exactly after fetch.
- No open PR or active request conflicts with this request.
- Full doctor reported zero errors; migration registry chain validates.
- `.spectacular.local/` is ignored, absent, untracked, and absent from reachable Git-history pathnames.
- D22 and D23 define the private boundary, schema numbering, compatibility window, branch isolation, and recovery model.
- Existing migration infrastructure can host a future `2.0 → 3.0` edge after the contract is frozen.

## What prevents migration apply

1. No approved schema-3 path/field delta exists.
2. The roadmap still describes the already-used v2 schema name as future work.
3. The local override prose is broader than the confirmed security boundary and has no general CLI implementation.
4. `status --against-latest` does not correctly explain a newer workspace.
5. A global fail-closed mutator guard for unknown newer schemas was not found.
6. Schema 2.0 currently tolerates tracked singular/root drift (`debug/`, `.last-mutation`, `.DS_Store`, and an undeclared `migrations.log`).

## Proposed next specification

Create `SPC-004 — Workspace schema 3 contract and migration safety` with this smallest scope:

- schema/product version independence and corrected roadmap terminology;
- explicit required/optional/forbidden shared-path contract;
- allowlisted private local-state contract and protected tracked-path detection;
- older/equal/newer schema behavior shared by all mutators;
- additive soak and dry-run stages before the breaking flip;
- migration manifest/rollback invariants and fixture matrix;
- explicit exclusions that keep GitHub collaboration behavior in SPC-003 and cleanup debt outside the schema migration.

The actual migration request must be created only after SPC-004 is grilled, approved, and its breaking delta is concrete.

## Open questions requiring the user

1. What real breaking change should schema 3 introduce? If there is none, should Spectacular keep schema 2.0 and implement these protections additively?
2. Should `.spectacular.local/` use one root compatibility marker or feature-owned versions?
3. Where should durable migration provenance live, and is current `migrations.log` durable or disposable?
4. What is the disposition of the legacy `debug/` trace after its outcome is classified?

## Traffic recheck requirement

Before SPC-004 execution or any migration branch begins, fetch again and recalculate traffic. A changed branch/PR/request graph invalidates this report's `parallel` assessment without erasing its evidence.
