---
status: planned
priority: high
owner: alex
updated: 2026-05-30
summary: "POLICY.md — a practice-layer doc paired with PRINCIPLES.md. Policies are filed under named work-phase hooks (@Init, @Planning, @Implementation, @Verification, @Archive, @Remember, @Snapshot, @SessionEnd); the skill retrieves only the active hook's policies and injects them. Ships always-set with 8 prefilled defaults (4 block, 4 warn)."
related:
  - PRD.md
  - ../../PRINCIPLES.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
target_version: v1.12.0
---

# Plan — policy-engine

> **URGENT — immediately after verify-walk (v1.11).** Captured 2026-05-30 from a live need: an agent should not start implementation without first establishing *how the current system works, what changes, and what stays the same*. That specific gate is one built-in policy; the broader deliverable is the **practice layer** that carries it.

## 1. Goal

Give Spectacular a **practice layer** — `POLICY.md`, the operational sibling to `PRINCIPLES.md`:

- **PRINCIPLES.md = theory** (the *why* — beliefs). Untouched.
- **POLICY.md = practice** (the *how we actually work* — executable rules).

A **policy** is filed under a named **work-phase hook** (e.g. `@Implementation`). When the skill enters a phase, it retrieves *only that hook's* policies and injects them into context — progressive disclosure (Principle 6) applied to the rule layer itself. Policies ship prefilled with sensible defaults; users customize, disable, or add their own via `config.yaml`. A policy may optionally tag the **principle** it enforces, so retrieval pulls the one relevant theory line alongside it.

The headline built-in: `understand-before-change` @ `@Implementation` — a request can't move `planned → active` until the agent has written down how the system works now, what changes, and what stays the same.

## 2. Constraints

- **POLICY.md is one file, always-set.** A single canonical doc (no soft-DB folder in v1), created on *every* `spectacular init` with the 8 prefilled defaults enabled. The enforcement spine exists from day one — deliberately asymmetric with PRINCIPLES.md (optional): theory is optional reading, practice is the operational floor.
- **Hooks name moments, not transitions.** The hook reads as "at \<moment\>" (`@` = "at"); the *policy name* carries the before/after verb (`understand-before-change`). A human reading POLICY.md parses "this policy applies @implementation" as plain English.
- **Skill-enforced, retrieve-by-hook.** Policy *evaluation* that needs judgment is skill-side; mechanical presence-checks go to CLI/`doctor`. The skill never loads all policies — it calls `spectacular policy @<hook>` on entering a phase. No Claude Code harness hooks (`hooks.json` stays empty); this is skill-native and works in bare-CLI and installed-plugin sessions alike.
- **Severity is per-policy.** `block` (refuse to proceed) vs `warn` (surface + continue). Safety-critical gates block; hygiene nudges warn. Config can change any policy's severity.
- **Reuses existing machinery.** Lifecycle transitions (`promote`/`archive`), the doctor severity model, frontmatter signals, and the rules-file dispatch pattern already exist — policies compose with them.

## 3. The hook set (locked)

`@` reads "at"; every hook completes "this policy applies **@\_\_\_**" as natural English.

### Lifecycle spine

| Hook | Reads | Fires at | Maps to |
|---|---|---|---|
| `@Init` | at init | workspace scaffolded | `spectacular init` |
| `@Planning` | at planning | PLAN + TASKS being shaped | `spectacular new` / request authoring |
| `@Implementation` | at implementation | planned → active | `promote → active` |
| `@Verification` | at verification | review → verified | `promote → verified` |
| `@Archive` | at archive | verified → archived | `archive` |

### Key moments

| Hook | Reads | Fires at | Maps to |
|---|---|---|---|
| `@Remember` | at remember | memory written | `spectacular remember` |
| `@Snapshot` | at snapshot | canonical doc overwritten | `snapshot` / overwrite |
| `@SessionEnd` | at session end | skill hands off | `session end` |

