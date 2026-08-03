---
updated: 2026-08-03
status: open
related:
  - PLAN.md
  - artifacts/workspace-inventory.md
---

# Risks — workspace-migration-readiness

| Risk | Impact | Evidence | Required mitigation |
|---|---|---|---|
| Product version and workspace schema use the same “v2” language for different contracts | Wrong migration target or overwritten history | Schema 2.0 is live while roadmap says future v1→v2 | Adopt D23 language everywhere before authoring migration code |
| No frozen breaking delta exists | A schema bump could become change-for-change's-sake | Planned roadmap fields/layout are unnamed | Approve a schema-3 capability spec before any apply path |
| Old CLI may mutate a newer workspace | Corruption of schema-3 state | Doctor warns on newer schema; a global mutator gate was not found | Add a shared fail-closed schema guard and tests before schema 3 ships |
| `status --against-latest` mislabels newer workspaces as behind | User may run the wrong recovery command | Command only checks equality then prints “behind” | Distinguish behind/newer and point newer users to CLI update only |
| Local override precedence is broader in prose than D22 allows | Private state could silently change authority or project truth | `init-workflow.md` says local “takes precedence” while CLI only enforces ignore | Replace broad precedence with an allowlisted supplement contract |
| Private local material becomes tracked | Confidentiality and credential exposure | No current leak; risk begins when local stores are implemented | Filename-only fail-closed detector, security authority, rotation/history procedure |
| Migration marks schema current without validating shape | False completion and later agent confusion | Schema 2.0 coexists with singular `debug/` and root ephemera | Validate required/forbidden shape before and after every schema bump |
| Automatic local backup duplicates secrets | Larger sensitive-data footprint | D22 permits protected local state | Initial local migration must be additive; destructive local conversion is separately authorized |
| Concurrent request/schema edits drift | Conflicts or silently incompatible assumptions | Planned stance/pack work may later touch config/frontmatter | Re-run traffic at activation, before apply, and before PR readiness |
| Non-reversible monolithic migration | Difficult recovery after partial rewrite | Existing 0.6→2.0 registry edge is marked non-reversible | Split schema-3 work into soak, guard, dry-run, and final flip; preserve branch/snapshots/manifest |
| Root and singular-path cleanup is mixed into schema design | Scope creep and misleading migration necessity | `.DS_Store`, `.last-mutation`, `migrations.log`, and `debug/` differ in meaning | Route cleanup separately; only schema-relevant rules enter SPC-004 |

## Security gate

Any discovered tracked `.spectacular.local/` path, protected-content hint, credential, or sensitive provider reference changes traffic to `unknown` and stops normal output. This request found no such tracked/history path.
