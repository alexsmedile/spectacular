---
type: mission-repair-delta
mission: scenario-a-cold-recovery
version: 1.0
status: authorized
authorized_by: central-orchestration
authorized_at: 2026-08-10
returned_head: 816c34aec885747c5f4a7d40274a35b34f91becf
returned_tree: 5a79f91dc98949292d8cdef20005017d4a52b347
branch: codex/feat/v2-scenario-a-cold-recovery
repair_rounds_authorized: 1
topology: primary-owner-plus-one-fresh-read-only-reviewer
next_action: implement-one-cohesive-scenario-a-repair-round
---

# Scenario A authorized repair delta

## Central disposition and continuity

Central orchestration bounced Scenario A. The existing Mission remains active but unaccepted;
Scenario B remains unauthorized. Before mutation, the primary agent revalidated the clean returned
head `816c34aec885747c5f4a7d40274a35b34f91becf`, tree
`5a79f91dc98949292d8cdef20005017d4a52b347`, on branch
`codex/feat/v2-scenario-a-cold-recovery`.

This evidence extends the exhausted charter repair budget by exactly one cohesive,
hypothesis-changing round. It does not replace or broaden the Scenario A Mission charter.

## Authorized repair envelope

1. Validate matched known-command arguments through the registry before workspace discovery so
   malformed invocations return usage exit `2` even outside a workspace.
2. Refuse a symlinked `.spectacular/workspace.yaml`, a symlinked `.spectacular` metadata directory,
   and equivalent authority-marker escapes before reading marker bytes. Prove the boundary with
   real filesystem tests.
3. Make Proposal/current-truth drill-down mechanically honest without adding a Proposal command.
   Proposal record metadata must retain exact noun, typed ref, canonical path, and fingerprint;
   `show_command` is omitted because the Scenario A registry has no Proposal operation.

The exact ten-command public surface, schemas, exit classes `0/2/3`, read-only invariant,
dependencies, baseline, prohibited paths/effects, and B+C non-expanding join remain unchanged.
No v1, provider, cache, compatibility, lifecycle, Scenario B, or Scenario C behavior is authorized.

## Hypothesis, effects, recovery, and proof

The three P1 findings share one cause: authority was being consumed before its interface boundary
was mechanically proven. The repair moves command argument authority ahead of discovery, proves
manifest containment ahead of parsing, and derives pointer commands only for registry-supported
nouns. Effects remain local owned-path edits, tests, evidence, and coherent commits. Recovery is the
clean returned head above plus the last green repair commit; no reset, stash, or destructive cleanup
is authorized.

Required proof is the three focused regressions; existing Scenario A focused, adversarial,
nonmutation, deterministic cold, and self-hosted exercises; the complete Go format/module/vet/test/
race/build matrix; Bash syntax and version guard; the full v1 suite; exact branch/head/tree and
baseline scope-diff verification; valid evidence sidecars; and one fresh independent reviewer who
does not edit. The reviewer may report findings only. Ordinary in-envelope findings may be repaired
only if this single round still covers them; otherwise the return is blocked.

Design Sufficiency remains `sufficient` and Slice Quality remains `coherent`: the delta resolves
three implementation defects without selecting new product semantics or changing the slice.

## Terminal action

Return one updated `spectacular.handoff-return.v2` packet to central orchestration for
`accept | bounce | escalate`. Do not claim Scenario A accepted or authorize Scenario B.
