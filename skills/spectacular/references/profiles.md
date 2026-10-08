# Document profiles and metadata

Use this when: selecting metadata for a new or maintained context document,
scaling its detail, or distinguishing context conventions from governed schemas.

## Choose the smallest useful profile

The machine-readable agreement is [knowledge-folders.yaml](../knowledge-folders.yaml).
Compact is the default for maintained soft context. Profiles are conventions,
not lifecycle states, file-size limits, or mechanical schema claims. No profile,
layer, UUID, or schema field is required to select one.

| Profile | Use | Metadata |
|---|---|---|
| Freeform | raw/, sketch/, scratchpad/ | None required |
| Minimal | Small maintained context | type, version, created, updated |
| Compact (default) | Most requirements, plans, specs, scenarios, Atlas, and soft Decisions | Minimal plus recommended description; title when useful |
| Extended | Context with independently useful retrieval detail | Compact plus relevant optional tags, sources, status, or domain fields |
| Governed | Explicitly governed records; Contracts retain this route | Exact emitting command's schema and validator |

All maintained soft profiles share the same four required fields. Missing
recommended description or optional fields is not a defect. Freeform is exempt;
a raw draft can be detailed without acquiring metadata duties. Minimal is valid
without being upgraded. Extended can still be a short document.

## Select and evolve by meaning

Start new maintained documents compact unless a minimal record is enough. Add a
short description that explains what readers will find; use title only when the
filename or first heading is insufficient. Add extended fields only when they
help retrieval, attribution, or interpretation. Keep explanations and substantial
relationships in linked Markdown, not a large frontmatter database.

Reuse the same file as its needs change; preserve creation time and unknown keys.
Do not bulk-upgrade existing documents or trim their metadata to satisfy a profile.
Profiles do not mandate a heading outline or a sequence of promotions. Split by
domain, module, phase, delivered tranche, independent meaning, or lifetime when
those parts stand alone; often 50–120 lines is useful guidance.

Type is foundational and the folder agreement defines its role. Extending context
metadata does not accept a spec as a Contract or activate a Mission. Add schema
only when the emitting command mechanically validates that governed frontmatter;
retrieve its template from --schema and round-trip through that validator. Existing
identities and schema claims retain governed checks regardless of profile wording.
The keys id, ref, human_ref, schema, and schema_version claim governed treatment;
keep domain identifiers in unambiguous domain fields or the linked body instead.
Contracts keep mechanical validation and amendment/version rules. Governed
Decisions are recommended for consequential choices, not required for all context.

## Versions and dates

Version is the document revision, independent of CLI releases and schema versions.
Start at 0.1 and increment the minor number for meaningful content changes, not
routine typo fixes. Keep created immutable; update updated on meaningful edits.
Use RFC3339 timestamps with timezone. Recover historical dates from reliable
history; leave unknown dates unresolved rather than fabricating them. Git preserves
edit history. Preserve additional metadata already owned by the document.

## Inspection

The optional knowledge checker checks the four baseline fields, folder types,
and file targets. It does not require description or optional extended fields,
infer a persisted profile label, certify section anchors, or validate governed
bindings. Findings are advisory and do not block unrelated authorized work.
