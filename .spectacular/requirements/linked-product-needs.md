---
type: Requirement
version: "0.1"
created: "2026-10-07T19:27:19Z"
updated: "2026-10-07T19:27:19Z"
description: Retain independently useful PRD needs as linked context without adding ceremony.
---

# Linked product needs

## Need

People and agents need to navigate and evolve individual product needs without
maintaining a monolithic PRD or duplicating requirements across delivery plans.
The owner selected this addition for Spectacular v3.1.

## Scope and constraints

Use an optional requirements/ container with existing compact metadata. Keep
existing Requirement files in specs/ valid. Plans and prototypes link the needs
they address; specs own behavior and interfaces. Contracts retain their enforced
agreement and amendment rules. No new command, required identity, lifecycle,
method file, or bulk migration is introduced.

## Acceptance

- Requirement and Reference files are recognized by optional folder diagnostics.
- Linked requirements, specs, and plans resolve; existing specs/ Requirements pass.
- A misplaced Plan in requirements/ produces an advisory folder-type finding.
- Planning and prototyping can consult optional skeleton and anatomy guidance.
- Product docs show the container and its distinction from specs and plans.

[Delivery plan](../archive/plans/linked-requirements-and-drafting.md) records verification.
[Runtime guide](../../skills/spectacular/references/requirements.md) owns usage.
