---
type: scratch
created: 2026-05-23
expires: when Tier 1 cleared and Tier 2 promoted to requests
summary: "Working priority stack after v0.6.0 ship + TODO/DECISIONS pruning"
---

# Priority Tiers — 2026-05-23

Snapshot of open work after pruning TODO.md and rejecting request-numbering (DECISIONS 2026-05-23). Delete this file once Tier 1 is verified and Tier 2 items become real requests.

## Tier 1 — Live validation (no design, just do)

Requests already in `review` status; code shipped, awaiting interactive scenarios:

- **convention-pack-fabricator** — live grill scenarios for `pack-overrides.md` + `alex-default` dogfood
- **convention-pack-application** — live three-mode (suggest / scaffold / enforce) scenarios for CLI `pack` subcommand + init/doctor wiring
- **doctor** — interactive skill-side scenarios for substrate self-check (shipped v0.3.1)

Run these → mark `verified` → archive. Cheapest path to a clean board.

## Tier 2 — Real design open (need a request first)

- **Phase-end auto-continue** (TODO L6) — biggest UX gap. Skill should drive task to done without stopping at every transition. Design needed: which transitions auto-advance, which gate, how to surface autonomous decisions.
- **Catch new plans/ideas → suggest wiring** (TODO L4) — detect "this sounds like a new request" mid-conversation, offer `spectacular new`. Small + high leverage.
- **Interview mode / `grill-me`** (TODO L7) — inverse of `grill`: skill interviews you to build a doc from scratch. Overlaps PRD kit prompts but more conversational.

## Tier 3 — Gated, don't touch

- **public-docs-advanced** — explicit 2-of-3 trigger gate (none fired)
- **convention-pack-modules** — gated on composition pain from v1 use
- **cli-bootstrap** — parked v0.2.x maintenance

## Tier 4 — Deferred / external

- Subagents (2026-05-11 DECISION defers to v2)
- Octopus tasks/sessions integration (need external context)

## Recommendation

Clear Tier 1 first (afternoon's work, 3 requests → verified), then pick **phase-end auto-continue** as next real build. Most-mentioned friction; reshapes how the skill feels.
