---
type: mission-design-evidence
mission: scenario-a-cold-recovery
version: 1.0
status: sufficient-and-coherent
author: scenario-a-investigator
observed_at: 2026-08-10
baseline_commit: 8d5722c
next_action: implement-locked-slice
---

# Scenario A implementation design evidence

The read-only Investigator inspected the accepted program, Scenario sequencing, M1 closure and
APIs, relevant immutable contracts, owner dispositions, and active charter. It made no file or Git
mutation and exposed no new product or architecture choice.

## Gate revalidation

- Design Sufficiency: `sufficient`
- Slice Quality: `coherent`
- Blockers: none

Exact record and JSON field layouts are authorized Type-2 implementation pins under the accepted
interface and execution defaults. They must be locked in typed Go structures and golden tests
before relying on their behavior.

## Directional recovery graph

Expected-fingerprint links are directional so canonical hashes cannot form an authority cycle:

```text
Mission --current_run + expected_run_fingerprint--> Run
Run --latest_checkpoint + expected_checkpoint_fingerprint--> Checkpoint
Decision --mission + expected_mission_fingerprint + operation + target--> Mission chain
Run --mission identity only--> Mission
Gap --scope identity only--> owning scope
```

A Mission must not store a Decision backlink that would alter the fingerprint the Decision
authorizes. Freshness is `current` only when the expected fingerprint matches normalized canonical
content, `stale` on mismatch, and `unknown` when no usable assertion exists.

## Implementation boundary

- The CLI executable remains a thin exit/stdout/stderr adapter.
- One command registry owns the ten accepted forms and their read-only effects/schema identifiers.
- Discovery selects the nearest explicit v2 marker, scans declared roots only, and rejects lexical
  or symlink escape before parsing records.
- Domain changes add static read grammars/refusals only; they do not apply transitions or infer
  lifecycle authority.
- Index retains M1 identity-first behavior and adds a separate path-first CLI view.
- Projection uses typed output structures and an injected clock.
- No Scenario A call path reaches persistence.

## Deterministic refusal order

1. marker, schema, and path-safety failure;
2. parse/schema failure in sorted canonical-path order;
3. duplicate identity/path;
4. relationship, freshness, and authority conflict;
5. lookup, noun, and scope mismatch.

Missing authorization on an otherwise valid chain yields a successful exact owner gate. A
structurally missing required source refuses.

## Primary risks retained for proof

- M1 index ordering must not be silently changed.
- New record grammars must preserve every M1 test and must not imply lifecycle operations.
- Symlink containment requires resolved-path validation, not lexical checks alone.
- Production generation time is real; deterministic golden output uses a fixed injected clock.
- Performance is measured on the built non-race binary; functional behavior is separately proven
  under the race detector.
