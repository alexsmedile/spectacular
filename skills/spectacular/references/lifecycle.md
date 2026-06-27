---
description: Lifecycle state transitions, signal detection, proactive proposals.
when_to_use: A request changes state (planned/active/review/verified).
---

# Lifecycle — State Transitions

State machine: `planned → active → review → verified → archived`

State lives in `PLAN.md` frontmatter: `status: planned | active | review | verified`

> **Policy gates on transitions.** Two transitions consult the policy engine before flipping state:
> - `planned → active` → run `spectacular policy @Implementation` (blocker: `understand-before-change`).
> - `review → verified` → run `spectacular policy @Verification` (blocker: `verification-present`).
>
> Satisfy every `block` policy or refuse the transition; surface `warn` policies and continue. See [policy-injection.md](policy-injection.md).

---

## Transition rules

The skill detects signals and **proposes** transitions. The user can also force transitions explicitly.

| From | To | Signal | Skill action |
|---|---|---|---|
| `planned` | `active` | User starts working on request | Run `spectacular advance <slug>`; create SESSION.md |
| `active` | `review` | All TASKS.md items checked | Propose `spectacular advance <slug>`; also propose VERIFY.md if 2-of-6 rule triggers (see [[verify]]) |
| `active` | `review` | User says "done" / "ready to review" | Run `spectacular advance <slug>` |
| `review` | `verified` | All `- [x]` in VERIFY.md (when present) | Propose `spectacular advance <slug>` |
| `review` | `verified` | All TASKS § Verification + PLAN § Validation items confirmed (when no VERIFY.md) | Propose `spectacular advance <slug>` |
| `review` | `verified` | User confirms everything works | Run `spectacular advance <slug>` |
| `verified` | `archived` | User confirms archive | Run `spectacular archive <slug>`; see [[archive]] |

> **Verb name (v1.19.0):** the lifecycle verb is `spectacular advance` (was `promote`). `promote` still works as a deprecated alias and prints a one-line notice. Distinct from `spectacular idea promote`, which promotes an *idea* into a request — that one keeps its name.

**Mutation principle (v0.7.0+):** state changes use `spectacular advance <slug>`. The CLI:
- Reads current status from PLAN.md frontmatter
- Refuses backward transitions without `--force`
- Sets `status:` + `updated:` atomically in PLAN.md AND TASKS.md
- Optional `--to <state>` for explicit jumps (e.g. `--to verified`)
- Optional `--archive` chains into `spectacular archive` after promoting to verified

Never hand-edit `status:` in PLAN.md frontmatter when the CLI verb covers it. Manual edits are for edge cases the verb doesn't handle.

**Never auto-transition.** Always propose and wait for user confirmation before running the verb.

---

## Signal detection details

### All tasks checked → propose review

When reading TASKS.md and all `- [ ]` items are now `- [x]`:

> "All tasks in `<slug>` are checked. Ready to move to `review`?"

If moving to review, also evaluate the 2-of-6 rule (see [[verify]]) — only propose creating VERIFY.md when the request actually warrants one. Default for doc-only / refactor / spec requests is **no VERIFY.md** — use PLAN § Validation or add a `### Verification` group to TASKS.md instead.

### Verification artifact detection (review → verified)

Verification is **never skipped**. The skill always checks against some artifact before allowing `verified`. The only question is which artifact:

1. **VERIFY.md exists** — load-bearing. Every `- [ ]` blocks the transition. The skill never moves to verified with unchecked items.
2. **No VERIFY.md, TASKS has `### Verification` group** — every item in that group must be `- [x]`. Blocks transition until checked.
3. **Neither exists** — every PLAN.md § Validation item must be explicitly confirmed by the user before transition.

**Substrate check (auto-invoked):** when the skill proposes `verified`, also run `spectacular doctor lifecycle` scoped to that request — confirms the verification artifact exists per the convention. If doctor reports an error, abort the transition with the finding.

Per [[verify]], "opt-in" refers to **whether a standalone VERIFY.md file gets scaffolded** — not whether verification runs. The 2-of-6 rule decides the file; verification itself is mandatory.

### Stale request detection

If a request's `updated` date is >14 days old and `status` is `active`:

> "`<slug>` has been active since <date> with no updates. Still in progress? Want to add a blocker note or update the status?"

### Draft capability with no active request

If a `specs/` capability has `status: draft` and no `requests/` item references it:

> "`specs/<capability>` is still draft. Want to create a request to finish it, or mark it deprecated?"

---

## Forcing transitions

User can explicitly say:
- "mark `<slug>` as active" → `spectacular advance <slug> --to active`; create SESSION.md
- "move `<slug>` to review" → `spectacular advance <slug> --to review`
- "mark `<slug>` as verified" → `spectacular advance <slug> --to verified`
- "archive `<slug>`" → `spectacular archive <slug>` (also see [[archive]])

Backward transitions (e.g. `verified → active` because verification failed) require `--force`:
- `spectacular advance <slug> --to active --force`

The `--force` flag is intentionally awkward — backward moves should be rare and deliberate.

---

## Capability spec states

`specs/<capability>/SPEC.md` frontmatter tracks its own state: `status: stable | draft | deprecated`

- `draft` — capability spec exists but is being developed (often tied to an active request)
- `stable` — current canonical truth, no active changes
- `deprecated` — no longer in use, kept for reference

Skill proposes `specs/` updates (and a bullet edit to `SPEC.md` index) when a request is archived (see `spec-sync.md`).
