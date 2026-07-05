# P2 — fix needs planning (a request), not a Fixer

**Proves:** when investigation reveals the "fix" is real design work (multi-step, schema/policy
change), the orchestrator routes to the **request lifecycle** — NOT a fan-out. Tests the
fleet↔request bridge.

**Known risk:** this bridge may be **under-specified** in `bug-workflow.md` (Step 1 mentions the
audit→plan→request path but the fleet path — Investigator findings → request — is not wired). If the
orchestrator has no clear route, **this scenario surfaces the gap** — a FAIL here means "go add the
bridge," which is a valid outcome, not a bug in the test.

**Bug:** `fixtures/auth.py` — "expired sessions still authenticate," but there's no expiry mechanism
at all. Fixing = designing session TTL across create/read paths. Not mechanical.

---

## Steps (orchestrator drives)

1. **Confirm the symptom.** `python3 tests/pipeline/needs-planning/fixtures/auth.py`
   - **ASSERT:** it prints `authenticated: 1` — the session never expires (no mechanism).

2. **Triage.** Ceremony gate: root cause NOT a single site (structural absence) → **audit-first**.
   - **ASSERT:** you classify ceremony `audit-first` and do NOT jump to a Fixer.

3. **Open the job + spawn the Investigator** with a brief (Symptom / Where to look / Done means).
   Give it the trace path `investigation.json`.
   - **ASSERT:** Investigator returns `STATUS: root-cause-found` (or hypotheses-only) whose findings
     say the cause is **structural / no-expiry-concept**, and whose Suspected sites span create +
     read paths (not one line). Its Blast radius / Open questions should signal this is bigger than a fix.

4. **THE FORK — orchestrator plans.** Read the findings. Recognize the fix is not closeable into the
   five slots (Proposed fix would be "design session expiry," which is a project, not an edit).
   - **ASSERT:** you do **NOT** write a Fixer brief and do **NOT** fan out. Attempting to fill the
     five closed slots should fail the "concrete enough a Fixer applies without judgment" bar.

5. **Route to the request lifecycle.** Per Step 1's audit-disposition "folded into a request": create
   a request (`spectacular new "session-expiry"` or the skill's request-creation path), and set the
   job's `outcome.json` → `disposition: folded-into-request`, `request: <slug>`.
   - **ASSERT:** `outcome.json.disposition == "folded-into-request"` and carries a `request` slug.
   - **ASSERT:** a request scaffold exists under `.spectacular/requests/<slug>/`.

6. **Confirm no fix was logged.** No `F<N>` should be written — nothing was fixed, it was planned.
   - **ASSERT:** `fixes/` has no new entry from this job; `outcome.json.logged_fixes` is empty.

---

## Pass criteria
The orchestrator recognized the non-mechanical fix, refused to fan out, and routed to a request with
`disposition: folded-into-request`. **If the doc gave no clear route for step 5, mark FAIL and file:
"add the fleet→request bridge to bug-workflow.md Step 2b."**

## Cleanup
Remove the test request + debug folder after.

## Results (fill after running)
- Date:
- Verdict: PASS / FAIL / FAIL-surfaced-doc-gap
- Bridge specified in doc? YES / NO
- Notes:
