---
description: POLICY.md structure, policy anatomy, the locked 9-hook set, and the config override surface. Read to author or audit a policy.
when_to_use: Authoring a custom policy, auditing POLICY.md structure, or understanding how policies are filed under work-phase hooks.
---

# Policies Contract — the practice layer

`POLICY.md` is Spectacular's **practice layer** — the operational sibling to `PRINCIPLES.md`.

| Doc | Layer | Answers | Tier |
|---|---|---|---|
| `PRINCIPLES.md` | theory | *why we work this way* (beliefs) | optional (kit-triggered) |
| `POLICY.md` | practice | *how we actually work* (executable rules) | **always-set** (every init) |

The asymmetry is deliberate: theory is optional reading, **practice is the operational floor**. A workspace can hold zero stated principles, but it always ships with a policy contract — 19 prefilled defaults, enabled, from the first `spectacular init`.

A **policy** is a rule filed under a named **work-phase hook**. When the skill enters a phase, it retrieves *only that hook's* policies and injects them — progressive disclosure (Principle 6) applied to the rule layer itself. The skill never loads all policies at once.

## POLICY.md structure

The file is a set of `## @<hook>` sections. Each section holds zero or more policy blocks (`### <policy-id>`):

```markdown
## @Implementation

### understand-before-change
- principle: 7
- severity: block
- check: PLAN.md has a filled `## Understanding` section
         (How it works now / What changes / What stays the same),
         OR a UNDERSTANDING.md exists with the same three subheads

A request must not move planned → active until the agent has written
down how the system works today, what this change touches, and what
it leaves alone.
```

Section heading = the hook. Policy id = the `###` slug (kebab-case, carries the before/after verb). Everything after the field list, until the next `###` or `##`, is the policy's **prose** — the human-readable rationale and the instruction the skill follows when injected.

## Policy anatomy

A policy block has five parts:

| Part | Required | Form | Meaning |
|---|---|---|---|
| `hook` | yes (implicit) | the `## @<hook>` it lives under | when this policy fires |
| `principle` | optional | `- principle: N` | the PRINCIPLES.md § it enforces; retrieval pulls that one line |
| `severity` | yes for blockers | `- severity: block \| warn` | how a failed check resolves |
| `check` | yes for blockers | `- check: <text>` | the condition that must hold |
| `directive` | recommended | `- directive: <one imperative sentence>` | what the CLI injects at the hook gate — the practice layer, verbatim |
| prose | yes | free text after the fields | rationale + the instruction injected into context |

**`hook`** is not written as a field — it *is* the section heading the block sits under. Moving a policy = moving its `###` block under a different `## @<hook>`.

**`principle`** is the one optional link between practice and theory. `- principle: 7` means "this policy enforces Principle 7." On retrieval, `spectacular policy @<hook>` pulls Principle 7's heading + one line alongside the policy — not the whole PRINCIPLES.md. Omit it for policies with no theory backing (e.g. continuity hygiene).

**`severity`** is **opt-in to blocking.** A policy blocks *only* if it explicitly declares `severity: block`. **Absent, `warn`, or any unrecognized value → non-blocking** (surface the finding and continue). This is a safe default: a half-written or custom policy can never accidentally hard-stop a user — you must opt in to a gate.

**`check`** is the condition. A **block** policy must declare one (the gate needs something concrete to evaluate); a **warn** policy may omit it and rely on its prose instruction alone (most hygiene nudges do). Whether a check is *mechanical* (doctor-verifiable, e.g. "`## Understanding` slot exists and is filled") or *judgment* (skill-evaluated, e.g. "the goal is well-formed") is **not a schema field** — it's simply *how the check gets evaluated*:

- **Mechanical** → `doctor policies` presence-checks it; the result is deterministic.
- **Judgment** → enforced purely by the injected instruction; the skill reads the check prose and decides.

No `check-kind` field exists. The same policy may even be checked both ways over time.

**`directive`** is the one-sentence injection line — the policy's own instruction, written imperatively, that `spectacular policy @<hook>` prints on every row (the gate output agents actually read). Without it the CLI falls back to the linked principle's title, so migration is per-policy incremental — but the directive is the *practice* layer speaking in its own voice, so author one for every policy that matters at a gate. Tiering: **warn** rows show directive + `P<n> — <title>`; **block** rows show directive + the full principle line (a refusal must carry its reasoning); `--full` restores full paragraphs everywhere. The prose stays the full rationale — the directive is its sharpest sentence, not a replacement.

## The locked hook set (9)

