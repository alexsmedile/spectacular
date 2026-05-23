# Archive — Completing a Request

Triggered by: `spectacular archive <slug>`, or skill proposing archive after request reaches `verified`.

---

## Archive sequence

1. **Verify state** — confirm request is `verified` (or ask user to confirm if skipping verification)
2. **Propose spec sync** — see `spec-sync.md` for what to update in `SPEC.md` + `specs/`
3. **Propose memory entries** — review the request for lessons worth keeping (see `memory.md`)
4. **Human confirms** both spec and memory proposals before any writes
5. **Move the request directory** — `requests/<slug>/` → `archive/<slug>/`
6. Update any `specs/<capability>/SPEC.md` files (and a bullet in `SPEC.md`) that reference this request

---

## What to look for when reviewing for memory

Read through the archived request for:
- Blockers that weren't obvious upfront
- Architecture decisions made mid-request (not captured in DECISIONS.md)
- Bugs or failures discovered during implementation
- Implementation patterns worth reusing
- Risks from RISKS.md that actually materialized

Propose concrete, specific entries. Avoid vague lessons. If nothing notable happened, say so — don't manufacture memory entries.

---

## Archive rules

- Keep original slug — `archive/add-team-billing/`
- Never modify archived content after moving
- Skill does not read `archive/` during normal operation
- Promoted idea files go to `archive/ideas/<filename>.md`

---

## Example archive proposal

> "Ready to archive `add-team-billing`. Before I move it:
>
> **Spec sync proposed:**
> - `.spectacular/SPEC.md` — add bullet for team-billing capability
> - `specs/billing/SPEC.md` — update to reflect new team tier (status: draft → stable)
>
> **Memory entries proposed:**
> - The Stripe webhook idempotency issue we hit — worth noting in `memory/failures.md`
>
> Confirm to proceed, or adjust either proposal."
