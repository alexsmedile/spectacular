---
status: planned
priority: high
owner: alex
updated: 2026-05-30
summary: "Configurable policy mechanism — projects declare action-gating policies; ships with a built-in 'understand before you change' pre-implementation gate (now/changes/stays-same) as the default"
related:
  - PRD.md
  - ../../PRINCIPLES.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
target_version: v1.12.0
---

# Plan — policy-engine

> **URGENT — consider immediately after verify-walk (v1.11).** Captured 2026-05-30 from a live need: an agent should not start implementation without first establishing *how the current system works, what changes, and what stays the same*. That specific gate is the built-in default; the broader ask is a general, configurable policy mechanism.

## 1. Goal

Give Spectacular a **policy layer**: declarative, configurable rules that gate actions at defined lifecycle points. Ship with sensible defaults (the headline one: a **pre-implementation understanding gate** — before a request moves `planned → active`, the agent must produce a *current-system / what-changes / what-stays-the-same* analysis) while letting any project enable, disable, or add its own policies via `config.yaml`.

## 2. Constraints

- **Defaults provided, fully configurable.** Built-in policies ship enabled with safe defaults; a project overrides via `.spectacular/config.yaml` (enable/disable/add/tune). No policy is hard-coded into the skill flow.
- **Distinct from PRINCIPLES.md.** PRINCIPLES are *beliefs* with prose "how the skill enforces this" notes. Policies are *executable gates* with a trigger point + a check + a pass/block outcome. Principles inspire; policies enforce. A policy may *implement* a principle (the understanding-gate implements Principle 7's intent layer).
- **Advisory vs blocking is per-policy.** Each policy declares its severity — `block` (refuse to proceed) vs `warn` (surface + continue). Default for the understanding gate: configurable, leaning `warn` first (humans decide), with opt-in `block`.
- **Reuses existing machinery.** Lifecycle transitions (`promote`), the doctor severity model, frontmatter signals, and the rules-file dispatch pattern already exist — policies should compose with them, not duplicate.
- **Skill-enforced, not just CLI.** Like grill/verify, policy *evaluation* that needs judgment is skill-side; mechanical presence-checks (does the analysis note exist?) can be CLI/doctor.

## 3. Milestones

- M1 — **Policy contract + schema.** `references/policies-contract.md`: what a policy is (id, trigger point, check, severity, default state), where they live, how `config.yaml` enables/disables/tunes them. 4-tier scope precedence mirroring convention-packs (project → user → app-store → bundled) is a candidate.
- M2 — **Built-in: understanding gate.** The default `understand-before-change` policy — triggered at `planned → active`, requires a `## Current system / Changes / Stays the same` analysis (in PLAN, a new slot, or a dedicated note). Defines what "satisfied" means and the warn/block outcomes.
- M3 — **Enforcement points.** Wire policy evaluation into the lifecycle (`promote`) + a `doctor policies` area for presence-checkable policies. Skill consults active policies at each declared trigger.
- M4 — **Config surface.** `config.yaml` `policies:` block — enable/disable built-ins, set per-policy severity, register custom policies. Document precedence + a worked custom-policy example.
- M5 — **Dogfood + ship.** Enable the understanding gate on this repo; drive 1+ real request through it; CHANGELOG + plugin bump to v1.12.0.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- **Sequenced after [[verify-walk]] (v1.11).** verify-walk establishes the skill-side "walk + gate + lifecycle-flip" pattern at `review → verified`; policy-engine generalizes that gate idea to *any* trigger point and makes it configurable. Building verify-walk first gives a concrete second example to generalize from (avoids premature abstraction).
- Composes with the convention-pack scope model ([[packs-contract]]) — the 4-tier precedence may be directly reusable.
- Touches [[lifecycle]] (`promote` gains policy consultation) and [[doctor]] (new `policies` area).
- Relates to PRINCIPLES.md — policies are the *executable* counterpart to its prose enforcement hooks.

## 6. Validation

- M1 — `policies-contract.md` documents the full schema; a reader can author a new policy from it alone.
- M2 — The understanding gate fires on `planned → active`; a request lacking the analysis is flagged (warn or block per config); a request with it passes clean.
- M3 — `doctor policies` reports policy status; `promote` consults active policies at the right transition.
- M4 — `config.yaml` can disable a built-in, change a severity, and register a custom policy — each takes effect; precedence resolves predictably.
- M5 — This repo runs with the gate enabled; a real request shows the analysis was required before `active`; manifests at v1.12.0.

## 7. Deliverables

- `references/policies-contract.md` — policy schema + scope/precedence + config surface
- Built-in `understand-before-change` policy (the now/changes/stays-same pre-implementation gate)
- `config.yaml` `policies:` block + precedence resolution
- Lifecycle wiring (`promote` consults policies) + `doctor policies` area
- A worked custom-policy example in docs
- CHANGELOG [1.12.0] entry

## Open questions (resolve in M1)

- **Where does the understanding analysis live?** A new PLAN slot (8th)? A dedicated `ANALYSIS.md` per request? A free note the gate just checks for presence of? (Leaning: a PLAN slot for small requests, escalating to a note for big ones — mirrors the VERIFY.md 2-of-N pattern.)
- **Default severity** of the understanding gate — `warn` or `block` out of the box?
- **Policy vs principle boundary** — should each built-in policy link to the principle it enforces, making PRINCIPLES.md the "why" and policies the "how it's checked"?
- **Trigger vocabulary** — what set of trigger points exists? (`planned→active`, `active→review`, `review→verified`, `pre-archive`, `pre-commit`…?) Start minimal.
- **Overlap with verify-walk** — is `review→verified` verification itself just *a policy*? If so, policy-engine might eventually subsume verify-walk's gate. Worth noting, not forcing.
