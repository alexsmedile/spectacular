# Quickstart

This walks through one Mission: from an idea, to an approved plan, to a record
the next agent can pick up. It takes about fifteen minutes.

You need the `spectacular` CLI on your PATH and a project that is a git
repository. See [Installation](#installation) below if you do not have the CLI
yet.

## The shape of the work

Spectacular runs one bounded Mission at a time. A Mission says what work is
allowed and what success looks like. Once you start it, that agreement is
locked. The agent does the work; you decide when it is complete.

```text
idea ──▶ Proposal ──▶ Mission (frozen at activation) ──▶ Run ──▶ Review ──▶ owner completes
          optional        the agreement                  the work    proof     the gate
```

A Proposal is an optional place to explore. A Mission is the agreement. A Run
is one attempt to do the work. Only you close the loop.

## 0. Initialize the workspace (greenfield)

If starting on a new project, initialize the Spectacular workspace boundary:

```sh
spectacular init
```

This creates `.spectacular/workspace.yaml` and seeds `.spectacular/PROJECT.md` safely without overwriting any existing files.

## 1. Orient before doing anything

Ask your agent to orient. It reads the project Anchors — `PROJECT.md`,
`ARCHITECTURE.md`, `PRODUCT.md`, `STACK.md` — and reports what is known, open,
or blocked, without changing anything.

```sh
spectacular mission show M1
```

If the workspace has no missions yet, there is nothing to show. That is the expected
starting state.

## 2. Explore with a Proposal (optional)

When the approach is not obvious, write a Proposal. It is mutable, carries no
authority, and exists to be argued with:

```sh
spectacular proposal check P1
```

A Proposal never grants permission to act. It is the cheapest place to be wrong.
Skip it when the work is small and the approach is settled.

## 3. Freeze a Mission and activate it

A Mission plan is a Markdown file that names the goal, the limits, and what
must be true at the end. Starting it locks that agreement:

```sh
spectacular mission start plan.md
```

What the Mission says now is what it is judged against later. Spectacular saves
a fingerprint of that text. Do not edit a live Mission to match what happened;
if the agreement is wrong, update it deliberately with `contract amend`.

**Branch before you activate.** A Mission that runs on `main` destroys the review
and isolation boundary it depends on.

## 4. Do the work in a Run

A Run is one attempt at the frozen Mission:

```sh
spectacular run start M1 --title "First implementation pass"
spectacular objective promote M1/O1     # start an Objective
spectacular objective finish M1/O1      # mark it implemented
```

Objectives are *earned*: expand one when the work is real, not upfront. Check
state at any point:

```sh
spectacular mission check M1
```

`check` validates the record graph and reports per-claim drift — which claims are
repaired, which evidence is stale, and where the next flag is.

## 5. Record proof, then hand it off

Evidence and Reviews are records, not chat messages:

```sh
spectacular review record M1 review.md
spectacular handoff record M1 handoff.md --by alex
```

A Handoff binds the exact commit and tree it was sent against, verified against
the repository, and separates `asserted` — what the sender verified — from
`assumed`, what they are carrying on trust. The receiver re-verifies everything
under `assumed` before acting.

## 6. The owner completes it

```sh
spectacular mission complete M1 --by alex
```

Completion is an owner act. An agent cannot complete its own Mission, and
completion refuses while a declared Gap is still open. This is the point of the
whole system: the agent produces work and proof; a human decides it is done.

## Installation

The CLI installs from a locally verified release directory. It does not fetch a
binary or require Go on the consuming machine.

Download the archive and `SHA256SUMS` for your platform from the
[latest release](https://github.com/alexsmedile/spectacular/releases/latest),
then install from the directory holding them:

```sh
install/install.sh install \
  --prefix "$HOME/.local" \
  --source "$PWD" \
  --runtime claude \
  --version "$VERSION"
```

`--source` is the directory **containing** the `.tar.gz` — the installer
verifies the checksum and extracts it for you, so do not unpack it first.

Confirm with `spectacular --version`. Full options, platform selection, and
update steps: [Installation](installation.md).

## Where to go next

- [Architecture](architecture.md) — what the pieces are and why they are separate.
- [Process](process.md) — the Mission lifecycle in detail, and the gates that hold it.
- [Mechanical interface](../skills/spectacular/generated/mechanical-interface.md) — the generated command catalog.
