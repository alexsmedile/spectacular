# Scenario D — new findings invalidate an active plan

You are operating the Spectacular skill. The request `webhook-retries` is
`active`. The user says:

> "We were wrong about the retry bug. Update the plan with what we learned
> today."

What was learned today (from the user + logs):

- The plan's diagnosis (§2) blamed the exponential-backoff calculator for
  skipping retries. Today's tracing shows the calculator is correct — retries
  are *scheduled* fine. The real cause: the worker pool drops scheduled jobs
  when it restarts on deploy, and deploys happen daily. So the fix direction
  changes: persist the retry queue across restarts (Redis-backed), not fix the
  backoff math.
- Milestone M2 ("backoff calculator hardened, property-tested") is therefore
  pointless as scoped.
- A new design choice was made in today's session: Redis-backed persistent
  queue over a DB-table queue, because Redis is already in the stack for
  rate-limiting and the DB table would need a poller.

Here is the current `requests/webhook-retries/PLAN.md`:

```md
---
status: active
updated: 2026-07-03
summary: "Fix webhook retries silently skipping — suspected backoff calculator bug"
related:
  - PRD.md
---

# Plan — webhook-retries

## 1. Goal

Webhook deliveries that fail are reliably retried until delivered or
dead-lettered — no silent skips (PRD G4: no customer-visible event loss).

## 2. Diagnosis

The backoff calculator (`retry/backoff.ts`) appears to skip the 3rd+ retry
when the computed delay exceeds 2^31 ms due to integer overflow. Logs show
retries 1–2 firing and 3+ never scheduled.

## 3. Milestones

- M1 — Repro harness: failing webhook target that forces retries 1–5
- M2 — Backoff calculator hardened, property-tested for overflow
- M3 — Dead-letter queue after max retries, visible in admin UI

## 4. Tasks

See TASKS.md.

## 5. Dependencies

- none

## 6. Validation

- M1 — run: harness script exits 0 having observed retries 1–5 scheduled
- M2 — run: property test suite for backoff passes 10k cases
- M3 — observable: exhausted webhook appears in admin dead-letter view

## 7. Deliverables

- Fixed backoff module + property tests
- Dead-letter queue + admin view
```

Produce the complete updated PLAN.md exactly as you would write it to disk,
following ONLY your reference docs. Then note (one line each) any other
file/frontmatter updates your docs require.
