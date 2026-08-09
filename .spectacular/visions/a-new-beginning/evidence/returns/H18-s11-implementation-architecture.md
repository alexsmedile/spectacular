---
type: handoff-return
schema_version: spectacular.handoff-return.v2
handoff_id: H18
session: S11
status: complete
central_disposition: accepted
baseline_commit: d36bc5bb14367b0566bcda11a8835010dbb942fa
baseline_tree: 55643587b23c12015cb18d0b45b668daa3855571
baseline_dirty: false
date: 2026-08-09
---

# H18 return — S11 implementation architecture

## Central disposition

**Accepted.** H18 verified its issued baseline and all thirteen immutable inputs, completed each
of the five owner-decision clusters, and satisfied S11's exit gate. The authoritative result is
[`../../IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md`](../../IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md).

This acceptance changes the refactor program's architecture contract only. It authorizes neither
implementation, v1 migration, retirement, nor S12A by itself.

## Verified gate

The delegated checkout was detached and clean at dispatch commit
`d36bc5bb14367b0566bcda11a8835010dbb942fa`, tree
`55643587b23c12015cb18d0b45b668daa3855571`. The H18 handoff hash was
`6691471ce74dc31352c7ebb8ec1d9e531912814d656a9124c92e6d578b4f21a4`.
All thirteen accepted contracts matched both their supplied hashes and sidecars. H18 made no
files, branches, commits, migration, lifecycle, or external-system changes.

## Owner dispositions

| Cluster | Accepted disposition |
|---|---|
| A | A2-DW: clean Go core; checksummed macOS/Linux native binaries; Windows deferred without lockout |
| B | B1: domain-centered modular executable |
| C | Markdown authority; deterministic in-memory index; narrow on-demand projections; optional earned modular cache |
| D | D1: isolated in-repository Go migration capsule with a separate build graph |
| E | E1: layered, risk-based test and proof architecture |

## Result and retained boundaries

- Go owns deterministic v2 semantics; the Skill retains judgment; host runtimes execute; native
  providers perform/attest effects.
- Canonical Markdown remains authoritative. Indexes, projections, and any cache remain disposable.
- One registry derives mechanical dispatch, help, effect classification, and machine interfaces.
- Guardrails select owner prose mechanically; they cannot mint authority or weaken invariants.
- The migration capsule is a removable v1-only boundary; v2 imports no legacy reader, aliases,
  fallback behavior, or migration dependency.
- The proof order includes cold recovery, bounded work formation, reconciliation/resume, and
  migration rollback/capsule-removal evidence.

## Open Type-2 work

Exact implementation libraries/layout, locking, serialization, cache mechanics, adapter protocol,
CI and benchmark thresholds, Windows delivery, and release-signing details remain deferred.
S12A may now be considered only through its own authorization and approval gate.
