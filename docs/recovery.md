# Cutover recovery manifest

- v2 root-cutover baseline: `7775d44b72e9d31e535f0d399d032ad313ccc8e6`
- baseline tree: `cdd749b538b3d0f289949d8b83bb6113e6ba0ef4`
- accepted release-readiness merge: `7837bf3d210985d018c2471eda3c232be4bea643`
- v2 root-cutover merge: `7710a332d1d9c61394473691deea1422e0034114`
- frozen v1 public surface: annotated tag `v1.37.3` at `9ac53358d391cb3a1d7405b0553ff836347fd738`
- constitutional and design record: `refactor/a-new-beginning` at `7775d44b72e9d31e535f0d399d032ad313ccc8e6` (tree `cdd749b538b3d0f289949d8b83bb6113e6ba0ef4`)

The v1 tree is intentionally not retained in this working tree or archive.
Normal Git history and the immutable `v1.37.3` tag are the recovery substrate.
Restore or inspect v1 from that tag; `main` is the v2 product line.

Cutover classification: root v1 CLI, Skill, plugins, installer, hooks, tests,
packs, agents, migration machinery, broad docs and loose tracking files were
replaced by the root v2 module or are history-only. The live `.spectacular/` is
the small v2 self-hosted RC workspace with current Contract, Evidence, Mission,
and owner-disposition gate.
