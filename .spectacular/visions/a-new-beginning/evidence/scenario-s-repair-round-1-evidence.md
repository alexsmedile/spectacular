---
type: repair-evidence
mission: skill-and-runtime-prerequisites
status: ready-for-independent-rereview
recorded_at: 2026-08-10
recorded_by: primary-mission-owner
reviewed_commit: 568ae88c5ea40938384656907094e9bf6bc3b5d6
reviewed_tree: 67ac6416dea58b71b3dec2bfc597034a3f00f2fb
repaired_commit: 7bede74ea4de75e1d6f887bc44d41c5a932e95cc
repaired_tree: 2f02997cc6b69978382ebe079c594d312d1a81ce
repair_round: 1
repair_budget: 2
---

# Scenario S repair round 1 evidence

## Independent-review disposition

The fresh reviewer reproduced the primary checks but returned `BLOCK` with two High findings and
one Medium finding:

1. Autopilot checked source existence and fingerprints but did not prove that its Decision
   authorized `mission.autopilot`, the requested providers/actions, or that its budgets, expiry,
   stops, and effects remained inside the bound Mission envelope.
2. Mission preparation was optional. When present, activation checked selected duplicated fields
   rather than revalidating the stored receipt content and its content fingerprint.
3. Disposable staging compared raw path spelling, allowing the real home through `$HOME/.` or an
   equivalent symlink.

The reviewer also noted assurance limits: cold-runtime equality uses independently reconstructed
Runners with a fixed clock rather than launched Codex/Claude hosts; staging proves discovery files,
not host launch. Those are retained as honest limits and do not imply a provider or host effect.

## Hypothesis and bounded correction

The repair hypothesis was that all three defects came from validating the requested representation
without first deriving an immutable effective boundary from canonical source identity.

- Autopilot now resolves the exact active Mission and owner Decision, validates their current
  fingerprints, requires Decision operation `mission.autopilot`, exact Mission target/fingerprint,
  owner disposition, scope, authorized actions and provider effects, conditions, freshness, and
  non-supersession, then passes Mission-derived hard limits to the runtime compiler. The compiler
  refuses outcome, provider, action, forbidden-effect, budget, repair, expiry, or stop expansion.
- Mission creation now requires a ready preparation receipt, validates all bound source bytes,
  binds selected outcome, evidence claims, stops, baseline, proposal, and both verdicts, and stores
  the full receipt. Activation parses that stored receipt, recomputes and checks its fingerprint,
  freshness and readiness, rebinds its Mission fields and exact source list, and revalidates every
  source fingerprint before the owner-authorized transition.
- The staging script canonicalizes physical destination and home paths before comparison. The
  installer test proves both `$HOME` and `$HOME/.` refuse with no writes.

## Adversarial proof

The clean Scenario S dogfood now creates a separate exact owner Decision for Autopilot. A replay of
the reviewer exploit—GitHub provider, `merge` and `push`, budget `999`, and the required merge
forbiddance—returns an unauthorized structured refusal. Unit tests also refuse forbidden/allowed
overlap, an action outside Mission authority, and a budget above the Mission. Governance tests
refuse Mission creation without preparation and preserve stale-receipt activation refusal.

## Verification

All checks passed at repaired commit `7bede74`:

```text
cd v2
gofmt -w internal
GOFLAGS=-mod=readonly go mod verify
GOFLAGS=-mod=readonly go vet ./...
GOFLAGS=-mod=readonly go test ./...
GOFLAGS=-mod=readonly go test -race ./...
GOFLAGS=-mod=readonly go build ./...
bash -n install/preflight.sh install/stage-local.sh install/test.sh
bash install/test.sh

cd ..
PYTHONPATH=<cached-pyyaml> python3 <skill-creator>/scripts/quick_validate.py v2/skills/spectacular
PYTHONPATH=<cached-pyyaml> python3 <plugin-creator>/scripts/validate_plugin.py v2
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
./cli/spectacular --help
bash tests/run.sh
git diff --check
```

The v1 suite passed all 31 files, including 15/15 in the previously order-sensitive sequencer.
This does not erase the earlier observed nondeterminism; it confirms the repaired Scenario S diff
does not modify v1 and the full regression was green at this checkpoint.

## Repair accounting and recovery

This is hypothesis-changing repair round `1/2`. It changes only nine v2 runtime/governance/command/
installer code and test files. The recovery point is commit
`7bede74ea4de75e1d6f887bc44d41c5a932e95cc`, tree
`2f02997cc6b69978382ebe079c594d312d1a81ce`.

No finding is self-closed. The same independent reviewer must inspect and reproduce the repaired
snapshot before the Mission can return for central acceptance.
