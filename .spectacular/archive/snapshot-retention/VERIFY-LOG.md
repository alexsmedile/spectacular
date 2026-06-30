---
request: snapshot-retention
build: b16
verified: 2026-06-30
verifier: alex (via spectacular skill)
artifact: PLAN § Validation + per-milestone tests
result: pass
---

# Verify log — snapshot-retention

Validation walk against `PLAN.md § 3. Milestones` (Validation criteria per milestone).
Every milestone has a runnable check in `tests/cli/snapshot-retention.test.sh`
(31 assertions, all passing) plus the existing suite (445 assertions across 12
files, 0 failures after the mutator snapshot scenario was updated to the new
contract).

| Milestone | Check | Evidence | Result |
|---|---|---|---|
| M1 — allowlist +DESIGN.md | DESIGN.md snapshot-able; non-canonical refused | `is_canonical_doc` adds `DESIGN.md`; test Scenario A | ✅ |
| M2 — version coupling | versioned doc → `@v<ver>`; version-less → counter; idempotence + `--major` intact | `cmd_snapshot` names by current version, counter fallback, no version injected on version-less docs; Scenario A + mutator Scenario 4 | ✅ |
| M3 — tiered retention + prune | origin ∪ periodic ∪ recent union; prune dry-run/apply; git-rm-if-tracked-else-trash | `snapshot_retained_set` + `cmd_snapshot_prune`; Scenarios B (off), C (monthly churn), D (git) — match PLAN's worked examples exactly | ✅ |
| M4 — folder rename + gitignore | 15 `snapshots/` refs → `$snap_root`; migration `--fix`; gitignore toggle | `snapshot_store_root`, config parsers, doctor migration + gitignore-drift `--fix`; Scenarios E (migration lossless), F (gitignore toggle) | ✅ |
| M5 — doctor + docs | retention/migration/gitignore checks in `check_snapshots`; docs updated | doctor `snapshots` area flags all three; commands.md, configuration.md (`snapshots:` block), ARCHITECTURE.md, SPEC.md, versioning.md updated | ✅ |

## Decisions implemented

1. **Filename = content's current version** — snapshot named `@v<current>`, then live doc bumps. ✅ (M2, Scenario A: PRD@1.3 → `@v1.3.md`, live → 1.4)
2. **gitignore default false** — store stays committed unless opted out. ✅ (`config_snapshots_gitignore` defaults false; Scenario F)
3. **Dogfood-migrate this repo** — `snapshots/` → `_snapshots/` run live. ✅ (18 files git-mv'd losslessly; doctor warning cleared)
4. **Tiered retention (origin + periodic + recent)** — generational union. ✅ (Scenarios B/C verify the worked examples)
5. **Prune target = git-rm if tracked, else .trash/** (resolved pre-build). ✅ (Scenario D: tracked → 2 staged deletions, no `.trash`; Scenario B: untracked → `.trash`)

## Bugs found + fixed during build (test-before-claim caught all)

1. **Version injection on version-less docs** — first cut set `version: 1.0` on docs with no `version:` field, switching them from counter to version mid-stream. Fixed: never impose a version; version-less docs stay on the counter, live frontmatter untouched.
2. **bash 3.2 has no `globstar`** — `shopt -s globstar` errored, doc-dir scan found nothing. Fixed: enumerate known nesting depths explicitly (`*/`, `*/*/`, `*/*/*/`), matching the existing `check_snapshots` pattern.
3. **`grep -v` exit-1 on empty result** — removing the sole `.gitignore` line left it in place (the `&& mv` never fired). Fixed: `|| true` before `mv`.
4. **`set -e` killed `init`** — `snapshot_gitignore_sync` returns 1 as an "already in sync" status; under `set -e` that aborted init after `.gitignore`. Fixed: `snapshot_gitignore_sync || true` at the init call site.
5. **Generic `.gitignore` --fix handler shadowed the snapshots one** — the `.spectacular.local/` handler matched snapshots-area findings first and `continue`d past them. Fixed: guard it with `[[ "$area" != "snapshots" ]]`.
6. **False version-sequence gap on mixed counter+version dirs** — after the first version-coupled snapshot, a dir with old `@v1/@v2/@v3` + new `@v1.6` tripped the gap check (parsed `@v1`→1.1, saw a gap to 1.6). Fixed: skip gap detection when a dir mixes counter (`@vN`) and version (`@vX.Y`) names — a b16 transition guard. Scenario G covers it.

## Ponytail review

- **One shared retention function** (`snapshot_retained_set`) feeds both `cmd_snapshot_prune` and the doctor retention check — no duplicated tier logic.
- **One shared gitignore function** (`snapshot_gitignore_sync`) feeds init + doctor --fix.
- **One config-field parser** (`config_snapshots_field`) with four thin wrappers, not four awk blocks.
- `snapshot_store_root` resolves the configured-vs-legacy folder in one place; the 15 hardcoded paths collapsed to it.

## Verdict

All 5 milestones pass; all 5 decisions implemented; 6 bugs caught + fixed before
ship; dogfood migration ran clean on this repo's 18 snapshots. Request → verified.
