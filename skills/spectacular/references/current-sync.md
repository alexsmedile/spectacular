# Current Sync — Updating Canonical Capability Specs

Triggered by: archive sequence, or when skill detects a request has meaningfully changed a capability.

---

## Purpose

When a request is completed and archived, the canonical capability specs in `current/` should reflect the new system truth. The skill proposes what to update — the human confirms before any write.

---

## What to check

For each `related:` entry in the request's `PLAN.md`, and for any capability specs referenced in `TASKS.md` or `SESSION.md`:

1. Does the capability spec exist? If not, propose creating it.
2. Is the spec still accurate given what was built? If not, propose updating it.
3. Should the `status` change? (e.g., `draft` → `stable` after implementation)

---

## Proposal format

> "Before archiving, here's what I'd update in `current/`:
>
> - `current/billing/plans.md` — update to reflect new team tier; change status: draft → stable
> - `current/auth/sessions.md` — no changes needed, still accurate
>
> Want me to proceed with these updates?"

---

## Snapshot rule

Always snapshot before editing any `current/` file. See `versioning.md`.

Sequence:
1. Human confirms current/ sync proposal
2. For each file to edit: snapshot → edit → bump version
3. Then proceed with moving request to archive

---

## Creating new capability specs

If a request introduced a net-new capability with no spec in `current/`:

Propose creating `current/<capability-group>/<capability>.md` with frontmatter:

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
