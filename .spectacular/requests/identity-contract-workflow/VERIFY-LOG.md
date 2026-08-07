---
updated: 2026-08-07
summary: "Identity-contract workflow verification evidence"
---

# Verify log — identity-contract-workflow

**Outcome:** verified

## Evidence

- `bash tests/run.sh` — all 29 CLI and pipeline test files passed.
- `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit` — passed.
- `scripts/hooks/pre-commit --check` — all guarded version strings consistent.
- Focused lifecycle coverage proves a merged spec commit is required, an unmerged contract is rejected, and an execution branch must contain the merge.
- Migration coverage proves preview does not write, apply archives a UUIDv7 mapping receipt, rewrites references, and preserves the numeric alias as a read-only resolver input.

## Compatibility boundary

Legacy numeric records remain readable until explicitly migrated. New durable records receive UUIDv7 IDs and slug paths; `spec` and `request` vocabulary remains unchanged.
