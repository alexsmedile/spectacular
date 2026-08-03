---
updated: 2026-08-03
status: dry-run-design
source_schema: "2.0"
target_schema: "3.0"
---

# Schema-3 migration manifest proposal

No step below is authorized to apply. Ordering is dependency-first.

| Stage | Proposed change | Class | Recovery / gate | Verification authority |
|---|---|---|---|---|
| P0 | Fetch and compare exact target/local commits; inspect traffic | Read-only preflight | Stop on divergence, overlap, undeclared access | Git SHAs, open PR/request evidence |
| P1 | Approve the schema-3 contract and resolve open questions | Human judgment | Product/security HITL | Approved SPC and DEC references |
| P2 | Correct roadmap/schema terminology | Canonical docs | Snapshot; ordinary revert | Link/spec/roadmap doctors |
| P3 | Add older/equal/newer schema guard shared by mutators | Additive production guard | Revert commit | CLI fixtures for all three relations |
| P4 | Introduce optional schema-3 fields/path validators under schema 2.0 | Additive soak | Revert commit; absence stays valid | Old/new-shape doctor tests and dogfood |
| P5 | Add filename-only tracked-local detection and protected refusal | Additive security guard | Non-bypassable; security review | Synthetic path fixtures; no body output |
| P6 | Add `v20-to-v30` registry entry in dry-run mode | Migration scaffold | No live apply; manifest comparison | Registry doctor and dry-run golden output |
| P7 | Classify legacy `debug/` and root artifacts | Judgment cleanup, separate scope | User/security decision per artifact | Clean root/plural-path doctor checks |
| P8 | Flip `CURRENT_SCHEMA` to 3.0 and enable apply only after soak | Breaking transition | Exact baseline, isolated branch, snapshots, commit-level revert/abandon | Full Bash/CLI suite and schema-2 fixture migration |
| P9 | Migrate this repository and open PR | Governed apply | Human PR/merge gate | Doctor zero schema errors; request verification |

## Mechanical candidates

- Version comparisons and refusal paths after the contract is frozen.
- Creating missing non-sensitive shared directories declared required by schema 3.
- Adding optional frontmatter fields only when deterministic defaults exist.
- Bumping `workspace_schema` after all shape checks pass.

## Judgment or separate-authority work

- Choosing the breaking schema delta.
- Interpreting/moving legacy content and conflicting destinations.
- Handling any tracked private path or Git-history exposure.
- Rewriting canonical prose and accepted decisions.
- Destructive local conversion, history rewriting, remote deletion, merge, or disclosure.

## Required fixtures

- Valid schema 2.0 with no optional additions.
- Schema 2.0 carrying optional schema-3 fields/paths.
- Valid schema 3.0.
- Schema newer than the running CLI.
- Missing `.spectacular.local/` ignore protection.
- Tracked local pathname with synthetic non-secret content.
- Offline GitHub and missing authentication.
- Clean fast-forward, diverged main, isolated unpublished rebase, and conflict.
- Partial/failed migration with unchanged schema marker.
- Legacy singular directory and root-artifact collisions.

## Completion invariant

Never bump `workspace_schema` merely because steps ran. Bump only after required shape, forbidden paths, links, permissions, and collection invariants pass. Failure preserves the branch and manifest and leaves the prior schema marker unchanged.

