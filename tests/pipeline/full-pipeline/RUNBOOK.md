# P5 — full pipeline (capstone)

**Proves:** the whole debug arc runs end-to-end on one open bug — trigger → Step 0 → ceremony gate →
open job → Investigator discovers → orchestrator plans → apply → verify → resolve → log — with every
handoff well-formed and every honesty valve available. Re-exercises P1 (resolve→ledger) and touches
the plan-then-fix seam. **Run this last.**

**Bug:** `fixtures/cart.py` — `remove_item` doesn't recompute `_total`; `get_total` returns a stale
cache. OPEN (cause + site not given). Shared root across mutators (add/remove/clear) → planning judgment.

---

## Steps (orchestrator drives — this is the real workflow, not a shortcut)

1. **Trigger + confirm.** Treat as a user bug report ("cart total wrong after remove"). Run the
   fixture.
   - **ASSERT:** it fails (`get_total()` returns 15, expected 10). Repro in hand, `reporter: user`.

2. **Step 0 — seen it before?** Grep `.spectacular/fixes/` signatures.
   - **ASSERT:** no prior match (or if P1's discount fix is still around, it's unrelated). Proceed.

3. **Step 1 — ceremony gate.** Cause not yet known + possible multi-site → **audit-first**.
   - **ASSERT:** classified `symptom_class: wrong_behavior`, ceremony `audit-first`.

4. **Step 1c — open the job.** Scaffold `.spectacular/debug/cart-stale-total/`, write `job.json` with
   the `brief` (Symptom / Where to look / Done means), `status: investigating`.
   - **ASSERT:** `job.json` parses, has `brief` + `symptom_class`.

5. **Spawn the Investigator** with the brief + `investigation.json` trace path.
   - **ASSERT:** returns a well-formed block (`STATUS` on-enum, non-empty Root cause + Suspected sites
     + EVIDENCE). Expected `root-cause-found`.
   - **ASSERT (the diagnosis is right):** root cause = "_total not recomputed on remove"; Suspected
     sites name `remove_item` AND flag the **shared root / blast radius** across mutators (add is fine,
     clear is latent). Plausible solutions describe the space (recompute-on-mutation vs derive-total-
     on-read) **without prescribing the literal diff.**
   - **ASSERT (channel):** the block IS the tool result; `investigation.json` exists + mirrors it.

6. **Symmetry check.** Compare returned STATUS to the brief's Done means.
   - **ASSERT:** Done means was "root cause + site" and STATUS is `root-cause-found` → met, proceed to plan.

7. **Step 2b — orchestrator plans.** Turn findings into a closed fix. Decide the approach among the
   Investigator's plausible solutions (e.g. recompute `_total` in `remove_item`, or make `get_total`
   derive from `_items`). Write the exact 5-slot brief. This is YOUR job, not the Investigator's.
   - **ASSERT:** you produce a closed brief with all 5 slots concrete — a Fixer could apply it with no
     judgment call.
   - **DECISION POINT:** if you fix the shared root as ONE change → 1 fix, self-serve inline (Step 1b:
     1–2 fixes → inline). Record which path you took.

8. **Apply the fix** (inline, 1 site after choosing the root approach). Copy fixture → `runs/`, edit,
   run the check.
   - **ASSERT:** the check prints `ok`, exit 0. Fix is `low`/`medium` risk (shared helper) — verify
     the OTHER mutators too (add still works, and the latent clear() gap is noted or fixed).
   - **ASSERT (risk/verify match):** if you report `medium` risk, your verification exercised more than
     the single symptom (per the risk-scaled verify rule).

9. **Step 3 — log if reusable.** Cache-coherence-on-mutation is a reusable footgun → `spectacular fix
   new` with a signature, `debug_job: cart-stale-total` back-link.
   - **ASSERT:** `F<N>` created; `fix list` shows it; back-link present.

10. **Resolve.** Write `outcome.json` (`disposition: resolved`, `logged_fixes: [F<N>]`), flip
    `job.json` → `status: resolved`. Trace folder kept.
    - **ASSERT:** `outcome.json` parses, `disposition == resolved`; debug folder still present.

11. **Prove the loop.** Grep `fixes/` for the signature.
    - **ASSERT:** the new `F<N>` is findable — a future occurrence would hit Step 0.

---

## Pass criteria
Every assertion holds across all 11 steps. The capstone passes only if: **the Investigator discovered
the right root (with blast radius, no literal diff), the orchestrator planned the fix, the fix applied
+ verified at the right depth for its risk, and the job resolved + graduated to the library with a
working back-link.** Any broken handoff (malformed block, symmetry miss, band-aid fix, unlogged
reusable lesson) is a FAIL with a specific pointer.

## Cleanup
Remove the debug folder, runs copies, and test `F<N>`.

## Results (fill after running)
- Date:
- Verdict: PASS / FAIL
- Which fix path (root inline / other):
- Any handoff that needed doc changes:
- Notes:
