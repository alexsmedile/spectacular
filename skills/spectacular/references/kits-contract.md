# Kits Contract — diff-only extensions over a base doc

Loaded by `grill.md` / `review.md` when the active doc has `kit-support: true` in the registry (currently only PRD).

## Core principle

**Kits are typed deltas, not standalone templates.**

A kit declares:
1. What slots it **adds** to the base
2. What slots it **modifies** in the base (description-level only — never removes or renames base slots)
3. What other docs it **triggers** (always-scaffold + suggested-scaffold)

The base PRD remains the contract. Kits extend it. The review gate always runs base checks; kit-specific checks run only when that kit is declared in the PRD's frontmatter (`kit: <name>`).

## Kit file shape

A kit is a markdown file with structured frontmatter + optional human-readable body explaining the kit's intent.

```yaml
---
kit: <kit-id>
version: 1.0
extends: prd
adds-slots:
  - name: <Slot Name>
    after: <existing-slot-name>     # where in slot order to insert (after this slot)
    required: true | false          # does the review gate enforce this slot?
    prompt: |                       # grill question for this slot
      <question text>
    example: |                      # optional example answer
      <example>
modifies-slots:
  - name: <existing-slot-name>
    note: |                         # additional guidance layered onto the base prompt
      <kit-specific guidance>
triggers-docs:
  always:                           # docs smart-init MUST scaffold when this kit is selected
    - <doc-id>
    - <doc-id>
  suggested:                        # docs smart-init OFFERS to scaffold (interactive y/n)
    - <doc-id>
    - <doc-id>
description: |
  <one-paragraph description of when to pick this kit>
---

# <Kit Name>

<Optional longer-form description, examples, tips.>
```

## Field semantics

**`kit`** (required) — kebab-case identifier. Used in frontmatter (`kit: coding`) of the resulting PRD and in registry lookups.

**`extends`** (required) — the base doc this kit extends. Currently only `prd`. Future docs may opt into kit support by setting `kit-support: true` in the registry.

**`adds-slots`** (optional, list) — slots inserted into the base slot order. Each entry has:
- `name` — slot heading text (becomes `## N. <name>` in scaffolded file, where N is computed by insertion order)
- `after` — name of the base slot this should immediately follow. Use `last` to append after all base slots.
- `required` — `true` means review gate flags missing/empty; `false` means optional (skipped if user leaves blank)
- `prompt` — grill question text (engine surfaces this verbatim)
- `example` — optional good-answer example shown to user

