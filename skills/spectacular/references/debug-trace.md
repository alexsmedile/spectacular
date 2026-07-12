---
doc-id: debug-trace
kind: reference
summary: "Traceability for a live debug job: one folder per job under .spectacular/debugs/, one JSON artifact per agent turn — all persisted by the orchestrator from returned blocks. Schemas + the single-writer rule. In-flight process state, distinct from the fixes/audit ledger."
status: active
---

# Debug trace — in-flight process state for a debug job

A debug job that runs the fleet (Investigator → orchestrator → Fixer/Researcher) needs a place to
record *what happened* while it happens — so a resuming agent, the human, or a `doctor` check can
answer: **what bug, what did each agent find/do, what's the state now, what's the outcome.** That
place is `.spectacular/debugs/<job-slug>/`, holding one JSON artifact per agent turn.

## debugs/ vs audits/ vs fixes/ — the pipeline and the two summaries it produces

`debugs/` is the raw pipeline *while a job runs* (the workbench — tolerates open + messy);
`audits/` + `fixes/` are the two permanent summaries distilled from it at resolution (the library —
requires resolved + clean):

| Folder | Answers | Scope | Lifetime | Curation |
|---|---|---|---|---|
| **`debugs/<job>/`** | "what's happening **right now** on this live job?" | one in-flight job | **kept as trace** (marked resolved, never pruned) | raw — verbose JSON, per-agent, may be unresolved |
| **`audits/A<N>.md`** | "what did we find?" (the **examination**) | any deep examination, not just bugs | permanent | written only if findings are worth keeping |
| **`fixes/F<N>.md`** | "have we solved this before?" (the **remedy**) | solved, reusable problems | permanent | written only if reusable |

The two summaries cross-link (`A<N>` ↔ `F<N>`) and both carry `debug_job: <job-slug>` back to the
raw trace that produced them.

## Trace vs ledger — the write rule

| Layer | Where | Writer | When |
|---|---|---|---|
| **Debug trace** | `.spectacular/debugs/<job>/*.json` | **orchestrator** — persists each agent's returned block | live, as each agent returns |
| **Ledger** | `fixes/F<N>.md`, `audits/A<N>.md` | **orchestrator only** (CLI verbs) | at resolution, distilled from the trace |

The invariant is **total**: agents only *return* their block (the Agent-tool result); the
orchestrator persists it as the leaf artifact and stays the single writer of both trace and ledger
— `F<N>`/`A<N>` go through `spectacular fix new` / `audit new`. (Agents writing their own leaf via
heredoc was the old contract — dropped because it paid the same content twice in agent output.) All of `.spectacular/` is **committed to git**, so a live job survives a
session loss, is team-visible, and a fresh session can resume from it.

## Folder layout

```
.spectacular/debugs/<job-slug>/
├── job.json              # the spine — orchestrator owns + updates
├── investigation.json    # Investigator's findings (0 or 1)
├── research/             # Researcher verdicts (0..n — one per external question)
│   └── research-01.json
├── fixes/                # one per fanned-out fix (0..n)
│   ├── fix-01.json
│   └── fix-02.json
└── outcome.json          # final disposition + graduation to the ledger
```

**The folder's *shape* tells the story before you read anything:** no `investigation.json` → it was
just-fix, no discovery needed. Three files in `fixes/` → a 3-way fan-out. A `research/` entry → it
smelled external. Missing `outcome.json` → still in flight.

## Slot assignment — no collision on fan-out

The orchestrator owns slot assignment (it knows the job folder and the fan-out count) **and** the
write: when it fans out 3 Fixers it earmarks `fixes/fix-01.json`, `fixes/fix-02.json`,
`fixes/fix-03.json`, and as each Fixer's block returns, persists it to that slot. Agents never
touch the trace folder; collision is impossible by construction.

## Schemas

All timestamps ISO-8601 UTC. All paths relative to the job folder. `null` for a slot that doesn't
apply (not omitted — presence of the key with `null` is explicit).

### `job.json` — the spine (orchestrator writes + updates)

```json
{
  "id": "cache-staleness-2026-07-05",
  "status": "investigating",
  "symptom_class": "wrong_behavior",
  "symptom": "get_user returns stale name after update_name",
  "reported_at": "2026-07-05T12:00:00Z",
  "reporter": "user",
  "ceremony": "audit-first",
  "brief": {
    "symptom": "get_user returns stale name after update_name; expected the written value",
    "where_to_look": "users.py cache path — update_name / get_user",
    "done_means": "root cause + suspected site established"
  },
  "artifacts": {
    "investigation": "investigation.json",
    "research": ["research/research-01.json"],
    "fixes": ["fixes/fix-01.json", "fixes/fix-02.json"]
  },
  "timeline": [
    {"at": "2026-07-05T12:01:00Z", "phase": "investigation", "agent": "debug-investigator", "result": "root-cause-found"},
    {"at": "2026-07-05T12:05:00Z", "phase": "fixing", "agent": "debug-fixer", "result": "applied", "artifact": "fixes/fix-01.json"}
  ],
  "outcome": "outcome.json"
}
```

