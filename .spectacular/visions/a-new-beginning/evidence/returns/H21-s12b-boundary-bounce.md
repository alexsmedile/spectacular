---
type: central-program-correction
source_session: H21
session: S12B
status: bounced
scope: mission-boundaries-only
foundation_contracts: unchanged
program_contract: EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.0
reason: independent-slicing-review
date: 2026-08-09
---

# H21 boundary correction — bounce and reissue

## Central disposition

The original H21 Mission boundaries are **bounced**. This correction does not reopen or supersede
the accepted foundation contracts, S12A specification topology, or the acceptance criteria that
produced them. It corrects program slicing before any implementation activation.

## P0 evidence disposition

**P0 is unproven and required as a narrow prerequisite Mission.** The current authoritative
baseline `1a86b15594ece9d4ac4e8ed8c2aae6739e6ccbb0` / tree
`66f1feb986d49dafccbc387d1cbaf12e5131a8dd` still contains both identified v1 defects:

- `cli/spectacular` reads Wayfinder readiness through direct legacy `type` at line 5293, while
  canonical records use `kind`; PZL-047 remains `disposition: pending`.
- Workspace cleanup and AFK cleanup construct and execute `git push origin --delete` for matching
  remote branches (lines 5667–5684 and 6009–6029) under the ordinary apply confirmation. The
  current CLI help and tests (`tests/cli/workspace.test.sh` and
  `tests/cli/afk-git-hygiene.test.sh`) explicitly expect remote deletion; PZL-048 remains pending.

No immutable accepted evidence proves either correction is implemented, reconciled, and tested.
P0 therefore cannot be replaced by a citation.

## Reissue boundary

H21-R1 must seek fresh owner approval for:

```text
P0 v1 safety stabilization
  → shared-scaffold Design Sufficiency gate
  → M1 semantic records + canonical workspace substrate
  → M2 durable governed Mission loop
  → M3 guided Skill + registry-generated CLI + retrieval/integrity
  → M4 frozen v1 release + final reviewed mapping inventory
  → M5 isolated capsule + v2 release readiness
```

Real project migrations and cutovers remain later, project-specific Missions. The earlier H21
Cluster-D approval does not transfer to these revised boundaries.
