---
description: The runtime loop — how the skill retrieves a hook's policies on entering a phase, injects them, evaluates each check, and resolves by severity.
when_to_use: Entering any work phase that carries a policy gate (init, planning, implementation, debugging, verification, archive, remember, snapshot, session-end).
---

# Policy Injection — the runtime loop

POLICY.md defines the practice layer; this doc is how it *fires*. The skill never loads all policies — on entering a phase it retrieves only that phase's hook and injects those rules. Progressive disclosure (Principle 6) applied to the rule layer itself.

## The mechanic: the ref doc *is* the phase boundary

There is no event bus and no `hooks.json` wiring. The skill already loads exactly one reference doc per phase (the SKILL.md routing model), so each phase's ref doc opens with a **gate block**:

> **@\<Hook\> policy gate.** First, run `spectacular policy @\<Hook\>` and follow every active policy returned…

When the skill loads that doc to do the phase, the first thing it reads is the instruction to consult policies. **The gate fires precisely when the phase begins, by construction.**

| Hook | Gate lives in | Fires at |
|---|---|---|
| `@Init` | `init-workflow.md` | `spectacular init` |
| `@Planning` | `new-request.md` | `spectacular new` / request authoring |
| `@Implementation` | `active-request.md` + `lifecycle.md` | `planned → active` |
| `@Debugging` | `bug-workflow.md` | a bug/quirk/regression is reported |
| `@Verification` | `verify.md` + `lifecycle.md` | `review → verified` |
| `@Archive` | `archive.md` | `spectacular archive` |
| `@Remember` | `memory.md` | memory written |
| `@Snapshot` | `versioning.md` | canonical doc overwritten |
| `@SessionEnd` | `sessions-rules.md` | session end / handoff |

## The loop

```
1. Phase entered     skill loads the phase ref doc; reads the gate block
2. Retrieve          spectacular policy @<Hook> --json
3. Inject            returned policies + linked principle lines enter context
                     (nothing from other hooks loads — Principle 6)
4. Evaluate each     mechanical check → doctor / CLI presence-check
                     judgment check  → skill reads the check prose and decides
5. Resolve by severity
                     block  fails → REFUSE the action; cite the policy + principle
                     warn   fails → SURFACE the finding; continue
6. Done              context holds only this hook's rules
```

## Retrieving

`spectacular policy @<Hook>` is the workhorse. It returns the hook's **enabled** policies (config overrides already merged) plus, for each, the one linked principle line — not the whole PRINCIPLES.md.

```
$ spectacular policy @Implementation
@Implementation
  ⛔ understand-before-change     block
      → P7. Three layers: intent → execution → validation — Every unit of work passes through all three.
```

For machine consumption use `--json`:

```json
[{"hook":"@Implementation","id":"understand-before-change","principle":"7",
  "severity":"block","blocking":true,"enabled":true,
  "check":"PLAN.md has a filled `## Understanding` section …"}]
```

A policy with `"enabled": false` (disabled via config) is still listed but marked — **do not enforce it**.

## Resolving by severity

**Severity is opt-in to blocking.** A policy halts the action **only** if `"blocking": true` (i.e. it literally declared `severity: block`). Everything else — `warn`, absent, unrecognized — is non-blocking.

| Severity | Check passes | Check fails |
|---|---|---|
| `block` | proceed silently | **refuse the action**; state which policy blocked + cite its principle; offer the fix |
| `warn` (or absent) | proceed silently | surface a one-line finding; proceed |

A `block` refusal is never a dead end — name the unmet `check` and the path to satisfy it. Example for `understand-before-change`:

> ⛔ Can't move `planned → active` — policy `understand-before-change` (Principle 7) requires a filled `## Understanding` in PLAN.md (How it works now / What changes / What stays the same), or a `UNDERSTANDING.md` with those three subheads. Want me to fill it now by walking the three questions?

## Mechanical vs judgment checks

Whether a check is *mechanical* or *judgment* is not a schema field — it's how the check is evaluated:

- **Mechanical** (e.g. "`## Understanding` slot exists and is filled", "a `<DOC>@v<N>.md` snapshot exists") → deterministic; `doctor policies` reports it and `--fix` may help.
- **Judgment** (e.g. "the goal is well-formed", "the memory contains no secrets") → the skill reads the check prose and decides.

The same policy can be evaluated both ways over time. The injected instruction is authoritative for judgment checks; the doctor is authoritative for mechanical ones.

## Why no harness hooks (v1)

Enforcement is **skill-side + doctor**, not `hooks.json`. This keeps policies working identically in a bare-CLI session and an installed-plugin session, and avoids a runtime the bare CLI doesn't have. Kernel-level locks (harness-enforced, un-skippable) are the **v2** upgrade path. The gate-block-in-ref-doc mechanic is the v1 contract.

## Related

- [policies-contract.md](policies-contract.md) — POLICY.md structure + policy anatomy + the 9 hooks
- [policy-rules.md](policy-rules.md) — doc-verb dispatch (grill/refine/review) for POLICY.md
- [lifecycle.md](lifecycle.md) — the two spine transitions that consult policies
- [verify.md](verify.md) — the `verification-present` policy's home phase (Part 2, the 2-of-6 rule)
