---
status: planned
updated: 2026-05-30
related:
  - PLAN.md
---

# Tasks — policy-engine

> URGENT — pick up immediately after verify-walk (v1.11). Design fully locked 2026-05-30 (see PLAN § Locked decisions).

## M1 — Practice-layer contract
- [ ] Write `references/policies-contract.md`: POLICY.md = sections keyed by `## @hook`; theory(PRINCIPLES)→practice(POLICY) framing
- [ ] Document policy anatomy: `hook + principle?(optional) + severity(block|warn) + check + prose`
- [ ] Document the locked hook set (8): spine @Init/@Planning/@Implementation/@Verification/@Archive + moments @Remember/@Snapshot/@SessionEnd
- [ ] Document config surface (config-only for v1; note 4-tier precedence as v2 candidate)
- [ ] Define "satisfied": mechanical presence-check (doctor) vs skill judgment

## M2 — POLICY.md scaffold + 8 defaults
- [ ] Add `doc_policy()` to `cli/spectacular`; emit POLICY.md with the 8 prefilled policies (4 block / 4 warn per PLAN §4)
- [ ] Add POLICY.md to the **always-set** list (created on every init); confirm new always-set count
- [ ] Spec the `## Understanding` PLAN slot (How it works now / What changes / What stays the same) + `UNDERSTANDING.md` escalation; gate satisfied by either
- [ ] Register the `## Understanding` slot in `plan-overrides.md` + `scaffold-reference.md`

## M3 — `spectacular policy` verb + injection loop
- [ ] Implement `spectacular policy` (forms: bare list / `@<hook>` / `<id>` / `--principle N` / `--json`)
- [ ] `@<hook>` output pulls the hook's policies + each linked principle's heading + one line (skim, read-verbs style)
- [ ] Reference doc + SKILL.md routing: the injection loop (enter phase → retrieve → inject → evaluate → resolve by severity)

## M4 — Enforcement + config
- [ ] Wire policy consultation into `promote`/`archive` at the spine hooks
- [ ] Add `doctor policies` area (presence-checks: POLICY.md exists, `## Understanding` filled, etc.)
- [ ] `config.yaml` `policies:` block — per-policy enable/disable/severity + register custom
- [ ] Worked custom-policy example in docs

## M5 — Dogfood + ship
- [ ] Enable POLICY.md on this repo (write the 8 defaults into `.spectacular/POLICY.md`)
- [ ] Drive 1+ real request through `@Implementation` with a filled `## Understanding` before active
- [ ] CHANGELOG [1.12.0] entry; plugin bump to v1.12.0
- [ ] docs/commands.md + configuration.md: document `spectacular policy` + the `policies:` config block

## Resolved before building (2026-05-30 — see PLAN § Locked decisions)
- [x] POLICY.md as practice layer paired with PRINCIPLES.md (theory→practice)
- [x] Hook naming → `@`="at" + moment-noun; 8 hooks locked; @Request/@RequestTask folded into @Planning
- [x] POLICY.md → single file, no soft-DB in v1
- [x] POLICY.md → always-set (every init, 8 prefilled, enabled)
- [x] Severity split → 4 block / 4 warn
- [x] Understanding content → `## Understanding` slot, escalates to UNDERSTANDING.md (no ANALYSIS.md)
- [x] Principle link → optional `principle:` tag
- [x] Enforcement → skill-side + doctor; no hooks.json in v1
- [x] Scope model → config-only v1; 4-tier deferred v2
- [x] verify-walk → absorbed as `verification-present` policy; not refactored in v1.12
