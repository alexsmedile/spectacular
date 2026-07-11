---
status: planned
updated: 2026-07-11
related:
  - PLAN.md
---

# Tasks — stance-layer

<!-- All design is pre-locked in PLAN.md + [[ideas/stance-layer]]. These tasks are execution only. -->

## v1

### M1 — architectural-stance policy live
- [ ] Append `### architectural-stance` block under `## @Planning` in POLICY.md (use the final text in ideas/stance-layer.md Part 2 "Final policy shape")
- [ ] Decide the `principle:` tag — reuse 11, or add a new judgment principle to PRINCIPLES.md
- [ ] Add an `**Override:**` clause (L5 style) and confirm NO `⛔` marker (warn, not block)
- [ ] → check: `spectacular policy @Planning` lists it; `spectacular doctor policies` exits 0

### M2 — grade label recognized by status
- [ ] Add `grade` to the `status` fleet-row render path
- [ ] Add `grade` to the `status <slug>` card render path
- [ ] Add `grade` to `status --json` output
- [ ] Confirm absent `grade:` resolves to standard (no error, no crash)
- [ ] → check: request with `grade: mvp` shows `mvp` in card + `--json`; absent renders clean

### M3 — doctor lifecycle validates grade enum
- [ ] Add closed-enum check (`prototype|mvp|standard|production`) to `doctor lifecycle`
- [ ] Warning-class, active requests only, mirrors the existing `status:`/`hold:` enum checks
- [ ] → check: `grade: protoype` emits exactly one warning; the 4 valid values emit none

### M4 — skill offers spectacular decide on a real fork
- [ ] Wire the offer-to-`decide` prompt into the @Planning skill flow (active-request.md or plan-rules.md)
- [ ] Trigger only on a real fork (crosses boundary / sets precedent / two viable structures); silent otherwise
- [ ] → check: boundary-crossing change surfaces the offer; trivial edit stays silent

### M5 — docs + tests synced
- [ ] Note optional `grade:` in plan-rules.md + scaffold-reference.md
- [ ] Note `architectural-stance` in specs/index.md + the policy-engine spec
- [ ] Add a `doctor lifecycle` grade-enum test + a `status` grade-render test
- [ ] CHANGELOG entry; plugin version bump
- [ ] → check: grep confirms each doc mention; new tests pass in `tests/run.sh`

## v2 (deferred)

- [~] L2 legibility pass — imperative `check:` phrasing (separate scoped request; touches machine-read check lines)
- [~] grade escalation verb / mid-flight transition logging (only if a real need surfaces)
- [~] `config.yaml` project-wide grade default (dropped from v1 as unneeded for a label)
