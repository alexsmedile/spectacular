# Quickstart

This walks through one full Mission: from an idea, to an activated plan, to a
completed record you can hand to the next agent. It takes about fifteen minutes
and leaves a workspace you can keep using.

You need the `spectacular` CLI on your PATH and a project that is a git
repository. See [Installation](#installation) below if you do not have the CLI
yet.

## The shape of the work

Spectacular runs **one bounded Mission at a time**. A Mission is frozen at an
activation gate, executed against that frozen text, and completed by the owner —
never by the agent. Everything else exists to make that loop safe.

```text
idea ──▶ Proposal ──▶ Mission (frozen at activation) ──▶ Run ──▶ Review ──▶ owner completes
          optional        the agreement                  the work    proof     the gate
```

A Proposal is optional exploration. A Mission is the agreement. A Run is the
attempt. Only the owner closes the loop.

## 1. Orient before doing anything

Ask your agent to orient. It reads the project Anchors — `PROJECT.md`,
`ARCHITECTURE.md`, `PRODUCT.md`, `STACK.md` — and reports what is known, open,
or blocked, without changing anything.

```sh
spectacular mission show M1
```

If the workspace is empty, there is nothing to show yet. That is the expected
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

A Mission plan is a Markdown file with typed frontmatter: an outcome, a stop
condition, and the claims that must hold. Activation freezes it:

```sh
spectacular mission start plan.md
```

This is the gate. What the Mission says at this moment is what it is judged
against later — the text is fingerprinted, and the fingerprint is recorded. An
activated Mission is not edited to match what happened; if the agreement turns
out to be wrong, it is amended deliberately through `contract amend`.

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
binary or require Go on the consuming machine:

```sh
install/install.sh install \
  --prefix /absolute/prefix \
  --source /absolute/release \
  --runtime codex
```

Verify the SHA-256 checksum in `SHA256SUMS` before installing a release archive.
The installer refuses an archive whose checksum does not match.

## Where to go next

- [Architecture](architecture.md) — what the pieces are and why they are separate.
- [Process](process.md) — the Mission lifecycle in detail, and the gates that hold it.
- [Mechanical interface](../skills/spectacular/generated/mechanical-interface.md) — the generated command catalog.
