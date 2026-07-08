---
version: 1.0
updated: <DATE>
summary: "Operating policies — the practice layer paired with PRINCIPLES.md"
---

# <NAME> — Operating Policies

<!-- Pairs with PRINCIPLES.md (theory). PRINCIPLES is optional, so it is NOT a
     declared 'related:' dependency — policies link individual principles via the
     per-policy 'principle: N' tag instead. -->

<!--
  POLICY.md is the PRACTICE layer (PRINCIPLES.md is the THEORY layer).
  Policies are filed under named work-phase hooks (## @<hook>).
  The skill retrieves only the active hook's policies on entering a phase.

  Anatomy of a policy block:
    ### <verb>-<noun>
    - principle: N        (optional — the PRINCIPLES.md § it enforces)
    - severity: block     (block = refuse to proceed | warn = surface + continue)
    - check: <condition>
    <prose: rationale + the instruction the skill follows when injected>

  Severity is OPT-IN to blocking: a policy blocks ONLY if it explicitly
  says `severity: block`. Absent/warn/unrecognized → non-blocking.

  Hooks (the only valid 9): @Init @Planning @Implementation @Verification
  @Archive @Debugging @Remember @Snapshot @SessionEnd.

  Full spec: skill references/policies-contract.md
-->

## @Init

### scaffold-contract
- principle: 4
- severity: warn
The workspace must satisfy its scaffold contract — README present, `.spectacular/` committed, `.spectacular.local/` gitignored, always-set docs in place. Surface any gap on init; don't block.

## @Planning

### request-shape
- principle: 3
- severity: warn
A new request's PLAN must be well-shaped before work begins: a one-sentence Goal, explicit Constraints, and demoable Milestones — not a vague wish. Surface a thin or unfocused plan; let the author proceed.

### scope-down
- principle: 10
- severity: warn
Before fixing milestones, name the smallest high-impact slice that delivers the core value now, and push the rest to ROADMAP as `v2+`. Prefer a finished MVP of the features actually needed today over a complete build of features that might be. Flag speculative generality and any feature without a current need. Surface the leaner cut; the human chooses the scope.

## @Implementation

### understand-before-change
- principle: 7
- severity: block
- check: PLAN.md has a filled `## Understanding` section (How it works now / What changes / What stays the same), OR a `UNDERSTANDING.md` exists with the same three subheads

A request must not move `planned → active` until the agent has written down how the system works today, what this change touches, and what it leaves alone. Establish understanding before touching code — half-understood changes are how regressions enter. Satisfied by either the PLAN slot or a dedicated UNDERSTANDING.md.

**Law: no flip to `active` without written understanding.**

| Excuse | Reality |
|---|---|
| "The change is small enough to skip this" | Small changes in half-understood systems are how regressions enter. Smallness is not understanding. |
| "I'll fill Understanding once I've started" | After you start, it becomes a justification of what you already did, not an analysis. |

Red flags — stop: editing code while `## Understanding` still holds placeholders; describing current behavior with "probably" or "should".

## @Verification

### verification-present
- principle: 7
- severity: block
- check: every check in VERIFY.md (or PLAN § Validation) is satisfied before `review → verified`

A request must not reach `verified` while any verification check is unmet. Verification always happens; the only question is which artifact carries the checks. If VERIFY.md exists it is load-bearing; otherwise PLAN § Validation blocks the transition. (Absorbs verify-walk's gate.)

**Law: no `verified` without every check satisfied — on fresh evidence.**

| Excuse | Reality |
|---|---|
| "The tests passed earlier" | Earlier is not now. Run the named check fresh; only its result decides. |
| "This check is basically covered by that one" | "Basically" means unverified. Each check runs as written. |

Red flags — stop: claiming a pass from memory; "should pass"; skipping a manual check because it's tedious.

## @Archive

### spec-sync
- principle: 2
- severity: warn
On archiving a request, propose the SPEC.md / `specs/` updates the shipped work implies. Intent and truth are different files — keep truth current. Surface the proposed sync; the human confirms.

### memory-propose
- principle: 5
- severity: warn
On archiving, propose any operational lesson worth keeping as a memory. Operational memory compounds across sessions. Surface the candidate; never write memory without confirmation.

## @Remember

### confirm-before-write
- principle: 8
- severity: block
- check: the user has confirmed the memory text before it is written to `.spectacular/memories/`

Memory is team-visible and git-committed. Humans decide, agents propose: never write a memory the user has not seen and confirmed. Show the proposed entry, then write on confirmation.

**Law: no memory write the user hasn't seen.**

| Excuse | Reality |
|---|---|
| "It's obviously worth remembering" | Obvious to you is not confirmed by them — memory is team-visible and permanent. |
| "I'll show it right after writing" | After is not confirmation; it's a fait accompli. |

Red flags — stop: writing to `memory/` in the same turn the lesson emerged, with no draft shown.

## @Snapshot

### snapshot-before-overwrite
- principle: 8
- severity: block
- check: a `<DOC>@v<N>.md` snapshot exists before a canonical doc is overwritten in place

Canonical documents are never overwritten without a snapshot first. The unversioned filename always points to current; history is preserved. Snapshot, then overwrite.

**Law: no canonical overwrite without a snapshot.**

| Excuse | Reality |
|---|---|
| "Tiny edit — git has the history" | Git history is not the workspace convention; `<DOC>@vN.md` snapshots are what agents and humans actually read. |
| "I'll snapshot if the edit grows" | By then the overwrite already happened. |

Red flags — stop: editing PRD/SPEC/any canonical in place with no `@vN` snapshot created first.

## @SessionEnd

### summarize-before-handoff
- severity: warn
Before handing off, summarize what changed, what's left, and what's next, so the next session (human or agent) resumes without re-deriving context. Surface the summary; don't block the handoff.
