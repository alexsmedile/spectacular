---
title: Visual Conventions
description: "ASCII rendering conventions for Spectacular's visual layer — progress bars, summary dashboard, roadmap arc, and app-UI mockup blocks."
updated: 2026-06-07
---

# Visual Conventions

Spectacular's visual layer (v1.15.0+) gives read surfaces a scannable ASCII rendering. This page documents the conventions so the skill can produce consistent output and so you know what to expect in the terminal.

---

## Guiding principles

1. **Degrade cleanly.** Set `NO_COLOR=1` or pipe output to a file — every surface falls back to plain ASCII with no escape codes, no block characters, no color. The information is always there; color is an enhancement only.
2. **One shared helper.** All bar math goes through `_ascii_bar`; all box drawing through `_ascii_box`. No bespoke rendering per command.
3. **`--format json` is always byte-stable.** The visual layer touches only the default `text` format. JSON consumers are unaffected.
4. **Monochrome-correct.** Output reads correctly with color stripped — bars use `#`/`.` in plain mode, `█`/`░` with color.

---

## Progress bars

`spectacular progress <slug>` renders each milestone as a filled bar with a percentage and count:

```
  M1 — Contract                        ████████████████ 100% ✓
  M2 — vision/ soft-folder substrate   ████████░░░░░░░░  50%  3/6
  M3 — Ship                            ░░░░░░░░░░░░░░░░   0%  0/4

  overall                              ████████████░░░░░░░░  60%
```

**Plain-text mode** (`NO_COLOR=1` or non-TTY):

```
  M1 — Contract                        [################] 100% ✓
  M2 — vision/ soft-folder substrate   [########........]  50%  3/6
  M3 — Ship                            [................]   0%  0/4

  overall                              [############........]  60%
```

- Milestone bar width: 16 fill characters
- Overall roll-up bar width: 20 fill characters
- Completed milestones (done == total) show `✓` instead of the count
- Empty milestones (0 tasks) are shown as `0% ✓` (vacuously complete)

---

## Summary dashboard

`spectacular summary` renders request-state counts as proportional mini bars:

```
── Spectacular Workspace Summary ──

Project:    my-app

Requests:   6 total
            planned    ██████░░░░ 66%  4
            active     ███░░░░░░░ 33%  2

Decisions:  3
Memories:   5
Sessions:   1 (1 open)
```

Bar width is 10 fill characters; each status row is omitted when its count is 0.

---

## Roadmap arc

`spectacular roadmap` renders the version timeline grouped by tier:

```
── Roadmap ──

  ── Runway (near-term) ──
  ▶  v1.15.0       Visual layer + Imagine mode
  ·  v1.16.0       Contract prep ①: v2 contract spec

  ── Major (mid-term) ──
  ·  v1.17.0       Contract prep ②: v2 frontmatter fields

  ── Vision (future) ──
  ·  v2.0.0        The major: file-contract evolution
```

**Tier legend:**

| Tier in ROADMAP.md | Label | Meaning |
|---|---|---|
| `full` | Runway (near-term) | Detailed, near-term, planned |
| `themed` | Major (mid-term) | Themed, mid-term |
| `vision` | Vision (future) | Direction only, no schedule |

**Status indicators:** `✓` shipped · `▶` active · `·` planned

Shipped versions are hidden by default. Use `--all` to include them.

---

## App-UI mockup blocks

The skill can render an ASCII mockup of an application screen, dialog, or output surface directly in a request's `PLAN.md` or a capability spec under `specs/<capability>.md`. The format uses light box-drawing characters and a consistent layout convention so mockups are readable in any markdown viewer.

### Format

A mockup block is a fenced code block with language tag `mockup`:

````
```mockup
┌─ <Title> ───────────────────────────────────────────────┐
│                                                          │
│  <Content area — use plain ASCII layout>                 │
│                                                          │
│  [Primary action]          [Secondary]   [Cancel]        │
└──────────────────────────────────────────────────────────┘
```
````

**Rules:**
- Language tag `mockup` (not `ascii` or `text`) — lets renderers and the skill identify it.
- Width: ≤ 64 characters per line (fits most terminal widths and PR previews).
- Left-border box style matches `_ascii_box` output from the CLI.
- Buttons / actions shown in `[square brackets]`.
- Input fields shown as `[____________]` (underscores fill the field width).
- Labels left-aligned, values right-aligned or indented 2 spaces.
- Never use actual escape codes or ANSI color — this is stored in markdown.

### Real example — visual-layer workspace progress view

This mockup lives in `.spectacular/requests/visual-layer/vision/ui/dashboard.md` (the approved vision fragment). Reproduced here as the canonical format reference:

````
```mockup
┌─ spectacular summary ────────────────────────────────────┐
│                                                          │
│  Project:   my-app                                       │
│                                                          │
│  Requests   6 total                                      │
│    planned  ██████░░░░ 66%  4                            │
│    active   ███░░░░░░░ 33%  2                            │
│                                                          │
│  Decisions  3    Memories  5    Sessions  1 (open)       │
│                                                          │
│  Active requests:                                        │
│    ▶ visual-layer     v1.15.0  visual layer + imagine    │
│    ▶ cross-req-links  v1.13.0  advisory link awareness   │
│                                                          │
└──────────────────────────────────────────────────────────┘
```
````

### When to use

- In a request `PLAN.md § Validation` or `specs/<capability>.md` to sketch the expected output of a new CLI verb or surface
- As a vision fragment (`vision/ui/<name>.md`) during `spectacular imagine` — the render→react→derive loop uses these as artifacts for human approval
- Never in `CHANGELOG.md` or `README.md` — those are prose documents, not spec artifacts

---

## Skill rendering guidance

When the skill renders a visual mockup in conversation (e.g. as part of `imagine`), it should:

1. Open with a ` ```mockup ` fence.
2. Keep width ≤ 64 characters.
3. Use `█`/`░` for bar fills (terminal context) or `#`/`.` when the output is going into a markdown file.
4. After the mockup, ask the user to react (approve / redirect / reject) before continuing.

See `references/imagine.md` for the full render→react→derive loop.
