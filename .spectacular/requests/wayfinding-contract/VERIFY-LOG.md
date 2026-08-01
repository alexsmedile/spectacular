---
updated: 2026-08-01
---

# Verify log — wayfinding-contract

## 2026-08-01 16:43 CEST — walk (8 passed, 0 blocked, 0 skipped)

- ✓ [assert] Canonical aliases resolve to normalized IDs and ambiguous naked numbers require context — resolver and focused fixtures confirm canonical, shorthand, contextual, and refusal paths.
- ✓ [assert] Canonical references, padding, reserved prefixes, and compatibility behavior are documented — `canonical-ids.md` defines all eight prefixes, explicit-prefix priority, contextual numbers, 3+ digit padding, and legacy aliases.
- ✓ [exec] Typed record, spec lifecycle, fog/frontier, migration, and doctor fixtures pass — `bash tests/cli/wayfinding-contract.test.sh` exit 0; 53 passed, 0 failed.
- ✓ [exec] Bash syntax is valid — `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit` exit 0.
- ✓ [exec] Guarded version strings are consistent — `scripts/hooks/pre-commit --check` exit 0; all guarded surfaces report 1.35.0.
- ✓ [exec] Complete regression suite passes — `bash tests/run.sh` exit 0; 18 test files passed, 0 failed.
- ✓ [assert] Migration records its mapping and archived originals — apply path writes `mapping.tsv` and copies originals beneath `.spectacular/archive/id-migrations/<timestamp>/` before renaming.
- ✓ [assert] No automatic repository migration was performed — `.spectacular/archive/id-migrations/` has no repository migration run.

**Coherence:** All five PLAN decisions are present in the shipped contract: distinct decisions/questions/ideas domains; canonical IDs plus aliases; `IDEA` canonical with `IDE` compatibility; additive archive-first migration; and `unconfirmed | current | deprecated` specification state.

**Outcome:** verified
