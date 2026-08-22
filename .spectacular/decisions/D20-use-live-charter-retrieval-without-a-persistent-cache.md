---
type: Decision
id: 01a02a64-c13a-7b47-ba86-9bcf426c0d61
title: Compile charters from live governance without a persistent cache
created_by: Alex
created: "2026-08-22T16:54:04Z"
updated: "2026-08-22T16:54:04Z"
actor: Alex
actor_role: owner
ref: D20-use-live-charter-retrieval-without-a-persistent-cache
question: What should charter compilation retrieve, preserve, and cache?
disposition: live-frontmatter-retrieval-and-temporary-charters
rationale: >-
    A quick live frontmatter sweep provides the same retrieval capability as a persistent preflight
    cache while remaining current and avoiding invalidation state. The charter is a temporary
    governance-compilation helper for the Orchestrator, not the complete assignment or a new source
    of truth. The frozen Handoff preserves what was actually assigned.
alternatives:
    - trusted persistent preflight state under the canonical governance tree
    - storing every compiled charter as a canonical record
    - allowing the CLI to infer assignment meaning and owner intent
authority_basis: Owner chose live frontmatter retrieval, temporary charter output, complete Orchestrator-authored Decisions and Handoffs, and benchmark-gated reconsideration of caching.
authorized_effects:
    - contract.version-bump
conditions:
    - live-governance-sweep
    - no-persistent-cache
    - charter-envelope-target-1200-max-1440
scope:
    - v2
targets:
    - Proposal:01a029be-b7d3-703c-a7ee-50c6b8bae3a2
supersedes: D14-preflight-discovery-caching
---

# Compile charters from live governance without a persistent cache

## Decision

- The Orchestrator selects relevant Decision refs by meaning from a quick live frontmatter sweep. The compiler consumes explicit refs and behaves deterministically.
- `charter` is read-only temporary output used to help write the complete Handoff. It is not stored as a canonical record.
- The governance envelope targets 1,200 tokens. At 1,201–1,400 it warns; at 1,401–1,440 it strongly recommends splitting and requires Orchestrator approval; above 1,440 it refuses.
- Safe compaction runs before warning or refusal and never removes claims, authority, scope, stops, or proof requirements. Sequential stages remain Runs; independently provable outcomes become separate Objectives.
- `decide` accepts a complete Orchestrator-drafted Decision package from a file or stdin, records it atomically, and refreshes generated indexes.
- `decide` reports newly eligible work only for explicit Decision blockers and never mutates Objective or Run state.
- Persistent caching adds no unique capability and does not ship in M15–M17. A later Mission may reconsider a disposable optimization only after benchmarks demonstrate a material bottleneck.
