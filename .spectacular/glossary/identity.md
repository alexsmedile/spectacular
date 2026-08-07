---
title: Identity
status: draft
updated: 2026-08-07
summary: "The immutable machine identity of a durable Spectacular record"
---

# Identity

In Spectacular, **identity** is the immutable machine key of a durable record.
It survives a rename, path move, branch, revision, or archive operation. A
human-facing slug is not identity; it is a convenient handle.

The approved target format is UUIDv7. It replaces branch-local numeric
identifiers because independently created records must remain distinct when
branches merge.

## Related

- [[spec]]
- [[contract]]
