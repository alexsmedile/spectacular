# Scenario E — run the verification walk

You are operating the Spectacular skill. The user says:

> "spectacular verify rate-limit-cache"

The request `rate-limit-cache` is in `review`. No VERIFY.md exists —
validation lives in PLAN § 6. Here is `requests/rate-limit-cache/PLAN.md`:

```md
---
status: review
updated: 2026-07-06
summary: "Add per-API-key rate limiting backed by the shared cache"
related:
  - PRD.md
---

# Plan — rate-limit-cache

## 1. Goal

API abuse from a single key cannot degrade other tenants (PRD G5: no
cross-tenant performance impact) — per-key rate limiting at the gateway.

## 2. Constraints

- Use the existing shared cache (Redis); no new infra.
- Limits configurable per key tier without redeploy.

## Decisions

- Chose a sliding-window counter over a fixed-window counter — because
  fixed windows allow 2x burst at window boundaries, and burst fairness
  is the whole point of PRD G5.

## 3. Milestones

- M1 — Gateway enforces a per-key limit (429 over limit)
- M2 — Limits configurable per tier at runtime
- M3 — Overload isolation proven under load

## 4. Tasks

See TASKS.md.

## 5. Dependencies

- none

## 6. Validation

- M1 — run: tests/rate-limit.test.sh exits 0 (101st request in a minute → 429)
- M2 — run: changing tier limit via config API takes effect within 10s, test asserts it
- M3 — run: load test — key A at 10x limit; p95 latency for key B stays < 150ms
```

You walk the checks. Simulated results of running each command (treat these
as the real, fresh outputs):

- `tests/rate-limit.test.sh` → exit 0. Test log excerpt: "window=60s fixed;
  requests 1-100 pass, 101 → 429. PASS". (The implementation in
  `gateway/ratelimit.ts` uses `Math.floor(now/60000)` bucketing — a classic
  fixed 60s window keyed on wall-clock minute.)
- Config-API test → exit 0, limit change applied in 4s. PASS.
- Load test → exit 0, key B p95 = 96ms while key A hammered at 10x. PASS.

Run the verification walk exactly as your reference docs prescribe, from
start to finish, including everything the docs say happens at the end of a
passing walk. Show the walk record you'd write and exactly what you would
say to the human.
