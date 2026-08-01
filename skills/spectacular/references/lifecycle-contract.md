---
description: Canonical lifecycle contract for every Spectacular entity, including transition gates, evidence, ownership, and archival behavior.
when_to_use: Defining, implementing, reviewing, or migrating any status-bearing Spectacular record.
---

# Lifecycle Contract

This file is the single authority for lifecycle vocabulary. Other references explain workflows and link here; they do not invent competing enums. A `status:` value is always interpreted in the context of its entity type.

Actual code is the source of implemented behavior. Markdown specifications are execution context and may become stale. `implemented` therefore records a historical verified point, never continuous synchronization.

## Entity contracts

| Entity | States | Transition gate | Lifecycle owner | Terminal/archive behavior |
|---|---|---|---|---|
| Root anchors | none | Snapshot before meaningful edit; review before accepting the unversioned file | File location + `version:` | Snapshots preserve history |
| Idea | `parked → exploring → promoted` | Promotion is explicit | Idea entry | Promoted entry moves to `archive/ideas/` |
| Question | `open ↔ deferred → archived (from resolved)` | Deferral needs a reason; resolution records answer provenance before archive | QUE entry | `archive/questions/`; satisfied canonical dependency, excluded from active fog |
| Research/spike | `open → running | deferred | completed` | Completion requires evidence and `result: supported | refuted | inconclusive` | RES/SPK entry | Completed records remain durable; inconclusive does not clear fog |
| Decision | `verified → superseded` | User choice, unambiguous sourced user intent, or the narrow authorized AFK evidence gate | DEC entry | Immutable durable history |
| Memory | `active → superseded | retracted` | Correction links a replacement or records a retraction reason | Memory entry | Never silently edited or deleted |
| Specification | `draft | unconfirmed → approved → implemented`; any state may close to `archived`, while implemented may first become `superseded | deprecated` | See specification gates below | SPC entry | Archived file moves to `archive/specs/` with `archived_from` |
| Request | `planned → active → review → verified → archived` | Existing policy, verification, docs-impact, and archive closure gates | Request `PLAN.md` | Folder moves to `archive/<slug>/` |
| Feedback | `open → resolved | parked → archived` | Resolution needs a next action | Feedback entry | Explicit feedback archive verb |
| Audit | `open → resolved | folded` | Disposition required | Audit entry | Remains durable in `audits/` |
| Fix | terminal verified record | Concrete `verified-by` evidence required at creation | Fix entry | Permanent trusted corpus |
| Session | `open → closed` | At most one open session; docs impact assessed at end | Session entry | Durable time record |
| Debug run | `investigating → researching → planning → fixing → verifying → resolved | folded | wont-fix` | Existing debug trace contract | Trace remains durable |
| AFK run | `active → gated → active | completed | cancelled` | Scope, allowed actions, declared HITL gates, and approval to resume | AFK run entry | Durable authorization audit trail |
| Release ledger | `candidate → planned → active → shipped`, or `candidate|planned|active → cancelled` | Human release planning | Roadmap ledger row | Shipped/cancelled history remains indexed |

`PRT` and `TSK` are reserved identifiers only. They have no standalone collection or lifecycle until a concrete consumer is designed. Prototype artifacts inherit the lifecycle of their owning request, vision, feedback entry, or spike; tracer bullets use the approved specification/request lifecycle because their code is production code.

Advanced engineering collections may be reserved at init with `--with findings,fixes,bugs,security,benchmarks`. `FND`, `BUG`, `SEC`, and `BMK` define future identity/path contracts only; they have no lifecycle, mutator, or autonomous behavior yet. `FIX` is the reserved padded successor to the active legacy `F<N>` verified-fix ledger. Activating it requires an explicitly confirmed migration. Empty optional folders therefore indicate capability intent, not implemented workflow.

## Specification gates

```text
draft ─────────┐
               ├──→ approved ──→ implemented ──→ archived
unconfirmed ───┘                      │
      │             │                 ├──→ superseded ──→ archived
      └─────────────┴──→ archived     └──→ deprecated ──→ archived
```

