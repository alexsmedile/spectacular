---
title: Integrations
description: How Spectacular composes with agent runtimes, Git, and public-documentation tooling.
section: guides
type: explanation
status: stable
updated: 2026-08-02
---

# Integrations

Spectacular owns the operational workspace in `.spectacular/`. It does not replace the coding agent, version control, or public-documentation system around it.

## Agent runtimes

Claude Code, Codex, Cursor, and other agents can read the same committed workspace. Install the `/spectacular` skill through a supported plugin or as plain local/global files; see [Installation](installation.md).

The contract is tool-agnostic:

- `.agents/skills/spectacular/` is the skill source.
- `.claude/skills/spectacular/` is the Claude-compatible symlink.
- `AGENTS.md` is the neutral project guide; tool-specific overrides remain optional.

## Git

Commit `.spectacular/` so decisions, requests, verification evidence, and archived work travel with the code. Keep `.spectacular.local/` ignored for personal or machine-local state.

AFK branch isolation, archive-first cleanup, and PR handoff are documented in the [command reference](commands.md) and the internal [AFK Git hygiene contract](../skills/spectacular/references/afk-git-hygiene.md).

## Public documentation

[pageworks](https://github.com/alexsmedile/pageworks) can own the public `docs/` surface. Spectacular records internal product and execution context; pageworks audits and maintains reader-facing pages. Neither tool invokes or installs the other automatically.

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash
```

When a spec-changing request is archived, Spectacular may signal that public docs need review. The user chooses whether to hand that work to pageworks.

## Extending the contract

Tools can integrate by reading frontmatter and canonical IDs rather than loading every Markdown body. Use `spectacular summary`, filtered list commands, and `spectacular paths` for cheap discovery. If an integration needs a new persistent field or lifecycle transition, treat that as a contract change rather than inferring it from prose.
