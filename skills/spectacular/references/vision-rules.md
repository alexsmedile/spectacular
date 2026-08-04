---
doc-id: vision
mode: index
location: .spectacular/visions/<slug>/
entries-dir: .spectacular/visions/<slug>/fragments/
spine: .spectacular/visions/<slug>/VISION.md
scope: pre-request
template: templates/vision/
snapshot-on-edit: false
summary: "Human-approved pre-request direction — grounded understanding, proportional fragments, linked evidence, and an agentic handoff to a draft SPC."
status: active
---

# Vision Rules

Vision is an optional pre-request direction workspace. `imagine` is the
generative operation that builds it; `VISION.md` is the durable result. Vision
closes divergent exploration into a human-approved north star before a
specification or execution plan exists.

## Structure

```text
.spectacular/visions/<slug>/
├── VISION.md
├── fragments/
│   └── <name>.md
└── evidence/
    └── <presentation artifact or link note>
```

New fragments are flat because `kind:` classifies them and filenames remain
unique within one Vision. Legacy request-owned `vision/{stories,ui,arch}/`
folders remain readable and diagnosable; new writes use the top-level collection.

## Lifecycle

```text
draft ──► proposed ──► approved ──► draft SPC
  └──────────┴──────► rejected
```

| State | Meaning | Gate |
|---|---|---|
| `draft` | Understanding and alternatives are still changing | None |
| `proposed` | All reactions resolved; coherent direction ready for vibe check | Required spine sections filled, ≥1 approved fragment |
| `approved` | Human accepted the whole direction | Explicit actor/date; proposed state |
| `rejected` | Human declined the direction | Explicit actor/reason from draft or proposed |

Approval authorizes specification derivation only. It never activates a request
or authorizes implementation.

## Spine contract

Required H2 sections, in order:

1. Intent
2. North star
3. Understanding
4. Experience signature
5. Strategies considered
6. Chosen direction
7. Boundaries
8. Evidence
9. Manifest
10. Approval

Understanding contains Current reality, Users and needs, Constraints, and
Material uncertainties. Manifest is regenerated. Approval records the whole
Vision decision; fragment reactions do not substitute for it.

Frontmatter:

```yaml
---
doc: vision
status: draft | proposed | approved | rejected
slug: <slug>
approved_by: ""
approved_at: ""
derived_to: ""
updated: YYYY-MM-DD
---
```

## Fragment contract

Kinds: `strategy | story | flow | ui | arch | prototype`. No kind is mandatory.

```yaml
---
kind: <kind>
caption: <one-line reaction target>
reaction: pending | approved | revise | rejected | superseded
reaction_note: ""
updated: YYYY-MM-DD
related: []
---
```

Every fragment states its proposal and decision impact. Only approved fragments
may become SPC requirements. Evidence is linked, not approved: RES/SPK truth uses
its own result vocabulary.

## Commands

```text
spectacular imagine <slug>
spectacular vision list
spectacular vision show <slug>
spectacular vision add <kind> <name> --slug <slug> [--caption <text>]
spectacular vision react <slug> <fragment> --approve|--revise|--reject|--supersede
spectacular vision propose <slug>
spectacular vision approve <slug> --approved-by <human>
spectacular vision reject <slug> --rejected-by <human> --reason <why>
/spectacular vision derive <slug>
```

Derive is agentic because translating direction into a precise contract requires
judgment. It requires approved Vision and produces a draft SPC, never PLAN.

## Doctor area

`doctor vision` validates new and legacy locations:

- spine presence and canonical section headings;
- `draft|proposed|approved|rejected` state;
- fragment kind/caption/reaction shape;
- proposed/approved reaction coherence;
- approval or rejection provenance;
- manifest agreement with fragment files;
- legacy kind/subfolder agreement.

`--fix` only regenerates the Manifest. Content, reactions, and approval are
judgment and never auto-set.

## Concept boundaries

- Idea captures a maybe; Vision shapes an unsettled direction.
- Research establishes facts; Spike establishes feasibility; neither chooses.
- Prototype is a showable fragment/artifact.
- Feedback-loop acquires learning from built/prototyped behavior after there is a
  target; pre-spec direction approval belongs here.
- SPC is the convergent implementation contract.
- PLAN is request execution decomposition and retains technical Understanding.

**Related:** [[imagine]], [[idea-rules]], [[discovery-protocol]],
[[lifecycle-contract]], [[spec-lifecycle]], and [[plan-rules]].
