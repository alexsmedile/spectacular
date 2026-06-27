---
status: active
updated: 2026-06-02
related:
  - PLAN.md
  - TASKS.md
---

# Session — cross-request-links

## Current state

Promoted `planned → active` (2026-06-02), the next release in roadmap order (v1.13.0). `## Understanding` filled; @Implementation gate satisfied. imagine-mode shipped as preview just before pivoting here.

This is the advisory cross-request awareness feature: `depends-on:` / `blocks:` frontmatter siblings to `related:`, an inverse-link resolver (computed, never stored), `doctor links` extension, and a `status` advisory surface. **Advisory only** — no locking.

## Active task

**M1 — Schema extension** (doc-only):
- Document `depends-on:` / `blocks:` in ARCHITECTURE.md alongside `related:`
- Specify the computed-not-stored inverse rule

Then M2 (resolver) → M3 (`doctor links` + doctor-memory staleness side-rider) → M4 (`status` advisory + `new` prompt) → M5 (examples + ship v1.13.0).

## Blockers

None.

## Next actions

1. M1: write the ARCHITECTURE.md frontmatter-schema section (3 relationship fields + inverse rule).
2. M2: implement the inverse-link resolver in `cli/spectacular` (mechanical — reads all PLAN frontmatter, computes the graph).
