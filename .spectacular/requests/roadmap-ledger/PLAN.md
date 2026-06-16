---
status: review
priority: high
owner: alex
updated: 2026-06-16
summary: "Single-ledger roadmap — requests carry a stable build id, not a hardcoded target version; one table in ROADMAP.md maps build → planned version. Inserting a request edits one row, not ~14 scattered refs."
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../PRINCIPLES.md
depends-on:
  - cross-request-links
build: b7
---

# Plan — roadmap-ledger

> **Origin (2026-05-30):** field feedback after inserting `policy-engine` at v1.12 forced manual edits to ~14 cross-references + 3 request frontmatters. Root cause: the roadmap uses the absolute version number as a request's *identity*, written into headings, prose, dependency chains, and range labels ("the v1.15→v1.17 ladder"). A single target version is currently duplicated across PLAN frontmatter + a block heading + N prose refs (measured: 61 inline version mentions for 8 future versions). Same class of bug as the "18 docs" hardcoded count — a *derived* fact written down as source.

## Understanding

### How it works now

Every request PLAN carries `target_version: v1.16.0` in frontmatter — a hand-written absolute version number. The same number then leaks into milestone prose ("plugin bump to v1.16.0"), validation lines ("manifests at v1.13.0"), dependency chains ("depends-on the v1.18→v1.20 ladder"), and ROADMAP.md block headings — measured at 61 inline version mentions for 8 future versions across the live ROADMAP alone. When a request is reordered (the `policy-engine` insert at v1.12 required edits to ~14 refs), every copy must be found and updated. Misses are guaranteed; we had three in today's renumber.

The ROADMAP.md is also purely prose — blocks are identified by their `## vX.Y.0` heading, so both the version number and the block text are the source of truth simultaneously.

### What changes

- **Requests get a stable build id** (`build: b7`) stamped by `spectacular new` at creation — immutable, never changes even if the version shifts.
- **A single ledger table** in ROADMAP.md maps `build | slug | title | tier | target-version | status`. This is the *only* place a version number is written. Everything else references requests by slug or build id.
- **`target_version:` removed from request frontmatter** — it was the source of drift. `build:` is the new stable identity field. The version is a ledger read, not a stored copy.
- **ROADMAP prose** converts from `## v1.16.0 — Cross-request awareness` to slug/label references; block headings become derived/rendered output, not source.
- **`spectacular new`** increments a `last_build:` counter in `config.yaml` and stamps `build:` on the new request.
- **`doctor` gains a links check** for stray hardcoded version refs outside the ledger.

### What stays the same

The slug is still the human identity and the primary cross-reference handle. `depends-on:` / `blocks:` (from cross-request-links) use slugs, not build ids. Build ids are the machine-stable handle; slugs are the readable one. The ROADMAP still renders version-labelled blocks — they just derive from the ledger rather than being hand-typed. Shipped history stays in CHANGELOG, not the ledger (ledger = planned runway only). `last_build` lives in `config.yaml`; single-user concurrency is acceptable to ignore for v1.

## 1. Goal

Make the **planned target version a derived value, not a hardcoded string**, by separating *build identity* from *release label* — Apple-style (`CFBundleVersion` build vs `CFBundleShortVersionString` marketing version).

- A request gets a **monotonic build number** (`b1`, `b2`, …) stamped at `spectacular new`, immutable thereafter. The build is the request's stable identity.
- The **version** (`v1.12.0`) is **derived** — the ledger maps build(s) → version. Many builds can share one version (they ship together); default is one request = one build.
- **One ledger** — a single table, living in `ROADMAP.md` itself — is the only place `build → version` lives. Everything else (block prose, dependency chains) references the request by **build or slug**, never by version.

