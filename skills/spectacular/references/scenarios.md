# User scenarios and field validation

Use this when: creating end-to-end user scenarios, preparing beta tests, or
checking feasibility and actual behavior with real software.
For deliberate adversarial breakage campaigns, start with [breakage.md](breakage.md);
this file owns shared case storage, attempt provenance, and verdicts.

## Store the journey separately from its attempts

A Scenario is a reusable user journey: a starting situation, an objective, and
an observable final output. Keep independently useful cases as `type: Scenario`
in `.spectacular/scenarios/`, using descriptive filenames and the existing
[soft metadata profiles](profiles.md). No schema claim, governed identity, CLI,
or Mission is required. Create the folder when the first case needs storage.

Update an existing case before adding a duplicate. A small one-off case can
stay in its plan; extract it when it needs independent reuse or lifetime.
Requirements own needs and specs own intended behavior; link them instead of
copying their authority into the case. A Scenario tests those expectations.

The plan owns campaign scope, coverage, priority, execution order, and progress.
An Audit in audits/ owns attempt-specific conditions, observations, evidence,
and limitations. Small findings can stay in the plan. Link each attempt to its
Scenario path and document version or inspected Git revision; capture the exact
criteria used if the case may change. Reuse a case without overwriting earlier
results. These attempts are ordinary observations, not governed Run records.

## 1. Select representative goals

Read relevant product context and identify software, users, and the outcomes
being tested. Use the requested count; otherwise aim for 5–10 cases when breadth
is useful. Select by frequency, value, and risk, rather than filling a quota.
Use administrative, professional, business, and developer roles only where they
fit the product. Include cross-software journeys when integrations matter.

Record a compact coverage table in the plan: case link, role, capability,
priority, and reason for inclusion. Mark exclusions and assumptions explicitly.
When software access or product context is missing, draft bounded cases with
those unknowns stated; inspection alone cannot prove actual behavior.
Done when every selected case tests a distinct useful goal and coverage gaps
are visible.

## 2. Make each case executable

Describe these parts in the body, adapting detail to the case:

- User, role, goal, and trigger; explain why the outcome matters.
- Software involved and applicable version constraints; required accounts,
  permissions, integrations, initial state, and representative synthetic data.
- Actions from the starting situation through completion, including handoffs.
- Final output or state and observable criteria defined before execution.
- Relevant failure variants and expected recovery: incomplete data, denied
  permission, duplicates, interrupted work, or unavailable integration.
- Evidence to collect, cleanup, and any action requiring explicit authorization.

Keep the expected outcome independent from the current implementation. Choose
failure variants by risk; do not add every failure type to every case.
Done when another tester can start, execute, and judge the case without guessing
its prerequisites or finish line.

## 3. Execute and preserve observations

Scenario preparation authorizes drafting, not sending messages, spending money,
publishing, or other external effects. Apply the user's actual authorization
before those steps. Use synthetic data and reversible test environments where
available; redact credentials and personal data from retained evidence.

Record actual software versions, environment, date, case basis, executed actions,
expected versus observed output, evidence locations, and cleanup outcome. Give
variants their own outcome when they differ from the main journey. Retain
failures and blocked steps; on retry add a new attempt or clearly separate entry.
If a required tool or access is unavailable, record the blocked step and continue
independent authorized cases. Never substitute imagined execution for evidence.

Distinguish the basis (hypothesis, documented feasibility, or observed execution)
from the execution verdict:

| Verdict | Meaning |
|---|---|
| Passed | Executed; evidence supports every stated criterion for this attempt |
| Failed | Executed; at least one stated criterion is contradicted |
| Blocked | An identified prerequisite or permission prevented a required step |
| Not verified | Not executed, or evidence is insufficient to decide |

A passing main journey with an untested required variant is not a full case pass.
An illustrative case or documentation review remains unverified in the field.
Done when each planned case and required variant has a verdict, attributable
basis, evidence or an explicit evidence gap, and remaining work is stated.

## Example: reconcile supplier invoices

Illustrative case only; software availability and behavior are not verified.
An accountant needs to reconcile a synthetic batch of supplier invoices with
purchase orders, then produce a discrepancy report for approval.

Prerequisites: an accounting test tenant, an order export, import permission,
and three synthetic invoices: one matching, one with a mismatched amount, and
one duplicate. Record the selected software and versions before execution.

Journey: import the batch, associate invoices with orders, inspect exceptions,
resolve the mismatch through the permitted workflow, and export the report.
Success: one payable entry per accepted invoice, no duplicate payable entry,
a traceable mismatch resolution, and a report matching the final ledger totals.
Recovery variant: interrupt import and retry; previously accepted entries must
remain unique. Collect redacted import receipts, ledger counts, and the report.
Sending approval requests or paying invoices requires separate authorization.
Clean up synthetic records according to the test tenant's supported procedure.

The case belongs in scenarios/reconcile-supplier-invoices.md; a campaign plan
links it and prioritizes its variants. An Audit records the actual attempt and
verdict. These paths illustrate placement; create only the artifacts needed.

## Retire by relevance

Follow [navigation.md](navigation.md): keep reusable cases active after a test
campaign, and retire obsolete cases to archive/scenarios/ with provenance and
repaired links. A completed campaign does not make its reusable cases obsolete.
