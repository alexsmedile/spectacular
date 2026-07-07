---
status: review
updated: 2026-07-07
related:
  - PLAN.md
---

# Verify Log — archive-closure-gate

Validation walk for the b22 closure-gate work. Checks map to PLAN § 6.

## Walk 2026-07-07

| Check | Kind | Evidence | Result |
|---|---|---|---|
| M1 — spec-sync/archive proposals are structured deltas, not prose | observable | `spec-sync.md` § Proposal format + `archive.md` closure-gate step rewritten to `SPEC-DELTA.md` (ADDED/MODIFIED/REMOVED/NONE); grep clean of the old "Want me to proceed with these updates" free-form line | ✅ |
| M2 — open TASKS box blocks archive, naming the box | run | `tests/cli/archive-closure-gate.test.sh` scenario 1 (open box → exit 1, output `✗ tasks`); confirmed in scratch | ✅ |
| M2 — VERIFY without a ✅ walk row blocks | run | scenario 2 (unwalked VERIFY.md → exit 1 `✗ verify`; adding a ✅ VERIFY-LOG row → exit 0) | ✅ |
| M2 — missing SPEC-DELTA blocks; NONE passes | run | scenario 3 (no delta → exit 1 `✗ spec`; `NONE — …` → exit 0) | ✅ |
| M2 — `--override <check> --reason` records into archive_overrides: | run | scenario 4 (double override → exit 0; archived PLAN contains one valid `archive_overrides:` list with both `{check, reason, date}` entries; YAML valid) | ✅ |
| M2 — `--force` does NOT bypass closure checks | run | scenario 6 (`--force` on planned+open → exit 1, `closure gate blocks`) | ✅ |
| M2 — already-archived request is refused, not revalidated | run | scenario 7 (second archive → exit 1 `already archived`) | ✅ |
| M3 — `doctor specs` flags a MODIFIED delta quoting a missing bullet | run | scenario 8 (bad quote → `not found in SPEC.md` warning; valid quote → no warning) | ✅ |
| M3 — full suite green | run | `bash tests/run.sh` → 15/15 files pass (incl. mutator 68, undo 37, new gate file 28) | ✅ |
| Regression — undo after overridden archive drops archive_overrides cleanly | run | `_fm_unset_block` added; scenario 4 undo tail asserts block gone + `related:` intact + no orphan items | ✅ |

## Coherence pass (PLAN § Decisions vs built artifact)

- **block-with-recorded-override** — shipped: each check blocks, `--override` records `archive_overrides:`. ✅
- **delta blocks (ADDED/MODIFIED/REMOVED)** — shipped: `SPEC-DELTA.md` format + `doctor specs` structural validation. ✅
- **`[~]`-with-reason valid closure state** — shipped: scenario 1b proves reasonless `[~]` blocks, `[~] … — reason` passes. ✅

All decisions materialized in the build. All checks pass.
