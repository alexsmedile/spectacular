---
type: Plan
version: "0.2"
created: "2026-10-07T19:26:45Z"
updated: "2026-10-07T19:31:13Z"
description: Add linked requirements and optional drafting methods, then publish v3.1.0.
---

# Linked requirements and optional drafting

## Outcome

Distribute independently useful PRD needs in requirements/, preserving existing
Requirement files in specs/. Support optional skeletons and anatomy guides during
planning and prototyping without adding ceremony, commands, or governed schemas.

[Product need](../requirements/linked-product-needs.md) owns acceptance criteria.

## Implementation

- Add Requirement and Reference to the optional requirements/ folder agreement.
- Link [requirements guidance](../../skills/spectacular/references/requirements.md)
  and [drafting methods](../../skills/spectacular/references/drafting-methods.md)
  from ordinary context, planning, and prototyping routes.
- Keep compact metadata as default; retain raw freedom and Contract mechanics.
- Update public docs, README, workspace navigation, versions, and plugin manifests.
- Verify new folder diagnostics, existing-file compatibility, and links, then run
  the full release gate and publish annotated v3.1.0 from main.

## Verification and result

Seven context diagnostic tests pass, including linked delivery, legacy specs/
Requirements, and misplaced-type detection. Both modified Skills pass frontmatter
validation. Context diagnostics report 28 maintained documents and zero findings;
changed Markdown file targets resolve. The full local release gate passed
(race, acceptance, distribution), with staged changes passing the secret scan.
Publish v3.1.0 through the tag-triggered release workflow and verify its assets.
Release authorization is explicit; historical governed records remain untouched.
