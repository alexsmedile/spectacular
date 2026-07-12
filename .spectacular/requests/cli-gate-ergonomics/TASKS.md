---
status: verified
updated: 2026-07-12
related:
  - PLAN.md
---

# Tasks — cli-gate-ergonomics

## v1

### M1 — Policy gate tiered output + directive field
- [x] Parser: read optional `- directive:` field per policy (alongside principle/severity/check)
- [x] Add `_policy_principle_title` helper (number + title: text before the first ` — ` of the principle line)
- [x] Hook form: warn rows print `— <directive>` + `→ P<n> — <title>` (title-only fallback when no directive); block rows print directive + full principle line
- [x] Add `--full` flag restoring full paragraphs; `--json` gains `directive` key
- [x] Author `- directive:` for all existing policies in `.spectacular/POLICY.md` (one imperative sentence each)
- [x] Document the field in policies-contract.md anatomy + templates/policy base
- [x] Confirm `_policy_consult_transition` (advance) inherits the tiering unchanged
- [x] Write `tests/cli/policy-output.test.sh`
- [x] → check: PLAN §Validation M1 — 18/18 asserts green; 21/21 policies carry a directive

### M2 — advance scaffolds SESSION.md
- [x] Embed SESSION.md heredoc template (frontmatter + 4 H2s, mirrors active-request.md's)
- [x] Write it on planned→active when absent; print `✓ scaffolded: SESSION.md`; never overwrite
- [x] Add advance case to `tests/cli/mutator.test.sh`
- [x] → check: PLAN §Validation M2 — scenario 11 green (8 asserts); doctor lifecycle 0 warnings; overwrite-never proven

### M3 — doctor findings block
- [x] `doctor_emit_text`: collect non-pass rows during the loop; print `── findings ──` block after the count line
- [x] Clean run prints no block
- [x] Add case to `tests/cli/doctor.test.sh`
- [x] Update doctor.md / status.md only where output shape is quoted (doctor.md report example; status.md quotes no doctor shape — untouched)
- [x] → check: PLAN §Validation M3 — scenario 22 green (error + warning rows repeat in block; clean-half env-guarded)

## v2 (deferred)

- [~] `doctor --quiet` (findings-only output) if agent usage shows the repeat block isn't enough
- [~] `doctor policies` info-nudge for enabled policies missing a `directive:` field
