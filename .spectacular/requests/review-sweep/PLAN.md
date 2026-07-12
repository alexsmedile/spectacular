---
status: verified
priority: high
owner: alex
updated: 2026-07-12
build: b31
summary: "Native review-sweep protocol: request-auditor agent (small/fast model) cross-checks review + ticked-active requests (claims vs code/tests/evidence), quick-checks planned requests for overlap with already-built work, appends sweep entries to VERIFY-LOG, hands off via SESSION.md; VERIFY-LOG evidence rows gain against:-stamps + pending-reverify semantics with doctor warnings."
related:
  - PRD.md
---

# Plan — review-sweep

## Goal

Make the manually-proven "small fast model audits the request fleet" workflow native: a repeatable, read-only sweep that cross-checks each request's claims against its code/tests/evidence, records findings durably, and hands off next-agent instructions — without ever mutating lifecycle itself.

## Constraints

- The sweep is read-only on the workspace: agents return findings; **the orchestrator is the only mutator** (fleet contract) — it appends VERIFY-LOG entries and updates SESSION.md, human confirms lifecycle proposals.
- The sweep never promotes. It *feeds* the verify walk; `review → verified` still goes only through the walk.
- No new STATUS.md file — per-request handoff lives in SESSION.md § Next actions; workspace-level in `spectacular status` / `next` (rejected: duplicate source of truth).
- No mechanical git-staleness heuristics in the CLI — staleness is stamped data + agent/human judgment; doctor only validates the stamps' presence and flags `pending-reverify` rows.
- Auditor agent follows the existing fleet grid conventions (`agents/*.md` source of truth, relative symlink in `.claude/agents/`, closed contract, bounces on judgment outside its brief).
- VERIFY-LOG stays append-only; sweep entries interleave with walk entries in the same file.

## Understanding

### How it works now

Requests reach `review` when TASKS are ticked; the interactive verify walk (verify.md) is the only defined protocol from there, and it's human-in-the-loop per check. `code-reviewer` reviews a *diff*, `spec-reviewer` a *doc*, `test-verifier` a *named check* — no agent's brief is "cross-check a request's claimed state vs its actual evidence". VERIFY-LOG `[manual]`/`[observe]` rows record evidence + date but carry no build/identity stamp and no staleness state, so old evidence is indistinguishable from current. Planned requests are never checked against what already shipped, so a stale planned request can duplicate built work.

### What changes

verify.md's VERIFY-LOG shape gains `against: <commit/build> · <identity>` stamps on `[manual]`/`[observe]` rows plus a `pending-reverify` row marker; doctor (lifecycle area) warns on missing stamps or `pending-reverify` rows in `review`/`verified` requests. A new `references/review-sweep.md` defines the three-tier sweep loop; a new read-only `request-auditor` agent (haiku-class) carries the per-request audit brief; SKILL.md gains the `spectacular sweep` trigger row; the CLI gains a `sweep` redirect stub (like `verify`).

### What stays the same

The verify walk and its gate; the lifecycle states and `advance` semantics; all existing agents' briefs; VERIFY-LOG's append-only walk entries; SESSION.md's template and update rhythm; the closure gate at archive.

## Decisions

- Chose stamp + judgment flip over mechanical `git log <stamp>..HEAD` staleness detection — because path tracking per request is a new schema surface and git heuristics misbehave on rebases/monorepos; doctor validates stamps mechanically, judgment flips `pending-reverify`.
- Chose VERIFY-LOG as the sweep-findings home over a separate SWEEP-LOG.md — one audit trail per request, walks and sweeps interleaved and dated; one fewer file type.
- Chose `spectacular sweep` as a skill-side trigger with a CLI redirect stub over a CLI-native verb — the sweep is judgment work (agent dispatch), same shape as `verify`.
- Chose coverage = `review` (full audit) + ticked-`active` (audit + may PROPOSE `advance --to review`) + `planned` (quick overlap check vs specs/index.md capabilities + archive slugs — flags work already built, added mid-grill by user) over review-only — catches both the "Implemented but unproven" state and the "about to re-build shipped work" state.
- Chose one auditor dispatch per request (parallel fan-out) over one agent walking the whole fleet — bounded context per agent, matches the fleet's closed-brief contract.
- Chose subagent delegation over inline-with-lower-effort (user raised both 2026-07-12) — the audit's cost is *reading* (PLAN/TASKS/VERIFY-LOG/code/tests), which inline would carry in the orchestrator's context permanently; a subagent reads off-context and returns only findings, and the small-model pin lives on the agent def. Orchestration (dispatch, collect, persist, relay) stays with the main agent.
- Refinement: the planned-tier overlap check is one *batched* light dispatch (all planned requests in a single auditor call — it only compares summaries vs specs/index.md + archive), not per-request fan-out; only review/ticked-active requests get individual deep auditors.
- Rejected STATUS.md — duplicates SESSION.md + `spectacular next`.

## Milestones

- M1 — Evidence freshness schema live: VERIFY-LOG shape documents `against:` stamps + `pending-reverify`; doctor warns on missing stamp / pending-reverify rows in review+verified requests
- M2 — Sweep protocol + trigger: `references/review-sweep.md` (three-tier loop, findings→VERIFY-LOG, handoff→SESSION.md), SKILL.md routing row, CLI `sweep` redirect stub
- M3 — `request-auditor` agent def in the fleet grid (read-only, small-model, closed audit brief) + fleet docs updated

## Tasks

See `TASKS.md`.

## Dependencies

- None (builds on b29's directive gates and b30's walk-only verify.md; VERIFY-LOG schema edit touches verify.md § VERIFY-LOG shape only)

## Validation

- M1 — run: `tests/cli/doctor.test.sh` passes with new scenario: seeded review request whose VERIFY-LOG has an unstamped `[manual]` row → warning; a `pending-reverify` row → warning; clean stamped log → no warning. Assertable: verify.md § VERIFY-LOG shape shows `against:` + `pending-reverify` markers.
- M2 — assertable: `./cli/spectacular sweep` prints the skill-redirect text (exit 1 — matches the `verify` skill-verb convention; superseded from "exit 0" 2026-07-12); SKILL.md routes `sweep` → review-sweep.md; review-sweep.md covers the three tiers (review / ticked-active / planned-overlap) and states the never-promotes + orchestrator-mutates rules; `doctor links` clean.
- M3 — assertable: `agents/request-auditor.md` exists with read-only tools (Read/Grep/Glob/Bash), `.claude/agents/request-auditor.md` is a relative symlink to it; CLAUDE.md fleet list includes it; the agent brief forbids edits/lifecycle moves and requires the findings-block return shape defined in review-sweep.md.

## Deliverables

- `skills/spectacular/references/review-sweep.md` — the three-tier sweep protocol
- `agents/request-auditor.md` + `.claude/agents/` symlink — the small-model auditor
- verify.md § VERIFY-LOG shape extension (`against:` stamps, `pending-reverify`)
- `cli/spectacular` — `sweep` redirect stub + doctor lifecycle stamp checks
- SKILL.md routing row, docs/commands.md sweep section, roadmap ledger row b31 (v1.35.0), CHANGELOG [Unreleased] entry
- Test scenarios in `tests/cli/doctor.test.sh` + a `sweep` redirect assert
