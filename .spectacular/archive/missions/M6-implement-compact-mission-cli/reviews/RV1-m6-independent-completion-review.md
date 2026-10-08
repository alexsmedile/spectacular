---
type: Review
id: 01a00a34-6cb8-7082-8673-d8191806b194
title: M6 independent completion review
status: passed
created: "2026-08-16T12:30:44Z"
claims:
    - claim: typed-bundle
      verdict: pass
    - claim: schema-validation
      verdict: pass
    - claim: typed-commands
      verdict: pass
    - claim: atomic-stress
      verdict: pass
findings: []
limitations:
    - YAML shape and UUID/ref collisions are covered by exact table-driven cases rather than generated inputs, so refusal-path defects reachable only from unenumerated field combinations would not be discovered. The owner accepted this scope because those inputs are written by Spectacular itself rather than supplied by an attacker.
    - The activation fingerprint matches the current bundle, but the tree alone cannot prove it was frozen before implementation began.
    - Cold resume is proven by process separation within a test run, not by a genuinely interrupted and restarted session.
    - validateReviews short-circuits unless the Mission is completed, so checks=14 on an active Mission includes one validator that returns immediately.
mission: M6
ref: RV1
reviewed:
    activation_fingerprint: sha256:2bb4d9ff6f84db040db5eb7ecdbeb392f93aac0c242cca7bce1cfe04679ff7c5
    commit: 7fc3436b00642e0d7ccaec8cbd6e8b5e52d69f21
    tree: dea4919c66559562b1ab80eaec68bce81ed68a56
reviewer:
    actor: Claude independent reviewer
    evidence:
        - task:a2330a761f2615db4
        - task:a41f0795a864d0619
    implemented_reviewed_scope: false
    independence_basis: Two separate reviewers ran in isolated clean clones of the exact reviewed commits, made no edits, and restored the tree after every mutation probe.
    operator: Claude primary session
    relation_to_operator: independent
---
# Review

## Basis

Two independent passes. The first reviewed commit `cd4c2cc` and returned
pass-with-findings: three claims passed, `atomic-stress` was partial because its
proof_requirement names fuzz coverage of dependency graphs, YAML shape, and
collisions, and the only fuzz target was a two-seed no-panic check on plan
parsing. The bounded repair added generative dependency-graph coverage. The
second pass reviewed commit `7fc3436` and confirmed the gap closed with no
regression.

## Claim verdicts

**typed-bundle — pass.** One decoder, one package root; `decodeLegacy` is a
branch inside it, not a second reader. `TestSelfHostedMissionGoldenDecodingAndRoundTrip`
golden-decodes live M5 and M6, asserts canonical round-trip stability, and loads
legacy M3 asserting `Legacy==true` with its source pointer preserved. Promoted
review records resolve into typed values with bodies retained.
`TestExpandedBundleResolvesRecordsWithoutExpandingPointers` proves unknown fields
survive on all four document types and that pointers stay pointers on write.
`mission check M5` validates the completed Mission and leaves the tree
byte-clean.

**schema-validation — pass.** The registry is a package-private slice of 14
validators iterated unconditionally; no Mission-authored field can shorten the
loop. `TestSchemaRegistryOwnsEveryMandatoryValidation` pins the exact names.
`TestMandatoryValidatorsReturnTypedZeroMutationRefusals` runs 13 single-property
mutations from pristine clones, asserting stable code, exact field, concrete
problem, and safe correction on each, with a whole-tree content digest taken
before and after to prove zero writes.
`TestReviewedGitBindingRejectsFabricatedCoordinates` runs real `git rev-parse`
and rejects both nonexistent commits and mismatched trees.

**typed-commands — pass.** `TestPublicRegistryIsMinimalAndTyped` pins the
surface to the ten accepted commands and asserts the superseded ceremony
commands are absent. The acceptance layer builds a real binary and drives it via
`exec.Command` through the full lifecycle — stdin and file input, generated
identities and bindings, promotion without identity change, R2 creation,
independent review, owner completion, and cold resume against a workspace the
process never held in memory. Compact and `--json` output are both asserted, as
are typed usage (exit 2) and typed refusal (exit 3, `mutation: none`).

**atomic-stress — pass.** `TestMutationCommandsRollbackAtEveryInstallBoundary`
probes each mutating command's transaction width, then re-runs it failing at
every install index, asserting a byte-identical `.spectacular` digest and an
openable workspace after each — genuine fault injection at every write boundary
across all five mutating commands. Path escape, symlinked targets, stale trees,
concurrent writers, derived-target collisions, and retry convergence each have
exact adversarial cases.

`FuzzObjectiveDependencyGraph` closes the dependency-graph gap. Its oracle is
Kahn's algorithm — iterative indegree peeling — against `validateDAG`'s
recursive tri-colour marking, so the two share no logic and a common bug cannot
hide. It asserts refusal code, exact field, problem, and correction on invalid
graphs and acceptance on valid ones. Reachable inputs include multi-node cycles,
diamonds, repeated edges, self-loops, and dangling references. The reviewer ran
7.7M generated executions with no disagreement.

Mutation testing established the target is not vacuous. Four mutants were
applied to `validateDAG`: disabling the cycle guard was caught; returning `nil`
instead of a refusal — a mutant that terminates cleanly and can therefore only
be caught by a genuinely independent oracle — was caught by name on the cycle
seeds; dropping the dangling-reference check was caught; removing memoisation
survived, correctly, being behaviour-preserving. The tree was restored after
each.

## Verification

`go build ./...` and `go vet ./...` clean. `go test ./...` passes in every
package. `bash test/verify.sh` passes at `mode=all`. `mission check` returns
`valid=true schema=mission.v2 checks=14` for both M5 and M6, and `git status`
stays clean afterward.

## Findings

None blocking.

## Limitations

- YAML shape and UUID/ref collisions are covered by exact table-driven cases
  rather than generated inputs. `TestMandatoryValidatorsReturnTypedZeroMutationRefusals`
  mutates one field at a time from two real bundles, so it will not discover
  refusal-path defects reachable only from field combinations nobody enumerated.
  The owner accepted this scope on the grounds that these inputs are written by
  Spectacular itself rather than supplied by an attacker. The O3 body states
  which subjects are fuzzed and which are table-driven, so the record does not
  overclaim.
- The activation fingerprint was verified to match the current bundle, but the
  tree alone cannot prove it was frozen before implementation began.
- Cold resume is proven by process separation within a test run, not by a
  genuinely interrupted and restarted session.
- `validateReviews` short-circuits unless the Mission is completed, so
  `checks=14` on an active Mission includes one validator that returns
  immediately. Correct by design, but the count is not 14 substantive assertions
  for every Mission.
