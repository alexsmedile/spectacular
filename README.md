# Spectacular v2

Spectacular is a canonical, pointer-first workspace for governed work.
`2.0.0-rc.1` is the first externally testable v2 release candidate.

## Install a verified local release

The installer accepts only a locally verified release directory and does not
require Go on the consumer machine:

```sh
install/install.sh install --prefix /absolute/prefix --source /absolute/release --runtime codex
```

## Orient in a workspace

Start from the project Anchor, compile bounded context, then follow every
emitted `show_command` pointer rather than guessing a path or command.

```sh
spectacular anchor show project --json
spectacular workspace context project --event @Orient --json
spectacular mission list --json
```

The generated [mechanical command catalog](skills/spectacular/generated/mechanical-interface.md)
is the authoritative command reference; do not maintain a duplicate.

For recovery pointers and the v1 freeze, see [RECOVERY.md](RECOVERY.md).
