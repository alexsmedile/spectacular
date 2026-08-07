---
type: source-card
source: source-008
provided_as: source8
received: 2026-08-07
authority: unsourced-protocol-synthesis
status: ingested
scope: [orchestration, retrieval, memory, harness, verification, infrastructure, state]
duplicate_sections: [source-003-selected, source-006-selected, source-007-expanded]
completeness: substantial
---

# Source 008 — Graph, harness, and state protocols

## Thesis

Replace unconstrained agent loops with bounded dependency graphs, isolated work and
verification contexts, deterministic engineering gates, explicit memory structures,
and crash-resumable external state. Use graph-shaped knowledge and a consolidated
PostgreSQL substrate to support that operating model.

## Source integrity

This is a high-density synthesis, not a cited research report. It provides no author,
date, benchmark, project sample, failure data, or provenance for its named frameworks.
Density improves review efficiency but does not increase authority. Claims phrased as
universal comparisons or infrastructure facts were therefore checked separately.

## Claim audit

| Claim | Assessment | Correction or boundary |
|---|---|---|
| Graph engineering should replace loop engineering | useful contrast, false dichotomy | DAGs model dependency and fan-out; execution, retry, and repair still require bounded state transitions or loops. |
| GraphRAG is better than vector RAG | overstated | Microsoft's GraphRAG work reports gains for a class of global sensemaking questions over large corpora, not universal superiority: [research paper](https://www.microsoft.com/en-us/research/publication/from-local-to-global-a-graph-rag-approach-to-query-focused-summarization/). |
| An ontology creates an infallible world model | refuted | An ontology can make accepted semantics explicit, but extraction, modeling, freshness, and implementation drift remain fallible. |
| Independent verifier agents should always grade work | directionally useful | Independent context reduces some self-review bias; it does not create independent evidence or justify universal parallel/high-tier-model cost. |
| PostgreSQL natively provides the whole proposed stack | mixed and time-sensitive | JSON/JSONB are core PostgreSQL types; `SKIP LOCKED` is suitable for queue-like access but exposes an inconsistent view. `pgvector` is an extension, not core PostgreSQL. Property graphs appear in PostgreSQL 19 documentation, while current stable PostgreSQL 18 does not expose that feature: [JSON types](https://www.postgresql.org/docs/current/datatype-json.html), [`SKIP LOCKED`](https://www.postgresql.org/docs/current/sql-select.html), [pgvector](https://github.com/pgvector/pgvector), [PostgreSQL 19 property graphs](https://www.postgresql.org/docs/19/ddl-property-graphs.html). |

## Durable contributions

- Model complex work as explicit dependencies and joins, while retaining bounded
  attempt loops for execution and recovery.
- Split independent work, check it against declared evidence, then merge only the
  accepted outputs.
- Give specialized nodes bounded context and a shared contract rather than one
  accumulated conversation.
- Store durable state and large evidence outside the active prompt; retrieve details
  through stable pointers.
- Surround nondeterministic agent judgment with deterministic tests, parsers, policy
  checks, and explicit state transitions.
- Treat prototypes as evidence for uncertain behavior, not as substitutes for accepted
  contracts or production quality.
- Separate creation from verification and calibrate human gates to consequence.
- Make crash recovery a state-model requirement rather than a chat-history feature.

## Assumptions and collisions

1. Source 006 explicitly defers multi-agent orchestration and semantic/vector retrieval;
   Source 008 treats both as core architecture.
2. Source 004 recommends retiring Wayfinder; Source 008 recommends dependency
   wayfinding before persistent execution without providing Spectacular usage evidence.
3. Source 007 already proposes graph-first retrieval and multi-run assurance. Source 008
   strengthens the rationale but does not validate graph adoption Levels 2 or 3.
4. Clean contexts need a compact shared contract; isolation without it can omit safety-
   critical information.
5. Raw logs may contain secrets, personal data, huge output, or transient paths. External
   storage needs redaction, retention, and ignore rules; Git-tracked Markdown is not a
   universal answer.
6. A deterministic gate can verify a declared property, not all semantic correctness.
7. A database choice is an implementation option for a server-backed product, not yet a
   requirement for Spectacular's Markdown/Git-native distribution.

## Provisional assessment

**Strong:** explicit dependency structure; split/work/check/merge; bounded contexts;
deterministic harnesses; classical modularity; evidence-bearing prototypes; independent
review as a risk-triggered gate; risk-calibrated human approval; targeted repair loops;
external resumable state.

**Promising but unproven:** graph memory, an ontology layer, hierarchical memory, and
pre-execution wayfinding.

**Rejected as written:** graph replaces loops; GraphRAG universally beats vector
retrieval; an ontology is infallible; every result needs parallel high-tier verification;
PostgreSQL is a universal substrate; all raw output belongs in versioned Markdown.

No graph platform, multi-agent default, database, memory hierarchy, or verifier authority
is accepted by this ingestion.

## New concept pieces

- [PZL-103 — Dependency DAG plus bounded attempts](concepts/PZL-103-dependency-dag-bounded-attempts.md)
- [PZL-104 — Split, work, check, merge](concepts/PZL-104-split-work-check-merge.md)
- [PZL-105 — Isolated bounded node context](concepts/PZL-105-isolated-node-context.md)
- [PZL-106 — Graph memory for structural retrieval](concepts/PZL-106-graph-memory.md)
- [PZL-107 — Explicit semantic ontology](concepts/PZL-107-semantic-ontology.md)
- [PZL-108 — Pointer-first externalized evidence](concepts/PZL-108-pointer-first-evidence.md)
- [PZL-109 — Layered compounding memory](concepts/PZL-109-layered-memory.md)
- [PZL-110 — Deterministic agent harness](concepts/PZL-110-deterministic-agent-harness.md)
- [PZL-111 — Classical modularity as anti-slop constraint](concepts/PZL-111-classical-modularity.md)
- [PZL-112 — Prototype as uncertainty evidence](concepts/PZL-112-prototype-evidence.md)
- [PZL-113 — Dependency wayfinding before execution](concepts/PZL-113-dependency-wayfinding.md)
- [PZL-114 — Independent evidence review](concepts/PZL-114-independent-evidence-review.md)
- [PZL-115 — Risk-calibrated human gates](concepts/PZL-115-risk-calibrated-human-gates.md)
- [PZL-116 — Targeted error-feedback repair](concepts/PZL-116-targeted-error-feedback.md)
- [PZL-117 — PostgreSQL consolidation option](concepts/PZL-117-postgres-consolidation.md)
- [PZL-118 — Crash-resumable external state](concepts/PZL-118-crash-resumable-state.md)

## Decision packets seeded

- Which workflows need an explicit DAG, and which remain a single bounded attempt loop?
- Under what risk and independence conditions is verifier fan-out worth its cost?
- Is graph memory an accepted-contract projection, a retrieval index, or both?
- What durable state is canonical, and what raw evidence is local, redacted, or expired?
- Does Spectacular need any server substrate, or should database-backed capabilities stay
  outside its portable core?
