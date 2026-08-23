---
type: Contract
id: 0199b000-0000-7000-8000-000000000001
human_ref: CC-axmej5
title: Capability Contract
status: current
created_by: owner
created: "2026-08-01T00:00:00Z"
updated: "2026-08-01T00:00:00Z"
accepted_proposal: Proposal:0199b000-0000-7000-8000-000000000099
applies_when:
  - A governed change is requested.
authority_freshness: Exact owner Decision and current fingerprint are required.
conformance_checks:
  - Current truth remains recoverable.
contract_version: "1"
does_not_apply_when:
  - The request is read-only.
does_not_provide:
  - Provider authority.
freshness_checked_at: "2026-08-10T10:00:00Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: 89a199584ff24912c402e7bd47914b88f4bda0237ab040734d90631c2d9965d1
freshness_valid_until: "2026-12-31T23:59:59Z"
gaps:
    - ref: smoke-gap
      problem: the installed release loop has no proven amendment path
      blocked_on: an owner decision recorded through the declaring Mission
operating_cases:
  - Existing capability evolution.
outcome: A durable current capability truth exists.
persistent_information:
  - Accepted Proposal and Decision provenance.
purpose: Govern capability evolution.
related_material:
  - Scenario B+C charter.
required_behavior:
  - Preserve Scenario A recovery.
---
# Capability Contract

This is current truth. A Proposal is only a base-bound candidate delta against it.
