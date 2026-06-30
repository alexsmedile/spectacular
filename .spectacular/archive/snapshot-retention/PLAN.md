---
status: archived
priority: medium
owner: alex
updated: 2026-06-30
build: b16
summary: "Snapshot system: bound the canonical allowlist (add DESIGN.md), couple @vN to frontmatter version: so they stop drifting, add configurable retention (default keep 3) with a doctor-driven prune, and rename the store to _snapshots/ with a configurable folder name + opt-in gitignore so the history layer stops bloating both the tree and git."
related:
  - PRD.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
  - ../../SPEC.md
archived: 2026-06-30
---

# Plan — snapshot-retention

> **Origin (2026-06-28):** The snapshot system (`spectacular snapshot <file>`,
> `cmd_snapshot`) freezes canonical docs to `.spectacular/snapshots/<DOC>/@v<N>.md`
> and bumps the live doc's `version:`. It works, but a v1.22 audit surfaced three
> issues: (1) the `@vN` filename and the doc's `version:` frontmatter **drift**
> apart — `@vN` is a plain counter, `version:` is MAJOR.MINOR, so `@v3.md` can hold
> version `2.0`; (2) snapshots accumulate **forever** — every canonical-doc edit can
> leave a new `@v<N>.md`, with no retention/cleanup, bloating the tree with files no
> one reads; (3) the snapshot-able set is an allowlist (good — keep it) but doesn't
> include `DESIGN.md`, which a project may want to snapshot even though Spectacular
> doesn't otherwise manage it. Full audit lives in repo-root `TODO.md` §"Snapshot
> cleanup / retention".

## 1. Goal

Make the snapshot system tidy and bounded:

1. **Keep the allowlist** — snapshotting stays restricted to canonical docs (no free-for-all). **Add `DESIGN.md`** to it as a snapshot-able doc even though it isn't otherwise managed by Spectacular.
2. **Stop the drift** — when the doc has a `version:` frontmatter field, the snapshot filename **corresponds to it** (`@v1.2.md` ⇄ version `1.2`). Only fall back to the plain counter when no `version:` is present (e.g. DESIGN.md without frontmatter).
3. **Configurable tiered retention** — keep `@v1` (origin) + one snapshot per calendar period (default **monthly**) + the most-recent **N (default 3)**; auto-clean the rest. Configurable in `config.yaml`. The recent tier sorts by `@vN` ordinal; the periodic tier buckets by each snapshot's recorded `updated:` frontmatter date — **never** by filesystem mtime (drifts on clone/restore).
4. **Rename + configurable folder, opt-in gitignore** — rename the store `snapshots/` → **`_snapshots/`** (the `_`-prefix marks it as a scanner-skip / non-content layer, consistent with `_archive/`). Make both the **folder name** and whether to **gitignore it by default** configurable in `config.yaml`. Snapshots are history that may not need to live in git — let projects opt out.

## 2. Constraints

- **Bash CLI only** — no new runtime deps (STACK.md). Snapshot + retention are deterministic mutators → CLI-owned (v0.7.0 mutation principle). The skill may *suggest* a prune; it never deletes.
- **Allowlist stays closed.** This request does NOT open snapshotting to arbitrary docs. It adds exactly one entry (`DESIGN.md`) to `is_canonical_doc`. Request docs (PLAN/TASKS), POLICY, memory, ideas remain non-snapshot-able by design — they archive whole with their request.
- **Deletion is destructive — recoverable by default.** Prune moves to `.trash/` (or git-rm in a git repo), not `rm -f`. Always dry-run / show-what-would-go before deleting. Never touch the live canonical file.
- **Back-compat.** Existing `@v1.md` / `@v1.0.md` snapshots on disk must keep working — the parser already tolerates both; retention must too. No forced rename of historical snapshot *files*.
- **Folder rename is a migration, not a break.** All 15 hardcoded `snapshots/` path references collapse to one config-sourced `$snap_root` (default `_snapshots`). `cmd_snapshot` reads/writes the configured folder; `check_snapshots` recognizes both old (`snapshots/`) and new (`_snapshots/`) as it already does for the pre-v1.5 legacy layout, and `--fix` migrates old → new. A workspace mid-migration must never lose snapshots.
- **Folder name + gitignore are config, with sane defaults.** Default folder `_snapshots`; default gitignore = **false** (snapshots stay in git unless opted out) — flag it as an M-question, since the user raised gitignoring as desirable. `init` and `doctor --fix` honor the config when scaffolding `.gitignore`.
- **Version coupling is best-effort, not a hard fail.** If `version:` is malformed or absent, fall back to the counter and carry on — never `die`.

