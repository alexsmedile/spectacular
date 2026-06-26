---
description: Proposing SPEC.md + specs/ updates when archiving a SPEC-touching request.
when_to_use: Archive flow on a request that changed what is built.
---

# Spec Sync — Updating the System Spec After a Request Ships

Triggered by: archive sequence, when skill detects a request has meaningfully changed a capability, or when `spectacular doctor specs` / `status` flags SPEC.md drift (its `updated` predates the newest archived request — see [[doctor-areas]] `specs` area). In the drift case there's no active request to sync from: read the flagged `archive/<slug>/PLAN.md` and check each `related:` capability against SPEC.md, same as below.

---

## Purpose

When a request is completed and archived, `.spectacular/SPEC.md` (and any per-capability `specs/<capability>/SPEC.md`) should reflect the new system truth. The skill proposes what to update — the human confirms before any write.

---

## What to check

For each `related:` entry in the request's `PLAN.md`, and for any capability specs referenced in `TASKS.md` or `SESSION.md`:

1. Does the top-level `SPEC.md` capabilities list mention this capability? If not, propose adding a bullet.
2. Does the capability have its own per-capability file? Only required when the bullet outgrows one line.
3. Is the per-capability spec still accurate given what was built? If not, propose updating it.
4. Should the `status` change? (e.g., `draft` → `stable` after implementation)

---

## Proposal format

> "Before archiving, here's what I'd update:
>
> - `.spectacular/SPEC.md` — add bullet: 'team-billing — Stripe-backed seats with idempotent webhooks'
> - `specs/billing/SPEC.md` — update to reflect new team tier; change status: draft → stable
> - `specs/auth/SPEC.md` — no changes needed, still accurate
>
> Want me to proceed with these updates?"

---

## Snapshot rule

Always snapshot before editing `.spectacular/SPEC.md` or any `specs/<capability>/SPEC.md`. See `versioning.md`.

Sequence:
1. Human confirms spec sync proposal
2. For each file to edit: snapshot → edit → bump version
3. Then proceed with moving request to archive

---

## Creating new capability specs

If a request introduced a net-new capability with no spec in `specs/`:

Propose adding a bullet to `.spectacular/SPEC.md` capabilities list. Only create a per-capability file when that bullet outgrows one line.

When you do break out per-capability, create `specs/<capability>/SPEC.md` (or `specs/<capability-group>/<capability>/SPEC.md` for nested) with frontmatter:

```md
---
status: stable
updated: <today>
summary: "<one-line description of what this capability does>"
---

# <Capability name>

## Purpose
<What this capability does>

## Requirements
- 

## Scenarios
- 

## Security considerations
- 

## Performance expectations
- 
```