**Dropped/deferred:** `@Request`/`@RequestTask` → folded into `@Planning` (one authoring span; slug rule already lives in the `new` verb). `@Review` (active→review) → no distinct policy yet. `@SessionStart` → deferred to v2 (wants a real harness runtime). `@Decide` → v2 candidate. `@Doctor`/`@BeforeCommit` → rejected (circular / wrong runtime).

## 4. Prefilled default policies (locked) — 4 block · 4 warn

| Hook | Policy | →Principle | Severity |
|---|---|---|---|
| `@Init` | `scaffold-contract` | 4 | warn |
| `@Planning` | `request-shape` | 3 / 7 | warn |
| `@Implementation` | **`understand-before-change`** | 7 | **block** |
| `@Verification` | **`verification-present`** (absorbs verify-walk's gate) | 7 / 9 | **block** |
| `@Archive` | `spec-sync` + `memory-propose` | 2 / 5 | warn |
| `@Remember` | **`confirm-before-write`** | 8 | **block** |
| `@Snapshot` | **`snapshot-before-overwrite`** | 8 | **block** |
| `@SessionEnd` | `summarize-before-handoff` | continuity | warn |

`verification-present` is verify-walk's gate re-expressed as a prefilled policy — absorbed, not special-cased. (verify-walk itself is **not** refactored onto the engine in v1.12; this is the architectural direction, noted not forced.)

## 5. Policy anatomy

A policy block under a `## @<hook>` heading in POLICY.md:

```markdown
## @Implementation

### understand-before-change
- principle: 7          ← optional theory link
- severity: block       ← block | warn
- check: PLAN.md has a filled `## Understanding` section
         (How it works now / What changes / What stays the same),
         OR a UNDERSTANDING.md exists with the same three subheads

A request must not move planned → active until the agent has written
down how the system works today, what this change touches, and what
it leaves alone.
```

Schema: `hook` (the heading it lives under) + `principle?` (optional) + `severity` + `check` + prose.

**Severity is opt-in to blocking.** A policy blocks **only** if it explicitly declares `severity: block`. **Absent, `warn`, or unrecognized severity → non-blocking** (surface + continue). This is a safe default: a half-written or custom policy can never accidentally hard-stop a user — you must opt in to a gate. The 8 defaults ship with explicit `severity:` on all of them; the 4 blockers (`understand-before-change`, `verification-present`, `confirm-before-write`, `snapshot-before-overwrite`) say `block`, the other 4 say `warn`.

Whether a `check` is *mechanical* (doctor-verifiable, e.g. "`## Understanding` slot exists") or *judgment* (skill-evaluated, e.g. "goal is well-formed") is not a schema field — it's just how the check gets evaluated. `doctor policies` reports on the mechanical ones; judgment ones are enforced purely by the injected instruction (§8).

The `understand-before-change` check is backed by a new **`## Understanding` slot in PLAN.md** (`### How it works now` / `### What changes` / `### What stays the same`), escalating to a dedicated `requests/<slug>/UNDERSTANDING.md` for large requests. Satisfied by *either* location — the VERIFY.md 2-of-N pattern. **No `ANALYSIS.md`.**

## 6. Retrieval script

```
spectacular policy                    # all policies, grouped by hook (skim)
spectacular policy @Verification      # one hook's policies + their linked principle lines
spectacular policy <id>               # one policy, full text + its principle
spectacular policy --principle 7      # reverse: which policies enforce principle 7
spectacular policy --json             # machine form (skill-consumed)
```

`spectacular policy @<hook>` is the workhorse — the skill calls it on entering a phase. Output pulls **only** the hook's policies + each linked principle's heading and one line (not the whole PRINCIPLES.md). Skim-by-default, matching the `read-verbs` convention.

## 7. Injection loop

```
1. Phase entered      skill about to propose promote <slug> → verified
2. Retrieve           spectacular policy @Verification --json
3. Inject             returned policies + principle lines enter context
                      (nothing from other hooks loads — Principle 6)
4. Evaluate each check  mechanical → doctor/CLI presence-check
                        judgment  → skill reads the check prose
5. Resolve by severity  block fails → refuse transition, cite policy + principle
                        warn fails  → surface finding, continue
6. Done               context holds only this hook's rules
```

## 8. Mechanics

**Enforcement = an injected instruction at the head of each phase's reference doc.** No event bus, no new wiring. The skill already loads exactly one reference doc per phase (the SKILL.md routing model), so **the reference doc *is* the phase boundary.** Each phase ref doc opens with a 2-line gate block:

```markdown
> **@Planning policy gate.** Before anything else, run `spectacular policy @Planning`
> and follow every active policy returned. `block` → satisfy or stop; otherwise → surface and continue.
```

When the skill loads that doc to do the phase, the first thing it reads is the instruction to consult policies. The gate fires precisely when the phase begins, by construction. Per-phase placement:

| Hook | Gate block goes at the top of |
|---|---|
| `@Init` | `init-workflow.md` |
| `@Planning` | `new-request.md` |
| `@Implementation` | `active-request.md` / `lifecycle.md` (→active) |
| `@Verification` | `verification.md` / `lifecycle.md` (review→verified) |
| `@Archive` | `archive.md` |
| `@Remember` | `memory.md` |
| `@Snapshot` | `versioning.md` |
| `@SessionEnd` | `sessions-rules.md` (session end flow) |

**Source of truth: POLICY.md; config is an override layer.** POLICY.md *defines* each policy (hook, check, severity, prose). `config.yaml`'s `policies:` block *tunes* it for this project (enable/disable/change severity). They are layers, not competing copies — `spectacular policy` reads POLICY.md, applies config overrides, returns the merged result. No "which wins" ambiguity.

**POLICY.md (and PRINCIPLES.md) are CLI-managed *and* hand-editable, but structure-bound.** The CLI scaffolds and can mutate them; a human can edit directly; both must obey the declared structure (`## @hook` sections + policy-block anatomy). The `review` verb on POLICY.md = a **structure check** (every block has required fields, hooks are from the known set, no orphan sections) — richer than `principles-rules.md`'s placeholder check, same spirit.

## 9. Milestones

- **M1 — Practice-layer contract.** `references/policies-contract.md`: POLICY.md structure (sections keyed by `## @hook`), policy anatomy (`hook + principle? + severity + check + prose`), the locked hook set, and how `config.yaml` enables/disables/tunes. Scope model: **config-only** for v1 (4-tier precedence noted as v2 candidate).
- **M2 — POLICY.md scaffold + 8 defaults.** New `doc_policy()` in `cli/spectacular` emitting POLICY.md with all 8 prefilled policies; add POLICY.md to the **always-set** list. Spec the `## Understanding` PLAN slot + `UNDERSTANDING.md` escalation; register the slot in `plan-overrides.md` + `scaffold-reference.md`.
- **M3 — `spectacular policy` verb + injection.** Implement the retrieval script (all five forms). Add the **gate block** (`run spectacular policy @<hook>, follow active policies`) to the head of each phase's reference doc per the §8 mapping. Write `references/policy-injection.md` documenting the loop.
- **M4 — Enforcement + config.** Wire policy consultation into `promote`/`archive` at the spine hooks; add a `doctor policies` area for presence-checkable policies. `config.yaml` `policies:` block — per-policy enable/disable/severity + register custom; worked custom-policy example.
- **M5 — Dogfood + ship.** Enable POLICY.md on this repo (it's already dogfooding); drive 1+ real request through `@Implementation` with a filled `## Understanding`; CHANGELOG + plugin bump to v1.12.0.

## 10. Tasks

See `TASKS.md`.

## 11. Dependencies

- **Sequenced after [[verify-walk]] (v1.11).** verify-walk established the skill-side "walk + gate + lifecycle-flip" pattern at `review → verified`; policy-engine generalizes it into the named-hook practice layer. Its gate becomes the `verification-present` prefilled policy.
- Touches [[lifecycle]] (`promote`/`archive` consult policies), [[doctor]] (new `policies` area), and the always-set scaffold (POLICY.md joins it).
- Relates to [[PRINCIPLES]] — POLICY.md is its executable practice sibling; the optional `principle:` tag links the two.

## 12. Validation

- M1 — `policies-contract.md` documents POLICY.md structure + policy anatomy + hook set; a reader can author a custom policy from it alone.
- M2 — `spectacular init` creates POLICY.md with 8 prefilled policies; the `## Understanding` slot is registered and scaffolded.
- M3 — `spectacular policy @<hook>` returns that hook's policies + linked principle lines; `spectacular policy --principle N` reverse-resolves; `--json` is machine-parseable.
- M4 — `promote → active` consults `@Implementation`; a request lacking `## Understanding` is **blocked**; a request with it passes. `config.yaml` can disable a policy and change a severity — each takes effect. `doctor policies` reports status.
- M5 — This repo runs with POLICY.md enabled; a real request shows `## Understanding` was required before `active`; manifests at v1.12.0.

## 13. Deliverables

- `POLICY.md` scaffold (8 prefilled policies) + always-set wiring
- `references/policies-contract.md` — POLICY.md structure, policy anatomy, hook set, config surface
- `spectacular policy` verb (5 forms) + injection-loop reference doc
- `## Understanding` PLAN slot + `UNDERSTANDING.md` escalation; registered in `plan-overrides.md` + `scaffold-reference.md`
- `promote`/`archive` policy consultation + `doctor policies` area
- `config.yaml` `policies:` block + worked custom-policy example
- CHANGELOG [1.12.0] entry; plugin bump

## Locked decisions (2026-05-30, decision-tree walk)

1. **POLICY.md is the practice layer** — sibling to PRINCIPLES.md (theory→practice). Policies filed under named work-phase hooks; skill retrieves by hook.
2. **Hook naming** — `@` = "at" + moment-noun, reads as natural English ("applies @implementation"). 8 hooks locked (§3). `@Request`/`@RequestTask` folded into `@Planning`.
3. **POLICY.md = single file, no soft-DB** in v1.
4. **POLICY.md = always-set** — every init, 8 prefilled policies enabled. Asymmetric with optional PRINCIPLES.md, on purpose.
5. **Severity split** — 4 block (understand-before-change, verification-present, confirm-before-write, snapshot-before-overwrite), 4 warn (scaffold-contract, request-shape, spec-sync+memory-propose, summarize-before-handoff).
6. **Understanding content** — `## Understanding` PLAN slot, escalates to `UNDERSTANDING.md`; gate checks either (2-of-N). No `ANALYSIS.md`.
7. **Principle link** — optional `principle:` tag; retrieval pulls the one relevant line.
8. **Enforcement** — skill-side + doctor; no `hooks.json` harness wiring in v1 (v2 upgrade path for kernel-level locks).
9. **Scope model** — config-only in v1; 4-tier precedence deferred to v2.
10. **verify-walk** — absorbed as the `verification-present` policy; not refactored onto the engine in v1.12.
11. **Severity is opt-in to blocking** — a policy blocks only if it explicitly says `severity: block`; absent/`warn`/unrecognized → non-blocking. Safe default: no policy accidentally hard-stops.
12. **Enforcement mechanic** — a gate block at the head of each phase's reference doc tells the skill to run `spectacular policy @<hook>` first. The ref doc *is* the phase trigger; no event bus, no SKILL.md special-casing. (Replaces the earlier "skill-side + doctor, somehow wired" hand-wave.)
13. **Source of truth + editability** — POLICY.md defines policies; `config.yaml` is an override layer (not a competing copy). Both POLICY.md and PRINCIPLES.md are CLI-managed *and* hand-editable but **structure-bound**; `review` on POLICY.md is a structure check. No `check-kind` schema field — mechanical-vs-judgment is just how a check is evaluated (doctor vs skill).
