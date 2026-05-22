---
version: 1.0
updated: <DATE>
summary: "The .spectacular/ directory — layers, file roles, frontmatter conventions, lifecycle, versioning"
related:
  - PRD.md
  - PRINCIPLES.md
  - AGENTS.md
---

# <Project Name> — Architecture

<!--
  Architecture describes the project's own structural conventions.
  Distinct from STACK.md (which describes the host project's tech choices).
  This file is mode: freeform — no slot grill, edit directly.
-->

This document defines the **structure of `.spectacular/`** — what each folder and file is for, how they relate, and the conventions every file must follow.

---

# Layout

```txt
.spectacular/
├── PRD.md
├── PRINCIPLES.md
├── ARCHITECTURE.md
├── ROADMAP.md
├── AGENTS.md
├── STACK.md
├── DECISIONS.md
├── config.yaml
│
├── ideas/
├── current/
├── requests/
├── skills/
├── memory/
└── archive/
```

---

# Root layer

<!-- Describe each root file's purpose. -->

| File | Purpose |
|---|---|
| `PRD.md` | <PURPOSE> |
| `PRINCIPLES.md` | <PURPOSE> |
| ... | ... |

---

# Frontmatter conventions

<!-- Required + optional fields per doc type. -->

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "..."
related:
  - <sibling-file>.md
---
```

---

# Layers

<!-- Describe each layer (ideas, current, requests, skills, memory, archive). -->

## Ideas

<DESCRIPTION>

## Current

<DESCRIPTION>

## Requests

<DESCRIPTION>

## Skills

<DESCRIPTION>

## Memory

<DESCRIPTION>

## Archive

<DESCRIPTION>

---

# Lifecycle

<!-- State machine for request lifecycle. -->

```txt
planned → active → review → verified → archived
```

---

# Versioning

<!-- Snapshot-before-overwrite convention. -->

Canonical documents are **never overwritten in place**. Snapshot first using `<DOC>@vN.md` naming.

---

# Related

- [PRD.md](PRD.md)
- [PRINCIPLES.md](PRINCIPLES.md)
- [AGENTS.md](AGENTS.md)
