# Side-session handoffs

Copy these prompts into separate tasks only in the order authorized by
[`../ORCHESTRATION.md`](../ORCHESTRATION.md). A prompt is a bounded contract, not permission to read
the whole workbench or expand scope.

## Initial queue

| ID | Prompt | State | Dependency |
|---|---|---|---|
| H01 | [Product-boundary evidence audit](H01-product-boundary-evidence-audit.md) | ready after planning-baseline commit | none |
| H04 | [Independent foundation adversarial review](H04-independent-foundation-adversarial-review.md) | ready after planning-baseline commit | none; may run with H01 |
| H02 | [S01 Product Constitution lead](H02-s01-product-constitution-lead.md) | waiting | H01 return |
| H03 | [Product Constitution skeptic](H03-product-constitution-skeptic.md) | waiting | H02 draft/return |

## Templates

- [Implementation Mission](TEMPLATE-implementation-mission.md) — do not use before S12 produces an
  approved spec and Mission boundary.

H04 should receive no conversational summary and should run on a different model when possible.
Review its blocking findings here before H02 begins. All returns come back to this orchestration
task for checkpoint review. Side sessions do not reconcile the program or authorize successors.
