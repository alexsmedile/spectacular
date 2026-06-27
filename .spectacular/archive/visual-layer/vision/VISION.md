---
doc: vision
request: visual-layer
updated: 2026-06-02
---

# Vision — workspace-wide `spectacular progress` (a visual-layer slice)

> **Origin:** imagined as `progress-view` during the imagine-mode M6 dogfood (2026-06-02), then folded into [[visual-layer]] — it's a concrete spec for that milestone's `spectacular progress` rendering theme. The standalone `progress-view` request was retired; these artifacts live on as visual-layer planning input.

> Imagination-backed planning spine. Produced by `spectacular imagine`.
> The skill renders the slots below + ASCII fragments in `stories/` `ui/` `arch/`;
> the human reacts per-fragment (`approved:`); then a draft PLAN is derived.
> This file is the **spine** — narrative + a regenerable manifest. It never holds
> the fragments themselves.

## End goal

Opening the project, you type `spectacular progress` (no slug) and see the **whole workspace at a glance**: every active and planned request, each with its milestone tick-rate as a little bar, sorted so the in-flight work is on top and the stalled work is visible. Today `spectacular progress <slug>` answers "how far is *this* request?" — the end goal answers "how far is *everything*, and what's stuck?" without reading six TASKS.md files by hand. It's the cold-start glance a maintainer takes after a week away.

## Macro dev phases

- **Make it work** — `spectacular progress` with no slug aggregates every request's tick-rate into one table. Reuses the existing per-request parser; just loops.
- **Make it readable** — sort by status (active first), render a compact ASCII bar per request, dim the archived/verified, surface "stalled" (active but 0 ticks since N days).
- **Make it answer "what's stuck"** — a one-line summary header (`N active · M planned · K stalled`) and a `--stalled` filter, so the view doubles as a triage tool.

## Flow walk

1. Maintainer returns to a project after a week. Runs `spectacular progress`.
2. Sees a table: each request as a row — slug, status, a `[████░░░░] 50%` bar, last-touched date.
3. Active requests sort to the top; one is flagged `⚠ stalled` (active, no tick in 8 days).
4. They eyeball it, spot the stalled one, and run `spectacular request <slug>` to dig in.
5. Total time to "where do I pick up?": one command, ~3 seconds — vs. opening each TASKS.md.
## Manifest

- `[story]` Maintainer returns after a week — **true** (`stories/week-away-maintainer.md`)
- `[ui]` Workspace progress table (default view) — **true** (`ui/dashboard.md`)
- `[ui]` Stalled-only triage view (--stalled) — **true** (`ui/stalled-filter.md`)
- `[arch]` How the aggregate view is computed — **true** (`arch/data-flow.md`)
