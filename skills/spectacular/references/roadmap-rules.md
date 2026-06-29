---
doc-id: roadmap
mode: grill-each
location: .spectacular/ROADMAP.md
scope: project-wide
template: templates/roadmap/base.md
slots: [Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests]
snapshot-on-edit: true
summary: "Per-version scope + phase + exit criteria (per-block grill; 9-phase chain)"
status: active
---

# ROADMAP Rules — roadmap-specific rules consumed by the skill

Loaded by `grill.md` / `refine.md` / `review.md` when the active doc is `roadmap` (per doc-index).

This file declares everything roadmap-specific. The skill handles the rest.

## What this doc produces

A `.spectacular/ROADMAP.md` containing:

1. One **version block** per planned/active/shipped release, each tagged with a **precision tier** (`full | themed | vision`) that controls which slots are required
2. An optional **Icebox section** at the end for ideas not yet tied to any version

Frontmatter unchanged from prior versions (`version`, `updated`, `summary`, `related`).

## Precision tiers (the gradient)

Roadmaps start granular (short-term) and get fuzzier (long-term). Forcing the same 6-slot precision on a v3+ vision block as on the active v0.7.1 block obscures the truth. Three tiers let "we don't know yet" be a first-class state.

| Tier | Required slots | Use for |
|---|---|---|
| `full` | Status, Phase, Scope-in, Scope-out, Exit criteria, Linked requests (6) | Active + near-term planned versions (current + next 1) |
| `themed` | Status, Phase, **Themes** (list), Exit criteria (directional) | Mid-term: 2-3 versions out from current |
| `vision` | Status, **Direction** (free-text paragraph) | Long-term + speculative; no scope/exit lists |

**Tier inference during grill** (skill picks tier based on what the user can answer):
1. Ask: "Can you name 3+ concrete capabilities this version ships?" — yes → `full`
2. Else: "Can you name 2-3 themes (areas of work)?" — yes → `themed`
3. Else: "Just directions, no commitments?" — yes → `vision`

User can override the inferred tier explicitly.

**Icebox section:** a flat `## Icebox` section at the end. Items are ideas not yet tied to any version — no phase, no scope, no exit criteria. Promoting an icebox item to a versioned block requires picking a version label + `Tier: vision` minimum. See "Icebox-promotion ritual" below.

## Beginner pattern (smallest useful first roadmap)

Inspired by GitHub Projects, Trello, Notion onboarding patterns and GIST methodology. Beginner users facing the full 3-tier × 6-slot model often skip the grill entirely — too much ceremony for a first roadmap.

**The pattern**, surfaced when the user has zero version blocks (empty or scaffold-only ROADMAP):

1. **Start at `vision` tier with one block** — pick a project label (often `v1.0` or just `future`), Tier: vision, Status: planned, Phase: intent. The Direction slot is the only required content. One paragraph: "where is this heading?" That's it.
2. **Graduate to `themed` when a second version is added** — by then you have enough scope clarity to name Themes for the new version. Outcome becomes required.
3. **Unlock `full` when the first request links via `target_version:`** — having concrete linked work means you can fill Scope-in / Scope-out / Exit criteria meaningfully. Outcome stays required.

**The skill recommends this progression in prose** when grilling an empty ROADMAP, but doesn't enforce it. Users can always force `--tier full` from day one if they have the clarity.

**Why this matters:** the convergent beginner pattern across roadmap tools (GitHub Projects: 3 default columns; Trello/Notion: NNL templates with empty cards; Aha!: pick one outcome before any rows) is *progressive disclosure*. Forcing all 6 slots on first use is the same anti-pattern as showing all 25 Aha! fields to a first-time PM.

## Pre-flight (before slot loop)

1. **Resolve target** — `.spectacular/ROADMAP.md`. Scaffold from `templates/roadmap/base.md` if absent.
2. **Detect existing version blocks** — parse current ROADMAP for `## v<...>` headings. List them. Ask user which version block to grill (existing or new).
3. **For new version blocks** — ask the version label first (e.g. `v0.7.1`, `v1.0.0`, `v2.x`). Use the label as the section heading.
4. **Snapshot rule** — per registry override: ROADMAP is `snapshot-on-edit: true`. The skill snapshots `.spectacular/ROADMAP.md` via `spectacular snapshot` before any structural change.

