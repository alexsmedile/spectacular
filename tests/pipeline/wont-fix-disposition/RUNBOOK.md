# P6 — won't-fix disposition

**Proves:** the third disposition branch (`resolved | folded-into-request | wont-fix`) that P1–P5
never reached. A bug can be **real, reproducible, and one-line to fix** — and the honest resolve is
still *don't apply it*. The orchestrator must recognize decline-as-correct, write
`disposition: wont-fix` with a stated reason, close the spine, and log **no** `F<N>`.

**Bug:** `fixtures/legacy_export.py` — `export_csv_v1` emits MM/DD/YYYY (ambiguous for EU readers).
The fix is a trivial format swap. BUT v1 is deprecated: every live caller moved to `export_csv_v2`
(ISO-8601), and v1 survives only for one frozen legacy integration that *depends* on MM/DD/YYYY.
"Fixing" it breaks the sole remaining consumer. → decline, point at v2.

**This is a routing / judgment test, not a repro-crash test.** Nothing throws; both functions run.
The test is whether the orchestrator resolves it as `wont-fix` instead of reflexively applying the
tempting one-liner.

---

## Steps (orchestrator drives)

1. **Trigger + confirm.** Treat as a user bug report ("old CSV export uses US date format"). Run the
   fixture, read the code.
   - **ASSERT:** the ambiguity is real (`03/04/2026` is genuinely MM/DD), AND the docstring/context
     shows v1 is deprecated with a frozen consumer + a v2 migration target.

2. **Step 0 — seen it before?** Grep `.spectacular/fixes/`.
   - **ASSERT:** no prior match. Proceed.

3. **Step 1 — ceremony gate.** Cause is understood (format string), change is local — this reads like
   **just-fix** on the surface. The judgment is whether to fix at all.
   - **ASSERT:** classified `symptom_class: wrong_behavior`. No fleet needed (cause known, single site) —
     this is an inline judgment, **no debug/ folder required** (ceremony scales with uncertainty).

4. **The decline decision (the thing under test).** Weigh apply-vs-decline. Applying the one-liner
   breaks the frozen consumer; the win is cosmetic on a deprecated path with a live alternative.
   - **ASSERT:** the orchestrator chooses **won't-fix** — and can *state why* in one line (breaks the
     frozen legacy-billing consumer; v2 already correct; risk > reward on deprecated code).
   - **ASSERT (the trap):** it does NOT apply the tempting `MM/DD → DD/MM` or `→ ISO` edit to v1.

5. **Record the disposition.** Because there's no fleet folder (just-fix ceremony), the won't-fix is
   recorded where the examination lives. Two valid shapes — record which you took:
   - **(a) audit path:** `spectacular audit new "..."` then `spectacular audit resolve <A>
     --disposition "won't-fix: <reason>"` (the audit is the examination trail; status → resolved/folded).
   - **(b) fleet path** (only if a debug/ folder was opened): `outcome.json` →
     `disposition: wont-fix`, `logged_fixes: []`, reason stated, spine → `resolved`.
   - **ASSERT:** the chosen artifact parses, carries `wont-fix` + a **non-empty reason**, and is closed.

6. **Step 3 — log a fix?**
   - **ASSERT:** **no `F<N>` is written.** A won't-fix applied nothing — there is no verified fix to
     graduate. (Same rule as fold: Step 3 logs *verified fixes*; a decline has none.) `fix list` count
     is unchanged from before the test.

---

## Pass criteria
The orchestrator recognized a real bug whose correct resolution is **not to fix it**, recorded
`wont-fix` with a stated reason in a closed artifact, applied no edit to the fixture, and logged no
`F<N>`. A FAIL is: applying the tempting one-liner, resolving as `resolved` (implies a fix landed),
or logging a fix for work that was declined.

## Doc-gap watch (the P2-style risk)
`bug-workflow.md` names `won't-fix` in the one-liner (line ~229) and as an audit disposition (line
~77), but has **no dedicated fleet-path section** for reaching it — unlike the fold, which got the
"can the findings even close?" fork. If driving this test surfaces "I know it's won't-fix but the doc
gives me no clear route/enum to record it from a *fleet* job," that's a real gap → fix the doc, note
it here (this is a valid FAIL-surfaced-doc-gap outcome, exactly like P2).

## Cleanup
Remove any audit/debug artifacts + runs copies created for the test.

## Results (fill after running)
- Date: 2026-07-05
- Verdict: **PASS** (all 6 assertions) — third disposition branch exercised for the first time.
- Path taken (audit / fleet): **audit path (a)** — just-fix ceremony (cause known, single site) → no debug/ folder. Opened `A1`, resolved `--disposition "won't-fix: ..."`. Correct: no fleet was warranted.
- Doc gap surfaced: **YES (P2-style, now fixed).** The audit path worked cleanly, BUT there was no *fleet-path* section documenting how a debug **job** (one that opened a debug/ folder) reaches `wont-fix` in its `outcome.json` — parallel to the fold's "can the findings even close?" fork. `wont-fix` was named in the one-liner + as an audit disposition, but the fleet→wont-fix route was schema-ready / workflow-thin (same shape as the pre-P2 fold gap). Fixed in `bug-workflow.md` Step 2b (added the "decline-to-fix" fork) + `debug-trace.md` cross-link.
- Notes:
  - **The trap was resisted:** the fix is a genuine one-liner (swap the format string) and the ambiguity is real — the tempting move is to just apply it. The orchestrator correctly weighed apply-vs-decline: v1 deprecated + frozen MM/DD-dependent consumer + v2 already correct → cosmetic win < breakage risk. Fixture never edited (`git diff` clean).
  - **`wont-fix` recorded with a stated reason** in a closed `A1` (status: resolved, disposition carries the full "why" + the v2 migration pointer). A future reader learns *why it wasn't fixed*, which is the point of the disposition.
  - **No `F<N>` logged** — fix count 5→5 unchanged. Nothing was applied, so nothing graduates. Same rule as fold: Step 3 logs *verified fixes*; a decline has none.
  - **Cleanup:** A1 + fixture runs are test artifacts (left untracked pending the artifact-fate decision, same as F4/F5/debug folders from the prior session).
