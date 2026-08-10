---
type: mission-independent-review
mission: scenario-a-cold-recovery
version: 1.0
status: complete
reviewer: /root/scenario_a_independent_reviewer
reviewed_at: 2026-08-10
reviewed_commit: 2e9fe497fb5d069cbaae6fa8d6cbd3e340518d63
reviewed_tree: ef2653aae01a71150b4ac040b1c0d40071739b35
baseline_commit: 243ae3e376e6eb30e59dc28bb691a98dfc6a7b92
verdict: bounce
next_action: central-orchestration-bounce-for-owner-authorized-repair-delta
---

# Scenario A independent review

## Reviewer identity and independence

Reviewer: Codex fresh-context independent reviewer
`/root/scenario_a_independent_reviewer`. I did not implement Scenario A, author its implementation
design, or author its primary proof. I inspected primary code, fixtures, outputs, and checks rather
than relying on the implementation return alone. I made no product, Git, lifecycle, provider, or
remote mutation.

## Reviewed revision

- Exact reviewed commit: `2e9fe497fb5d069cbaae6fa8d6cbd3e340518d63`
- Exact reviewed tree: `ef2653aae01a71150b4ac040b1c0d40071739b35`
- Exact baseline commit: `243ae3e376e6eb30e59dc28bb691a98dfc6a7b92`
- Baseline tree: `c0922b35d84617ea82e30c63b6b70729dfa87110`
- Product implementation commit recorded by primary proof:
  `231a3775bfe498135dfba60fdaa2a82ad9ff57f0`
- Product implementation tree recorded by primary proof:
  `a8b48ebd0bb95154c53fec75efb5a0d4eec49e57`
- Branch: `codex/feat/v2-scenario-a-cold-recovery`

The reviewed head is the proof-recording commit above the product implementation. The complete
`243ae3e..2e9fe49` diff contains 36 changed paths: Scenario A charter/evidence plus `v2/` code,
tests, and fixtures. It contains no v1/root CLI, Skill, hook, schema, test, or command change.

## Read set

- repository `AGENTS.md`, `.spectacular/AGENTS.md`, and `skills/spectacular/SKILL.md`;
- `SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md@1.0`;
- `evidence/scenario-a-owner-dispositions.md@1.0`;
- `evidence/scenario-a-implementation-design.md@1.0`;
- `evidence/scenario-a-primary-proof.md@1.0`;
- accepted `EXECUTABLE-REFACTOR-PROGRAM-CONTRACT.md@1.7`,
  `M1-SEMANTIC-SUBSTRATE-MISSION-CHARTER.md@1.2`, and
  `MVP-SCENARIO-CLI-SEQUENCING-DECISION.md@1.0`;
- linked immutable authority/retrieval/interface/evidence contracts needed for this review:
  `PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0`, `EXECUTION-AUTHORITY-CONTRACT.md@1.0`,
  `EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0`,
  `RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0`,
  `PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md@1.0`, and
  `SUCCESS-EVIDENCE-CONSTITUTION.md@1.0`;
- full baseline-to-head diff, all changed `v2/` implementation and test files, both Scenario A
  record fixtures, and direct CLI outputs from disposable and self-hosted runs.

The charter, owner-disposition, design, and primary-proof SHA-256 sidecars all verified.

## Methods and independently reproduced evidence

I inspected registry-to-operation bindings, discovery containment, lookup grammar, projection and
continuation code, record fixtures, adversarial tests, nonmutation snapshots, and the absence of a
write path from command dispatch. I built a fresh binary with a disposable Go cache and exercised
success, owner-gate, usage, lookup, source drill-down, self-hosting, and refusal paths.

All required automated matrices passed at the reviewed tree:

| Check | Independent result |
|---|---|
| `gofmt -l .` from `v2/` | pass; empty output |
| `GOFLAGS=-mod=readonly go mod verify` | pass |
| `GOFLAGS=-mod=readonly go vet ./...` | pass |
| `GOFLAGS=-mod=readonly go test -race -count=1 ./...` | pass |
| `GOFLAGS=-mod=readonly go test -count=1 ./...` | pass |
| `GOFLAGS=-mod=readonly go build ./...` | pass |
| `bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit` | pass |
| `scripts/hooks/pre-commit --check` | pass |
| `bash tests/run.sh` | pass; 31 files, 0 failed |

The exact ten accepted fixture commands completed in ten invocations and emitted 20,839 aggregate
JSON bytes; project orientation was 3,029 bytes. Self-hosted project orientation from
`v2/.spectacular/` was also 3,029 bytes and returned the same fingerprint-bound `resume`
continuation. Twenty warm fixture orientation processes were each below the timer's 10 ms
resolution, far below 500 ms. These independently pass the 12-call, 64 KiB aggregate, 24 KiB
orientation, and 500 ms host-local limits.

