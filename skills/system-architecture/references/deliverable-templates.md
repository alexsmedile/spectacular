# Deliverable Templates

Use this when: producing a complete architecture package, architecture decision record, architecture review, focused decision, or evolutionary upgrade report.

Select only the template matching the routed deliverable. Replace every placeholder and remove sections that do not help the audience decide or act.

## Contents

- [Complete architecture package](#complete-architecture-package)
- [Architecture decision record](#architecture-decision-record)
- [Architecture review](#architecture-review)
- [Focused decision response](#focused-decision-response)
- [Evolutionary upgrade](#evolutionary-upgrade)

## Complete architecture package

```markdown
# <System> Architecture

## Executive summary
<Purpose, proposed shape, strongest drivers, and most important trade-off.>

## Plain-language explanation
<What the system does, for whom, what it protects, and what complexity is deliberately hidden.>

## Level 0 product map
<Include for mixed or nontechnical audiences; provide editable source, rendered SVG when requested or supported, and a textual equivalent.>

## Scope
- In scope:
- Out of scope:
- Current state:
- Target horizon:

## Drivers and constraints
| Priority | Driver or constraint | Measure/evidence | Architectural implication |
|---|---|---|---|

## Assumptions and open questions
| Type | Statement | Confidence/owner | Validation or due date |
|---|---|---|---|

## Quality-attribute scenarios
| Quality | Scenario | Target | Validation |
|---|---|---|---|

## Architecture overview
<Responsibilities, boundaries, and why this shape fits the drivers.>

## Technical views
<Include the smallest set that answers the architecture questions: context and/or container, one critical behavior view, and optional trust, deployment, state, decision, roadmap, or delta views.>

## Visual artifact manifest
| View | Audience/question | Editable source | Rendered artifact | Text equivalent | Validation |
|---|---|---|---|---|---|

## Responsibilities and contracts
| Boundary | Owns | Exposes/consumes | Guarantees | Owner |
|---|---|---|---|---|

## Data strategy
| Data/owner | Access pattern and invariants | System of record | Consistency | Derived stores | Retention/recovery | Rationale |
|---|---|---|---|---|---|---|

## Security and privacy
<Assets, trust boundaries, identity, authorization, secrets, encryption, audit, privacy lifecycle.>

## Reliability and operations
<SLOs, failure modes, overload, observability, deployment, rollback, backup, DR.>

## Decisions
| Decision | Drivers | Alternatives | Consequences | Validation |
|---|---|---|---|---|

## Risks and validation
| Risk/unknown | Impact | Likelihood/confidence | Mitigation or experiment | Owner |
|---|---|---|---|---|

## Evolution roadmap
| Phase | Outcome | Dependencies | Compatibility/migration | Rollback | Exit criteria |
|---|---|---|---|---|---|
```

## Architecture decision record

```markdown
# ADR-<number>: <Decision title>

- Status: Proposed | Accepted | Superseded | Deprecated
- Date: <YYYY-MM-DD>
- Owners: <names or roles>

## Context
<Problem, scope, constraints, and evidence.>

## Decision drivers
1. <Ranked driver>

## Considered options
### <Option>
- Benefits:
- Costs and risks:
- Fit to drivers:

## Decision
<Choice and precise scope.>

## Consequences
- Positive:
- Negative:
- Follow-on work:

## Validation and reversal
- Validate by:
- Revisit when:
- Reverse or migrate by:
```

## Architecture review

```markdown
# Architecture Review: <System or proposal>

## Verdict
<Fit for purpose, conditionally acceptable, or not yet supportable, with the main reason.>

## Evidence reviewed
- <Artifact, version/date, and relevant scope>

## Driver coverage
| Driver | Evidence in design | Assessment | Gap |
|---|---|---|---|

## Findings
### <Severity>: <Finding>
- Affected scenario:
- Evidence:
- Consequence:
- Recommendation:
- Validation:

## Strengths
- <Evidence-backed strength>

## Open questions
- <Question that could change the verdict>

## Recommended next actions
1. <Ordered action with owner or exit criterion>
```

## Focused decision response

```markdown
## Recommendation
<Decision in one or two sentences.>

## Why
<Drivers and evidence.>

## Alternatives considered
- <Alternative and why it is weaker here.>

## Consequences
- <Operational, delivery, data, or cost implication.>

## Assumptions and validation
- <Assumption and next test.>
```

## Evolutionary upgrade

```markdown
# Architecture Evolution: <Upgrade>

## Outcome and constraints
<User or business outcome, scope, compatibility window, and non-goals.>

## Baseline seam
<Existing owner, contract, event, store, or module that absorbs the change.>

## Delta
| Element | NEW / CHANGED / UNCHANGED | Behavioral or ownership change | Compatibility |
|---|---|---|---|

## Behavior
<Add a focused sequence or state view only when runtime behavior or lifecycle changes materially.>

## Migration and rollback
| Phase | Change | Coexistence | Rollback | Exit criterion |
|---|---|---|---|---|

## Risks and validation
| Risk or unknown | Validation | Owner |
|---|---|---|
```
