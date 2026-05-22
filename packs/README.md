# Spectacular pack app store

This folder is Spectacular's distribution point for community convention packs. Each subfolder is one installable pack.

A **pack** is a small folder that declares how a repo should be shaped — naming rules, folder taxonomy, required root files, gitignore defaults, file placement rules, project-type templates. Packs are mini-skills: same shape (`pack.md` + `templates/` + `references/`), different purpose. Full schema in [`skills/spectacular/references/packs-contract.md`](../skills/spectacular/references/packs-contract.md).

Spectacular ships with `minimal` bundled (in `skills/spectacular/templates/packs/minimal/`). It enforces only a `.gitignore` baseline + a README contract. Stronger opinions are opt-in — install a pack from here when you want them.

## Available packs

| Pack | Status | Description |
|---|---|---|
| `alex-default` | placeholder (planned in [convention-pack-fabricator](../.spectacular/requests/convention-pack-fabricator/PLAN.md)) | Opinionated defaults: kebab-case naming, mono-collection taxonomy, role suffixes, AGENTS.md pattern, project-type scaffolds for cli/library/skill/plugin/content/research/vault-project |

`alex-default` is reserved as the dogfood pack — request 2 (convention-pack-fabricator) produces it by grilling against the maintainer's existing conventions, then commits the result here.

## Installing a pack

```bash
spectacular pack install <pack-id>
```

Lands at `~/.spectacular/packs/<pack-id>/`. Activate per-repo via `config.yaml`:

```yaml
convention_pack:
  source: <pack-id>
  mode: suggest | scaffold | enforce
```

> **Note:** `spectacular pack install`, `pack list`, `pack remove`, and the `convention_pack:` config block all land in request 3 ([convention-pack-application](../.spectacular/requests/convention-pack-application/PLAN.md)). Until that ships, packs in this folder are reference material — install them manually by copying to `~/.spectacular/packs/<name>/`.

## Modes (when application layer lands)

- **`suggest`** — skill mentions pack opinions when relevant during interactive work, never auto-applies
- **`scaffold`** — init and new-request actively apply the pack's rules
- **`enforce`** — doctor flags drift from the pack's rules during normal use

## Contributing a pack

1. Run `spectacular pack new <name>` (lands in request 2) — interactive grill produces `~/.spectacular/packs/<name>/`
2. Validate with `spectacular pack review`
3. Copy the produced folder to `packs/<name>/` in this repo
4. Open a PR — include the rationale in `references/why-<name>.md` so users can read what the pack stands for before installing
