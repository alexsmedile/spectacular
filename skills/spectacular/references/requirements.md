# Linked requirements

Use this when: distributing a PRD into small linked objects, retaining needs and
acceptance criteria, or finding what a plan or prototype should satisfy.

## Choose one content owner

Use optional `.spectacular/requirements/` for independently useful needs,
outcomes, constraints, and acceptance criteria. Use type Requirement; a short
linking overview may use Reference. Follow the compact default in
[profiles.md](profiles.md); minimal remains valid. No required numeric identifier,
status, schema claim, or approval lifecycle is introduced.

A small plan can keep its requirements inline. Extract only when a need has an
independent reader, lifetime, reuse, or delivery scope. Existing Requirement files
in specs/ remain valid: link them without duplicating or bulk-moving them. If a
move becomes useful, preserve metadata and repair incoming links.

## Distribute a PRD by meaning

Keep the overview short: problem, desired outcomes, boundaries, and links to the
relevant requirements, domain objects, decisions, and specifications. A PRD can
remain one file when that is enough. Avoid a parallel overview that repeats every
requirement or a separate file for every acceptance bullet.

Each requirement should explain the need and who benefits, its scope and relevant
constraints, and observable acceptance criteria. Include rationale, sources,
unknowns, or dependencies only when they help. These are drafting prompts, not
mandatory sections or mechanically validated body fields. Distinguish proposed
needs from accepted constraints; agents do not invent stakeholder acceptance.

Resolve product choices progressively: ask the few questions needed for the
next usable outcome, using the owner’s language and concrete consequences.
Carry accepted constraints forward. Recommend reversible technical defaults
within scope and explain their effect; do not ask a nontechnical owner to choose
a database, queue, or protocol library without a product consequence requiring
their judgment. Defer provider, privacy, irreversible-effect, and publication
choices until their gate, while continuing independent authorized work.

Use descriptive filenames such as guest-purchase.md. Group by domain only when
several files justify it. Prefer roughly 50–120 lines; link shared constraints
instead of repeating them. Preserve creation dates and evolve the same file.

## Acceptance and value hypotheses

When a requirement affects launch or MVP acceptance, consult
[milestone-readiness.md](milestone-readiness.md). Distinguish observable behavior
and integrity constraints from the value hypothesis to be tested with users.
State the audience/environment and consequence of omission for Must requirements;
keep proposed thresholds distinct from accepted ones. Link delivery gates in the
plan rather than duplicating their state here.

## Connect needs to delivery

- Requirements answer what is needed and why, with observable acceptance.
- Specs define behavior, interfaces, and design detail that satisfy those needs.
- Plans describe delivery and verification, linking the applicable requirements.
- Atlas owns domain meanings and lasting relationships.
- Contracts retain mechanically governed agreements and amendment/version rules.

Traceability uses relative Markdown links or path-qualified wikilinks, not a new
mandatory matrix or numbered registry. A prototype links the requirements and
open question it explores; successful exploration does not prove every criterion.
Implementation reports the checked outcomes and unresolved needs in the plan.

For optional outlines and reasoning aids, see [drafting-methods.md](drafting-methods.md).
Requirements do not activate a Mission or replace an accepted Contract. The
optional knowledge checker covers baseline metadata, folder types, and file
links; it does not certify completeness, acceptance, or compliance with a PRD.
