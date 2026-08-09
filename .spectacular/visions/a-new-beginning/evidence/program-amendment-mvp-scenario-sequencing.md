---
type: refactor-program-amendment-proposal
status: proposed
authority: advisory
date: 2026-08-09
target: EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.6
requested_by: owner
disposition_required_from: central-orchestrator-and-owner
observed_feature_branch: codex/feat/v2-semantic-substrate
observed_feature_head: a488b2efe7828f59724f730b9a590b9a644e6885
review_handoff: H28
review_target: 489bd6008e1720e4b0310b999a0bac02c62df6dc
next_action: orchestrator-verify-and-present-owner-decision-packet
---

# Proposed program amendment — prove the MVP through scenarios

## Authority boundary

This is an owner-requested proposal to the central orchestrator. It does not revise an accepted
constitution, accept M1, supersede H28, activate M2, authorize implementation, or mutate lifecycle
state. The central orchestrator must verify the evidence, reconcile it with the accepted contracts,
and obtain explicit owner dispositions before changing the executable program.

## Recommendation

Keep the v2 architecture and amend the execution program.

The identity model, strict path boundary, typed relationships, canonical semantic records, stable
refusal taxonomy, clean v1/v2 break, owner authority, and independent fresh-context review are the
right foundation. The program should stop treating completion of horizontal layers as product
proof, however. It should drive the foundation through the three accepted rebuilt-product scenarios
before freezing downstream design.

The constitutional MVP is:

1. **Cold recovery:** a cold person or agent finds accepted direction, current state, evidence,
   blockers, and one safe next action or the exact unresolved decision, with source drill-down.
2. **Fuzzy intent to bounded accountable work:** an incomplete idea becomes an owner-accepted,
   bounded handoff with authority, responsibility, dependencies, stop conditions, expected
   evidence, and return expectations.
3. **Return, evaluate, reconcile, and resume:** returned evidence is evaluated against accepted
   expectations, receives an authorized disposition, is reconciled into project truth, and supports
   a new cold recovery.

Scenario C may be implemented after A and B, but it cannot be cut from the public MVP. It carries
the Product Constitution's evidence, reconciliation, and durable-continuation failure test.

## Evidence for reconsideration

### Constitutional mismatch

`SUCCESS-EVIDENCE-CONSTITUTION.md@1.0` defines the minimum rebuilt-product acceptance in Scenarios
A, B, and C. Its evidence table says protected behavior needs an end-to-end scenario plus focused
behavioral checks; file presence or an agent completion claim is insufficient alone. M1's focused
tests are necessary evidence for irreversible storage semantics, but cannot by themselves establish
that the rebuilt product works.

`PRODUCT-CONSTITUTION.md@1.0` says Spectacular has failed if it cannot preserve accepted intent,
bounded accountable work, declared authority, traceable handoffs, evidence-backed acceptance,
reconciliation, and durable continuation. The executable program lists these capabilities across
M2–M4 but does not create an early end-to-end contact point before the substrate is treated as
settled.

### Independent-review convergence

The supplied independent critiques converge on two material risks:

- governance ceremony is expensive relative to the behavior being proved; and
- strict M1 → M2 → M3 → M4 serialization delays user and agent feedback until the CLI/Skill stage.

They also identify a real undeclared product behavior: official writes normalize frontmatter that
the parser does not semantically own, including ordering, scalar presentation, comments, and line
endings. That behavior is implemented but is not yet an explicit user-facing promise.

The critiques' objections to exact path containment should **not** be adopted. Rejecting absolute
paths, traversal, backslashes, ambiguous bare names, and fuzzy matches inside the substrate is a
security and determinism boundary. Human-friendly discovery can be layered above the exact resolver
without weakening it.

### Live control-state observation

H28 correctly binds an immutable review of `489bd6008e1720e4b0310b999a0bac02c62df6dc` and correctly
cites the H27 evidence-return digest. The hash is not defective.

