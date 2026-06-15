---
status: active
updated: 2026-06-15
related:
  - PLAN.md
---

# Tasks — decisions-index

## v1

### M1 — Schema + detection rule
- [x] Update `decisions-rules.md` frontmatter: add `mode: index | flat` enum; update `location:` note to cover both modes
- [x] Document detection rule in `decisions-rules.md`: presence of `decisions/` subfolder = index mode; absence = flat mode (backwards compat)
- [x] Canonicalize index line format: `- **D42** — Short title — one-sentence rationale`
- [x] Canonicalize per-entry file format: `decisions/D<N>.md` with `# D<N> — Title` heading + Context/Decision/Consequences/Session fields
- [x] Update agent read pattern note: always load index (cheap); load `decisions/D<N>.md` on demand

### M2 — `decisions migrate` verb
- [x] Add `decisions migrate` sub-verb to CLI (`spectacular decisions migrate`)
- [x] Parse flat `DECISIONS.md`: extract each `## YYYY-MM-DD —` block, derive slug/number, write `decisions/D<N>.md`
- [x] Rewrite `DECISIONS.md` root as one-liner index after all per-entry files written (never before — no data loss on partial run)
- [x] `--dry-run`: print would-write paths + index preview; write nothing
- [x] Idempotent: if `decisions/` already exists, print already-migrated message and exit 0
- [x] `--help` updated

### M3 — `decide` index-mode write
- [x] Detect index mode in `decide` verb: check for `decisions/` folder
- [x] In index mode: write full ADR prose to `decisions/D<N>.md`; append one-liner to `DECISIONS.md` index
- [x] In flat mode: behavior unchanged (append full block to `DECISIONS.md`)
- [x] Auto-number: read highest D<N> from `decisions/` dir, increment by 1
- [x] `--dry-run` in index mode: preview both the index line and the per-entry file; write nothing

### M4 — `doctor decisions` area
- [x] Add `decisions` area to `doctor` (new area or extend existing)
- [x] Check 1 — mode consistency: if `decisions/` exists, DECISIONS.md must be index-only (no prose blocks)
- [x] Check 2 — no orphans: every index line has a corresponding `decisions/D<N>.md` file
- [x] Check 3 — no stale files: every `decisions/D<N>.md` has a corresponding index line
- [x] Check 4 — numbering: D-numbers are sequential, no gaps, no duplicates
- [x] Update `doctor-areas.md` with decisions area checks

### M5 — Dogfood + ship
- [x] Run `spectacular decisions migrate` on this repo's `.spectacular/DECISIONS.md`
- [x] Add one new decision via `spectacular decide` in index mode; verify D<N>.md written + index appended
- [x] Verify: `grep -c "^\*\*Context" .spectacular/DECISIONS.md` returns 0
- [x] Verify: `ls .spectacular/decisions/ | wc -l` returns correct count
- [x] CHANGELOG entry
- [ ] Plugin bump to target release (deferred — ships with v1.17.0 release commit)

## v2 (deferred)

- [ ] `spectacular decisions` (no args) renders the index as a formatted table (number, date, title, one-liner)
- [ ] `spectacular decision <N>` also accepts bare number (D42 or 42) in addition to slug
- [ ] `doctor decisions --fix` auto-repairs orphan/stale entries
