---
description: Scaffold a new request — slug rules, templates, target-version.
when_to_use: spectacular new <description> or spectacular idea promote.
---

# New Request — Scaffolding

Triggered by: `spectacular new <description>`, or when conversational context clearly indicates a new piece of work to track.

> **@Planning policy gate.** First, run `spectacular policy @Planning` and follow every active policy returned. A `block` policy must be satisfied or you stop; a `warn` policy you surface and continue. See [policy-injection.md](policy-injection.md).

---

## Autopilot vs interactive

| Context | Behavior |
|---|---|
| Clear description in conversation | Derive slug, draft PLAN.md, show user, confirm, write |
| Thin context | Ask targeted questions (goal? scope? why now?), then scaffold |

Always show the derived slug before writing. User can override before confirmation.

---

## Slug rules

1. Derive from the request description — kebab-case, concise (2–4 words)
2. Check `requests/` for existing slugs — if collision, propose appending `-2` or ask user to rename
3. Apply `naming.prefix` from `config.yaml` if set
4. User can override at any time before write

Example: "add team billing" → `add-team-billing`

---

## Scaffold sequence

**Mutation principle (v0.7.0+):** scaffolding goes through `spectacular new <slug>`. The CLI:
- Validates the slug (kebab-case, max 64 chars; respects pack naming rules if active)
- Refuses if slug already exists in `requests/` OR `archive/`
- Reads PLAN.md and TASKS.md from the bundled templates
- Prefills frontmatter via shared helpers (status, priority, owner from config, updated, summary)
- Writes both files atomically

Skill flow:

1. Derive the slug from the user's description (kebab-case, 2–4 words)
2. Apply slug rules + check for collisions (see below)
3. Confirm slug + summary with user
4. Run: `spectacular new <slug> --summary "<desc>" [--priority high|low]`
5. Inspect the scaffolded PLAN.md; edit slot bodies as needed (slots are skill territory; the *files* are CLI territory)

On confirmation, the verb creates:

```
requests/<slug>/
├── PLAN.md       ← always created
└── TASKS.md      ← always created
```

> **Tier-reveal (next step).** After scaffolding, surface exactly one suggestion: `Next: spectacular plan grill <slug> to stress-test the PLAN before building.` One line, only after the scaffold is confirmed — never mid-flow, never a menu.

Create only on demand (skill proposes, user confirms):
- `SESSION.md` — when request moves to `active`
- `RISKS.md` — when request touches auth, billing, migrations, or anything flagged sensitive in STACK.md
- `VERIFY.md` — when the **2-of-6 rule** in [[verification]] triggers. Default for doc/refactor/spec requests is **no VERIFY.md** — verification lives in PLAN § Validation or TASKS § Verification instead. The file is opt-in; **the practice is not** — every request must reach `verified` through some artifact.
- `artifacts/` — when screenshots, benchmarks, or research need storing. **Pack consultation:** if a convention pack is active (config.yaml `convention_pack:` declared), the skill reads the pack's `file-placement.request-artifacts:` rule for the artifact subdirectory layout. Default when no pack: `artifacts/<kind>/`.

### Verification routing at scaffold time

Before finalizing the PLAN.md, the skill should:
1. Apply the 2-of-6 rule from [[verification]] against the request's nature
2. If yes — propose `VERIFY.md` scaffold; populate PLAN § Validation with milestone checkpoints only
3. If no — populate PLAN § Validation with full per-milestone checks; optionally add `### Verification` group to TASKS.md for procedural items
4. Either way — verification is documented somewhere before `active` ends

---

## PLAN.md template (7-slot decomposition)

```md
---
status: planned
priority: medium
owner: 
updated: <today>
summary: "<one-line description>"
related:
  - ../../PRD.md
  - specs/<capability>/SPEC.md
---

# <Request title>

## Goal
<One sentence — compressed intent. Traces back to a success criterion in the root PRD.>

## Why (intent)
<Why now? What problem does this solve? The full why lives in PRD.md — keep this tight.>

## Constraints
- <What's fixed before starting — budget, tech, policy>

## Milestones
1. **<Milestone name>** — <demoable outcome, not a task>
2. ...

## Tasks
See [TASKS.md](TASKS.md).

## Dependencies
- <Other requests, skills, blocking decisions — or "none">

## Validation
| Milestone | How we verify it passed |
|---|---|
| 1 | <test, demo, review> |

## Deliverables
- <Artifact that ships out of this request>

## Open questions
- <Things you don't know yet>
```

**Anti-pattern:** never create `requests/<slug>/PRD.md`. Product intent is project-wide and lives at `.spectacular/PRD.md`. If a request needs to extend or revise product intent, edit the root PRD (snapshot first) — don't fork it into a request folder.

---

## TASKS.md template

```md
---
updated: <today>
---

# Tasks — <slug>

## <Group name>

- [ ] 
- [ ] 

## <Group name>

- [ ] 
- [ ] 
```

---

## Promoting from an idea

When user runs `spectacular idea promote <idea-file>`:

1. Read `ideas/<idea-file>.md` for content — use it to pre-fill PLAN.md goal, why, approach
2. Scaffold request as above (show slug, confirm)
3. Move `ideas/<idea-file>.md` → `archive/ideas/<idea-file>.md`
4. Note in PLAN.md: `promoted from ideas/<idea-file>.md`
