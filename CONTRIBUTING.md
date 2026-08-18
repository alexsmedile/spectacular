# Contributing

Use the root v2 module only. During development, run the quick or acceptance
layer as appropriate:

```sh
bash test/verify.sh quick
bash test/verify.sh acceptance
```

Before every release-facing change, run `bash test/verify.sh all`. See
[docs/testing.md](docs/testing.md) for the boundaries and evidence owned by each layer.

`internal/command.Registry` is the single source for the generated mechanical
interface. Regenerate it when changing the registry; do not hand-maintain a
command inventory.

Follow [docs/human-workspace-contract.md](docs/human-workspace-contract.md). Canonical
records keep UUIDv7 identity but persist under readable project Anchors and
Mission bundles; generated `index.md` files are committed projections and
never authority.

Use conventional commit messages. Do not reintroduce v1 compatibility,
migrations, generic record/search verbs, or a second package root.
