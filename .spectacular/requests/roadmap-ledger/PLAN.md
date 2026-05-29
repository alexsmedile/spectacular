---
status: planned
priority: high
owner: alex
updated: 2026-05-30
summary: "Single-ledger roadmap — requests carry a stable id, not a hardcoded target version; one table in ROADMAP.md maps position → planned version. Inserting a request edits one row, not ~14 scattered refs."
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../PRINCIPLES.md
  - cross-request-links
target_version: tbd
---

# Plan — roadmap-ledger

> **Origin (2026-05-30):** field feedback after inserting `policy-engine` at v1.12 forced manual edits to ~14 cross-references + 3 request frontmatters. Root cause: the roadmap uses the absolute version number as a request's *identity*, written into headings, prose, dependency chains, and range labels ("the v1.15→v1.17 ladder"). A single target version is currently duplicated across PLAN frontmatter + a block heading + N prose refs (measured: 61 inline version mentions for 8 future versions). Same class of bug as the "18 docs" hardcoded count — a *derived* fact written down as source.

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
- M2 — **De-duplicate references.** Convert ROADMAP prose + dependency chains from absolute versions to slug/label references ("the contract-prep ladder", "depends-on: workspace-v2-spec"). Remove `target_version:` from request frontmatter (or make it a derived/advisory mirror, not a source).
- M3 — **Insert/reorder is one edit.** Demonstrate: adding a request = one new ledger row + (if mid-sequence) a re-render; no prose touched. Compare against the policy-engine reslot (the motivating pain).
- M4 — **Render + check.** A read surface renders per-version blocks *from* the ledger (ties into visual-layer's `roadmap` render); `doctor` flags any stray hardcoded version reference outside the ledger (the enforcement that keeps it clean).
- M5 — **Migrate this ROADMAP + ship.** Convert the live ROADMAP.md to ledger-driven; dogfood by reslotting a request and confirming one-row-edit; CHANGELOG + plugin bump.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- **Shares foundation with [[cross-request-links]]** (slug-as-identity). Decide ordering: this could land first (data model), or merge with it. Likely sequence them adjacent.
- **Render half ties into [[visual-layer]]** (the `spectacular roadmap` render reads the ledger).
- Relates to `docs/versioning.md` (extends the single-source rule to roadmap targets).
- No hard blocker; could slot anywhere in the runway. Target left **tbd** deliberately — *this request is about not pinning versions prematurely*, so it dogfoods its own principle by staying unpinned until the ledger exists to assign it.

## 6. Validation

- M1 — The ledger schema is documented; a reader understands version-is-derived.
- M2 — `grep -c "v1\\.[0-9]" ROADMAP.md` outside the ledger table drops to ~0; request frontmatter no longer carries a source-of-truth `target_version`.
- M3 — Inserting a fixture request touches one ledger row + zero prose lines (vs. the ~14-ref policy-engine reslot).
- M4 — `doctor` flags a deliberately-hardcoded version reference outside the ledger; roadmap render reproduces the current block view from ledger data.
- M5 — Live ROADMAP.md is ledger-driven; a real reslot is a one-row edit; manifests bumped.

## 7. Deliverables

- The roadmap ledger (single position→version table) + its schema doc
- ROADMAP.md converted to slug/label references (no scattered absolute versions)
- `target_version` demoted from request-frontmatter source to derived/advisory (or removed)
- `doctor` check for stray hardcoded version references
- Roadmap render reads from the ledger (coordinated with visual-layer)
- CHANGELOG entry

## Open questions (resolve in M1)

**Resolved 2026-05-30:**
- ~~Stable id: sequence vs slug?~~ → **monotonic build counter** (`b1`, `b2`…) as identity + slug as label.
- ~~When is the id stamped?~~ → **at `spectacular new`**; gaps from merges are normal.

**Still open:**
- **Build column in PLAN frontmatter** — request frontmatter gains `build: b12` (stable); `target_version:` is demoted to derived/advisory or removed. Confirm the exact frontmatter shape.
- **Where does `last_build` live?** `config.yaml` counter that `spectacular new` increments. Concurrency (two `new` calls) — acceptable to ignore for a single-user/team workflow?
- **Merge with cross-request-links or stay separate?** Both need stable-id-as-identity (build/slug). One request or two adjacent ones?
- **How are gaps/buffers represented** in the ledger — explicit buffer rows, or just absent build numbers?
- **Does the ledger include shipped history** (builds → already-released versions), or only the planned runway (history → CHANGELOG)?
- **Version-column format for grouped builds** — how does the ledger show `b10`+`b11` both → v1.10.0?
