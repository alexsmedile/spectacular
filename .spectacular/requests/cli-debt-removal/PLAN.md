---
status: review
priority: medium
owner: alex
updated: 2026-06-16
summary: "Remove long-deprecated docs * verbs + --global alias + docs-* refs; shrink v2 to a single breaking concern"
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../DECISIONS.md
build: b4
---

# Plan — cli-debt-removal

## 1. Goal

Remove Spectacular's accumulated, banner-warned deprecation debt — the `docs *` verbs (extracted to `pageworks` in v1.2.0), the `--global` alias, and the `docs-*` reference docs — as a clean MINOR, so the v2.0.0 major is left with a single breaking concern (the file-contract change).

## 2. Constraints

- **Classified as MINOR, deliberately.** A strict SemVer reading would call verb removal MAJOR; the justification (recorded as an ADR) is that these have shown in-product deprecation banners pointing at v2.0.0 removal since v1.2.0, `pageworks` is the documented replacement, and no *current* documented surface changes behavior. See [[versioning]] and the ADR this request writes.
- **Keep `doctor docs` (discovery-only).** The deprecation removed the *verbs*, not the discovery awareness — `doctor docs` stays as folder/manifest presence + pageworks install hint.
- **`pageworks` hint must persist.** Where `docs *` used to run, users should still be pointed at the replacement skill.
- **No file-contract changes.** That's strictly v2.0.0. This request touches only CLI surface + reference docs.

## 3. Milestones

- M1 — Audit + ADR: enumerate every deprecated surface (verbs, aliases, refs, banner machinery); write a DECISIONS ADR justifying MINOR classification + listing exactly what's removed.
- M2 — Remove verbs + banner: delete `docs init|export|new|review|status` and the `deprecation_notice()` machinery; preserve the pageworks install hint.
- M3 — Remove refs + alias: delete `docs-contract` / `docs-rules` / `docs-renderer-adapters` + legacy back-compat PRD references; remove the `--global` alias for `--skill-scope global`.
- M4 — Update surface + tests: `--help`, usage text, and the test suite reflect the trimmed surface; `doctor docs` (discovery-only) still passes.
- M5 — Ship: CHANGELOG entry (Removed section, MINOR per ADR); plugin bump to v1.17.0.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None blocking. Independent of the other runway requests. Should land *before* v2.0.0 so the major stays single-concern (the roadmap's v2 scope-out explicitly assumes this shipped).

## 6. Validation

- M1 — ADR present in DECISIONS.md with the removal list + MINOR rationale.
- M2 — `spectacular docs <anything>` no longer runs the old verb; pageworks hint still surfaces.
- M3 — Deprecated reference docs gone; `--global` removed (only `--skill-scope global` remains).
- M4 — `spectacular --help` shows no deprecated surface; test suite green; `doctor docs` discovery-only passes.
- M5 — CHANGELOG Removed section; manifests at v1.17.0.

## 7. Deliverables

- DECISIONS ADR: MINOR classification of deprecated-surface removal
- CLI with `docs *` verbs + `deprecation_notice()` machinery removed
- `--global` alias removed
- `docs-contract` / `docs-rules` / `docs-renderer-adapters` + legacy PRD refs removed
- Updated `--help`, usage, tests
- CHANGELOG [1.14.0] entry (Removed)
