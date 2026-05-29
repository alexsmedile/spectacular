---
description: The 2-of-6 rule — when a request needs a VERIFY.md vs folding verification into PLAN.
when_to_use: Deciding where verification lives for a request.
---

# Verification — when VERIFY.md is needed vs folded into PLAN/TASKS

Loaded when the skill is deciding whether to scaffold a `VERIFY.md` for a request, or when moving a request from `review` → `verified`.

## Core principle

**Verification always happens. The standalone file is opt-in.**

To be precise about what "opt-in" means here — it refers to **whether you scaffold a `VERIFY.md` file**, not whether you verify the work. Every request runs verification before reaching `verified` status; the only question is *which artifact carries the checks*.

**Never skip verification because VERIFY.md is "optional."** If VERIFY.md exists, it is load-bearing — every unchecked item blocks the `verified` transition. If it doesn't exist, the checks live in PLAN § Validation or TASKS § Verification, and *those* block the transition instead. There is no path to `verified` without explicit verification of some artifact.

**Three kinds of checks exist; conflating them creates noise.**

| Check type | Question it answers | When it runs | Natural home |
|---|---|---|---|
| Task completion | Did I do the work? | During `active` | `TASKS.md` checkboxes |
| Acceptance criteria | Did we build the right thing? | At `review` entry | `PLAN.md` § Validation |
| QA / risk verification | Did we build it correctly + safely? Will it break? | During `review` | `VERIFY.md` (when needed) |

The first two are universal. The third is **conditional** — it earns a standalone file only when the request's risk and surface justify it.

## The 2-of-6 rule

Scaffold a standalone `VERIFY.md` when **at least two** of these are true:

1. **User-visible change** — anyone outside the team can observe behavior change. Demo it, screenshot it, click through it.
2. **High reversibility cost** — migrations, schema changes, destructive operations, anything where "oops" is expensive.
3. **Multi-surface verification** — needs more than running tests: manual QA, browser checks, prod smoke, etc.
4. **Risk surface non-trivial** — auth, billing, payments, PII, security, data integrity.
5. **External contract change** — public API, exported library shape, plugin interface, CLI flags users rely on.
6. **Rollback plan exists** — there's a non-trivial undo procedure that itself needs verification.

The common thread: checks that can't be expressed as `- [x]` next to a task because they're multi-step, time-sequenced, or require human judgment in a specific environment.

**Default to no file.** A new file must earn its keep. Spectacular's principle: small files over giant docs, but also fewer files over more files when content fits cleanly elsewhere.

**Reminder:** "no file" ≠ "no verification". The checks still exist — they just live in PLAN or TASKS instead of a dedicated file. See § Folded patterns below.

## When to fold into PLAN/TASKS instead

VERIFY.md is overkill — and the noise hurts — for these request types:

- **Doc-only changes** (README, ARCHITECTURE, PRD edits) — review = "does it read well?"
- **Internal refactor with test coverage** — review = "tests still pass + no behavior change"
- **Build-system / tooling changes** — review = "builds clean, CI green"
- **Spec/template additions** — review = "dogfood test passes"
- **Configuration changes** — review = "loads correctly, no breakage in staging"
- **Anything single-step** — if verification fits in one checkbox, it's a TASKS item

For these, use one of the two folded patterns below.

## Folded pattern A — PLAN § Validation (preferred for small requests)

Per-milestone validation criteria live in PLAN.md slot 6. Each milestone gets one or more lines describing how it's confirmed.

```md
## 6. Validation

- M1 — `references/kits-contract.md` exists with full schema documented
- M2 — All 5 kits versioned to 2.0; v1.1 snapshots in versions/
- M3 — Dogfood test: coding kit produces 10-slot merged sequence with `kit: coding` frontmatter
```

Best when validation is descriptive ("this property holds") rather than procedural ("run these steps").

## Folded pattern B — TASKS § Verification (preferred when checks are step-by-step)

Add a `### Verification` group at the end of TASKS.md. Use the same checklist syntax as the rest of the file.

```md
### Verification
- [ ] Scaffold a coding-kit PRD in /tmp/test-workspace/
- [ ] Confirm slot order is base 8 + Stack@8 + Interfaces@9
- [ ] Confirm `kit: coding` written to frontmatter
- [ ] Run kit-aware gate with 3 scenarios (required-missing, optional-empty, no-kit)
```

Best when verification has procedural steps that need to be ticked off — closest in shape to a "lite VERIFY.md" without the separate file.

## Standalone VERIFY.md shape

When 2-of-6 triggers it, scaffold from `templates/verify/base.md` (planned — see `requests/verify-doc/` when opened). Until that template exists, the inline stub in `scaffold-reference.md § VERIFY.md` is the source of truth.

Sections:
- **Manual QA checklist** — human-driven steps in a specific environment
- **Edge cases to verify** — non-obvious scenarios likely to break
- **Regression checklist** — confirm previously-working flows still work
- **Rollback validation** — if applicable, prove the undo procedure works

Location: `.spectacular/requests/<slug>/VERIFY.md` (per-request, colocated with PLAN + TASKS).

## How the lifecycle interacts

The `review → verified` transition checks one of two artifacts, depending on what the request scaffolded:

| If VERIFY.md exists | If it doesn't |
|---|---|
| All `- [x]` in VERIFY.md → propose `verified` | All `- [x]` in TASKS § Verification (if present) AND/OR all PLAN § Validation items confirmed by human → propose `verified` |

Same gate, different artifact. The skill detects which is present and adapts.

## Decision flow (when to scaffold VERIFY.md)

When a request enters `review`, the skill asks itself the 2-of-6 questions. If yes:

> "This request looks like it needs VERIFY.md (reason: <which axes hit>). Create one now?"

If no:

> "Validation for this request is in PLAN § Validation / TASKS § Verification. Ready to walk through and mark `verified`?"

User can always override either way.

## Examples from this project's history

- `prd-craft v1.1` — doc-only, dogfooded via /tmp simulation. **No VERIFY.md.** PLAN Validation was enough.
- `doc-writer` — internal spec + skill. **No VERIFY.md.** Dogfood test in TASKS.
- `kits-as-plugins` — spec + templates. **No VERIFY.md.** Validation in PLAN.
- `smart-init` (planned) — user-visible CLI behavior + overwrite risk + external contract (flags). **Will need VERIFY.md.** Manual QA on fresh dir + existing workspace + idempotency.
- `doctor` (planned) — user-visible behavior + opt-in repairs (reversibility matters). **Will need VERIFY.md.**

## Anti-patterns

- **Skipping verification because VERIFY.md is "opt-in"** — opt-in refers to the *file*, not the *practice*. Verification always runs against some artifact.
- **VERIFY.md exists but is ignored** — if the file is present, every unchecked item blocks `verified`. No "we'll come back to it" exits.
- **VERIFY.md for every request** — creates empty noise; ignored over time
- **VERIFY.md duplicating TASKS** — if it's just a copy of the implementation checklist, fold it back
- **PLAN § Validation as a wishlist** — must be checkable, not aspirational
- **No verification record at all** — `verified` status without any artifact (PLAN Validation, TASKS Verification, or VERIFY) is unsupported by the lifecycle rule

## Related

- [[lifecycle]] — `review → verified` transition rule
- [[archive]] — verified is a precondition for archive
- [[scaffold-reference]] — VERIFY.md inline stub (until `templates/verify/base.md` lands)
- [[principles]] — "small files over giant documents" + "humans decide, agents propose"
