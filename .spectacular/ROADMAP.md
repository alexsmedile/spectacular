---
version: 3.6
updated: 2026-06-28
summary: "Per-version scope, phase, and exit criteria. Shipped through v1.18.1 (SPEC.md drift check). Next runway is the coherence batch (v1.19–v1.22: naming-coherence, rules-files-audit, onboarding-dedup, lifecycle-undo), then the contract-prep ladder (v1.23–v1.25) leading to the v2.0.0 major — a single breaking concern, .spectacular/ file-contract evolution. Long-term gets fuzzier on purpose. (v3.4 of this doc: stale 'contract prep ①' v1.18 section corrected — v1.18.0 actually shipped the drift check; coherence batch slotted ahead of the ladder; ledger rows b11–b15 added.)"
related:
  - PRD.md
  - ARCHITECTURE.md
  - SPEC.md
---

# Spectacular — Roadmap

Per-version planning artifact. Uses a precision gradient: active + near-term versions are detailed (`full` tier), mid-term are themed, long-term are vision (just direction). Detail for in-flight work lives in `.spectacular/requests/<slug>/`. Shipped history is in `CHANGELOG.md`.

Phase chain: `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`. See [[roadmap-overrides]] for the full spec.

> **Versioning convention:** the number scheme follows [`docs/versioning.md`](../docs/versioning.md) — strict SemVer underneath, with this roadmap as the only place an optional marketing/arc narrative lives. **MAJOR is reserved for a breaking contract change** (here: the `.spectacular/` file-format break). The planned runway below is all backward-compatible MINOR work — including the contract-prep ladder (`workspace-v2-spec` → `workspace-v2-fields` → `workspace-v2-migration`) that *pre-stages* the contract change so **v2.0.0 is a near-mechanical flip, not a big-bang major.**
>
> *Reconciliation note (2026-05-29):* this map was rewritten to match shipped reality. The earlier version listed pre-convention bands (`v0.7.x`, `v1.5.x`–`v1.8.x`) as planned even though the soft-DB substrate, query/read verbs, and the v1.0 stable surface had already shipped (now in [Recently shipped](#recently-shipped)). Unshipped work was judged on substance — not its stale label — and re-targeted: kept items got real forward minors (v1.10–v1.12), non-urgent items moved to v2.x/vision.
>
> *Reconciliation note (2026-06-07):* v1.11, v1.12, and v1.15 marked shipped. v1.13 (cross-request-links) and v1.14 (CLI debt removal) renumbered to v1.16 and v1.17 respectively — the jump from v1.12 to v1.15 was real (visual-layer + imagine-mode shipped ahead of the advisory and debt work). Contract-prep ladder cascaded accordingly: v1.16→v1.18 becomes v1.18→v1.20.

---

## Roadmap ledger

The single source of truth for `build → version` mapping. Every planned request gets one row here when slotted. The `target-version` column is the **only place a version number is written** — everything else references requests by slug or build id.

| build | slug | title | tier | target-version | status |
|-------|------|-------|------|----------------|--------|
| b3 | convention-pack-modules | Convention pack v2 — modular packs | vision | tbd | planned |
| b4 | cli-debt-removal | CLI debt removal | themed | v1.17.0 | shipped |
| b5 | cross-request-links | Cross-request awareness | themed | v1.16.0 | shipped |
| b6 | imagine-mode | Imagine mode | full | v1.15.0 | shipped |
| b7 | roadmap-ledger | Roadmap ledger | full | v1.17.0 | shipped |
| b8 | visual-layer | Visual layer | full | v1.15.0 | shipped |
| b9 | decisions-index | Decisions index mode | full | v1.17.0 | shipped |
| b10 | skill-desc-length-check | Skill description length guard | themed | tbd | planned |
| b11 | spec-audit-mode | Content-aware spec audit | themed | tbd | planned |
| b15 | naming-coherence | Naming coherence (advance/feedback/pack/next) | themed | v1.19.0 | shipped |
| b13 | rules-files-audit | Rules-file body audit + verify-trio collapse | themed | v1.20.0 | shipped |
| b14 | onboarding-dedup | Onboarding dedup + guided first-run | themed | v1.21.0 | shipped |
| b12 | lifecycle-undo | Lifecycle undo (reverse gear) | full | v1.22.0 | planned |

> **Schema:** `build` = monotonic id (immutable); `slug` = human identity; `tier` = `full` · `themed` · `vision`; `target-version` = only mutable field (one-row edit to reslot); `status` = release-level `planned · active · shipped` (distinct from request lifecycle). See [ARCHITECTURE.md — Roadmap ledger](ARCHITECTURE.md).

---

## v1.9.0 — Versioning convention

**Tier:** full
**Status:** shipped (2026-05-29)
**Phase:** release

**Outcome:**
Spectacular has a written, enforceable versioning convention — strict SemVer, a documented breaking-change trigger, a single-canonical-version-source rule, and an optional roadmap-only marketing/arc layer — so version numbers stop drifting from reality (the very drift this ROADMAP rewrite fixes) and MAJOR stays a deliberate, asked-first event.

**Shipped:**
- `docs/versioning.md` — full convention (see intro note above)
- Registered in `docs/docs.yaml`
- Fixed long-stale `cli/spectacular` version constant (was `1.8.3`, predating v1.8.4) — the exact drift the doc warns about

**Linked requests:**
- (none — shipped directly)

---

> **Runway to v2.0.0.** The blocks below are all backward-compatible MINOR work, pinned in priority order: the validation walk (v1.11 ✓), the policy engine (v1.12 ✓), the visual layer + imagine mode (v1.15 ✓), cross-request links (v1.16 ✓), then `cli-debt-removal`, `roadmap-ledger`, and the contract-prep ladder (`workspace-v2-spec` → `workspace-v2-fields` → `workspace-v2-migration`). Numbers can reorder if priority shifts, but each is a real, self-contained release with no contract-breaking risk — the break is isolated to v2.0.0.

---

## v1.10.0 — SPEC.md density refactor (current)

**Tier:** full
**Status:** shipped (2026-05-30)
**Phase:** release
**Linked request:** `spec-refactor` (verified + archived; grew to include the skill-file audit — self-describing refs + catalog.sh + legacy cleanup)

**Outcome:**
`.spectacular/SPEC.md`'s densest 1-2 capability bullets are promoted to focused `specs/<capability>/SPEC.md` files, so an agent loading capability detail reads a standalone spec instead of a cramped index line — without bloating the index or promoting prematurely.

**Themes:**
- Density audit of all SPEC.md bullets (line count + cross-ref frequency + agent-load friction)
- Promote 1-2 (per "don't promote everything"): scaffold `specs/<capability>/SPEC.md`, lift + expand, compress index entry to one-line + link, snapshot first
- `spectacular doctor specs` stays green

**Exit criteria:**
- 1-2 dense bullets promoted with full frontmatter + standalone narrative
- SPEC.md index entries compressed to one-line + link; format otherwise unchanged
- `doctor specs` passes; new spec files appear in the report
- Snapshot + CHANGELOG entry; plugin bump to v1.10.0

**Linked requests:**
- spec-refactor (planned, retarget target_version: v1.10.0)

---

## v1.11.0 — Validation walk (`/spectacular verify`)

**Tier:** full
**Status:** shipped (2026-05-30)
**Phase:** release

**Shipped:** `references/verify.md` (typed check walk — 5 kinds, VERIFY-LOG); `spectacular archive` warns on verified-without-walk; CLI redirect for `verify`; typed check taxonomy documented in `docs/commands.md`. See CHANGELOG [1.11.0].

**Linked requests:**
- verify-walk (archived)

---

## v1.12.0 — Policy engine (configurable action gates)

**Tier:** full
**Status:** shipped (2026-05-31)
**Phase:** release

**Shipped:** POLICY.md always-set; 8 work-phase hooks; 8 prefilled policies (4 block · 4 warn); `spectacular policy` verb (5 forms); `doctor policies` area; `config.yaml` `policies:` override layer; `## Understanding` PLAN slot; policy injection loop across all 9 phase ref docs. Patches: v1.12.1 (scope-down policy + Principle 10), v1.12.2 (scaffold propagation). See CHANGELOG [1.12.0–1.12.2].

**Linked requests:**
- policy-engine (archived)

---

## v1.15.0 — Visual layer + Imagine mode (ASCII rendering)

**Tier:** full
**Status:** shipped (2026-06-07)
**Phase:** release

**Outcome:**
Spectacular's read surfaces stop rendering as flat text and gain a scannable visual layer — progress bars, a roadmap render, a summary dashboard, and ASCII mockup blocks for app-UI requests — so a human (or agent) understands workspace state at a glance instead of parsing `M1 — …: 0/5` lines. **Paired in the same release: `imagine` mode** — imagination-backed planning that *renders* see-able ASCII artifacts (UI/flow/stories/arch) the human reacts to, then derives a draft PLAN. Both are the ASCII-rendering milestone and share the rendering substrate; the data already exists (`progress`, `summary`, the milestone parser) — this is a **rendering + generation layer** over it, not a new subsystem.

**Shipped:** `_ascii_bar` / `_ascii_box` / `_ascii_color_enabled` helpers; visual `progress`, `summary`, `roadmap` renders; ASCII mockup block format; `docs/visual-conventions.md`; `references/imagine.md` + `vision-rules.md`; `vision/` soft-folder substrate; `doctor vision` area; 35-assertion test suite green. See CHANGELOG [1.15.0].

**Linked requests:**
- visual-layer (verified, archived)
- imagine-mode (verified, archived)

---

## v1.16.0 — Cross-request awareness (advisory)

**Tier:** themed
**Status:** shipped (2026-06-08)
**Phase:** release

**Shipped:** `depends-on:` / `blocks:` link schema in PLAN frontmatter; inverse-link resolver (computed, never written); `spectacular links [<slug>] [--json] [--all]` verb; link section in `spectacular request <slug>`; summary link advisory; `new` relationship prompt; `doctor links` with root-aware path resolution + `depends-on:`/`blocks:` validation; `doctor memory` staleness flag (180-day). See CHANGELOG [1.16.0].

**Linked requests:**
- cross-request-links (archived)

---

## v1.17.0 — Roadmap ledger + Decisions index + CLI debt removal

**Tier:** themed
**Status:** active
**Phase:** release-prep

**Outcome:**
Three housekeeping items that sharpen the operational substrate. Roadmap ledger makes `build → version` the single source of truth (one edit to reslot a request, not ~14 refs). Decisions index mode splits a flat `DECISIONS.md` into a cheap index + per-entry files when it grows large. CLI debt removal sheds the long-deprecated `docs *` verbs and `--global` alias. The first two ship in this release; cli-debt-removal is planned for this slot.

**Shipped (2026-06-16):**
- `roadmap-ledger` — ledger table as single source of truth; `spectacular roadmap` reads from ledger; `spectacular new` stamps build ids; `doctor links` flags stray version refs
- `decisions-index` — index mode (`decisions/D<N>.md` + cheap root index); `spectacular decisions migrate`; `doctor decisions` area

**Themes (cli-debt-removal, planned):**
- Remove `docs init|export|new|review|status` verbs + the `deprecation_notice()` banner machinery
- Remove the `docs-contract` / `docs-rules` / `docs-renderer-adapters` reference docs + legacy back-compat PRD references
- Remove the `--global` alias for `--skill-scope global`
- Update `--help`, usage, tests for the removed surface; `doctor docs` stays (discovery-only)

> **Versioning note (cli-debt-removal):** treating banner-warned, long-deprecated verb removal as MINOR is a deliberate call (per [`docs/versioning.md`](../docs/versioning.md) the strict reading would be MAJOR). The justification: the removal was announced in-product since v1.2.0, `pageworks` is the documented replacement, and no *current* documented surface changes behavior.

**Linked requests:**
- roadmap-ledger (shipped)
- decisions-index (shipped)
- cli-debt-removal (planned, medium)

---

> **Reconciliation note (2026-06-27):** v1.18.0 shipped the **SPEC.md drift check** (see [Recently shipped](#recently-shipped)), not "Contract prep ①" as an earlier draft of this doc labelled it — a stale section that predated the release. Corrected here: the contract-prep ladder moved to v1.23–v1.25, and a **coherence batch** (b15/b13/b14/b12) is slotted ahead of it at v1.19–v1.22. These four are smaller, ready, and sharpen daily use; the contract-prep ladder is still all `intent`-phase, so deferring it costs nothing.

---

## v1.19.0 — Naming coherence (advance · feedback · pack · next)

**Tier:** themed
**Status:** shipped (2026-06-28)
**Phase:** release
**Linked request:** `naming-coherence` (b15, verified)

**Outcome:**
The verb surface stops carrying near-synonyms and stutters: lifecycle move-forward is `advance` (was `promote`), human feedback is `feedback` (was `feedback-loop`), convention packs are `pack` only (drop `convention-pack`), and a new read-only `spectacular next` names the single highest-priority next action. All renames ship backward-compatible (deprecation aliases), and flow docs gain one-line tier-reveal suggestions so users discover the next verb in context.

**Themes:**
- `promote` → `advance` (alias kept, deprecation notice)
- `feedback-loop` → `feedback` (hidden alias); `feedback-rules.md` doc name unchanged
- `convention-pack` → silent alias of `pack`, flagged for removal next release
- `spectacular next` — read-only, mutates nothing, prints one next action
- Tier-reveal: one suggestion max in new-request/active-request/lifecycle/archive flows, never mid-flow

**Exit criteria:**
- Both old + new verb names work; deprecation notices on the old ones
- `spectacular next` returns a sensible action on active + empty workspaces, mutates nothing
- `--help`, docs/commands.md, SKILL.md routing all updated; tests cover both paths
- CHANGELOG entry; plugin bump to v1.19.0

---

## v1.20.0 — Rules-file body audit + verify-trio collapse

**Tier:** themed
**Status:** shipped (2026-06-28)
**Phase:** release
**Linked request:** `rules-files-audit` (b13, verified)

**Outcome:**
Skill-reference doc sprawl shrinks without touching dispatch: the 6 stub-mode `<doc>-rules.md` bodies (near-identical boilerplate) are thinned to frontmatter + a single pointer (or promoted where a real body is warranted), and the three-file verify trio (`verification.md` + `verify.md` + `verify-tests.md`) collapses into one sectioned `verify.md`. Frontmatter (the engine's dispatch) stays per-file; only duplicated/empty bodies stop being maintained per-file.

**Themes:**
- Audit 6 stub bodies; classify thin-to-pointer vs write-real-body; record policy in DECISIONS.md
- Shared "stub default behavior" section in doc-index.md
- Collapse verify trio → one `verify.md`; update SKILL.md routing
- No user-facing behavior change; `doctor docs` stays clean

**Exit criteria:**
- Each thinned file keeps complete frontmatter; verb behavior identical pre/post
- `verification.md` + `verify-tests.md` no longer standalone; content sectioned in `verify.md`
- `doctor docs` clean; CHANGELOG entry; plugin bump to v1.20.0

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
**Status:** planned
**Phase:** intent
**Linked request:** `lifecycle-undo` (b12)

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

> **Contract-prep ladder (`workspace-v2-spec` → `workspace-v2-fields` → `workspace-v2-migration`).** Three non-breaking MINORs that stage the v2.0.0 file-contract change so the major becomes a near-trivial "flip the switch." Each is backward-compatible on its own: spec the design → soak the fields → stage the migration. All hang off the existing `workspace_schema:` field and the already-shipped `spectacular migrate` registry infra.

---

## v1.23.0 — Contract prep ①: v2 contract spec (doc only)

**Tier:** full
**Status:** planned
**Phase:** intent

**Outcome:**
The v2 `.spectacular/` file format is fully *designed and frozen on paper* before any code changes — shipped as the first real per-capability spec under `specs/workspace-v2/SPEC.md` plus an ARCHITECTURE.md update — so the contract is reviewed and agreed before implementation, and `workspace-v2-fields`, `workspace-v2-migration`, and the major just execute against a fixed target.

**Themes:**
- `specs/workspace-v2/SPEC.md` — the v2 file-format contract (new/changed frontmatter fields, file layout, what breaks vs. what's additive)
- ARCHITECTURE.md updated with the v2 layout + the v1→v2 delta
- DECISIONS ADR per breaking element (what breaks, why, migration path) — written now, while it's a design discussion, not a code rush
- No code change — this is the design-freeze milestone (also the first inhabitant of the now-empty `specs/`, dovetailing with v1.10's spec-promotion pattern)

**Exit criteria:**
- `specs/workspace-v2/SPEC.md` ships; `doctor specs` validates it (frontmatter + non-empty body)
- ARCHITECTURE.md reflects the v2 contract + delta
- 1 ADR per breaking change in DECISIONS.md
- CHANGELOG entry; plugin bump to the target release

**Linked requests:**
<!-- autopopulated -->

---

## v1.24.0 — Contract prep ②: v2 frontmatter fields (optional/additive)

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

## v1.25.0 — Contract prep ③: v1→v2 migration scaffold (dry-run, no-op)

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

## Icebox

Ideas worth capturing but not yet tied to any version. Promoting an item via the 4-step ritual (see `roadmap-overrides.md`): pick item → choose target version → choose tier (default vision) → fill tier-appropriate slots → delete from Icebox.

- Hook-driven automation (auto-archive on merge to main, auto-propose lifecycle transitions on CI/PR/deploy signals, auto-update SESSION.md on commit once v1.5.x ships)
- Multi-agent orchestration + agent spawn method (subagent handoff conventions, parallel execution patterns, agent contracts, conflict resolution beyond v1.7.x advisory) — deferred until a real complex multi-agent request exercises the need
- Human-contributor onboarding surface (`spectacular tour` / `spectacular explain <slug>` — cold-start context for new devs, issue-reporting templates, contributing-rules generator) — currently covered by AGENTS.md aimed at agents, not humans
- Project management surface — kanban-style view over requests by lifecycle state, possibly with WIP limits per state — flagged as a re-examination of the PRD non-goal ("not a ticketing system"); only revisit if multi-request coordination pain surfaces in real use
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

- **v1.9.0** (2026-05-29) — Versioning convention (`docs/versioning.md`): SemVer + breaking-change trigger + single-canonical-version-source rule + optional roadmap-only marketing/arc layer. Fixed stale `cli/spectacular` version constant. *(This ROADMAP rewrite reconciles the version map against the convention — see intro note.)*
- **v1.8.x** (2026-05) — Read verbs: 11 read-only CLI verbs collapsing multi-step agent workflows into single deterministic calls; skim-by-default + `--full`; universal `--status|--since|--limit|--all|--json`. *(This is the shipped form of the old "v1.6.x query verbs" theme.)*
- **v1.7.0** (2026-05) — Ideas doc-type: `idea` registered as first-class doc-type; CLI `idea new|list|promote`; status `parked|exploring|promoted`; doctor `ideas` area
- **v1.6.0** (2026-05) — Feedback-loop: prototyping-stage human-feedback acquisition; `.spectacular/feedback/` substrate; doctor `feedback` area; PRINCIPLES §9 (feedback ≠ verification ≠ benchmark)
- **v1.5.0** (2026-05) — Soft-DB substrate: memory/sessions/decisions promoted to soft-folder DBs (`index.md` + `entries/`); CLI mutators `decide|remember|session start|end`; doctor `memory` + `sessions` areas. Also snapshot-tidy (snapshots under `snapshots/<DOC>/@v<N>.md`). *(This is the shipped form of the old "v1.5.x soft-DB" theme.)*
- **v1.4.0** (2026-05) — Substrate clarity: doc-registry → doc-index, grill sub-modes (`grill-wide|grill-each|grill-loop`), rules-file frontmatter dispatch schema
- **v1.3.0** (2026-05) — Personas: opt-in `PERSONAS.md` canonical doc with per-persona grill-each; triggered by product/content kits
- **v1.2.1** (2026-05-24) — patch on the pageworks deprecation
- **v1.2.0** (2026-05-24) — Deprecate public-docs surface; extracted to standalone [pageworks](https://github.com/alexsmedile/pageworks) skill. `doctor docs` becomes discovery-only; archive-time handoff prompt; `docs *` verbs banner + scheduled for removal in v2.0.0
- **v1.1.0** (2026-05-23) — `docs export` adapters: MkDocs Material + Docusaurus renderers (later moved to pageworks)
- **v1.0.1** (2026-05-23) — `spectacular --version` + top-level usage on CLI
- **v1.0.0** (2026-05-23) — First stable release; CLI + skill + 13 registered doc-types frozen for v1.x
- **v0.7.2** (2026-05-23) — ROADMAP refinement v2: required `outcome:` slot, `icebox` renaming (was "Bucket list"), tier-aware grill/review polish
- **v0.7.1** (2026-05-23) — Structured ROADMAP: per-version blocks, precision tiers, 9-phase chain, 18-check review gate
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
