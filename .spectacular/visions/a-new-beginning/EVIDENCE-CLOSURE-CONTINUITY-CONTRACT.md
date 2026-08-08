---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S06
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
next_session: S07
---

# Evidence, Closure, and Continuity Contract

This accepted S06 contract defines the minimum evidence-to-closure loop. It
does not choose exact schemas, storage, retention/redaction mechanics,
commands, public terms, or responsibility placement.

## Claim-scoped evidence

Every material acceptance claim needs attributable, fresh, claim-appropriate
evidence. No score, green suite, checklist, artifact presence, reviewer
confidence, or completion report proves success alone.

Before execution, a Mission declares its completion boundary and appropriate
proof method. Mechanical/configuration claims need reproducible checks;
behavior changes need demonstrated scenarios or contract/regression evidence;
migration/compatibility claims need affected-population and recovery evidence;
operational claims need attributable environment observation. A material change
to the method is an explicit Mission/evidence-plan delta.

For each material claim, the inspectable envelope includes its Mission and
Objective context; target and baseline pointers; evidence pointer, method,
result, actor/observer, revision/environment/time; authority, direct/derived
status, freshness basis, limitations, and contrary evidence; required-check and
review status; and assessment/reconciliation pointers.

Deterministic validation may reject malformed, missing, stale, misattributed,
failed, or baseline-mismatched envelopes. It cannot alone prove semantic,
behavioral, security, operational, or product sufficiency.

## Independent review

Independent review is required when a claim involves security, privacy, rights,
reserved or destructive effects; material user, compatibility, data, or
operational contract change; deployment or observed-environment behavior; or
materially conflicting, indirect, ambiguous, or executor-dependent evidence.
The owner may require it for any Mission.

An independent reviewer did not execute the work or author the disputed
evidence, can inspect primary evidence, and is not relying solely on the same
unsupported claim or correlated method. A fresh agent alone is insufficient.
A qualified human, distinct method, or independent observation can provide
stronger independence. Review reports assurance and blocking findings only; it
does not resolve a Mission or change current behavior.

Routine bounded reversible work with direct reproducible evidence and no
trigger may proceed to owner assessment without a separate reviewer.

## Bounded repair and escalation

A Mission declares a proportional repair budget. Each repair attempt requires a
new hypothesis, new evidence, or materially narrower corrective action, then
reruns the narrowest relevant check before broader verification. It never
broadens the authority envelope.

Unchanged/repeated findings, exhausted budget, unresolved conflict, or
authority/scope drift escalate with the disputed claim, evidence/freshness
pointers, attempts and outcomes, last known-good checkpoint, remaining
alternatives, authority needed, and one safe next action.

## Closure and reconciliation order

```text
assessment → owner disposition → authorized reconciliation → archival
```

Assessment determines whether evidence is ready for the owner, needs repair,
or must escalate. It cannot resolve a Mission. Only the owner resolves a
Mission and separately authorizes current Capability Contract change. Authorized
reconciliation updates or supersedes the responsible current authoritative
material while preserving prior-truth, evidence, and disposition pointers.
Archival retains the resolved record and receipts; it proves nothing and cannot
precede required disposition and reconciliation.

If a Mission is abandoned or superseded, current behavior remains unchanged
unless separately authorized reconciliation says otherwise.

## Continuity and refusal

A terminal continuity return is pointer-first and identifies the Mission,
Objective, target delta and current-contract pointers, baseline, current Run
boundary/checkpoint or terminal result, evidence/review pointers with
freshness/conflicts/limits, open Gaps/dependencies, authority and repair-budget
state, disposition/reconciliation state when present, and exactly one safe next
action or exact owner gate. It never duplicates canonical Mission truth.

Closure is refused when material evidence is missing, stale, unattributable, or
uninterpretable; authorities conflict; required checks/review/findings remain;
baseline, target, boundary, or authorization drifted; required provider
attestation is unknown; or repair has no new hypothesis/evidence or exhausted
its budget. Safe unaffected inspection, evidence collection, reproduction,
option preparation, revalidation, and bounded repair/owner-decision drafting
may continue. Acceptance, lifecycle promotion, reconciliation, and external
effects may not.
