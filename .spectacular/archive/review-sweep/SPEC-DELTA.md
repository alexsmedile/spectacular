# Spec delta — review-sweep (b31)

### ADDED
- specs/index.md :: **Review sweep (v1.35.0+)** — `spectacular sweep [<slug>]` (skill-only; CLI redirects) audits the request fleet read-only: one small-model `request-auditor` subagent per `review`/ticked-`active` request (claims vs code, cited tests, evidence freshness, PLAN coherence), one batched overlap check over `planned` requests vs `specs/index.md` + archive. Findings append to VERIFY-LOG as dated sweep entries; next-agent handoff to `SESSION.md § Next actions`; never promotes — feeds the verify walk. VERIFY-LOG `[manual]`/`[observe]` rows carry `against: <commit/build> · <identity>` stamps with judgment-flipped `pending-reverify` semantics; `doctor lifecycle` warns on unstamped/pending rows (presence checks only, no git heuristics). See [[review-sweep]].

### MODIFIED
- specs/index.md :: "**`test-verifier`** (apply-only, tests only; run a named check or write a test to a closed spec → honest pass/fail)" -> "**`test-verifier`** (apply-only, tests only; run a named check or write a test to a closed spec → honest pass/fail), **`request-auditor`** (read-only, haiku-pinned; one request's claimed state vs its actual evidence → findings block + next-agent handoff, dispatched by the review sweep)"
