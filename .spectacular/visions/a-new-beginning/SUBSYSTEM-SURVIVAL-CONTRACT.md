---
type: foundation-contract
contract: subsystem-survival
version: 1.0
status: accepted
owner: alex
source_session: H17
accepted: 2026-08-09
---

# Subsystem Survival and Responsibility Contract

## Purpose

This accepted S10 contract fixes which v1 jobs survive into Spectacular v2, which merge into the
smaller semantic core, which move behind replaceable boundaries, and which retire. It chooses
responsibilities and survival outcomes, not source language, modules, paths, schemas, commands,
storage, migration mechanics, or deletion operations. S11 owns those implementation choices.

Every survival decision applies the accepted test: distinct protected job, credible benefit,
simplest adequate replacement, sustainable maintenance/context/recovery cost, and a safe history
and recovery path. Current repository use is scoped proxy evidence, not external-value proof.

## Acceptance normalizations

- **Mission Director** is an internal role contract for Spectacular's control-plane behavior, not a
  new authority entity, lifecycle owner, or provider.
- **Receipt** is an attributable boundary/result shape attached to its semantic owner. It is not a
  top-level truth collection.
- **Project Direction** describes authoritative Anchor content—purpose, promises, boundaries,
  current focus, and critical constraints. S11 may choose its representation but may not create a
  competing authority.
- Roadmap extraction establishes an optional, replaceable boundary. It does not approve GitHub
  Projects or any companion as a mandatory dependency.
- Imagine, Grill, Refine, Review, Verify, Assess, and Audit are guided techniques or behaviors.
  Their exact public command exposure remains constrained by the accepted S09 grammar and is an
  S11 interface decision.

## Retained semantic core

### Direction and accountable change

- **Project Anchors** preserve purpose, promises, boundaries, current focus, critical constraints,
  current Contract entry points, and Project Guardrails.
- **Imagine** preserves divergent-first exploration. A simple idea stays within a draft Proposal;
  complex exploration earns a Proposal-owned workspace.
- **Proposal** owns a candidate behavior change or bounded target delta. Acceptance authorizes the
  target; it does not alter current Contract truth.
- **Mission** remains the sole durable accountable delivery unit, with Objectives, Runs, local
  Tasks, preparation, authority, Evidence, recovery, and closure.
- **Capability Contract** remains the authority for current owner-accepted observable behavior.

### Truth, learning, and retrieval

- **Decision** records a settled resolution; open matters remain Proposals or Gaps.
- **Gap** owns a consequential unknown, assumption, question, or discovery need.
- **Evidence** is an attributable scoped result or observation and never performs acceptance.
- Semantic records are embedded by default and promoted globally only when cross-boundary durable
  identity is earned.
- Stable identity, pointers, lean searchable metadata, deterministic rebuildable indexes,
  drill-down, and promotion continuity provide collection-like retrieval without collection
  sprawl. Generated indexes remain non-authoritative projections.
- **Glossary** survives as versioned global semantic authority.

### Execution and continuity

- **Run** owns one attributable resumable Mission attempt and remains independent of host chat
  sessions.
- **Checkpoint** is an embedded meaningful Run boundary, not a separate activity diary.
- **Receipt** attaches provider, execution, review, handoff, or result evidence to its natural
  semantic owner.
- **Autopilot** is an explicit Mission consent envelope. It grants no ambient authority and is not
  an executor.
- Execution topology is prepared before activation. Objectives expose candidate lanes; dependency,
  independence, cancellation-state, evidence, and join analysis determine valid concurrency.

## Mission orchestration model

Spectacular acts as Mission Director/control plane while the host runtime executes and native
providers perform their own effects.

- Small jobs run inline; delegation must be earned.
- Missions never nest. An Objective may have an optional Lead.
- Investigator and Builder are execution roles; Reviewer and Verifier supply assurance under the
  accepted risk and independence rules.
- Nested delegation is disabled unless depth and budget are explicitly authorized.
- The owner has one human-facing channel.
- A mutating Mission gets a branch when isolation is needed; workers do not receive branches merely
  because they exist.
