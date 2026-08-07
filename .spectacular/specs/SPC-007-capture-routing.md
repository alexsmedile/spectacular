---
id: SPC-007
type: specification
status: draft
target_version: "tbd"
supersedes: ""
updated: 2026-08-07
summary: "Govern source-neutral agent signal capture, explicit destination authority, and external adapter boundaries without duplicating local IDEA state"
related:
  - "GitHub Issue alexsmedile/spectacular#6"
  - "SPC-005-idea-destination-routing.md"
  - "../../STATUS.md"
  - "../../skills/spectacular/references/github-work-bridge.md"
---

# SPC-007 — Govern source-neutral agent signal capture, explicit destination authority, and external adapter boundaries without duplicating local IDEA state

## Intent

Let a human or agent turn a useful signal into a deliberate, appropriately
bounded destination without creating a second backlog, publishing sensitive
material, or treating capture as a commitment to build.

Capture is source-neutral: a person, a coding agent, an IDE integration, or a
provider-specific filing assistant may supply the same redacted proposal. The
destination—not the tool that observed the signal—owns its discussion and
lifecycle.

```text
capture → route to a deliberate destination
                 ├─ issue
                 ├─ shared
                 ├─ roadmap
                 ├─ vision
                 └─ request

local IDEA = compatibility/private-offline fallback, never the default authority
```

## Requirements

### R1 — A bounded capture proposal

The capture surface must produce or accept a proposal envelope before any
external publication. It contains only the information needed to evaluate the
route:

- a concise claim or observed opportunity;
- minimal, redacted evidence references (for example, a commit, file path,
  command/result, or reproducible observation);
- source and run provenance when available;
- confidence, sensitivity, and a proposed destination; and
- likely duplicate references when the proposed destination is an Issue.

The envelope must not include raw chat history, credentials, secrets, personal
data, or copied remote Issue/Discussion bodies. Missing evidence is allowed;
the resulting proposal must say so rather than inventing confidence.

### R2 — Destination authority and routing

Routing considers sensitivity first, then whether the signal belongs inside
current work, its audience/scope, uncertainty, and commitment. Current-work
checkpoints and audits resolve existing internal work before the five approved
capture destinations are selected. The selected destination becomes the
authority for its conversation and lifecycle.

| Destination | Use when | Authority and constraint |
|---|---|---|
| current request checkpoint | A defect or friction is introduced or found while doing already-owned work. | Keep the regression evidence in that request's TASKS/SESSION/verification flow; do not manufacture an Issue. |
| audit | An independently observed symptom is unclear. | Use the existing investigation path before proposing a fix or public report. |
| `issue` | An independently actionable, repository-specific collaborative report or proposal is ready for human review. | The remote Issue owns discussion; it is not a roadmap, SPC, or implementation commitment. |
| `shared` | Collaboration spans repositories or needs a provider-configured external surface. | The external record owns discussion. Spectacular stores only a stable reference when needed. |
| `roadmap` | A human has intentionally made a project-priority choice. | Explicit Icebox/version/tier placement is required; capture never infers a commitment. |
| `vision` | Product direction needs grounded alternatives and human reaction before an implementation contract. | Vision remains opt-in and can derive only a draft SPC after approval. |
| `request` | The outcome is defined and needs durable execution coordination. | Create only after an explicit human decision and normal lifecycle gates. |
| local `IDEA` | The thought is private, offline, or needs local refinement before another route is known. | Compatibility/private fallback only; never the default shared backlog. |

Protected or security-sensitive material stops ordinary routing. It must return
a redacted blocker and follow the project’s protected handling path; it must
never be emitted to ordinary Issue or shared adapters.

### R3 — External adapter boundary

Spectacular defines a destination-neutral proposal and routing contract, not a
default GitHub, chat, IDE, or tracker integration. An external provider adapter
may accept a proposal envelope and render a destination-specific draft, but it
must declare its provider, destination reference, and publication capability.

The initial contract may render/copy a proposal or record a stable external
reference. It must not silently create, comment on, label, close, synchronize,
or otherwise mutate an external record. A remote Issue create remains
dry-run/preview first and requires explicit human confirmation at the point of
publication. A `shared` destination is provider-configured external state, not
a new committed Spectacular collection.

No background daemon, polling loop, or ambient automatic capture is introduced.
Capture happens only at an explicit user/agent checkpoint or through a
separately invoked adapter.

