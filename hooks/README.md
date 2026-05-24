# hooks/ — Plugin runtime hooks

This directory contains **Claude Code / Codex plugin event handlers**, loaded by the plugin runtime after a user installs Spectacular as a plugin.

- `hooks.json` — Claude Code plugin hooks
- `hooks-codex.json` — Codex plugin hooks

These are *not* git hooks. For git hooks (pre-commit, etc.) used during local development of this repo, see [`scripts/hooks/`](../scripts/hooks/).

> **Same word, different runtimes.** Plugin hooks fire inside an installed Claude Code / Codex session. Git hooks fire during `git commit` on this repo. Don't conflate.

## Schema

Plugin hooks follow the standard Claude Code hooks schema. Use `${CLAUDE_PLUGIN_ROOT}` (never hardcoded paths) for any helper script references. See the [Claude Code plugin docs](https://docs.claude.com/en/docs/claude-code/plugins) for the full event list.

Both files currently contain empty `{"hooks": {}}` stubs — spectacular ships hooks-ready but doesn't wire any events yet.
