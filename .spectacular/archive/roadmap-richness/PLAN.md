---
status: archived
priority: medium
owner: alex
updated: 2026-05-23
target_version: 0.7.x
summary: "Promote ROADMAP from freeform → structured grill with per-version scope-in/out + release-phase taxonomy + alpha/beta/stable axis + exit criteria"
related:
  - ../../roadmaps/index.md
  - ../../../skills/spectacular/references/doc-registry.md
  - ../../../skills/spectacular/templates/roadmap/base.md
archived: 2026-05-23
---

# Plan — Roadmap Richness

## Goal

Transform ROADMAP.md from a coarse "v1/v2/v3 themes" freeform doc into a structured release-planning artifact. Each declared version should be a complete, agent-and-human readable answer to: "what's in this release, what's deferred, what phase is it in, what makes it 'done', and where on the alpha→stable axis does it sit?"

This makes ROADMAP a **single source of truth** for scope decisions across requests + commits + CHANGELOG.

## Why

Current state (as of v0.6.1):

- ROADMAP is `mode: freeform` in `doc-registry.md` — template scaffolds and exits, no grill, no review gate
- Template ships v1/v2/v3+ section headings with status field but no enforced scope/exit-criteria/phase structure
- The live `.spectacular/ROADMAP.md` reflects this: bullets per version, but no consistent shape across versions

The cost: every time a maintainer (or agent) asks "is this in scope for v1.0?" they have to read the bullets, infer intent, and guess at exit criteria. This grows worse as more requests stack up across more versions.

A structured ROADMAP collapses the question to: look at the section, read scope-in / scope-out / exit-criteria. Every release shaped the same way.

## Scope

### In scope (v1 of this request)

- **`roadmap-overrides.md`** — new reference doc, parallel to `prd-overrides.md` / `plan-overrides.md`:
  - Slot prompts for each version section
  - Mini-refine patterns (scope-in vs scope-out overlap detection; missing exit criteria; phase/version mismatch)
  - Vibe → spec rewrite table for vague release descriptions
  - Review gate checks (every declared version has scope-in/out + exit-criteria + phase)
- **Template rewrite** — `templates/roadmap/base.md` shifts to structured shape:
  - Required per-version slots: `Status`, `Phase`, `Scope (in)`, `Scope (out)`, `Exit criteria`, `Linked requests` (autopopulated from `requests/<slug>/PLAN.md` frontmatter)
  - Release-phase taxonomy: `MVP → Testing artifacts → Atomic concepts → Spec write → Build artifacts → Versioning → Release`
  - Alpha/beta/stable axis as part of `Phase` (e.g. `Phase: build-artifacts (alpha)`)
  - Comment hints inline showing what each slot is for + good/bad examples
- **Registry update** — `doc-registry.md` `roadmap:` entry switches `mode: freeform` → `mode: structured` with `overrides: roadmap-overrides.md`
- **CLI verbs (already exist via engine)** — `spectacular roadmap grill / refine / review` activate automatically once registry is structured
- **Dogfood** — rewrite the live `.spectacular/ROADMAP.md` against the new structure as the first user. Captures all current versions (v0.6.x, v0.6.2, v1.0.0) with proper scope/exit-criteria.
- **Doctor extension (lightweight)** — `check_workspace` adds info-level reminder if ROADMAP.md exists but uses old freeform shape (no `Phase:` / `Scope:` slots). Mechanical fix: re-scaffold from new template (with snapshot-before-edit).

### Out of scope (v2+)

- Cross-doc enforcement: e.g. "every active request must be linked to a roadmap version" — needs a registry of inverse links first
- Auto-detection of phase from request status (mapping `verified` requests to "build-artifacts" phase, etc.)
- Burndown / progress visualization
- Multi-product roadmaps (one ROADMAP.md per product line in a mono-repo)
- Time-based release predictions

### Explicit anti-patterns

