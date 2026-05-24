---
version: 1.0
updated: 2026-05-24
summary: "Audience profiles and the user stories that drive Spectacular's build decisions"
related:
  - PRD.md
---

# Personas

> Spectacular's own dogfood. Three personas, each ~10 lines. Stories drive what we build next.

## Solo OSS maintainer

**Who** — One-person shop shipping a CLI/library/skill in their spare time. Day job pays the bills; OSS work happens after hours.

**Wants to** — Ship coherent features without re-deriving context from scratch every weekend session.

**Pain** —
- Last weekend's decisions are invisible to the agent two weeks later — every session starts cold
- No time for formal product docs; PRD/SPEC discipline feels like overhead until something breaks

**Stories** —
- As a solo maintainer, I want to scaffold a `.spectacular/` workspace in seconds, so that I get coherence without a planning ceremony
- As a solo maintainer, I want the agent to read project state on its own, so that I stop pasting context into every prompt
- As a solo maintainer, I want one canonical file per concept (PRD, SPEC, PLAN), so that I never wonder which doc is the source of truth

**Not for** — Maintainers who already have a working spec/issue/wiki system they're happy with — Spectacular won't out-organize a healthy existing setup.

---

## Small-team tech lead

**Who** — Lead on a 2-5 person engineering team using Claude Code, Codex, or Cursor for real work. Owns the architecture, reviews PRs, runs the roadmap meeting.

**Wants to** — Give every teammate (human and agent) the same operational context, so that decisions made in one session compound across the team.

**Pain** —
- Architectural decisions made in a 1:1 evaporate; nobody can find the "why" three weeks later
- AI agents accelerate code generation but produce subtly inconsistent output across team members because they each prompt differently

**Stories** —
- As a tech lead, I want a per-request PLAN.md decomposition, so that any teammate (or their agent) can pick up the work mid-flight
- As a tech lead, I want a DECISIONS.md log written by whoever made the call, so that I stop being the human archive
- As a tech lead, I want PRINCIPLES.md to enforce house rules at agent runtime, so that style/architecture drift doesn't show up at PR time

**Not for** — Teams >10 where formal RFC/design-doc processes already enforce coherence.

---

## Tool builder using AI agents to build AI tools

**Who** — Builds skills, plugins, CLIs, or agent workflows for a living. The artifact they ship *is* an agent system or developer tool.

**Wants to** — Maintain coherent intent across nested layers (their own project, the agents they spawn, the tools their agents call) without losing the thread.

**Pain** —
- Recursive complexity: the project has agents that have agents; context loading becomes its own problem
- Convention drift across sibling projects — three CLIs, three different `.gitignore` shapes, three different naming schemes
- Specs and READMEs constantly outdated relative to actual shipped behavior

**Stories** —
- As a tool builder, I want convention packs that encode my repo-shape opinions, so that every new project I start follows the same conventions without manual setup
- As a tool builder, I want SPEC.md to be a present-tense index of what's actually built, so that my agent's plan never contradicts shipped state
- As a tool builder, I want frontmatter-driven progressive context loading, so that the agent reads the smallest correct slice of the project — not the whole repo
- As a tool builder, I want my agents to share the same `.spectacular/` workspace as me, so that we operate from one source of truth

**Not for** — One-off scripts and throwaway prototypes — Spectacular's overhead doesn't pay back at that scale.
