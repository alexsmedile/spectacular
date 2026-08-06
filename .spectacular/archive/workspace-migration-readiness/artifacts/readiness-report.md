---
updated: 2026-08-06
verdict: go-for-spec-no-go-for-migration
against: 9ba090cf260ed2ac7897d9a6063ae0afc6cb2bf8
traffic: unknown
---

# Workspace migration readiness report

## Verdict

**GO** to write and grill a focused additive workspace-boundary and migration-safety specification under schema 2.0.

**NO-GO** to implement or apply a schema-3 migration now. D24 reserves schema 3.0 for a future real breaking contract rather than manufacturing one for additive improvements.

The repository is not blocked by corruption or a detected private-data leak. It is blocked by an undefined breaking target: schema 2.0 already exists, while the roadmap's future “v2” plan names no concrete fields/layout and collides with the shipped schema number.

## What is ready

- Current `main` and `origin/main` matched exactly after fetch at `9ba090c`.
- The original 2026-08-03 `parallel` traffic assessment is expired; current traffic is `unknown` until a follow-up request declares its boundaries and runs preflight.
- The workspace and migration registry remain schema-2.0-current with the same three-entry registry chain.
- `.spectacular.local/` is ignored, absent, untracked, and absent from reachable Git-history pathnames.
- D22 defines the private boundary; D24 supersedes D23's assumed immediate 2.0-to-3.0 transition while retaining its compatibility, branch-isolation, and recovery protocol for any future real break.
- Existing migration infrastructure can host a future `2.0 → 3.0` edge after the contract is frozen.

## What prevents migration apply

1. No approved schema-3 path/field delta exists.
2. The roadmap still describes the already-used v2 schema name as future work.
3. The local override prose is broader than the confirmed security boundary and has no general CLI implementation.
4. `status --against-latest` does not correctly explain a newer workspace.
5. A global fail-closed mutator guard for unknown newer schemas was not found.
6. The remaining schema-2.0 concerns are contract gaps, not live root-artifact drift: the former `.last-mutation`, `.DS_Store`, and `migrations.log` claims are no longer current on `main`.

## Proposed next specification

If the safeguards are prioritized later, create a newly allocated SPC — never reuse an existing canonical ID — with this smallest scope:

- schema/product version independence, D24's schema-2 retention rule, and corrected roadmap terminology;
- explicit required/optional/forbidden shared-path contract;
- D27's intentionally tiny, non-authoritative migration receipt contract;
- allowlisted private local-state contract and protected tracked-path detection;
- older/equal/newer schema behavior shared by all mutators;
- additive schema-2 hardening now and reusable dry-run stages for a later real breaking flip;
- migration manifest/rollback invariants and fixture matrix;
- explicit exclusions that keep GitHub collaboration behavior in SPC-003 and cleanup debt outside the schema migration.

No schema migration request should be created unless a later approved specification names a real breaking delta. The future additive hardening SPC may seed a request only after grilling and approval.

## Interview resolutions

1. **Resolved — D24:** no current change earns schema 3.0; keep schema 2.0 and implement compatible protections additively.
2. **Resolved — D25:** use no local schema/version by default; add a narrow format marker only after a real incompatible change earns it.
3. **Resolved — D27:** keep the existing lightweight migration receipt; do not build a replacement ledger or command.
4. **Resolved:** remove the generated legacy trace; retain its synthetic fixture and Git history rather than creating an archive duplicate.

## Traffic recheck requirement

Before any additive-hardening execution or migration branch begins, fetch again and recalculate traffic. A changed branch/PR/request graph invalidates the original `parallel` assessment without erasing its historical evidence.
