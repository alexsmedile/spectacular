---
type: handoff-return-review
handoff: H05
status: accepted-with-bounded-ingestion
baseline: 7a85469618d984f35755a3b5db0684bfa3620c05
reviewed: 2026-08-08
source_promoted_as: source-015
authority: orchestration-disposition
---

# H05 return — competing-skills study

## Disposition

**Accepted with bounded ingestion.** The return satisfied its read-only scope, named its baseline,
separated facts from recommendations, supplied primary repository evidence, reported no owner
decisions, and identified no reversal-grade finding. Its recommendations remain advisory.

The orchestration review promoted the comparative study as Source 015, added two independently
decidable concepts, and merged repeated findings as provenance into existing cards. It did **not**
approve Spectacular's product identity, an execution companion, a reconciliation schema, or a
feasibility spike.

## Evidence checked

Primary sources were sampled again on 2026-08-08:

- [GSD gate taxonomy](https://raw.githubusercontent.com/open-gsd/gsd-core/next/gsd-core/references/gates.md)
  defines pre-flight, revision, escalation, and abort gates, including bounded revision.
- [GSD phase discussion](https://raw.githubusercontent.com/open-gsd/gsd-core/next/commands/gsd/discuss-phase.md)
  writes a phase context artifact for downstream research and planning and routes through
  on-demand workflow loading.
- [Superpowers subagent-driven development](https://raw.githubusercontent.com/obra/superpowers/main/skills/subagent-driven-development/SKILL.md)
  uses fresh task workers, task-local ledgers, review packages, bounded repair, and final review.
- [sdd-skill instructions](https://raw.githubusercontent.com/SpillwaveSolutions/sdd-skill/main/skill.md)
  prescribe artifact-presence status and fixed progress stages, while
  [`skill.json`](https://raw.githubusercontent.com/SpillwaveSolutions/sdd-skill/main/skill.json)
  reports version `1.0.0` against the instruction file's `2.1.0` claim.

These checks establish that the reported mechanisms and drift examples exist. They do not establish
comparative adoption, user outcomes, or the correct future boundary for Spectacular.

## Candidate-finding dispositions

| Candidate finding | Disposition | Destination |
|---|---|---|
| Lifecycle promotion needs a machine-checkable evidence reconciliation edge | ingest as a distinct pending proposal | PZL-172 |
| A run-local ledger must remain subordinate to canonical Mission state | ingest as a distinct pending proposal | PZL-173 |
| Fresh-context delegation is a host-runtime technique; portable value is the handoff contract | merge as corroborating provenance | PZL-069, PZL-105, PZL-119 |
| Verification repair needs bounded retry and escalation semantics | merge as implementation-pattern provenance | PZL-127 |
| Status must be derived from canonical state rather than manually shadowed | merge as drift evidence | PZL-129 |

The broader conclusion that Spectacular should remain a project control plane is supporting evidence
for S01 and S07, not a decision. The proposed execution/review companion remains needs-evidence.

## Gaps retained

- No target demonstrates a contract-delta-to-evidence reconciliation mechanism equivalent to the
  proposed Spectacular loop.
- No comparative user research or outcome measurement was supplied.
- A minimum reconciliation envelope is not yet defined.
- The value of a standalone implementation-runner companion is unproven.
- GitHub Spec-Kit itself was outside this study.

## Orchestration return

```yaml
return:
  handoff_id: H05
  status: complete
  baseline: 7a85469618d984f35755a3b5db0684bfa3620c05
  result: accepted with bounded ingestion as source-015
  decisions: []
  facts:
    - inspected targets implement durable artifacts, gates, or review loops
    - sdd-skill exposes a concrete version/status-authority drift example
  assumptions:
    - repository snapshots represent the inspected public revisions
  artifacts:
    - evidence/source-015-competing-skills-study.md
    - evidence/concepts/PZL-172-evidence-reconciliation-envelope.md
    - evidence/concepts/PZL-173-subordinate-run-ledger.md
  evidence:
    - primary repository files linked above
  conflicts: []
  scope_deviations: []
  next_action: carry PZL-172 and PZL-173 into S06; do not authorize a spike before S03 and S05 settle their prerequisites
```
