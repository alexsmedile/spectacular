---
type: source-card
source: source-009
provided_as: source9
received: 2026-08-07
authority: unsourced-expanded-synthesis
status: ingested
scope: [graph-orchestration, handoffs, harness, ontology, retrieval, memory, modules, state, alignment]
duplicate_sections: [source-008-substantial, source-007-selected, source-006-selected]
completeness: substantial
---

# Source 009 — Expanded graph engineering synthesis

## Thesis

Clarify the graph-oriented agent architecture from Source 008: plan work as an
explicit split/work/check/merge graph, exchange bounded typed payloads between
isolated nodes, ground reasoning in an ontology and structural graph, compress
context through hierarchical state and pointers, enforce software quality through
a deterministic harness, and align architecture before autonomous execution.

## Relationship to Source 008

Source 009 substantially restates Source 008. It does not independently corroborate
the Diamond pattern, graph memory, context rot, deep modules, autonomous error
feedback, resumable state, shift-left review, Wayfinder, or prototypes. Its useful
clarifications are:

- graph edges are proposed as validated handoff contracts, not merely dependencies;
- OWL/RDF are named as possible semantic technologies;
- the memory hierarchy is defined as raw, atoms, scenes, and persona;
- diagrams or node IDs are proposed as a progressive navigation layer;
- the harness-versus-model claim receives concrete but unattributed numbers;
- prototype code is promoted from evidence to “definitive spec,” exposing an
  important authority collision.

## Claim audit

| Claim | Assessment | Correction or boundary |
|---|---|---|
| Harness changes yield up to 48 points while model swaps yield about 5% | directionally supported, exact formulation unattributed | Recent benchmark work supports evaluating model and harness together and reports large harness-dependent variation, but effects vary by task. Harness-Bench reports substantial configuration variation; Agents' Last Exam reports a case where model choice moved results more than harness choice: [Harness-Bench](https://arxiv.org/abs/2605.27922), [Agents' Last Exam analysis](https://agents-last-exam.org/blogs/harness-matters). |
| A Skeptic node “holds no bias” | refuted | Context isolation can reduce anchoring but does not remove shared-model bias, rubric defects, or correlated failure. |
| OWL/RDF create deterministic business-rule enforcement | incomplete technology mapping | RDF represents graphs and OWL defines ontology semantics and inference. Constraint validation is more directly the role of SHACL or application checks; W3C notes assumptions and missing axioms can cause missed or false violations: [SHACL](https://www.w3.org/TR/shacl/), [OWL 2 semantics](https://www.w3.org/TR/owl2-rdf-based-semantics/). |
| Vector retrieval fails on topology questions | overstated | Similarity alone does not encode authoritative topology, but retrieved source text may still answer the question. Explicit edges make traversal and explanation more deterministic when fresh. |
| Deep modules improve AI retrieval | plausible, context-dependent | Runtime modules should hide cohesive complexity behind narrow interfaces. Atomic evidence cards remain intentionally small records and are not runtime modules. |
| A working prototype becomes the definitive spec | rejected | A prototype is high-fidelity evidence about selected behavior; it is still incomplete on failure, security, compatibility, operations, and maintainability. Accepted contracts remain authoritative. |

## How these ideas could improve the refactor method

1. Add a **derived visual map** above the atomic card database: Mermaid nodes use
   stable PZL IDs, show only decision-relevant dependencies/collisions, and link to a
   Markdown table of full cards. The map is a projection, never new authority.
2. Treat every future worker/reviewer handoff as a **typed packet**: objective,
   permitted inputs, authoritative sources, exclusions, output schema, evidence,
   confidence limits, and unresolved questions.
3. Evaluate **model plus harness configuration** on the same Spectacular tasks before
   solving reliability problems by model escalation. Record task success, retrieval
   volume, gate failures, retries, latency, and cost.
4. Use a **small semantic constraint layer** first: frontmatter schema, IDs, typed
   relationships, and deterministic validators. Do not introduce RDF/OWL/SHACL until a
   concrete multi-hop or validation problem defeats the simpler representation.
5. Apply **deep-module criteria to implementation**, not to every information record:
   a module must own one cohesive responsibility behind a small stable interface.
6. Use prototypes for unresolved CLI/UX or lifecycle behavior, then translate observed
   results and rejected behavior into an approved contract. Never promote prototype code
   merely because it is interactive.
7. Make architecture review proportional: settle execution authority, truth authority,
   portable-core boundaries, and state ownership before implementation; defer file-level
   predictions until the run inspects current code.

## Application questions

- Which 10–20 PZL nodes form the decision-critical graph, and which should remain
  searchable cards outside the visual map?
- Is the first graph artifact a human navigation projection, an agent retrieval index,
  or an executable scheduler? The current recommendation is navigation projection only.
- What exact schema must cross a worker-to-reviewer edge without passing conversation history?
- Which three representative Spectacular tasks can compare harness variants fairly?
- Which invariants need deterministic validation now, and which require human or evidence review?
- Where would an interactive prototype materially answer a question that prose cannot?
- Which state is canonical enough to resume, and which evidence is local, sensitive, or disposable?

## Provisional assessment

**Strong additions:** typed handoff packets; model+harness evaluation; derived Mermaid
navigation; explicit distinction between semantic models and validators.

**Existing cards strengthened:** graph/attempt composition, context isolation, graph
retrieval, layered memory, deterministic harness, deep modules, prototypes, wayfinding,
independent verification, shift-left review, targeted repair, and resumable state.

**Rejected as written:** bias-free Skeptic; exact universal 5%/48-point comparison;
OWL/RDF alone as business-rule enforcement; vectors always fail multi-hop questions;
prototype code as definitive specification.

No graph runtime, RDF stack, parallel fleet, memory hierarchy, or prototype authority is
accepted by this ingestion.

## New concept pieces

- [PZL-119 — Typed inter-node handoff contract](concepts/PZL-119-typed-handoff-contract.md)
- [PZL-120 — Derived Mermaid decision map](concepts/PZL-120-derived-mermaid-map.md)
- [PZL-121 — Evaluate the model-harness system](concepts/PZL-121-model-harness-evaluation.md)
- [PZL-122 — Separate semantics, inference, and validation](concepts/PZL-122-semantics-inference-validation.md)

## Decision packets seeded

- Approve or reject a derived Mermaid decision map for the intake database.
- Define the minimum typed packet for planner, worker, reviewer, and merger handoffs.
- Select a small repeatable harness-evaluation suite before execution architecture changes.
- Decide whether simple Markdown schemas satisfy the semantic layer before evaluating
  semantic-web technologies.
