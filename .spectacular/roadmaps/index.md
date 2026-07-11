---
version: 3.12
updated: 2026-07-11
summary: "Per-version scope, phase, and exit criteria. Shipped through v1.23.0 (roadmap ledger docs + index-mode pruning, b17+b18). Next up — priority order: cli-debt-removal → skill-desc-length-check → snapshot-retention → spec-audit-mode, all ahead of the contract-prep ladder (①→②→③, target tbd) → v2.0.0 major. Long-term gets fuzzier on purpose. (v3.12: contract ladder un-pinned to tbd — generic, so reslotting never forces a renumber.)"
related:
  - ../PRD.md
  - ../ARCHITECTURE.md
  - ../specs/index.md
---

# Spectacular — Roadmap

Per-version planning artifact. Uses a precision gradient: active + near-term versions are detailed (`full` tier), mid-term are themed, long-term are vision (just direction). Detail for in-flight work lives in `.spectacular/requests/<slug>/`. Shipped history is in `CHANGELOG.md`.

Phase chain: `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`. See [[roadmap-rules]] for the full spec.

> **Versioning convention:** the number scheme follows [`docs/versioning.md`](../docs/versioning.md) — strict SemVer underneath, with this roadmap as the only place an optional marketing/arc narrative lives. **MAJOR is reserved for a breaking contract change** (here: the `.spectacular/` file-format break). The planned runway below is all backward-compatible MINOR work — including the contract-prep ladder (`workspace-v2-spec` → `workspace-v2-fields` → `workspace-v2-migration`) that *pre-stages* the contract change so **v2.0.0 is a near-mechanical flip, not a big-bang major.** Shipped history lives in [`CHANGELOG.md`](../CHANGELOG.md) and (per-version prose) under [`## Shipped`](#shipped).

---

## Roadmap ledger

The single source of truth for `build → version` mapping. Every planned request gets one row here when slotted. The `target-version` column is the **only place a version number is written** — everything else references requests by slug or build id.

| build | slug | title | tier | target-version | status |
|-------|------|-------|------|----------------|--------|
| b3 | convention-pack-modules | Convention pack v2 — modular packs | vision | tbd | planned |
| b4 | cli-debt-removal | CLI debt removal | themed | v1.23.2 | shipped |
| b5 | cross-request-links | Cross-request awareness | themed | v1.16.0 | shipped |
| b6 | imagine-mode | Imagine mode | full | v1.15.0 | shipped |
| b7 | roadmap-ledger | Roadmap ledger | full | v1.17.0 | shipped |
| b8 | visual-layer | Visual layer | full | v1.15.0 | shipped |
| b9 | decisions-index | Decisions index mode | full | v1.17.0 | shipped |
| b10 | skill-desc-length-check | Skill description length guard | themed | v1.23.3 | shipped |
| b11 | spec-audit-mode | Spec frontmatter schema check (pivoted from semantic audit) | themed | tbd | review |
| b15 | naming-coherence | Naming coherence (advance/feedback/pack/next) | themed | v1.19.0 | shipped |
| b13 | rules-files-audit | Rules-file body audit + verify-trio collapse | themed | v1.20.0 | shipped |
| b14 | onboarding-dedup | Onboarding dedup + guided first-run | themed | v1.21.0 | shipped |
| b12 | lifecycle-undo | Lifecycle undo (reverse gear) | full | v1.22.0 | shipped |
| b16 | snapshot-retention | Snapshot retention + version coupling | themed | v1.24.0 | shipped |
| b17 | roadmap-contract-docs | Spec + document the roadmap ledger | themed | v1.23.0 | shipped |
| b18 | roadmap-pruning | Roadmap shipped-history pruning/scaling | themed | v1.23.0 | shipped |
| b20 | commit-discipline | Soft periodic-commit nudge for in-progress code (M1 shipped; M2 hook deferred) | themed | tbd | planned |
| b21 | builder-agent | Builder agent + build-workflow orchestrator arc (build-direction fleet) — shipped M1+M2 as `spec-builder`; M3 + fan-out deferred to ideas/builder-trace-and-fanout | themed | v1.30.0 | shipped |
| b22 | archive-closure-gate | Archive closure gate + delta-based spec-sync (fable review #1+#2) | themed | v1.28.0 | shipped |
| b23 | status-fleet-view | Deterministic `spectacular status` fleet view + enforced PLAN/TASKS structure | themed | v1.29.0 | shipped |
| b24 | cli-path-abstraction | Centralized path/variable declaration in the CLI | themed | tbd | planned |
| b25 | fleet-arc-wiring | Wire repo-explorer/code-reviewer/test-verifier into the workflow arcs (optional judgment-gated dispatch) | themed | v1.30.0 | shipped |
| b26 | stance-layer | architectural-stance @Planning policy + grade label (severity dial rejected; decisions pre-locked in ideas/stance-layer) | themed | tbd | planned |
| b27 | milestone-decomposition | Step 1.5 size-and-decompose gate + decompose-large-milestone policy — visible sub-steps inside a fat milestone (decisions pre-locked in ideas/milestone-decomposition) | themed | tbd | shipped |

> **Schema:** `build` = monotonic id (immutable); `slug` = human identity; `tier` = `full` · `themed` · `vision`; `target-version` = only mutable field (one-row edit to reslot); `status` = release-level `planned · active · shipped` (distinct from request lifecycle). See [ARCHITECTURE.md — Roadmap ledger](ARCHITECTURE.md).

### Next up — priority order (2026-06-28)

Operational substrate work ships **before** the contract-prep ladder. These are
`tbd` in the ledger (not version-pinned) — they take the next free MINOR slots in
this order as each is cut; the contract-prep ladder (`tbd`, ordered ①→②→③) follows behind them.

| Order | Build | Slug | Lifecycle | Why here |
|---|---|---|---|---|
| ✅ | b17 | `roadmap-contract-docs` | review → **shipped v1.23.0** | Ledger docs + tbd sentinel + ADR discoverability. |
| ✅ | b18 | `roadmap-pruning` | review → **shipped v1.23.0** | Index-mode pruning + `doctor roadmap`. |
| ✅ | b4 | `cli-debt-removal` | review → **shipped v1.23.2** | Verified live; swept a missed dead `templates/docs/` dir on archive. |
| ✅ | b10 | `skill-desc-length-check` | review → **shipped v1.23.3** | Verified live (53/53 tests, guard fires on commit); awk duplication reviewed + accepted (install boundary). |
| ✅ | b16 | `snapshot-retention` | planned → **shipped v1.24.0** | Version coupling + tiered retention + `_snapshots/` rename; dogfooded on this repo. |
| 1 | b11 | `spec-audit-mode` | review | Pivoted to a mechanical frontmatter schema check (semantic audit dropped); code + 7 tests shipped, awaiting archive. |
| — | b3 | `convention-pack-modules` | planned (gated) | Deferred until pack-composition pain surfaces. |

Then the runway: **contract-prep ①→②→③** (`target: tbd` — they take the next
free MINORs in order, not fixed numbers) → **v2.0.0** (the one deliberate break).
The ladder is all `intent`-phase and deliberately **not** version-pinned, so
operational work slotting ahead of it never forces a renumber. (v1.23.0 shipped
the roadmap work — b17 ledger docs + b18 index-mode pruning — ahead of the ladder.)

---

## v1.21.0 — Onboarding dedup + guided first-run

**Tier:** themed
**Status:** shipped (2026-06-28)
**Phase:** release
**Linked request:** `onboarding-dedup` (b14, verified)

**Outcome:**
First-contact gets two fixes. `onboarding.md` (existing-workspace orientation) references `status.md` for the shared read+briefing flow instead of duplicating it — one source of truth, no independent drift. And a **guided first-run** ushers a brand-new/empty workspace through new→PRD-grill→first request→`spectacular next`, one step at a time, instead of printing an empty briefing or dumping the verb surface.

**Themes:**
- `onboarding.md` references status.md for the shared spine; keeps only onboarding-specific deltas
- Guided first-run flow for empty workspaces (skill ushers; CLI entry via `init --walk` or empty-detect)
- Distinction preserved: onboarding = existing workspace; guided first-run = fresh/empty

**Exit criteria:**
- onboarding.md no longer restates the read sequence verbatim; warm-workspace status unchanged
- Empty workspace ushers one step at a time; no verb menu dumped
- CHANGELOG entry; plugin bump to v1.21.0

---

## v1.22.0 — Lifecycle undo (reverse gear)

**Tier:** full
**Status:** shipped (2026-06-28)
**Phase:** release
**Linked request:** `lifecycle-undo` (b12, verified)

**Outcome:**
Spectacular gains a reverse gear: `spectacular undo` reverts the last mutation — a lifecycle status transition, an archive (dir move + link rewrites — the hard part), or an idea promote — using a gitignored `.spectacular/.last-mutation` breadcrumb. Guardrails + `--dry-run` make it safe to preview before applying.

**Themes:**
- `.last-mutation` breadcrumb (gitignored) written by mutating verbs
- Undo status transition, archive (move back + rewrite links), idea promote
- Guardrails (only the last mutation, clean-tree check) + `--dry-run`
- **Has open questions — grill before building.**

**Exit criteria:**
- `undo` reverts each supported mutation type to the prior state
- `--dry-run` previews accurately; refuses when the breadcrumb is stale/absent
- Tests cover each undo path + the refuse cases
- CHANGELOG entry; plugin bump to v1.22.0

---

## v1.23.0 — Roadmap ledger docs + index-mode pruning

**Tier:** themed
**Status:** shipped (2026-06-29)
**Phase:** release
**Linked requests:** `roadmap-contract-docs` (b17) + `roadmap-pruning` (b18)

**Outcome:**
The roadmap's own build→version model is now documented and the file is kept lean as history grows. b17 specced the ledger (build ids, `target-version` single-source, the `tbd` sentinel, ledger-status-vs-request-lifecycle) in `specs/roadmap.md` + user docs, and made ADRs discoverable (DECISIONS.md is the home; store-worthy routing table). b18 added `spectacular roadmap migrate` + a `doctor roadmap` area — index mode that moves old shipped prose into `roadmaps/v*.md` behind a `## Shipped` index, keeping the newest 3 inline. Dogfooded here: ROADMAP.md 528 → 410 lines.

**Shipped:**
- `spectacular roadmap migrate [--dry-run] [--keep N]` + `doctor roadmap` area (orphan/stale/prune-nudge)
- `specs/roadmap.md` ledger + index-mode sections; `docs/versioning.md` ledger walkthrough; `last_build:` documented
- `tbd` sentinel documented; placeholder check scoped to prose slots
- ADR store-worthy routing table in `decisions-rules.md`; doc-index + SKILL triggers grep-match "ADR"
- `tests/cli/roadmap-migrate.test.sh` (22 assertions)

---

> **Contract-prep ladder (`workspace-v2-spec` → `workspace-v2-fields` → `workspace-v2-migration`).** Three non-breaking MINORs that stage the v2.0.0 file-contract change so the major becomes a near-trivial "flip the switch." Each is backward-compatible on its own: spec the design → soak the fields → stage the migration. All hang off the existing `workspace_schema:` field and the already-shipped `spectacular migrate` registry infra.

---

## Contract prep ① — v2 contract spec (doc only)  *(target: tbd)*

**Tier:** full
**Status:** planned
**Phase:** intent

**Outcome:**
The v2 `.spectacular/` file format is fully *designed and frozen on paper* before any code changes — shipped as the first real per-capability spec under `specs/workspace-v2.md` plus an ARCHITECTURE.md update — so the contract is reviewed and agreed before implementation, and `workspace-v2-fields`, `workspace-v2-migration`, and the major just execute against a fixed target.

**Themes:**
- `specs/workspace-v2.md` — the v2 file-format contract (new/changed frontmatter fields, file layout, what breaks vs. what's additive)
- ARCHITECTURE.md updated with the v2 layout + the v1→v2 delta
- DECISIONS ADR per breaking element (what breaks, why, migration path) — written now, while it's a design discussion, not a code rush
- No code change — this is the design-freeze milestone (also the first inhabitant of the now-empty `specs/`, dovetailing with v1.10's spec-promotion pattern)

**Exit criteria:**
- `specs/workspace-v2.md` ships; `doctor specs` validates it (frontmatter + non-empty body)
- ARCHITECTURE.md reflects the v2 contract + delta
- 1 ADR per breaking change in DECISIONS.md
- CHANGELOG entry; plugin bump to the target release

**Linked requests:**
<!-- autopopulated -->

---

## Contract prep ② — v2 frontmatter fields (optional/additive)  *(target: tbd)*

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
The new v2 frontmatter fields land as **optional, additive** in v1.19 — old workspaces keep reading and validating fine, new fields are written-when-present — so the schema soaks in real use (this repo's own `.spectacular/`) before v2.0.0 makes them load-bearing. Deprecation-in-reverse: introduce soft, harden later.

**Themes:**
- Add the v2 fields (from the `workspace-v2-spec` contract) to scaffolders + frontmatter helpers as optional
- Doctor recognizes them when present, never *requires* them (no warning on absence)
- Dogfood: this repo's workspace adopts the optional fields and runs on them through the rest of the runway
- `workspace_schema:` stays at v1 (fields are additive, not a schema break yet)

**Exit criteria:**
- New fields writable + readable; absence is silent (backward-compatible)
- Doctor validates shape when present; tests cover both old + new shape
- This repo's `.spectacular/` carries the optional fields and stays green
- CHANGELOG entry; plugin bump to the target release

**Linked requests:**
<!-- autopopulated -->

---

## Contract prep ③ — v1→v2 migration scaffold (dry-run, no-op)  *(target: tbd)*

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
The v1→v2 migration exists and is testable *before* the major — a registry entry under the already-shipped `spectacular migrate` infra whose `--dry-run` reports exactly what v2.0.0 will change, while the live apply is still effectively a no-op (fields already soaking from `workspace-v2-fields`). So when v2.0.0 flips `workspace_schema`, the migration that runs is one that's already been exercised.

**Themes:**
- `references/migrations/v1-to-v2.md` registry entry (frontmatter contract: `id`, `from`, `to`, `mechanical`, `reversible`, `apply-fn`, `affects`)
- `spectacular migrate --dry-run --to v2` reports the planned delta accurately
- `doctor` recognizes the v2 shape as a valid target (alongside v1)
- Judgment-walk path defined for any non-mechanical part (snapshot-before-edit + y/n/q)

**Exit criteria:**
- Migration registry entry ships; `--dry-run` output matches the `workspace-v2-spec` delta
- Apply path tested on a v1 fixture → produces valid v2 shape
- Doctor accepts both v1 and v2 `workspace_schema`
- CHANGELOG entry; plugin bump to the target release

**Linked requests:**
<!-- autopopulated -->

---

---

## v2.0.0 — The major: file-contract evolution

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
Spectacular evolves the `.spectacular/` file-format contract in one deliberate, asked-first major — the **only** breaking release on the map. By the time it lands, the change is nearly mechanical: the design was frozen in `workspace-v2-spec`, the new fields have been soaking since `workspace-v2-fields`, and the v1→v2 migration has been dry-run-tested since `workspace-v2-migration`. v2.0.0 is the **flip** — make the new fields load-bearing, remove the old layout, bump `workspace_schema`.

Per [`docs/versioning.md`](../docs/versioning.md), MAJOR is reserved for exactly this kind of contract break — and the agent confirms the target before bumping into it. (CLI deprecation debt is *not* here — it ships earlier as the v1.17.0 MINOR, deliberately keeping this major to a single breaking concern.)

**Themes (the flip — everything below was staged by the contract-prep ladder):**
- Make the v2 frontmatter fields **required / load-bearing** (additive + soaking since `workspace-v2-fields`)
- Remove the v1 file layout the new format supersedes (the actual break — old unmigrated workspaces stop validating)
- `workspace_schema:` bump to v2; doctor's v2-shape recognition (added in `workspace-v2-migration`) becomes the default expectation
- Promote the v1→v2 migration (scaffolded in `workspace-v2-migration`) from dry-run to live apply; exercise it on this repo's own `.spectacular/`

**Scope (out) — explicitly NOT in this major (separate future lines):**
- CLI debt removal → **already shipped in v1.17.0**
- Multi/nested workspaces (vision, below)
- Context-orchestration / semantic retrieval (vision, below)
- Convention pack v2 / modular packs (gated, below)

**Exit criteria:**
- `workspace-v2-spec`, `workspace-v2-fields`, and `workspace-v2-migration` are all shipped (hard dependency — v2.0.0 does not start until the ladder is complete)
- New fields required; v1 layout removed; `workspace_schema` = v2; tests updated for the v2-only shape
- v1→v2 migration flipped to live + exercised on this repo's `.spectacular/` with a clean doctor afterward
- Optional `-rc.N` soak per the versioning ladder if the migration warrants real-world bake time
- CHANGELOG `### Breaking` section; plugin bump to v2.0.0; agent-confirmed target

**Depends on:** `workspace-v2-spec` (spec) → `workspace-v2-fields` (optional fields) → `workspace-v2-migration` (migration scaffold). This block cannot start until all three ship.

**Linked requests:**
<!-- autopopulated — backfill when v2-planning request is cut -->

---

> **Beyond the major.** Everything below is post-v2.0.0 and stays at vision tier — direction, not commitment. None of it is bundled into the v2 major (it's explicitly scoped out above). Numbers like `v2.x` mean "somewhere in the v2 line after the major," not a pinned target.

---

## Workflows layer

**Tier:** vision
**Status:** planned (retargeted off stale v0.7.x — never shipped; non-urgent)
**Phase:** intent

**Direction:**
Projects capture procedural sequences (release cycle, hotfix protocol, migration runbook) as first-class registered docs the skill walks step-by-step — one file per workflow at `.spectacular/workflows/<name>.md`, workflow as a registered doc-type with its own grill (PLAN-shaped). Reduces release-mistake rate and onboarding time. Not needed for the v2 goal; lands when a real project feels the procedural-drift pain. (Originally deferred "to v2" in DECISIONS 2026-05-11 — that deferral predates the convention; re-homed here as vision, not pinned to the major.)

---

## Convention pack v2 — modular packs

**Tier:** vision
**Status:** planned (gated — retargeted off stale v0.11.x)
**Phase:** intent
**Linked request:** `convention-pack-modules` (planned, priority low, gated)

**Direction:**
Compose packs from smaller building blocks (inherit `minimal` + override specific rules) instead of forking the whole `alex-default`. Pack diff/merge; multi-pack per project. **Gated:** stays planned until real composition pain surfaces from v1 pack use (2+ packs in the wild that would benefit). Not part of the v2 major — it's a pack-system evolution, orthogonal to the CLI/contract cleanup.

---

## v2.x — Multi-workspace + nested workspaces

**Tier:** vision
**Status:** planned
**Phase:** intent

**Direction:**
Support multi-workspace setups: `.spectacular.<workspace>/` for named team workspaces alongside default `.spectacular/`, and nested workspaces (`apps/builder/.spectacular/`) for monorepos where separate teams own separate apps. Each workspace stays independent — no cross-workspace inheritance. Workspace discovery walks up from cwd to find the nearest `.spectacular/`. CLI verbs gain a `--workspace <name>` flag. Cross-workspace coordination is an explicit non-goal. *(Was labeled "the v2.x line" when v2 had no defined major; now the v2.0.0 major is CLI+contract cleanup, and this is a follow-on within the v2 line.)*

---

## v3+ — Context orchestration / Repository operating system

**Tier:** vision
**Status:** someday
**Phase:** intent

**Direction:**
Spectacular as the substrate for coordinated agent teams operating on long-running products. Smart retrieval across the full workspace structure — semantic search over PRD/PLAN/TASKS/memory with citations back to source files. Cross-project memory that travels between projects without leaking project-specific detail. Validated by ≥3 real consumer projects on v1+v2 surface before any v3 design starts. Anything not validated by real v1+v2 use first stays out.

---

## Shipped

> Older shipped versions — full prose moved to per-version files (`roadmaps/v*.md`); the most recent stay inline above. Facts also live in `CHANGELOG.md`.

- v1.9.0 → roadmaps/v1.9.0.md
- v1.10.0 → roadmaps/v1.10.0.md
- v1.11.0 → roadmaps/v1.11.0.md
- v1.12.0 → roadmaps/v1.12.0.md
- v1.15.0 → roadmaps/v1.15.0.md
- v1.16.0 → roadmaps/v1.16.0.md
- v1.17.0 → roadmaps/v1.17.0.md
- v1.19.0 → roadmaps/v1.19.0.md
- v1.20.0 → roadmaps/v1.20.0.md

## Icebox

Ideas worth capturing but not yet tied to any version. Promoting an item via the 4-step ritual (see `roadmap-rules.md`): pick item → choose target version → choose tier (default vision) → fill tier-appropriate slots → delete from Icebox.

- Hook-driven automation (auto-archive on merge to main, auto-propose lifecycle transitions on CI/PR/deploy signals, auto-update SESSION.md on commit once v1.5.x ships)
- Multi-agent orchestration + agent spawn method (subagent handoff conventions, parallel execution patterns, agent contracts, conflict resolution beyond v1.7.x advisory) — deferred until a real complex multi-agent request exercises the need
- Human-contributor onboarding surface (`spectacular tour` / `spectacular explain <slug>` — cold-start context for new devs, issue-reporting templates, contributing-rules generator) — currently covered by AGENTS.md aimed at agents, not humans
- Project management surface — kanban-style view over requests by lifecycle state, possibly with WIP limits per state — flagged as a re-examination of the PRD non-goal ("not a ticketing system"); only revisit if multi-request coordination pain surfaces in real use
- Schema-first request validation (parse YAML schemas declared in PLAN frontmatter, validate against contract docs)
- Multi-language convention packs (Python, Go, TypeScript variants of alex-default)
- Time-tracking integration (auto-log session duration to memory)
- `spectacular roadmap grill --icebox` CLI flow (currently defined in roadmap-rules.md but no separate verb)
- ROADMAP burndown / progress visualization (renders exit-criteria checkbox percentages per version)
- ROADMAP-as-source-of-truth enforcement (every active request must link to a roadmap version) — needs inverse-link registry first
- Confidence rating per row (GIST/ProductBoard pattern) — overlaps tier; revisit if disagreements about tier promotion become common
- Audience field on ROADMAP (internal vs external view, per Pichler) — needed only when ROADMAP gets published publicly
- Opportunity-Solution-Tree as separate registered doc-type (Torres methodology) — heavyweight; only worth it for product-discovery-heavy teams
- ICE/RICE scoring for icebox items (GIST signature) — too prescriptive for core; convention-pack territory

---

## Related

- [PRD.md](PRD.md) — what Spectacular is
- [SPEC.md](specs/index.md) — what's built right now
- [ARCHITECTURE.md](ARCHITECTURE.md) — structures that v0.x+ items extend
- [PRINCIPLES.md](PRINCIPLES.md) — principles every future addition must respect