## Per-version slot prompts

The skill uses these in the slot loop. One question per slot. Each version block is grilled independently. **Slots required depend on the version's tier** — see the tier table above.

### Tier 0 prompt (always asked first)
> What tier of precision can you commit to for this version?
>
> - `full` — you can name 3+ concrete capabilities, their scope-out, and checkable exit criteria
> - `themed` — you can name 2-3 themes / work areas with directional exit criteria
> - `vision` — just directions; no concrete commitments yet
>
> Skill infers if you describe the version; you can override.

**Slot 1 — Status**
> What's the lifecycle state of this version?
>
> Pick one: `planned` (not started) | `active` (work in progress) | `shipped` (released) | `cancelled` (abandoned).
>
> *Example:* `Status: active`

**Slot 2 — Phase**
> Where in the build pipeline is this version right now?
>
> Pick from the recommended 9-phase chain (see Phase taxonomy below):
> `intent → discover → prototype → spec-refine → mvp → iterate → test → release-prep → release`
>
> Skipped phases recorded explicitly: `Phase: spec-refine (skipped: discover, prototype)`.
>
> Optional alpha/beta/stable qualifier appended in parens: `Phase: mvp (alpha)`.
>
> *Example:* `Phase: spec-refine (skipped: prototype)` — discover happened, prototype was unnecessary, currently refining specs against discover findings.

