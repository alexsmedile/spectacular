---
status: planned
priority: medium
owner: alex
updated: 2026-07-11
build: b26
summary: "Add architectural-stance @Planning policy + grade label to PLAN frontmatter (decisions pre-locked in ideas/stance-layer.md)"
related:
  - PRD.md
  - ../../ideas/stance-layer.md
  - ../../POLICY.md
  - ../../PRINCIPLES.md
---

# Plan — stance-layer

> **All design decisions are pre-locked** via a grill session (2026-07-11). The full rationale
> lives in [[ideas/stance-layer]] Parts 2 & 3. This PLAN is the buildable extract — a builder
> should not need to re-open the decision tree.

## Goal

Add two independent, small "stance" additions: an `architectural-stance` @Planning warn policy that prompts senior-engineer reasoning only on real architectural forks, and a `grade` label on PLAN frontmatter (`prototype | mvp | standard | production`) that `status` surfaces so prototype-grade work never reads as production-grade.

## Constraints

- **No severity mechanism for `grade`.** The severity *dial* was explicitly rejected (see [[ideas/stance-layer]] Part 3). `grade` is a **label only** — no gate ever changes behavior on it. Strictness stays in declared `severity:` + `config.yaml` overrides.
- **Written legibly from birth.** This ships *after* the Part 4 legibility patch (already in POLICY.md v1.5). The new `architectural-stance` warn gets an `**Override:**` clause (L5 style) and, being a warn, **no** `⛔ BLOCKING` marker (L1 is blockers-only).
- **Parser-safe.** Any new policy block must keep its id parseable by `_policy_records` — the id is the full `### heading`, so no markers/suffixes on the heading line (learned during the legibility patch).
- **CLI mutates, skill orchestrates.** `grade` validation is a mechanical closed-enum check → belongs in `doctor lifecycle`. The architectural-stance "offer to `spectacular decide`" is agentic → skill-side.
- **Backward compatible.** Absent `grade:` resolves to `standard`; existing requests need no backfill.

## Understanding

### How it works now
- **Policies:** `.spectacular/POLICY.md` holds `### <id>` blocks under `## @<hook>` headings. `_policy_records` (cli/spectacular) parses id/principle/severity/check; `cmd_policy <hook>` renders a hook's policies; the skill injects them on phase entry. `@Planning` already carries `scope-down` + `milestones-in-build-order`.
- **`spectacular decide`:** an existing verb that writes a `decisions/D<N>` ADR entry (agentic — skill decides *when*, CLI mutates).
- **`grade`:** does not exist. PLAN frontmatter carries `status`, `priority`, `owner`, `updated`, `build`, `summary`, `related`, and optional `hold`.
- **`status` render:** `_status_*` awk paths render fleet row / card / `--json` from PLAN frontmatter + body signals (see [[specs/status]]).
- **`doctor lifecycle`:** already closed-enum-checks `status:` and `hold:` for active requests.

### What changes
1. **New policy** `architectural-stance` appended to `## @Planning` in POLICY.md (final text in [[ideas/stance-layer]] Part 2 "Final policy shape").
2. **New frontmatter field** `grade:` — recognized by `status` render paths (fleet/card/json) and validated by `doctor lifecycle`.
3. **Skill wiring** — on @Planning, when the stance's trigger fires, the skill *offers* `spectacular decide` (no new verb; reuses the existing one).
4. **Docs** — scaffold-reference/plan-rules note the optional `grade:` field; specs/index + policy-engine spec note the new policy.

### What stays the same
- No change to `_policy_records` or any severity resolution — `grade` touches none of it.
- The 4 blockers, the 2-of-6 verification rule, the archive closure gate — all unchanged.
- `standard` / absent behavior is identical to today.

## Decisions

- **architectural-stance = @Planning + warn** (not @Implementation, not block) — chose @Planning because the decision must exist before code; warn because "did you reason architecturally" is not mechanically falsifiable. Rejected: @Implementation (too late), block (unfalsifiable), both-hooks (dilutes).
- **architectural-stance fires only on a real fork** (crosses a boundary / sets a precedent / two viable structures) — chose conditional trigger over always-on because always-on manufactures fake alternatives on trivial edits. Same discipline as `ceremony-matches-uncertainty`.
- **architectural-stance offers, never requires, the ADR** — chose offer over require because require contradicts warn and can't be gated; Principle 8 (humans decide, agents propose).
- **grade = label only, no severity effect** — chose the "honesty half" over the full dial because config overrides already tune gates, and a computed `declared ± grade ± override` severity would undercut the mixed-model legibility work. Rejected: global ±1 shift, per-policy grade-severity, curated relax-list.
- **grade ladder = `prototype | mvp | standard | production`** (4 rungs) — chose four over a binary because `prototype` (throwaway) and `mvp` (smallest *shippable*) are a genuine distinction the system already implies (feedback-loop vs Principle 10). Absent ⇒ `standard`.
- **grade lives in PLAN frontmatter only** (no `config.yaml` project default) — a project-wide default was dropped as unneeded for a pure label.

## Milestones

- M1 — `architectural-stance` policy live: appears in `spectacular policy @Planning`, parses cleanly (`doctor policies` green), reads legibly (Override clause, no block marker).
- M2 — `grade` label recognized: `status` fleet row / card / `--json` render a request's grade; absent ⇒ standard.
- M3 — `doctor lifecycle` closed-enum-warns on an out-of-enum `grade` value; the four valid values pass.
- M4 — skill offers `spectacular decide` when the architectural-stance trigger fires (agentic wiring + a walkthrough in the skill reference).
- M5 — docs synced: plan-rules / scaffold-reference note `grade:`; specs/index + policy-engine spec note the new policy; CHANGELOG.

## Tasks

See `TASKS.md`.

## Dependencies

- Builds on the Part 4 legibility patch already shipped to POLICY.md (v1.5) — the new warn is authored in that style.
- No dependency on [[fleet-arc-wiring]] (b25) or any other active request. Independent.

## Validation

- M1 — run: `spectacular policy @Planning` lists `architectural-stance`; run: `spectacular doctor policies` exits 0 with the new block counted; observable: block has an `**Override:**` line and NO `⛔` marker.
- M2 — run: a request with `grade: mvp` shows `mvp` in `spectacular status <slug>` and in `status --json`; observable: a request with no `grade:` renders as standard (or blank), not an error.
- M3 — run: `spectacular doctor lifecycle` on a request with `grade: protoype` emits exactly one warning; the 4 valid values emit none.
- M4 — observable: in a @Planning walk where the change crosses a module boundary, the skill surfaces an "offer to `spectacular decide`" prompt; on a trivial edit it stays silent.
- M5 — grep: `grade` documented in plan-rules.md + scaffold-reference.md; `architectural-stance` in specs/index.md; CHANGELOG entry present.

## Deliverables

- `POLICY.md` — one new `### architectural-stance` block under `## @Planning`.
- `cli/spectacular` — `grade` recognized in `status` render paths + a closed-enum check in `doctor lifecycle`.
- Skill reference — architectural-stance offer-to-decide wiring (likely `references/active-request.md` or `plan-rules.md`).
- Docs — plan-rules.md, scaffold-reference.md, specs/index.md, policy-engine spec, CHANGELOG.
- Tests — a `doctor lifecycle` grade-enum test + a `status` grade-render test.
