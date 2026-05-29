---
status: planned
priority: high
owner: alex
updated: 2026-05-29
summary: "ASCII rendering layer — progress bars, roadmap render, summary dashboard, app-UI mockup blocks; shared ascii-render helper, --format/NO_COLOR aware"
related:
  - PRD.md
  - ../../ROADMAP.md
  - ../../ARCHITECTURE.md
target_version: v1.14.0
---

# Plan — visual-layer

## 1. Goal

Give Spectacular's read surfaces a scannable visual layer — progress bars, a roadmap render, a summary dashboard, and ASCII mockup blocks for app-UI requests — so workspace state is understood at a glance instead of parsed from flat `M1 — …: 0/5` text. This is a rendering layer over data the CLI already computes, not a new subsystem.

## 2. Constraints

- **Convention alignment, not new flags.** Follow the existing `--format text|json` family (visual = default `text`; `json` stays plain-machine-readable). Don't invent a parallel `--ascii`/`--render` flag.
- **Degrade cleanly.** Honor `NO_COLOR` and non-TTY (piped) output — fall back to ASCII-only, no color, no escape codes — so piping into files/tools stays clean. Existing JSON consumers must be byte-unaffected.
- **One shared helper.** A single `ascii-render` helper (bar math, box-drawing, width clamping) is reused by every surface — no per-command bespoke rendering.
- **Read-only, static output.** This is rendered output, not a live TUI. No cursor control, no interactivity, no image/SVG export (that's pageworks/renderer territory).
- **Monochrome-correct.** Color is an enhancement, never the only signal — output must read correctly with color stripped (accessibility).

## 3. Milestones

- M1 — `ascii-render` helper: bar rendering (`███░░ 60%`), box-drawing, width clamping; NO_COLOR / non-TTY detection; unit-tested in isolation.
- M2 — Data-backed renders: `spectacular progress <slug>` → milestone bars + roll-up %; `spectacular summary` → dashboard with request-state bars + substrate counts. `--format json` unchanged.
- M3 — Roadmap render: `spectacular roadmap` → version arc / timeline view (runway → major → vision), tier-aware.
- M4 — App-UI mockup blocks: a documented ASCII layout-block format the skill can drop into a request's PLAN/SPEC to sketch a UI (rendered, not just described); pairs with the AskUserQuestion preview pattern; 1+ real example.
- M5 — Docs + ship: `docs/` page for the visual conventions; CHANGELOG entry; plugin bump to v1.14.0.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- Builds on shipped read verbs (`progress`, `summary`) + the milestone parser + the structured ROADMAP — all already in place.
- Enables the deferred "visual link-graph render" in [[cross-request-links]] (v2 task there).
- Independent of [[verify-walk]] and [[cli-debt-removal]]; orderable freely on the runway.

## 6. Validation

- M1 — Helper unit tests cover bar math, width clamping, and NO_COLOR stripping; no escape codes leak when piped.
- M2 — `progress` and `summary` render bars on a TTY; `--format json` output is byte-identical to today; piped text is plain.
- M3 — `spectacular roadmap` renders the version arc with correct tiers; degrades to plain text non-TTY.
- M4 — A real request's PLAN/SPEC carries a rendered app-UI mockup block; the format is documented.
- M5 — `docs/` visual-conventions page ships; manifests at v1.14.0.

## 7. Deliverables

- `ascii-render` shared helper (in `cli/spectacular`) + unit tests
- Visual `progress`, `summary`, and `roadmap` renders (text default; json untouched)
- Documented app-UI mockup block format + 1 example in a real request
- `docs/<visual-conventions>.md` page (registered in `docs.yaml`)
- CHANGELOG [1.14.0] entry
