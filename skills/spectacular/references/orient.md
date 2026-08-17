# Orient

Answer one question: where does this project stand, and what happens next.

## Read

| What | Where |
|---|---|
| Project Anchor | `.spectacular/PROJECT.md` |
| direction sources | the Anchors and Contracts its `current_truth` field names |
| the selected Mission | `.spectacular/missions/<slug>/MISSION.md` |

If `.spectacular/PROJECT.md` does not exist, the repository is greenfield (uninitialized). Do not report an error: route immediately to One-Shot Genesis in [prepare.md](prepare.md).

To see what exists without reading every record, use the generated
`.spectacular/missions/index.md`. It maps each ref to its source path and is
non-authoritative — never cite it as proof.

Stop there. Do not preload history, every record, or generated catalogs.

## Report

Lead with plain project direction, then:

- current Mission outcome and owner
- Contract and Git baseline
- validation mode
- current Objective and Run
- blocking Gaps or stops
- one continuation, or one owner gate

## Inspect with the right tool

```bash
spectacular mission show  <ref> --json   # current state
spectacular mission check <ref> --json   # validation only, read-only
```

`<ref>` is the Mission's `ref` from its frontmatter — `M7`. Sub-records use a
slash: `M7/O1` for an Objective, `M7/R1` for a Run. The full command list is in
[../generated/mechanical-interface.md](../generated/mechanical-interface.md).

Under `manual-bootstrap`, read the canonical Markdown and its referenced files
directly — see [bootstrap.md](bootstrap.md). An incompatible CLI refusal is a
tooling Gap. It is not evidence that the Mission is invalid.

## When several Missions are active

Say which one this session will operate, and why.

Never silently combine two Missions' authority, criteria, scope, or mutable state.
