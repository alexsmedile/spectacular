# Handoff and Autopilot

Use this when: Orchestrator packaging delegation, Handoff context contracts, or Autopilot charters.

## Abstract Model Profiles

Spectacular is runtime-agnostic. Different host harnesses expose different model knobs (e.g. Antigravity subagent models, Claude Code model flags, Goose droids, OpenAI model tiers). To optimize both reasoning depth and token cost/latency, Missions, Objectives, and Handoffs declare abstract **Model Profiles**:

| Semantic Profile | Ideal Model Archetype | Spectacular Role | Typical Work |
|---|---|---|---|
| `reasoning` | Deep reasoning / thinking capability | **Orchestrator** | Genesis, Campaign planning, Claim design, Gap resolution, complex audits. |
| `fast-code` | Fast, high-throughput code fluency | **Worker / Runner** | Bounded Objective implementation, file edits, local test sweeps, refactoring. |
| `strict-verifier` | High instruction adherence, clean context | **Validator / Reviewer** | Adversarial verification, independent FROST review, regression suites. |

## Lean Autopilot Delegation Pattern (Subagents & Thread Linking)

For fast, zero-waste autonomous execution using host subagents (e.g. Antigravity/Claude `invoke_subagent`):

1. **Compact Dispatch Payload ($\le 300\text{--}500$ tokens)**:
   Provide the worker subagent strictly with:
   - **Mission Target**: Path to `.spectacular/missions/M<N>-<slug>/M<N>-<slug>.md`.
   - **Mechanical Scope**: Allowed paths (e.g. `cmd/`, `internal/pkg/`).
   - **Failable Verification**: The exact test command (e.g. `go test -v ./internal/pkg/...`).
   - **Fail-Fast Stop Triggers**: If an unrecorded architectural choice or external boundary conflict arises $\to$ stop immediately and yield back to Orchestrator.
   - **Thread Link**: Pass `conversation://<orchestrator-conversation-id>` for session continuity.

2. **Subagent Autopilot Execution**:
   - Subagent reads the target file and assigned code paths only (no workspace scans).
   - Writes the minimal coherent code satisfying the claim.
   - Executes the test command.
   - On pass $\to$ creates a clean Git commit and reports completion to the Orchestrator.
   - Zero sub-record files created.

3. **Sample Subagent Invocation Prompt**:
   ```text
   Task: Execute Mission M1 (SQLite Storage Engine) in autopilot.
   Mission File: .spectacular/missions/M1-sqlite-storage/M1-sqlite-storage.md
   Allowed Scope: internal/storage/, cmd/
   Verification Command: go test -v ./internal/storage/...
   Thread Context: conversation://2c67bc6b-f518-4aa3-951b-cddf1ba530b5
   Stop Triggers: Stop and yield back if schema migrations require external dependencies.
   ```

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

### Worker Execution Invariant (Charter Appendix)
Append this execution invariant to compiled Runner handoffs:
```text
INVARIANT: Implement the smallest coherent change that satisfies the assigned claims.
Check existing codebase utilities before authoring new abstractions.
Preserve validation, error handling, and security guards.
Touch only assigned mechanical scope; report unrelated defects without editing them.
```

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

- `.spectacular/missions/<slug>/handoffs/`: Governed task delegation records created by `spectacular handoff record`.
- `.spectacular/missions/<slug>/reviews/`: Canonical Reviews created by `spectacular review record`; never place an input `ReviewDraft` here.
- `scratch/` or a temporary directory: Ephemeral intake, review prompts, and `ReviewDraft` inputs before the CLI records them.

### Optional runtime pointer for threaded harnesses

For advanced operators delegating to host-managed threads or subagents (e.g. Claude Code `invoke_subagent`, Codex thread runs, OpenAI Agents SDK), the `HandoffDraft` can include an optional `runtime_pointer:` block:

```yaml
runtime_pointer:
  harness: claude-code # e.g. "claude-code", "codex", "openai-agents"
  thread_id: "agent-conv-903ea432" # host thread / conversation ID
  workspace_mode: "share" # "share" (git worktree), "branch" (isolated git branch), "inherit"
```

A host thread ID is ephemeral and owns nothing. The Handoff remains the immutable, governed artifact anchored to the repository commit and tree.

## Optional Thread and Subagent Execution (Advanced)

Single-session inline execution is the default. When an advanced operator opts into threaded subagent execution:

1. **Compile Context Sandwich**: Run `spectacular charter <mission-ref>/<objective-ref>` to obtain the exact 3-layer sandwich (Frozen Truth, Steering Sources, Perimeter Guardrails).
2. **Select Workspace Mode**: Determine workspace isolation before recording or dispatching:
   - `share` (Git Worktree): Best for parallel Runners operating on disjoint write reservations (`writes:`).
   - `branch` (Git Branch): Best for exploratory or destructive verifier passes.
   - `inherit` (Current Tree): Best for sequential bounded edits.
3. **Record Handoff First**: The Orchestrator records the delegation via `spectacular handoff record <mission> <draft> --by <sender>`, declaring write reservations and optional `runtime_pointer:` (`harness`, `thread_id`, `workspace_mode`). The governed Handoff must exist before the receiver begins acting.
4. **Dispatch Clean Sub-thread**: Dispatch the subagent runner in the selected workspace (e.g. Claude `invoke_subagent`, OpenAI thread run) with the compiled charter and recorded Handoff pointer as the fresh system/prompt context, avoiding parent history bleed.
5. **Intake & Re-verify**: On return, the Orchestrator re-verifies any `assumed` items, runs verification checks, and commits Evidence or Reviews.

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
