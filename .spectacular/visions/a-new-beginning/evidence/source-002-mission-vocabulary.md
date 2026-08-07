---
type: source-card
source: source-002
provided_as: source2
received: 2026-08-07
authority: proposal
status: ingested
scope: [vocabulary, product-model, lifecycle, positioning, truth]
---

# Source 002 — Mission-oriented vocabulary

## Thesis

Choose one memorable conceptual world for Spectacular rather than mixing flight,
trail, and landing metaphors. The strongest proposed chain is a minimal-metaphor
hybrid:

> Anchors → Capability Contracts → Missions → Journeys → Waypoints → Execution →
> Readiness → Launch → Operations → Evidence & Records

High-stakes terms such as contracts, evidence, authority, and release should stay
unambiguous even if the surrounding operational language uses a space/expedition
metaphor.

## Repository observations

Verified on `refactor/a-new-beginning` at `9ac5335`:

- **Anchors already exists** as the canonical name for stable project-wide root
  context in `ARCHITECTURE.md` and `lifecycle-contract.md`.
- The current anchor set is broader than the proposal: it includes PRD,
  PRINCIPLES, ARCHITECTURE, roadmap index, AGENTS, STACK, decisions, and config.
- Exact uses of Mission, Missions, Capability Contract, Atlas, Journey, Waypoint,
  Product Truth, and Mission Control are absent from the inspected live runtime,
  system-spec, architecture, and documentation surfaces.
- Specifications already function as convergent implementation contracts, and
  spec-derived requests record a `contract` UUID. The proposal may clarify or
  rename an existing concept rather than introduce a new one.
- Request lifecycle is currently `planned → active → review → verified →
  archived`. Release lifecycle is separate. A request may be verified before its
  release ships.
- Current implementation truth is production code plus executable tests after
  verified integration. Production behavior is observable reality but may be a
  defect, not intended truth.
- The PRD excludes becoming a project-management or multi-agent orchestration
  platform. “Mission control” therefore carries product-positioning risk.

## Vocabulary families

| Family | Character | Intake observation |
|---|---|---|
| Space/NASA | Flight Rules, Flight Plan, Maneuvers, Launch, Mission Operations | Memorable but can overstate control and deployment universality |
| Expedition | Basecamp, Terrain Map, Route, Waypoint, Field Log | Human and adaptive, but mixes with the proposed space positioning |
| Neutral | Product Charter, System Map, Mission Plan, Work Step, Evidence | Lower metaphor cost; some terms are less distinctive |
| Minimal hybrid | Anchors, Contracts, Missions, Journey, Waypoints, neutral high-stakes terms | Source recommendation; internally tensions with “choose one world” |

## Proposals, normalized by domain

| Domain | Proposal | Intake state |
|---|---|---|
| Naming system | Choose one coherent vocabulary world | promising |
| Context | Use Anchors for enduring project context | already present; scope disputed |
| Specifications | Rename/reframe specs as Capability Contracts | promising clarification |
| Spec collection | Call the contract collection Atlas | disputed metaphor |
| Requests | Rename bounded work to Missions | mixed |
| Adaptive path | Use Journey for the revisable solution path | promising, ownership unclear |
| Milestones | Use Waypoints for outcome checkpoints | promising |
| Delivery lifecycle | Separate Readiness, Launch, Operations, and Closeout | conceptually useful; scope-changing |
| Truth/history | Separate Product Truth from Evidence and Records | strong distinction; production nuance needed |
| Positioning | “Spectacular is mission control for software projects” | memorable; orchestration collision |
| Feedback | Use Telemetry or Signals for runtime learning | mixed; must preserve feedback boundaries |

## Assumptions and contradictions to resolve

1. A themed vocabulary improves recall more than it increases jargon and
   translation cost.
2. The proposal asks for one coherent world but recommends a hybrid. A rule is
   needed for where metaphor stops and neutral terminology begins.
3. Renaming mature concepts improves comprehension enough to justify migration,
   aliases, documentation churn, and existing-user relearning.
4. Mission has a stable meaning. It could mean project purpose, request, release,
   goal, or long-running operation.
5. Anchors should remain the current broad root layer or narrow to product intent,
   principles, architecture, and policy.
6. Atlas improves discovery rather than obscuring the familiar `specs/` path.
7. Journey owns adaptive execution without duplicating PLAN, TASKS, wayfinding,
   request lifecycle, or agent-native planning.
8. Waypoint represents a demoable outcome rather than merely an intermediate
   location.
9. Launch and Operations apply to libraries, local tools, research, and content
   projects that may have no deployment or live-production phase.
10. Production behavior can be called truth without implying a defect is the
    intended contract.
11. Mission Control can remain a metaphor for coherence rather than expanding
    Spectacular into orchestration, approval authority, or deployment control.
12. Telemetry does not collapse quantitative runtime signals, qualitative human
    feedback, verification, and evidence into one concept.

## Valuable ideas independent of the full vocabulary

- Use one explicit naming design rule instead of adding metaphors opportunistically.
- Prefer neutral language at irreversible, contractual, and authority boundaries.
- Keep review, release, observation, and closeout conceptually separate.
- Preserve the distinction between intended contracts, implementation truth,
  observed behavior, evidence, and durable history.
- Test names against multiple project types, not only deployed applications.

## Provisional assessment

**Strong existing foundations:** Anchors, contract semantics for specs, separate
request/release lifecycles, and code/tests as implementation truth.

**Promising comparison candidates:** Capability Contract, Journey, Waypoint, and
an explicit neutral-at-high-stakes naming rule.

**Needs stronger evidence:** Mission, Atlas, the full hybrid chain, and whether a
themed vocabulary improves cold-agent and new-user comprehension.

**Do not treat as a cosmetic rename:** Readiness/Launch/Operations/Closeout,
Product Truth including production behavior, or Mission Control positioning.
Those choices change lifecycle scope or product promise.

No product decision is recorded by this assessment.

## Concept pieces

- [PZL-011 — Coherent vocabulary system](concepts/PZL-011-coherent-vocabulary-system.md)
- [PZL-012 — Anchors as enduring context](concepts/PZL-012-anchors-as-enduring-context.md)
- [PZL-013 — Capability Contracts](concepts/PZL-013-capability-contracts.md)
- [PZL-014 — Atlas as contract collection](concepts/PZL-014-atlas-contract-collection.md)
- [PZL-015 — Missions as bounded requests](concepts/PZL-015-missions-as-requests.md)
- [PZL-016 — Journey as adaptive solution path](concepts/PZL-016-journey-adaptive-path.md)
- [PZL-017 — Waypoints as milestones](concepts/PZL-017-waypoints-as-milestones.md)
- [PZL-018 — Readiness through closeout phases](concepts/PZL-018-readiness-launch-operations-closeout.md)
- [PZL-019 — Product Truth and Records](concepts/PZL-019-product-truth-and-records.md)
- [PZL-020 — Mission Control positioning](concepts/PZL-020-mission-control-positioning.md)
- [PZL-021 — Minimal-metaphor hybrid](concepts/PZL-021-minimal-metaphor-hybrid.md)
- [PZL-022 — Telemetry and Signals](concepts/PZL-022-telemetry-and-signals.md)

## Decision packets seeded

- Should Spectacular use a metaphor system at all, and where must neutral terms win?
- Which current concepts deserve renaming versus clearer definitions?
- What exactly belongs to Anchors?
- Is Journey a useful adaptive layer or a duplicate execution abstraction?
- Are delivery and operations part of Spectacular's universal lifecycle?
- What are the distinct authorities for intended, implemented, and observed truth?
- Does Mission Control describe the product accurately without expanding its scope?
