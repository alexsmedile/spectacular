---
type: handoff-return
handoff_id: H01
status: complete
disposition: accepted-as-evidence
baseline: 7a85469618d984f35755a3b5db0684bfa3620c05
source_thread: 019fde9a-e43e-7c03-9ebf-d5c39dd66d9c
authority: evidence-only
---

# H01 — Product-boundary evidence audit

## Executive finding

The current PRD describes a lightweight file convention, an ongoing skill, and a
one-time bootstrap CLI. The shipped product demonstrably also governs specs and
requests, sequences dependencies, enforces lifecycle and verification gates,
coordinates agents and concurrent work, manages Git safety, and prepares GitHub
handoffs. The full CLI suite passed 31 of 31 test files during H01, so this is
implemented responsibility rather than documentation drift.

The strongest evidence-backed core candidates are shared human/agent state,
progressive retrieval, cross-session continuity, separation of accepted intent
from implementation truth, and evidence-backed safety. Repository tests establish
mechanisms, not unique user value. H01 therefore identifies the product-boundary
conflict without deciding whether Spectacular should become a knowledge workspace,
a project control plane, or an integrated executor.

## Claimed product versus actual responsibility

| Area | Claimed product | Verified actual responsibility | Assessment |
|---|---|---|---|
| Product layers | Convention, ongoing skill, CLI used once for bootstrap | Ongoing CLI mutators/readers for requests, specs, discovery, policies, migrations, Git, GitHub, AFK, and coordination | Direct contradiction |
| Work management | Not a ticketing or project-management system | Priorities, requests, milestones, holds, dependency graphs, roadmaps, sessions, lifecycle, traffic, next action | Strong category collision |
| Orchestration | Multi-agent orchestration is out of scope | Agent fleet, dispatch/reconciliation, audits, AFK, traffic preflight, worktree coordination, agentic activation | Direct responsibility contradiction |
| Authority | Humans decide irreversible actions; Git/provider collaboration remains separate | Spectacular is primary lifecycle authority; merge remains human-gated | Expanded scope with preserved human gate |
| Portability | Markdown-only, no mandatory tooling, text-editor operable | Canonical Markdown survives, but IDs, migrations, projections, graph checks, and lifecycle gates rely on CLI behavior | Manual parity unverified |

## Load-bearing behavior candidates

1. Shared repository-readable state for humans and agents.
2. Cheap progressive retrieval through frontmatter, status, and bounded projections.
3. Cross-session continuity without duplicating runtime micro-work.
4. Separation of intended contracts, implementation truth, evidence, and history.
5. Evidence-backed and reversible boundaries around closure and external effects.
6. Bounded uncertainty routing; implemented, but unique user value remains unproven.

## Inherited constraints requiring confirmation

| Constraint | Evidence assessment |
|---|---|
| Human/agent-readable canonical state | Load-bearing outcome candidate |
| Progressive disclosure and cheap signals | Load-bearing candidate |
| Intent/implementation/evidence authority separation | Load-bearing candidate |
| Human authority over destructive/external actions | Load-bearing trust candidate |
| Markdown with no database/server of any kind | Human-readable canon supported; absolute ban is reopenable |
| Git-committed by default | Likely load-bearing; exact visibility defaults remain open |
| Bash 3.2 | Current implementation constraint, not proven product identity |
| Exact taxonomy, lifecycle, IDs, and folders | Reopenable mechanisms |
| Snapshot scheme and packaging | Reversibility/context outcomes matter; exact mechanism is reopenable |
| One-time bootstrap CLI | Contradicted by shipped behavior; not load-bearing |

## Missing evidence that could change S01

- Direct cohort interviews and observed usage.
- Real-project frequency and retention evidence.
- A complete field trial of intent through cold resume.
- Comparative evidence for workspace, control-plane, and executor boundaries.
- Failure/recovery evidence from autonomous or parallel work.
- Maintenance cost per subsystem.
- Practical proof of tool-free operation beyond reading files.
- Evidence separating Spectacular's value from host runtime, Git, and GitHub.
- Messaging tests for workspace, mission-control, and control-plane language.

## Framing question

> Which single user outcome must Spectacular own end to end, and at what exact
> boundary must it hand execution or collaboration to the coding agent, Git, or
> GitHub?

## Return packet

```yaml
return:
  handoff_id: H01
  status: complete
  baseline: 7a85469618d984f35755a3b5db0684bfa3620c05
  result: >-
    The implemented product has crossed its documented boundary: it remains a
    Markdown/Git-native coherence workspace, but also acts as a deterministic
    lifecycle and coordination control plane. Evidence does not decide whether
    that expansion should be retained.
  decisions: []
  facts:
    - The PRD calls the CLI a one-time bootstrap tool.
    - The CLI exposes ongoing lifecycle, discovery, policy, Git, GitHub, autonomy, and coordination commands.
    - The PRD excludes PM and multi-agent orchestration products.
    - The shipped capability index assigns Spectacular those responsibilities.
    - Verified decisions make Spectacular primary lifecycle authority and require PR handoff.
  assumptions:
    - Passing tests demonstrate implemented behavior, not product-market value.
    - Inherited decisions are evidence, not refactor-preservation authority.
    - No external usage or interview evidence was available.
  artifacts: []
  evidence:
    - ./cli/spectacular --help
    - bash tests/run.sh: 31 test files passed, 0 failed
    - PRD, specs index, decisions DEC-018/DEC-019, CLI, skill, agents, and tests
  conflicts:
    - One-time CLI claim versus ongoing CLI responsibility
    - Orchestration non-goal versus shipped coordination behavior
    - PM non-goal versus owned work/lifecycle state
    - No-mandatory-tooling promise versus unverified manual parity
  scope_deviations: []
  next_action: Attach this evidence to H02 and H03; do not infer the future product boundary from H01.
```