- `status`: `investigating | researching | planning | fixing | verifying | resolved | folded | wont-fix`.
- `symptom_class`: how the bug *surfaces* — `test_failure | runtime_error | wrong_behavior | build_error | performance | unknown`. Orthogonal to routing (which agent) — this is the greppable symptom axis ("show me all perf bugs"). `unknown` is honest; set it at triage and refine if it becomes clear.
- `reporter`: `user | while-coding`.
- `ceremony`: `just-fix | audit-first` (from the Step-1 gate).
- `brief`: the investigation brief the orchestrator handed the Investigator (`symptom` / `where_to_look` / `done_means`) — persisted because the spawn prompt is gone once the window closes. Lets a resuming orchestrator check `investigation.json`'s STATUS against what was actually asked (the symmetry check). `null` for a just-fix job with no investigation.
- `artifacts`: forward-index to every artifact file (arrays grow as agents write).
- `timeline`: append-only; one entry per agent turn or status change. The audit trail.
- `outcome`: path to `outcome.json`, or `null` until resolved.

### `investigation.json` — Investigator findings (orchestrator persists from the returned block)

Mirrors the agent's findings block, one-to-one:

```json
{
  "agent": "debug-investigator",
  "at": "2026-07-05T12:01:00Z",
  "status": "root-cause-found",
  "reason": null,
  "symptom": "get_user returns stale 'alice' after update_name(1,'bob'); deterministic",
  "root_cause": "update_name writes _store but never invalidates _cache; a key read before the update returns the stale cached value",
  "hypotheses": [
    {"rank": 1, "claim": "cache not invalidated on write", "evidence_for": "trace + ran 8x", "evidence_against": null}
  ],
  "ruled_out": [
    {"hypothesis": "TTL expiry race", "evidence": "cache has no TTL config anywhere; grep 'ttl' → 0 hits"}
  ],
  "suspected_sites": ["users.py:23 (update_name)"],
  "plausible_solutions": [
    {"approach": "pop-on-write", "tradeoff": "simplest; next read repopulates"},
    {"approach": "write-through", "tradeoff": "one fewer store read"},
    {"approach": "centralize _invalidate() helper", "tradeoff": "worth it if more mutators arrive"}
  ],
  "blast_radius": "create() has the same latent gap; safe only because it runs before any read caches",
  "open_questions": [],
  "evidence": "ran read-only 8x, deterministic AssertionError at line 42; trace: get_user caches alice, update_name sets store only, next get_user hits stale cache"
}
```

- `status`: `root-cause-found | hypotheses-only`.
- `reason`: only non-null when `hypotheses-only` — `needs-reproduction | needs-research | needs-decision | needs-more-context`.
- `ruled_out`: hypotheses **tested and eliminated**, each with the evidence that killed it. Empty array means nothing was eliminated — suspicious for anything beyond a trivial bug. The orchestrator copies these into the `audit/A<N>` entry so no future walk re-opens a dead end (see [[bug-workflow]] § Coming back).
- `plausible_solutions`: the solution space (approaches + trade-offs) — **never the literal diff**.

### `research/research-NN.json` — Researcher verdict (orchestrator persists from the returned block)

```json
{
  "agent": "debug-researcher",
  "at": "2026-07-05T12:03:00Z",
  "verdict": "known-platform-bug",
  "confidence": "high",
  "summary": "Library X v2.3 has a documented regression matching this signature; fixed in 2.3.1",
  "evidence": [
    {"title": "X#1234 cache race", "url": "https://github.com/...", "match": "same error string + version 2.3; maintainer-confirmed"}
  ],
  "workaround": "pin X>=2.3.1, or call flush() after write (docs §cache)",
  "next": "upgrade the dep to 2.3.1"
}
```

- `verdict`: `known-platform-bug | genuinely-ours | no-strong-match`.
- `workaround`: string, or `"none found"`.

### `fixes/fix-NN.json` — Fixer result (orchestrator persists from the returned block)

Mirrors the Fixer's output block, plus the `diff` the orchestrator captures itself:

```json
{
  "agent": "debug-fixer",
  "at": "2026-07-05T12:05:00Z",
  "verdict": "applied",
  "site": "users.py:25",
  "brief": {
    "problem": "stale cache after update_name",
    "intended": "read after write returns written value",
    "root_cause": "update_name doesn't invalidate _cache",
    "proposed_fix": "self._cache.pop(uid, None) after the store write",
    "success_criteria": "python3 users.py prints ok, exit 0"
  },
  "changed": [
    {"file": "users.py", "what": "pop uid from _cache after the store write, so the next read repopulates"}
  ],
  "diff": "--- a/users.py\n+++ b/users.py\n@@ ...",
  "test": "none (trivial)",
  "risk": "low",
  "verify": {"check": "python3 users.py", "result": "pass"},
  "bounce_reason": null,
  "ledger": "not-written"
}
```

