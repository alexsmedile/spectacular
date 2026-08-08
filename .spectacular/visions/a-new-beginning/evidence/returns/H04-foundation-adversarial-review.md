---
type: handoff-return
handoff_id: H04
status: complete
verdict: sound-with-required-changes
disposition: accepted-with-repairs
reviewed_commit: c8ff3fd34a5cf508101e57e195ca58f66ff14f25
source_thread: 019fe0ef-af1f-77c3-87a0-eda8e5b034ef
authority: adversarial-review
---

# H04 — Independent foundation adversarial review

H04 reviewed immutable Git objects at `c8ff3fd`; the live checkout had advanced
to `7a85469` and was dirty, so later content was deliberately excluded.

## Verdict

**Sound with required changes.** The method is directionally strong and
repairable, but its spec/program boundary, decision order, anchoring controls,
and orchestration envelope were inconsistent with its own safeguards.

## Findings and central disposition

| ID | Finding | Central disposition |
|---|---|---|
| B1 | S12 produced draft specifications and executable requests despite the approved-spec-before-request invariant | **Accept.** Split specification approval from executable-program compilation. |
| B2 | S02 evidence rules preceded S03 truth/provenance; S10 deletion preceded a compatibility floor | **Accept.** Run minimal truth semantics before S02 and establish compatibility constraints before S10. |
| B3 | Repeated control-plane, Mission, runtime, verification, and companion recommendations anchored nominally undecided sessions | **Accept with temporal qualification.** The defect was real at `c8ff3fd`; H02 later supplied owner dispositions for the control-plane proposal. H03 must still audit anchoring, and every S04–S09 placement remains conditional. |
| B4 | Return packets and traffic controls lacked immutable input identity and actual Git/worktree/file conflict evidence | **Accept.** Version/hash handoffs, bind returns to inputs/contracts/tree/read set/reviewer, reject or revalidate baseline drift, and require real overlap inspection before mutation. |

## Required program repairs

- Add immutable handoff identity, accepted-contract versions, read set, reviewer,
  and reviewed tree to returns.
- Make baseline mismatch an explicit `reject | revalidate` gate.
- Preserve symmetric credible alternatives until the owner disposes the decision.
- Distinguish owner/proposal/self-hosting evidence from real-user evidence.
- Define minimal truth/provenance before the success/evidence constitution.
- Define a compatibility floor before deletion decisions.
- Keep S05–S06 capability/authority definitions responsibility-neutral until S07.
- Split S12 into approved specification topology and later program compilation.
- Require traceability from promise through contract/spec/Mission/evidence/rollback.
- Test stale-baseline, conflicting-return, and cold-orchestrator-succession cases.

## Non-blocking risks retained

God-context bottleneck; shadow truth in summaries; correlated-model bias; semantic
conflicts despite disjoint files; missing revalidation order; owner fatigue;
preventive-value blindness; metric gaming; secrets in committed evidence; and
software/Git assumptions that may narrow the claimed audience.

## Strongest repaired order

```text
immutable orchestration envelope
→ S01 product constitution
→ S03A minimal truth/provenance
→ S02 success/evidence constitution
→ S03B contract model
→ S04 work ontology
→ S05–S06 authority/evidence capabilities
→ coherence gate
→ S07–S09 responsibility and surfaces
→ compatibility floor
→ S10 subsystem survival
→ S11 architecture/migration
→ S12A specification approval
→ S12B executable Mission program
→ implementation
```

H04 accepts no product decision and authorizes neither H03 nor S02.

## Return packet

```yaml
return:
  handoff_id: H04
  status: complete
  baseline: c8ff3fd34a5cf508101e57e195ca58f66ff14f25
  result: sound-with-required-changes; four foundation blockers require repair
  decisions: []
  facts:
    - S12 contradicted the approved-spec-before-request invariant.
    - S02 preceded truth/provenance and S10 preceded compatibility constraints.
    - Pending product/responsibility recommendations were repeated as baselines.
    - Return and traffic contracts lacked immutable identity and actual overlap checks.
  assumptions:
    - Distributed decision work is required; its final product placement is open.
    - No real-user research corpus was included.
  artifacts: []
  evidence:
    - METHOD, ORCHESTRATION, decision-sessions, top-20, responsibility-boundaries, PRD, and CLI at c8ff3fd
  conflicts:
    - S12 versus METHOD phases 6–7
    - S10-before-S11 versus compatibility-before-reduction
    - pending decisions versus repeated working recommendations
    - exact baseline contract versus advanced dirty checkout
  scope_deviations: []
  next_action: repair B1–B4, then issue downstream handoffs from a validated immutable baseline
```
