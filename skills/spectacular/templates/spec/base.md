---
version: 1.0
updated: <DATE>
summary: "Index of what this system actually is and how it behaves right now"
related:
  - PRD.md
  - ARCHITECTURE.md
---

# <PROJECT NAME> — System Spec

<!--
  SPEC.md is the always-on index of system truth.
  PRD says what we want. ARCHITECTURE says how the workspace is shaped.
  This file says what's actually built and how it actually behaves right now.

  Keep it short. For small projects, one bullet list is enough.
  For complex projects, break out per-capability files at specs/<capability>/SPEC.md
  and reference them from here.
-->

## What this system is

<One paragraph — what the thing is today, in present tense. Skip vision, skip roadmap.
"X is a CLI that does Y, built in Z, distributed via W." Five sentences max.>

## Capabilities

<Bullet list. One line each. Each bullet is something the system can do right now.
Link out to specs/<capability>/SPEC.md only when a capability needs more than one line.>

- _no capabilities yet — write the first one when the first request ships_

## How to extend this file

- Add a bullet when a new capability ships (request → verified)
- Promote a bullet to `specs/<capability>/SPEC.md` when it grows past one line
- Snapshot before major rewrites: `spectacular snapshot .spectacular/SPEC.md`
