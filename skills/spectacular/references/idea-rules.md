---
doc-id: idea
mode: index
location: .spectacular/ideas/ (compatibility and local/private refinement)
entries-dir: .spectacular/ideas/
scope: project-wide
template: templates/idea/base.md
snapshot-on-edit: false
summary: "Compatibility workbench for private/offline or locally useful refinement before an explicit handoff to request, roadmap, or a shared destination."
status: active
---

# Idea Rules

Compatibility soft-folder database for captured local ideas. There is **no top-level `IDEAS.md` index file** — folder listing is the canonical local view. Each idea is a self-contained markdown file named `IDEA-NNN-<slug>.md`; users normally refer to it as `I<N>`. Legacy unnumbered files remain readable until migrated. See [[canonical-ids]].

**Mode: `index`** (no regenerated index file in v1.7.0). Entry files at `entries-dir`.

**Philosophy (carried over from ARCHITECTURE.md § ideas/):** an idea is a **pre-commitment compatibility workbench, not an execution stage or default backlog**. Capture it where it will actually be developed and shared whenever that place is known. Use a local IDEA only for private/offline capture or locally useful refinement. Nothing is acted on automatically: the skill never auto-routes an idea, re-prioritizes it, or starts grilling it unsolicited.

## Capture sources and authority

An idea may start in a GitHub Issue, GitHub Discussion, another shared tracker, phone note, `TODO.md`, `FEEDBACK.md`, `.spectacular.local/ideas/`, or a committed `IDEA`. Capture is deliberately flexible: use the lowest-friction place that is appropriate for its visibility and sensitivity.

Do **not** mirror every capture into `ideas/`. While an Issue or Discussion is being discussed, it remains authoritative for that conversation. Create a committed `IDEA` only when local, agent-readable refinement is useful; record the source in `origin:`/`related:` and link rather than copy it. A private phone note can remain local until the human chooses to publish or promote it.

**Verbs:**
- `grill` → interactive form to flesh out an individual idea's hypothesis, context, open questions, evidence links, alternatives, and optional working plan. Useful when an idea is moving from `parked` → `exploring`. User-initiated only.
- `refine` → rewrite vague slots into specific alternatives, decisions still needed, and plan assumptions. Optional.
- `review` → validate frontmatter shape across all entries; flag stale `exploring` (>90 days), orphan promoted entries, unknown status values.

**Mutator verbs (CLI, not skill):**
- `spectacular idea new <slug>` — scaffold one entry from template with `status: parked`, `updated:` today
- `spectacular idea list [--status <state>]` — list entries with status + last-updated date
- `spectacular idea promote <slug> --to request` — scaffold a new request from the idea (delegates to [[new-request]]), then archive the source.
- `spectacular idea promote <slug> --to roadmap --placement icebox` — snapshot the roadmap and add an explicit Icebox handoff, then archive the source. Choosing a version/tier remains the manual roadmap ritual.
- `spectacular idea promote <slug> --to shared --ref <stable-reference>` — archive the local source with the existing shared destination reference. It makes **no network call** and never creates, comments on, or updates a remote record.

Omitting `--to` remains a warning-producing compatibility alias for `--to request`.

**Snapshot-on-edit: false** — ideas are scratchpad records, not versioned canonical docs. They mature by being promoted, not by accumulating snapshots.

**Entry frontmatter (required shape):**

```yaml
---
type: idea
id: IDEA-NNN
status: parked | exploring | promoted
priority: low | medium | high
owner: <name>
origin: <free-text — where it came from: conversation, side-thought, abandoned request, etc.>
updated: YYYY-MM-DD
promoted_to: requests/<slug>/ | roadmaps/index.md#icebox | shared:<reference> | null
related: []
---
```

**Required body sections (template-enforced, not gate-checked):** Hypothesis, Context, Open questions, Working plan, Promoted-to (placeholder until promotion). A Working plan is exploratory: it must name assumptions and unanswered decisions rather than masquerade as approved execution scope.

## Status lifecycle

```
parked ──(start shaping)──► exploring ──(promote)──► promoted
                              │
                              └─(let cool)──► parked
```