At proposal time, however, `codex/feat/v2-semantic-substrate` points to
`a488b2efe7828f59724f730b9a590b9a644e6885`, three commits beyond H28's review target. The later
commits materially change semantic validation, canonicalization, refusals, tests, and the YAML
dependency. The branch now contains six M1 commits from its declared base despite the three-commit
ceiling. This does not invalidate H28's immutable target; it means H28 cannot establish acceptance
of the current feature head.

The latest head already resolves one review concern by pinning stable `go.yaml.in/yaml/v3 v3.0.5`
and explicitly refusing unsupported YAML graphs. Do not carry the stale prerelease-pin finding into
the amended program.

## Invariants that remain non-negotiable

The scenario-first sequence must not relax M1's irreversible substrate properties:

- canonical lowercase UUIDv7 identity is independent from mutable workspace-relative paths;
- the substrate resolver remains exact, containment-safe, deterministic, and ambiguity-refusing;
- relationships are typed and remain valid across rename or move;
- reads do not rewrite files;
- official writes use an explicit, idempotent canonical-normalization contract;
- unknown supported YAML values and the Markdown body survive semantically; unsupported structures
  refuse before mutation;
- fingerprints are stable over declared semantic content;
- refusal codes and the versioned JSON envelope are stable and machine-readable;
- refusals, failed validation, stale expected fingerprints, and unauthorized effects cause no
  workspace mutation;
- write concurrency uses an expected fingerprint or equivalent optimistic guard; and
- persistence claims distinguish process-crash atomicity from power-loss durability and are tested
  to the strength actually promised.

These are stored-data and authority boundaries: costly to reverse after real projects adopt v2.
Command spelling, card layout, prose wording, and implementation mechanics remain reversible.

## Proposed executable sequence

### G0 — repair the program controls

Before accepting M1 or dispatching successor work:

1. Preserve H28 as historical evidence and choose an explicit review generation:
   - review the original H28 target, then issue a separate repair-generation review; or
   - supersede H28 before dispatch and issue a fresh review against the exact current head.
   The recommendation is the second option because the current head contains material corrections.
2. Replace commit-count ceilings with controls that measure risk: owned-path scope diff, prohibited
   effects, changed-invariant inventory, dependency diff, clean review target, and hypothesis/repair
   accounting. Commits should remain coherent review units, not a quota.
3. Keep the accepted constitutional history intact, while reducing the live execution surface to
   two controlling projections: the current executable program and the current Mission charter.
   Older contracts remain immutable referenced inputs rather than repeatedly restated ceremony.
4. Bind every review to an exact commit and tree. A later branch head requires a new generation; it
   must never silently inherit an earlier verdict.

### S0 — accept explicit storage semantics

Finish M1 review against the selected exact head, but require owner disposition on two promises:

- **Normalization:** parsing is non-mutating; an explicit official write canonicalizes known and
  supported unknown frontmatter, may discard presentation-only YAML formatting/comments, preserves
  semantic values and the opaque Markdown body, and is idempotent thereafter.
- **Durability:** either promise only same-process/crash-safe atomic replacement and name the
  power-loss limitation, or implement and test the platform-specific durable replacement sequence,
  including containing-directory synchronization where supported. Do not imply zero-loss under
  power failure without that stronger contract.

### A — cold-recovery tracer

Pull a minimal read-only vertical slice forward before the governed state engine is frozen. Exact
command names are Type 2; the capability must:

- list and inspect records in deterministic human and versioned JSON projections;
- expose identity, type, path, state, fingerprint, typed relationships, blockers, evidence links,
  and one safe next action or the exact unresolved gate;
- drill every consequential conclusion to its source;
- mutate nothing; and
- run both on a clean disposable workspace and on Spectacular's own workspace.

Have a genuinely cold actor perform Scenario A. Use the misses to correct substrate and projection
contracts before treating them as frozen.

### B — governed-work tracer

