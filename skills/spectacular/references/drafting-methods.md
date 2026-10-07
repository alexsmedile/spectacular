# Optional drafting methods

Use this when: a planning or prototyping task benefits from a starting outline,
a guide to artifact structure, or a distinction between conventions and schemas.

## Select only what helps

| Method | Purpose | Obligation |
|---|---|---|
| Skeleton | Starting outline with unanswered parts | Optional drafting aid |
| Anatomy | Explains parts, relationships, and why they matter | Optional reasoning guide |
| Folder agreement | Recommended types and basic soft metadata | Maintained-context convention |
| Governed schema | Exact command-owned fields and constraints | Existing mechanical validator |

These are general techniques supported by Spectacular, not additional lifecycle
stages. A skeleton or anatomy does not govern a document. Reserve schema claims
for command-validated governed records; retrieve those templates through --schema
and round-trip through the emitting validator. Do not invent governed frontmatter.
Software/data schemas still belong with their owning implementation or spec.

## Small skeletons, adapted to the task

Use these prompts as needed; remove irrelevant headings and leave uncertainty
explicit. They are body outlines, not published frontmatter templates.

- Requirement: need and beneficiary; scope and constraints; observable acceptance;
  relevant sources or unresolved questions.
- Spec: applicable requirements; behavior or interface; states and failure cases;
  boundaries and verification examples.
- Plan: outcome and linked requirements; constraints; delivery slices; checks;
  open choices, progress, and result.
- Prototype: requirement and open axis; locked constraints; representative slice;
  success signals; observations and selected direction.

For an anatomy, explain how those parts relate: acceptance criteria make a need
observable, a spec describes satisfying behavior, a plan sequences delivery, and
a prototype tests an uncertain assumption. Read the existing Mission anatomy only
for explicitly selected governed work; it is not a default outline for soft plans.

## Store content, avoid companion sprawl

Start freeform drafts in the established raw/sketch/scratchpad home, or draft
directly in the intended requirements/, specs/, or plans/ file when its purpose
is clear. Respect host planning write restrictions. Evolve the draft into the
actual artifact, reusing its path when the role stays the same. When moving out
of raw is useful and authorized, preserve content and repair links; ignored raw
is not implicitly published. Do not keep a redundant skeleton beside the result.

Project-specific reusable methods can live as linked References in atlas/ when
they have lasting independent use. Otherwise keep them in the current plan or
conversation. Do not create methods/, skeletons/, or anatomies/ automatically.
Load this reference only when an outline or structural explanation helps; do not
require a prompt collection, template, checker, or extra file before starting.