- Date-based commitments in slots (`Ship: 2026-Q3`) — phase + scope is the contract, not dates
- Slots that duplicate `requests/` content (don't repeat the PLAN; just link by slug)
- Free-text status when the schema offers a value (`Status: "kinda planned-ish"` → must be one of: `planned`, `active`, `shipped`, `cancelled`)

## Decisions (locked 2026-05-23 via interview)

### Release-phase taxonomy — 9-phase chain

Spectacular's recommended phase progression for any major version. Skill suggests the next phase during grill; user can skip phases with reason. Skips recorded explicitly in the `Phase:` field (e.g. `Phase: spec-refine (skipped: discover, prototype)`).

| # | Phase | What's happening | Skippable? |
|---|---|---|---|
| 1 | `intent` | PRD exists; vision is set | No — given |
| 2 | `discover` | Interview user; surface scope/risks/options | Yes (when problem is well-known) |
| 3 | `prototype` | **Optional mid-flow phase** — produce any artifact that lets a decision be validated against real tooling **or against the user** before committing. Includes data/schema drafts run through parsers, fake datasets tested against downstream scripts, mock API responses, sketch implementations, **flow previews / ASCII wireframes / interactive mocks / screenshots / video walkthroughs**, sample CLI output. The artifact isn't the deliverable; the *decision it informs* is. | Yes (often skipped for small features) |
| 4 | `spec-refine` | Update PRD/SPEC/PLAN against discovery + prototype findings; lock the contract | No (locks the contract) |
| 5 | `mvp` | First shippable version of the feature/version | Sometimes (small features ship complete) |
| 6 | `iterate` | Build out beyond MVP based on use feedback | Yes (one-shot features) |
| 7 | `test` | Verification artifacts — VERIFY.md scenarios, tests, manual QA | No (verification mandatory per 2-of-6 rule) |
| 8 | `release-prep` | Semver bump, CHANGELOG, snapshots, docs sync | No (required to ship) |
| 9 | `release` | Ship, tag, archive request | No (terminal) |

**Skill behavior:**
- Recommends next phase but allows skip via `--skip-to <phase>` or interactive answer "skip to release-prep"
- Records skipped phases in the `Phase:` field so the trace stays auditable
- Doctor doesn't flag skips; only warns on **regressions** (backward phase without `--force`)
- Recommends but never enforces — strict-mode is opt-in via convention pack (`roadmap.strict: true`)

### Precision tiers — version blocks have a precision gradient

A more complete roadmap can usually be written clearly once the user end goal, expectations, project complexity, dependencies, and steps are understood. It usually starts more granular (short-mid term) and gets fuzzier (less precise, less strict, less definitive) the longer the timeframe, concluding with someday / bucket-list ideas.

Forcing the same 6-slot precision on a v3+ vision block as on the active v0.7.1 block produces artificial rigor that obscures the truth. The roadmap should let "we don't know yet" be a first-class state.

**Three tiers per version block:**

| Tier | Required slots | Use for |
|---|---|---|
| `full` | All 6 (Status, Phase, Scope-in, Scope-out, Exit criteria, Linked requests) | Active + near-term planned versions |
| `themed` | Status, Phase, **Themes** (free-text list, not concrete capabilities), Exit criteria (directional, not checkable) | Mid-term: 2-3 versions out |
| `vision` | Status (always `planned` or `someday`), **Direction** (free-text paragraph), no scope/exit slots | Long-term + speculative; just enough to anchor the direction |

The grill picks the tier based on what the user can answer:
- "Can you name 3+ concrete capabilities?" → `full`
- "Can you name 2-3 themes?" → `themed`
- "Just directions?" → `vision`

The review gate respects tiers — `vision` blocks never fail for missing scope/exit. Gate fails only when:
- A `full` block is missing any of the 6 slots
- A `themed` block is missing themes or has fewer than 1 theme
- A `vision` block has dates or concrete commitments (those don't belong there)

**Bucket list section:** separate `## Bucket list` at the end of `ROADMAP.md`. Flat list of ideas not yet tied to any version. No phase, no scope, no exit criteria — just "things we'd consider eventually". Promoting an item to a versioned block requires `Tier: vision` minimum + a version label.

### Other decisions

- **Alpha/beta/stable** is a sub-axis qualifier of `Phase:` not a separate field. Render as `Phase: mvp (alpha)` or `Phase: release-prep (beta)`.
- **`prototype` is a real phase, not just a tagged artifact directory** — when the version is actively producing prototype work, `Phase: prototype` is set. This makes "we're still figuring it out via prototyping" visible in ROADMAP. Prototype artifacts still live in `.spectacular/requests/<slug>/artifacts/prototype/` by convention.
- **`prototype` is broader than "throwaway code"** — any artifact produced to validate a decision against real tooling **or against the user** counts. Examples that qualify:
  - **Data/schema:** YAML/JSON schema drafted and run through a parser; fake example dataset exported and tested against downstream scripts or import flows
  - **API/integration:** mock API response tried against a client library to verify the contract
  - **Algorithm/logic:** throwaway implementation sketch run against a known input set
  - **UX/flow:** flow preview (numbered step-by-step walkthrough), ASCII wireframe, Figma/screenshot mock, interactive clickable prototype, video walkthrough — anything that lets the user *experience* the proposed design before it's built
  - **CLI/output:** sample output rendered as a markdown block to show how a command's results would look
  The artifact isn't the deliverable; **the decision it informs is**. Once the decision is made, the artifact's job is done — it either gets archived under `artifacts/prototype/` or deleted. The grill should ask: "what decision will this prototype let you make?" to keep the phase honest.
- **Skipping `prototype` is the common case** — most requests don't need it. The phase exists for the times when a major shape decision benefits from real-tooling validation before being committed to spec. Don't make people feel guilty for skipping.
- **Scope-in and scope-out are required** for every declared version. Empty arrays valid (`Scope (out): []`) but the slot must be present.
- **Exit criteria are checklists** (`- [x] X shipped`, `- [ ] Y validated`). Review gate validates ≥1 criterion exists per version.
- **Linked requests autopopulated** — engine scans `requests/` + `archive/` for any PLAN.md frontmatter `target_version:` matching this version. Read-only render; user can't edit (single source of truth is the PLAN frontmatter).
- **`target_version:` is optional, prompted during `spectacular new`** — `spectacular new <slug> --target-version <ver>` sets it; if omitted the field is absent and the request doesn't render in ROADMAP's Linked-requests section until set. Doctor can suggest; never errors. Backfilled by hand during M4 dogfood.
- **`spectacular new` flag addition** — `--target-version <ver>` added to the v0.7.0 `cmd_new` as part of M2 here.
- **Iteration is real, not hidden** — `iterate` is a distinct phase, not folded into `mvp`. Lets ROADMAP show "v0.5 shipped MVP and is currently in iterate" as a meaningful state.

## Validation

- `spectacular roadmap grill` on a fresh workspace walks all slots, produces a parseable ROADMAP.md
- `spectacular roadmap review` flags a roadmap missing exit criteria → fails gate
- Live `.spectacular/ROADMAP.md` migrated to new shape; all current versions (v0.6.x, v0.6.2, v1.0.0+) shape-clean
- Doctor on workspace with old-shape ROADMAP emits info line "ROADMAP uses pre-v0.7 shape — run `spectacular roadmap refine` to migrate"
- Adding a new request with `target_version: 0.6.2` shows up in v0.6.2's Linked-requests section after `spectacular roadmap refine`
- Mini-refine pattern: declaring `Scope (in): [feature-x]` AND `Scope (out): [feature-x]` flags the overlap

## Milestones

1. **M1 — roadmap-overrides.md draft** — write slot prompts, mini-refine patterns, vibe-rewrite table, gate checks. Follows prd-overrides.md structure exactly.
2. **M2 — Template rewrite** — `templates/roadmap/base.md` shifts to structured shape. Comment hints + good/bad examples.
3. **M3 — Registry switch** — `doc-registry.md` flips `roadmap:` from `freeform` → `structured` + overrides path.
4. **M4 — Dogfood: rewrite live ROADMAP** — apply the new shape to `.spectacular/ROADMAP.md` (snapshot-before-edit). Validates the template against real content.
5. **M5 — Doctor extension** — `check_workspace` detects old-shape ROADMAP, emits info line with migration verb.
6. **M6 — Tests + VERIFY** — doctor scenario for old-shape detection; tests/cli/roadmap-overrides scenarios (if grill engine surfaces fail modes worth scripting).

## Risks

- **Over-structuring kills adoption** — too many required slots means people skip the grill. Mitigation: only 6 required slots per version (Status, Phase, Scope-in, Scope-out, Exit criteria, Linked requests-autopopulated). Optional everything else.
- **9-phase chain feels heavy for small versions** — a one-day patch version doesn't need all 9 phases. Mitigation: skill recommends but allows skip; small versions typically land at `Phase: release-prep (skipped: discover, prototype, mvp, iterate)` after a quick spec-refine.
- **Phase taxonomy might be wrong for non-software projects** — the chain assumes a build pipeline. Mitigation: kits can extend/override the taxonomy via `roadmap-overrides.md` (a content kit might use `draft → review → publish`).
- **`prototype` phase confuses "this is throwaway" with "this is shippable"** — risk that prototype artifacts leak into production. Mitigation: prototype artifacts live in `.spectacular/requests/<slug>/artifacts/prototype/` (per convention) which is documented as "not part of the deliverable"; `release-prep` review gate checks that no `Phase: prototype` exists for the version being released.
- **Linked requests render breaks on stale `target_version:` fields** — a PLAN claiming `target_version: 0.5.0` after 0.5.0 already shipped looks weird in ROADMAP. Mitigation: doctor scans PLAN frontmatter and warns on past-version targets after release.
- **Backfilling existing requests with `target_version:`** is manual work. Mitigation: dogfood includes this work; agents can pattern-match from PLAN status + ROADMAP context.
- **Skip-tracking metadata grows ugly** — `Phase: spec-refine (skipped: discover, prototype, mvp)` is verbose. Mitigation: render only the current phase + skipped list in templates; the underlying frontmatter stays structured.

## Open questions

- Does ROADMAP need its own `snapshot-on-edit: true` (registry override)? Yes, probably — version bumps + scope changes are exactly the kind of change worth snapshotting.
- How does this interact with `public-docs-advanced`'s spec→doc sync? ROADMAP isn't a spec; it's a planning artifact. Probably no interaction.
- Should `archive/<slug>/PLAN.md` `target_version:` still count toward Linked-requests rendering? Yes — historical accuracy matters; archived items shown as `(shipped)` markers.
