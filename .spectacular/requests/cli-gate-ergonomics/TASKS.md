---
status: planned
updated: 2026-07-12
related:
  - PLAN.md
---

# Tasks — cli-gate-ergonomics

## v1

### M1 — Policy gate tiered output + directive field
- [ ] Parser: read optional `- directive:` field per policy (alongside principle/severity/check)
- [ ] Add `_policy_principle_title` helper (number + title: text before the first ` — ` of the principle line)
- [ ] Hook form: warn rows print `— <directive>` + `→ P<n> — <title>` (title-only fallback when no directive); block rows print directive + full principle line
- [ ] Add `--full` flag restoring full paragraphs; `--json` gains `directive` key
- [ ] Author `- directive:` for all existing policies in `.spectacular/POLICY.md` (one imperative sentence each)
- [ ] Document the field in policies-contract.md anatomy + templates/policy base
- [ ] Confirm `_policy_consult_transition` (advance) inherits the tiering unchanged
- [ ] Write `tests/cli/policy-output.test.sh`
- [ ] → check: PLAN §Validation M1

### M2 — advance scaffolds SESSION.md
- [ ] Embed SESSION.md heredoc template (frontmatter + 4 H2s, mirrors active-request.md's)
- [ ] Write it on planned→active when absent; print `✓ scaffolded: SESSION.md`; never overwrite
- [ ] Add advance case to `tests/cli/mutator.test.sh`
- [ ] → check: PLAN §Validation M2

### M3 — doctor findings block
- [ ] `doctor_emit_text`: collect non-pass rows during the loop; print `── findings ──` block after the count line
- [ ] Clean run prints no block
- [ ] Add case to `tests/cli/doctor.test.sh`
- [ ] Update doctor.md / status.md only where output shape is quoted
- [ ] → check: PLAN §Validation M3

## v2 (deferred)

- [~] `doctor --quiet` (findings-only output) if agent usage shows the repeat block isn't enough
- [~] `doctor policies` info-nudge for enabled policies missing a `directive:` field
