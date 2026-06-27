---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Verify Log — naming-coherence

Validation walk for the v1.19.0 naming-coherence work. Checks map to PLAN § 6.

| Check | Kind | Evidence | Result |
|---|---|---|---|
| M1 — `advance` works; `promote` warns + still routes | executable | `tests/cli/mutator.test.sh` scenario 7b (advance advances, promote prints "deprecated", state moves) | ✅ |
| M1 — help/docs teach `advance` | observable | `spectacular advance --help`; commands.md/lifecycle.md/verify.md grep clean of lifecycle `promote <slug>` | ✅ |
| M2 — `feedback` canonical; `feedback-loop` alias | executable | `spectacular feedback list` + `spectacular feedback-loop list` both succeed | ✅ |
| M2 — backtick heredoc bug fixed | executable | bare `spectacular feedback` prints usage, no longer executes `spectacular remember` (verified via `bash -x`) | ✅ |
| M3 — doc-id `pack` (alias `convention-pack`) | observable | `doc-id: pack` + `doc-id-aliases: [convention-pack]` in pack-rules.md; doc-index.md row = `pack` | ✅ |
| M3 — `pack new/grill/refine/review` redirect | executable | `spectacular pack grill x` → skill-flow message (was "unknown pack verb") | ✅ |
| M4 — `spectacular next` read-only, all states | executable | `tests/cli/mutator.test.sh` scenario 7c (empty→CTA, planned→names request, status unchanged, rejects args) | ✅ |
| M5 — one-line tier-reveal, never mid-flow | observable | new-request.md (grill suggestion after scaffold), active-request.md (advance→review at all-checked), archive.md (spec-sync, pre-existing) | ✅ |
| Regression — full suite | executable | `./tests/run.sh` → 9/9 areas pass; `doctor` → 0 errors / 0 warnings | ✅ |

All checks pass. SKILL.md description re-measured at 1014 chars (under Codex's 1024 limit).
