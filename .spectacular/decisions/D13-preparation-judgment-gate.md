---
type: Decision
id: 01a029c7-e300-748e-babb-b669d9378101
title: Adopt light gate for Mission design sufficiency and slice coherence
created_by: Alex
created: "2026-08-22T16:02:40Z"
updated: "2026-08-22T16:02:40Z"
actor: Alex
actor_role: owner
ref: D13-preparation-judgment-gate
question: How strict should Mission preparation and sizing checks be prior to activation?
disposition: orchestrator-guided-sufficiency-and-coherence
rationale: >-
    To prevent runaway, under-specified, or bloated Missions without introducing bureaucratic
    refusals for minor tasks, the Orchestrator checks for design sufficiency (`sufficient | needs-evidence | needs-decision`)
    and slice coherence (`coherent | too-broad | fragmented | dependency-bound`). If a Mission is
    found too broad, the Orchestrator proposes a decomposed sequence of smaller Missions before
    presenting for activation.
alternatives:
    - strict CLI-enforced mechanical blocking on all preparation criteria
    - unstructured informal prose with no sizing evaluation
authority_basis: Owner explicitly approved Option A (Light Gate with guided decomposition) in the design interview.
---

# Adopt light gate for Mission design sufficiency and slice coherence

## Decision
- The Orchestrator evaluates `prepare.md` criteria to ensure the approach is `sufficient` and the slice is `coherent`.
- If an outcome is too large or fragmented, the Orchestrator automatically splits it into a campaign of smaller, tightly bounded Missions before activation.
