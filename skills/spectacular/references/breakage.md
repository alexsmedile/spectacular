# Deliberate breakage testing

Use this when: the user asks Spectacular to break software or its structure,
find failure modes, run adversarial scenarios, or prepare a breakage test suite.
This is a Skill route for ordinary work, not a CLI command or Mission activation.

## 1. Learn the target and define the boundary

Inspect the requested projects, relevant instructions, entry points, tests,
commands, persistence, roles, dependencies, and recovery mechanisms. Follow one
representative path from input to stored state or output; use a bounded baseline
check when available. Identify invariants from accepted requirements, interfaces,
and domain rules rather than treating current behavior as the expected result.
For multiple products, map each separately and their integration if applicable.
Distinguish observed behavior, documented expectations, and hypotheses.

Record target revision/version, available environment, affected surfaces,
excluded surfaces, execution budget, and permitted effects. Use existing project
limits; otherwise choose and state bounded time/resource limits before trials.
Missing access still permits a provisional suite; state what must be inspected
or supplied before its cases can run. Resume an existing campaign from its last
recorded attempt rather than replaying successful or destructive actions blindly.
Done when the target, invariant sources, test boundary, and access gaps are explicit.

## 2. Present principles and plausible tests before execution

First write a target-specific principles list, then a ranked list of 12–24
independent adversarial cases, unless the user requests a different count.
Publish both in the conversation or a saved plan linked before trials begin.
If the request is principles/list-only, stop here. If testing is authorized,
continue after presenting the list; preparation adds no approval gate.

Apply these principles to the target:

- Challenge guarantees at boundaries and invalid transitions, rather than
  repeating happy paths with cosmetic input changes.
- Define an observable failure oracle before running each trial. A clean,
  actionable rejection may be a pass; accepting invalid state may be a failure.
- Prefer plausible user mistakes, developer assumptions, and operational faults
  with high impact and inexpensive reproduction; rank risk and cost separately.
- Test integrity and recovery as well as the immediate response. Verify persisted
  state, duplicate effects, residual resources, and the ability to resume work.
- Separate violated requirements from usability concerns and design limitations;
  when intended behavior is unclear, record a hypothesis rather than inventing a bug.
- Use controlled, attributable attempts with reproducible inputs and bounded
  fault injection. Minimize a failing case before escalating complexity.

Use the following as a coverage menu, not a mandatory fixed checklist. Select
relevant combinations and explain exclusions; do not invent unavailable features.

| Surface | Plausible disruption | Observable property to challenge |
|---|---|---|
| Input | Empty, truncated, oversized, malformed, Unicode, boundary values | Validation, precision, actionable errors |
| Commands/workflows | Wrong order, repeated invocation, conflicting options, stale references | Valid transitions, rejection without mutation |
| State/structure | Missing links/files, duplicate identities, corrupt metadata, conflicting owners | Detection, authority, preservation of valid state |
| Concurrency | Simultaneous writers, stale reads, duplicate submissions | Atomicity, isolation, idempotency |
| Interruption | Kill mid-write, timeout, disconnect, retry after partial completion | Durability, bounded recovery, unique effects |
| Dependencies/server | Unavailable service, slow response, malformed response, version mismatch | Failure containment, honest status, useful fallback |
| Permissions/security | Wrong role, cross-account reference, traversal, untrusted instructions | Authorization, containment, secret protection |
| Resources | Small disk/quota, growing batch, repeated failure, leaked handles | Bounded time/memory/storage and cleanup |
| Accuracy | Rounding, time zones, sorting, missing values, contradictory inputs | Domain-correct output and traceable calculations |
| Usability/feasibility | Ambiguous error, impossible prerequisite, misleading success, hard recovery | User can understand, finish, or recover within stated bounds |

Give each case a campaign-local label, target surface, plausible actor/fault,
invariant source, priority with rationale, and estimated execution cost.
For multiple products, show per-product coverage and applicable integration cases.
Done when every case challenges a distinct failure mechanism and every applicable
quality dimension has coverage or an explicit gap. Reduce a forced quota when the
actual target supports fewer meaningful cases, explaining why.

## 3. Make the suite executable

Use [scenarios.md](scenarios.md) for storage, metadata, provenance, and verdicts.
Keep one-off cases in the plan; extract independently reusable cases into
.spectacular/scenarios/. Store campaign scope in plans/ and attempts in audits/
when persistence is requested or established. Without a workspace or write access,
return the suite and evidence in chat with the unsaved state explicit.

Each case states synthetic fixtures and initial state; exact actions or commands;
the injected fault or misuse; expected output/state and failure oracle; timeout,
resource and stopping bounds; evidence to collect; and reset/cleanup procedure.
Use concrete values and runnable steps after inspection. Mark unresolved steps
provisional. Performance claims need a stated threshold, workload, and environment;
usability claims need an observable task/error criterion rather than taste alone.
Done when another tester can run and judge each executable case without guessing.

## 4. Run bounded trials and preserve failures

Execute cheapest high-impact cases first in a disposable copy, test tenant, or
isolated fixtures. Match authorization to actual effects. Local reversible trials
may proceed within the request; live destructive actions, external messages,
charges, or publication require existing explicit authorization or remain blocked.
Use supported test fault injection; never corrupt the user's real workspace,
bound Contracts, frozen Missions, credentials, or production data to prove a case.
Stop an attempt at its declared limit or unexpected external/data-loss effect;
preserve evidence, contain the effect, and continue independent safe cases.

Record the actual revision, environment, exact inputs/actions, exit status,
expected versus observed output/state, timing where relevant, and cleanup outcome.
Redact secrets from logs; keep evidence locatable. Use the shared verdicts in
scenarios.md: passed, failed, blocked, or not verified. Execution failure alone
is not a product bug: distinguish harness/setup faults from violated guarantees.
Retries are separate attempts; retain the original failure and changing conditions.

For each failure, minimize the reproduction and retry within budget after reset.
If it cannot be reproduced, report the observed failure with uncertainty intact.
A finding includes case label, severity rationale, expected versus actual behavior,
minimal reproduction, evidence, affected revision, and suspected cause separately
from confirmed facts. Group the same root cause without losing case-level outcomes.
Testing alone does not authorize product fixes; when repairs are already requested,
preserve failure evidence first, repair, and rerun the reproduction plus relevant
regressions against the changed revision. Passing retests do not erase earlier failures.
Done when each planned case has an attributable verdict or explicit evidence gap,
with cleanup verified or residual state reported.

## 5. Close with actionable results

Return counts by verdict, highest-impact confirmed findings, reproduction/evidence
links, limitations and untested surfaces, cleanup status, and the next repair or
verification action. Label design concerns and hypotheses separately from bugs.
A bounded suite cannot certify that software is unbreakable. Keep useful cases for
regression; retire obsolete cases using the existing scenario retirement rules.
Done when the reader can reproduce findings and distinguish tested guarantees
from unknowns without relying on a blanket success claim.
