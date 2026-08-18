# Installation

Spectacular ships as **two halves that install and update separately**:

| Half | What it is | Where it travels |
|---|---|---|
| **Skill** | `skills/spectacular/` — the method, the routing, the guidance | Every host that reads the plugin manifests |
| **CLI** | `spectacular` — validation, fingerprints, atomic writes | Installed from a verified release directory |

The Skill alone is a supported, reduced mode: you can read records, learn the
method, and draft a Mission plan. Governed execution — starting, promoting,
completing, amending, recording — needs the CLI, because those produce freeze
points and transactional writes.

Check what you have:

```sh
spectacular --version
```

## Install the Skill

**Claude Code**

```sh
/plugin marketplace add alexsmedile/spectacular
/plugin install spectacular@spectacular
```

**Codex**

```sh
npx codex-marketplace add alexsmedile/spectacular --plugin
```

**Agent Plugins hosts** (ChatGPT, Cursor, GitHub Copilot, Kiro, VS Code) read the
root `plugin.json`. Point the host at the repository; skills are discovered as
immediate children of `skills/` containing a `SKILL.md`.

**Skills only, any host**

```sh
npx skills add alexsmedile/spectacular
```

**Antigravity** — clone into the plugin directory:

```sh
# workspace
git clone https://github.com/alexsmedile/spectacular .agents/plugins/spectacular
# global
git clone https://github.com/alexsmedile/spectacular ~/.gemini/config/plugins/spectacular
```

## Install the CLI

The installer works from a **locally verified release directory**. It does not
fetch a binary, require Go on your machine, or publish anything on your behalf.

```sh
install/install.sh install \
  --prefix /absolute/prefix \
  --source /absolute/release \
  --runtime claude
```

`--runtime` accepts `claude` or `codex`. Use `--os` and `--arch` to select a
platform explicitly when staging for another machine.

Verify the checksum before installing a release archive. `SHA256SUMS` ships with
every release, and the installer refuses an archive whose selected checksum does
not match — a failed verification is a refusal, not a warning.

Confirm:

```sh
spectacular --version
```

## Update

The two halves update independently.

**Skill** — re-run the marketplace command for your host:

```sh
/plugin marketplace update spectacular     # Claude Code
codex plugin marketplace upgrade           # Codex
npx skills add alexsmedile/spectacular     # npx skills
```

Antigravity installs are git clones — `git pull` in the plugin directory.

**CLI** — point the installer at the newer release directory:

```sh
install/install.sh update --prefix /absolute/prefix --source /absolute/release --runtime claude
```

Keep them roughly in step. The Skill names commands the CLI provides, so a Skill
much newer than the binary may reference a command that is not there yet.

## Recover from a bad install

```sh
install/install.sh rollback  --prefix /absolute/prefix   # restore the previous version
install/install.sh recover   --prefix /absolute/prefix   # restore after an uninstall
install/install.sh uninstall --prefix /absolute/prefix   # remove, recoverably
```

`uninstall` is recoverable by design: it leaves enough behind for `recover` to
restore the installation.

## Working without the CLI

The Skill bundles read-only shell fallbacks at `skills/spectacular/scripts/`.
They need no toolchain and run wherever the Skill lands:

```sh
sh scripts/doctor.sh          # which mode you are in, and what is unavailable
sh scripts/orient.sh          # workspace, Missions and status, what is live
sh scripts/where.sh M12       # resolve a ref to its record path
```

Where Node is available, two helpers parse frontmatter properly:

```sh
node scripts/show.mjs M12       # state, outcome, objectives, dependency edges, gaps
node scripts/check.mjs          # structural validation across every record
```

All of them read and report only — no writes, no fingerprints, and no fingerprint
verification. The shell tier reads flat fields; the Node tier checks structure;
neither verifies a binding. Prefer the CLI whenever it is installed.

## What travels, and what does not

The Agent Plugins standard covers `skills/` and `mcp.json`. Everything else in
this repository — `cmd/`, `install/`, and the `.spectacular/` governance
workspace — is outside the standard and does not travel with a plugin install.

This is why the CLI is a separate step rather than an oversight: a conformant
host installs the method, and the binary is installed alongside it when you want
governed execution.

## See also

- [Quickstart](quickstart.md) — run one Mission end to end.
- [Testing](testing.md) — verifying a build before release.
- [Release recovery](recovery.md) — the cutover baseline and recovery point.
