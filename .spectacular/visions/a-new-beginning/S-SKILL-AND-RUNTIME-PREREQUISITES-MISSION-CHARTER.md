---
type: mission-charter
mission: skill-and-runtime-prerequisites
version: 1.0
status: ready-for-independent-rereview
prepared_at: 2026-08-10
prepared_by: mission-owner-gap-audit
accepted_by: owner
accepted_at: 2026-08-10
activated_at: 2026-08-10
activation_baseline_commit: eaef96577c0cf4bbe19debbada407f464c2aa7bb
activation_baseline_tree: 927230b65e74dec652e9b30e63231a397b7f9c40
branch: codex/feat/v2-s-skill-runtime-prerequisites
source_thread: 019feb82-33b4-7961-aebc-5d9471939a5a
design_sufficiency: sufficient
slice_quality: coherent
repair_budget: 2
repair_rounds_used: 1
implementation_completed_at: 2026-08-10
implementation_commit: e64f488bcf4f3279947d8d3fe9660b939f1be956
implementation_tree: ba69f34e858807c3e49d6e7b035614f52b3236aa
implementation_evidence: evidence/scenario-s-primary-implementation-evidence.md
first_reviewed_commit: 568ae88c5ea40938384656907094e9bf6bc3b5d6
first_reviewed_tree: 67ac6416dea58b71b3dec2bfc597034a3f00f2fb
first_review_disposition: block
repaired_commit: 7bede74ea4de75e1d6f887bc44d41c5a932e95cc
repaired_tree: 2f02997cc6b69978382ebe079c594d312d1a81ce
repair_evidence: evidence/scenario-s-repair-round-1-evidence.md
next_action: independent-rereview
upstream:
  - EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.9
  - RESPONSIBILITY-PLACEMENT-CONTRACT.md@1.0
  - RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0
  - PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0
  - MISSION-PREPARATION-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - IMPLEMENTATION-ARCHITECTURE-AND-MIGRATION-CONTRACT.md@1.0
  - SHARED-SCAFFOLD-CONTRACT.md@1.0
  - BC-GOVERNED-LOOP-MISSION-CHARTER.md@1.6
---

# Scenario S Mission Charter — Skill and runtime prerequisites

## Outcome

Deliver the canonical v2 guided Skill, progressive context compiler, proportional Mission
preparation and Autopilot charter compilation, runtime-neutral Handoff workflow, v2 Codex and
Claude manifests, and local installation prerequisites. Prove the integrated behavior in an
isolated clean-v2 Spectacular dogfood workspace.

## Preparation verdict

- Design Sufficiency: `sufficient`. Accepted contracts fix public language, responsibility,
  authority, preparation, retrieval, Handoff, evidence, and implementation boundaries.
- Slice Quality: `coherent`. The components share one Skill/runtime integration boundary and one
  end-to-end clean-runtime proof. Remaining encoding, package, and fixture choices are reversible.
- Decision delta: no unresolved Type-1 choice remains.

## Objectives

1. Derive the Skill-facing mechanical catalog from the Go command registry and compile bounded,
   source-backed project or Mission context with freshness, omissions, Gaps, and drill-down.
2. Implement guided A/B/C workflows, proportional preparation, explicit non-default Autopilot
   charters, and runtime-neutral Handoff/return behavior without creating authority.
3. Provide v2-only Codex and Claude manifests and local installation preflights, then prove the
   complete workflow through isolated dogfood, cold recovery, and independent review.

## Completion boundary

- The v2 Skill routes `orient`, `propose`, `define`, `decide`, `start`, `resume`, `handoff`,
  `assess`, `reconcile`, `resolve`, and `audit` while preserving the Skill/CLI judgment boundary.
- Unsupported Link/Message persistence is reported as an explicit capability Gap rather than
  invented by this Mission.
- Compiled context is deterministic, bounded, non-authoritative, freshness-aware, source-backed,
  and refuses missing, stale, conflicting, or scope-invalid inputs.
- Preparation represents separate Design Sufficiency and Slice Quality verdicts without a new
  lifecycle state or mandatory design document.
- Autopilot is explicit and non-default; its charter binds outcome, non-goals, sources, allowed and
  forbidden effects, budgets, checks, expiry, stops, recovery, and return destination.
- A Handoff remains within one Mission, survives host-runtime replacement, and never transfers
  Mission accountability or provider authority.
- Disposable Codex and Claude roots discover the same canonical v2 Skill and local binary without
  loading v1 surfaces or changing real user configuration.

## Invariants and dependency diff

- The Skill owns judgment; Domain and governed operations retain semantic and lifecycle authority.
- The Go registry remains the only owner of command arguments, effects, dispatch, and machine facts.
- Context, catalogs, cards, and charters are projections or validated inputs, never competing truth.
- Installation grants no credentials, provider permission, or external-effect authority.
- V1 CLI, Skill, manifests, hooks, tests, and collections remain outside v2 runtime and retrieval.
- No new dependency is planned. A required dependency or public-contract expansion is a charter
  delta and stop condition.

## Owned paths

```text
.spectacular/visions/a-new-beginning/S-SKILL-AND-RUNTIME-PREREQUISITES-MISSION-CHARTER.md
.spectacular/visions/a-new-beginning/S-SKILL-AND-RUNTIME-PREREQUISITES-MISSION-CHARTER.md.sha256
.spectacular/visions/a-new-beginning/evidence/scenario-s-*.md
v2/skills/spectacular/
v2/internal/command/
v2/internal/context/
v2/internal/guardrails/
v2/internal/runtime/
v2/internal/governance/
v2/cmd/spectacular/
v2/.codex-plugin/
v2/.claude-plugin/
v2/install/
v2/testdata/scenario-s/
```

Existing v2 Domain, Workspace, Index, Projection, and Scenario A/B+C surfaces may receive only
backward-compatible extensions required by this outcome.

## Authority and prohibited effects

Authorized effects are inspection, local edits, checks, disposable runtime roots, coherent local
commits, and work on `codex/feat/v2-s-skill-runtime-prerequisites`. Prohibited are v1 compatibility
or migration, release/distribution, global installation, real providers, deployment, push, PR,
merge, publication, credentials/configuration changes, destructive effects, permanent agent
fleets, generic schedulers, and duplicate status or projection authority.

## Proof and repair

Proof covers registry parity, deterministic context and refusals, compressed and deep preparation,
reslicing after a false assumption, shared-interface work, Autopilot envelope enforcement,
runtime replacement, disposable observed completion and incident recovery, Codex/Claude discovery,
and cold resume without chat or provider history. Required checks are focused tests, `gofmt`, module
verification, vet, full Go tests, race tests, build, Bash syntax, version guard, full v1 tests,
sidecar verification, exact scope diff, and one fresh read-only final reviewer bound to the exact
commit/tree who also reproduces the cold-context proof.

Two hypothesis-changing repair rounds are authorized. Stop only for a new Type-1 decision,
authority/safety/provider/irreversible-effect conflict, exhausted repair budget, or an unresolved
required check. The terminal return binds final commit/tree, claim-mapped Evidence, reviewer result,
repair usage, recovery point, remaining Gaps, and exactly one next action.
