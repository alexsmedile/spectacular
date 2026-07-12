---
title: Installation
description: Install, upgrade, and uninstall Spectacular (CLI, Claude Code plugin, Codex plugin).
section: getting-started
type: how-to
status: stable
updated: 2026-07-12
---

# Installation

Spectacular ships in three independently-installable forms. Most users want the CLI plus one plugin.

| Form | What you get | When to use |
|---|---|---|
| **CLI** | `spectacular` binary in `~/.local/bin/` | Always — required to scaffold `.spectacular/` |
| **Claude Code plugin** | `/spectacular` slash command in Claude Code | If you use Claude Code |
| **Codex plugin** | `/spectacular` in Codex | If you use Codex |
| **Skill only** (manual) | Skill files under `~/.claude/skills/spectacular/` or project-local | Air-gapped installs / customization |

The CLI is always the first install. Plugins layer on top.

---

## CLI

### Install

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/spectacular/main/cli/install.sh | bash
```

Installs `spectacular` to `~/.local/bin/spectacular`. Make sure `~/.local/bin` is on your `PATH`:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
```

### Verify

```bash
spectacular --version
# spectacular 1.2.1
```

### Bootstrap a project

```bash
cd your-project
spectacular init
```

This scaffolds the 6-file always-set (`PRD.md`, `POLICY.md`, `config.yaml`, an agents file, `requests/`, `specs/index.md`) and installs the skill locally into `.agents/skills/spectacular/` with a symlink at `.claude/skills/spectacular/`.

Pass `-i` for interactive setup or `--kit coding` for a kit with `STACK.md` + `ARCHITECTURE.md`.

### Upgrade

```bash
spectacular init --update
```

Re-downloads the latest skill release into `.agents/skills/spectacular/` and updates `.spectacular/skills.lock`. Your `.spectacular/` content is **never overwritten** — `--update` only refreshes the skill files.

To upgrade the CLI binary itself, re-run the curl one-liner:

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/spectacular/main/cli/install.sh | bash
```

### Uninstall

```bash
rm ~/.local/bin/spectacular
```

The `.spectacular/` directory in each of your projects is plain markdown — leave it, or remove it manually if no longer needed. The skill install in `.agents/skills/spectacular/` is also plain files (symlinked from `.claude/skills/spectacular/`); remove both if you want to fully disengage.

---

## Claude Code plugin

### Install

From Claude Code, run:

```text
/plugin marketplace add alexsmedile/spectacular
/plugin install spectacular@spectacular
/reload-plugins
```

Or use the Claude CLI:

```bash
claude plugin marketplace add alexsmedile/spectacular
claude plugin install spectacular@spectacular
```

After install, `/spectacular` is available in any Claude Code session.

### Upgrade

```text
/plugin marketplace update spectacular
```

### Uninstall

```text
/plugin uninstall spectacular@spectacular
/plugin marketplace remove spectacular
```

---

## Codex plugin

### Install

Add the marketplace, then install the plugin:

```bash
codex plugin marketplace add alexsmedile/spectacular
codex plugin add spectacular@spectacular
```

### Upgrade

```bash
codex plugin marketplace upgrade spectacular
codex plugin add spectacular@spectacular
```

If adding the marketplace fails after a prior install, clear its stale local registration and add it again:

```bash
codex plugin marketplace remove spectacular
codex plugin marketplace add alexsmedile/spectacular
codex plugin add spectacular@spectacular
```

Codex's `/plugins` screen offers the same marketplace and plugin-management actions.

### Uninstall

Open `/plugins` in Codex and disable/remove `spectacular`.

---

## Skill only (manual)

If you don't want the CLI or plugin marketplace, copy the skill files directly:

```bash
git clone https://github.com/alexsmedile/spectacular /tmp/spectacular
cp -r /tmp/spectacular/skills/spectacular ~/.claude/skills/
mkdir -p ~/.agents/skills
cp -r /tmp/spectacular/skills/spectacular ~/.agents/skills/
```

You won't get `spectacular init`, `spectacular doctor`, or any other CLI verbs — only the `/spectacular` skill inside Claude Code / Codex.

---

## Install locations

| Scope | Source | Claude symlink |
|---|---|---|
| Project-local (default) | `.agents/skills/spectacular/` | `.claude/skills/spectacular/` |
| Global (`spectacular init --skill-scope global`) | `~/.agents/skills/spectacular/` | `~/.claude/skills/spectacular/` |

`.agents/` is the source of truth; `.claude/` is always a symlink. This is intentional — it keeps the skill toolchain-agnostic and lets Codex, Cursor, and any other tool that respects `.agents/` use the same files.

---

## Pairing with pageworks

If you want a public-facing `docs/` surface, install [pageworks](https://github.com/alexsmedile/pageworks) separately:

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/pageworks/main/cli/install.sh | bash
```

Spectacular's `doctor docs` will detect pageworks's presence and skip discovery hints when it's installed. See the [Pairing with pageworks](../README.md#pairing-with-pageworks) section in the README.

---

## Troubleshooting

See [troubleshooting.md](troubleshooting.md) for install issues, symlink problems, and skill-discovery failures.
