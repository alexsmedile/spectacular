---
doc-id: idea
mode: index
location: .spectacular/ideas/
entries-dir: .spectacular/ideas/
scope: project-wide
template: templates/idea/base.md
snapshot-on-edit: false
summary: "Thinking scratchpad — captured ideas not yet committed to a request. Soft-folder DB; promotion to a request is explicit and one-way."
status: active
---

# Idea Rules

Soft-folder database of captured ideas. There is **no top-level `IDEAS.md` index file** — folder listing is the canonical view. Each idea is a self-contained markdown file.

**Mode: `index`** (no regenerated index file in v1.7.0). Entry files at `entries-dir`.

**Philosophy (carried over from ARCHITECTURE.md § ideas/):** ideas are a **thinking scratchpad, not a workflow stage**. Nothing in `ideas/` is acted on automatically by the skill. The skill creates, lists, doctor-checks, and promotes — it never auto-routes ideas into requests, never re-prioritizes them, never grills them unsolicited.

**Verbs:**
- `grill` → interactive form to flesh out an idea's slots (Hypothesis / Context / Open questions). Useful when an idea is moving from `parked` → `exploring`. User-initiated only.
- `refine` → rewrite vague slots into specific ones. Optional.
- `review` → validate frontmatter shape across all entries; flag stale `exploring` (>90 days), orphan promoted entries, unknown status values.

**Mutator verbs (CLI, not skill):**
- `spectacular idea new <slug>` — scaffold one entry from template with `status: parked`, `updated:` today
- `spectacular idea list [--status <state>]` — list entries with status + last-updated date
- `spectacular idea promote <slug>` — scaffold a new request from the idea (delegates to [[new-request]] flow), move source file to `archive/ideas/<slug>.md`, set its `status: promoted`

**Snapshot-on-edit: false** — ideas are scratchpad records, not versioned canonical docs. They mature by being promoted, not by accumulating snapshots.

**Entry frontmatter (required shape):**

```yaml
---
type: idea
status: parked | exploring | promoted
priority: low | medium | high
owner: <name>
origin: <free-text — where it came from: conversation, side-thought, abandoned request, etc.>
updated: YYYY-MM-DD
promoted_to: requests/<slug>/ | null
related: []
---
```

**Required body sections (template-enforced, not gate-checked):** Hypothesis, Context, Open questions, Promoted-to (placeholder until promotion).

## Status lifecycle

```
parked ──(start shaping)──► exploring ──(promote)──► promoted
                              │
                              └─(let cool)──► parked
```

- **`parked`** — captured. Not actively being shaped. Default state at creation.
- **`exploring`** — actively thinking; slots being filled. Doctor flags this state if `updated:` is >90 days old (decide: promote, demote to parked, or delete).
- **`promoted`** — became a request. File should live in `archive/ideas/<slug>.md`, not `.spectacular/ideas/`. Doctor flags promoted entries still in the live folder as orphans.

## Promotion to request

When the user runs `spectacular idea promote <slug>`:

1. CLI reads `.spectacular/ideas/<slug>.md` and extracts frontmatter + body content
2. Hands off to the [[new-request]] flow — request scaffolded with PLAN.md pre-filled from idea content
3. Sets `promoted_to: requests/<slug>/` on the idea file
4. Sets `status: promoted`
5. Moves file to `.spectacular/archive/ideas/<slug>.md`
6. Notes in the new PLAN.md: `promoted from ideas/<slug>.md`

Promotion is **one-way and explicit** — never automatic, never reversible by tool (manual revert if needed).

## Doctor area

`spectacular doctor ideas` is **judgment-only** (no `--fix`). Flags:

| Check | Severity | Condition |
|---|---|---|
| Stale exploring | warning | `status: exploring` + `updated:` >90 days |
| Orphan promoted | warning | `status: promoted` but file still in `.spectacular/ideas/` (should be in `archive/ideas/`) |
| Missing required frontmatter | warning | `type`, `status`, `updated` absent or empty |
| Unknown status value | warning | `status:` not one of `parked\|exploring\|promoted` |

No `--fix` because every finding requires a human decision (promote? demote? delete? move?). Mechanical auto-moves on `promoted` orphans would conflict with the "explicit and one-way" promotion contract.

## Aliases

None in v1.7.0. The word "idea" is short enough and the verb surface clear enough that hidden routing isn't warranted (unlike feedback-loop's `iterate|experiment|test|probe|try`).

## What this is **not**

- **Not a request.** Ideas don't carry PLAN/TASKS/lifecycle state. Promotion produces a request; the idea itself never gains those structures.
- **Not memory.** MEMORY.md is for durable preferences and decisions. Ideas are *speculative* — they may never matter again.
- **Not a backlog.** ROADMAP.md's `## Icebox` section holds vision-tier items tied (loosely) to versions. Ideas have no version, no scope commitment, no exit criteria.
- **Not feedback.** Feedback entries answer "was that the right thing to ship?" — they're tied to something already built. Ideas are pre-commitment scratchpad for things that *might* get built.

**Related:** [[new-request]] (promotion flow), [[archive]] (where promoted ideas land), [[doc-index]], [[scaffold-reference]], [[roadmap-rules]] (icebox vs ideas), [[memory-rules]] (durable vs speculative), [[feedback-rules]] (post-ship vs pre-commit).
