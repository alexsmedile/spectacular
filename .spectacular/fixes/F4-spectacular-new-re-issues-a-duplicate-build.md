---
type: fix
opened: 2026-07-11
verified: 2026-07-11
severity: medium
from_audit: null
debug_job: null
signature: new duplicate build id, last_build drift, ledger max reconciliation, bN collision
related: []
---

# F4 — spectacular new re-issues a duplicate build id when config.last_build drifts behind the roadmap ledger

## Problem
cmd_new computed _next_build = config.last_build + 1, reading only config.yaml and never the roadmap ledger. When a request was slotted into the ledger without last_build catching up, new re-issued an id the ledger already used (observed live: a duplicate b25).

## Intended behavior
new assigns a build id strictly greater than every id already in use, self-healing against counter/ledger drift regardless of how the drift arose.

## Root cause
The build-id counter was a single source (config.last_build) that silently trusted itself against a separately hand-editable ledger. Derived-counter-vs-source-of-truth drift with no reconciliation.

## Fix
Before computing _next_build, grep the roadmap ledger (roadmaps/index.md, legacy ROADMAP.md fallback) for the max bN and reconcile: _last_build = max(config.last_build, ledger_max). No-drift path unchanged.

## Success criteria
In a workspace where config.last_build=5 and the ledger has b9, new assigns b10 (not b6) and writes last_build: 10 back.

## Verified by
tests/cli/mutator.test.sh scenario 8b (drift → b10 not b6) + full suite 16/16

## Signature
new duplicate build id, last_build drift, ledger max reconciliation, bN collision
