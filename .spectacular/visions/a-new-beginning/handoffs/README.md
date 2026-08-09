# Side-session handoffs

Copy these prompts into separate tasks only in the order authorized by
[`../ORCHESTRATION.md`](../ORCHESTRATION.md). A prompt is a bounded contract, not permission to read
the whole workbench or expand scope.

## Initial queue

| ID | Prompt | State | Dependency |
|---|---|---|---|
| H01 | [Product-boundary evidence audit](H01-product-boundary-evidence-audit.md) | accepted as [primary evidence](../evidence/returns/H01-product-boundary-audit.md) | none |
| H04 | [Independent foundation adversarial review](H04-independent-foundation-adversarial-review.md) | accepted with repairs; [return](../evidence/returns/H04-foundation-adversarial-review.md) | none; fresh-context review |
| H02 | [S01 Product Constitution lead](H02-s01-product-constitution-lead.md) | revision 2 accepted; [return](../evidence/returns/H02-product-constitution-r2.md) | H03 required changes resolved |
| H03 | [Product Constitution skeptic](H03-product-constitution-skeptic.md) | complete; [return](../evidence/returns/H03-product-constitution-skeptic.md) | H01 + H02 + reconciled H04 |
| H06 | S03A truth/provenance owner-decision task | accepted; [return](../evidence/returns/H06-s03a-truth-provenance.md) | accepted S01 |
| H07 | S02 success/evidence owner-decision task | accepted; [return](../evidence/returns/H07-s02-success-evidence.md) | accepted S01 + S03A |

S01, S03A, and S02 are centrally accepted in [`../PRODUCT-CONSTITUTION.md`](../PRODUCT-CONSTITUTION.md),
[`../TRUTH-PROVENANCE-FLOOR.md`](../TRUTH-PROVENANCE-FLOOR.md), and
[`../SUCCESS-EVIDENCE-CONSTITUTION.md`](../SUCCESS-EVIDENCE-CONSTITUTION.md). S03B is the next
authorized owner-decision session.
| H05 | [Competing-skills research study](H05-competing-skills-research-study.md) | accepted; [return review](../evidence/returns/H05-competing-skills-study.md) | ingested as source-015 |
| H16 | [SDLC coherence adversarial review](H16-sdlc-coherence-adversarial-review.md) | accepted; [return](../evidence/returns/H16-sdlc-coherence-adversarial-review.md) | accepted contracts + source-016 |
| H17 | S10 subsystem survival | accepted; [final return](../evidence/returns/H17-s10-subsystem-survival.md), [interim A/B checkpoint](../evidence/returns/H17-s10-interim-clusters-a-b.md), [Cluster C delta](../evidence/returns/H17-s10-interim-cluster-c.md) | thirteen accepted contracts; S11 next-ready |
| H18 | [S11 implementation architecture](H18-s11-implementation-architecture.md) | active in Codex thread `019fe697-77b5-7f63-ab75-702044c02ae1` | thirteen accepted contracts; S12A blocked |

## Templates

- [Implementation Mission](TEMPLATE-implementation-mission.md) — do not use before S12 produces an
  approved spec and Mission boundary.

H04 should receive no conversational summary and should run on a different model when possible.
Review its blocking findings here before H02 begins. All returns come back to this orchestration
task for checkpoint review. Side sessions do not reconcile the program or authorize successors.

H05 is deliberately research-only. It does not import ideas or modify the workbench. The
orchestration task reviews its evidence and explicitly selects any atomic findings worth ingesting
as `source-015` or later.