## Understanding

### How it works now

`cmd_snapshot` (cli/spectacular:4395):
1. **Gate** — `is_canonical_doc "$file"` (cli:1601) allows: `PRD/SPEC/ARCHITECTURE/PRINCIPLES/ROADMAP/STACK/DECISIONS/AGENTS.md` + `config.yaml` at `.spectacular/` root, plus any `specs/<cap>/SPEC.md`. Everything else → `return 1` → snapshot refused. **DESIGN.md is not in the list.**
2. **Path** — derives `snap_dir=.spectacular/snapshots/<rel_dir>/<DOC>` mirroring sub-paths (cli:4446-4453).
3. **Next N** — scans existing `@v<N>.<ext>` in new + legacy locations, takes `max(N)+1` as a **plain integer counter** (cli:4456-4481). This is the source of drift.
4. **Idempotence** — body-only md5 compare (frontmatter excluded); no-op if unchanged (cli:4486-4495).
5. **Version bump** — reads `version:` frontmatter, bumps `minor+1` (or `major+1` with `--major`), writes snapshot to `@v<next_n>.<ext>`, then `fm_set version` + `fm_touch` on the live doc (cli:4497-4520).

The bug: step 3 (`@vN` = file counter) and step 5 (`version:` = MAJOR.MINOR) are computed **independently**. First snapshot aligns them (`@v1` ↔ `1.0`) but any `--major` bump or hand-set `version:` desyncs them permanently.

`check_snapshots` (cli:7043) — doctor's `snapshots` area: detects legacy root-level snapshots (pre-v1.5 layout) and version-sequence gaps (`@v1` + `@v3` without `@v2`). **No retention check.**

`config.yaml` — has `naming:`, `required_files:`, `agents:`, `skills:` blocks. **No `snapshots:` block.**

### What changes

1. **`is_canonical_doc`** (cli:1601) — add `DESIGN.md` to the root-doc case alongside `PRD.md|ARCHITECTURE.md|…`. One line.
2. **`cmd_snapshot` filename** (cli:4480-4481) — when the live doc has a parseable `version:`, name the snapshot `@v<current_ver>.<ext>` (the version the copied content **is**, per decision 1), then bump the live doc. When absent (DESIGN.md without frontmatter), fall back to the `max(N)+1` integer counter. The "next N" scan stays only as the fallback path.
3. **Retention** — new `cmd_snapshot_prune` (or `doctor --fix snapshots`): per doc, compute the union of three tiers — origin (`@v1`), periodic (newest per `period` bucket by `updated:` date), recent (newest `keep` by `@vN`) — and move everything *outside* the union to `.trash/` (git-rm if tracked). Dry-run by default; `--fix`/`--confirm` to apply.
4. **`config.yaml`** — read a `snapshots:` block:
   ```yaml
   snapshots:
     folder: _snapshots   # store dir name (default _snapshots)
     keep: 3              # recent-tier count per doc (default 3)
     period: month        # periodic tier: month | week | off (default month)
     gitignore: false     # gitignore the store by default? (default false)
   ```
   Parsers mirroring `config_workspace_schema()` (cli:1622): `config_snapshots_folder()`, `config_snapshots_keep()`, `config_snapshots_period()`, `config_snapshots_gitignore()`. All default when the block/field is absent. `period: off` collapses retention to origin + recent (the simple two-tier scheme).