- `draft` is collaboratively authored and unfinished.
- `unconfirmed` is produced during an authorized AFK run and awaits human review.
- `approved` is explicitly authorized for implementation. Only this state may seed execution requests.
- `implemented` requires `implemented_at` and `verified_against` plus verified linked requests and closed documentation impact. It says only that the spec was proven at that commit/build.
- `superseded` means a newer implemented SPC declares `supersedes: <old-SPC>`.
- `deprecated` remains visible because behavior may still exist but should not be used.
- `archived` means the detailed execution contract has left active context: implemented/merged, superseded, deprecated, rejected, or abandoned. `archived_from`, dates, reason, and any verification reference preserve which terminal fact occurred.

Minor wording, formatting, typo, or evidence-link corrections do not create a revision. A behavior-changing revision creates a new draft/unconfirmed SPC. The old SPC remains implemented until the replacement becomes implemented; that transition atomically marks the old record superseded.

Requests own `source_spec`. Spec-to-request backlinks are derived so one SPC can safely produce several requests. Replacement specs own `supersedes`; `superseded_by` is derived. An archived implemented predecessor remains immutable; the replacement's `supersedes` link is sufficient.

## Decision authority

A DEC may be written from an explicit user choice or from unambiguous user intent already recorded in a conversation, interview, PRD, or approved spec. Derived capture may summarize but cannot add a choice. Ambiguity becomes a high-priority QUE with options, a recommendation, evidence for/against, and the exact user decision needed.

The autonomous AFK exception requires all of: an active goal-scoped AFK run; permission for decisions; a technical, in-scope, reversible choice; validated research/spike/vendor/interview evidence; recorded alternatives; no unresolved product/business trade-off; and no HITL boundary. Failure creates or updates a QUE instead of a DEC.

## Discovery evidence

`status` records progress; `result` records learning. A completed, refuted experiment can satisfy the question it tested by ruling out a path. A completed inconclusive record remains fog and must link a narrower successor question, research record, or spike before dependents proceed.

Discovery is progressive, not mandatory ceremony. Apply [[discovery-protocol]] before creating a node: inspect current code/tests/docs or ask the user directly when that is sufficient; research bounded facts; spike technical feasibility; attach prototypes only when human interaction is the evidence; execute a tracer bullet only from an approved `SPC` with `execution_mode: tracer`. Generic artifacts and technical debt do not receive standalone lifecycles.

## Retention and freshness

[[artifact-retention]] is the authority for live, stale-safe, temporary, and throwaway classification. Retention is derived from entity/status/path rather than duplicated in frontmatter. Open questions and live indexes are checkpoint-synchronized; detailed specs become stale-safe archive after verified integration; code and executable invariant/unit tests are implementation truth. Durable Markdown is archived rather than purged. Disposable branches/mocks require extracted evidence and a recovery boundary before deletion.

## Documentation impact

Every request begins with `docs_impact: pending`. Assessment is mandatory after a major behavior/architecture update, before a heavy request enters verification, and at session end.

- `docs_impact: none` requires `docs_impact_reason`.
- `docs_impact: required` requires `docs_impact_evidence` before verification/archive.
- Intentional deferral uses the existing recorded override mechanism; it is never silent.

Documentation is a completion gate, not another lifecycle state. Pageworks owns public documentation work; Spectacular owns the impact declaration and evidence/handoff.

## Compatibility and migration

Legacy records and singular collection folders remain readable until an explicitly confirmed migration. Migration is preview-first and archive-first. Specification mappings are conservative:

```text
draft                    → draft
unconfirmed              → unconfirmed
stable/published/current → approved
deprecated               → deprecated
```

No legacy label becomes `implemented` without a user-supplied commit/build reference. Canonical folders are plural collections: `ideas/`, `questions/`, `research/`, `spikes/`, `decisions/`, `specs/`, `memories/`, `sessions/`, `feedbacks/`, `audits/`, `fixes/`, and `debugs/`. Optional advanced paths are `findings/`, `bugs/`, `security/` (singular domain name), and `benchmarks/`; `fixes/` is both an existing collection and part of that advanced init set.

`spectacular collections migrate` previews singular-to-plural moves. `spectacular lifecycle migrate` previews conservative label mappings (`stable|published|current → approved`; statusless legacy memories → `active`). Both require `--apply --yes` and archive originals plus a mapping before mutation. Ambiguous `verified` discovery records are reported for result/evidence review and never guessed.

## Enforcement

CLI verbs own mechanical transitions. `doctor lifecycle`, `doctor specs`, `doctor wayfinding`, and collection-specific doctor areas validate closed enums, required evidence, unique slugs within each collection, relationship integrity, and archive placement. Fog/frontier and inverse links are derived views, never stored status.
