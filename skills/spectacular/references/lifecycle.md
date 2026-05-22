# Lifecycle — State Transitions

State machine: `planned → active → review → verified → archived`

State lives in `PLAN.md` frontmatter: `status: planned | active | review | verified`

---

## Transition rules

The skill detects signals and **proposes** transitions. The user can also force transitions explicitly.

| From | To | Signal | Skill action |
|---|---|---|---|
| `planned` | `active` | User starts working on request | Create SESSION.md, update frontmatter |
| `active` | `review` | All TASKS.md items checked | Propose move to review; also propose VERIFY.md if 2-of-6 rule triggers (see [[verification]]) |
| `active` | `review` | User says "done" / "ready to review" | Propose move to review |
| `review` | `verified` | All `- [x]` in VERIFY.md (when present) | Propose verified status |
| `review` | `verified` | All TASKS § Verification + PLAN § Validation items confirmed (when no VERIFY.md) | Propose verified status |
| `review` | `verified` | User confirms everything works | Update frontmatter |
| `verified` | `archived` | User runs `spectacular archive <slug>` | See `archive.md` |

**Never auto-transition.** Always propose and wait for user confirmation.

---

## Signal detection details

### All tasks checked → propose review

When reading TASKS.md and all `- [ ]` items are now `- [x]`:

> "All tasks in `<slug>` are checked. Ready to move to `review`?"

If moving to review, also evaluate the 2-of-6 rule (see [[verification]]) — only propose creating VERIFY.md when the request actually warrants one. Default for doc-only / refactor / spec requests is **no VERIFY.md** — use PLAN § Validation or add a `### Verification` group to TASKS.md instead.

### Verification artifact detection (review → verified)

Verification is **never skipped**. The skill always checks against some artifact before allowing `verified`. The only question is which artifact:

1. **VERIFY.md exists** — load-bearing. Every `- [ ]` blocks the transition. The skill never moves to verified with unchecked items.
2. **No VERIFY.md, TASKS has `### Verification` group** — every item in that group must be `- [x]`. Blocks transition until checked.
3. **Neither exists** — every PLAN.md § Validation item must be explicitly confirmed by the user before transition.

Per [[verification]], "opt-in" refers to **whether a standalone VERIFY.md file gets scaffolded** — not whether verification runs. The 2-of-6 rule decides the file; verification itself is mandatory.

### Stale request detection

If a request's `updated` date is >14 days old and `status` is `active`:

> "`<slug>` has been active since <date> with no updates. Still in progress? Want to add a blocker note or update the status?"

### Draft capability with no active request

If a `current/` spec has `status: draft` and no `requests/` item references it:

> "`current/<capability>` is still draft. Want to create a request to finish it, or mark it deprecated?"

---

## Forcing transitions

User can explicitly say:
- "mark `<slug>` as active" → update frontmatter, create SESSION.md
- "move `<slug>` to review" → update frontmatter
- "mark `<slug>` as verified" → update frontmatter
- "archive `<slug>`" → route to `archive.md`

---

## Capability spec states

`current/<capability>.md` frontmatter tracks its own state: `status: stable | draft | deprecated`

- `draft` — capability spec exists but is being developed (often tied to an active request)
- `stable` — current canonical truth, no active changes
- `deprecated` — no longer in use, kept for reference

Skill proposes `current/` updates when a request is archived (see `current-sync.md`).