5. **`check_snapshots`** (cli:7043) — add a retention check (info/warning when a doc has prunable snapshots outside the kept tiers) AND folder-migration detection (old `snapshots/` present while config wants `_snapshots/` → `--fix` renames). Surfaced in `doctor snapshots`, relayed by `status`.
6. **Path resolution** — replace the 15 hardcoded `snapshots/` references in `cmd_snapshot` + `check_snapshots` + the legacy-migration block (cli:8719) with a single `$snap_root` resolved from `config_snapshots_folder()`.
7. **`.gitignore` handling** — when `gitignore: true`, `init` scaffolds and `doctor --fix` ensures a `.spectacular/<folder>/` ignore line; when false, ensure it's NOT ignored.

### Folder rename: `snapshots/` → `_snapshots/`

The store moves to **`_snapshots/`** (underscore). Rationale: the `_`-prefix marks a non-content / scanner-skip layer (`_archive/`, `_backups/`), and snapshots are exactly that — a history sidecar, not live workspace content. This **reverses** the original "keep `snapshots/`" lean: doctor still checks the folder (by explicit path, not by scanning), and gitignorability now matters more than the "doctor reads it" argument. Migration: 15 path references in `cli/spectacular` move to a single resolved `$snap_root` (sourced from config, default `_snapshots`), and `doctor --fix snapshots` renames an existing on-disk `snapshots/` → `_snapshots/` (git-mv if tracked), the same pattern as the pre-v1.5 legacy-layout migration already in `check_snapshots`.

### What stays the same

The allowlist model (still closed, one entry added), the `<store>/<DOC>/@v<N>.md` layout shape, the body-only idempotence compare, the legacy-back-compat parsing, and the CLI-mutates/skill-suggests split.

## 3. Milestones

### M1 — Allowlist: add DESIGN.md
- Add `DESIGN.md` to `is_canonical_doc` root-doc case.
- Test: `is_canonical_doc .spectacular/DESIGN.md` → 0; a non-canonical doc still → 1.

### M2 — Version coupling (stop the drift)
- `cmd_snapshot`: when `version:` is parseable, snapshot filename = `@v<current_version>.<ext>` (the version the content **is**, per decision 1), then bump the live doc; else fall back to integer counter.
- Idempotence + `--major` paths still work.
- Tests: snapshot a doc at `version: 1.3` → produces `@v1.3.md`, live doc then reads `1.4`; snapshot DESIGN.md with no frontmatter → `@v1.md` then `@v2.md` (counter).

### M3 — Tiered retention config + prune
- `config_snapshots_keep()` + `config_snapshots_period()` parsers; defaults 3 / `month`.
- Tier computation: origin (`@v1`) ∪ periodic (newest per `updated:`-date bucket — `YYYY-MM` for month, `%G-W%V` for week; `off` skips this tier) ∪ recent (newest `keep` by `@vN`). Bucket dates read from snapshot frontmatter `updated:`, parsed as plain string prefixes (no date math needed — `2026-06-28` → `2026-06` is a substring; week needs `date -j`/`date -d`, guard both BSD/GNU).
- `cmd_snapshot_prune` (and/or `doctor --fix snapshots`): move snapshots outside the union to `.trash/` (git-rm if tracked), dry-run default.
- Tests:
  - `period: off`, `@v1..@v6` + keep 3 → keep `@v1 @v4 @v5 @v6`, prune `@v2 @v3`.
  - `period: month` with 4 same-month snapshots + keep 1 → keep `@v1` (origin) + that month's newest + last 1 (may overlap).
  - Snapshots spanning 3 months → one survivor per month retained even when outside the recent window.

