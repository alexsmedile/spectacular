---
description: The review sweep — a read-only fleet audit that cross-checks each request's claims against its code, tests, and evidence, then hands off findings. Never promotes.
when_to_use: spectacular sweep [<slug>]; user asks to audit/review the request fleet, check what's really done, or catch stale evidence / duplicate planned work.
---

# Review sweep — audit the fleet, promote nothing

Triggered by: `spectacular sweep` (whole fleet) or `spectacular sweep <slug>` (one request).

The sweep answers one question per request: **does the claimed state match the evidence?** It is the cheap, repeatable middle layer between "TASKS are ticked" and the verify walk — it *feeds* the walk ([[verify]]), never replaces it.

**Two hard rules:**

1. **The sweep never promotes.** No `advance`, ever — not even `--to review`. It may *propose* moves; the human confirms.
2. **The orchestrator is the only mutator** (fleet contract). Auditor agents read and return findings; the orchestrator persists them (VERIFY-LOG, SESSION.md) and relays proposals.

## Three tiers

| Tier | Requests | Depth | Dispatch |
|---|---|---|---|
| **review** | `status: review` | Full audit (below) | one `request-auditor` per request, parallel |
| **ticked-active** | `status: active` with every v1 box `[x]` | Full audit + may propose `advance --to review` | one `request-auditor` per request, parallel |
| **planned-overlap** | `status: planned` | Quick overlap check only: PLAN Goal/summary vs `specs/index.md` capabilities + `archive/` slugs + `fixes/` signatures — flags work that already shipped | ONE batched auditor call for all planned requests |

Mid-flight `active` requests (open boxes) are deliberately out of scope — findings on unfinished work are mostly noise.

## The full audit (per review / ticked-active request)

The auditor receives a closed brief: the request slug + paths. It reads PLAN.md, TASKS.md, VERIFY-LOG.md (if any), SESSION.md, and the code/tests the deliverables name. It cross-checks:

1. **Claims vs code** — does each `[x]` task have a corresponding artifact in the repo (file exists, function present, behavior implemented)? A ticked box with no artifact is a finding.
2. **Claims vs tests** — do the tests the PLAN/TASKS cite exist and pass? (Run them read-only if cheap; otherwise verify they exist and assert what they claim.)
3. **Evidence freshness** — for each `[manual]`/`[observe]` ✓ row in VERIFY-LOG: is the `against:` stamp still current (has the request's code moved past it)? Stale → recommend a `pending-reverify` flip. Missing stamp → finding. An old ✗ row is **not** a current bug and an old ✓ is **not** current proof — the stamp decides ([[verify]] § VERIFY-LOG shape).
4. **PLAN coherence** — do `## Decisions` entries match what was actually built (the walk's coherence pass, run cheaply)?
5. **Blockers** — anything that would stop the verify walk from passing today.

**Return shape (findings block)** — the auditor returns exactly this, no prose around it:

```md
SLUG: <slug>
VERDICT: clean | findings | blocked
FINDINGS:
- [claims|tests|evidence|coherence|blocker] <one-line finding — file:line or row quoted>
PROPOSALS:
- <e.g. "advance --to review" | "flip row X to pending-reverify" | none>
NEXT-AGENT:
- <1-3 imperative lines for whoever picks this request up next>
```

## The loop (orchestrator)

1. `spectacular status --json` → partition the fleet into the three tiers.
2. Dispatch: one auditor per review/ticked-active request (parallel, single message) + one batched planned-overlap auditor. The auditor agent def pins a small/fast model — dispatch is cheap by design; don't inline the reads into your own context.
3. Collect findings blocks. For each audited request:
   - Append a **sweep entry** to `requests/<slug>/VERIFY-LOG.md` (shape below; scaffold the file if absent).
   - Update `SESSION.md § Next actions` with the NEXT-AGENT lines.
   - Apply confirmed `pending-reverify` flips (append `⟳` rows — never edit old rows; the log is append-only).
4. Relay all PROPOSALS to the human in one summary. Execute only what they confirm.
5. Planned-overlap findings → report as "possible duplicate of <shipped capability / archived slug>"; proposing to drop or re-scope a planned request is always human-confirmed.

## Sweep entry shape (VERIFY-LOG)

Interleaves with walk entries in the same append-only file:

```md
## <YYYY-MM-DD HH:MM> — sweep (<C> clean, <F> findings, <B> blockers)

- ✓ [claims] all M1-M3 boxes map to shipped artifacts
- ✗ [evidence] "install flow" ✓ row stamped against b41; code moved (b43) → flipped pending-reverify
- ✗ [tests] TASKS cites tests/foo.test.sh — file missing
**Proposals:** <relayed | none>   **Next-agent:** see SESSION.md
```

## Related

- [[verify]] — the walk this sweep feeds; VERIFY-LOG shape + `against:` stamps
- [[verify-authoring]] — when a request should have a VERIFY.md at all
- [[active-request]] — SESSION.md template + handoff rhythm
- `agents/request-auditor.md` — the auditor's closed contract
