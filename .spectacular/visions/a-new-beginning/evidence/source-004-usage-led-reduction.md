---
type: source-card
source: source-004
provided_as: source4
received: 2026-08-07
authority: proposal
status: ingested
scope: [usage-evidence, protected-core, subsystem-cuts, policy, agents, cli-size]
duplicate_sections: [source-003, source-001-selected]
completeness: partial
---

# Source 004 — Usage-led reduction

## Thesis

Use Spectacular's self-hosting workspace as a usage log: preserve the heavily used
convention/request/spec/decision core, delete sparse or unused collections and
reserved identities, consolidate retrieval fragments, then review wayfinding,
discovery, Vision, AFK, policy, agents, and the Bash implementation as increasingly
contested simplification decisions.

## Source integrity

The source repeats Source 003's full reduction plan and selected Source 001
cleanup proposals. That repetition is tracked as duplicate provenance, not
independent support. Source 004 supplies the previously missing sessions 7–9,
including discovery-type collapse, policy review, and the Bash strategy question.
Some tables remain truncated in the supplied text.

## Measurement audit

The source's counts mix incompatible units:

| Claim | Current observation | Intake judgment |
|---|---|---|
| 71 archived requests / 257 archive files | Not loaded during normal operation; source-provided only | unverified here |
| decisions 28 | `summary --json` reports 27 decision records | close but stale/different counting |
| memories 4 | read verb reports 3 records | likely counted index file |
| sessions 1 | summary reports 0 sessions; folder has index only | structural file, not usage record |
| visions 0 | active refactor Vision now exists | stale |
| ideas 8 | current idea list reports 8 | supported |
| specs 12 / SPC 1001 | live spec list reports 9 records | record count and text occurrence conflated |
| DEC unused | 27 current decisions exist | internally contradicted |
| FIX unused | no FIX IDs allocated; five live legacy F1–F5 fixes exist | reservation and workflow conflated |
| policy hooks 8 | current POLICY defines 9 | stale |
| agents 9 | `agents/*.md` contains 9 | supported |
| doctor areas 19 | doctor help lists 19 | supported |
| CLI 16,471 lines | current CLI is 16,507 lines | directionally current, exact count stale |

Counts must specify whether they measure records, files, indexes, textual
occurrences, live state, archived history, or external use. Without that schema,
the survival table cannot support deletion decisions.

## Proposed protected core

- `.spectacular/` directory-as-contract convention;
- project intent plus request and archive flow;
- PLAN/TASKS frontmatter lifecycle and cheap briefings;
- decision records;
- specifications and archive/spec synchronization;
- deterministic init/status/doctor/closure mechanics around that core.

The supplied must-have table is truncated after specs, so later intended rows are
not reconstructed.

## New proposals, normalized

| Domain | Proposal | Intake state |
|---|---|---|
| Evidence | Treat self-hosting record counts as survival evidence | useful input; invalid as sole rule |
| Protected core | Define the load-bearing core before cuts | strong |
| Discovery | Collapse QUE/RES/SPK/PRT into QUE with `kind:` | disputed; behavior differs |
| Policy | Reduce nine current hooks to three risk gates | promising goal; mapping unproven |
| Agents | Reduce nine agents to five | unsupported by usage evidence |
| Bash | Treat CLI line count as a surface budget; decide freeze/extract/rewrite last | useful metric; target speculative |
| Vision | Delete Imagine/Vision due zero use | refuted by current active use and recency |
| AFK | Remove AFK/worktree/traffic infrastructure as speculative | contested; recently implemented and used locally |
| Wayfinding | Remove scheduler-like surface for solo queue | contested; live differential test shows unique output |
| Historical target | Prefer the simpler v0.4-sized shape over v1.37 accumulation | useful comparison, not rollback proof |

## Assumptions and contradictions to resolve

1. Self-hosting frequency is representative of all target users and future value.
2. File/occurrence counts can stand in for entity records and workflow outcomes.
3. Rare safety and coordination capabilities are low-value because they produce
   few records.
4. Recently shipped capabilities had enough observation time to establish non-use.
5. QUE, RES, SPK, and prototype behavior differs only by name rather than
   authorization, evidence contract, or execution semantics.
6. Three generic risk hooks preserve the distinctions currently enforced across
   planning, implementation, debugging, verification, and closure.
7. Agent-definition count measures runtime cost even when agents load only on dispatch.
8. A 6k-line Bash target follows from the chosen behavior rather than driving cuts.
9. Earlier-version size correlates with adoption more strongly than later capability quality.

## Valuable ideas independent of the cut list

- Define a protected core before evaluating deletions.
- Normalize usage telemetry by artifact type and lifecycle.
- Distinguish reserved compatibility surface from active workflow surface.
- Require unique behavioral value, not merely code existence, for a subsystem to survive.
- Separate reducing default scaffold cost from deleting runtime capability.
- Make implementation-language strategy follow product-surface decisions.

## Provisional assessment

**Strong:** protect the proven core; normalize evidence; review each subsystem
independently; postpone Bash strategy until behavior settles.

**Needs evidence:** discovery collapse, policy reduction, agent reduction, and
the 6k-line target.

**Current deletion cases fail their own evidence standard:** Vision is active,
AFK has local runs, and wayfinding produces a distinct next action. They may still
be removed, but not as “provably dead.”

No product or deletion decision is recorded by this assessment.

## Repeated concept provenance

Source 004 repeats PZL-023 through PZL-036. Its stub/registry/doctrine and init-ID
content also overlaps PZL-003, PZL-004, PZL-005, and PZL-009.

## New concept pieces

- [PZL-037 — Normalize self-hosting usage evidence](concepts/PZL-037-normalize-usage-evidence.md)
- [PZL-038 — Define a protected core](concepts/PZL-038-protected-core.md)
- [PZL-039 — Collapse discovery types](concepts/PZL-039-collapse-discovery-types.md)
- [PZL-040 — Three risk-gate policy model](concepts/PZL-040-three-risk-gate-policy.md)
- [PZL-041 — Reduce the agent fleet](concepts/PZL-041-reduce-agent-fleet.md)
- [PZL-042 — Bash surface budget](concepts/PZL-042-bash-surface-budget.md)
- [PZL-043 — Retire Imagine and Vision](concepts/PZL-043-retire-imagine-vision.md)
- [PZL-044 — Retire AFK coordination](concepts/PZL-044-retire-afk-coordination.md)
- [PZL-045 — Retire wayfinding](concepts/PZL-045-retire-wayfinding.md)
- [PZL-046 — Use the earlier simple shape as a comparison](concepts/PZL-046-earlier-shape-comparison.md)

## Decision packets seeded

- What normalized usage evidence can legitimately support removal?
- What exact capabilities form Spectacular's protected core?
- Are discovery types behaviorally distinct enough to keep separate?
- Which policy distinctions prevent real failures?
- Which agent roles produce unique value when dispatched?
- What CLI size follows from the accepted surface, and what architecture fits it?
