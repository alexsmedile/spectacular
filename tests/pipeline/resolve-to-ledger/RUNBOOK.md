# P1 — resolve → ledger graduation

**Proves:** a resolved debug job graduates into the permanent library — `fixes/F<N>` (and
optionally `audit/A<N>`) with cross-links and a `debug_job` back-link — via the CLI verbs, single-writer.

**Bug:** `fixtures/discount.py` — `apply_discount` uses `max()` where it needs `min()`, so a 95%
coupon charges more than list. Closed, single-site, reusable (clamp-direction footgun).

---

## Steps (orchestrator drives)

1. **Confirm the bug.** `python3 tests/pipeline/resolve-to-ledger/fixtures/discount.py`
   - **ASSERT:** exits non-zero with an AssertionError (bug is real, repro in hand).

2. **Triage (Step 0–1).** Grep `.spectacular/fixes/` for a prior match; there is none. Ceremony gate:
   root cause clear (max vs min) + single site → **just-fix**, but we want the ledger graduation, so
   treat it as a fleet job to exercise the full back half.
   - **ASSERT:** you classify it `symptom_class: wrong_behavior`, ceremony `just-fix`.

3. **Open the job (Step 1c).** Scaffold `.spectacular/debug/discount-clamp/` and write `job.json`
   (spine): `symptom_class`, `symptom`, `reporter: user`, `ceremony`, `status: fixing`, the `brief`
   field, empty `artifacts` + `timeline`, `outcome: null`.
   - **ASSERT:** `python3 -m json.tool .spectacular/debug/discount-clamp/job.json` parses, and it has
     a non-null `brief` and a `symptom_class`.

4. **Apply the fix** (inline — it's 1 site, self-serve per Step 1b). Copy the fixture to a scratch
   `runs/` copy, change `max` → `min`, run the fixture check.
   - **ASSERT:** the check now prints `ok`, exit 0.

5. **Write `outcome.json`** — `disposition: resolved`, `fix_ids`, `logged_fixes: ["F<N>"]` (the id you
   will create in step 6), `audit: null` (obvious cause, no separate examination earned), a `note`.
   - **ASSERT:** file parses; `disposition == "resolved"`.

6. **Graduate to the ledger** — `spectacular fix new` with `--problem/--intended/--cause/--fix/--criteria
   --verified-by "<the check>" --signature "<clamp direction, max-vs-min, discount cap>"`. Stamp the
   `debug_job: discount-clamp` back-link on the created `F<N>.md`.
   - **ASSERT:** `spectacular fix list` shows the new `F<N>`.
   - **ASSERT:** the `F<N>.md` frontmatter carries `debug_job: discount-clamp`.
   - **ASSERT:** `outcome.json.logged_fixes` and the actual `F<N>` id agree.

7. **Flip the spine to resolved.** Update `job.json` → `status: resolved`, `outcome: outcome.json`,
   `artifacts.fixes` populated.
   - **ASSERT:** `job.json.status == "resolved"` and the debug folder is **still present** (kept as
     trace, not pruned — the whole point).

8. **Prove the loop closed.** Grep `.spectacular/fixes/` for the signature keyword you just logged.
   - **ASSERT:** the new `F<N>` is found — a *future* Step 0 would catch this bug. Self-learning loop verified.

---

## Pass criteria
All 8 assertions hold. Specifically: **the fix graduated to `F<N>`, the back-link exists both
directions (`outcome.logged_fixes` ↔ `F<N>.debug_job`), and the debug trace survived resolution.**

## Cleanup
Delete `.spectacular/debug/discount-clamp/` and the test `F<N>` after — this is a test, not a real
fix. (Or keep as a documented sample; note which in results.)

## Results (fill after running)
- Date: 2026-07-05
- Verdict: **PASS** (all 8 assertions)
- Notes:
  - Real integration run — `debug-fixer` subagent spawned live, applied `max→min`, verified `ok`/exit 0, wrote `fixes/fix-01.json`. Orchestrator graduated to `F4` (severity medium), stamped `debug_job: discount-clamp`, wrote `outcome.json`, flipped spine to `resolved`. Back-links agree both directions; debug folder kept; signature findable → loop closed.
  - **Agent-contract win (fixture bug found):** the Fixer's FIRST pass correctly **bounced** — the fixture's own `__main__` used exact float equality (`== 10.0` vs `100*(1-0.9) == 9.999...`), a latent test defect independent of the clamp bug. The Fixer applied the exact brief fix, verified honestly, and refused to freelance a second edit the brief didn't authorize. Exactly the apply-only contract working. Fixture corrected to `round(...,2)` tolerance (the clamp bug — what P1 tests — is unchanged), Fixer re-run → `applied`/pass.
  - **CLI note:** `spectacular fix new --severity` requires the full word (`medium`), not `m`. The `-i` init menu and doctor use `l|m|h` elsewhere — minor inconsistency, not a blocker.
  - **CLI gap:** `fix new` does not stamp `debug_job:` — the orchestrator adds it to frontmatter manually. Candidate for a `--debug-job <slug>` flag if graduation-from-debug becomes common.
  - Kept as documented sample (not cleaned up) so P5 has a prior fix to grep against in its Step 0.