Drive one incomplete idea through the real authority boundary. A thin implementation may be guided
or partly manual, but the durable records must expose uncertainty and boundaries, capture explicit
owner acceptance, define expected evidence, and compile one bounded handoff. Test illegal and
unauthorized transitions as structured refusals with zero mutation.

This is where the Proposal/Mission ontology earns its place. Do not substitute a state-transition
matrix alone for the end-to-end scenario.

### C — closure tracer

Use a provider-neutral fake receipt to return evidence, assess it against the accepted expectation,
obtain an owner disposition, reconcile accepted changes into project truth, and perform Scenario A
again from cold context. The second cold recovery is the closure proof.

### R — release hardening

Only after A, B, and C pass end to end, complete the reversible productization work: generated
noun-first CLI registry, guided Skill, self-contained macOS/Linux binary, install and checksum flow,
race and integration checks, clean-workspace recovery, self-hosting, and proof that v2 core has no v1
runtime dependency.

## Test-first acceptance shape

Focused tests should be written before each slice, but they support rather than replace scenario
evidence:

- deterministic list/inspect JSON and refusal-envelope tests;
- ID/path equivalence, rename-safe relationships, duplicate-order independence, and source drill-down;
- read operations proven byte- and metadata-non-mutating;
- canonicalization idempotence and diff-visible first-write normalization;
- legal, illegal, prerequisite-gated, unauthorized, and stale-fingerprint transitions;
- write/sync/close/replace failure injection with original preservation and no temporary residue;
- an explicitly graded power-loss durability test or an explicit limitation;
- cold Scenario A on disposable and self-hosted workspaces;
- one real Scenario B; and
- Scenario C followed by a fresh Scenario A.

MVP acceptance requires all three scenarios plus their focused checks on an exact reviewed build.
Test count, commit count, file count, and agent completion claims are not acceptance metrics.

## Explicit MVP cuts

Defer unless a scenario proves the need:

- fuzzy or bare-name substrate lookup and path aliases;
- persistent caches, daemons, and graph scheduling;
- elaborate Fog visualizations and card polish;
- packs, app-store behavior, and broad collections;
- decision companions and multiplexer UI;
- real provider integrations beyond a deterministic fake;
- generic v1 migration or v1 compatibility in v2 core;
- Windows distribution; and
- broad release polish unrelated to install, recovery, and the three scenarios.

Keep the clean v1/v2 break.

## Owner decisions the orchestrator must grill

1. **Review generation:** preserve and run H28 at `489bd600...`, or supersede it with a fresh
   current-head review? Recommendation: preserve as historical, supersede before dispatch, reissue.
2. **Risk control:** retire the three-commit ceiling in favor of scope/effect/invariant/dependency
   checks? Recommendation: yes.
3. **Program shape:** amend horizontal M2/M3/M4 delivery into A → B → C scenario tracers followed by
   release hardening? Recommendation: yes.
4. **Normalization promise:** may official writes intentionally erase YAML presentation details
   while preserving declared semantics and body content? Recommendation: yes, only as an explicit,
   diff-visible, idempotent contract.
5. **Durability grade:** is process-crash atomicity sufficient for MVP, or does the product claim
   power-loss durability? Recommendation: implement the strongest practical macOS/Linux replacement
   protocol before the first write-bearing public slice; otherwise state the narrower guarantee.
6. **Live governance surface:** keep immutable constitutional history but operate from only the
   current program and current Mission charter? Recommendation: yes.
7. **MVP boundary:** confirm Scenario C is mandatory for MVP rather than post-MVP. Recommendation:
   mandatory, because the accepted constitutions make evaluation, reconciliation, and continuation
   identity-defining behavior.

## Recommended disposition and one next action

Recommended orchestrator disposition: **adapt** the executable program; do not discard the accepted
architecture and do not accept the external critiques wholesale.

Exactly one next action: **the central orchestrator verifies the branch and constitutional citations,
presents decisions 1–7 as one owner amendment packet, and stops for explicit disposition.**
