# Test Architecture

Spectacular tests observable contracts in four layers. A lower layer may make
diagnosis faster, but it never substitutes for a higher boundary.

| Layer | Owns | Required boundary |
|---|---|---|
| Invariants | identity, parsing, fingerprints, layout, transactions | Go package APIs and real temporary filesystems |
| Workspace contracts | discovery, authority, projections, governed mutation | canonical human-layout workspaces reopened after every mutation |
| CLI acceptance | human output, JSON envelopes, exit classes, pointer traversal, cold resume | separately built `spectacular` processes with no shared memory |
| Distribution | archives, checksums, installation, update, rollback, recovery | assembled native artifacts in disposable prefixes |

## Principles

- Assert outcomes, refusals, and preserved state rather than private call order.
- Official fixtures use the current human workspace layout. Old flat record
  directories are not test conveniences.
- Every emitted `show_command` must run successfully in a fresh process against
  the unchanged workspace.
- Read-only operations must preserve bytes, modes, and modification times.
- Every successful mutation is followed by cold rediscovery from disk.
- Failure tests use real filesystem and process boundaries where the operating
  system behavior is part of the claim.
- Generated indexes, help, and schemas are compared with their deterministic
  generators; checked-in projections never become test authority.
- Coverage is diagnostic. Release acceptance is scenario and invariant based,
  not a percentage target.

## Canonical runner

```sh
bash test/verify.sh quick
bash test/verify.sh acceptance
bash test/verify.sh release
bash test/verify.sh all
```

`all` is the pre-release gate. It runs static checks, package tests, the real
binary acceptance suite, race detection, installation tests, and reproducible
release proof.
