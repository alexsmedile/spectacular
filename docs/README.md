# Spectacular documentation

Product documentation for people using Spectacular. For contributor rules see
[`AGENTS.md`](../AGENTS.md); for the agent-facing runtime guidance see
[`skills/spectacular/SKILL.md`](../skills/spectacular/SKILL.md).

## Start here

- **[Installation](installation.md)** — install and update the Skill and the CLI,
  which ship as two separate halves.
- **[Quickstart](quickstart.md)** — install the CLI and run one Mission end to
  end, from idea to owner completion.

## Concepts

- **[Architecture](architecture.md)** — where Spectacular keeps its files,
  what the CLI does, and how to adopt surfaces modularly (Decisions, Projections, Interview Mode).
- **[Process](process.md)** — how a piece of work moves from idea to completion,
  including standalone interview and unbundled workflows.
- **[GitHub Integration](github-integration.md)** — how Spectacular leverages native
  GitHub features (Issues, Pull Requests, Projects, Discussions, Actions, and `gh` CLI).

## Reference

- **[Human workspace contract](human-workspace-contract.md)** — the normative
  rules for what a canonical workspace must look like on disk.
- **[Testing](testing.md)** — the verification tiers and what each one proves.
- **[Release recovery](recovery.md)** — the v2 cutover baseline and v1 recovery point.
- **[Mechanical interface](../skills/spectacular/generated/mechanical-interface.md)**
  — the generated command catalog. Generated from the command registry, so it
  cannot drift from the binary. When a document and this catalog disagree, the
  catalog wins.

## Diagrams

- [Mission lifecycle](diagrams/lifecycle.svg)
- [Architecture](diagrams/architecture.svg)
- [Division of labor](diagrams/division-of-labor.svg)

---

These pages explain the product. The generated command reference and your
workspace records define the exact behavior.
