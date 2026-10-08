# Naming, reading routes, and retirement

Use this when: navigating a workspace, naming or moving a document, distinguishing
configuration from knowledge, rebuilding inventories, or retiring finished work.

## Filenames and one content owner

Root anchors use uppercase semantic names: PROJECT.md, ONTOLOGY.md, GUARDRAILS.md,
and other independently useful concerns. ONTOLOGY.md is preferred for new domain
models; retain an existing VOCABULARY.md until a rename is authorized, without
creating duplicate authority. Preserve identity and repair incoming soft links.

Soft objects use descriptive lowercase kebab-case: guest-purchase.md or
checkout-first-delivery.md. Use domain subfolders only when several objects justify
them. Avoid numeric sorting prefixes, dates, or identifiers without a real need.
Governed filenames and references follow their existing writer and archival rules.

A folder owns a role, not a mandatory stage. Reuse the existing topic owner; split
by independent meaning, lifetime, domain, module, phase, or delivered tranche.
Requirements own needs and acceptance, specs own behavior, and plans own delivery.
Skeletons and anatomies are optional methods; do not create containers for them.

## Read by task, not filesystem order

A useful route is INDEX.md → PROJECT.md → relevant Requirement → Spec → Plan,
consulting ONTOLOGY, Atlas objects, Decisions, and guardrails when they matter.
The index and intermediate objects are optional: a direct task can start from its
own plan or raw draft. Read only the useful linked context, not entire collections.

Manual INDEX.md separates current work, reusable knowledge, and historical work.
Update it when active work changes; do not label delivered plans as current.
Relative Markdown links resolve from their source; path-qualified wikilinks resolve
from .spectacular. Repair incoming links when moving soft files. A directory listing
or completed Mission's presence does not establish current execution authority.

## Knowledge, configuration, and generated data

| Surface | Role | Needed? |
|---|---|---|
| workspace.yaml | Existing governed CLI marker, scan roots, project Anchor pointer | Governed CLI only |
| config.yaml | Optional supported runtime overrides | Only for actual overrides |
| INDEX.md | Authored task-oriented routes through context | Optional |
| index.json | Generated non-authoritative governed inventory | Rebuildable |
| Collection index.json | Optional filtered view of the same governed graph | Writer-generated, rebuildable |
| .engine/ | Temporary internal execution state | Only during supported execution |
| .cache/ | Disposable internal cache, never an authority | Only if an actual consumer needs it |

YAML is human-editable configuration, consistent with frontmatter; JSON is generated
data. Neither configuration file activates a Mission. Ordinary context needs no
manifest, config, or CLI. Do not copy built-in defaults into config.yaml or present
an undocumented field as an enforced override. Generated inventories state their
scope; they do not catalog all soft context or certify bindings and completion.

For this repository, `go run ./scripts/rebuild-workspace-index.go` prints the root
governed inventory; add --write to atomically refresh index.json. It requires the
source checkout and Go. It uses existing discovery/projection code, refuses an
authored or unrecognized destination, and never writes INDEX.md. This developer
script adds no public CLI command. Other workspaces use their available tooling.

Generate JSON, never Markdown indexes or a duplicate catalog.json cache. Existing
legacy workspaces can retain historical navigation until cleanup is authorized.

## Retire by outcome, preserve provenance

Use archive/ as the retirement home. Completed Mission bundles move together and
retain identity, status, activation and Contract bindings, and frozen proof. Record
owner authorization and the original canonical fingerprint under existing archival
policy; never re-point a completed Mission or edit its frozen claims to fit today.

Retire a Proposal only when its question is answered: write resolved_by and accepted
status before moving, with authorization and source fingerprint. Partial delivery
stays open. Archive independently useful completed plans and campaign inputs under
archive/plans/ and archive/campaigns/; do not keep stale activation instructions live.

Legacy imports may use archive/raw/ to preserve their original bytes and unsupported
historical frontmatter outside discovery. They remain source material, not accepted
truth or runtime instructions. Do not migrate their schemas or revive a v1 reader.
Keep historical governed bodies intact; repair links in mutable context and use
current navigation/typed references to find archived proof. No second history/ root
is needed, and no archive folder is created until something needs retention.
