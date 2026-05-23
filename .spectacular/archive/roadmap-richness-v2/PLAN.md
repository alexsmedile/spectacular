---
status: archived
priority: high
owner: alex
updated: 2026-05-23
target_version: 0.7.2
summary: "Apply roadmap-research findings: add Outcome slot, rename Icebox, date guards on themed/vision, row-cap tiered warning, scope-out warning when scope-in>3, meta-phase aliases, beginner-mode docs, icebox-promotion ritual docs"
related:
  - PRD.md
  - ../../archive/roadmap-richness/PLAN.md
  - ../../../skills/spectacular/references/roadmap-overrides.md
provenance:
  - source: web research (2026-05-23) — Pichler GO, Torres OST, Cagan SVPG, Gilad GIST, Productside, GitHub Projects, Aha! guide
    captures: "Convergent slot/field set across roadmap frameworks; long-term fuzziness handling; anti-patterns; beginner onboarding patterns. Strongest signal: Outcome/goal as #1 missing slot."
archived: 2026-05-23
---

# Plan — Roadmap Richness v2

## Goal

Apply convergent findings from roadmap-methodology research (Pichler GO, Torres OST, Cagan SVPG, Gilad GIST, Productside, beginner-tool onboarding patterns) to tighten Spectacular's structured ROADMAP. Closes the gap between v0.7.1 and what dominant frameworks recommend, without ballooning the slot set or losing the precision-tier gradient that's v0.7.1's differentiator.

## Why

v0.7.1 shipped the 3-tier model + 6 slots + 9-phase chain. Research validates structural decisions but exposes specific tightenings:

- **Goal/outcome is the #1 missing slot** — Pichler, Torres, Cagan, Gilad all converge that "what outcome does this version move?" is the most-cited gap in feature-list-shaped roadmaps. We have Scope-in (what ships) and Exit criteria (how we know it's done) but not the goal in between.
- **"Bucket list" reads consumer-y** — convergent dev-tool idiom is "Icebox" (GitHub Projects, Pivotal, Linear) or GIST's "Idea Bank".
- **Dates in long-term tiers** is Cagan's "#1 sin" — our gate forbids dates in scope/exit slots but not anywhere in themed/vision blocks.
- **Row cap (5-7)** — Cagan's "a roadmap with 40 rows is a backlog" applies; current gate doesn't push back.
- **Scope-out is universal in prose, rare in templates** — Productside flags "scope ambiguity"; our slot is good but not pressured when scope-in grows.
- **9-phase chain is on the long end** — discoverable via 3 meta-phase groupings (discover/build/release) matching Cagan's discovery-vs-delivery split.
- **Beginner-tool pattern** — GitHub/Trello/Notion start with smallest useful roadmap (1 row, 3 columns). Our 6-slot full-tier upfront is too much for a first roadmap.
- **Icebox-promotion ritual** — GIST's Idea Bank → Step-project requires explicit re-scoring + tier-up. Our doc mentions promotion but doesn't walk it.

## Scope

### In scope (v0.7.2)

**Structural changes:**
1. **Add `Outcome:` slot** between Phase and Scope-in. Required for full + themed tiers; absent for vision (Direction covers). One paragraph: "what business or product outcome does this version move?"
2. **Rename "Bucket list" → "Icebox"** across template, overrides, live ROADMAP.

**Review gate additions** (all warnings/info, not errors — preserve recommend-not-enforce stance):
3. **Date guards extended to themed/vision blocks** — currently only checked in scope/exit slots; extend to scan entire themed/vision blocks for `YYYY-MM-DD`, `Q[1-4] YYYY`, `MMM YYYY`. Warning on hit.
4. **Row-cap tiered warning** — count full-tier `## v` blocks:
   - ≤4: silent
   - 5-7: silent (sweet spot)
   - 8-10: info "consider demoting older versions to themed tier"
   - 11+: warning "full-tier count high — likely roadmap-as-backlog anti-pattern"
5. **Scope-out warning at 4+ scope-in** — when Scope-in has ≥4 items AND Scope-out is empty/`[]`, warning "consider what you're explicitly deferring."

**Phase taxonomy extension:**
6. **Meta-phase aliases in frontmatter** — `Phase:` accepts both individual values (e.g. `mvp`) AND coarser meta-phase values:
   - `discover` = intent | discover | prototype
   - `build` = spec-refine | mvp | iterate
   - `release` = test | release-prep | release
   - Gate accepts both styles. Document "start coarse, refine as work crystallizes."

**Documentation-only:**
7. **Beginner pattern in `roadmap-overrides.md`** — "start at vision tier, add themed when 2nd version exists, unlock full when first request links via target_version:". No automation; pure doc guidance.
8. **Icebox-promotion ritual in `roadmap-overrides.md`** — explicit 4-step walk (pick item → choose version → choose tier → fill slots → delete from Icebox). Skill executes on `/spectacular roadmap` invocation; no new CLI verb.