### M4 — Folder rename + configurable name + gitignore
- Collapse 15 hardcoded `snapshots/` refs → one `$snap_root` from `config_snapshots_folder()` (default `_snapshots`).
- `config_snapshots_folder()` + `config_snapshots_gitignore()` parsers.
- `check_snapshots`: detect old `snapshots/` while config wants `_snapshots/`; `--fix` renames (git-mv if tracked) without data loss.
- `gitignore: true` → `init` + `doctor --fix` add `.spectacular/<folder>/` to `.gitignore`; `false` → ensure not ignored.
- Tests: snapshot with default config lands in `_snapshots/`; migration moves existing `snapshots/` → `_snapshots/` losslessly; gitignore toggle adds/removes the line.

### M5 — Doctor retention check + docs
- `check_snapshots`: flag docs over `keep`.
- Update `docs/commands.md`, `docs/configuration.md` (new `snapshots:` block — folder/keep/period/gitignore + the tiered-retention explanation), ARCHITECTURE.md versioning section + any `snapshots/` path mentions, SPEC.md verb list (+ `specs/cli/SPEC.md` if a CLI capability spec exists by then).
- VERIFY-LOG.

## Decisions (resolved at grill, 2026-06-28)

1. **Snapshot filename = the content's *current* version.** The snapshot is a copy of the doc *before* the bump, so it's named the version the copied content actually is (`@v1.3.md` for a doc at `1.3`), *then* the live doc bumps to `1.4`. Honest: filename = content's version. (M2.)
2. **`gitignore` default = `false`.** Snapshots stay committed unless a project opts out (`gitignore: true`). Least surprise vs current behavior where all of `.spectacular/` is committed. (M4.)
3. **Migrate THIS repo's `snapshots/` → `_snapshots/` with b16.** Run M4's `--fix` migration on this workspace when b16 ships and commit the rename — dogfoods the migration path on real data before users hit it. (M4, ship step.)

4. **Retention = tiered (origin + periodic + recent).** A generational scheme, like backup rotation. Three tiers, all survivors unioned (a snapshot kept by *any* tier stays):
   - **Origin** — always keep `@v1`.
   - **Periodic** — keep the newest snapshot in each calendar period, period configurable, **default monthly** (`weekly` also valid). Keyed off the snapshot's recorded `updated:` frontmatter date (every snapshot has one — stable across clone, unlike mtime), bucketed by `YYYY-MM` (monthly) / ISO week (weekly).
   - **Recent** — always keep the newest `keep` (default 3) by `@vN` ordinal.

   Prune deletes only snapshots no tier claims. Example, monthly + keep 3, doc with snapshots across 4 months `@v1(Jan) @v2(Jan) @v3(Feb) @v4(Mar) @v5(Apr) @v6(Apr) @v7(Apr)`:
   - origin → `@v1`
   - periodic (newest per month) → `@v2`(Jan) `@v3`(Feb) `@v4`(Mar) `@v7`(Apr)
   - recent (last 3) → `@v5 @v6 @v7`
   - **retained** = `@v1 @v2 @v3 @v4 @v5 @v6 @v7` minus none here is small; with more intra-month churn (e.g. 10 snapshots in Apr) only Apr's newest + last-3 + monthlies survive, the rest go.

   Rationale: origin = "where it started," periodic = "a thread through the whole history at coarse resolution," recent = "fine detail of how it got here lately." Bounds growth to ≈ `1 + months_alive + keep` per doc instead of unbounded.

## Decisions (cont.)

5. **Prune target = `git rm` if tracked, else `.trash/`** (resolved 2026-06-30, before build). In a git repo the snapshot's content survives in history, so `git rm` is sufficient and avoids `.trash/` clutter; untracked snapshots or non-git workspaces fall back to `.spectacular/.trash/`. Either way the live canonical file is never touched and prune is dry-run by default. (M3.)

## Deferred (TODO, not this request)

6. **Does git make snapshots redundant?** The bigger "do we need snapshots vs `git show`" question stays in repo-root `TODO.md`. b16 only makes the existing mechanism tidy + bounded.
