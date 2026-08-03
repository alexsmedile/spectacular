---
title: Integrations
description: How Spectacular composes with agent runtimes, Git, and public-documentation tooling.
section: guides
type: explanation
status: stable
updated: 2026-08-03
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

## GitHub

GitHub Issues are collaborative job cards; Spectacular is used only when destination design or durable coordination adds value. `spectacular github triage <issue>` routes work to:

- `direct` — Issue → bounded agent session → draft/ready PR, with no Spectacular artifact;
- `request` — an already-defined destination receives lean PLAN/TASKS through `--from-issue` or `--from-goal`;
- `spec-first` — consequential or unsettled behavior receives an approved SPC before execution.

GitHub remains authoritative for Issue/Discussion content, labels, comments, assignees, branches, PRs, checks, permissions, and merge state. Spectacular stores canonical references and the accepted interpretation, never copied remote bodies or a second triage inbox.

The bridge opens draft PRs with a reviewer-facing integration manifest, gates ready-for-review on exact-head verification and checks, and performs read-only reconciliation. It never merges. Use raw `gh` for GitHub-only work; use the Spectacular wrapper when local lifecycle interpretation, provenance, safety gates, or reconciliation are required.

Committed `.spectacular/` is shared project knowledge. Keep incomplete ideas, machine/account preferences, and undeclared sensitive material in gitignored `.spectacular.local/`; authentication remains owned by `gh`. Protected vulnerability material never enters ordinary Issues, PR metadata, logs, or committed Spectacular files.

## Public documentation

[pageworks](https://github.com/alexsmedile/pageworks) can own the public `docs/` surface. Spectacular records internal product and execution context; pageworks audits and maintains reader-facing pages. Neither tool invokes or installs the other automatically.

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash
```

When a spec-changing request is archived, Spectacular may signal that public docs need review. The user chooses whether to hand that work to pageworks.

## Extending the contract

Tools can integrate by reading frontmatter and canonical IDs rather than loading every Markdown body. Use `spectacular summary`, filtered list commands, and `spectacular paths` for cheap discovery. If an integration needs a new persistent field or lifecycle transition, treat that as a contract change rather than inferring it from prose.
