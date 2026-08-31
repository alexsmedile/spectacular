# Quickstart

This walks through one Mission: from an idea, to an approved plan, to a record
the next agent can pick up. It takes about fifteen minutes.

You need the `spectacular` CLI on your PATH and a project that is a git
repository. See [Installation](#installation) below if you do not have the CLI
yet.

## The shape of the work: The Lean 3-Layer Model

Spectacular is designed for fast, token-efficient autonomous execution with minimal ceremony. All governance reduces to 3 simple layers:

```text
1. Ground Truth & Decisions ──▶ 2. Flight Plan Roadmap ──▶ 3. Single-File Mission ──▶ Autonomous Autopilot
   PROJECT.md + decisions/         campaigns/flight-plan.md    missions/M1.md (≤500t)     Tests pass = Proof
```

- **Layer 1: Living Ground Truth**: `PROJECT.md` (scope/boundaries) + `.spectacular/decisions/` (bulk-ideated architectural choices locked with `spectacular decide`).
- **Layer 2: Topological Flight Plan**: Multi-session roadmap in `.spectacular/campaigns/flight-plan.md` (4–8 macro milestone blocks).
- **Layer 3: Single-File Execution Envelopes**: Compact, self-contained Mission files (`missions/M1.md`, $\le 500$ tokens) with inline deliverables checklist and failable test boundaries.

## 0. Initialize the workspace (greenfield)

If starting on a new project, initialize the Spectacular workspace boundary:

```sh
spectacular init [--name <project>]
```

This safely creates `.spectacular/workspace.yaml` and seeds `.spectacular/PROJECT.md` without overwriting existing files.

## 1. Upfront Bulk-Decide (Settle Architecture Early)

Before starting implementation, brainstorm the technical stack and lock key architectural forks with `spectacular decide`:

```sh
spectacular decide decision.md
```

This writes an immutable record to `.spectacular/decisions/D1-<slug>.md`. Subagents and downstream sessions read these permanent rulings and never re-debate or hallucinate conflicting choices.

## 2. Sequence the Flight Plan Roadmap

Map 4–8 milestone blocks in `.spectacular/campaigns/flight-plan.md`. Unstarted blocks remain lightweight 4-line draft cards.

## 3. Freeze a Single-File Mission and activate it

Create a compact Mission draft (`plan.md`, $\le 500$ tokens) naming the outcome, deliverables checklist, and automated test command (`pass_boundary`). Activate it:

```sh
spectacular mission start plan.md
```

This creates `.spectacular/missions/M1-<slug>/M1-<slug>.md` as a frozen execution envelope.

## 4. Run with Autonomous Subagent / Autopilot

Delegate work to an autonomous worker subagent or execute in autopilot. Provide the worker with a compact charter ($\le 300\text{--}500$ tokens):

```text
Task: Execute Mission M1 in autopilot.
Mission File: .spectacular/missions/M1-sqlite/M1-sqlite.md
Allowed Scope: internal/storage/, cmd/
Verification Command: go test -v ./internal/storage/...
Thread Context: conversation://<parent-id>
Stop Triggers: Stop and yield back if schema migrations require external dependencies.
```

- **Zero Sub-Record Sprawl**: The worker does not create extra checkpoint or handoff files.
- **Fail-Fast Decision Gates**: If the worker hits an unrecorded fork, it halts and sends an escalation $\to$ Orchestrator records `spectacular decide` $\to$ Worker resumes.
- **Tests Pass = Proof**: Passing the test suite (`exit 0`) and creating a clean Git commit is the completion proof.

## 5. The owner completes it

Verify the mission and close the loop:

```sh
spectacular mission check M1                               # Read-only verification
spectacular mission complete M1 --by alex                  # Owner completion gate
```

Completion is an attributable owner act. The agent produces work and proof; the human approves completion.

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