**Slot 2.5 — Outcome** *(required for `full` + `themed` tiers; absent for `vision`)*
> What business or product outcome does this version move?
>
> One paragraph. Names the user/business state that's different after this version ships. This is the goal that scope serves — the "why" between Phase (where in the pipeline) and Scope-in (what concretely changes).
>
> Per Pichler/Torres/Cagan/Gilad: the single most-cited missing slot in feature-list-shaped roadmaps. Adding it forces goal-before-features discipline. Don't confuse with Exit criteria (how we know we're done) or Scope-in (what concretely ships).
>
> *Good:* "Reduce time-to-first-spec from 2+ hours to <30 minutes for solo devs writing their first PRD via /spectacular prd."
> *Good:* "Existing Spectacular workspaces on v0.4-shape can self-upgrade without maintainer intervention."
> *Bad:* "Improve UX" (no outcome named)
> *Bad:* "Ship the migrate verb" (that's Scope-in, not Outcome)

**Slot 3 — Scope (in)**
> What concrete capabilities ship in THIS version?
>
> Bullet list. Each item should be a noun phrase naming a deliverable. Anchor to capability bullets from SPEC.md when possible.
>
> Empty array (`Scope (in): []`) is allowed only for `Status: cancelled`. All other statuses require ≥1 item.
>
> *Example:*
> ```
> Scope (in):
>   - workspace_schema field in config.yaml
>   - spectacular migrate verb (CLI + dry-run)
>   - 2 backfilled migrations: v0.4→0.5, v0.5→0.6
> ```

**Slot 4 — Scope (out)**
> What's explicitly deferred to a later version?
>
> Bullet list. Same shape as Scope (in). Empty list valid (`Scope (out): []`) but the slot must be present — that's the contract.
>
> Lists deferrals you'd be tempted to slip in but agreed not to. Future-self will thank you.
>
> *Example:*
> ```
> Scope (out):
>   - migration registry (deferred to v0.6.2)
>   - judgment-migration skill walk (deferred to v0.6.2)
>   - --to/--from flags (deferred to v0.6.2)
> ```

**Slot 5 — Exit criteria**
> What concrete, checkable conditions mark this version "done"?
>
> Checklist. Minimum 1 item required for the gate to pass. Each item should be a verifiable predicate — "X exists", "Y passes", "Z is documented" — not aspirational vibes.
>
> *Example:*
> ```
> Exit criteria:
>   - [x] All Stage 1 milestones in workspace-migrations/TASKS.md complete
>   - [x] Live audit on the Octopus repo passes after migrate
>   - [ ] CHANGELOG entry written
> ```

**Slot 6 — Linked requests (autopopulated)**
> NOT a grilled slot — the skill fills this automatically by scanning `requests/<slug>/PLAN.md` and `archive/<slug>/PLAN.md` for `target_version:` fields matching this version label.
>
> Read-only render. To add a request to a version, set `target_version: <label>` in that request's PLAN.md frontmatter. To remove, unset.
>
> Archived requests render with a `(shipped)` marker; active/review render with their status.
>
> *Example (rendered):*
> ```
> Linked requests:
>   - workspace-migrations (shipped — Stage 1)
>   - cli-mutator-verbs (active)
>   - roadmap-richness (active)
> ```

### Themed-tier slot variant

When `Tier: themed`: Outcome (Slot 2.5) is still required (and answered the same way — one paragraph). Slots 3 (Scope-in) and 4 (Scope-out) are replaced with a single **Themes** slot. Slot 5 (Exit criteria) becomes directional.

**Themed slot 3 — Themes**
> What 2-3 areas of work does this version cover?
>
> Themes are coarser than capabilities — name a category, not a specific deliverable.
>
> *Example:*
> ```
> Themes:
>   - workspace switching UX
>   - nested workspace discovery
>   - cross-workspace coordination contract
> ```

**Themed slot 5 — Exit criteria (directional)**
> What directional signs indicate this version is "done enough"?
>
> Can be aspirational ("X works in 3 example projects"). Don't need to be checkable predicates yet — that's `full`-tier shape.
>
> *Example:*
> ```
> Exit criteria:
>   - Example monorepo project demonstrates the pattern
>   - 2-3 consumer projects validate the design
> ```

### Vision-tier slot variant

When `Tier: vision`, only Slots 1 (Status), 2 (Phase — usually `intent`), and a single **Direction** slot apply. No scope, no exit criteria.

**Vision slot 3 — Direction**
> One paragraph: where is this heading? What problem space does it occupy?
>
> Free text. No commitments. No dates. No specifics. Just enough to anchor "yes, we'd consider this someday".
>
> *Example:*
> ```
> Direction:
>   Spectacular as the substrate for coordinated agent teams operating on
>   long-running products. Semantic retrieval across the full workspace.
>   Cross-project memory that travels without leaking project-specific detail.
>   Validated by ≥3 real consumer projects on v1+v2 before any v3 design starts.
> ```

### Icebox (no tier — separate section)

The `## Icebox` section at the end of `ROADMAP.md` holds ideas not yet tied to any version. Flat list, free-text items.

> **Why "Icebox"** — convergent dev-tool idiom: GitHub Projects, Pivotal Tracker, Linear all use this name (or close cognates: "Idea Bank" in GIST, "Someday/Maybe" in GTD). Distinguishes "unbound idea" from "planned but vague" (which is what `vision`-tier version blocks are for).

**Icebox prompt** (when user invokes `spectacular roadmap grill --icebox`):
> What ideas are worth capturing but not yet tied to a version?
>
> *Example:*
> ```
> ## Icebox
> - Schema-first request validation (parse YAML schemas in PLAN frontmatter)
> - Multi-language convention packs (Python, Go, TypeScript variants of alex-default)
> - Time-tracking integration (auto-log session duration to memory)
> ```

**Icebox-promotion ritual** (the 4-step walk; skill executes on `/spectacular roadmap` invocation):

1. **Pick the item** — by exact text or fuzzy substring match. Confirm with user.
2. **Choose target version** — must reference a version label already declared in the ROADMAP, OR a new version label (the grill will scaffold the new block). Free-form labels (`v0.7.2`, `v2.x`, `v3+`) all valid.
3. **Choose tier** — defaults to `vision` (least commitment). User can promote to `themed` or `full` if they're ready to fill those slots now.
4. **Fill tier-appropriate slots** — grill walks the slots for the chosen tier. Outcome required for full/themed. Direction required for vision.
5. **Delete from Icebox** — only after the version block is fully written. Failure mid-walk leaves the icebox item intact.

The ritual is intentionally manual: promotion is a commitment, and the friction is the point.

## Phase taxonomy (the 9-phase chain + 3 meta-phase aliases)

Spectacular's recommended phase progression. Each version's `Phase:` field declares one current phase. Skipped phases recorded in parens. Skill recommends the next phase; never enforces.

The 9 specific phases group into 3 **meta-phases** that match Cagan's discovery-vs-delivery split (extended with explicit release-prep). Frontmatter accepts either style:

| Meta-phase | Specific phases | When to use which |
|---|---|---|
| `discover` | intent → discover → prototype | Coarse: "we're figuring it out" — pick when you don't know which specific sub-phase fits, or when reporting up to a non-technical audience |
| `build` | spec-refine → mvp → iterate | Coarse: "we're building" — pick early in build before MVP boundary is clear |
| `release` | test → release-prep → release | Coarse: "we're shipping" — pick when about to ship and the sub-phase doesn't matter for the audience |

**Coexist rule:** `Phase: build` and `Phase: mvp` are both valid. Start coarse when uncertain; refine to the specific phase as work crystallizes. Example progression for a single version:
- Early planning → `Phase: discover` (don't yet know if you're researching or prototyping)
- Mid planning → `Phase: prototype` (specific — actively producing artifacts to validate decisions)
- Build start → `Phase: build` (specific MVP boundary unclear)
- MVP shipped → `Phase: iterate` (specific — past MVP, refining)
- Pre-release → `Phase: release` (specific phase doesn't matter; close to shipping)

The two styles never conflict — they're aliases. Gate check 8 (Phase valid) accepts both.

### The 9 specific phases

| # | Phase | Meta | What's happening | Skippable? |
|---|---|---|---|---|
| 1 | `intent` | discover | PRD exists; vision is set | No — given |
| 2 | `discover` | discover | Interview user; surface scope/risks/options | Yes (when problem is well-known) |
| 3 | `prototype` | discover | Produce artifact that validates a decision against real tooling or against the user — schema drafts, fake data, mock APIs, ASCII wireframes, interactive mocks, sample CLI output. **Artifact isn't the deliverable; the decision it informs is.** | Yes (often skipped for small features) |
| 4 | `spec-refine` | build | Update PRD/SPEC/PLAN against discovery + prototype findings | No — locks the contract |
| 5 | `mvp` | build | First shippable version of the feature/version | Sometimes (small features ship complete) |
| 6 | `iterate` | build | Build out beyond MVP based on use feedback | Yes (one-shot features) |
| 7 | `test` | release | Verification artifacts — VERIFY.md scenarios, tests, manual QA | No — verification mandatory per 2-of-6 rule |
| 8 | `release-prep` | release | Semver bump, CHANGELOG, snapshots, docs sync | No — required to ship |
| 9 | `release` | release | Ship, tag, archive request | No — terminal |

**Alpha/beta/stable** is a sub-axis qualifier appended in parens, not a separate field: `Phase: mvp (alpha)`, `Phase: release-prep (beta)`, `Phase: release (stable)`.

**Skip notation:** `Phase: spec-refine (skipped: discover, prototype)`. Comma-separated. Multiple skips render in order.

**Grill heuristic for prototype:** ask "what decision will this prototype let you make?" — keeps the phase honest. If the user can't name the decision, the phase shouldn't be `prototype`.

## Mini-refine patterns

Applied inline by the grill after each slot answer.

| Pattern | Slot scope | Trigger | Proposed action |
|---|---|---|---|
| Scope overlap | 3 + 4 | An item appears in both Scope (in) AND Scope (out) | "`<item>` is listed both in-scope and out-of-scope. Pick one — usually in-scope wins if work has started, out-of-scope if it's been explicitly deferred." |
| Vague scope item | 3 only | Item matches the vague-scope blocklist (see below) | "`<item>` is too vague to be a contract. Replace with the specific capability — e.g. 'workspace_schema field' not 'schema improvements'." |
| Aspirational exit criterion | 5 only | Exit criterion contains aspirational verbs (`improve`, `enhance`, `optimize`) without measurable predicate | "Exit criteria are checkable. '`<item>`' is aspirational — rephrase as 'X exists' or 'Y passes test'." |
| Phase regression detected | 2 only | Version's new `Phase:` is earlier in the chain than its previous `Phase:` value (and not just adding a skip-list entry) | "Going backward from `<old-phase>` to `<new-phase>`. That's a regression — usually means scope changed or verification failed. Confirm with a one-line note (will be appended to the version block)." |
| Date in slot | 3, 4, or 5 | Slot text contains a date pattern (`YYYY-MM-DD`, `Q[1-4]`, `MMM YYYY`) | "Dates don't belong in ROADMAP slots — phase + scope is the contract, not the schedule. Remove the date or move it to a separate `Target ship:` note outside the version block." |
| Empty scope-in for non-cancelled | 3 only | `Scope (in): []` but `Status:` is not `cancelled` | "Empty in-scope only makes sense when cancelled. Either name what ships or change Status to cancelled." |
| Missing exit criterion | 5 only | Exit criteria has 0 items | "Every version needs at least one exit criterion — what makes this version 'done'? Cannot pass review gate without it." |
| Prototype without decision named | 2 only | `Phase: prototype` set, but no comment or body section names the decision being validated | "Prototype phase active. What decision will this prototype let you make? Name it explicitly so the phase has a clear exit." |
| Phase skip without recommended-first phases | 2 only | Skipping `spec-refine` (declared as 'No — locks the contract' in the chain) | "`spec-refine` is non-skippable in the recommended chain — locking the contract is what makes the rest of the work focused. If you really want to skip, use `--force-skip` and add a one-line rationale." |

## Vibe → spec rewrite tables (refine mode)

### Vague scope items → concrete capabilities

| Vibe | Spec |
|---|---|
| "improve the CLI" | "[NEEDS CLARIFICATION: name the specific CLI capability — which verb? what behavior change?]" |
| "better error messages" | "Doctor area errors include suggested-fix command + judgment/mechanical tag" |
| "performance work" | "[NEEDS CLARIFICATION: which operation? what's the current latency vs target?]" |
| "polish" | "[NEEDS CLARIFICATION: polish is vague — name 3 concrete improvements]" |
| "bug fixes" | "[NEEDS CLARIFICATION: name the bugs by issue # or symptom]" |

### Vague exit criteria → checkable predicates

| Vibe | Spec |
|---|---|
| "users love it" | "[NEEDS CLARIFICATION: how measured — NPS, feedback count, active usage?]" |
| "works well" | "[NEEDS CLARIFICATION: works well = passes which test suite or scenario?]" |
| "is stable" | "Zero panics observed over 7 days of normal use" |
| "is documented" | "README.md + docs/<feature>.md exist, link from main index, exercise via copy-paste" |
| "ready for production" | "[NEEDS CLARIFICATION: production = ships to which users? what's the rollout plan?]" |
| "MVP complete" | "User can run X command and see Y output without manual intervention" |

### Vague status → concrete status

| Vibe | Spec |
|---|---|
| "in progress" | `Status: active` |
| "almost done" | `Status: active` (with exit criteria checked off as you go) |
| "soon" | `Status: planned` (use Phase: discover or spec-refine to show what's blocking) |
| "maybe v2" | `Status: planned` + `Scope (out): [<feature>]` on the current version |
| "shipped-ish" | `Status: shipped` (with any incomplete exit criteria moved to a follow-up version) |

### Vague phase → concrete phase

| Vibe | Spec |
|---|---|
| "we're building it" | `Phase: mvp` (if first shippable) or `Phase: iterate` (if MVP shipped) |
| "still figuring out" | `Phase: discover` (interviewing) or `Phase: prototype` (validating via artifact) |
| "almost ready" | `Phase: test` or `Phase: release-prep` (be specific — test still in progress or release prep underway) |
| "released" | `Phase: release (stable)` |
| "beta" | `Phase: release (beta)` or `Phase: mvp (beta)` depending on whether iteration is still happening |

## Review gate checks (in addition to base)

**Tier-aware:** every check applies to a specific tier (or all tiers). Gate respects the version block's declared `Tier:` and skips checks that don't apply.

| # | Check | Tier | How |
|---|---|---|---|
| 4 | At least 1 version block declared | all | ROADMAP body contains ≥1 `## v<...>` heading |
| 5 | Tier is declared | all | Each version block has a `Tier:` field with value `full | themed | vision` |
| 6 | Full-tier has all 6 slots | full | Status, Phase, Scope (in), Scope (out), Exit criteria, Linked requests — each present |
| 6a | Themed-tier has Themes + directional exit | themed | Themes (≥1 item) AND Exit criteria (≥1 item) — Scope-in/out NOT required |
| 6b | Vision-tier has Direction | vision | Direction (free-text, ≥30 chars) — no scope or exit required |
| 7 | Status is valid | all | Status value is one of `planned`, `active`, `shipped`, `cancelled`, `someday` |
| 8 | Phase is valid | all | Phase value (stripping qualifier and skip-list) is one of the 9 specific phase values OR one of the 3 meta-phase aliases (`discover`, `build`, `release`). Both styles accepted; coexist per the "start coarse, refine later" rule. |
| 9 | Scope-in / Scope-out have no overlap | full | No item appears in both lists for the same version |
| 10 | Non-cancelled full-tier has non-empty Scope (in) | full | If Status ≠ cancelled AND Tier = full, Scope (in) must have ≥1 item |
| 11 | No vague-scope items in scope slots | full + themed | Slot 3 (Scope-in or Themes), Slot 4 (Scope-out, full only), Slot 5 (Exit criteria) contain none of the vague-scope blocklist |
| 12 | No date patterns anywhere in themed/vision blocks | themed + vision | Scan entire block (not just scope/exit slots) for `YYYY-MM-DD`, `Q[1-4] YYYY`, or `MMM YYYY` date strings. Cagan's "#1 sin" — long-term tiers should not carry date commitments. Warning, not error. |
| 12a | No date patterns in scope/exit slots | full | Original check, retained for full-tier blocks where dates are sometimes legitimate elsewhere (Exit criteria can name a launch date when scope is concrete) |
| 13 | Vision-tier has no concrete commitments | vision | Direction slot contains no checklists, no numbered exit criteria — just paragraph prose |
| 14 | Phase is consistent with Status | all | `Status: shipped` implies `Phase: release` (any qualifier); `Status: planned` implies Phase is one of `intent | discover | prototype | spec-refine` (or meta-phase `discover`); `Status: someday` implies `Tier: vision`; warn otherwise |
| 15 | Icebox items have no version tag | icebox | Items in `## Icebox` section don't reference specific versions (those should be vision-tier version blocks instead) |
| 16 | Outcome required for full + themed | full + themed | Outcome slot is present and ≥1 sentence (~20+ chars). For vision tier, Outcome must be absent (Direction covers). Pichler/Torres/Cagan/Gilad convergent: the #1 missing slot. |
| 17 | Full-tier row count in healthy range | all (counts full blocks) | Tiered warning based on number of `## v` blocks tagged `Tier: full`: ≤4 silent; 5-7 silent (sweet spot); 8-10 info "consider demoting older versions to themed tier"; 11+ warning "full-tier count high — likely roadmap-as-backlog anti-pattern (Cagan)". Never errors. |
| 18 | Scope-out push at 4+ scope-in items | full | When Scope-in has ≥4 items AND Scope-out is empty or `[]`: warning "consider what you're explicitly deferring — every item you add implies others you're not building (Productside: scope ambiguity)." Never errors. |

### Universal base checks still run

Placeholder check (no `<TODO>`, `<TBD>`, `???` in any **prose slot**), clarification check (no `[NEEDS CLARIFICATION: ...]` markers remaining), frontmatter integrity (frontmatter parseable as YAML).

> **Ledger `tbd` is not a placeholder.** The placeholder check scans the per-version **prose slots** (Status, Phase, Scope, Exit criteria, …). It does **not** apply to the **ledger table** at the top of ROADMAP.md, where `tbd` is a legitimate `target-version` value (see below). Never flag a `tbd` in the ledger column; only flag `<TBD>`/`<TODO>`/`???` left in a prose slot.

## The ledger (build → version)

Above the first version block, ROADMAP.md carries a **ledger table** — the single source of truth mapping each build id (`b1..bN`) to its `target-version`. Schema, columns, and the human-adds-rows / gaps-are-normal rules are canonical in `ARCHITECTURE.md` § Roadmap ledger; the two-layer model is summarized in [[specs/roadmap/SPEC]].

**`tbd` rule (behavioral):** when a build is real and prioritized but you don't yet know which release it lands in, set its `target-version` to **`tbd`** — not a guessed version, not a blank, not a `<TBD>` placeholder. `tbd` is a committed "slotted, not pinned" state; pin it to a concrete `vX.Y.Z` when the release is decided (a one-row edit). Prefer `tbd` over inventing a speculative version number — false precision on unpinned work is the exact anti-pattern the precision gradient exists to prevent.

**Never cascade-renumber the runway (the renumber anti-pattern).** When a work item ships into a version slot that *unstarted, intent-phase* future work was loosely "pinned" to (e.g. roadmap work takes v1.23.0 that the contract ladder was sketched at), **do NOT shift that future work's numbers** (v1.24 → v1.25 → …). That cascade is busywork and churns the doc on every reslot. Instead, the unstarted runway should already be `target: tbd` (ordered, not numbered) so nothing needs to move. The rule:

- **Pin a version number only when a release is actually decided/imminent.** Shipped versions keep their real numbers forever; in-flight (`active`) work may carry its target; everything further out stays `tbd`.
- **Reslotting is a one-cell edit, never a renumber.** If you ever find yourself editing 3+ version numbers because one slot moved, stop — convert the affected future items to `tbd` instead.
- Version blocks for `tbd` runway use a label header (`## Contract prep ① — … *(target: tbd)*`), not a `## vX.Y.Z` header. Only give a block a `## vX.Y.Z` header once its version is pinned.

This is why the ledger stores `build → version` (not the reverse): the build id is the stable identity; the version is a late-bound, single-source cell that defaults to `tbd`.

## Vague-scope blocklist (slots 3, 4, 5)

`improvements`, `polish`, `enhancements`, `optimizations`, `better`, `cleanup`, `refactoring`, `tweaks`, `bug fixes` (without specifics), `performance work`, `quality work`, `dx improvements`, `ux improvements`.

These describe categories of work, not what concretely changes. Push for the specific capability.

**Tokenization rule** — same as prd-rules: hyphenated compounds count as single tokens. `bug-fixes-for-edge-cases` is one token (not vague); "bug fixes for edge cases" is multi-token and hits the blocklist on `bug fixes`.

## Phase taxonomy override

A convention pack can override the phase chain by declaring `roadmap.phase-chain:` in its pack.md. Useful for non-software projects:

```yaml
# content pack might declare:
roadmap:
  phase-chain:
    - draft
    - review
    - publish
```

When a pack is active and declares its own chain, the grill uses the pack's chain instead of the default 9-phase chain. Skip notation still works the same way.

The bundled default chain is the 9-phase software chain. Most users never override.

## Related

- [[doc-index]] — catalog entry pointing to this file (`mode: grill-each` declared in frontmatter above)
- [[grill]] — consumes the slot prompts + mini-refine patterns from this file
- [[refine]] — consumes the vibe→spec tables for full refine passes
- [[review]] — consumes the gate checks from this file
- [[prd-rules]] — sibling rules file (structural model)
- [[packs-contract]] — convention packs that can override the phase chain
- [[verify]] — 2-of-6 rule that gates whether a version's `test` phase needs a separate VERIFY.md