- `verdict`: `applied | bounced`.
- `changed`: one entry per file touched — `{file, what}` (the human-readable "explain the change" the diff alone doesn't give). Empty array if bounced.
- `diff`: the unified diff — **captured by the orchestrator** via `git diff -- <the files in
  CHANGED>` at collection time (the Fixer's return block carries no diff; its edit lives in the
  working tree). Empty string if bounced.
- `test`: the regression test added/updated (`file:name`), or `"none (<reason>)"` — `trivial | no-framework | brief-didn't-ask`. A test that pins the fixed bug is part of the fix, not scope-widening.
- `risk`: blast radius of the change — `low` (one site, isolated), `medium` (shared helper, few callers), `high` (shared root, wide blast radius). Feeds the orchestrator's fan-out-vs-inline call and flags what to watch on verify. `null` if bounced.
- `verify.result`: `pass | fail | not-run`.
- `bounce_reason`: string when bounced, else `null`.
- `ledger`: always `"not-written"` — the Fixer never writes the ledger.

### `outcome.json` — final disposition (orchestrator writes)

The bridge from trace to permanent ledger:

```json
{
  "at": "2026-07-05T12:06:00Z",
  "disposition": "resolved",
  "fix_ids": ["fixes/fix-01.json"],
  "logged_fixes": ["F4"],
  "audit": null,
  "note": "single-site cache invalidation; logged as reusable (cache-coherence class)"
}
```

- `disposition`: `resolved | folded-into-request | wont-fix`.
- `logged_fixes`: the `F<N>` ids the reusable remedy graduated into — empty if nothing was reusable
  (`log-only-verified-reusable`).
- `audit`: the `A<N>` id if the examination was worth keeping, else `null` (a job may earn neither
  `audit` nor `logged_fixes`; the trace itself is always retained regardless).
- `folded-into-request` also carries a `request` slug. Reached from the fleet when the Investigator's findings are design work too big to close into fix slots — the route is [[bug-workflow]] Step 2b's "can the findings even close?" fork. (A disposition in this schema is only usable once the *workflow* wires a path that reaches it — schema-ready ≠ workflow-ready.)
- `wont-fix` carries a stated `reason`. Reached when the findings *do* close into a clean fix but applying it is the wrong call (frozen consumer, deprecated path with a live alternative, a deliberate trade-off, cost > symptom) — the route is [[bug-workflow]] Step 2b's "should the fix even be applied?" fork. `logged_fixes` stays empty: a decline applies nothing, so nothing graduates. The `reason` is the durable why-not, so the bug isn't re-litigated later. (Just-fix ceremony records the same decline on the audit via `audit resolve --disposition "won't-fix: ..."` — no `outcome.json` when no folder was opened.)
- Both `A<N>` and `F<N>` (when written) cross-link each other and carry `debug_job` back to this
  folder.

> **These spines are hand-written by the orchestrator (no `debug` write-verb) — so `spectacular
> doctor debug` validates them.** It checks `job.json` `status` and `outcome.json` `disposition`
> against the enums above, and enforces the invariant that a `wont-fix`/`folded-into-request` job
> logs no `F<N>` (`logged_fixes: []`). An LLM drifts on closed enums — e.g. leaking a `reason` value
> (`needs-more-context`) into the `status` slot — so this is the guardrail that catches it at check
> time instead of a resume failure. If you add a status/disposition value here, add it to
> `check_debug` in `cli/spectacular` too, or doctor will flag valid spines.

## Lifecycle

1. **Open** — orchestrator scaffolds `debugs/<job-slug>/` + writes `job.json` (`status: investigating`
   or `fixing` for a just-fix bug that skips investigation).
2. **Returns append** — as each agent returns, the orchestrator persists its block to the assigned
   slot, appends a `timeline` entry + updates `status` + the `artifacts` index in `job.json`.
   Jobs may sit **unresolved** here — the trace holds mid-flight state with no final clarity yet.
   That's fine; debugs/ tolerates open.
3. **Resolve — summarize into the library (each only if earned).** When clarity exists, the
   orchestrator writes `outcome.json` and distills the trace into the two permanent summaries:
   - **`audits/A<N>.md`** — the examination — written **only if the findings are worth keeping**
     (a non-obvious cause, a real investigation). A shorter, human-facing distillation of
     `investigation.json`, *not* a copy of it.
   - **`fixes/F<N>.md`** — the remedy — written **only if reusable** (`log-only-verified-reusable`).
   Cross-link them (`A<N>` ↔ `F<N>`) and stamp both with `debug_job: <job-slug>`. A trivial job may
   earn neither. Both go through the CLI verbs (`spectacular audit new` / `fix new`) — the ledger
   stays single-writer.
4. **Keep the trace in place.** The resolved `debugs/<job>/` is **not pruned** — `job.json` flips to
   `status: resolved` with `audits`/`fixes` back-links, and the folder stays as the full "how the job
   ran" record. The library is the distilled summary; debugs/ is the durable raw feed behind it.

A `just-fix` 1–2 line bug that never runs the fleet **needs no debug folder** — it's fixed inline
and, if reusable, logged straight to `F<N>`. The trace exists for jobs that fan out; ceremony scales
with uncertainty (`ceremony-matches-uncertainty`).

**Related:** [[bug-workflow]] (the flow that writes these), [[fixes-rules]] / [[audit-rules]] (the
ledger they graduate into), [[soft-db-index]] (where knowledge belongs).
