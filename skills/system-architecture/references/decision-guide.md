# Architecture Decision Guide

Use this when: comparing system shape, communication, consistency, storage, caching, reliability, or security patterns for a consequential architecture choice.

SKILL.md owns the decision workflow and completion criteria. Use these heuristics only for the pattern categories in scope.

## Contents

- [System shape](#system-shape)
- [Communication](#communication)
- [Consistency and transactions](#consistency-and-transactions)
- [Storage selection](#storage-selection)
- [Caching](#caching)
- [Reliability](#reliability)
- [Security and privacy](#security-and-privacy)

## System shape

### Prefer a modular monolith when

- Domain boundaries are still evolving.
- One team or a small number of teams owns the product.
- Strong transactions span several capabilities.
- Independent scaling or deployment is not a demonstrated need.
- Operational simplicity and fast feedback dominate.

Enforce module APIs, dependency rules, data ownership, and observability so modules can evolve independently.

### Introduce services when

- A capability needs independent release cadence, scaling, isolation, security, or technology.
- Ownership boundaries are durable and teams can operate services end to end.
- The value of distribution exceeds network, consistency, testing, deployment, and incident costs.

Prefer coarse, cohesive services. Avoid a shared database that lets services bypass their contracts. Extract along a proven seam, not a speculative entity-per-service map.

## Communication

| Need | Prefer | Watch for |
|---|---|---|
| Immediate result or validation | Synchronous HTTP/gRPC | Latency coupling, cascading failures, versioning |
| Decoupled reaction or fan-out | Events | Duplicate/out-of-order delivery, schema evolution, observability |
| Durable work distribution | Queue/command message | Poison messages, backpressure, ownership |
| Bulk transfer or analytics | Files, streams, or change feeds | Freshness, replay, privacy, reconciliation |

Use events to state facts that occurred; use commands to request one owner to act. Do not hide synchronous business dependencies behind a broker. Define delivery semantics in terms of business outcomes; most practical systems use at-least-once delivery with idempotent consumers.

Use a transactional outbox or equivalent when a database update and event publication must be atomic. Define inbox/deduplication scope, replay behavior, dead-letter handling, and repair tooling. Use sagas only when a business process genuinely spans independent transaction owners; define compensations and irreversible steps.

## Consistency and transactions

Start from invariants:

- Keep data in one transaction boundary when an invariant must hold immediately.
- Use optimistic concurrency for low-contention collaborative changes.
- Use pessimistic locking or serialized ownership for short, high-value contention.
- Use reservations with expiry for scarce resources that span workflows.
- Use eventual consistency only with an explicit staleness window and user-visible behavior.
- Reconcile when atomic coordination is unavailable or too costly.

Discuss partition behavior separately from normal-operation latency. State what accepts writes, what becomes stale, how conflicts resolve, and how the system converges.

## Storage selection

Choose storage after documenting access patterns and invariants.

| Capability | Good fit | Poor fit / caution |
|---|---|---|
| Relational database | Transactions, constraints, joins, mature operations | Unbounded horizontal write scale without partition design |
| Key-value store | Known-key access, extreme scale, simple state | Ad hoc queries, multi-record invariants |
| Document database | Aggregate-shaped documents, flexible fields | Cross-document transactions and joins as a dominant pattern |
| Wide-column store | Large partitioned write workloads, predictable queries | Flexible querying and small operational teams |
| Graph database | Deep, changing relationship traversal | Ordinary CRUD or shallow fixed joins |
| Search engine | Full text, relevance, faceting | System of record or transactional integrity |
| Columnar analytics | Aggregation over large historical data | Point transactions and low-latency OLTP |
| Object storage | Durable blobs, archives, data lake | Fine-grained transactional updates |
| In-memory cache | Recomputable hot data, rate limits, coordination with care | Authoritative durable state |

Do not use polyglot persistence as a goal. Add a datastore only when its benefit exceeds expertise, backup, security, monitoring, migration, and incident-response costs.

For money or regulated value, prefer an append-only, balanced ledger model with stable identifiers, explicit corrections, auditability, reconciliation, and controlled posting. Do not infer that one database product creates financial correctness.

## Caching

Define the cache purpose, key, owner, source of truth, TTL, invalidation, stampede protection, consistency tolerance, and behavior on loss. Prefer cache-aside for simple read acceleration. Avoid caching until a measured bottleneck or availability need justifies it.

## Reliability

- Derive component targets from the end-to-end SLO and dependency chain.
- Use timeouts on remote calls and propagate cancellation.
- Retry only transient, safe, and idempotent operations with bounded exponential backoff and jitter.
- Use circuit breaking or concurrency limits to prevent resource exhaustion.
- Apply backpressure, admission control, quotas, and load shedding before saturation.
- Isolate high-risk workloads with pools, queues, cells, or bulkheads when justified.
- Design degraded modes that preserve the highest-value functions.
- Test recovery; replication is not backup.

## Security and privacy

- Model assets, actors, entry points, data flows, and trust boundaries.
- Authenticate workloads and users; authorize at the resource or capability boundary.
- Minimize privileges, retained data, exposed interfaces, and shared credentials.
- Encrypt sensitive data in transit and at rest; define key ownership and rotation.
- Protect administrative paths and record tamper-resistant audit events.
- Define tenant isolation, deletion, residency, retention, and incident response.
- Treat third-party services and build pipelines as dependencies with failure and compromise modes.
