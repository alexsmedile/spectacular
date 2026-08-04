---
doc-id: plan
mode: grill
location: .spectacular/requests/<slug>/PLAN.md
scope: per-request
template: templates/plan/base.md
slots: [Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables]
snapshot-on-edit: false
summary: "Request-scoped plan — 7-slot decomposition (owns lifecycle state)"
status: active
---

# PLAN Rules — routing card

Keep the registry/template contract above. Load exactly one companion:

- Grill/scaffold → [[plan-authoring]]
- Refine → [[plan-refine]]
- Review → [[plan-review]]

The lifecycle and verification safeguards remain mandatory; this card prevents
an author from loading refine/review detail unnecessarily.