**`modifies-slots`** (optional, list) — kit-specific guidance layered onto a base slot's prompt. Never changes the slot's name, position, or required status. Each entry has:
- `name` — must match a base slot name exactly
- `note` — appended to the base prompt during grill (e.g. coding kit's "Deliverable" note: "Name the binary, library, or package format.")

**`triggers-docs`** (optional, two lists) — declares which root docs this kit needs alongside PRD. Consumed by [[smart-init]] to drive doc scaffolding decisions:
- `always` — smart-init creates these without asking
- `suggested` — smart-init prompts y/n per doc in interactive mode; ignores in non-interactive

**`description`** (required) — one-paragraph summary used in the grill's kit-selection menu.

## Base slot list (the contract)

PRD base slots, in canonical order (from `prd-rules.md` slot prompts):

1. Vision
2. Problem
3. Target users
4. Deliverable
5. Goals & success criteria
6. Non-goals
7. Constraints
8. First milestone

Kits insert via `after: <base-slot-name>`. Example: coding kit's "Stack" slot uses `after: Constraints` → becomes slot 8, pushing "First milestone" to slot 9.

## What kits CAN do

- Add new slots after any base slot
- Mark added slots as required or optional
- Layer kit-specific guidance onto base slot prompts (without renaming)
- Declare triggered root docs (always + suggested)
- Provide kit-specific examples and prompts

## What kits CANNOT do

- Rename base slots (`Vision`, `Problem`, etc. are fixed names)
- Remove base slots (all 8 base slots always present)
- Change base slot order relative to each other
- Override the base review gate's universal checks (placeholder + clarification + frontmatter)
- Apply two kits simultaneously to one PRD (v1: single-kit only — see Composition below)

## Composition (v1: single-kit only)

Each PRD has exactly one kit declared in frontmatter. Multi-kit composition is **out of scope for v1** because:

1. **Slot insertion conflicts** — two kits could both want to insert after the same base slot; no clean resolution rule.
2. **Trigger-docs conflicts** — two kits could disagree on whether a doc is `always` or `suggested`.
3. **Modify-slot conflicts** — two kits could both layer guidance onto the same base slot.
4. **Review gate complexity** — combining kit-specific checks across multiple kits introduces edge cases.

**v2 multi-kit composition sketch** (deferred):
- Apply kits in declaration order; later kits' modifications override earlier ones for the same slot
- Trigger-docs unions: any kit's `always` wins; `suggested` accumulates
- Slot insertion: error if two kits target the same `after:` value without explicit ordering hints
- Frontmatter: `kits: [coding, content]` (array) instead of `kit: coding`

## Engine integration

### Grill flow (kit-aware)

1. **Kit selection** (pre-flight) — engine shows registry-discovered kits + descriptions; user picks one (or `blank`)
2. **Slot resolution** — engine merges base slots + kit's `adds-slots` in declared order
3. **Prompt resolution** — for each slot, engine uses (in priority): kit `prompt` (for added slots) > base prompt + kit `modifies-slots.note` (for layered slots) > base prompt alone
4. **Grill loop** — walks merged slot list normally; mini-refine applies to base slots only unless kit declares slot-specific patterns
5. **Frontmatter write** — engine sets `kit: <kit-id>` in the resulting PRD.md frontmatter

### Review gate (kit-aware)

1. **Base checks always run** — slot presence, placeholder, clarification, frontmatter
2. **Read PRD frontmatter** — if `kit: <name>` present, load that kit file
3. **Kit-required slots checked** — every `adds-slots` entry with `required: true` must be filled
4. **Kit-optional slots ignored** — `required: false` slots may be empty without failing the gate
5. **Engine ignores kit additions if no kit declared** — base-only PRDs still pass without ever knowing about kit slots

## Triggered-docs contract (consumed by smart-init)

```yaml
triggers-docs:
  always:
    - stack         # smart-init: scaffolds .spectacular/STACK.md unconditionally
    - architecture  # smart-init: scaffolds .spectacular/ARCHITECTURE.md unconditionally
  suggested:
    - principles    # smart-init: in -i mode asks "scaffold PRINCIPLES.md? [y/n]"; in non-interactive mode skips
    - roadmap
```

Each entry must be a valid `<doc-id>` from `doc-registry.md`. smart-init errors on unknown doc IDs.

The always-set (`PRD.md`, `SPEC.md`, `requests/`, `specs/`, `config.yaml`, `AGENTS.md`) is scaffolded by smart-init regardless of kit. Kits add to this set; they never subtract.

## Override path (project-local kits)

A project can override a bundled kit by placing a same-named kit file at `.spectacular/templates/prd/kits/<kit-id>.md`. Project-local takes precedence over the skill's bundled kit. Used for customizing kit defaults per project.

Rules file must conform to the same kits-contract schema; the engine validates frontmatter on load.

## Worked example: coding kit

```yaml
---
kit: coding
version: 1.0
extends: prd
adds-slots:
  - name: Stack
    after: Constraints
    required: true
    prompt: |
      What tech does this run on? Language, runtime, key deps, distribution mechanism.
    example: |
      Bash 5+, macOS/Linux. Single-file CLI installed to ~/.local/bin via curl one-liner.
  - name: Interfaces
    after: Stack
    required: false
    prompt: |
      What surfaces do users interact with? CLI commands, API endpoints, UI screens.
    example: |
      CLI: `spectacular init`, `spectacular new`, `spectacular <doc> <verb>`.
modifies-slots:
  - name: Deliverable
    note: |
      For coding projects: name the concrete artifact (CLI binary, library, npm package, etc.).
      Don't just say "a tool" — name what users install/import.
triggers-docs:
  always:
    - stack
    - architecture
  suggested:
    - principles
    - roadmap
    - decisions
description: |
  Coding projects: CLIs, libraries, apps, services, SDKs.
  Adds Stack + Interfaces slots. Triggers STACK.md + ARCHITECTURE.md scaffolding.
---

# Coding kit

For software projects shipping installable or runnable artifacts.

Pairs naturally with `spectacular init --kit coding`, which scaffolds STACK.md and ARCHITECTURE.md alongside the always-set.
```

When applied during grill:
- Resulting PRD has slots in order: Vision / Problem / Target users / Deliverable / Goals / Non-goals / Constraints / **Stack** / **Interfaces** / First milestone
- Frontmatter includes `kit: coding`
- Slot 4 (Deliverable) prompt includes the coding kit's "name the concrete artifact" note
- Review gate enforces all 8 base slots + Stack (required) + skips Interfaces if empty (optional)
- smart-init scaffolds STACK.md + ARCHITECTURE.md alongside PRD.md

## Adding a new kit

1. Pick a kit-id (e.g. `mobile-app`)
2. Create `templates/prd/kits/<kit-id>.md` with full frontmatter per schema
3. (Optional) Document in kit-selection menu via SKILL.md or scaffold-reference.md
4. No engine changes needed — discovery is automatic

## Related

- [[doc-registry]] — declares which docs support kits (currently only `prd: kit-support: true`)
- [[grill]] — consumes kit's `adds-slots` and `modifies-slots`
- [[review]] — runs base checks always, kit checks only when kit declared
- [[prd-rules]] — base PRD slot prompts that kits layer onto
- [[smart-init]] (planned) — consumes `triggers-docs` to drive scaffolding
