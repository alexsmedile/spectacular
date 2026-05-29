---
status: planned
updated: 2026-05-30
related:
  - PLAN.md
---

# Tasks — policy-engine

> URGENT — pick up immediately after verify-walk (v1.11).

## M1 — Policy contract + schema
- [ ] Write `references/policies-contract.md`: policy = id + trigger-point + check + severity + default-state
- [ ] Define trigger vocabulary (start minimal: `planned→active`, maybe `pre-archive`)
- [ ] Decide scope/precedence model (reuse convention-pack 4-tier? or simpler config-only)
- [ ] Define how a policy declares "satisfied" (presence check vs judgment check)

## M2 — Built-in: understand-before-change gate
- [ ] Specify the default policy: trigger `planned → active`, requires current-system / changes / stays-same analysis
- [ ] Decide where the analysis lives (PLAN slot 8? dedicated note? presence-only?)
- [ ] Define warn vs block outcomes + the default severity

## M3 — Enforcement points
- [ ] Wire policy consultation into `promote` at declared transitions
- [ ] Add `doctor policies` area for presence-checkable policies
- [ ] Skill consults active policies at each trigger; surfaces warn/block

## M4 — Config surface
- [ ] `config.yaml` `policies:` block — enable/disable built-ins, set severity, register custom
- [ ] Precedence resolution (project overrides bundled defaults)
- [ ] Worked custom-policy example in docs

## M5 — Dogfood + ship
- [ ] Enable the understanding gate on this repo
- [ ] Drive 1+ real request through the gate before `active`
- [ ] CHANGELOG [1.12.0] entry; plugin bump to v1.12.0

## Resolve before building (from PLAN open questions)
- [ ] Analysis location: PLAN slot vs dedicated note vs presence-only
- [ ] Default severity: warn vs block
- [ ] Policy↔principle linkage (does each policy cite the principle it enforces?)
- [ ] Does verify-walk's review→verified gate eventually become "just a policy"?
