---
status: review
updated: 2026-08-03
related:
  - PLAN.md
---

# Tasks — workspace-migration-readiness

## v1

### M1 — Reproducible baseline
- [x] Record local/remote commit identities, worktree state, branches, open PRs, active/planned requests, and traffic assessment
- [x] Run the full workspace doctor, migration registry listing, schema currency check, and focused Git ignore/tracking checks
- [x] Capture current warnings and distinguish migration blockers from unrelated maintenance
- [x] → check: every baseline claim cites a command/output identity and can be repeated from the recorded commit

### M2 — Shared/private boundary inventory
- [x] Enumerate live `.spectacular/` paths by purpose, authority, lifecycle, freshness duty, and retention class without loading archive bodies
- [x] Enumerate declared and present `.spectacular.local/` paths by purpose, creator, sensitivity, permissions, and tracked/ignored status without printing protected content
- [x] Map which local settings may supplement operation and prove that none may override shared truth, authority, policy, lifecycle, or verification
- [x] Define the fail-closed response for tracked local paths, suspected leakage, and history exposure
- [x] Write `artifacts/workspace-inventory.md` and `RISKS.md`
- [x] → check: the inventory has no unclassified path and the leakage procedure exposes filenames only

### M3 — Schema-3 migration contract proposal
- [x] Reconcile the roadmap’s future-v2 wording with the already-live workspace schema 2.0 and product/schema version independence
- [x] Draft `artifacts/schema-3-delta.md` covering additions, changes, removals, optional soak behavior, old-client refusal, and lazy local creation
- [x] Draft `artifacts/migration-manifest.md` classifying every step as mechanical, judgmental, local-only, shared, reversible, or separately security-gated
- [x] Define fixture coverage for schema 2.0, additive soak, schema 3.0, missing ignore protection, tracked local paths, offline GitHub, divergence, and rollback
- [x] → check: the proposal changes no production implementation surface and every migration step names an authority and recovery boundary

### M4 — Readiness decision
- [x] Review artifacts against D22, D23, SPC-003 boundaries, the PRD, and migration registry contract
- [x] Resolve or surface every ambiguity with an owner; do not convert unanswered product/security choices into autonomous decisions
- [x] Write `artifacts/readiness-report.md` with one go/no-go result and the smallest proposed schema-3 SPC scope
- [x] Re-run traffic and baseline checks before recommending the implementation path
- [x] → check: the report is decision-ready, evidence-linked, placeholder-free, and leaves the request in review rather than executing migration

## v2 (deferred)

- [~] Write or approve the schema-3 capability specification
- [~] Implement or apply any workspace migration
- [~] Change roadmap state, workspace schema, production code, or public GitHub configuration
