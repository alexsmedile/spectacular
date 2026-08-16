---
type: Contract
id: 019fe381-5d61-7223-b362-03a5f99a7b10
human_ref: CC-v2prod
title: Spectacular v2 product contract
status: current
created_by: owner
created: "2026-08-10T00:00:00Z"
updated: "2026-08-16T09:55:54Z"
contract_version: "2"
revision_basis: Owner-confirmed Mission model in the 2026-08-16 Codex design session.
purpose: Govern the Spectacular v2 public product surface.
outcome: A human and a cold agent can recover exact current work without hidden chat context.
applies_when:
  - Spectacular stores, retrieves, or presents governed project work.
does_not_apply_when:
  - A native provider performs its own external effect.
does_not_provide:
  - Owner authority, provider credentials, or evidence sufficiency by itself.
required_behavior:
  - Preserve human-readable anchors, Mission bundles, and scoped references alongside UUID identity.
  - Expose source, freshness, gaps, conflicts, and exactly one continuation or owner gate.
  - Treat Proposal as optional, mutable exploration rather than execution authority or current truth.
  - Freeze the selected outcome, completion criteria, Objectives, boundaries, and initial Run in MISSION.md when a Mission starts.
  - Begin a Mission bundle with only MISSION.md; keep simple Objectives and the initial Run inline.
  - Promote an Objective to a dedicated file without changing its identity when detail, delegation, or independent review makes the file useful.
  - Create separate Run files only when a distinct execution job or recovery boundary makes them useful.
  - Treat MISSION.md as the Mission entry point and derive navigation from its frontmatter and referenced files.
  - Record durable architectural or product choices as ADR-like Decisions, not lifecycle approvals.
  - Update applicable product specifications as ordinary Mission work and close them through the same Mission completion gate.
  - Use adaptive preparation only while success criteria, scope, dependencies, risk, or blocking Gaps remain unresolved.
operating_cases:
  - Cold project orientation, Mission resume, governed execution, assessment, and closure.
persistent_information:
  - Anchors, applicable Contracts, Missions, Evidence, ADR-like Decisions, Gaps, continuity boundaries, and optional Proposals.
conformance_checks:
  - Canonical Markdown is understandable from the filesystem.
  - Every emitted pointer drills down through an executable noun-first command.
  - A simple newly started Mission is recoverable from MISSION.md alone.
  - Expansion preserves Objective and Run identity and does not introduce a mandatory Mission index.
  - Mission completion checks its frozen criteria once and presents one owner gate without a separate reconciliation ritual.
authority_freshness: Owner-confirmed Mission boundaries and current source fingerprints are required for consequential mutation.
related_material:
  - HUMAN-WORKSPACE-CONTRACT.md
freshness_checked_at: "2026-08-10T00:00:00Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: d8b24fe7cfef0986a4b48e7f4e6dd8c7373b451d4c54bde425a904889539b4d3
freshness_valid_until: "2027-08-10T00:00:00Z"
---
# Spectacular v2 product contract

Human operability and machine rigor are one behavior: navigation labels may be
readable, while durable identity and revision proof remain exact.

## Working model

Proposals are optional places to explore a problem, alternatives, questions,
and draft specifications. They may live in Spectacular, an issue, or the
current preparation conversation. A Proposal is neither authority nor current
product truth.

A Mission is the frozen execution plan. Its `MISSION.md` owns the selected
outcome, completion criteria, inline Objectives, current Run, execution
boundaries, and activation. Additional Objective and Run files exist only when
their independent detail or lifecycle earns the extra file.

Product and Capability Contracts state accepted behavior. When a Mission
changes that behavior, editing the relevant specification is part of the same
work as editing code and tests. One Mission completion gate closes the work;
there is no separate reconciliation step for the owner to remember.