- **`parked`** — captured. Not actively being shaped. Default state at creation.
- **`exploring`** — actively thinking; questions, evidence, alternatives, and a working plan may be filled in. Doctor flags this state if `updated:` is >90 days old (decide: promote, demote to parked, or delete).
- **`promoted`** — was explicitly handed to a request, roadmap Icebox, or existing shared destination. File should live in `archive/ideas/<slug>.md`, not `.spectacular/ideas/`. Doctor flags promoted entries still in the live folder as orphans.

## Shaping versus promotion

An idea is cheap capture. Before promotion, decide whether the destination is
accepted:

- unsettled product/experience/system direction → start a pre-request Vision and
  keep the idea as its origin link;
- accepted consequential behavior/contract → draft and approve an SPC, then
  create its request;
- accepted bounded execution destination → direct request promotion remains a
  compatibility route.

Vision approval does not itself promote or archive the idea. After the derived
SPC/request is accepted, record `promoted_to` and archive explicitly.

## Explicit handoff destinations

Choose the destination from uncertainty, visibility, cross-repo scope, commitment, and sensitivity:

| Destination | Use when | Authority after handoff |
|---|---|---|
| `request` | The outcome is accepted and durable execution coordination is warranted. | The request owns planning and execution. |
| `roadmap` | The human deliberately wants a project-priority signal but no version commitment yet. | The roadmap Icebox owns the local priority record. |
| `shared` | Collaborators need to shape the idea in a provider-neutral shared place. | The referenced Issue, Discussion, tracker, or document owns the discussion. |

Protected or sensitive material stays private/protected; never use ordinary shared publication for it. An Issue is a shared destination subtype, not a request or roadmap commitment.

### Request

When the user runs `spectacular idea promote <slug> --to request`:

1. CLI reads `.spectacular/ideas/<slug>.md` and extracts frontmatter + body content
2. Hands off to the [[new-request]] flow — request scaffolded with PLAN.md pre-filled from idea content
3. Sets `promoted_to: requests/<slug>/` on the idea file
4. Sets `status: promoted`
5. Moves file to `.spectacular/archive/ideas/<slug>.md`
6. Notes in the new PLAN.md: `promoted from ideas/<slug>.md`

### Roadmap

`--to roadmap` requires `--placement icebox`. The CLI snapshots the canonical roadmap before adding one Icebox item and archives the IDEA. It does not select a version, tier, scope, date, or request. Use the roadmap Icebox-promotion ritual for that separate commitment.

### Shared

`--to shared` requires an existing `--ref`. It records that stable reference in the archived IDEA and does not contact a network service. If a shared record does not yet exist, provide a copyable proposal to the human's chosen filing tool; wait for explicit external approval before anything is filed.

Handoff is **one-way and explicit**. A draft implementation plan alone is not a request handoff, and no destination is ever chosen automatically. `undo` remains available for the compatible request destination only.

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

None in v1.7.0. The word "idea" is short enough and the verb surface clear
enough that hidden routing isn't warranted. Feedback-loop retains
`iterate|test|probe|try`; ambiguous `experiment` routes through
[[discovery-protocol]].

## What this is **not**

- **Not a request.** Ideas don't carry PLAN/TASKS/lifecycle state. Their working plan is exploratory and may change with research or decisions; promotion produces a request only after an execution outcome is accepted.
- **Not memory.** MEMORY.md is for durable preferences and decisions. Ideas are *speculative* — they may never matter again.
- **Not a backlog.** ROADMAP.md's `## Icebox` section is the local project-priority surface. Ideas have no version, no scope commitment, and no exit criteria; prefer direct Icebox capture when that is already the intended home.
- **Not feedback.** Feedback entries answer "was that the right thing to ship?" — they're tied to something already built. Ideas are pre-commitment scratchpad for things that *might* get built.
- **Not Vision.** An idea captures a maybe. Vision is the structured, human-reactable process for choosing a direction before specification.

**Related:** [[new-request]] (promotion flow), [[archive]] (where promoted ideas land), [[doc-index]], [[scaffold-reference]], [[roadmap-rules]] (icebox vs ideas), [[memory-rules]] (durable vs speculative), [[feedback-rules]] (post-ship vs pre-commit).
