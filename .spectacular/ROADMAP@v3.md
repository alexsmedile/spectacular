---
version: 2.1
updated: 2026-05-23
summary: "Per-version scope, phase, and exit criteria. Active line is v0.7.1. Long-term gets fuzzier on purpose."
related:
  - PRD.md
  - ARCHITECTURE.md
  - SPEC.md
---

# Spectacular — Roadmap

Per-version planning artifact. Uses a precision gradient: active + near-term versions are detailed (`full` tier), mid-term are themed, long-term are vision (just direction). Detail for in-flight work lives in `.spectacular/requests/<slug>/`. Shipped history is in `CHANGELOG.md`.

Phase chain: `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`. See [[roadmap-overrides]] for the full spec.

---

## v0.7.1 — Structured ROADMAP + roadmap-richness

**Tier:** full
**Status:** active
**Phase:** mvp (skipped: discover, prototype)

**Outcome:**
Spectacular maintainers can plan, review, and ship versions through a single structured ROADMAP that captures intent (Outcome), scope (in + out), and verification (Exit criteria) consistently across every release — replacing the prior freeform shape that made cross-version comparison and "what's in scope for v1?" answerable only by reading the bullets and guessing.

**Scope (in):**
- `roadmap-overrides.md` reference doc (slot prompts, mini-refine, vibe→spec, tier-aware gate)
- `templates/roadmap/base.md` rewrite to structured per-version shape with tier examples
- `doc-registry.md` switch: ROADMAP `mode: freeform` → `mode: structured`
- Live `.spectacular/ROADMAP.md` rewritten against new shape (this file) with tier gradient
- Doctor extension: `check_workspace` info line for pre-v0.7.1 ROADMAP shape
- `--target-version <ver>` flag added to `spectacular new` (M5 dependency, shipped)
- Precision tiers: `full | themed | vision` per version block + `## Icebox` section

**Scope (out):**
- Cross-doc enforcement ("every active request must link to a version") — needs inverse-link registry first
- Auto-detection of phase from request status (mapping verified → release-prep)
- Burndown / progress visualization
- Multi-product roadmaps (one ROADMAP per product line in monorepo)
- Time-based release predictions
- `spectacular roadmap grill --icebox` CLI wiring (deferred; flow defined in roadmap-overrides.md but not shipped as a separate verb)

**Exit criteria:**
- [ ] `roadmap-overrides.md` ships with tier-aware grill + tier-aware review gate
- [ ] Template rewrite shows all 3 tiers + Icebox as examples
- [ ] Live `.spectacular/ROADMAP.md` migrated with correct tier per version; doctor passes
- [ ] `spectacular roadmap review` would exit clean on the migrated file (manual review since no grill yet)
- [ ] `--target-version` flag wired into `cmd_new` (shipped this milestone)
- [ ] CHANGELOG entry + plugin bump to v0.7.1
- [ ] 7 test files still green; new asserts cover `--target-version` + ROADMAP shape detection

**Linked requests:**
<!-- autopopulated; backfill via target_version: in request frontmatter -->
- roadmap-richness (active, target_version: 0.7.1)

---

## v0.7.x — Workflows layer

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
Spectacular projects can capture their procedural sequences (release cycle, hotfix protocol, migration runbook) as first-class registered docs that the skill walks step-by-step — reducing release-mistake rate and onboarding time for new maintainers who today have to reconstruct the procedure from CHANGELOG entries.

**Themes:**
- `spectacular workflows` — project-specific procedural sequences (release cycles, hotfix flows, migration procedures)
- One file per workflow at `.spectacular/workflows/<name>.md`
- Workflow as registered doc-type with its own grill (similar shape to PLAN.md)

**Exit criteria:**
- Workflow file format defined and registered in `doc-registry.md`
- 2-3 example workflows shipped (release-checklist, hotfix-procedure)
- Skill flow exercises a workflow end-to-end on a real release

**Linked requests:**
<!-- autopopulated -->

---

## v0.11.x — Convention pack v2 (modular)

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
Convention pack authors can compose packs from smaller building blocks (inherit `minimal` + override specific rules) instead of copying the entire `alex-default` to make minor variants — making pack maintenance practical at scale and lowering the bar for sharing project-shape opinions across teams.

**Themes:**
- Pack composition (inherit + override) — gated on real composition pain surfacing from v1 use
- Pack diff/merge for resolving conflicts between inherited packs
- Multi-pack per project (currently single-pack only)

**Exit criteria:**
- `convention-pack-modules` request transitions from planned to active when composition pain surfaces
- 2+ packs in the wild that would benefit from composition

**Linked requests:**
<!-- autopopulated -->
- convention-pack-modules (planned)

---

## v1.0.0 — Stable surface

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
External users (not just the maintainer) can adopt Spectacular knowing the v1 surface is frozen and changes will follow semver discipline — turning Spectacular from a maintainer-facing experiment into a tool consumer projects can depend on without being on the upgrade treadmill.

