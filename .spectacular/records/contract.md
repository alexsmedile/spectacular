---
type: Contract
id: 019fe381-5d61-7223-b362-03a5f99a7b10
title: Spectacular v2 public product contract
status: current
created_by: owner
created: "2026-08-10T00:00:00Z"
updated: "2026-08-10T00:00:00Z"
accepted_proposal: Proposal:019fe381-5d61-7223-b362-03a5f99a7b01
applies_when:
    - A public Spectacular v2 capability or release action is considered.
authority_freshness: Exact owner Decision and current source fingerprints are required.
conformance_checks:
    - Pointer-first retrieval remains executable and read-only.
    - Release artifacts and installer checks pass without a user-side Go toolchain.
contract_version: 2.0.0-rc.1
does_not_apply_when:
    - A request concerns the frozen v1 product surface.
does_not_provide:
    - Provider authorization or publication permission.
freshness_checked_at: "2026-08-10T00:00:00Z"
freshness_source: .spectacular/workspace.yaml
freshness_source_fingerprint: 5b9af18ac4b8e0f891306a59e3d6e4d81c8b0669ab220847a22b3267d778d1f0
freshness_valid_until: "2027-08-10T00:00:00Z"
operating_cases:
    - Governed v2 work and release readiness.
outcome: One unambiguous v2 public product surface is recoverable and release-ready.
persistent_information:
    - Accepted contract, release evidence, recovery pointers, and the next owner gate.
purpose: Govern the Spectacular v2 public product surface.
related_material:
    - RECOVERY.md
    - commit:1812a2578d2cf2a873d80e235734e38e56fd091e
required_behavior:
    - Preserve root-only v2 identity and explicit provider boundaries.
---
# Spectacular v2 public product contract

This Contract captures the accepted root cutover and the RC release boundary.
