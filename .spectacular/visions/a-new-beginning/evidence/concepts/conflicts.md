---
type: concept-conflict-matrix
updated: 2026-08-07
---

# Contradiction and dependency matrix

| Piece | Collision or dependency | Kind | Current state |
|---|---|---|---|
| PZL-001 | Removing command rows depends on PZL-002 providing equivalent deterministic semantics | dependency | open |
| PZL-002 | Current `spectacular help <verb>` falls back to global help | implementation gap | verified |
| PZL-003 | A new YAML registry may duplicate rules frontmatter and `doc-index.md` | authority duplication | open |
| PZL-003 | YAML consumption must remain Bash 3.2-compatible | technical constraint | open |
| PZL-004 | Moving 77 references changes links, tests, packaging, and retrieval contracts | migration impact | open |
| PZL-005 | Root `docs/` is currently the pageworks-owned public surface | contract collision | open |
| PZL-006 | Removing specs from init conflicts with the current spec-first lifecycle as stated | lifecycle collision | open |
| PZL-006 | Making config machine-owned conflicts with user-owned policy and tool overrides | authority collision | open |
| PZL-006 | Removing AGENTS requires another cold-agent instruction authority | safety dependency | open |
| PZL-007 | Lazy creation must not make absent optional collections indistinguishable from unsupported capabilities | substrate ambiguity | open |
| PZL-008 | Coding-kit declared triggers disagree with observed fresh-init output | existing drift | verified |
| PZL-009 | Removing flags bundles unrelated compatibility decisions | decision coupling | open |
| PZL-010 | Line/token reduction can be gamed at the expense of correctness, safety, or extra tool calls | measurement risk | open |
| PZL-011 | “Choose one world” tensions with PZL-021's space/expedition/neutral hybrid | internal contradiction | open |
| PZL-012 | Anchors already exists, but source and current architecture disagree on membership | scope drift | verified |
| PZL-013 | Capability Contract may clarify existing spec semantics or narrow non-capability specs | taxonomy collision | open |
| PZL-014 | Atlas adds metaphor where familiar `specs/` already provides a neutral path | comprehension trade-off | open |
| PZL-015 | Mission may mean project purpose, goal, request, release, or long-running operation | semantic ambiguity | open |
| PZL-016 | Journey may duplicate PLAN, TASKS, wayfinding, native plans, or AFK run state | abstraction duplication | open |
| PZL-017 | Waypoint can weaken “milestone = demoable validated outcome” into a progress marker | contract dilution | open |
| PZL-018 | Universal Launch/Operations conflicts with the separate request and release lifecycles | lifecycle collision | verified |
| PZL-018 | Some supported project types have no deployment or live-production phase | applicability gap | open |
| PZL-019 | Observed production behavior may be a defect rather than intended product truth | authority ambiguity | open |
| PZL-020 | Mission Control suggests orchestration and operational authority excluded by the PRD | positioning collision | verified |
| PZL-021 | Minimal hybrid needs a principled boundary or recreates the mixed-world defect | design-rule gap | open |
| PZL-022 | Telemetry/Signals may collapse feedback, verification, observation, and benchmark evidence | evidence-boundary collision | open |
| PZL-006 + PZL-012 | The minimum scaffold determines which Anchors are universal versus earned | cross-source dependency | open |
| PZL-010 + PZL-011 | Vocabulary should be tested for retrieval correctness, not selected by taste alone | cross-source evidence need | open |
| source-003 | Sessions 7–9 are incomplete or missing from the supplied nine-session plan | source completeness | verified |
| PZL-023 | Sessions are described as independent but carry explicit dependencies and reordering | sequencing contradiction | verified |
| PZL-024 | Archive is called mechanical despite closure, spec-sync, policy, and human gates | authority collision | verified |
| PZL-025 + PZL-036 | Shipping 57 planned tasks may expand the surface during reduction | scope collision | open |
| PZL-026 | “Three records or first-hour stranger use” is arbitrary and partly unmeasured | evidence gap | open |
| PZL-027 | Visions is now active and AFK has two run records plus a branch ledger | stale deadness claim | verified |
| PZL-027 | Local sparse usage does not establish product-wide non-value | inference gap | open |
| PZL-028 | FIX reservation and the live legacy F1–F5 fix workflow are separate contracts | identity collision | verified |
| PZL-030 | Nine reduction requests may violate spec-first routing and recreate scaffold overhead | process collision | open |
| PZL-031 | `contrib/` preserves dead code in live retrieval; a Git ref/tag does not | retention collision | open |
| PZL-032 | Predicted midpoint counts depend on unapproved and now-stale deletion assumptions | measurement gap | open |
| PZL-033 | Wayfind and status returned different next actions from the same workspace | non-equivalence | verified |
| PZL-034 | The source refers to eight hooks while current POLICY defines nine | stale count | verified |
| PZL-035 | The port target and sessions surrounding the final strategy choice are missing | source completeness | open |
| PZL-027 + current Vision | Deleting visions would remove the active refactor evidence database | self-hosting collision | verified |
| source-004 | Usage table mixes records, files, indexes, and token occurrences | measurement invalidity | verified |
| source-004 | DEC is called unused despite 27 current decision records | internal contradiction | verified |
| source-004 | SPC is counted as 1001 while the live spec list contains 9 records | unit mismatch | verified |
| PZL-037 | Self-hosting frequency cannot prove external-user or preventive value alone | evidence boundary | open |
| PZL-038 | Protected core must preserve outcomes without freezing current file shape | abstraction risk | open |
| PZL-039 | QUE/RES/SPK differ in authority and evidence, not only labels | contract collision | verified |
| PZL-040 | Three risk gates may erase phase-specific correctness safeguards | policy collision | open |
| PZL-041 | Nine agent files do not establish dispatch cost or role redundancy | evidence gap | open |
| PZL-042 | The ~6k Bash target is not derived from accepted behavior | measurement gap | open |
| PZL-043 | Vision is active and shipped only three days before the source date | deadness claim refuted | verified |
| PZL-044 | AFK is sparse but has two runs; safety value may be preventive | deadness claim disputed | open |
| PZL-045 | Wayfinding and status produce different live decisions | equivalence claim refuted | verified |
| PZL-046 | Earlier file count is not evidence of better adoption or outcomes | nostalgia risk | open |
| source-005 | Its “local checkout is two commits behind” note is stale on the refactor branch | source staleness | verified |
| PZL-047 | Canonical `kind` records are ranked through direct legacy `type` reads | implementation defect | verified |
| PZL-047 | Compatibility fallback preserves legacy shape temporarily rather than enforcing immediate normalization | migration trade-off | open |
| PZL-048 | Cleanup code and tests permit remote deletion while specs and approved authority records forbid it | canonical contract collision | verified |
| PZL-048 + PZL-054 | Requiring `--delete-remote` conflicts with removing Spectacular-owned provider mutations entirely | ownership alternative | open |
| PZL-049 | Two public names share one dispatcher path but may have distinct documented promises | compatibility collision | open |
| PZL-050 | “Worktree coordination” currently inspects only the active checkout | scope overclaim | verified |
| PZL-051 | Remember, imagine, status, and spec flows currently cross the proposed mechanical/agentic boundary | existing contract collision | verified |
| PZL-052 | Alias implementation is verified, but external usage is unmeasured | compatibility evidence gap | open |
| PZL-053 | Overlapping state does not prove that specialized projections answer the same user question | semantic-equivalence gap | open |
| PZL-054 | Native Git/gh mutations may reduce wrapper code while distributing safety and recovery semantics | boundary trade-off | open |
| PZL-055 | A command registry may duplicate dispatcher logic just as PZL-003 may duplicate doc metadata | authority duplication | open |
| PZL-055 + PZL-003 | Combining command and document registries could centralize unrelated schemas | centralization risk | open |
| PZL-056 | Generated contract tables must not overwrite pageworks-owned human documentation | ownership collision | open |
| PZL-057 | Executable helper is unused, but canonical prose still reflects the older identity model | documentation drift | verified |
| PZL-058 | A coherent noun-first grammar requires many independently reviewable compatibility choices | decision coupling | open |
| PZL-059 | Source modularity introduces assembly to a repository whose current contributor contract says no build step | architecture collision | verified |
| PZL-060 | One minor-release window is proposed without adoption or script-usage evidence | migration evidence gap | open |
| PZL-061 + PZL-045 | Source 005 retains dependency Wayfinder; Source 004 deletes the subsystem | direct proposal conflict | open |
| PZL-061 | Wayfinder uniqueness is supported, but its present ranking correctness is not | evidence ordering | verified |
| source-006 | Proposed `mission plan|run|close` makes the CLI agentic while the PRD calls it a one-time bootstrap | product-boundary collision | verified |
| source-006 | One coding agent mutates directly while current agent contracts reserve general mutation for the orchestrator | authority collision | verified |
| PZL-006 + PZL-077 | Source 001's minimal PRD/config/requests floor and Source 006's PROJECT/SYSTEM/capabilities/missions floor are incompatible shapes | direct proposal conflict | open |
| PZL-063 | PROJECT/SYSTEM overlap PRD, AGENTS, PRINCIPLES, ARCHITECTURE, STACK, and specs index | authority duplication | verified |
| PZL-063 | Repository inspection can infer implementation facts more reliably than product purpose or forbidden operations | inference boundary | open |
| PZL-064 | `capabilities/` may duplicate or rename the existing `specs/` authority | taxonomy collision | open |
| PZL-065 | Capability-local sections can duplicate interfaces and state shared across several capabilities | normalization trade-off | open |
| PZL-066 | Mission combines approved direction, execution plan, autonomy contract, and closure delta currently owned by separate artifacts | authority compression | open |
| PZL-067 | One plan lock conflicts with distinct current gates for direction, activation, irreversible action, verification, closure, PR readiness, and merge | consent collision | verified |
| PZL-068 | Outcome objectives align with current milestone doctrine, but still require explicit risky technical seams | completeness boundary | open |
| PZL-069 | Reusing the host runtime aligns with the PRD, but `mission run` may recreate a runtime abstraction inside the CLI | boundary ambiguity | open |
| PZL-070 | Derived run context lowers retrieval cost but may omit safety-critical context | compilation risk | open |
| PZL-071 + PZL-051 | Source 006 places autonomous judgment behind CLI; Source 005 requires agentic work under the skill | direct proposal conflict | open |
| PZL-071 | Optional draft PR delivery does not apply to every supported project type | applicability gap | open |
| PZL-072 | Current safety contracts resemble the proposed stops, but the remote-cleanup regression shows prose alone does not enforce them | enforcement gap | verified |
| PZL-073 | A new EVIDENCE file may duplicate VERIFY, VERIFY-LOG, SPEC-DELTA, check output, and PR evidence | authority duplication | open |
| PZL-074 | Closure reconciliation converges with current archive/spec-sync; artifact replacement is not required to preserve the outcome | implementation independence | verified |
| PZL-075 | Status and SESSION support continuity, but complete cold-agent reconstruction has not been demonstrated | acceptance evidence gap | open |
| PZL-076 | Absorbing research, fixes, tasks, verification, and sessions assumes their authority and evidence differences are incidental | semantic-collapse risk | open |
| PZL-076 + PZL-018 | Source 006 closes after proof; Source 002 keeps Mission open through Launch and Operations | lifecycle conflict | open |
| PZL-077 | Git-only document history conflicts with current explicit snapshot and staged-migration contracts | recovery collision | verified |
| PZL-078 | Draft-PR acceptance needs a provider-neutral equivalent for non-GitHub or non-code projects | test portability | open |
| PZL-079 | Several “deferred” capabilities already ship; deferral does not specify freeze, default-off, deprecation, or deletion | action ambiguity | verified |
| source-007 | Its three attachments support ingredients but do not establish the typed graph or reconciliation architecture | provenance overreach | verified |
| source-007 | “Independent agent approval” conflicts with its rule that reviewer confidence is not evidence | internal authority ambiguity | open |
| PZL-080 | A maintained graph may describe accepted contracts while code/tests and production remain stronger truth authorities | truth-label collision | open |
| PZL-080 | No current graph relationship vocabulary or thin-slice retrieval result exists | architecture evidence gap | verified |
| PZL-081 | Capability and architecture are distinct axes, but separate storage is not implied by that distinction | implementation independence | verified |
| PZL-082 + PZL-079 | Source 007 adds six stable types while Source 006 defers separate registries and policy language | scope tension | open |
| PZL-083 | A universal contract envelope can become boilerplate when fields do not apply to a type | schema tax | open |
| PZL-084 | Naming every missing field may recreate the question/research/spike databases being collapsed | taxonomy recurrence | open |
| PZL-085 | Product/engineering authority separation may hide architecture decisions with product consequences | authority boundary | open |
| PZL-086 | Repository facts, requirements, assumptions, and recommendations need enforceable provenance classes | schema dependency | open |
| PZL-087 + PZL-067 | Source 006 requires one human lock; Source 007 requires product lock plus engineering assurance | direct approval conflict | open |
| PZL-087 | Independent evidence report and final approval authority are not the same thing | authority ambiguity | open |
| PZL-088 | Graph transaction semantics require baseline identity, concurrency rules, and partial-reconciliation behavior | architecture dependency | open |
| PZL-089 | Objective DAG may duplicate milestones, tasks, wayfinding dependencies, and validation | model duplication | open |
| PZL-090 | Graph-first retrieval is only as safe as edge completeness and freshness | retrieval risk | open |
| PZL-090 | Semantic search can locate evidence but cannot establish authority | authority boundary | supported |
| PZL-091 | Generated run manifests risk becoming stale shadow copies of source contracts | projection drift | open |
| PZL-092 | Structural gates can be gamed through artifact presence without substantive behavior | enforcement risk | open |
| PZL-093 | Reconnaissance evidence can become performative for bounded familiar changes | ceremony risk | open |
| PZL-094 | Exact path fences conflict with plans that intentionally derive file edits after repository inspection | planning collision | open |
| PZL-095 | Evidence-first generalizes TDD more honestly, but routing by change class can become another rules engine | complexity trade-off | open |
| PZL-096 | New-hypothesis retries improve diagnosis; a fixed universal attempt count does not follow | policy granularity | open |
| PZL-097 | Completion-boundary field can preserve delivery honesty while still keeping request and release lifecycles separate | model boundary | open |
| PZL-098 + PZL-071 | Source 006 proposes one coding agent; Source 007 proposes a specialist run chain | direct execution conflict | open |
| PZL-098 + PZL-079 | A multi-agent run chain tensions with the explicit MVP deferral of generalized orchestration | scope contradiction | open |
| PZL-099 | Mission/run separation adds a second lifecycle while the proposal claims taxonomy simplification | complexity transfer | open |
| PZL-100 | Mission-local records may lose reusable fix knowledge, durable decisions, and fleet-wide queries | retention collision | open |
| PZL-101 + PZL-077 | Source 007's graph workspace is materially larger than Source 006's minimal Mission workspace | direct workspace conflict | open |
| PZL-101 | The proposed tree recreates several collections before the thin slice earns them | scaffold recurrence | open |
| PZL-102 | A disposable vertical slice must not silently become production architecture | promotion boundary | open |
| source-008 | High-density synthesis provides no author, citations, benchmarks, or observed project evidence | provenance gap | verified |
| PZL-103 | “Graph over loop” is a false dichotomy: DAGs model dependencies while attempts and retries require cycles or state transitions | model correction | verified |
| PZL-104 + PZL-079 | Default parallel fan-out conflicts with Source 006's explicit MVP deferral of multi-agent orchestration | scope contradiction | open |
| PZL-104 | Split/work/check/merge only pays when units are separable and semantic integration has an owner | applicability boundary | open |
| PZL-105 | Clean isolated context can omit shared safety, policy, interface, or dependency constraints | context omission risk | open |
| PZL-106 | Published GraphRAG gains are scoped to global sensemaking questions, not universal retrieval superiority | evidence overreach | verified |
| PZL-106 + PZL-079 | Graph/semantic retrieval tensions with Source 006's explicit first-release deferral | scope contradiction | open |
| PZL-107 | An ontology is accepted semantics, not an infallible world model; it can be incomplete, stale, or wrong | truth-model correction | verified |
| PZL-108 | Raw logs may expose secrets or personal data and may be too large or transient for committed Markdown | retention and safety risk | open |
| PZL-109 | A systemic persona layer can harden subjective or stale inference into apparent project authority | memory-authority risk | open |
| PZL-110 | Deterministic checks verify declared mechanical properties, not complete semantic correctness | assurance boundary | verified |
| PZL-112 | A prototype is uncertainty evidence, not an accepted contract or production-quality implementation | promotion boundary | verified |
| PZL-113 + PZL-045 | Source 008 favors Wayfinder-like dependency mapping while Source 004 recommends deleting wayfinding | direct subsystem conflict | open |
| PZL-113 | Mentioning Wayfinder as a category does not establish current Spectacular usage or value | evidence gap | verified |
| PZL-114 | Context-isolated agent review still shares model biases and cannot make confidence into evidence | independence boundary | open |
| PZL-114 + PZL-071 | Risk-triggered independent verification conflicts with a strict one-agent Mission model | execution conflict | open |
| PZL-115 | Shift-right review for a small diff cannot bypass early consent for sensitive or irreversible effects | safety boundary | verified |
| PZL-116 | Feeding unchanged failures back without a new hypothesis automates repetition rather than diagnosis | loop-control risk | verified |
| PZL-117 | pgvector is an extension, not native core PostgreSQL | infrastructure factual correction | verified |
| PZL-117 | Property graphs appear in PostgreSQL 19 documentation but not current stable PostgreSQL 18 | version-bound claim | verified |
| PZL-117 | `SKIP LOCKED` is a queue-like concurrency primitive with an inconsistent view, not a complete queue contract | infrastructure scope correction | verified |
| PZL-117 | PostgreSQL consolidation may suit a server stack but conflicts with a portable Markdown/Git-native core as a universal requirement | product-boundary collision | open |
| PZL-118 | Canonical resume state and raw transient run artifacts require different persistence and retention rules | state-authority boundary | open |
| source-009 | Source 009 substantially restates Source 008 and is clarification, not independent corroboration | duplicate provenance | verified |
| source-009 | The exact “about 5% model versus up to 48% harness” comparison is unattributed and varies by task and benchmark | numerical evidence overreach | verified |
| PZL-104 + source-009 | A context-isolated Skeptic can reduce anchoring but does not “hold no bias” | verifier claim correction | verified |
| PZL-106 + source-009 | Similarity retrieval lacks explicit topology but does not categorically fail every multi-hop question | retrieval claim correction | verified |
| PZL-107 + PZL-122 | RDF/OWL semantic modeling and closed-world constraint validation are different functions | technology-boundary correction | verified |
| PZL-109 | L3 persona can encode user preference but must not outrank project facts, explicit instructions, or current evidence | memory-authority risk | open |
| PZL-111 | Deep runtime modules and small atomic evidence records solve different problems; file size alone is not module quality | category correction | verified |
| PZL-112 | Interactive prototype code cannot be the definitive spec for security, failure, compatibility, operations, or maintainability | authority collision | verified |
| PZL-119 | Typed handoff packets can become stale shadow copies unless they reference rather than reproduce primary evidence | projection drift | open |
| PZL-120 | A manually maintained Mermaid map can become a second, stale authority over concept metadata | projection drift | open |
| PZL-121 | Harness effects can be large, but some workloads show model choice dominating harness choice | evaluation scope boundary | supported |
| PZL-122 | Semantic-web technology may be unjustified when frontmatter schemas and ordinary validators express the required constraints | complexity threshold | open |
| source-010 | Source 010 substantially restates Sources 008/009 and supplies no independent benchmark or implementation evidence | duplicate provenance | verified |
| source-010 | The referenced interactive visualization was not included with the source | source completeness gap | verified |
| source-010 | “300-token handoff,” AC/DC, and Gauntlet are source-specific labels rather than verified universal standards | provenance boundary | verified |
| PZL-123 | An abstraction can reduce context while also hiding a safety-critical constraint or stale authority | abstraction leak risk | open |
| PZL-124 | A semantic adapter can become a leaky query language or stale duplicate of storage semantics | adapter drift | open |
| PZL-124 + PZL-077 | A database-oriented semantic substrate conflicts with the proposed portable Markdown/Git-native minimum as a universal feature | product-boundary collision | open |
| PZL-125 | Encapsulation limits blast radius but does not make messy internal code acceptable | quality correction | verified |
| PZL-125 | Human-designed interfaces can themselves encode premature or incorrect architecture | authority limitation | open |
| PZL-126 | A universal 300-token limit may omit baseline, negative evidence, security constraints, or unresolved failures | arbitrary budget risk | verified |
| PZL-126 | Failed trial chronology should be omitted only when it does not affect the next hypothesis or prevent repeated work | evidence-retention boundary | open |
| PZL-127 | “Until the Skeptic is entirely satisfied” creates nontermination and makes model confidence an authority | termination defect | verified |
| PZL-127 | A centralized reviewer can become a latency, cost, and correlated-failure bottleneck | orchestration bottleneck | open |
| source-011 | The open-issue corpus is a time-bound owner backlog, not independent evidence or blanket implementation authority | provenance boundary | verified |
| PZL-128 + PZL-129 | Progressive-loading measurement and soft-DB briefing design can create duplicate schemas and benchmarks | authority duplication | open |
| PZL-130 + PZL-079 | Semantic retrieval is deliberately downstream of measured deterministic retrieval and outside the first-release fence | sequencing conflict | open |
| PZL-134 + PZL-015 | Mission is proposed both as the bounded request replacement and as a portfolio layer above requests | work-unit collision | open |
| PZL-139 + PZL-071 | Executable graph concurrency conflicts with the single-agent checkpointed MVP | execution conflict | open |
| PZL-132 + PZL-135 | A generic renderer can proceed, but a live graph view depends on unsettled Mission and node contracts | dependency boundary | open |
| PZL-136 | Goal, change summary, and artifact title serve different reader questions and should not be collapsed into one field | semantic boundary | verified |
| PZL-137 + PZL-044 | Commit discipline is a valid run invariant whose implementation owner depends on whether AFK remains in Spectacular | ownership collision | open |
| PZL-141 + PZL-006 | Default GitHub scaffolding conflicts with minimal earned structure | initialization collision | open |
| PZL-142 | Commit enforcement needs a chosen boundary: commits, PR title, merge commit, or release tooling | enforcement ambiguity | open |
| PZL-143 + PZL-010 | Always-loaded generic principles consume context while mostly duplicating model knowledge | context-budget collision | open |
| PZL-145 | Treating `proposed` as a decision state can make unresolved proposals appear authoritative; QUE is existing prior art | lifecycle ambiguity | open |
| PZL-146 | Agent provenance does not establish human decision authority | authority boundary | verified |
| PZL-147 + PZL-038 | Debt scoring was previously rejected from core and has no usage evidence to overcome the protected-core threshold | survival evidence | verified |
| PZL-148 + PZL-043 | Current Vision usage friction is evidence for simplifying entry and also evidence against calling Vision dead | subsystem evidence | verified |
| PZL-149 + PZL-054 | Generic capture may create a local proposal, but network mutations remain explicit native-provider actions | side-effect boundary | open |
| source-012 | Exact claims of 10–20 decisions or 20+ auto-resolved defects per pass are unsourced and not independent evidence | numerical evidence gap | verified |
| source-012 | “Never use plain text” conflicts with its own text-based grooming interview and with abstract authority/policy decisions | internal contradiction | verified |
| PZL-150 + PZL-010 | Maximizing raw decision count can reduce correctness, retention, and context quality | measurement collision | open |
| PZL-151 | Requiring exactly three variants can manufacture weak options or triple prototype cost without adding information | arbitrary cardinality | verified |
| PZL-152 | One-question-per-turn traversal and 10–15-decisions-per-pass are distinct interaction modes, not one universal protocol | interaction contradiction | verified |
| PZL-153 | A simplified logic harness can omit the persistence, concurrency, or provider behavior that determines the real decision | fidelity risk | open |
| PZL-154 + PZL-127 | “Until a high bar” is not a terminal state; critic repair needs a bounded acceptance and escalation contract | termination defect | verified |
| PZL-154 + PZL-115 | Evaluators may reject defects but cannot silently resolve human-owned product trade-offs | authority boundary | verified |
| PZL-155 | Higher-fidelity artifacts can cost more than the uncertainty they resolve; cheapest sufficient evidence remains the routing rule | proportionality boundary | open |
| PZL-156 + PZL-041 | Five command-shaped capabilities do not establish five standalone skills and may recreate fleet/surface bloat | responsibility collision | open |
| source-013 | Source 013 substantially extends Source 012's decision-density proposal and is clarification, not independent empirical corroboration | duplicate provenance | verified |
| source-013 | Four-to-six variants, five/ten personas, and 10–20 resolved questions are arbitrary unsourced counts | numerical evidence gap | verified |
| PZL-157 | A tension matrix exposes trade-offs but cannot prove that a generated reconciliation is usable or ethically acceptable | decision-authority boundary | verified |
| PZL-158 | Parametric prototypes can create combinatorial noise and false fidelity when axes do not model real conditions | prototype-fidelity risk | open |
| PZL-159 | Simulated personas are model-generated hypotheses, not human research, accessibility evidence, emotion, or trust measurement | evidence-authority boundary | verified |
| PZL-159 | Cognitive Friction Score and Time-to-Trust lack operational definitions and would be pseudo-precision if generated by an LLM | metric validity gap | verified |
| PZL-160 | Consequence analysis must remain proportional or every reversible interaction inherits high-risk ceremony | proportionality boundary | open |
| PZL-161 + PZL-115 | Blocking intent previews on trivial actions can recreate micro-approval fatigue | interaction-cost collision | open |
| PZL-162 | Maintaining separate AI and manual paths can duplicate business rules unless both share one underlying contract | implementation-drift risk | open |
| PZL-163 | Exposing raw tokens, embeddings, or hidden retrieval internals may confuse users without providing meaningful control | abstraction-boundary risk | open |
| PZL-164 | Silent escalation from assistive to autonomous behavior violates explicit authority and permission continuity | authority collision | verified |
| PZL-165 + PZL-156 | A deep decision companion may unify the workflow, but must not become a universal orchestrator or five shallow skills behind one name | product-boundary risk | open |
| PZL-166 | AI UX stress testing may be a domain profile, pack, or standalone skill; packaging depends on distinct substrate, tools, and users | extraction ambiguity | open |
| source-014 | Type 1/Type 2, 30-minute decisions, one-day spikes, and five-day reviews are heuristics with arbitrary thresholds, not validated universal rules | calibration boundary | verified |
| source-014 | A formal PRD/ADR pipeline cannot guarantee alignment; it can only make authority, disagreement, and evidence more inspectable | process overclaim | verified |
| PZL-167 | Reversibility is contextual and changes when consumers, stored data, operations, regulation, or adoption make migration harder | classification drift | open |
| PZL-167 + PZL-115 | A nominally reversible choice still requires early authority when it touches security, privacy, rights, or irreversible side effects | risk-classification collision | open |
| PZL-168 + PZL-067 | One decision owner does not imply one universal approval; product, engineering, security, provider, and irreversible authorities may remain separate | authority collision | open |
| PZL-168 | “Disagree and commit” becomes coercive when the owner lacks delegated authority, omits affected advisers, or offers no escalation path | governance boundary | open |
| PZL-169 | Combining product intent, UX, schemas, dependencies, KPIs, SLOs, and architecture in one PRD creates a mixed authority and retrieval burden | document-boundary risk | open |
| PZL-170 + PZL-111 | Abstracting every third party for hypothetical replacement conflicts with earned structure and can manufacture shallow indirection | premature-abstraction risk | open |
| PZL-170 | “Monolith first” and “boring technology” are defaults, not substitutes for repository constraints, failure analysis, or measured fit | heuristic boundary | verified |
| PZL-171 + PZL-110 | Schema diffs detect declared structural compatibility, not every semantic, behavioral, security, or operational break | assurance boundary | verified |
| PZL-171 | Generated mocks become a second truth if they are not reproducibly derived from the accepted interface schema | projection-drift risk | open |

Add cross-source contradictions here as soon as they are observed. A row records
the collision; it does not choose the winner.
