---
updated: 2026-08-04
related:
  - PLAN.md
  - TASKS.md
---

# Verify — vision-workflow

This request earns a standalone verification file through two axes: the CLI
and workflow are user-visible, and the commands plus `.spectacular/visions/`
layout are an external contract.

## Automated {run}

- [x] bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
- [x] bash tests/cli/vision.test.sh
- [x] bash tests/cli/lifecycle-contract.test.sh
- [x] bash tests/cli/doctor.test.sh
- [x] bash tests/cli/request-workflow.test.sh
- [x] bash tests/cli/roadmap-migrate.test.sh
- [x] bash tests/run.sh
- [x] scripts/hooks/pre-commit --check

## Contract assertions {assert}

- [x] A new Vision can exist without a request; legacy request-owned Vision remains readable and diagnosable.
- [x] Fragment reactions and whole-Vision approval are separate, explicit records; approved Vision records actor and date.
- [x] Vision derives only a draft SPC, while approved SPC remains the source of request PLAN/TASKS.
- [x] Prototype remains an owned fragment/evidence artifact, spike remains feasibility evidence, and feedback remains post-build learning.
- [x] Backend, maintenance, and already-specified work can bypass Vision without a warning or mandatory gate.

## Coherence review {judge}

- [x] `skills/spectacular/references/imagine.md`, `vision-rules.md`, `discovery-protocol.md`, `spike-rules.md`, `feedback-loop.md`, `roadmap-rules.md`, and `spec-lifecycle.md` assign one non-competing role to each concept.
- [x] `.spectacular/ARCHITECTURE.md` and `.spectacular/specs/index.md` describe the implemented workflow and preserve the legacy compatibility boundary.
- [x] User-facing documentation impact is assessed without directly editing Pageworks-owned `docs/`; any substantial follow-up has a usable handoff prompt.

## Rollback {assert}

- [x] Reverting this feature commit restores the previous request-owned Vision behavior without requiring a workspace data migration.
