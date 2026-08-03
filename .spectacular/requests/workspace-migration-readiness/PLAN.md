---
status: active
priority: medium
owner: alex
updated: 2026-08-03
build: b41
docs_impact: none
docs_impact_reason: "Internal readiness artifacts only; public documentation waits for an approved schema-3 contract"
summary: "Audit the shared and private workspace boundaries, reconcile schema numbering, and produce a safe schema-3 migration proposal without changing production behavior"
related:
  - ../../PRD.md
  - ../../specs/SPC-003-github-native-lifecycle.md
  - ../../decisions/DEC-022-keep-spectacular-local-subordinate-private-and-lazily.md
  - ../../decisions/DEC-023-migrate-workspace-schema-2-0-to-3-0.md
---

# Plan — workspace-migration-readiness

## Goal

Serve the PRD goals of reducing context rot and preserving human-agent coherence by producing an evidence-backed readiness package for evolving the shared/private workspace contract from schema 2.0 toward schema 3.0 without exposing private state or changing production behavior during discovery.

## Constraints

- Discovery only: beyond the CLI's required `last_build` allocation for this request, do not modify production code, migration functions, workspace schema, roadmap state, existing requests, or GitHub configuration.
- Treat `.spectacular/` as committed shared project knowledge and `.spectacular.local/` as ignored, private, subordinate, and lazily created per `DEC-022`.
- Preserve the existing `workspace_schema: 2.0`; the next breaking workspace layout is schema 3.0 and product versions remain independent per `DEC-023`.
- Never inspect or print protected local contents merely to inventory them. If any `.spectacular.local/` path is tracked, report filenames only and stop for explicit security handling.
- Do not load archive bodies during the normal audit. Inventory paths and metadata only unless a later named question justifies historical retrieval.
- Baseline is `main@7074541` with `origin/main@7074541`, a clean tree, no open pull requests, and traffic assessed `parallel` on 2026-08-03.
- Any remote divergence, unexpected overlap, undeclared access, or evidence that invalidates this baseline changes traffic to `unknown` and stops the audit for user review.

## Understanding

### How it works now

The repository already uses workspace schema 2.0 and a three-entry migration registry. Init and doctor protect `.spectacular.local/` through `.gitignore`, but the planned private idea, GitHub-account, and protected-security paths are not yet an implemented local-state contract. The roadmap still describes a future “v2” workspace migration even though schema 2.0 is already live.

### What changes

This request gathers current Git, workspace, migration, path, and security-boundary evidence. It produces the schema-3 delta, classifies mechanical versus judgment migration work, and recommends whether a separately approved migration specification can safely proceed.

### What stays the same

CLI behavior, workspace layout, schema value, public documentation, GitHub repository configuration, existing request state, private local contents, and production code remain unchanged. The request scaffold's mechanical `last_build` update is the only project-config delta.

## Decisions

- Keep readiness separate from migration apply instead of combining audit and mutation, because evidence must define the later approved specification.
- Record traffic as `parallel` instead of `conditional`, `serialized`, or `unknown`, because no active request or open PR overlaps and every readiness output stays inside this request.

## Milestones

- M1 — Exact Git/workspace baseline and current health are reproducibly documented
- M2 — Shared/private boundary inventory identifies every relevant path, owner, retention class, and leakage condition
- M3 — Schema-3 contract delta and migration classification are complete without modifying implementation surfaces
- M4 — Readiness report gives a reviewable go/no-go recommendation, risks, open questions, and next-spec scope

## Tasks

See `TASKS.md`.

## Dependencies

- Governed by `DEC-022` and `DEC-023`; informed by `SPC-003` without becoming its implementation request.
- `blocked_by`: none at scaffold time.
- `blocks`: the future schema-3 migration specification and its separately approved implementation request.
- `conflicts_with`: none at scaffold time.
- Traffic boundary: request-owned evidence only; any proposed edit outside this folder remains a recommendation until separately authorized.

## Validation

- M1 — run: `git status --short --branch`, `git rev-list --left-right --count origin/main...main`, `spectacular doctor`, `spectacular migrate --list`, and `spectacular status --against-latest`; the recorded report cites exact outputs and commit identities.
- M2 — assert: the boundary matrix covers every discovered live shared/local path, marks tracked/ignored status, reads no protected bodies, and gives a disposition for any leakage signal.
- M3 — judge: the delta names each schema-3 addition/change/removal, compatibility behavior, mechanical/judgment owner, rollback boundary, and verification authority; `git diff --name-only main...HEAD` contains no production implementation surface.
- M4 — judge: the readiness report contains one explicit go/no-go result, unresolved questions with owners, a proposed later SPC scope, and evidence links; no placeholder remains in PLAN/TASKS or delivered artifacts.

## Deliverables

- `artifacts/workspace-inventory.md` — current shared/private path and ownership matrix
- `artifacts/schema-3-delta.md` — proposed contract delta and compatibility window
- `artifacts/migration-manifest.md` — dry-run change classes, recovery boundaries, and verification map
- `RISKS.md` — leakage, compatibility, collaboration, rollback, and tooling risks
- `artifacts/readiness-report.md` — final go/no-go recommendation and proposed next-spec scope
