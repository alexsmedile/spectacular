---
type: implementation-evidence
mission: skill-and-runtime-prerequisites
status: ready-for-independent-review
recorded_at: 2026-08-10
recorded_by: primary-mission-owner
baseline_commit: eaef96577c0cf4bbe19debbada407f464c2aa7bb
baseline_tree: 927230b65e74dec652e9b30e63231a397b7f9c40
implementation_commit: e64f488bcf4f3279947d8d3fe9660b939f1be956
implementation_tree: ba69f34e858807c3e49d6e7b035614f52b3236aa
branch: codex/feat/v2-s-skill-runtime-prerequisites
review_authority: fresh-independent-reviewer
---

# Scenario S primary implementation evidence

## Outcome

Scenario S supplies the v2 guided Skill and runtime prerequisites without activating providers,
installing into a real user runtime, or changing v1. The implementation compiles bounded,
source-backed context; validates proportional Mission preparation; binds preparation to Mission
creation and activation; compiles explicit non-default Autopilot charters; represents Handoff as a
runtime-neutral contract; generates the mechanical Skill interface from the Go registry; and stages
the same canonical Skill and binary into disposable Codex and Claude roots.

## Material Type-2 choices

- Guardrail input is a strict, sectioned text format using the accepted `@Event` names and an
  optional `$domain.verb` selector. Owner prose is preserved verbatim; malformed, duplicate,
  unknown, or empty sections refuse.
- Context bundles are deterministic projections with explicit source bindings, generation basis,
  omissions, conflicts, Gaps, and one exact continuation or owner gate. They never become truth.
- Preparation uses a fingerprinted, expiring receipt with separate Design Sufficiency and Slice
  Quality verdicts. Mission activation revalidates the receipt, baseline, and every source.
- Autopilot is a compiled, validated, non-default charter. Its required forbidden-effect ceiling
  rules out provider-specific authority, and its recovery and return clauses remain explicit.
- The v2 mechanical catalog is generated from the Go registry into checked-in JSON and Markdown;
  an exact-parity test prevents hand-maintained command drift.
- Installation is local staging only: an explicit disposable non-empty root is required, `/`, the
  real home, and non-empty destinations refuse, and no credentials or provider effects are used.

## Claim mapping

| Claim | Primary evidence |
|---|---|
| S1 guided Skill workflows and progressive loading preserve the Skill/CLI authority boundary | `v2/skills/spectacular/SKILL.md`, routed references, registry-derived interface parity test |
| S2 context is bounded, deterministic, source-backed, freshness-aware, and non-authoritative | context compiler unit tests and the two-runtime Scenario S dogfood test |
| S3 preparation separates design and slice judgments and blocks stale or tampered activation | runtime compiler tests plus governance receipt/source revalidation tests |
| S4 Autopilot is explicit, bounded, expiring, recoverable, and provider-neutral | runtime compiler validation and command tests for source fingerprints and forbidden effects |
| S5 Handoff survives runtime replacement without moving accountability or authority | runtime-neutral Handoff contract validation and byte-identical second-runtime context proof |
| S6 Codex and Claude discover the same v2 Skill and local binary without v1 leakage | manifest validation, `v2/install/test.sh`, and cold outside/inside-workspace checks |
| S7 unsupported Message persistence remains visible as a Gap | guided Skill contract and generated interface contain no invented Message/status authority |

## Verification

The following checks passed against implementation commit `e64f488`:

```text
cd v2
gofmt -w internal cmd
GOFLAGS=-mod=readonly go mod verify
GOFLAGS=-mod=readonly go vet ./...
GOFLAGS=-mod=readonly go test ./...
GOFLAGS=-mod=readonly go test -race ./...
GOFLAGS=-mod=readonly go build ./...

cd ..
bash -n v2/install/preflight.sh v2/install/stage-local.sh v2/install/test.sh
bash v2/install/test.sh
PYTHONPATH=<cached-pyyaml> python3 <skill-creator>/scripts/quick_validate.py v2/skills/spectacular
PYTHONPATH=<cached-pyyaml> python3 <plugin-creator>/scripts/validate_plugin.py v2
bash -n cli/spectacular cli/install.sh scripts/hooks/pre-commit
scripts/hooks/pre-commit --check
./cli/spectacular --help
git diff --check
```

Focused proof covers malformed Guardrails, scope and selector refusals, deterministic context,
source staleness, preparation reslicing and readiness, receipt tampering and expiry, Mission
activation revalidation, Autopilot effect ceilings, generated catalog parity, clean-workspace
dogfood, runtime replacement, zero-mutation reads, and disposable installation discovery.

The unchanged v1 suite has one known order-sensitive test. In repeated full-suite runs on this
branch, 30 of 31 files passed and `tests/cli/wayfinding-sequencer.test.sh` failed its equal-priority
ordering assertion; that same file passed 15/15 alone. An export of the exact accepted baseline ran
the full 31/31 suite successfully. No v1 file differs from the baseline, so this nondeterminism is
recorded rather than misclassified as a Scenario S repair or silently claimed green.

## Scope and safety

The implementation diff contains 34 files: the Mission charter plus v2 Skill, runtime, command,
governance, discovery, manifest, installation, generated-interface, and test surfaces. It contains
no v1 modification, network effect, global installation, provider action, release change, push,
PR, merge, credential change, scheduler, or duplicate status authority. Known unrelated untracked
`.agents/caveman*`, `.qwen/`, and `skills-lock.json` paths were neither modified nor used as proof.

## Repair accounting and recovery

No hypothesis-changing repair round was used (`0/2`). The recoverable implementation checkpoint is
commit `e64f488bcf4f3279947d8d3fe9660b939f1be956`, tree
`ba69f34e858807c3e49d6e7b035614f52b3236aa`.

No known implementation Gap remains inside the Scenario S charter. Unsupported Message
persistence remains an explicit product Gap; release/distribution, v1 migration, real provider
effects, and central acceptance remain intentionally outside this Mission.

## Owner gate

Run one fresh, read-only independent review against the exact evidence-head commit and tree. The
review must reproduce the cold-context proof and inspect authority, stale-source handling,
Autopilot effects, runtime-neutral Handoff, manifest/install safety, v1 isolation, and generated
registry parity before the Mission can return for central acceptance.
