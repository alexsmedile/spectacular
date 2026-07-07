---
description: Archive a completed request — move to archive/, propose spec sync + memory entries.
when_to_use: spectacular archive <slug>.
---

# Archive — Completing a Request

Triggered by: `spectacular archive <slug>`, or skill proposing archive after request reaches `verified`.

> **@Archive policy gate.** First, run `spectacular policy @Archive` and follow every active policy. Defaults (`spec-sync`, `memory-propose`) are `warn`: propose the SPEC/`specs/` sync and any memory worth keeping, then continue. See [policy-injection.md](policy-injection.md).

---

## Archive sequence

1. **Verify state** — confirm request is `verified` (or ask user to confirm if skipping verification)
2. **Write the spec delta** — `SPEC-DELTA.md` in the request folder (structured ADDED/MODIFIED/REMOVED, or `NONE — <why>`). See `spec-sync.md` § Proposal format. **The gate blocks archive if this file is missing** — write it here, not after.
3. **Propose memory entries** — review the request for lessons worth keeping (see `memory.md`)
4. **Propose bug-lifecycle captures** — if this request fixed a bug: offer a `fixes/` entry (`spectacular fix new … --signature …`) so the resolution is reusable next time (see [[bug-workflow]] Step 3). If it closes an open `audit/`, resolve it (`spectacular audit resolve <A> --disposition "requests/<slug>"`) so no investigation is left dangling. Skip silently when neither applies — most requests aren't bugs.
5. **Human confirms** the spec delta, memory, and bug-lifecycle proposals before any writes
6. **Run `spectacular archive <slug>`** — the CLI verb (v0.7.0+) runs the **closure gate**, then does the move + frontmatter bump + inbound link rewriting atomically. Do NOT manually `git mv` and then sed link paths — that's fragile and easy to leave half-broken.
7. Update any `specs/<capability>/SPEC.md` files (and a bullet in `SPEC.md`) that reference this request

### What `spectacular archive <slug>` does

- Refuses unless `status:` is `verified` or `review` (use `--force` to override — rare; usually a sign of a wrong call)
- **Runs the closure gate** (see below) — three checks that block a drifting archive
- Sets PLAN.md frontmatter: `status: archived`, `archived: <today>`, `updated: <today>`
- Moves `.spectacular/requests/<slug>/` → `.spectacular/archive/<slug>/` (via `git mv` if in a repo)
- Rewrites every inbound `related:` link in other request files from `../<slug>/...` to `../../archive/<slug>/...`

The CLI is the canonical mutator. Manual editing is for unusual cases the verb doesn't cover.

---

## The closure gate

Archive is where docs and reality reconcile — so `cmd_archive` refuses to close a request that left its own tracking half-finished. Three mechanical checks, each **block-with-recorded-override** (never a silent bypass):

| Check | Blocks when | Fix, or override |
|---|---|---|
| `tasks` | a TASKS.md box is open (`- [ ]`) or deferred without a reason (`- [~]` with no ` — <reason>`) | tick it, or write `- [~] <task> — <reason>`; else `--override tasks --reason "<why>"` |
| `verify` | VERIFY.md exists but VERIFY-LOG.md has no passed (`✅`) walk row — i.e. it was never actually walked | run `spectacular verify <slug>`; else `--override verify --reason "<why>"` |
| `spec` | no `SPEC-DELTA.md` in the request folder | write the delta (step 2 above), or `NONE — <why>`; else `--override spec --reason "<why>"` |

**Override etiquette.** An override is a recorded decision, not an escape hatch. Each `--override <check> --reason "<text>"` stamps `archive_overrides:` into the archived PLAN frontmatter (a list of `{check, reason, date}`) — so the shortcut is auditable forever. Prefer fixing the underlying gap; reach for `--override` only when the gap is *legitimate* (deferred task with a real reason, a request that genuinely walked its checks by hand, a refactor with no spec impact you'd rather record as `NONE`). If you find yourself overriding the same check every archive, the gate is telling you something — raise it, don't route around it.

`--force` is unrelated: it bypasses the **status gate** only (archiving a `planned`/`active` request), never a closure check. The two flags don't overlap.

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
- Use `spectacular archive <slug>` — never manual `git mv` + sed (v0.7.0+)
- **Only requests archive.** The append-only soft-DB collections (`memory/`, `decisions/`, `sessions/`, `audit/`, `fixes/`) are *not* archived — they stay live. A resolved `audit/` entry keeps `status: resolved` in place; `fixes/` is a permanent corpus by design (that's the point of the self-learning loop). `feedback/` is the one exception — it has its own `spectacular feedback-loop archive` verb → `archive/feedback/<year>/`. See [[soft-db-index]].
- **Reversible (v1.22.0+):** an accidental archive is undone with `spectacular undo` — it moves the dir back, restores status, and reverses the inbound-link rewrites. Surface this as a one-line tier-reveal after archiving, not a separate flow. Undo is single-level and refuses on a stale breadcrumb.

---

## Example archive proposal

> "Ready to archive `add-team-billing`. Before I move it:
>
> **Spec delta** (written to `SPEC-DELTA.md`):
> ```
> ### ADDED
> - specs/billing/SPEC.md :: team-billing — Stripe-backed seats; a duplicate webhook never double-charges
> ### MODIFIED
> - SPEC.md :: "billing — single-seat only" -> "billing — multi-seat, Stripe-backed"
> ```
>
> **Memory entries proposed:**
> - The Stripe webhook idempotency issue we hit — worth noting in `memory/failures.md`
>
> All TASKS closed, verify walk recorded — the gate is green. Confirm to proceed, or adjust either proposal."
