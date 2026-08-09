---
type: program-amendment-decision
decision: mvp-scenario-cli-sequencing
version: 1.0
status: accepted
authority: owner
decided_at: 2026-08-09
owner_disposition: "OK CONTINUE"
source:
  - evidence/program-amendment-mvp-scenario-sequencing.md
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.1
m1_head: a488b2efe7828f59724f730b9a590b9a644e6885
m1_tree: 6cce1468c13e51ef007b0e23ed6b5295cdefd87b
next_action: integrate-m1-and-prepare-scenario-a
---

# MVP scenario and CLI sequencing decision

## Decision

Accept M1 at the exact independently reviewed head and tree above. Preserve its minimal substrate
boundary: official writes produce idempotent semantic normalization; Markdown body content and
supported tree-shaped YAML values survive semantically; YAML presentation details and comments are
not authoritative; anchors, aliases, shared graphs, and cyclic graphs refuse deterministically.

For the MVP, promise atomic replacement and process-crash safety only to the strength demonstrated
by tests. Do not claim power-loss durability until platform-specific directory synchronization and
recovery behavior are implemented and proved.

Replace the remaining horizontal implementation sequence with four vertical outcomes:

```text
A — cold recovery
→ B — fuzzy intent to bounded governed work
→ C — evidence, disposition, reconciliation, and cold resume
→ R — release and distribution hardening
```

Scenario C remains mandatory for MVP because evidence-backed acceptance, reconciliation, and
durable continuation are constitutional product behavior.

## Execution controls

Retire fixed commit-count ceilings. Use exact reviewed commit/tree, owned-path scope, prohibited
effects, changed-invariant inventory, dependency diff, clean target, and bounded
hypothesis-changing repair accounting. Commits remain coherent review units, not a success metric.

Operate from two live projections: the current executable program and the current Mission charter.
Immutable constitutional contracts remain authoritative referenced inputs rather than duplicated
live checklists.

## Mechanical CLI sequence

Scenario A begins with:

```text
spectacular anchor show project [--json]
spectacular mission list [--json]
spectacular mission show <ref> [--json]
```

Consequential pointers drill down through noun-first `gap`, `run`, `checkpoint`, `evidence`, and
`decision` `list`/`show` operations plus scoped `workspace validate`. One internal generic
record/projection engine may serve these commands, but `Record` is not a primary public noun.

Scenario B adds Proposal and Handoff inspection/validation, then confirmed record creation and
authorized Mission transitions. Scenario C adds Evidence and Decision persistence, Contract
inspection and authorized reconciliation, Mission resolution, and post-closure archival.

Every mutating transition or reconciliation requires an explicit Decision reference and expected
fingerprint. Guided `/spectacular propose|define|assess|decide|resolve|reconcile` operations retain
judgment. Mechanical CLI commands only validate and apply already-authorized effects.

Canonical v2 rejects generic `status`, `inspect`, `new`, `advance`, universal `doctor`, and bare
mechanical `assess`, `decide`, or `resolve` commands. `workspace validate` is scoped and read-only;
correction remains owned by the relevant operation.

## Consequences

- M1 may be integrated without pretending it already contains a CLI or lifecycle engine.
- Scenario A validates the substrate through cold use before governed writes are frozen.
- Scenario B proves accountable work formation and authorization end to end.
- Scenario C proves closure and makes a second cold recovery the final MVP loop test.
- Release hardening follows behavioral proof rather than substituting for it.