**Themes:**
- Freeze v0.x capabilities as the v1 surface
- README + docs reorganized around user-facing capabilities (not the dev arc)
- Pinned compatibility note: Claude Code / Codex versions supported
- Semver discipline kicks in (no more breaking changes without major bump)

**Exit criteria:**
- Every v0.x archive request still resolves cleanly via doctor
- CHANGELOG v1.0.0 entry written as "first stable release" + capability index
- README rewritten as user-facing intro
- Compatibility matrix documented in `docs/`

**Linked requests:**
<!-- autopopulated -->

---

## v2.x — Multi-workspace + nested workspaces

**Tier:** vision
**Status:** planned
**Phase:** intent

**Direction:**
Support multi-workspace setups: `.spectacular.<workspace>/` for named team workspaces alongside default `.spectacular/`, and nested workspaces (`apps/builder/.spectacular/`) for monorepos where separate teams own separate apps. Each workspace stays independent — no cross-workspace inheritance. Workspace discovery walks up from cwd to find the nearest `.spectacular/`. CLI verbs gain a `--workspace <name>` flag. Cross-workspace coordination is an explicit non-goal.

---

## v3+ — Context orchestration / Repository operating system

**Tier:** vision
**Status:** someday
**Phase:** intent

**Direction:**
Spectacular as the substrate for coordinated agent teams operating on long-running products. Smart retrieval across the full workspace structure — semantic search over PRD/PLAN/TASKS/memory with citations back to source files. Cross-project memory that travels between projects without leaking project-specific detail. Validated by ≥3 real consumer projects on v1+v2 surface before any v3 design starts. Anything not validated by real v1+v2 use first stays out.

---

## Icebox

Ideas worth capturing but not yet tied to any version. Promoting an item via the 4-step ritual (see `roadmap-overrides.md`): pick item → choose target version → choose tier (default vision) → fill tier-appropriate slots → delete from Icebox.

- Hook-driven automation (auto-update SESSION.md on commit, auto-archive on merge to main, auto-propose lifecycle transitions on CI/PR/deploy signals)
- Multi-agent orchestration (subagent handoff conventions, parallel execution patterns, agent contracts) — deferred until a real complex request exercises the need
- Schema-first request validation (parse YAML schemas declared in PLAN frontmatter, validate against contract docs)
- Multi-language convention packs (Python, Go, TypeScript variants of alex-default)
- Time-tracking integration (auto-log session duration to memory)
- `spectacular roadmap grill --icebox` CLI flow (currently defined in roadmap-overrides.md but no separate verb)
- ROADMAP burndown / progress visualization (renders exit-criteria checkbox percentages per version)
- ROADMAP-as-source-of-truth enforcement (every active request must link to a roadmap version) — needs inverse-link registry first
- Confidence rating per row (GIST/ProductBoard pattern) — overlaps tier; revisit if disagreements about tier promotion become common
- Audience field on ROADMAP (internal vs external view, per Pichler) — needed only when ROADMAP gets published publicly
- Opportunity-Solution-Tree as separate registered doc-type (Torres methodology) — heavyweight; only worth it for product-discovery-heavy teams
- ICE/RICE scoring for icebox items (GIST signature) — too prescriptive for core; convention-pack territory

---

## Recently shipped

- **v0.7.0** (2026-05-23) — CLI mutator verbs (new, promote, snapshot, archive, touch); skill orchestrates, CLI mutates
- **v0.6.2** (2026-05-23) — Workspace migrations Stage 2: registry pattern + judgment skill walk + chain validation
- **v0.6.1** (2026-05-23) — Workspace migrations Stage 1: workspace_schema + migrate verb + flat contract docs + scaffold suggestion
- **v0.6.0** (2026-05-23) — Public `docs/` as first-class surface; `docs.yaml` + page frontmatter + doctor docs area
- **v0.5.0** (2026-05-23) — SPEC.md + `specs/` replace legacy `current/`; per-capability subfolder support
- **v0.4.0** (2026-05-23) — Convention pack system end-to-end (schema + fabricator + application)
- **v0.3.x** (2026-05-22 → 2026-05-23) — Doctor substrate + smart-init kits + doc-writer engine + verification 2-of-6 rule
- **v0.2.0 → v0.1.x** (2026-05-11 → 2026-05-21) — Initial scaffold, PRD-craft flow, foundational skill structure

Full CHANGELOG: [`CHANGELOG.md`](../CHANGELOG.md)

---

## Related

- [PRD.md](PRD.md) — what Spectacular is
- [SPEC.md](SPEC.md) — what's built right now
- [ARCHITECTURE.md](ARCHITECTURE.md) — structures that v0.x+ items extend
- [PRINCIPLES.md](PRINCIPLES.md) — principles every future addition must respect
