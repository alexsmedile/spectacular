# Debug-fleet pipeline tests

Integration tests for the **whole debug pipeline** — not single-agent judgment (that's
`tests/agents/`). These exercise the *orchestrator's* choreography: trace-folder setup, briefs out,
agents back, the resolve→ledger graduation, and the routing forks.

**Run these next session, after `/reload-plugins` + `/reload-skills`** so the current (edited) agent
defs and `bug-workflow.md` are live. They can't be run inline — they need the real registered agents.

## What each scenario proves

| # | Scenario | Gap it closes |
|---|---|---|
| P1 | `resolve-to-ledger/` | **resolve→ledger graduation** — a job runs, resolves, and graduates to `fixes/F<N>` + `audit/A<N>` with cross-links + `debug_job` back-link |
| P2 | `needs-planning/` | **fix needs a request, not a Fixer** — Investigator finds a big change; orchestrator routes to the request lifecycle, not a fan-out |
| P3 | `researcher-run/` | **Researcher actually runs** — an external-smelling bug; debug-researcher returns a verdict with citations |
| P4 | `concurrent-fixers/` | **concurrent Fixers write disjoint `fix-NN.json`** — 3-way fan-out, no trace collision, spine collects all three |
| P5 | `full-pipeline/` | **end-to-end** — open bug → Investigator → plan → fan out → verify → resolve → log, one clean run |
| P6 | `wont-fix-disposition/` | **won't-fix disposition** — a real, one-line-fixable bug the orchestrator deliberately declines (deprecated path, frozen consumer); records `wont-fix` + reason, logs no `F<N>` |
| P7 | `malformed-return/` | **inbound safety valve** — Investigator returns an honest under-determined block; orchestrator's well-formedness + symmetry backstop fires, no fix fabricated |

## Convention (same as `tests/agents/`)

- `fixtures/` — pristine bug code, NEVER edited. Agents copy to a scratch `runs/` dir.
- Each scenario has a `RUNBOOK.md` — numbered orchestrator steps + a **pass/fail assertion** per step.
- The orchestrator (you, main thread) drives; agents are spawned per the runbook.
- A scenario **passes** only if every assertion holds. Record results at the bottom of each RUNBOOK.

## Order to run

P3 and P4 are independent — run in any order. P1 depends on nothing but proves the back half. **P5
is the capstone — run it last**, it re-exercises P1+P4 in one flow. P2 is a pure routing check (no
fix applied), quickest to run first as a warm-up.

Suggested: **P2 (warm-up) → P3 → P4 → P1 → P5 (capstone).**