### R4 — Duplicate handling

Before an Issue proposal is presented for publication, the responsible adapter
or human filing workflow searches current Issues and returns likely duplicates
with its evidence. The human may choose a candidate, revise the proposal, or
continue with a new Issue. The workflow must not autonomously comment, react,
close, merge, or alter the candidate records.

### R5 — Autonomy and approval boundary

The baseline required by this SPC is proposal-only. Capture itself never grants
authority to create an Issue, create a request/specification, open a PR, merge,
close an Issue, disclose protected content, or expand current scope.

A future managed-publication grant, if desired, needs a separate approved
contract that defines trusted repository scope, reproducibility threshold,
sensitivity exclusions, audit provenance, revocation, and human override. It
is explicitly out of scope for this SPC.

### R6 — IDEA deprecation path and compatibility

Existing `IDEA` records, IDs, archives, undo behavior, and local refinement
remain compatible. Guidance should favor the authoritative eventual
destination, while local IDEA storage is reduced to private/offline or
compatibility use. There is no automatic migration, mirroring, copy, or
synchronized lifecycle between local and external stores.

Private capture must be designed as private state before it is added (for
example, under `.spectacular.local/`); it must never be silently committed or
included in shared listings. Any later removal or numerical reduction of local
IDEA behavior requires adoption evidence and a separate major-release decision.

## Boundaries

- This replaces neither the current Issue/job-card triage contract nor the
  request, audit, Vision, roadmap, or protected-security lifecycle.
- It does not introduce a GitHub-only API, provider credentials, a generic
  remote database, managed repository setup, labels/rulesets, or webhook/event
  synchronization.
- It does not classify signals by silently inspecting chats, telemetry, or
  private workspaces.
- It does not make a captured idea a roadmap commitment, an approved SPC, or a
  request.

## Acceptance criteria

- A source-neutral, redacted proposal envelope can be routed without copying
  chat or remote-body content into a second local inbox.
- The route distinguishes current-work checkpoints and audit investigation
  from externally collaborative Issue/shared destinations.
- Each listed destination states its authority and the explicit gate that
  prevents accidental priority or execution commitment.
- `issue` publication is duplicate-aware, dry-run/preview first, and requires
  explicit confirmation; no ordinary external mutation happens implicitly.
- Protected signals are stopped before ordinary publication and return only a
  redacted protected-path result.
- `shared` remains a provider-configured external destination, while local
  IDEA records remain compatible private/offline fallbacks with no mirroring.
- The design introduces no background daemon, polling, automatic capture, or
  managed-autonomous publication path.

## Proposed implementation and validation

The implementation must be sequenced only after this draft is explicitly
approved. The likely slice is limited to routing guidance, a proposal
envelope/rendering boundary, compatible IDEA guidance, and deterministic test
fixtures; it must not add a remote-write adapter by implication.

Validation must demonstrate at least:

1. a private/offline capture, an Issue proposal with duplicate candidates, a
   shared-reference handoff, roadmap placement, Vision handoff, and request
   handoff;
2. an in-request defect and an unclear independent symptom routing to their
   existing checkpoint/audit paths rather than a new Issue;
3. redaction and protected-signal rejection before ordinary rendering or
   publication;
4. absence of network mutation without a distinct explicit publication gate;
   and
5. preservation of existing IDEA IDs, archives, undo behavior, and
   no-mirroring guarantees.

## Evidence and decisions

- [GitHub Issue #6](https://github.com/alexsmedile/spectacular/issues/6),
  including its expanded proposal for governed agent signal capture.
- `STATUS.md` (2026-08-07): approved source-neutral capture-routing model and
  constraints for external destinations, confirmation-gated remote creation,
  and local IDEA fallback.
- `SPC-005-idea-destination-routing.md`: implemented destination-routing
  precedent; this successor preserves its no-mirroring and compatibility
  decisions while covering the broader source/adapter boundary excluded there.
- `skills/spectacular/references/github-work-bridge.md`: GitHub owns capture
  and discussion; Spectacular owns local reasoning and coordinated execution.
- `skills/spectacular/references/idea-rules.md` and
  `skills/spectacular/references/soft-db-index.md`: existing local IDEA and
  destination-routing contract.

## Confirmation

draft — created from the approved routing model. Not eligible for implementation
until the maintainer explicitly approves this exact SPC.
