---
status: verified
updated: 2026-05-31
related:
  - PLAN.md
---

# Tasks — policy-engine

> URGENT — pick up immediately after verify-walk (v1.11). Design fully locked 2026-05-30 (see PLAN § Locked decisions).

## M1 — Practice-layer contract
- [x] Write `references/policies-contract.md`: POLICY.md = sections keyed by `## @hook`; theory(PRINCIPLES)→practice(POLICY) framing
- [x] Document policy anatomy: `hook + principle?(optional) + severity(block|warn) + check + prose`
- [x] Document the locked hook set (8): spine @Init/@Planning/@Implementation/@Verification/@Archive + moments @Remember/@Snapshot/@SessionEnd
- [x] Document config surface (config-only for v1; note 4-tier precedence as v2 candidate)
- [x] Define "satisfied": mechanical presence-check (doctor) vs skill judgment

## M2 — POLICY.md scaffold + 8 defaults
- [x] Add `doc_policy()` to `cli/spectacular`; emit POLICY.md with the 8 prefilled policies (4 block / 4 warn per PLAN §4)
- [x] Add POLICY.md to the **always-set** list (created on every init); confirm new always-set count
- [x] Spec the `## Understanding` PLAN slot (How it works now / What changes / What stays the same) + `UNDERSTANDING.md` escalation; gate satisfied by either
- [x] Register the `## Understanding` slot in `plan-overrides.md` + `scaffold-reference.md`

## M3 — `spectacular policy` verb + injection loop
- [x] Implement `spectacular policy` (forms: bare list / `@<hook>` / `<id>` / `--principle N` / `--json`)
- [x] `@<hook>` output pulls the hook's policies + each linked principle's heading + one line (skim, read-verbs style)
- [x] Add the gate block (`run spectacular policy @<hook>, follow active policies`) to the head of each phase ref doc (PLAN §8 mapping: init-workflow / new-request / active-request / verification / archive / memory / versioning / sessions-rules)
- [x] Write `references/policy-injection.md`: the loop (enter phase → retrieve → inject → evaluate → resolve by severity); severity-default = non-blocking

## M4 — Enforcement + config
- [x] Wire policy consultation into `promote`/`archive` at the spine hooks
- [x] Add `doctor policies` area (presence-checks: POLICY.md exists, `## Understanding` filled, etc.)
- [x] `config.yaml` `policies:` block — per-policy enable/disable/severity + register custom
- [x] Worked custom-policy example in docs

## M5 — Dogfood + ship
- [x] Enable POLICY.md on this repo (write the 8 defaults into `.spectacular/POLICY.md`)
- [x] Drive 1+ real request through `@Implementation` with a filled `## Understanding` before active
- [x] CHANGELOG [1.12.0] entry; plugin bump to v1.12.0
- [x] docs/commands.md + configuration.md: document `spectacular policy` + the `policies:` config block

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
- [x] Severity opt-in → blocks only if explicit `severity: block`; absent/warn → non-blocking (safe default)
- [x] Enforcement mechanic → gate block at head of each phase ref doc (ref doc = phase trigger); no event bus
- [x] Source of truth → POLICY.md defines, config.yaml overrides; both CLI-managed + hand-editable, structure-bound; no check-kind field
