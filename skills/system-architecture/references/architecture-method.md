# Architecture Method

Use this when: exploring consequential architecture options or performing a complete design, modernization, or evidence-based review that needs detailed evidence prompts.

SKILL.md owns route selection, workflow order, and completion criteria. Use this file as the working evidence worksheet for the selected steps.

## Contents

- [1. Frame the problem](#1-frame-the-problem)
- [2. Model domain and ownership](#2-model-domain-and-ownership)
- [3. Compare system shapes](#3-compare-system-shapes)
- [4. Trace runtime behavior](#4-trace-runtime-behavior)
- [5. Design data ownership](#5-design-data-ownership)
- [6. Design operations and security](#6-design-operations-and-security)
- [7. Record decisions and uncertainty](#7-record-decisions-and-uncertainty)
- [8. Plan evolution](#8-plan-evolution)
- [Architecture review lens](#architecture-review-lens)

## 1. Frame the problem

Capture:

- Business outcomes and success measures
- Actors, workflows, and external systems
- System boundary and explicit non-goals
- Current-state assets and constraints
- Regulatory, residency, privacy, and security obligations
- Delivery timeline, budget, team topology, and operational maturity
- Evidence sources and unresolved questions

For an existing system, inspect source code, runtime manifests, database migrations, schemas, API definitions, telemetry, incident history, and ADRs. Do not infer the deployed architecture from directory names alone.

Express important quality attributes as scenarios:

| Field | Meaning | Example |
|---|---|---|
| Source | Who or what causes the event | Authenticated customer |
| Stimulus | Event to handle | Submits an order |
| Environment | Conditions | Peak sale; one zone unavailable |
| Artifact | Affected system part | Checkout path |
| Response | Required behavior | Accept once or fail clearly; preserve cart |
| Measure | Testable threshold | p99 under 800 ms; no duplicate order |

Quantify at least the drivers likely to shape the design:

- Traffic: average, peak, burst, read/write mix, payload size, growth
- Performance: latency percentiles, throughput, concurrency, freshness
- Reliability: SLO, error budget, RTO, RPO, durability, degraded modes
- Security: threat actors, sensitive assets, trust boundaries, audit needs
- Changeability: release frequency, compatibility window, ownership
- Cost: target range, major scaling dimensions, fixed versus variable costs

Record the source and confidence of each number. Use ranges when evidence is weak.

## 2. Model domain and ownership

Identify business capabilities, language, entities, invariants, commands, events, and ownership. Distinguish:

- **Bounded context:** a consistency boundary for a domain model and language
- **Module:** a code organization and encapsulation boundary
- **Container:** an independently runnable or deployable unit in C4
- **Service:** a network-accessible capability with an operational lifecycle
- **Team boundary:** a durable ownership and communication boundary

Do not make these boundaries identical without a reason. Prefer high cohesion, low coupling, explicit contracts, and ownership aligned with the ability to change and operate the capability.

## Module boundary review

Use this lens for modularization or module reviews; scale it to the changed
capabilities. Inspect imports, call paths, data access, and transactions before
judging folder layout.

| Boundary fact | Evidence to capture | Verification |
|---|---|---|
| Responsibility and ownership | Reason to change; authoritative writer | Changes remain local; writes go through the owner |
| Public surface | Minimal operations, IDs, detached read views | Contract compatibility; consumers cannot mutate canonical state |
| Dependency direction | Allowed imports/calls; hidden implementation | Negative import checks; acyclic production graph |
| Policy versus effects | Pure rules; consumer ports; concrete adapters | Rules tested without filesystem, database, network, or processes |
| Transaction integrity | Invariants requiring one coordinated commit | Failure injection and rollback preserve those invariants |

Classify each important boundary as **declared**, **implemented**, or **verified**,
citing the enforcement location and command for the last category. Language
visibility rules may protect external consumers while permitting sibling imports;
verify their actual scope. Check platform-specific sources where relevant.

Extract a read model through the existing decoder before duplicating parsing.
Copy nested mutable data or use immutable values. Introduce ports for actual
consumer needs, avoiding interfaces that merely mirror every implementation.
Keep one composition point for concrete adapters. Extract pure rules separately
from persistence while retaining one coordinator for coupled writes.

Compare a coherent package, a modular monolith, and distribution only when their
drivers apply. Do not create one module per entity, a generic shared-business
kernel, or an event bus solely to reduce imports. An invariant requiring immediate
consistency can justify a shared transaction boundary.

Pass the ownership/consistency facts to `data-modeling`, public operation semantics
to `action-contracts`, and executable expectations to `test-sentinel`. Load
`event-spine` only for durable asynchronous delivery. Follow existing project
artifacts rather than creating a second boundary catalog.

## 3. Compare system shapes

Record the credible candidates required by workflow step 3. When selecting patterns, load only the relevant section of [decision-guide.md](decision-guide.md).

Compare them against the same drivers:

- Domain fit and transaction integrity
- Independent change and deployment
- Scaling profile and isolation needs
- Availability and failure propagation
- Security and compliance isolation
- Team cognitive load and operational readiness
- Migration complexity and reversibility
- Cost now and at the expected scale

When the request authorizes a decision, record the preferred shape, rejected shapes, evidence, and conditions that would reverse the choice. During non-decisional exploration, record the leading candidates and the evidence needed to choose without naming a winner.

## 4. Trace runtime behavior

Trace the happy path and at least these failures where applicable:

- Client retry or duplicate submission
- Downstream timeout, throttling, or partial outage
- Concurrent update or stale read
- Message duplication, delay, reordering, or poison payload
- Process crash between state change and event publication
- Lost cache, failover, deployment rollback, or schema mismatch

For each boundary, define:

- Contract and versioning policy
- Authentication and authorization
- Timeout budget and cancellation
- Retry owner, limit, backoff, and jitter
- Idempotency scope and retention
- Ordering and delivery guarantees
- Backpressure and overload behavior
- Observability and correlation identifiers
- Fallback, compensation, reconciliation, and operator action

Budget end-to-end latency rather than assigning each dependency the full request timeout. Avoid layered retries that amplify load.

## 5. Design data ownership

For each data set, record:

- Authoritative owner and writers
- Primary keys, access patterns, and query shapes
- Invariants and transaction boundary
- Consistency and staleness tolerance
- Volume, velocity, retention, and lifecycle
- Classification, residency, encryption, and audit requirements
- Backup, restore, RPO, and RTO
- Replicas, indexes, caches, search projections, and analytical copies
- Schema evolution and migration strategy

Keep a single authoritative source for each fact. Treat caches, search indexes, warehouses, and read models as derived data with explicit rebuild or reconciliation paths.

## 6. Design operations and security

Cover only what materially shapes the architecture:

- Edge, ingress, network zones, egress, and trust boundaries
- Workload identity, human identity, authorization, and secrets
- Compute topology, autoscaling signals, quotas, and capacity headroom
- Availability zones or regions and dependency failure domains
- Deployment, progressive delivery, compatibility, and rollback
- Logs, metrics, traces, audit events, alerting, and runbooks
- Backup verification, disaster recovery exercises, and dependency recovery
- Software supply chain, patching, and vulnerability response

State whether multi-region operation is active-active, active-passive, or recovery-only. Explain data conflict and failover behavior.

## 7. Record decisions and uncertainty

Turn uncertainty into validation work:

- Prototype risky integration or consistency behavior
- Benchmark the suspected bottleneck with representative data
- Load-test critical paths and overload controls
- Run failure injection or game days for recovery assumptions
- Threat-model sensitive flows
- Test backup restoration and regional recovery
- Validate migration and rollback on production-like data

Bound each experiment by a question, representative evidence, effort limit and
decision. End with selection plus integration, rejection, or reasoned deferral
and revisit condition. Extend only for identified missing evidence; successful
synthetic tests do not establish real integration or user acceptance.

For each consequential decision, record status, context, ranked drivers, credible options, choice, consequences, confidence, reversal conditions, and validation evidence. Give each material assumption or open question an owner or next action.

## 8. Plan evolution

Sequence implementation as thin, observable slices. Prefer expand-and-contract migrations, dual-read or dual-write only with reconciliation, compatibility windows, and explicit exit criteria. Include rollback for every material migration step.

For product delivery, identify the shortest end-to-end workflow that can be
used with a real dependency, the minimum integrated batch that enables it,
and its closure evidence. Show which later capabilities can wait. Keep chosen
constraints; recommend Docker, OS isolation or distribution only for a stated
need, never as inferred startup gates. Disclose unenforced boundaries and do
independent local work while live-access authorization remains pending.

## Architecture review lens

When reviewing, report findings by severity and evidence:

- **Critical:** likely security, integrity, or systemic availability failure
- **High:** major unmet driver or difficult-to-reverse design risk
- **Medium:** meaningful operability, evolvability, or cost concern
- **Low:** clarity, maintainability, or documentation improvement

For each finding, state the affected scenario, evidence, consequence, recommendation, and validation method. Do not manufacture findings to fill every category.