Inserting or reordering a request edits **one row in one table** (and its version column), with zero cascade into prose. This is the **single-canonical-version-source rule** (already in `docs/versioning.md` for the *product's* released version) extended to the *roadmap's planned per-request targets*.

## Resolved design (from 2026-05-30 feedback)

| Concept | What it is | Stability |
|---|---|---|
| **Slug** | `policy-engine` | stable human label |
| **Build** | `b12` — monotonic counter, stamped at `spectacular new` | **immutable identity** |
| **Version** | `v1.12.0` | **derived** from the ledger |
| **Grouping** | many builds → one version (e.g. `b10`+`b11` → v1.10.0) | a ledger column |

- **Counter, not hash** — sortable + human-readable; `last_build` source lives in `config.yaml`.
- **Stamped at creation** — every `spectacular new` takes the next build. If a request later merges into another's release, its build becomes a **gap** — gaps are normal (like skipped Xcode builds), not errors.
- The ledger column for version is the *only* mutable target; the build and slug never change.

## 2. Constraints

- **One ledger, one source.** The position→version mapping lives in exactly one place. No version number is hand-written anywhere a tool could derive it.
- **Slug is identity; number is derived.** Requests reference each other by slug (`depends-on: workspace-v2-spec`), never by `v1.16`. Version is computed from ledger position at read/render time.
- **Stable build/sequence id (optional but preferred).** A request may carry an immutable sequence number assigned at creation — survives reordering, gives a stable handle even before a version is pinned. The *version* is the volatile, derived field; the *sequence id* and *slug* are stable.
- **Backward-compatible read.** Existing `## vX.Y.0` block prose can stay as rendered output, but it must become *generated/derived from the ledger*, not the source. Migration path needed for the current hand-numbered ROADMAP.
- **Don't over-engineer.** The ledger is a markdown table first, not a database. Tooling (a render/check step) can come later; the discipline (one source) is the win.
- **Dovetails with cross-request-links.** That request adds slug-based `depends-on:`/`blocks:` + a doctor links area + ROADMAP-as-source-of-truth enforcement. This request is the *data-model* half (version as derived); cross-request-links is the *relationship* half. They should share the slug-as-identity foundation.

## 3. Milestones

- M1 — **Ledger schema.** Define the single table: columns `seq | slug | title | tier | planned-target | status`. Decide where it lives (top of ROADMAP.md) and whether `seq` (stable) or position (derived) drives ordering. Document the "version is derived, never written elsewhere" rule.
- M2 — **De-duplicate references.** Convert ROADMAP prose + dependency chains from absolute versions to slug/label references ("the contract-prep ladder", "depends-on: workspace-v2-spec"). Remove `target_version:` from request frontmatter entirely; add `build:` as the stable identity field.
- M3 — **Insert/reorder is one edit.** Demonstrate: adding a request = one new ledger row + (if mid-sequence) a re-render; no prose touched. Compare against the policy-engine reslot (the motivating pain).
- M4 — **Render + check.** A read surface renders per-version blocks *from* the ledger (ties into visual-layer's `roadmap` render); `doctor` flags any stray hardcoded version reference outside the ledger (the enforcement that keeps it clean).
- M5 — **Migrate this ROADMAP + ship.** Convert the live ROADMAP.md to ledger-driven; dogfood by reslotting a request and confirming one-row-edit; CHANGELOG + plugin bump.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- **Depends on [[cross-request-links]]** (ships in the same v1.16.0 release; cross-request-links builds the `depends-on:`/`blocks:` slug-identity infra that this request's dependency tracking relies on). cross-request-links ships first or in parallel; roadmap-ledger's M5 migration finalizes both.
- **Render half ties into [[visual-layer]]** (the `spectacular roadmap` render already shipped; M4 extends it to read from the ledger).
- Relates to `docs/versioning.md` (extends the single-source rule to roadmap targets).

## 6. Validation

- M1 — The ledger schema is documented; a reader understands version-is-derived.
- M2 — `grep -c "v1\\.[0-9]" ROADMAP.md` outside the ledger table drops to ~0; request frontmatter no longer carries a source-of-truth `target_version`.
- M3 — Inserting a fixture request touches one ledger row + zero prose lines (vs. the ~14-ref policy-engine reslot).
- M4 — `doctor` flags a deliberately-hardcoded version reference outside the ledger; roadmap render reproduces the current block view from ledger data.
- M5 — Live ROADMAP.md is ledger-driven; a real reslot is a one-row edit; manifests bumped.

## 7. Deliverables

- The roadmap ledger (single position→version table) + its schema doc
- ROADMAP.md converted to slug/label references (no scattered absolute versions)
- `target_version:` removed from request frontmatter; `build: bN` added as stable identity field
- `doctor` check for stray hardcoded version references
- Roadmap render reads from the ledger (coordinated with visual-layer)
- CHANGELOG entry

## Design decisions (2026-06-07, extended 2026-06-14)

- ~~Stable id: sequence vs slug?~~ → **monotonic build counter** (`b1`, `b2`…) as identity + slug as label.
- ~~When is the id stamped?~~ → **at `spectacular new`**; gaps from merges are normal.
- ~~Build column in PLAN frontmatter~~ → `target_version:` **removed entirely**; `build: bN` added as stable identity. No derived mirror — mirrors drift.
- ~~Where does `last_build` live?~~ → `config.yaml`; `spectacular new` increments it. Single-user concurrency ignored in v1.
- ~~Merge with cross-request-links or stay separate?~~ → **Two requests, sequential releases.** cross-request-links shipped the link schema (v1.16.0, 2026-06-08); roadmap-ledger ships the version-as-derived model in the next available slot (v1.17.0, co-shipping with or following cli-debt-removal).
- ~~Gaps/buffers~~ → absent build numbers; no explicit buffer rows in the ledger.
- ~~Shipped history in ledger?~~ → **Planned runway only.** Shipped history stays in CHANGELOG.
- ~~Version-column format for grouped builds~~ → **two rows, same `target-version` value.** Flat table; render groups visually. No merged cells, no comma-separated values.
- ~~Ledger `status` column values?~~ → **`planned | active | shipped`** — release-level states, distinct from request lifecycle (`planned | active | review | verified`). Flips to `shipped` when the version tags.
- ~~Does `spectacular new` auto-add a ledger row?~~ → **No.** `new` stamps `build: bN` + increments `last_build:` only. Human adds the ledger row manually when slotting the request into a version.
- ~~`last_build:` initialization~~ → missing treated as `0`; first `new` stamps `b1`, writes `last_build: 1` silently.
- ~~Build id assignment for existing requests (M2)~~ → sorted by `updated:` date ascending, alpha tiebreaker; `last_build:` set to N after assignment.
- ~~Remove `--target-version` flag in M1 or M2?~~ → **M2.** M1 adds `build:` stamping only; M2 is the `target_version:` removal sweep.
- ~~Tier legend — where does it live?~~ → **ARCHITECTURE.md** (ledger schema section). Not in `docs/versioning.md`. Values: `full` = near-term detailed, `themed` = mid-term directional, `vision` = long-horizon direction-only.
- ~~`spectacular new` output~~ → prints `✓ build id: bN` + "add a row to the ledger in ROADMAP.md when slotting" hint.

## M4 — Insert/reorder demonstration (2026-06-14)

**Before (policy-engine reslot, 2026-05-29):** inserting `policy-engine` at v1.12 required edits to ~14 refs — 3 request PLAN frontmatters (`target_version:`), ROADMAP block heading, prose references in ~6 dependency-chain lines, range labels ("v1.15→v1.17 ladder"), and the ROADMAP summary frontmatter. Three misses were confirmed afterward.

**After (ledger model):** to insert a new fixture request `spec-completeness` at v1.17.0 (bumping `roadmap-ledger` to v1.18.0), the full edit is:

```diff
 | b4 | cli-debt-removal | CLI debt removal | themed | v1.17.0 | planned |
+| b9 | spec-completeness | Spec completeness check | full | v1.17.0 | planned |
-| b7 | roadmap-ledger | Roadmap ledger | full | v1.17.0 | active |
+| b7 | roadmap-ledger | Roadmap ledger | full | v1.18.0 | active |
```

Two rows in one table. Zero prose touched. The `roadmap-ledger` PLAN.md keeps `build: b7` — no frontmatter change. No dependency chains to update. No version refs to grep for.

**Reduction:** ~14 scattered edits → 2 ledger rows (1 insert + 1 version column update). The motivating pain — guaranteed misses on reslots — is structurally eliminated.
