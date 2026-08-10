---
type: repair-evidence
mission: skill-and-runtime-prerequisites
status: ready-for-final-independent-rereview
recorded_at: 2026-08-10
recorded_by: primary-mission-owner
reviewed_commit: e6b8c54cbb999687271bce1ab04bd82f566883b7
reviewed_tree: 75b3e3ac4b6f5cedad513eb77a0f30860126c719
repaired_commit: b815ef4a154c0d682de6f663abbf06fb675f98d5
repaired_tree: 4e5cbe84f6cd657f139e9aa5e5466566f0540dcd
repair_round: 2
repair_budget: 2
repair_budget_exhausted: true
---

# Scenario S repair round 2 evidence

## Independent-rereview disposition

The independent rereviewer confirmed that all three round-1 findings were repaired, but returned
`BLOCK` on one remaining High finding: Autopilot checked the exact Mission and Decision record
fingerprints without revalidating either record's declared canonical freshness source. Changing
`.spectacular/workspace.yaml` left the Decision bytes unchanged, so the stale authority still
compiled. This also made the round-1 evidence claim about authority freshness inaccurate.

## Hypothesis and bounded correction

The final repair hypothesis was narrow: source-bound freshness is a separate validity dimension
from record-byte identity and must be checked before deriving any effective Autopilot limit.

`governance.Service.CompileAutopilot` now invokes the existing canonical freshness validator for
the exact bound Mission and Decision immediately after their record fingerprints match. Compilation
therefore requires parseable freshness times, an unexpired validity window, a resolvable freshness
source, and an unchanged freshness-source fingerprint before status, envelope, or Decision effects
can authorize a charter.

The dogfood suite adds independent cases for both sides of the boundary:

- changing the Mission's Evidence freshness source while leaving Mission bytes unchanged refuses;
- changing the Decision's workspace-manifest freshness source while leaving Decision bytes
  unchanged refuses.

Both return structured `insufficient_evidence` refusals with zero runtime/provider effects. The
previous successful bounded charter and all round-1 adversarial refusals remain green.

## Verification

All checks passed at repaired commit `b815ef4`:

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

The v1 suite again passed all 31 files, including the previously order-sensitive sequencer.

## Repair accounting, limitations, and recovery

This is hypothesis-changing repair round `2/2`; the declared repair budget is exhausted. It changes
only `v2/internal/governance/service.go` and `v2/internal/command/scenario_s_test.go`. The recovery
point is commit `b815ef4a154c0d682de6f663abbf06fb675f98d5`, tree
`4e5cbe84f6cd657f139e9aa5e5466566f0540dcd`.

Assurance remains limited to deterministic fixed-clock Runner reconstruction and disposable staged
runtime roots rather than launched Codex/Claude hosts. That is the accepted runtime-neutral local
prerequisite boundary, not evidence of provider or host effects.

No finding is self-closed. The independent reviewer must reproduce this exact final repaired
snapshot. A new blocker cannot receive another Mission repair without new owner authority.
