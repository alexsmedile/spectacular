---
type: source-card
source: source-001
provided_as: suggestion1
received: 2026-08-07
authority: proposal
status: ingested
scope: [skill-routing, references, init, kits, measurement]
---

# Source 001 — Skill and scaffold simplification

## Thesis

Spectacular exposes too much mature-system structure at cold start. Reduce the
skill to judgment-only routing, move deterministic command knowledge into the
CLI, tier the reference layer, and make init grow structure only when use earns
it. Measure success as context loaded before an agent's first real action.

## Repository observations

Verified on `refactor/a-new-beginning` at `9ac5335`:

- `skills/spectacular/SKILL.md` is exactly 309 lines.
- It contains 132 Markdown table rows across routing, state, and explanatory
  tables; the source's “roughly 80 routing rows” needs a precise counting rule.
- `skills/spectacular/references/` contains 77 flat Markdown files, including 23
  `*-rules.md` files.
- `ALWAYS_SET="prd spec config agents policy"` is current CLI behavior. A blank
  init creates PRD, specs/index, config, AGENTS, POLICY, requests/, specs/, and a
  root `.gitignore` entry/change.
- `spectacular help status` currently prints global help. Per-command help is not
  yet an authority that can replace skill routing rows.
- The coding kit declares STACK and ARCHITECTURE as always-triggered, but a fresh
  observed coding-kit init created only the always-set. This is contract/runtime
  drift to diagnose, not a stable premise for the redesign.

## Proposals, normalized by domain

| Domain | Proposal | Intake state |
|---|---|---|
| Skill boundary | Keep intent judgment, an eight-domain router, canonical rules, and the execution state machine | promising |
| CLI authority | Delete command-documentation rows and teach the CLI to answer per-verb help | needs contract design |
| Doc registry | Replace thin rules stubs with one machine-readable registry | needs authority decision |
| Reference retrieval | Replace the flat directory with workflow/engine/rules/contracts tiers | promising, migration-heavy |
| Doctrine placement | Remove doctrine, glossary, and retention material from the model path | goal promising; destination disputed |
| Init floor | Default to PRD, config, and requests only | needs lifecycle compatibility analysis |
| Lazy growth | Create collections and canonical docs only on first demonstrated use | promising, triggers disputed |
| Kit semantics | Make kits PRD presets plus contextual suggestions rather than eager doc triggers | needs product decision |
| CLI compatibility | Make minimal the default; add `--full`; remove `--minimal` and reserved-collection flags | separate compatibility decisions |
| Measurement | Benchmark context loaded before the first code-relevant action; target under 150 lines | strong direction, metric needs definition |

## Assumptions and contradictions to resolve

1. **CLI help dependency.** Router rows cannot be deleted safely before
   per-command CLI help exists and covers the semantic information agents need.
2. **Registry authority.** Rules-file frontmatter and `doc-index.md` already act
   as registries. Adding `registry.yaml` without retiring an authority would
   create the duplicate truth this source aims to remove.
3. **Bash cost.** A YAML registry needs a Bash 3.2-compatible reader or a generated
   representation. Its operational cost must be lower than the stubs it replaces.
4. **Docs boundary.** Root `docs/` is currently owned by pageworks for public
   documentation. Moving internal runtime doctrine there would supersede that
   boundary or require a different non-runtime location.
5. **Spec-first lifecycle.** Current consequential work begins with an approved
   spec before request creation. Creating `specs/` only during archive is
   incompatible as written and may conflate draft contracts with system truth.
6. **Configuration ownership.** `config.yaml` currently exposes user-owned kit,
   policy, snapshot, and tool overrides. “Machine-owned, user never edits” is a
   broader configuration redesign, not merely a smaller scaffold.
7. **Agent substrate.** Removing AGENTS from init changes how a cold agent learns
   workspace rules. The skill or another deterministic surface must absorb that
   responsibility without increasing context elsewhere.
8. **Usage thresholds.** “Third request” and “fifth archive” are arbitrary until
   user evidence shows those moments predict a need for principles or roadmap.
9. **Kit reality.** Declared coding-kit triggers and observed init behavior
   disagree. Diagnose current behavior before choosing what compatibility means.
10. **Metric gaming.** Line count is a useful proxy, not the outcome itself. A
    shorter route that causes wrong intent classification, excess tool calls, or
    unsafe action is a regression.

## Valuable ideas independent of the full design

- Treat cold-start context as a measured product budget.
- Remove deterministic command syntax from model instructions once an equivalent
  CLI help surface is proven.
- Route by a small number of user-intent domains rather than enumerating verbs.
- Use paths and names as retrieval signals instead of one flat reference folder.
- Make empty workspace structure justify itself through real use.
- Test the whole retrieval path on realistic tasks, not only file line counts.

## Provisional assessment

**Adopt as evaluation goals:** context-budget measurement, command-help authority,
domain-level routing, and earned structure.

**Compare before selecting architecture:** YAML registry, exact reference tiers,
the three-item init floor, lazy-creation triggers, kit semantics, and flag removal.

**Do not implement as written:** moving internal doctrine into public `docs/`,
creating specs only at archive, arbitrary request/archive thresholds, or deleting
router knowledge before the CLI supplies equivalent semantics.

No product decision is recorded by this assessment.

## Concept pieces

- [PZL-001 — Domain-routed lean skill](concepts/PZL-001-domain-routed-lean-skill.md)
- [PZL-002 — CLI command-help authority](concepts/PZL-002-cli-command-help-authority.md)
- [PZL-003 — Single dispatch registry](concepts/PZL-003-single-dispatch-registry.md)
- [PZL-004 — Tiered reference layout](concepts/PZL-004-tiered-reference-layout.md)
- [PZL-005 — Runtime/doctrine separation](concepts/PZL-005-runtime-doctrine-separation.md)
- [PZL-006 — Minimal default scaffold](concepts/PZL-006-minimal-default-scaffold.md)
- [PZL-007 — Lazy grow-on-write substrates](concepts/PZL-007-lazy-grow-on-write.md)
- [PZL-008 — Kits as PRD presets](concepts/PZL-008-kits-as-prd-presets.md)
- [PZL-009 — Init flag redesign](concepts/PZL-009-init-flag-redesign.md)
- [PZL-010 — Cold-start retrieval budget](concepts/PZL-010-cold-start-retrieval-budget.md)

## Decision packets seeded

- What belongs in model context versus deterministic CLI help?
- What is the single authority for document dispatch metadata?
- What reference hierarchy minimizes retrieval cost without harming maintenance?
- What is the honest minimum viable workspace substrate?
- Which structures should be eager, lazy-on-first-write, or suggested?
- What do kits own: PRD shape, workspace structure, or both?
- Which init compatibility changes justify a breaking release?
- What benchmark measures fast, correct, safe cold-start behavior?
