# P4 — concurrent Fixers, disjoint trace writes

**Proves:** the fan-out invariant holds — 3 independent, closed, **disjoint-file** fixes spawn 3
`debug-fixer`s in parallel, each writes its own `fixes/fix-NN.json` to the slot the orchestrator
assigned, **no trace collision**, and the orchestrator collects all three into the spine. This is the
≥3-independent-closed-disjoint rule (Step 1b) exercised for real, and the fan-out trace-slot
assignment (`debug-trace.md`) proven.

**Bugs (3 disjoint files):**
- `fixtures/alpha.py` — off-by-one in `last_n`
- `fixtures/beta.js` — `> 18` should be `>= 18`
- `fixtures/gamma.py` — `per_page` default 0 should be 10

Each is closed (root cause obvious, single site) and in a different file → the fan-out case.

---

## Steps (orchestrator drives)

1. **Confirm all three bugs.** Run each fixture; each fails.
   - **ASSERT:** all 3 exit non-zero (real repros).

2. **Triage + Step 1b decision.** 3 independent, closed, disjoint-file fixes → **fan out** (this is
   exactly the ≥3 threshold).
   - **ASSERT:** you choose fan-out (not inline) — the table's ≥3-independent row.

3. **Open the job (Step 1c).** Scaffold `.spectacular/debug/multi-clamp/`, write `job.json`, and
   **assign three disjoint slots**: `fixes/fix-01.json` (alpha), `fixes/fix-02.json` (beta),
   `fixes/fix-03.json` (gamma). Record each brief.
   - **ASSERT:** `job.json` has all 3 slots pre-assigned in `artifacts.fixes`.

4. **Fan out — spawn 3 `debug-fixer`s in ONE message** (parallel), each with: its closed 5-slot brief,
   its fixture→runs copy path, and its assigned `fix-NN.json` trace path.
   - **ASSERT:** all 3 spawn concurrently (single message, 3 Agent calls).

5. **Collect.** Each returns `VERDICT: applied` + diff + a `RISK` + `VERIFY: pass`.
   - **ASSERT:** 3 blocks returned, all `applied`, all `VERIFY → pass`.
   - **ASSERT (the collision test):** 3 distinct files exist — `fix-01.json`, `fix-02.json`,
     `fix-03.json` — each parses, each describes its OWN fix (alpha in 01, beta in 02, gamma in 03).
     **No file was overwritten; no two Fixers wrote the same path.**

6. **Update the spine.** Append 3 timeline entries; `artifacts.fixes` lists all three.
   - **ASSERT:** `job.json.timeline` has 3 fix entries; `status` advanced to `verifying`/`resolved`.

7. **Confirm no cross-contamination.** Diff each `runs/` copy vs its fixture.
   - **ASSERT:** each fix touched ONLY its own file — alpha's Fixer never touched beta/gamma, etc.
     (Disjoint-file isolation held; this is what makes parallel safe.)

---

## Pass criteria
3 Fixers ran in parallel, each wrote a distinct `fix-NN.json` (no collision), each fix is isolated to
its own file, and the spine collected all three. **The disjoint-slot assignment prevented any trace
write race.**

## Cleanup
Remove `.spectacular/debug/multi-clamp/` and the `runs/` copies.

## Results (fill after running)
- Date: 2026-07-05
- Verdict: **PASS** (all assertions)
- All 3 trace files distinct + correct? **YES** — fix-01→alpha.py:6, fix-02→beta.js:4, fix-03→gamma.py:4; each parses, each describes only its own fix. No overwrite.
- Notes:
  - Real fan-out: 3 `debug-fixer` subagents spawned in ONE message (parallel), all returned `applied` / `VERIFY: pass` / `RISK: low`.
  - **Collision test held:** disjoint-slot pre-assignment (`artifacts.fixes` seeded before spawn) meant no two Fixers targeted the same trace path. Each wrote its assigned `fix-NN.json`, none clobbered another.
  - **Isolation held:** diff of each `runs/` copy vs its pristine fixture shows exactly ONE changed line, in that file only — alpha's Fixer never touched beta/gamma, etc. Disjoint-file is what makes parallel safe, confirmed empirically.
  - Spine collected all 3 (timeline has 3 fix entries), status → resolved. Not graduated to ledger — these are mechanical one-liners (off-by-one / operator / default), not reusable-signature footguns.
  - Every Fixer independently proposed dropping the now-stale `# bug:` comment alongside its fix — consistent "leave no stale marker" behavior, no prompting.