`@` reads "at"; every hook completes "this policy applies **@___**" as natural English. The before/after verb lives in the *policy name*, never the hook.

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
| `@Debugging` | at debugging | a bug/quirk/regression is reported | `spectacular audit` / `fix`, loads [[bug-workflow]] |
| `@Remember` | at remember | memory written | `spectacular remember` |
| `@Snapshot` | at snapshot | canonical doc overwritten | `snapshot` / overwrite |
| `@SessionEnd` | at session end | skill hands off | `session end` |

These 9 are the **only** valid hooks. `review` on POLICY.md flags any `## @<hook>` section outside this set as an orphan. `@Debugging` is not part of the lifecycle spine — it's a key-moment hook entered whenever a bug surfaces, independent of a request's lifecycle state.

**Deferred (v2):** `@SessionStart` (wants a real harness runtime), `@Decide`. **Folded:** `@Request`/`@RequestTask` → `@Planning`. **Rejected:** `@Doctor` (circular), `@BeforeCommit` (wrong runtime).

## The prefilled defaults — 4 block · 15 warn

Every `spectacular init` writes these enabled:

| Hook | Policy | →Principle | Severity |
|---|---|---|---|
| `@Init` | `scaffold-contract` | 4 | warn |
| `@Planning` | `request-shape` + `scope-down` + `milestones-in-build-order` | 3 / 7 / 10 / 11 | warn |
| `@Implementation` | **`understand-before-change`** | 7 | **block** |
| `@Implementation` | `build-order` + `earn-the-verification` + `prefer-cli-mutator` | 11 / 6 | warn |
| `@Debugging` | `check-prior-fixes` + `ceremony-matches-uncertainty` + `fix-root-not-symptom` + `log-only-verified-reusable` + `use-audit-fix-verbs` | 5 / 11 / 6 | warn |
| `@Verification` | **`verification-present`** | 7 / 9 | **block** |
| `@Archive` | `spec-sync` + `memory-propose` | 2 / 5 | warn |
| `@Remember` | **`confirm-before-write`** | 8 | **block** |
| `@Snapshot` | **`snapshot-before-overwrite`** | 8 | **block** |
| `@SessionEnd` | `summarize-before-handoff` | continuity | warn |

`verification-present` is verify-walk's gate re-expressed as a prefilled policy — absorbed, not special-cased. (verify-walk itself is *not* refactored onto the engine in v1.12; this is the architectural direction, noted not forced.)

The headline blocker — `understand-before-change` — is backed by a `## Understanding` slot in `PLAN.md` (`### How it works now` / `### What changes` / `### What stays the same`), escalating to a dedicated `requests/<slug>/UNDERSTANDING.md` for large requests. Satisfied by *either* location (the VERIFY.md 2-of-N pattern). There is no `ANALYSIS.md`.

## Config surface — POLICY.md defines, config.yaml overrides

POLICY.md is the **source of truth**: it defines each policy's hook, check, severity, and prose. `config.yaml`'s `policies:` block is an **override layer** that tunes the contract for this project — it never duplicates the definition.

```yaml
# config.yaml
policies:
  understand-before-change:
    enabled: true        # default true; set false to disable a prefilled policy
  scaffold-contract:
    severity: block      # override the shipped severity (warn → block here)
  no-secrets-in-memory:  # register a custom policy by id
    hook: "@Remember"
    severity: warn
    check: "memory entry contains no API keys, tokens, or passwords"
```

`spectacular policy` reads POLICY.md, applies these overrides, and returns the merged result. They are **layers, not competing copies** — there is no "which wins" ambiguity: POLICY.md is the definition, config is the per-project tune.

**Scope model: config-only for v1.** A single `policies:` block in `.spectacular/config.yaml`. The 4-tier scope precedence (project → user → app-store → bundled) used by convention packs is a **v2 candidate**, noted not built.

## Editability

POLICY.md (and PRINCIPLES.md) are **CLI-managed *and* hand-editable, but structure-bound.** The CLI scaffolds and can mutate them; a human can edit directly; both must obey the structure above (`## @hook` sections from the known set + the policy-block anatomy). The `review` verb on POLICY.md is a **structure check** — every block has a `check`, severities are `block|warn`, hooks are from the known 8, no orphan sections — richer than `principles-rules.md`'s placeholder check, same spirit.

## Authoring a custom policy

1. Pick the hook — when should it fire? (one of the 8)
2. Add a `### <verb>-<noun>` block under that `## @<hook>` section in POLICY.md.
3. Write the `check` (the condition) and prose (rationale + instruction).
4. Set `severity: block` only if a failure must hard-stop; otherwise leave it `warn` or absent.
5. Optionally tag `- principle: N` to link the theory it enforces.
6. (Optional) override or disable it later via `config.yaml`'s `policies:` block.

See [policy-injection.md](policy-injection.md) for how a policy is retrieved and enforced at runtime, and [policy-rules.md](policy-rules.md) for the doc-verb dispatch.
