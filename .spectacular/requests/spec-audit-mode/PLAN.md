---
status: planned
priority: medium
owner: alex
updated: 2026-06-26
build: b11
summary: "Spec-audit mode: cross-check SPEC.md bullets ↔ specs/ files ↔ archives so status answers 'is the spec current?', not just 'what's the lifecycle?'"
related:
  - PRD.md
  - ../../SPEC.md
  - ../../ROADMAP.md
---

# Plan — spec-audit-mode

## 1. Goal

Give Spectacular a real answer to "is my spec current with what's built?" — beyond the single date heuristic shipped in the drift check (build b10/b11), which only compares SPEC.md's `updated` against the newest archive. A user opening status to "review specs, files, requests" should see *where* SPEC.md and `specs/` have fallen out of sync with reality, not just a smoke-alarm date warning.

## 2. Constraints

- Bash CLI only — no new runtime deps (STACK.md).
- Doctor checks the **substrate**, content reconciliation stays with the skill + spec-sync. This request adds *detection* signals to `doctor specs`; the *fixes* remain judgment (skill walks them, never auto-edits SPEC.md).
- Must not duplicate the existing date-drift check — extend the same `check_specs` area.
- Heuristics, not semantic NLP. Slug/name matching at most; flag for human, never auto-resolve.

## Understanding

### How it works now

`doctor specs` (build b11) checks: SPEC.md present + parseable, specs/ dir present, per-capability SPEC.md frontmatter, flat contract docs, and **one** drift signal — SPEC.md `updated` older than newest `archive/*/PLAN.md`. `status` relays that warning. spec-sync reconciles content, but only fires at archive time (per-request).

### What changes

Add three content-aware audit signals to `check_specs`, each a `⚠️ judgment` finding routing to spec-sync:
1. **Orphan capability bullet** — a `**Name**` bullet in SPEC.md's Capabilities list with no matching `specs/<slug>/SPEC.md` *and* no archived request whose slug/summary mentions it. (Tunable: bullets are allowed to be spec-file-less until they outgrow one line — so this is info-level unless the bullet is long.)
2. **Orphan spec file** — a `specs/<cap>/SPEC.md` not referenced anywhere in SPEC.md's body. Dead capability doc.
3. **Stale capability spec** — a `specs/<cap>/SPEC.md` whose own `updated` predates the newest archive touching it (per the archive PLAN's `related:`).

### What stays the same

The date-drift check, spec-sync's content logic, the CLI-mutates/skill-orchestrates split, and the rule that doctor never edits SPEC.md. `status` keeps relaying — it just has more signals to relay.

## 3. Milestones

- M1 — `doctor specs` flags orphan capability bullets (SPEC.md bullet ↔ no spec file ↔ no archive mention), info/warn by bullet length.
- M2 — `doctor specs` flags orphan spec files (specs/<cap> not referenced in SPEC.md body).
- M3 — `doctor specs` flags stale per-capability specs (cap SPEC.md older than its newest related archive).
- M4 — `status` + spec-sync reference docs updated; `--json` audit summary so CI can gate on spec coverage.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Builds on the SPEC.md date-drift heuristic already shipped in `check_specs` (build b11) — shares the same area + `_date_to_epoch`. No separate request; that work landed ad-hoc.

## 6. Validation

- M1 — test: SPEC.md with a bullet for a capability that has no spec file and no archive → warning; with a matching spec file → clean.
- M2 — test: a `specs/ghost/SPEC.md` not named in SPEC.md body → warning; referenced → clean.
- M3 — test: cap spec `updated` older than a related archive → warning.
- M4 — `doctor specs --json` includes per-signal findings; status briefing relays them.

## 7. Deliverables

- Extended `check_specs` in `cli/spectacular` (3 new signals).
- New scenarios in `tests/cli/specs.test.sh`.
- Updated `doctor-areas.md` (specs table), `status.md` (signal table), `spec-sync.md` (standalone audit trigger).
- ROADMAP ledger row mapping build → version.
