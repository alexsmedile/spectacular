---
schema_version: spectacular.project-anchor.v1
title: Event-Sourced Transaction Ledger Engine
boundaries:
  in_scope:
    - Event sourcing with append-only ledger in events.jsonl
    - CLI commands for deposit, transfer, balance, and reconcile
    - Idempotency via transaction ID (--tx-id) tracking
    - Lossless balance replay and reconstruction
  out_of_scope:
    - Remote distributed databases or cloud APIs
constraints:
  - Must pass tests/check.sh deterministically
---
# Project Anchor
