---
type: issue-evidence
source: source-011
github_issue: 38
url: https://github.com/alexsmedile/spectacular/issues/38
actionability: handoff-candidate
maps_to: [PZL-005, PZL-056, PZL-144]
retrieved: 2026-08-07
---

# GH-038 — Dublin Core documentation mapping

## Plain-language focus

Document how existing documentation metadata maps to a recognized external vocabulary.

## Problem and evidence

Docs frontmatter is informally similar to Dublin Core but the mapping and useful missing fields
are not stated.

## Proposed direction

Map existing fields without mechanically renaming everything; let renderers such as pageworks
translate to Dublin Core/Open Graph/JSON-LD at export.

## Relationships and collisions

Public documentation rendering is already assigned to pageworks. This may be a pageworks
contract or shared metadata envelope rather than Spectacular runtime behavior.

## Actionability

Confirm ownership with pageworks before editing Spectacular schemas.

## Refactor relevance

Concrete companion-skill boundary test: file handoff and standard mapping, not duplicated rendering.
