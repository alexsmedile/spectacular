---
type: source-card
source: source-005
provided_as: source5
received: 2026-08-07
authority: code-audit-proposal
status: ingested
scope: [cli-contract, defects, git-safety, command-grammar, documentation, implementation-architecture]
duplicate_sections: []
completeness: substantial
---

# Source 005 — CLI contract audit

## Thesis

Spectacular's CLI contains valuable lifecycle interpretation, but its surface has
outgrown the PRD's “one-time bootstrap tool” boundary. Stabilize two concrete
defects first, then define a noun-first mechanical command contract, consolidate
overlapping entry points, delegate provider mutations to Git and GitHub CLI, and
modularize the implementation without changing the zero-dependency artifact.

## Source integrity

This is the first supplied source centered on code-level CLI evidence rather than
primarily product-shape inference. Its final checkout note is stale for this
branch: local `HEAD` now equals `origin/main`. Claims about Git and GitHub behavior
and the linked external manuals remain source-provided evidence; they were not
independently re-audited during intake.

## Verification audit

| Finding | Intake observation | Judgment |
|---|---|---|
| PRD calls CLI one-time bootstrap while CLI is 16,507 lines | Both are current | supported architectural drift |
| Wayfinder writes/uses `kind:` inconsistently with readers of `type:` | Direct `type` reads exist in status/next; UUIDv7 records use `kind` | supported defect |
| Intermittent Wayfinder ordering failure | focused contract suite passed 7/7; sequencer suite failed 1/15 on spike-vs-research ranking | reproduced |
| Cleanup deletes matching remote branch under `--apply --yes` | workspace and AFK paths do so; no `--delete-remote` | supported safety/contract defect |
| Remote deletion contradicts approved contracts | specs forbid autonomous remote deletion; AFK hygiene reference and tests now expect it | verified canonical collision |
| `workspace preflight` and `workspace plan` share one dispatcher path | current dispatcher confirms | supported overlap |
| workspace does not inspect sibling worktrees | no `git worktree list` use found | supported scope mismatch |
| feedback aliases dispatch identically | `feedback-loop`, `iterate`, `test`, `probe`, and `try` are present | supported |
| `_next_canonical_id()` is unused | only definition found in executable code | supported dead code |
| docs-command compatibility statement is stale | decision says removed; dispatcher lacks commands; system spec says retained until v2 | verified drift |

The Debugging policy and prior-fix signature search were run because two findings
are bug-class. No matching reusable fix record exists. Intake remains diagnostic;
no fix was implemented.

## Stabilization candidates

1. Introduce a compatibility reader for record kind and use it at every shared
   reader before changing the v2 surface.
2. Restore explicit consent for remote deletion and reconcile code, tests,
   references, specs, and decisions around one safety contract.
3. Remove the duplicate workspace entry point and dead canonical-ID helper after
   their public and documentation impacts are checked.

## Proposed v2 contract

### Keep first-class

`init`, `doctor`, `status`, request/spec namespaces, typed records, durable record
namespaces, `wayfind`, `snapshot`, `migrate`, `policy`, and `pack`.

### Fold

- summary/next into status;
- progress/links/traffic/docs-impact into request;
- verb forms such as decide/remember into noun-level add operations;
- list plurals into noun-level list operations;
- migrations into one `migrate <area>` namespace;
- mechanical imagine creation into a vision namespace.

### Remove or deprecate

Top-level lifecycle aliases, surprising feedback aliases, `workspace plan`, AFK
PR and duplicated Git mutation helpers, public internal helpers, and dead identity
allocation code.

### Reframe provider boundaries

Spectacular should own semantic interpretation, ownership/readiness manifests,
authorization records, and reconciliation. Git and `gh` should own branch and PR
mutations, with mutation remaining explicit and narrowly authorized.

## Proposed implementation sequence

1. stabilization patch;
2. accepted v2 command contract and deprecation table;
3. command registry as the authority for help, classification, docs, and tests;
4. modular Bash development sources assembled into one portable executable;
5. one compatibility release before v2 removals.

## Assumptions and contradictions to resolve

1. “CLI mechanical, skill agentic” can classify every hybrid command cleanly.
2. Native Git/gh recipes preserve Spectacular's safety guarantees and adoption
   experience better than wrapped mutation commands.
3. A command registry can become the sole authority rather than a new duplicate.
4. Noun-first consolidation does not hide useful specialized projections.
5. Modular source assembly is worth adding a build/release step to a project that
   currently advertises no build step.
6. One minor compatibility window is sufficient for existing users and scripts.
7. Wayfinding's unique dependency semantics earn a first-class namespace after
   its record-kind defect is fixed.

## Provisional assessment

**Urgent and supported:** compatible record-kind reads; explicit remote-deletion
consent; canonical contract reconciliation.

**Strong simplification candidates:** remove literal duplicate entry points,
surprising aliases, dead code, and stale documentation claims.

**Needs product decisions:** the mechanical/agentic boundary, consolidated read
views, provider mutation ownership, noun-first v2 grammar, registry authority,
modular source layout, and compatibility duration.

No command removal, compatibility promise, or implementation architecture is
accepted by this assessment.

## New concept pieces

- [PZL-047 — Compatible record-kind reader](concepts/PZL-047-compatible-record-kind-reader.md)
- [PZL-048 — Explicit remote-deletion consent](concepts/PZL-048-explicit-remote-deletion-consent.md)
- [PZL-049 — Remove duplicate workspace plan](concepts/PZL-049-remove-workspace-plan.md)
- [PZL-050 — Honest worktree inspection scope](concepts/PZL-050-honest-worktree-scope.md)
- [PZL-051 — Mechanical CLI, agentic skill boundary](concepts/PZL-051-mechanical-cli-agentic-skill.md)
- [PZL-052 — Remove surprising aliases](concepts/PZL-052-remove-surprising-aliases.md)
- [PZL-053 — Consolidate request read views](concepts/PZL-053-consolidate-read-views.md)
- [PZL-054 — Native provider mutation boundary](concepts/PZL-054-native-provider-mutations.md)
- [PZL-055 — Command registry authority](concepts/PZL-055-command-registry-authority.md)
- [PZL-056 — Generated command documentation](concepts/PZL-056-generated-command-docs.md)
- [PZL-057 — Remove dead canonical-ID allocator](concepts/PZL-057-remove-dead-id-allocator.md)
- [PZL-058 — Noun-first v2 grammar](concepts/PZL-058-noun-first-v2-grammar.md)
- [PZL-059 — Modular Bash source, single artifact](concepts/PZL-059-modular-bash-single-artifact.md)
- [PZL-060 — Compatibility release before v2 removal](concepts/PZL-060-compatibility-release.md)
- [PZL-061 — Keep dependency wayfinding, remove generic next](concepts/PZL-061-keep-wayfinding-remove-next.md)

## Decision packets seeded

- Which stabilization defects should be fixed before the larger refactor begins?
- What exact boundary separates mechanical CLI behavior from agentic skill work?
- Which projections are distinct user questions versus aliases over the same data?
- Which Git/GitHub safety semantics must Spectacular preserve when delegating mutation?
- What is the sole command-contract authority, compatibility window, and release artifact?
