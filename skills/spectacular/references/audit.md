# Audit

Audit one named claim, or one bounded scope. That is the default — do not widen it
without being asked.

## Inspect

Read the canonical sources, not a projection of them. `index.md` and any generated
catalog are navigation aids and prove nothing.

```bash
spectacular mission show  <ref> --json   # frozen envelope + current state
spectacular mission check <ref> --json   # what the validators say today
git log --oneline <baseline>..HEAD       # what actually changed since activation
```

Then work through:

- canonical Mission and specification sources under `.spectacular/`
- activation and source fingerprints
- the Git tree against the recorded baseline
- Evidence: its method, and its freshness
- contrary evidence
- reviewer independence
- authority
- dependencies and Gaps
- stops

### A completed Mission's Contract binding is a freeze point, not a stale pointer

`mission check` reports `contract-drift` when a Mission's `contract.fingerprint`
no longer matches the Contract on disk. On a **completed** Mission this is a
notice, the Mission stays `valid=true`, and **nothing is wrong**. The binding
states which agreement that Mission was executed against, which is what lets you
derive why the work was shaped the way it was.

Read the Contract as it was, not as it is:

```bash
git log -S <bound-fingerprint> -- .spectacular/contracts/<file>.md
```

Never re-point a completed Mission to silence the notice. That overwrites the
freeze point with today's answer and destroys the thing the audit depends on.
Amendments re-point only the live Mission, by design — see `D10-repoint`.

Drift on a **live** Mission is different and is worth investigating: its binding
is supposed to track the Contract it is working against.

## Apply FROST

| | Check | Across |
|---|---|---|
| **F** | Frozen fit | outcome, claims, non-goals |
| **R** | Risk | security, privacy, data, authority, irreversible effects |
| **O** | Operability | failures, diagnostics, recovery, maintenance |
| **S** | System integrity | dependencies, regressions, generation, distribution |
| **T** | Truth of proof | attributable methods, freshness, contrary evidence, no overclaim |

## Report

- which claims are supported
- findings, split into blocking and non-blocking
- what proof is missing
- bounded corrective options

Return the reviewed tree, the claims, the findings, the limitations, and one
verdict: `pass | repair | owner-gate`.

An audit reports in the session by default. It becomes a durable record only when
the Mission's review level asked for one — then it goes through
`spectacular review record <mission-ref> <review.md|-> --json`, and the
independence rule in [close.md](close.md) applies.

## Stay inside the audit

An audit inspects. It does not:

- mutate the lifecycle
- redefine criteria
- pick product trade-offs
- treat a fresh agent as independent Evidence

Batch compatible review and reuse Evidence that has not changed. Do not repeat an
unchanged review loop.