- Git/provider effects remain inside the Mission consent envelope and provider permissions.
- AFK, workspace preflight, and traffic jobs merge into Autopilot, Mission start/resume/closure,
  prepared/live topology, and provider adapters.

## Project Guardrails

The fixed v1 POLICY hook/schema subsystem retires. Human-authored Markdown **Project Guardrails**
survive as authoritative project operating constraints selected mechanically and interpreted by
the agent.

- A heading identifies a guardrail; no explicit guardrail ID is required.
- The first non-empty line under the heading is an order-independent selector line containing one
  or more defined `@Event` tokens and optional guided/extensible `$domain.verb` tokens.
- Remaining Markdown is unrestricted and injected verbatim.
- Guardrails cannot create authority, weaken accepted invariants, prove provider facts, or enforce
  unrelated harnesses.
- Native hooks are optional and are never edited automatically.
- S11 owns event vocabulary, command catalog, parser, result protocol, paths, and adapters.

## Integrity, repair, and migration

- Keep a **scoped read-only v2 integrity capability** supporting trustworthy prepare, resume,
  execute, resolve, assess, and reconcile operations.
- Retire the universal doctor-repair subsystem. Deterministic correction belongs to the operation
  that owns the invariant; judgment repair routes through Gap, Proposal, Decision, or Mission;
  code repair belongs to Mission execution; provider repair belongs to the provider.
- Extract all v1→v2 conversion into the isolated S11 migration capsule. v2 core contains no legacy
  parser, aliases, fallback reads, dual writes, lazy conversion, historical archive tree, or
  general legacy migration registry.
- No retirement authorizes deletion until unique truth, replacement, rollback, and recovery are
  proven through the immutable v1 release/tag and validated snapshots.

## Guided authoring kernel

The v1 universal mode/slot engine, per-document rule sprawl, and template-as-truth architecture do
not survive. Their useful job merges into one proportional guided-authoring kernel:

```text
Orient → Diagnose → Question → Compare → Propose → Resolve → Review → Receipt
```

Artifact semantic contracts own sufficiency. Minimal representations, optional scaffolds,
deterministic validation, and conformance tests support the kernel. Imagine, Grill, Refine, and
Review remain natural techniques:

- Imagine diverges before convergence.
- Grill challenges and converges through explicit owner choices.
- Refine clarifies without changing accepted meaning.
- Review evaluates against a named contract or rubric.

Exact schemas, paths, templates, rules, and public routes remain S11/S12 decisions.

## Evidence and assessment loop

The alternate VERIFY.md/VERIFY-LOG authority and fixed verification ceremony retire. Their useful
job merges into one claim-centered loop:

- **Verify** gathers and evaluates claim-appropriate Evidence.
- **Review** produces findings and attributable assurance.
- **Assess** determines whether Evidence is sufficient for the requested readiness or disposition
  question.
- **Audit** tests historical, fleet, reconciliation, or other consequential claims on demand.

Providers and executors produce primary evidence; Spectacular maps claims, validates envelopes,
and guides assessment; independent actors or methods review when S06 requires it; only the owner
resolves a Mission and separately authorizes Contract reconciliation. Code, commits, PRs, checks,
and provider receipts are scoped Evidence, never sufficient closure by themselves.

## Small role-contract library

The fixed nine-agent fleet retires. v2 keeps a small progressively loaded role-contract library:

| Role | Unique job | Boundary |
|---|---|---|
| Mission Director | Compile and govern one Mission's preparation, authority, topology, continuity, and owner channel | Does not execute provider effects or gain owner authority |
| Objective Lead | Optionally coordinate one Objective when its internal complexity earns a lead | Does not own Mission lifecycle |
| Mission Slice Advisor | Generate and compare coherent Mission slices | Cannot approve or activate a Mission |
| Design Sufficiency Reviewer | Skeptically test whether the proposed direction is sufficiently understood | Cannot declare sufficiency or choose product trade-offs |
| Investigator | Gather repository, localization, diagnosis, or external-research evidence through loaded lenses | Does not mutate unless separately authorized |
| Builder | Implement accepted change or bounded repair through loaded lenses | Cannot widen Mission scope or authority |
| Reviewer | Produce contract-, code-, security-, audit-, or experimental narrative findings | Cannot accept or reconcile |
| Verifier | Gather and evaluate independent or deterministic proof against claims | Cannot resolve the Mission |

