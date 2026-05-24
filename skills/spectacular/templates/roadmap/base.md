---
version: 2.0
updated: <DATE>
summary: "Per-version scope, phase, and exit criteria. Active line is v<X.Y.Z>."
related:
  - PRD.md
  - ARCHITECTURE.md
---

# <Project Name> — Roadmap

<!--
  ROADMAP is the per-version planning artifact. It uses a precision gradient:
  active + near-term versions are detailed (full tier); mid-term are themed;
  long-term are vision (just direction). Icebox at the end for ideas not yet
  tied to any version.

  Mode: reps (v0.7.1+; renamed from "structured" in v1.3.0). Use `spectacular roadmap grill <version>` to walk
  the slots interactively. Use `spectacular roadmap refine` to apply vibe→spec
  rewrites + autopopulate Linked requests. Use `spectacular roadmap review` to
  check the gate.

  Tiers (per version block):
    full    — Status, Phase, Outcome, Scope (in), Scope (out), Exit criteria, Linked requests
              Use for active + near-term planned (current + next 1)
    themed  — Status, Phase, Outcome, Themes, Exit criteria (directional)
              Use for mid-term: 2-3 versions out
    vision  — Status, Direction (free-text paragraph)
              Use for long-term + speculative

  Phase taxonomy (recommended 9-phase chain + 3 meta-phase aliases):
    Specific:   intent → discover → prototype → spec-refine → mvp → iterate
                → test → release-prep → release
    Meta:       discover (intent/discover/prototype) | build (spec-refine/mvp/iterate)
                | release (test/release-prep/release)

  Both styles accepted in frontmatter. Start coarse, refine to specific phase
  as work crystallizes. Skill recommends next phase; user can skip with reason.
  Skips render as `Phase: <current> (skipped: <list>)`.
  Alpha/beta/stable qualifier: `Phase: mvp (alpha)`.

  Beginner pattern: start at vision tier (one paragraph), graduate to themed
  when 2nd version exists, unlock full when first request links via
  target_version: in PLAN frontmatter. See roadmap-rules.md.

  See `references/roadmap-rules.md` for the full spec.
-->

## v<X.Y.Z> — <Active or near-term theme>

**Tier:** full
**Status:** active
**Phase:** <intent | discover | prototype | spec-refine | mvp | iterate | test | release-prep | release> <(qualifier)> <(skipped: ...)>

**Outcome:**
<One paragraph. What business or product outcome does this version move?
Names the user/business state different after this version ships. The goal
between Phase (where in the pipeline) and Scope-in (what concretely changes).>

**Scope (in):**
- <Capability shipping in this version>
- <Capability shipping in this version>

**Scope (out):**
- <Explicitly deferred to a later version>
- <Empty list valid — write: Scope (out): []>

**Exit criteria:**
- [ ] <Concrete, checkable predicate (X exists, Y passes, Z is documented)>
- [ ] <Minimum one required>

**Linked requests:**
<!-- Autopopulated by `spectacular roadmap refine`. Reads request PLAN.md
     frontmatter `target_version: <version>` and renders matching slugs here.
     Do NOT hand-edit this section — it's regenerated. -->
- <slug> (<status>)

---

## v<X.Y.Z+1> — <Mid-term theme>

**Tier:** themed
**Status:** planned
**Phase:** intent

**Outcome:**
<One paragraph — same shape as full-tier Outcome. What business or product
outcome will this version move? Required even at themed tier; without it,
themed reduces to a feature wishlist.>

**Themes:**
- <Area of work, not a specific deliverable>
- <Area of work — 2-3 items typical>

**Exit criteria:**
- <Directional sign that this version is "done enough" — can be aspirational>
- <Not required to be checkable predicates yet — that's full-tier shape>

**Linked requests:**
<!-- autopopulated -->

---

## v<X+1>.x — <Long-term direction>

**Tier:** vision
**Status:** planned
**Phase:** intent

**Direction:**
<One paragraph. Where is this heading? What problem space does it occupy?
Free text. No commitments. No dates. No specifics. Just enough to anchor
"yes, we'd consider this someday".>

---

## Icebox

<!-- Ideas not yet tied to any version. Flat list. No phase, no scope, no
     exit criteria.

     Promoting an item (the 4-step ritual; skill walks via /spectacular roadmap):
       1. Pick the item (by text or fuzzy match)
       2. Choose target version (existing or new)
       3. Choose tier (default: vision)
       4. Fill tier-appropriate slots
       5. Delete from Icebox (only after the version block is complete)

     See roadmap-rules.md § Icebox-promotion ritual for details.

     Why "Icebox": convergent dev-tool idiom (GitHub Projects, Pivotal, Linear,
     GIST's Idea Bank). Distinguishes "unbound idea" from "planned but vague"
     (which is what vision-tier version blocks are for). -->

- <Idea worth capturing but not yet tied to a version>
- <Another icebox item>

---

## Recently shipped

<!-- Optional. Last 1-3 versions as a quick-glance reference. Older entries get
     compressed into a single line or removed entirely as the roadmap grows. -->

- v<previous> — shipped <DATE>: <one-line summary>
