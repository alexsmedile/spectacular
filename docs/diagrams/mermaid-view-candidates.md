# Mermaid Views

Mermaid renderings for Spectacular Mission and Objective views. ASCII is the default
surface and carries the straightforward case; mermaid is for graphs ASCII cannot render
legibly.

**M2 and M4 are approved.** M1 and M3 are declined, because the ASCII Objective graph
and ASCII timeline read better for their cases. All four are kept here so the decision
is judged from rendered output rather than description.

The approved ASCII views (Mission timeline, compact chain, Objective graph, Objective
level sets) and the signals that justify escalating to mermaid are specified in
`.spectacular/contracts/CC-projsurf-derived-state-and-dependency-shape.md`.

Everything below is an example of a notation, not a fixed template. Where a Mission's
shape is not served well by these drawings, adapt the layout to communicate that
Mission clearly. The glyph vocabulary stays stable; layout does not have to.

All examples below are drawn from a real fixture bundle:

| Ref | Outcome | Status | After | Claims |
|---|---|---|---|---|
| O1 | Extract the shared derivation layer over the typed bundle. | implemented | — | — |
| O2 | Render the compact state line and NEXT line. | in-progress | O1 | state-line |
| O3 | Compute per-claim drift flags and default the audit target. | planned | O1 | drift-flags |
| O4 | Answer authority questions by lookup. | planned | — | authority-lookup |
| O5 | Render the multi-Mission timeline. | planned | O2, O3 | timeline |

Derived readiness: `O3` and `O4` are startable, `O5` is blocked on `O2, O3`.

## M1 — Objective graph, status by colour, claims per node — DECLINED

Colour carries status and claim names sit inside the node, which the ASCII graph cannot
show without a second block. Declined because the ASCII Objective graph already reads
well for this case, and the extra information does not justify a second surface.

```mermaid
flowchart LR
  O1["O1 ✓<br/>derivation layer"]
  O2["O2 ◐<br/>state line<br/><i>state-line</i>"]
  O3["O3 ▶<br/>drift flags<br/><i>drift-flags</i>"]
  O4["O4 ▶<br/>authority lookup<br/><i>authority-lookup</i>"]
  O5["O5 ·<br/>timeline<br/><i>timeline</i>"]

  O1 --> O2 --> O5
  O1 --> O3 --> O5

  classDef done   fill:#d4edda,stroke:#28a745,color:#000
  classDef run    fill:#cce5ff,stroke:#004085,color:#000
  classDef ready  fill:#fff3cd,stroke:#856404,color:#000
  classDef block  fill:#f0f0f0,stroke:#999,color:#666

  class O1 done
  class O2 run
  class O3,O4 ready
  class O5 block
```

## M2 — Grouped by level, parallel work enclosed — APPROVED

One subgraph per level. The band of independent Objectives is visually enclosed and
labelled rather than inferred from spacing, which is what ASCII level sets cannot do.
Use when the graph branches or runs deeper than about two levels.

```mermaid
flowchart TB
  subgraph L0[" "]
    O1["O1 ✓ derivation layer"]
  end
  subgraph L1["startable in parallel"]
    O2["O2 ◐ state line<br/>Codex session"]
    O3["O3 ▶ drift flags"]
    O4["O4 ▶ authority lookup"]
  end
  subgraph L2[" "]
    O5["O5 · timeline<br/>waits O2, O3"]
  end

  O1 --> O2
  O1 --> O3
  O2 --> O5
  O3 --> O5

  classDef done  fill:#d4edda,stroke:#28a745,color:#000
  classDef run   fill:#cce5ff,stroke:#004085,color:#000
  classDef ready fill:#fff3cd,stroke:#856404,color:#000
  classDef block fill:#f0f0f0,stroke:#999,color:#666

  class O1 done
  class O2 run
  class O3,O4 ready
  class O5 block
```

## M3 — Mission timeline as a native Gantt — DECLINED

Mermaid resolves `after m7` itself, so ordering is declared rather than positioned by
hand. Declined because the ASCII timeline reads better for this case and needs no
duration for an active Mission, which the Gantt does.

```mermaid
gantt
  title Multi-Mission job
  dateFormat YYYY-MM-DD
  axisFormat %b %d

  section Done
  M5 · Compact Missions   :done, m5, 2026-08-15, 1d
  M6 · Mission CLI        :done, m6, 2026-08-16, 1d

  section Active
  M7 · Derived state      :active, m7, 2026-08-16, 2d

  section Blocked
  M8 · Frozen schema      :crit, m8, after m7, 2d
  M9 · Timeline render    :crit, m9, after m7, 2d
```

## M4 — Mission edges by kind — APPROVED

Solid means the Mission needs the produced artifact. Dotted means it needs only the
frozen interface and is therefore startable at activation. Use when a Mission splits
into related Missions that follow it, or wherever both edge kinds appear together —
the distinction is clearer here than as `╌╌╌` in ASCII.

```mermaid
flowchart LR
  M7["M7 · Derived state ◐"]
  M8["M8 · Frozen schema ·"]
  M9["M9 · Timeline render ·"]

  M7 -->|artifact| M8
  M7 -.->|interface only<br/>startable at activation| M9

  classDef run   fill:#cce5ff,stroke:#004085,color:#000
  classDef block fill:#f0f0f0,stroke:#999,color:#666

  class M7 run
  class M8,M9 block
```

## Open questions

- Colour and glyph both encode status in M2. Keeping both is redundant but survives
  monochrome rendering and copy-paste into plain text; dropping glyphs is cleaner and
  fails silently where colour is lost.
- M4 presumes the artifact-versus-interface edge split, which is M8 work. It cannot be
  built before that lands.
- These are emitted by flags on existing commands, never by a new one. CC-missioncli
  freezes a ten-command surface and M6 stops on a growing surface, so views arrive as
  `mission show <ref> --graph` and `mission show <ref> --timeline`.
- Views are written to stdout and never cached in the workspace, because `show` and
  `check` are read-only. A rendered graph is a projection, not an artifact.
- Rendering must be byte-identical whether an Objective is inline or promoted. The two
  decode paths differ, and a graph is where that difference would surface quietly.
