---
description: Maintain Spectacular's project glossary as one concise file per Spectacular-specific term.
when_to_use: Adding, revising, splitting, or deciding whether to add a term to .spectacular/glossary/.
---

# Glossary

The glossary preserves only Spectacular-specific vocabulary. It is a retrieval
aid, not a general software dictionary.

## Admission rule

Add a term only when Spectacular gives it a precise meaning that changes how
work is routed, stored, approved, retrieved, or interpreted. Do not define
general industry primitives such as UUIDv7, YAML, Git, Markdown, pull request,
branch, or slug unless Spectacular gives the word an uncommon operational role.

## Shape

- One term per `.spectacular/glossary/<term>.md` file.
- Use kebab-case file names and a one-word or concise term title.
- Include one definition, its boundary, and `## Related` `[[links]]`.
- Keep `.spectacular/glossary/index.md` as a navigation-only list.
- Link shared rules; never copy another glossary term's definition.
