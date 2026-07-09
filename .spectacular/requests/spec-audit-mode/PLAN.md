---
status: review
priority: medium
owner: alex
updated: 2026-07-09
build: b11
summary: "Spec frontmatter schema check: doctor specs validates required keys, ISO date, closed status enum, and conditional version on flat specs/<cap>.md files — mechanical, deterministic, no semantic NLP"
related:
  - ../../PRD.md
  - ../../specs/index.md
  - ../../roadmaps/index.md
---

# Plan — spec-audit-mode (frontmatter schema check)

> **Pivot (2026-07-09, grill-me).** This request was originally a *semantic
> coverage audit* (orphan capability bullets / orphan spec files / stale specs —
> heuristic, fuzzy, and written against the pre-OKF `specs/<slug>/SPEC.md`
> layout). That design never settled ("heuristics still to settle") and its paths
> went stale after the OKF v2.0 restructure flattened specs to `specs/<cap>.md`.
> Grilled down to what "good `doctor specs`" actually means: a **mechanical
> frontmatter/format schema check** — deterministic, high-confidence, zero false
> positives, no NLP. The semantic audit is dropped, not deferred.

## Goal

Give `spectacular doctor specs` a real answer to "is this spec file *well-formed*?"
— validate the frontmatter schema of each flat capability spec so drift in the
signal layer (missing keys, bad dates, typo'd status, unversioned published
contracts) is caught mechanically, the way `doctor` catches substrate problems
everywhere else.

## Constraints

- **Mechanical only.** No semantic matching, no NLP. Every finding is a
  deterministic frontmatter fact — the class of check with zero false positives.
- **Consistency with existing code governs.** Severity model mirrors
  `check_frontmatter`: broken/unparseable frontmatter = `error`; missing key or
  bad value = `warning`. No stricter bar for specs than the rest of `doctor`.
- **Findings surface under `specs`.** A user running `doctor specs` sees spec
  problems — so this extends `check_specs`, not `check_frontmatter`.
- **No `related:` resolution here.** Path resolution is already fully owned by
  `check_links` (root-aware target checks). We assert the key is *present*, never
  re-resolve targets — no duplication.
- **Two-schema fork is intentional and permanent.** Capability specs use their
  own required set (`status, updated, summary, related`, conditional `version`),
  which deliberately differs from the root-anchor set in `check_frontmatter`
  (`version, updated, summary`). Inlined into `check_specs`, no shared helper, no
  schema-unify debt marker — we accept the fork.
- Bash CLI only, no new runtime deps (STACK.md).

## Understanding

### How it works now

`check_specs` validates specs/ *structure* (dir present, index.md parseable, flat
OKF layout, SPEC-DELTA integrity, date-drift vs archive) and checks that each flat
`specs/*.md` merely *has* a frontmatter delimiter — but never validates the keys
inside it. `check_frontmatter` validates key contents, but scans **only
`.spectacular/*.md`** root anchors; it never descends into `specs/`. So the
frontmatter *contents* of capability specs are unvalidated — a real gap, already
visible in the live repo (`doc-engine.md` has no `version:`, `roadmap.md` does).

### What changes

`check_specs`'s flat-file loop gains a schema layer per `specs/*.md` (excluding
`index.md`):

1. **Required keys** — `status, updated, summary, related` present → `warning` if missing.
2. **ISO date** — `updated` matches `YYYY-MM-DD` → `warning` if not (reuses the `check_frontmatter` pattern).
3. **Closed status enum** — `status ∈ {draft, published, deprecated}` → `warning` if outside.
4. **Conditional version** — `version:` required **iff** `status: published` (a published spec is a contract; a draft has nothing to version) → `warning` if published-without-version.

Structurally-broken frontmatter (no delimiter) stays the existing `warning`
(delimiter branch) — the file can't be parsed for keys, so it short-circuits.

### What stays the same

Structure checks, SPEC-DELTA integrity, the date-drift signal, `check_links`
owning `related:` resolution, `check_frontmatter` owning root anchors, and the
rule that `doctor` never edits SPEC files. `index.md` is skipped — it's a catalog
doc-class, not a capability, and legitimately carries no `status:`.

## Milestones

- M1 — **Schema check shipped.** `check_specs` validates required keys + ISO date
  + closed status enum + conditional version on flat `specs/*.md`; `index.md`
  skipped. This is the whole feature — one mechanical slice.

## Tasks

See `TASKS.md`.

## Dependencies

- Reuses the frontmatter-extraction + ISO-date pattern already in
  `check_frontmatter` (copied inline, not shared — see the two-schema constraint).
- `check_links` continues to own `related:` resolution.

## Validation

- 7 scenarios in `tests/cli/specs.test.sh` (one failing case per rule, both
  branches of conditional-version, + an index-skip regression guard):
  clean draft · missing summary · non-ISO date · status outside enum · published
  without version · published with version (clean) · index.md lacking status (clean).
- End-to-end: `doctor specs` on this repo — `doc-engine.md` (draft, no version) →
  clean; `roadmap.md` (published, versioned) → clean; `index.md` → skipped.

## Deliverables

- Extended `check_specs` in `cli/spectacular` (schema layer on the flat-file loop). ✅
- 7 new scenarios in `tests/cli/specs.test.sh` (11–17). ✅
- ROADMAP ledger row mapping build b11 → target version.
- Doc note in `doctor` area reference (specs table) describing the schema.
