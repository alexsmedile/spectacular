# Reduced Mode and Fallback Guidance

Use this when: Spectacular CLI is unavailable or incompatible, and the session must operate in read/draft fallback mode.

## Core Principle

Spectacular Skill and CLI are distributed separately:
- The **Skill** travels with the plugin / harness context.
- The **CLI** binary provides typed validation, deterministic SHA-256 fingerprints, UUIDv7 identity allocation, and atomic transitions.

When the CLI is not on `PATH` (or incompatible), the workspace is in **reduced mode**. This is a valid, supported operating tier—not a broken state.

```text
┌─────────────────────────────────────────────────────────────┐
│ REDUCED MODE RULES                                         │
│ • Read canonical Markdown records                          │
│ • Orient on project direction and mission status            │
│ • Explain method and draft Mission plans for activation     │
│ ─────────────────────────────────────────────────────────── │
│ ✕ Do not fabricate command-owned records or fingerprints   │
│ ✕ Do not claim an edit was atomic when it was a plain write │
│ ✕ Do not simulate CLI verification or certification        │
└─────────────────────────────────────────────────────────────┘
```

Without the CLI, you cannot run `mission start`, `objective promote`, `objective finish`, `mission complete`, `contract amend`, `review record`, or `handoff record` because those perform transactional writes and cryptographic binding.

---

## Installing the CLI

To restore full governed execution:

1. Select the release matching the generated mechanical interface. Download its
   platform archive and `SHA256SUMS` from the corresponding GitHub release into
   a dedicated directory. Record that absolute directory as `<release-dir>`.
2. In that directory verify the downloaded archive with
   `shasum -a 256 --check SHA256SUMS --ignore-missing`; stop on failure.
3. Obtain a repository checkout or source distribution of Spectacular at the
   same release tag (`v<VERSION>`). Its `install/install.sh` is the installer;
   the platform archive is the payload, not an installer checkout. Record the
   absolute checkout/source root as `<repo-dir>`. Keep its install/ helpers together.
4. After installation is authorized, replace the placeholders below. Choose
   `codex` for Codex or `claude` for Claude; other hosts can continue ordinary work
   without this runtime installer. The command works from any directory:

```bash
bash "<repo-dir>/install/install.sh" install \
  --prefix "$HOME/.local" \
  --source "<release-dir>" \
  --runtime <codex-or-claude> \
  --version <VERSION>
```

The installer consumes the archive and checksum manifest from `<release-dir>`;
leave the archive packed. If `$HOME/.local/bin` is absent from PATH, add it before
using `spectacular` by name.

5. Confirm with `spectacular --version --json`, then apply the kernel's mechanical
   mode compatibility check before resuming command-owned work.

See the [installation guide](https://github.com/alexsmedile/spectacular/blob/main/docs/installation.md) for full installation and platform troubleshooting.

---

## Bundled Read-Only Fallbacks

The Skill bundles standalone zero-dependency scripts in `scripts/` that operate without the CLI:

Run from the target project directory so workspace discovery finds that project.
Resolve `<skill-dir>` from the loaded SKILL.md location, not the current directory;
replace it below with that absolute path. Do not change into the installed Skill
merely to locate a helper. The shell tier needs standard shell utilities.

### Shell Tier

```bash
sh "<skill-dir>/scripts/doctor.sh"          # inspect mode and report available capabilities
sh "<skill-dir>/scripts/orient.sh"          # workspace orientation, active Missions, and status
sh "<skill-dir>/scripts/where.sh" <ref>     # resolve human ref (e.g. M1, M1/O2) to record path
```

### Node.js Tier (Best-effort extraction)

Where Node.js is available:

```bash
node "<skill-dir>/scripts/show.mjs" <ref>    # state, outcome, objectives, dependency edges, gaps
node "<skill-dir>/scripts/check.mjs" [<ref>] # structural check; checks all records when ref omitted
```

### Fallback Guarantee & Limits

All bundled fallback scripts **read and report only**. None writes files,
calculates fingerprints, or verifies cryptographic bindings. `check.mjs` performs
best-effort field inspection, not YAML certification; exit 2 means validation
remains unverified, and exit 1 means field problems were observed. These scripts
are not a mechanical equivalent to the binary. Ordinary work remains available
through [knowledge.md](knowledge.md) regardless of CLI availability.

---

## Declared Manual-Bootstrap Exception

If self-development needs to draft a future record shape that the current CLI cannot represent, follow [bootstrap.md](bootstrap.md). That exception requires explicit owner authorization, stays outside governed lifecycle transitions, and never cites the legacy CLI as proof.
