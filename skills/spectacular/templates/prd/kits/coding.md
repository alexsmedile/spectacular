---
kit: coding
version: 2.0
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

Common usage:
- CLI tools (Bash, Go, Rust, Node binaries)
- Libraries (npm, pip, gem, cargo)
- Web apps and APIs
- SDKs and integrations

Skip this kit when:
- The project ships content rather than code (→ `content`)
- The project is research feeding a decision (→ `research`)
- The project doesn't ship runnable artifacts (→ `blank`)
