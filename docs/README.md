# Spectacular documentation

Product documentation for humans. For contributor rules see
[`AGENTS.md`](../AGENTS.md); for the agent-facing runtime guidance see
[`skills/spectacular/SKILL.md`](../skills/spectacular/SKILL.md).

## Start here

- **[Quickstart](quickstart.md)** — install the CLI and run one Mission end to
  end, from idea to owner completion.

## Concepts

- **[Architecture](architecture.md)** — the four surfaces, the thirteen record
  types, and why the CLI deliberately decides nothing.
- **[Process](process.md)** — the Mission lifecycle and the gates that hold it:
  activation, review, completion.

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

Nothing in `docs/` is authoritative for behavior. It explains; the generated
interface and the records define.
