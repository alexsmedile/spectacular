---
id: SPC-002
type: specification
status: implemented
target_version: "v1.36.0-execution"
supersedes: ""
updated: 2026-08-02
summary: "Request creation, implementation context, command grammar, and verified execution handoff"
related:
  - SPC-001-wayfinding.md
  - ../archive/request-workflow-interface/PLAN.md
version: 1.1
approved_at: 2026-08-02
approved_by: user
implemented_at: 2026-08-02
verified_against: commit bf0756f
---

# SPC-002 — Request creation, implementation context, command grammar, and verified execution handoff

## Intent

Give humans and coding agents one predictable path from approved specification to bounded implementation context, native session planning, evidence-backed verification, and archival without duplicating project history.

## Requirements

- Keep persistent request milestones in `TASKS.md`, current-session micro-planning in native Codex/Claude planning tools, and narrower closed briefs for subagents.
- Use noun-first canonical CLI entity commands and verb-first conversational document commands while preserving existing aliases.
- Make `spectacular request <slug>` the overview, `--brief` the self-retrieved implementation starting prompt, and `--full` the ordered request-owned Markdown bundle.
- Permit `--brief` only for active requests. Planned requests point to activation, review requests point to verification, and verified requests point to closure.
- Support milestone selection as `--milestone M2`, `--milestone 2`, `-m 2`, and `-m2`; defer generic `--artifact` and `--section` selectors.
- Keep human output concise and expose stable equivalent JSON output with explicit errors and nonzero exit status for blocked reads or mutations.
- Make `spectacular request new <slug> --from SPC-NNN` require an approved specification and generate both PLAN and milestone-level TASKS without activating implementation.
- Permit initial TASKS generation from an approved specification, but forbid silent scope addition, removal, or reordering after activation.
- Make `/spectacular act SPC-NNN` the canonical short agentic execution command; accept `/spectacular spec act SPC-NNN` as an explicit alias and `/spectacular SPC-NNN` only when the target is unambiguous.
- Have the agentic act flow find or create exactly one request, resolve blockers, review execution readiness, record activation provenance, move to active, retrieve the brief, inspect current code/Git, initialize native planning, and begin implementation.
- Record flat activation fields for actor, timestamp, source-spec version/digest, and Git baseline; do not copy the specification body into the request.
- Make `/spectacular verify` the normal authority for `review → verified`; exceptional direct transitions require an explicit reason recorded as an override.
- Fold documentation-impact assessment into verification and closure guidance while retaining the low-level compatibility command.
- Keep `spectacular afk cleanup` as AFK branch hygiene; do not add another namespace level until another cleanup category exists.
- Use collision-safe aliases: `fnd1` for `FND-001`, `f1` for fixes, `bug1` for `BUG-001`, and `b1` for roadmap builds.
- Compact decisions after 50 entries: newest 50 individual, preceding 50 in blocks of 10, older history in blocks of 50; preserve strong frontmatter and index summaries for selective retrieval.
- Preserve the five-state request lifecycle, code/tests as implementation truth, archive-first behavior, Bash 3.2 support, and all existing HITL and execution-scope boundaries.

## Evidence and decisions

- `DEC-010` — durable milestones, native session plans, and narrow subagent briefs.
- `DEC-011` — noun-first CLI and verb-first conversational commands with aliases.
- `DEC-012` — generated implementation brief, never stored BRIEF.md.
- `DEC-013` — mechanical spec scaffolding separated from agentic act authorization.
- `DEC-014` — verification evidence owns the verified transition; docs impact moves into closure guidance.
- `DEC-015` — collision-safe FND/FIX/BUG/build aliases.
- `DEC-016` — rolling decision-index compaction tiers.
- Explicit user confirmation and implementation authorization recorded in the 2026-08-02 request-workflow interview.

## Confirmation

Confirmed by the user on 2026-08-02. Approve this specification before activating the linked request.

**Approved 2026-08-02 by user** — Explicit user confirmation of the complete request-workflow contract on 2026-08-02
