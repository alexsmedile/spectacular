# Handoff and Autopilot

Use this when: Orchestrator packaging delegation, Handoff context contracts, or Autopilot charters.

## Abstract Model Profiles

Spectacular is runtime-agnostic. Different host harnesses expose different model knobs (e.g. Antigravity subagent models, Claude Code model flags, Goose droids, OpenAI model tiers). To optimize both reasoning depth and token cost/latency, Missions, Objectives, and Handoffs declare abstract **Model Profiles**:

| Semantic Profile | Ideal Model Archetype | Spectacular Role | Typical Work |
|---|---|---|---|
| `reasoning` | Deep reasoning / Thinking (Claude Sonnet w/ Thinking, Gemini Pro, o1/o3) | **Orchestrator** | Genesis, Campaign planning, Claim design, Gap resolution, complex audits. |
| `fast-code` | Fast, high-throughput code fluency (Gemini Flash, Claude Haiku, GPT-4o-mini) | **Worker / Runner** | Bounded Objective implementation, file edits, local test sweeps, refactoring. |
| `strict-verifier` | High instruction adherence, clean context | **Validator / Reviewer** | Adversarial verification, independent FROST review, regression suites. |

When dispatching subagents or delegating tasks:
- The Skill maps `fast-code` to faster, lightweight runner models (e.g. `Model: flash` in Antigravity).
- The Skill maps `reasoning` or `strict-verifier` to deep reasoning models (e.g. `Model: pro` or isolated reviewer contexts).

## Autopilot is explicit and non-default

Never assume it. When the owner turns it on, bind the charter to:

- the exact Mission activation fingerprint
- Objective and claim scope
- Contract and Git baseline
- allowed operator actions
- effects that still require the owner
- budgets and checks
- expiry, stops, recovery
- the return destination

State how resources are actually enforced, as one of:

| Level | Means |
|---|---|
| `hard` | independently verified measurement, and real cancellation |
| `observed` | measured and reported, but not enforced |
| `unsupported` | not measured at all |

Only claim `hard` when both the measurement and the cancellation are verified.

## Promote before delegating

```bash
spectacular objective promote <mission-ref>/<objective-ref> --json   # e.g. M7/O2
```

Promote an inline Objective to its own file before independent delegation. It
lands at `.spectacular/missions/<slug>/objectives/<ref>-<slug>.md` and keeps its
identity. The file then carries the exact:

- outcome and claims
- dependencies and inputs
- semantic and mechanical scope
- authority and stops
- return contract

Accountability stays with the Mission owner. A host task or thread is only a
destination pointer — it owns nothing.

### Add a Runner context contract

Every independent Runner Handoff carries this compact section in its Markdown
body. It is guidance, not a new schema field: use judgment to keep the read set
small and exact.

```md
## Runner context contract

Read:
- `M15/O2`
- `M15/R1`
- `internal/auth/...`
- `STACK.md`

Do not load:
- Campaigns
- other Missions
- archive
- workspace catalog

If blocked:
- Ask the Orchestrator for one named authoritative source.
```

A Runner follows this contract instead of scanning the workspace. A Campaign's
current block is roadmap context for an Orchestrator, never an assignment to a
Runner. If a Mission body explicitly cites a Campaign, read only the cited
context and only when it is relevant to the assigned work.

## Record the delegation as a Handoff

```bash
spectacular handoff record <mission-ref> <handoff.md|-> --by <sender> --json
```

Run `spectacular handoff record --help` to output the exact `HandoffDraft` YAML frontmatter template.

The Handoff lands in the Mission bundle and binds the exact commit and tree it
was sent against, verified against the repository (if `tree` is omitted from the draft,
it is auto-derived from the commit). A delegation that lives only
in a chat message or a temp file leaves no record of what was asked or what state
it was asked against.

Separate what you checked from what you are carrying over:

| Field | Means |
|---|---|
| `asserted` | the sender verified this |
| `assumed` | the sender is taking this on trust |

Both are required; an empty list is a legal statement, an absent one is not.
**Neither is ever scored** — nothing verifies that an `asserted` item was really
checked. The split records a claim its sender signs. **The receiver re-verifies
everything under `assumed` before acting on it.**

A recorded Handoff is frozen. Correct it by recording a new one carrying
`supersedes:`; the original survives as what its sender believed at the time, and
`mission show` points a reader of the superseded record forward to the one that
is current. Never edit a Handoff in place.

The receiving agent inspects incoming Handoffs via `spectacular mission show <ref> --json`
or directly in `.spectacular/missions/<mission>/handoffs/`.

### Handoff and review directory architecture

To keep multi-agent artifacts clean and unambiguous:

- `.spectacular/missions/<slug>/handoffs/`: Governed task delegation records (`spectacular handoff record`) AND review handoff prompts (`review-handoff-prompt.md`) for external or subagent reviewers.
- `.spectacular/missions/<slug>/reviews/`: Formal review output (`ReviewDraft` records recorded via `spectacular review record`).
- `scratch/` or project root: Ephemeral scratchpad files (e.g. 1-shot intake PRDs before Genesis digestion).

## Fan out sparingly

Delegate only cohesive mid-to-long work whose claim ownership is disjoint.

Avoid:

- tiny sessions
- recursive critic loops
- repeated full reviews

Finish working code, run focused checks, and batch compatible review at the
Mission's frozen review level.

## What the receiver returns

- status and actor
- final baseline and result
- changed files
- checks that ran
- native-provider receipts
- Evidence
- remaining Gaps
- budget use
- recovery point
- one next action, or one owner gate

**The receiver never** changes Mission criteria, declares Evidence sufficient, or
gains provider permission it did not already have.
