# Spectacular v2

Spectacular is a canonical, pointer-first workspace for governed work.
`2.0.0-rc.2` is the human-operable v2 release candidate.

Canonical work is readable directly from `.spectacular/`: named project
Anchors live at the root, and each Mission is one cohesive directory with
scoped references such as `M1/O1`, `M1/R1`, and `M1/R1/C1`. UUIDv7 and
SHA-256 remain available for durable identity and exact revision proof.

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
spectacular anchor show project
spectacular workspace context project --event @Orient
spectacular mission list
spectacular mission show M1
```

The generated [mechanical command catalog](skills/spectacular/generated/mechanical-interface.md)
is the authoritative command reference; do not maintain a duplicate.

## Verify the implementation

Run the complete invariant, real-binary acceptance, installation, and release
matrix with:

```sh
bash test/verify.sh all
```

The test boundaries and shorter development modes are defined in
[TESTING.md](TESTING.md).

For recovery pointers and the v1 freeze, see [RECOVERY.md](RECOVERY.md).