Bug, repository, research, build, repair, code, Contract, security, audit, and narrative expertise
become progressively loaded lenses. Dispatch is earned; no current agent filename survives merely
because a role remains useful.

## External and companion boundaries

- Spectacular owns Project Direction reasoning, accountable-work formation, Mission
  orchestration, authority/evidence/recovery context, and owner communication.
- Ongoing roadmap sequencing, portfolio views, shared release horizons, assignments, dates, and
  notifications move behind an optional replaceable companion/provider boundary. Spectacular
  retains typed pointers and a minimal handoff; exact integration remains unapproved.
- Feedback inboxes and capture remain provider-owned. Spectacular routes consequential input into
  the correct Gap, Proposal, Evidence, Decision, Contract, Anchor, or Mission path.
- Git, GitHub, CI, deployment, and other providers perform and attest their native effects.
- Pageworks, Bugworks, and Specwright retain only their previously accepted optional-specialist
  boundaries; this contract does not authorize new dependencies or implementations.

## Final v1 disposition map

| v1 surface | S10 disposition | v2 replacement |
|---|---|---|
| requests | merge | Missions, Objectives, Runs |
| specs | merge | Capability Contracts and Proposals |
| ideas/questions/research/spikes/audits/fixes | merge | Draft Proposals, Gaps, discovery/evidence methods, violations, repair Missions |
| decisions | keep-simplify | Settled-current, rejected, superseded Decisions |
| glossary | keep | Versioned semantic authority |
| archive collection | retire from live v2 tree | Natural-owner history, closure/recovery pointers, frozen v1 recovery |
| Vision | simplify-merge | Imagine plus earned Proposal exploration workspace |
| Wayfinding/generic next | merge-retire | Imagine, Fog, Links, preparation, orientation, resume, safe continuation |
| roadmap/goals/portfolio manager | extract | Project Anchors in core; optional replaceable roadmap companion/provider |
| feedback/memory/session collections | merge-retire | Consequence-owned semantic records and Run Checkpoints |
| AFK/workspace/traffic bridges | merge-simplify | Autopilot, Mission operations, topology, provider adapters |
| POLICY hooks/schema | simplify-replace | Markdown Project Guardrails |
| broad doctor/repair | keep scoped integrity; retire universal repair | Operation-owned correction and Mission/provider repair |
| resident migrations | extract | Disposable v1→v2 capsule |
| universal document engine/rule/template matrix | simplify-merge | Guided-authoring kernel plus artifact contracts |
| VERIFY artifacts and review sweep | simplify-merge | Claim-centered Evidence and Assessment loop |
| nine named agents | merge-simplify | Small role library plus progressive lenses |

## Context and maintenance budget

Surviving capabilities must load only for their owned question and representative scenario.
Acceptance requires cold-start, interrupted-Mission, and material-claim-audit evaluation paired
with correct authority, source drill-down, conflict/Gap exposure, safe recovery, owner
comprehension, and output quality. File, line, token, agent, command, or collection counts are
diagnostics only. A smaller implementation is a regression if it hides constraints, weakens
recovery, increases unsafe inference, or creates a mandatory provider.

## Reserved S11 decisions

S11 owns implementation language and module boundaries; storage and schemas; stable-key and query
syntax; generated index mechanics; promotion continuity; Project Direction representation;
roadmap handoff; Project Guardrails parser and catalogs; role packaging and dispatch interfaces;
provider adapters; verification representation; migration capsule; v1 tag/release and recovery
proof; and exact Skill/CLI routes consistent with S09.

## Exit condition

Every surviving job has one semantic owner, one authority boundary, labeled evidence, an explicit
context/maintenance budget, and a recoverable disposition path. No v1 subsystem survives by name,
current implementation, usage count, or sunk cost alone.
