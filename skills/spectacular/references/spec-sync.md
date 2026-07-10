---
description: Proposing SPEC.md + specs/ updates when archiving a SPEC-touching request.
when_to_use: Archive flow on a request that changed what is built.
---

# Spec Sync — Updating the System Spec After a Request Ships

Triggered by: archive sequence, when skill detects a request has meaningfully changed a capability, or when `spectacular doctor specs` / `status` flags SPEC.md drift (its `updated` predates the newest archived request — see [[doctor-areas]] `specs` area). In the drift case there's no active request to sync from: read the flagged `archive/<slug>/PLAN.md` and check each `related:` capability against SPEC.md, same as below.

---

## Purpose

When a request is completed and archived, `.spectacular/specs/index.md` (and any per-capability `specs/<capability>.md`) should reflect the new system truth. The skill proposes what to update — the human confirms before any write.

---

## What to check

For each `related:` entry in the request's `PLAN.md`, and for any capability specs referenced in `TASKS.md` or `SESSION.md`:

1. Does the top-level `SPEC.md` capabilities list mention this capability? If not, propose adding a bullet.
2. Does the capability have its own per-capability file? Only required when the bullet outgrows one line.
3. Is the per-capability spec still accurate given what was built? If not, propose updating it.
4. Should the `status` change? (e.g., `draft` → `stable` after implementation)

> **Optional — dispatch [[spec-reviewer]] for an arms-length currency pass (judgment-gated).** Before you write the delta, consider dispatching `spec-reviewer` over `specs/index.md` (and any touched `specs/<cap>.md`). It greps each capability claim against the code and cross-checks the roadmap/CHANGELOG ledger, returning a punch list of **stale / gap / premature** claims with cited evidence. This matters most when the request shipped a substantial capability or the spec hasn't been reviewed in a while — you write the SPEC-DELTA against a *known-current* baseline instead of one that may have drifted from earlier requests. It's the specs' guardian doing the currency check you'd otherwise do by hand; it returns findings, you (the orchestrator) fold them into the delta. Same worth-it economics as the other optional fleet gates — skip it for a one-line spec touch. It never edits the spec; that's your `refine` / delta write.

---

## Proposal format — the spec delta

The proposal is a **structured delta**, written to `.spectacular/requests/<slug>/SPEC-DELTA.md`, not free-form prose. This is what makes the impact mechanically mergeable and lets `spectacular doctor specs` validate it structurally (was: a date-only drift heuristic). The archive gate (`cmd_archive`) blocks unless this file exists — so writing it is part of the archive flow, not optional polish.

Three sections, each a bullet list of `<file> :: <payload>` lines:

```md
### ADDED
- specs/billing.md :: team-billing — seats are Stripe-backed; a duplicate webhook never double-charges

### MODIFIED
- SPEC.md :: "billing — single-seat only" -> "billing — multi-seat, Stripe-backed"

### REMOVED
- SPEC.md :: "auth — password-only login"
```

- **ADDED** — payload is the new bullet text. Prefer observable, SHALL-strength phrasing (see § Creating new capability specs below).
- **MODIFIED** — payload is `"<exact current bullet>" -> "<replacement>"`. The current side is quoted **verbatim** (character-for-character) so the mechanical merge can find and replace it.
- **REMOVED** — payload is `"<exact current bullet>"`, quoted verbatim.

**Empty impact is valid.** If the request changed nothing a spec should record, the entire file is a single line — no section headings:

```md
NONE — internal refactor of the retry queue; no capability changed
```

Present the delta to the human for confirmation before archiving. The merge itself is mechanical (ADDED appends under the file's capabilities list, MODIFIED does a quoted-string replace, REMOVED deletes the quoted line) — you propose, the human confirms, then `cmd_archive` applies it.

> "Before archiving, here's the spec delta I've written to `SPEC-DELTA.md`:
> [show the file]
> Confirm to proceed, or tell me what to change."

---

## Snapshot rule

Always snapshot before editing `.spectacular/specs/index.md` or any `specs/<capability>.md`. See `versioning.md`.

Sequence:
1. Write the delta to `SPEC-DELTA.md` and human confirms it
2. For each file the delta touches: snapshot → apply the delta's ADDED/MODIFIED/REMOVED lines → bump version
3. Then proceed with `spectacular archive <slug>` (the gate confirms `SPEC-DELTA.md` is present)

---

## Creating new capability specs

If a request introduced a net-new capability with no spec in `specs/`:

Propose adding a bullet to `.spectacular/specs/index.md` capabilities list. Only create a per-capability file when that bullet outgrows one line.

**Bullet phrasing:** prefer an observable, SHALL-strength statement of behavior over a feature label — "team-billing — seats are Stripe-backed; a duplicate webhook never double-charges" beats "team-billing support". The acid test: *could someone who has never seen the code tell whether the bullet holds?* An important behavior earns a GIVEN/WHEN/THEN scenario — in the per-capability spec's `## Scenarios` section, never in SPEC.md itself (the index stays one line per capability).

When you do break out per-capability, create `specs/<capability>.md` (or `specs/<capability-group>/<capability>.md` for nested) with frontmatter:

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
