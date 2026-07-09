---
request: spec-audit-mode
build: b11
updated: 2026-07-09
---

# Spec delta — spec-audit-mode (b11)

The `specs` doctor area gains a mechanical frontmatter schema check on flat
capability specs. Updates the **Substrate doctor** bullet, which already carried
a forward-reference to this build.

### MODIFIED
- specs/index.md :: "The `specs` area validates `SPEC-DELTA.md` integrity (MODIFIED/REMOVED quote an existing bullet; ADDED isn't a duplicate) as the primary drift signal (v1.28.0+), with SPEC.md date-drift vs the newest archive (v1.18.0+) kept as a backstop — a heuristic blind to work that ships outside the request lifecycle (see `b11: spec-audit-mode` on the roadmap)." -> "The `specs` area validates `SPEC-DELTA.md` integrity (MODIFIED/REMOVED quote an existing bullet; ADDED isn't a duplicate) as the primary drift signal (v1.28.0+), plus a **frontmatter schema check** on each flat `specs/<cap>.md` (b11): required keys `status, updated, summary, related`, ISO `updated`, closed status enum `draft|published|deprecated`, and `version:` required iff `status: published` — mechanical, warning-class, `index.md` skipped, `related:` resolution left to the `links` area. SPEC.md date-drift vs the newest archive (v1.18.0+) is kept as a backstop."
