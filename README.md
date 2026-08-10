# Spectacular v2

Spectacular is a canonical, pointer-first workspace for governed agent work.
This repository contains the v2 release candidate: `2.0.0-rc.1`.

Use the native CLI to orient from an Anchor, inspect a named record, prepare a
Mission, record Evidence and an Assessment, then reconcile and close. Public
read projections include an executable `show_command` for every emitted
pointer.

## Install from a verified release directory

The installer accepts only a locally verified release directory and never
requires Go on the consumer machine:

```sh
install/install.sh install --prefix /absolute/prefix --source /absolute/release --runtime codex
```

Build maintainers assemble the four deterministic archives with:

```sh
go run ./cmd/assemble-release --output /absolute/output --commit <commit>
```

No tag, archive upload, or provider action is performed by this repository.
For the v1 recovery point and cutover classification, see [RECOVERY.md](RECOVERY.md).
