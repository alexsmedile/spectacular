---
version: 1.1
updated: 2026-05-31
summary: "Operating policies — the practice layer paired with PRINCIPLES.md"
---

# Spectacular — Operating Policies

<!--
  POLICY.md is the PRACTICE layer (PRINCIPLES.md is the THEORY layer).
  Policies are filed under named work-phase hooks (## @<hook>). The skill
  retrieves only the active hook's policies on entering a phase.

  Anatomy:  ### <verb>-<noun>
            - principle: N      (optional — the PRINCIPLES.md § it enforces)
            - severity: block   (block = refuse | warn = surface + continue)
            - check: <condition>
            <prose: rationale + the instruction injected into context>

  Severity is OPT-IN to blocking: blocks ONLY if it explicitly says
  'severity: block'. Absent/warn/unrecognized → non-blocking.

  Valid hooks (8): @Init @Planning @Implementation @Verification @Archive
  @Remember @Snapshot @SessionEnd.
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
A new request's PLAN must be well-shaped before work begins: a one-sentence Goal, explicit Constraints, and demoable Milestones — not a vague wish. Surface a thin plan; let the author proceed.

### scope-down
- principle: 10
- severity: warn
Before fixing milestones, name the smallest high-impact slice that delivers the core value now, and push the rest to ROADMAP as `v2+`. Prefer a finished MVP of the features actually needed today over a complete build of features that might be. Flag speculative generality and any feature without a current need. Surface the leaner cut; the human chooses the scope.

## @Implementation

### understand-before-change
- principle: 7
- severity: block
- check: PLAN.md has a filled `## Understanding` section (How it works now / What changes / What stays the same), OR a `UNDERSTANDING.md` exists with the same three subheads

A request must not move `planned → active` until the agent has written down how the system works today, what this change touches, and what it leaves alone. Establish understanding before touching code. Satisfied by either the PLAN slot or a dedicated UNDERSTANDING.md.

## @Verification

### verification-present
- principle: 7
- severity: block
- check: every check in VERIFY.md (or PLAN § Validation) is satisfied before `review → verified`

A request must not reach `verified` while any verification check is unmet. Verification always happens; the only question is which artifact carries the checks. (Absorbs verify-walk's gate.)

## @Archive

### spec-sync
- principle: 2
- severity: warn
On archiving a request, propose the SPEC.md / `specs/` updates the shipped work implies. Intent and truth are different files — keep truth current. The human confirms.

### memory-propose
- principle: 5
- severity: warn
On archiving, propose any operational lesson worth keeping as a memory. Operational memory compounds. Surface the candidate; never write memory without confirmation.

## @Remember

### confirm-before-write
- principle: 8
- severity: block
- check: the user has confirmed the memory text before it is written to `.spectacular/memory/`

Memory is team-visible and git-committed. Humans decide, agents propose: never write a memory the user has not seen and confirmed.

## @Snapshot

### snapshot-before-overwrite
- principle: 8
- severity: block
- check: a `<DOC>@v<N>.md` snapshot exists before a canonical doc is overwritten in place

Canonical documents are never overwritten without a snapshot first. The unversioned filename always points to current; history is preserved.

## @SessionEnd

### summarize-before-handoff
- severity: warn
Before handing off, summarize what changed, what's left, and what's next, so the next session resumes without re-deriving context. Surface the summary; don't block.
