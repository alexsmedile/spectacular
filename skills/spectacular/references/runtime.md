# Handoff and Autopilot

## Autopilot is explicit and non-default

Never assume it. When the owner turns it on, bind the charter to:

- the exact Mission activation fingerprint
- Objective and claim scope
- Contract and Git baseline
- allowed operator actions
- effects that still require the owner
- budgets and checks
- expiry, stops, recovery
- the return destination

State how resources are actually enforced, as one of:

| Level | Means |
|---|---|
| `hard` | independently verified measurement, and real cancellation |
| `observed` | measured and reported, but not enforced |
| `unsupported` | not measured at all |

Only claim `hard` when both the measurement and the cancellation are verified.

## Promote before delegating

```bash
spectacular objective promote <mission-ref>/<objective-ref> --json   # e.g. M7/O2
```

Promote an inline Objective to its own file before independent delegation. It
lands at `.spectacular/missions/<slug>/objectives/<ref>-<slug>.md` and keeps its
identity. The file then carries the exact:

- outcome and claims
- dependencies and inputs
- semantic and mechanical scope
- authority and stops
- return contract

Accountability stays with the Mission owner. A host task or thread is only a
destination pointer — it owns nothing.

## Fan out sparingly

Delegate only cohesive mid-to-long work whose claim ownership is disjoint.

Avoid:

- tiny sessions
- recursive critic loops
- repeated full reviews

Finish working code, run focused checks, and batch compatible review at the
Mission's frozen review level.

## What the receiver returns

- status and actor
- final baseline and result
- changed files
- checks that ran
- native-provider receipts
- Evidence
- remaining Gaps
- budget use
- recovery point
- one next action, or one owner gate

**The receiver never** changes Mission criteria, declares Evidence sufficient, or
gains provider permission it did not already have.
