---
title: Troubleshooting
description: Common setup, install, skill discovery, update, symlink, and workspace state issues.
section: ""
status: stable
since: 0.1.0
updated: 2026-05-23
---

# Troubleshooting

This guide covers common setup and usage problems.

---

## `spectacular: command not found`

The CLI is not on your `PATH`.

Install it:

```bash
curl -fsSL https://raw.githubusercontent.com/alexsmedile/spectacular/main/cli/install.sh | bash
```

Then open a new shell and check:

```bash
spectacular --help
```

The installer writes the CLI to `~/.local/bin/spectacular`. If your shell still cannot find it, add this to your shell profile:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

---

## `/spectacular` does not appear in the agent

Run initialization from the project root:

```bash
spectacular init
```

Then verify the skill exists:

```text
.agents/skills/spectacular/
.claude/skills/spectacular/
```

For project-local installs, open the agent from the same project root. Some tools only discover project skills when the current workspace is the repository root.

For Claude Code, reload plugins or restart the session if the skill was installed while Claude Code was already running.

---

## Claude sees the skill but Codex does not

Spectacular installs the source skill in:

```text
.agents/skills/spectacular/
```

The Claude path is a symlink:

```text
.claude/skills/spectacular/ -> .agents/skills/spectacular/
```

If Claude works but Codex does not, check that `.agents/skills/spectacular/SKILL.md` exists and that Codex is opened at the repository root.

---

## Codex sees the skill but Claude does not

Check the Claude symlink:

```text
.claude/skills/spectacular/
```

If it is missing, rerun:

```bash
spectacular init
```

If you installed globally, verify:

```text
~/.agents/skills/spectacular/
~/.claude/skills/spectacular/
```

Then restart or reload Claude Code.

---

## `spectacular init` skips files

This is expected. The initializer uses "write if missing" behavior for workspace files.

If a file already exists, it is not overwritten:

```text
⊘  .spectacular/PRD.md already present, leaving alone
⊘  .spectacular/config.yaml already present, leaving alone
```

If a file is empty (0 bytes or whitespace only), init fills it with a fresh stub:

```text
✓  .spectacular/PRD.md (filled empty stub)
```

If a file exists but is malformed (no frontmatter, etc.), init skips it with a hint to run doctor:

```text
⊘  .spectacular/PRD.md (issues detected — run `spectacular doctor frontmatter` for details)
```

This protects project context from accidental replacement. Edit existing files manually when needed, or use `/spectacular doctor --fix` for agent-driven repair of detected drift.

---

## Running `spectacular doctor`

Doctor is the workspace self-check. Run it any time something feels off:

```bash
spectacular doctor                    # full sweep, exit 0/1/2 per severity
spectacular doctor frontmatter        # scoped to one area
spectacular doctor --fix              # apply mechanical repairs (gitignore, missing dirs, dangling symlinks)
```

Doctor surfaces issues in 7 areas: `skill`, `workspace`, `frontmatter`, `snapshots`, `links`, `lifecycle`, `kits`. For repairs requiring judgment (frontmatter drift, missing snapshot, lifecycle mismatch), use the skill side in your AI agent:

```text
/spectacular doctor --fix
```

The skill reads doctor's JSON report, proposes context-aware fixes per finding, and walks you through `y/n/q` confirmation.

---

## `spectacular init --update` does not change project docs

`--update` updates the installed skill and `.spectacular/skills.lock`.

It does not rewrite:

- `.spectacular/PRD.md`
- `.spectacular/PRINCIPLES.md`
- `.spectacular/ARCHITECTURE.md`
- `.spectacular/roadmaps/index.md`
- `.spectacular/STACK.md`
- `.spectacular/decisions/index.md`
- `.spectacular/AGENTS.md`
- `.spectacular/config.yaml`
- existing request folders

This is intentional. Workspace content belongs to the project.

---

## Network failure during install or update

`spectacular init` downloads the latest skill release from GitHub. If the network is unavailable, install can fail while fetching the skill.

Try again when network access is available:

```bash
spectacular init
```

Or install the skill manually from a local checkout:

```bash
mkdir -p .agents/skills .claude/skills
cp -r skills/spectacular .agents/skills/spectacular
ln -s "$PWD/.agents/skills/spectacular" .claude/skills/spectacular
```

For global manual install:

```bash
mkdir -p ~/.agents/skills ~/.claude/skills
cp -r skills/spectacular ~/.agents/skills/spectacular
ln -s ~/.agents/skills/spectacular ~/.claude/skills/spectacular
```

---

## Skill triggers do not work in the shell

Only `spectacular init` is a shell CLI command.

These are skill triggers for an agent conversation:

```text
/spectacular
spectacular new <description>
spectacular advance <slug>
spectacular next
spectacular archive <slug>
spectacular remember this
spectacular snapshot <file>
spectacular idea promote <idea>
spectacular status
spectacular prd
spectacular prd refine
spectacular prd review
```

If you run `spectacular new ...` in a terminal, the CLI will reject it. Use the trigger inside Claude Code, Codex, or another agent that has loaded the skill.

---

## `.spectacular.local/` was committed by accident

`.spectacular.local/` is for personal overrides and should stay uncommitted.

`spectacular init` adds it to `.gitignore`. If it was committed before that, remove it from git tracking while keeping local files:

```bash
git rm -r --cached .spectacular.local
```

Then commit the `.gitignore` update.

---

## The workspace has stale active requests

Run:

```text
/spectacular
```

The skill should surface active, planned, review, and verified requests from `.spectacular/requests/*/PLAN.md`.

For each stale request, decide whether to:

- resume it
- move it to `review`
- mark it `verified`
- archive it
- delete it only if it was never real project history

Prefer archiving completed work over deleting it.

---

## A canonical document needs editing

Canonical docs should be snapshotted before edits:

- `.spectacular/PRD.md`
- `.spectacular/PRINCIPLES.md`
- `.spectacular/ARCHITECTURE.md`
- `.spectacular/roadmaps/index.md`
- `.spectacular/STACK.md`
- `.spectacular/decisions/index.md`
- `.spectacular/AGENTS.md`
- `.spectacular/config.yaml`
- `.spectacular/specs/index.md`
- `.spectacular/specs/**`

Use the skill trigger:

```text
spectacular snapshot .spectacular/PRD.md
```

Then edit the current file.

---

## `current/` feels empty after init

That is normal. `current/` starts empty because Spectacular should not invent system truth.

Add capability specs when:

- the project already has stable behavior worth documenting
- a completed request changes behavior
- the skill proposes a `specs/index.md` / `specs/` sync (via `SPEC-DELTA.md`) during archive

Keep `current/` behavior-oriented. It should describe what the system does now, not speculative future plans.
