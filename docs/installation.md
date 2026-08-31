# Installation

Spectacular has **two parts that install and update separately**:

| Half | What it is | Where it travels |
|---|---|---|
| **Skill** | `skills/spectacular/` — the method, the routing, the guidance | Every host that reads the plugin manifests |
| **CLI** | `spectacular` — validation, fingerprints, atomic writes | Installed from a verified release directory |

With the Skill alone, you can read records, learn the method, and draft a
Mission plan. You need the CLI to start, update, or complete a Mission because
it safely writes and checks those records.

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

**Antigravity** has no marketplace or install command — a plugin is any folder
carrying a root `plugin.json` that sits in a scanned directory. Place it
yourself, by clone or by symlink:

```sh
# workspace — read by both the desktop app and the agy CLI
git clone https://github.com/alexsmedile/spectacular .agents/plugins/spectacular

# global — Antigravity 2.0 desktop
git clone https://github.com/alexsmedile/spectacular ~/.gemini/config/plugins/spectacular

# global — agy CLI (a separate state tree; installing to one does not serve the other)
git clone https://github.com/alexsmedile/spectacular ~/.gemini/antigravity-cli/plugins/spectacular
```

The desktop app and the `agy` CLI scan **different** global directories. A
plugin present in one is not visible to the other.

Working from a local checkout, symlink it instead of cloning twice — the link
tracks your tree with no update step:

```sh
ln -s /absolute/path/to/spectacular ~/.gemini/config/plugins/spectacular
```

Verify what Antigravity actually loads:

```sh
agy plugin validate ~/.gemini/config/plugins/spectacular
```

Note that `agy plugin list` reports *imported* plugins only, and stays empty for
plugins discovered by directory scan. An empty list is not evidence that the
install failed — `validate` is the check that answers it.

## Install the CLI

The installer never reaches the network: it takes a local directory that already
holds the release archive, verifies the checksum, and extracts it. It does not
fetch a binary, require Go on your machine, or publish anything on your behalf,
so piping it from a URL refuses — there is no download step to trigger. Download
first, then install from that directory:

```sh
# 1. download the archive and SHA256SUMS for your platform
VERSION=2.10.0
PLATFORM=darwin-arm64          # or darwin-amd64, linux-amd64, linux-arm64
BASE=https://github.com/alexsmedile/spectacular/releases/download/v$VERSION

curl -LO $BASE/spectacular-v$VERSION-$PLATFORM.tar.gz
curl -LO $BASE/SHA256SUMS

# 2. verify before unpacking
shasum -a 256 --check SHA256SUMS --ignore-missing

# 3. install from the directory holding the archive
install/install.sh install \
  --prefix "$HOME/.local" \
  --source "$PWD" \
  --runtime claude \
  --version $VERSION
```

Pass `--version` explicitly, or also download the release's `VERSION` asset into
the same directory — the installer reads the version from that file when the
flag is absent, and refuses with `release version is invalid` if neither is
present.

`--source` is the directory that **contains** the `.tar.gz` and `SHA256SUMS` —
not an unpacked tree. The installer re-verifies the checksum, inspects the
archive inventory, extracts to a staging area, and checks that the binary,
generated interface, Skill, and runtime manifest all report the same version
before anything is placed. Unpacking it yourself first is not a step.

`--prefix` must be an absolute path, and cannot be `/` or your home directory
itself. `$HOME/.local` is the usual choice: the binary lands at
`$HOME/.local/bin/spectacular`, which is already on `PATH` on most systems. If
`spectacular --version` is not found afterwards, add that `bin` directory to
`PATH`.

The installer refuses when any target path is an existing **symbolic link**,
rather than following it and writing through. A leftover link from an older
install must be removed before installing:

```sh
ls -l "$HOME/.local/bin/spectacular"   # inspect what it points at first
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
# Claude Code — refresh the marketplace, then update the plugin.
# Updating alone does not pull a newer marketplace snapshot, and the
# plugin id must be namespaced or it is reported as not found.
claude plugin marketplace update spectacular
claude plugin update spectacular@spectacular   # restart to apply

# Codex — refresh the snapshot, then re-add to install it.
# There is no `codex plugin update`; upgrading the marketplace moves the
# snapshot without touching the installed plugin, so the `add` is what
# actually lands the new version.
codex plugin marketplace upgrade spectacular
codex plugin add spectacular@spectacular

# npx skills
npx skills add alexsmedile/spectacular
```

Confirm which version a host actually resolved before assuming it updated —
a registered marketplace can stay pinned to an old revision indefinitely:

```sh
codex plugin list | grep spectacular
```

Antigravity installs are git clones — `git pull` in the plugin directory. A
symlinked install tracks its target and needs no update step.

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
