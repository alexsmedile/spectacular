---
status: verified
priority: high
owner: alex
updated: 2026-08-03
build: b40
docs_impact: required
summary: "Implement SPC-003: route GitHub work through direct, request, or spec-first paths and conclude coordinated work through PRs"
related:
  - PRD.md
source_spec: SPC-003
source_type: spec
source_ref: SPC-003
source_spec_version: 1.0
source_spec_digest: "sha256:3b5644e7e4a6337741f1bb9e4dbb6a2ca0f83a4922991178b2a1227c6ef891ad"
scaffolded_against: 142c983bef0b170ea96104a0b42cd27419d185dc
activated_at: 2026-08-03
activated_by: alex
activated_against: 142c983bef0b170ea96104a0b42cd27419d185dc
docs_impact_evidence: Updated README, command/workflow/integration/configuration guides, skill routing/reference contract, PRD/architecture/agent guidance, capability index, and changelog for v1.37.0
github_pr: "https://github.com/alexsmedile/spectacular/pull/14"
github_pr_opened_at: 2026-08-03
github_pr_head: de5a26087e0ebd9f012307c3a2e0a8e804b54544
issue_resolution: on_merge
---

# Plan — github-work-bridge

## Goal

Join GitHub's collaborative work queue to Spectacular's durable coordination layer while keeping short, self-contained Issue work lean.

## Constraints

- Preserve the approved `direct | request | spec-first` routing and existing spec-derived request behavior.
- GitHub remains authoritative for remote records; Spectacular stores accepted meaning and stable references without copying bodies or comments.
- Remote mutation, security, readiness, and merge gates must fail safely.
- Keep Bash 3.2 compatibility and avoid passthrough-only `gh` wrappers.

## Understanding

### How it works now

Requests can only be scaffolded from approved specs, AFK owns the only PR helper, and the CLI has no general Issue-readiness, GitHub handoff, or reconciliation path. Traffic semantics and GitHub lifecycle principles exist in decisions/specs but are not exposed as a coherent implementation workflow.

### What changes

Add the smallest complete GitHub bridge: semantic work routing, Issue/goal request provenance, reviewer-facing PR manifests, safe draft/ready handoff, and read-only reconciliation. Update the skill and user docs so agents choose direct work, a lean request, or spec-first execution consistently.

### What stays the same

Existing spec-derived request commands and provenance remain compatible. GitHub owns Issues, comments, PRs, checks, permissions, and merge state. Merge, security disclosure, and governance changes remain human-gated; managed repository setup and broader GitHub administration stay deferred.

## Decisions

- Reuse existing request, AFK, traffic, and verification primitives instead of creating a second scheduler or PR engine.
- Add GitHub wrappers only where they combine local lifecycle state with remote evidence or enforce a gate.
- Keep the first slice `observe`/`adapt`; defer managed repository configuration.

## Milestones

- M1 — Add work routing and Issue-readiness interpretation.
- M2 — Extend request provenance to Issue and goal sources.
- M3 — Generalize PR handoff and add read-only reconciliation.
- M4 — Document, dogfood, and verify the complete bridge.

## Tasks

See `TASKS.md`.

## Dependencies

- Approved `SPC-003` and `DEC-018` through `DEC-021`.
- Existing request, AFK Git hygiene, wayfinding traffic, verification, and GitHub CLI primitives.
- `workspace-migration-readiness` is in review and does not mutate these runtime surfaces; recheck before integration.

## Validation

- Bash 3.2 syntax and version guard pass.
- Focused GitHub/request tests cover all three routes, provenance compatibility, closing relationships, readiness gates, offline/permission ambiguity, and reconciliation.
- Full test suite passes.
- Spectacular doctors report no spec, roadmap, lifecycle, link, or GitHub-contract drift.

## Deliverables

- Production CLI/skill behavior for the SPC-003 core bridge.
- Focused automated tests and refreshed workflow/command/configuration documentation.
- Verified request evidence and a reviewer-facing pull request linked to this request and SPC-003.