**Dogfood:**
9. Live `.spectacular/ROADMAP.md`: backfill Outcome for v0.7.x/v0.11.x/v1.0.0 (themed) + v0.7.1 (full); rename Bucket list → Icebox; verify no dates in themed/vision; confirm row count.

### Out of scope (deferred per interview)

- **Confidence rating per row** (GIST/ProductBoard) — overlaps tier
- **Audience field** (Pichler internal-vs-external) — over-engineered for solo/small-team
- **Opportunity-Solution-Tree as separate doc** (Torres) — heavyweight new doc-type
- **ICE/RICE scoring for icebox** (GIST signature) — too prescriptive; convention-pack territory
- **CLI verb for icebox promotion** (`spectacular roadmap promote <item>`) — manual via skill enough
- **`--beginner` flag** — auto-detection needs state-machine; doc-only enough
- **Row-cap as error** — stays warning
- **Scope-out as error** — stays warning

### Anti-patterns (explicit)

- Don't add a slot that overlaps an existing one — Outcome is added because it can't be inferred from Scope-in or Exit criteria.
- Don't force meta-phase aliases as exclusive — both styles coexist. Document recommendation; don't enforce.
- Don't add row-cap or scope-out hard errors — preserve recommend-not-enforce stance.

## Decisions (locked 2026-05-23 via interview)

- **Outcome slot added** — required for full + themed; absent for vision (Direction covers).
- **Icebox rename** — convergent dev-tool idiom; pure section-header change.
- **Date guards** — warnings (not errors) on themed/vision block contents.
- **Row cap** — tiered warning: silent ≤7, info 8-10, warning 11+. Never errors.
- **Scope-out push** — warning at 4+ scope-in items with empty scope-out.
- **Meta-phase aliases** — frontmatter accepts both `Phase: mvp` and `Phase: build`. Documented with "start coarse, refine as work crystallizes."
- **Beginner mode** — documentation-only. No automation.
- **Icebox-promotion ritual** — documentation-only. Skill walks 4 steps; no new CLI verb.
- **All 4 deferred items** stay deferred. Revisit only if real usage surfaces need.

## Validation

- `spectacular roadmap review` flags themed/vision blocks containing dates (warning).
- 11 full-tier blocks emits warning; 8-10 emits info; ≤7 stays silent.
- Scope-in 4+ items + empty Scope-out emits warning.
- Gate accepts both `Phase: mvp` and `Phase: build` without error.
- Live ROADMAP.md passes gate after backfill: Outcome populated for themed blocks, Icebox renamed, no dates in themed/vision.
- Doctor on workspace whose ROADMAP still says "Bucket list": info line suggesting rename (mechanical fix).
- All existing tests (7 files, 198+ asserts) pass.

## Milestones

1. **M1 — `roadmap-overrides.md` updates** — Outcome slot prompt + tier rules; beginner-mode section; icebox-promotion ritual section; extend date-pattern gate to all themed/vision contents; add gate checks for row cap + scope-out + Outcome required-by-tier
2. **M2 — Template rewrite** — `templates/roadmap/base.md` adds Outcome in full + themed examples; renames Bucket list → Icebox
3. **M3 — Live ROADMAP dogfood** — snapshot + update `.spectacular/ROADMAP.md`: add Outcome to v0.7.1 (full) + v0.7.x/v0.11.x/v1.0.0 (themed); rename Bucket list → Icebox; verify no date drift
4. **M4 — Meta-phase aliases** — extend Phase taxonomy table to declare aliases; update gate check 8 to accept both styles; document the coexist rule
5. **M5 — Doctor extension (light)** — `check_workspace` detects "Bucket list" heading in ROADMAP, suggests rename (mechanical)
6. **M6 — Tests + v0.7.2 release** — extend `doctor.test.sh`; CHANGELOG; plugin bump 0.7.2; SPEC.md update; CLAUDE.md active-requests; archive request

## Risks

- **Outcome slot feels like ceremony for tiny versions** — adds one more slot per full/themed. Mitigation: template shows one-sentence example; "Outcome: ship v0.7.2 capabilities to consumer projects" is enough.
- **Meta-phase aliases create two ways to say the same thing** — exactly the anti-pattern. Mitigation: document "start coarse, refine later"; gate accepts both but template recommends specific phases.
- **Row-cap warning at 11+ may annoy users with legitimately many active versions** — solo devs unlikely to hit; large orgs will. Mitigation: warning only, never error.
- **Backfilling Outcome on existing themed blocks** — manual work during M3 dogfood. Acceptable; that's the point of dogfood.

## Open questions

- Should Outcome also appear in TASKS.md per-milestone for symmetry? Probably not — TASKS is mechanical; outcomes belong at version level.
- Should doctor surface gate warnings via `doctor lifecycle` too? Probably yes — defer until M5 wiring.
- Does meta-phase alias need its own gate check or extend check 8? Extend check 8 — simpler.
