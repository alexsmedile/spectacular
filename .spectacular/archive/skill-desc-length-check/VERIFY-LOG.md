---
request: skill-desc-length-check
build: b10
verified: 2026-06-30
verifier: alex (via spectacular skill)
artifact: PLAN § Validation (no VERIFY.md — 2-of-6 rule)
result: pass
---

# Verify log — skill-desc-length-check

Validation walk run against `PLAN.md § 6. Validation`. Work was shipped to `main`
prior to this walk; verification confirmed each check against the live CLI, the
test suite, and the running pre-commit guard.

| # | Check (PLAN § Validation) | Kind | Evidence | Result |
|---|---|---|---|---|
| M1 | Written note: Codex's measured field + cap, with v1.17.2 numbers | assertable | `DECISIONS.md` D7 — gates on `description` alone at 1024 (error)/1000 (warn); cites v1.17.2 (1146→986 cleared while combined 1253 stayed over → proves concatenation isn't measured) | ✅ |
| M2 | doctor: >1024 error · ~1010 warning · ≤1000 pass; live SKILL.md (now 983) passes; shows count | executable | `check_skill_desc_len()` cli:7010; `tests/cli/doctor.test.sh` scenario 17 — **53/53 assertions pass**; `doctor skill` on live SKILL.md → "description length ok (983 chars, under 1024 Codex cap)" | ✅ |
| M3 | Staging SKILL.md >1024 → commit rejected; trim → commit succeeds; shared logic | observable | `scripts/check-skill-desc.sh` (sourceable functions + standalone), wired via `scripts/hooks/pre-commit-wrapper` + `.active/` symlink so git-guard's `pre-commit` survives regen. Guard **fired live this session** ("✓ description 983 chars" on every commit) | ✅ |

## Ponytail confrontation (review pass)

M3 specified "shared logic, no duplication." The awk parser appears in **two**
places — `cli/spectacular:7013-7026` and `scripts/check-skill-desc.sh:32-56`.

- **Verified byte-identical logic** (same rules, thresholds, YAML-literal-block
  handling; only line-wrapping differs).
- **Duplication is justified, not slop:** `cli/install.sh` ships *only* the
  `spectacular` binary to `~/.local/bin` — `scripts/` does not travel, so the
  installed doctor cannot source the helper. The CLI copy must be self-contained;
  the helper serves the dev-repo pre-commit path. Legitimate install-boundary seam.
- **One residual gap (accepted):** parity is guarded only by a `# Kept in
  lock-step` comment — no test asserts the two blocks stay equal. Low-churn code,
  currently in sync; comment is a defensible ceiling. Logged here rather than
  fixed (user decision: archive as-is). Future drift would surface as a scenario-17
  failure only if the *helper* side changed, not the CLI side — so a true parity
  test is the v2 follow-up if this code ever churns.

## Verdict

All 3 validation checks pass. Duplication reviewed and accepted. Request →
verified, ready to archive.
