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

## Where the tests live

Unit tests sit next to the code they cover, as `*_test.go`. Everything the
runner drives lives under `test/`:

```text
test/
├── verify.sh          # the runner: preflight | quick | acceptance | release | all
├── acceptance/        # CLI acceptance suite, built binary against fixtures
├── evals/             # skill-behaviour benchmark, not part of the release gate
└── release.sh         # distribution gate: build, checksums, install, recovery
cmd/release-smoke/     # drives an installed binary through the governed lifecycle
testdata/              # shared fixtures (Go treats this name specially)
```

`testdata/` stays at the repository root because the Go toolchain ignores that
directory name when resolving packages, and several scripts resolve fixtures
from the root.

## Canonical runner

```sh
bash test/verify.sh preflight
bash test/verify.sh quick
bash test/verify.sh acceptance
bash test/verify.sh release
bash test/verify.sh all
```

`preflight` is a read-only, sub-second sanity gate. It emits a
`spectacular.preflight-receipt.v1` receipt and fails fast, so a heavier tier is
never spent on a workspace that is already broken.

`all` is the pre-release gate. It runs static checks, package tests, the real
binary acceptance suite, race detection, installation tests, and reproducible
release proof.

## The installed-binary lifecycle proof

`cmd/release-smoke` is the only place a released binary is driven through the
governed lifecycle as a user would drive it: separate processes, a disposable
workspace, and a fixture contract. The `release` and `all` tiers run it.

It exists because unit coverage cannot see the boundary a user meets first. A
command whose flags parse, whose validation passes, and whose typed result says
success can still write nothing, or write to the wrong path. Every mutating step
here therefore asserts that the record it reported actually landed on disk.

The order of the steps is a claim in itself. A run goes terminal when the last
objective finishes, so run transitions are exercised while the run is still
live; `contract amend` runs while the Mission that declared the Gap is live,
which is the one exemption an owner actually uses. Reordering the steps to make
the driver read more neatly will fail against those rules, correctly.
