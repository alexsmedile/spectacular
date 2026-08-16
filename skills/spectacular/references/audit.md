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