Code and tests substantiate deterministic path-first discovery, exact UUIDv7/typed/full-path
lookup, stable fixed-clock JSON, source fingerprints and generation basis, stale/missing/duplicate/
conflicting/noun/scope refusals, strict Decision-bound continuation, owner gates, disposable
in-memory projections, and byte/mode/mtime nonmutation on representative success, gate, and refusal
paths. The changed tree adds no dependency, v1 behavior, persistence call from dispatch, provider
effect, Scenario B/C mutation, compatibility alias, or generated workspace surface. The typed
envelopes and registry remain a plausible non-expanding B+C join; no B/C behavior was implemented.

## Findings by severity

### P0 critical

None.

### P1 blocking

1. **Known-command argument errors depend on workspace discovery and can return the wrong exit
   class.** `Runner.Run` opens the workspace before validating the matched command's arguments
   (`v2/internal/command/command.go:72-132`). From `/tmp`, independent executions of
   `spectacular mission show --json` and `spectacular anchor show project extra --json` returned
   `workspace_not_found`, exit 3. They are usage errors and the charter requires usage exit 2
   (`SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md:128-129`). This also means the registry is not the
   complete first authority for argument handling outside a valid workspace. Existing tests cover
   unknown-command and duplicate-JSON usage, but not malformed arguments without a workspace.

2. **The explicit v2 marker can escape through a symlink before discovery validates
   containment.** `Open` uses `Lstat` only to detect the marker and `load` immediately reads and
   unmarshals it (`v2/internal/discovery/discovery.go:56-82`); containment checks begin later with
   the metadata directory and record roots. In a copied fixture whose
   `.spectacular/workspace.yaml` symlink targeted a file outside `.spectacular`,
   `spectacular mission show not-a-ref --json` reached lookup and returned `invalid_reference`,
   proving the escaped marker had been accepted and its records parsed. The charter requires
   nearest explicit v2 discovery, declared-root confinement, and refusal of path/symlink escape
   (`SCENARIO-A-COLD-RECOVERY-MISSION-CHARTER.md:120-122,135-138`); the design specifically requires
   escape rejection before parsing records. The current symlink test covers a record symlink, not
   the authority-bearing marker.

3. **Authoritative `current_truth` accepts record nouns whose generated drill-down command is not
   in the exact registry.** `Project` skips non-Mission refs only while selecting a continuation,
   then accepts every typed record in authoritative `current_truth` and passes it to the generic
   pointer renderer (`v2/internal/projection/projection.go:127-143,174-186,650-654`). In an
   independently copied fixture with a valid Proposal ref in `current_truth`, `anchor show project
   --json` exited 0 and emitted `spectacular proposal show Proposal:...`. No such operation exists
   in the exact ten-command registry (`v2/internal/command/command.go:43-53`), so direct drill-down
   is false and registry/interface parity is broken. The committed Mission fixture also represents
   its typed Proposal source as path/fingerprint only (`projection.go:211-216`), without the
   noun/ref/show-command fields required for a consequential record pointer. This must be resolved
   without silently adding an eleventh public command: either the allowed authoritative pointer
   nouns and consequential-source rule need an owner-approved pin, or projection/schema validation
   must refuse unsupported pointer nouns while preserving an exact authorized drill-down route.

### P2 nonblocking

- The primary proof correctly distinguishes the implementation commit/tree from the reviewed
  proof-recording head, but future evidence should lead with the exact review target to prevent a
  cold reviewer from accidentally reviewing only `231a3775`.
- The independent ten-command proof used 20,839 bytes rather than the primary proof's selected
  seven-call 16,517-byte self-host sample. Both are within the accepted limit; the difference is a
  sampling choice, not a contract failure.

## Verdict

`bounce`

The core safety calculus, deterministic outputs, nonmutation proof, self-host recovery, limits,
v1 isolation, and non-expanding B+C boundary are substantially supported, but the three P1 findings
violate exact Scenario A interface and discovery contracts. They are blocking. The declared two
repair rounds are exhausted, so this reviewer did not repair or authorize further implementation.

## Exact next action

Central orchestration must bounce Scenario A and obtain an owner-authorized, narrowly scoped repair
budget/charter delta for the three P1 corrections; only then may a new exact commit/tree be produced
and sent to a fresh independent reviewer. Do not accept Scenario A or authorize Scenario B from this
review.
