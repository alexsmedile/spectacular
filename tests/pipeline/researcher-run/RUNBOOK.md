# P3 — Researcher actually runs

**Proves:** the `debug-researcher` agent runs end-to-end — takes an external-smelling bug, searches
(scrapekit or harness fallback), judges relevancy, and returns a `VERDICT` block with citations +
writes `research/research-NN.json`. This agent has **never been run once**; P3 is its first live test.

**No code fixture** — the "bug" is a real, well-documented external issue so we can check the verdict
against known ground truth. Use a genuinely researchable symptom, e.g.:

> "Python `json.dumps` raises `TypeError: Object of type datetime is not JSON serializable`."

This is a famous, documented stdlib behavior with a known cause + workaround (`default=str` / a custom
encoder). A correct Researcher will find it fast and return `known-platform-bug` (well, known-stdlib-
behavior) with a documented workaround — a clean, checkable result.

---

## Steps (orchestrator drives)

1. **Frame the external question.** Write a research brief: the symptom in others' vocabulary
   (the exact error string + the stdlib surface), and the question: "is this a known
   library/platform behavior with a documented cause/workaround?"

2. **Spawn `debug-researcher`** with that brief and the trace path
   `.spectacular/debug/<job>/research/research-01.json`.
   - **ASSERT:** it runs without a tool error (scrapekit *or* harness WebSearch/WebFetch — either is
     a pass; note which path it took).

3. **Check the returned block.**
   - **ASSERT:** `VERDICT` is from the enum (`known-platform-bug | genuinely-ours | no-strong-match`);
     for this fixture it should be `known-platform-bug`.
   - **ASSERT:** `CONFIDENCE` present; `EVIDENCE` has ≥1 entry with a real URL and a match-reason.
   - **ASSERT:** `WORKAROUND` names the documented fix (a `default=` handler / custom encoder), not "none found".
   - **ASSERT:** `NEXT` is actionable (apply the documented workaround).

4. **Check the trace artifact.**
   - **ASSERT:** `research/research-01.json` exists, parses, and mirrors the block (same verdict + evidence).

5. **Skeptic check (relevancy judging is the skill).** Read the EVIDENCE URLs' match-reasons.
   - **ASSERT:** each cited source is actually about *this* error, not a superficial keyword match —
     the Researcher's protocol step 4 ("judge relevancy hard") held.

---

## Pass criteria
The Researcher ran, returned a well-formed `known-platform-bug` verdict with a real citation and the
documented workaround, and wrote a matching `research-01.json`. Note which fetch path (scrapekit vs
harness) it used.

## Results (fill after running)
- Date: 2026-07-05
- Verdict: **PASS** (all 5 assertions) — `debug-researcher`'s first-ever live run.
- Fetch path used: **harness WebSearch + WebFetch** (scrapekit not needed — sources returned clean full content; the agent noted it fell back correctly by choice, not failure).
- Notes:
  - Ran with no tool error. `VERDICT: known-platform-bug` (correct for this stdlib behavior), `CONFIDENCE: high`.
  - **Relevancy judging (the skill) held:** 3 EVIDENCE entries, each with a real URL + a specific match-reason. Primary citation is the **authoritative stdlib doc** (docs.python.org/3/library/json.html), which documents both the cause (datetime absent from the encodable-types table, by design) AND the workaround — the spec itself, not a keyword match. The two community sources (bobbyhadz exact-error-string, PyNative) corroborate, correctly labeled as corroboration not speculation.
  - `WORKAROUND` names all three documented fixes (`default=str`, a `default=` isoformat handler, `JSONEncoder` subclass). `NEXT` actionable: apply at call site, no pin/upgrade needed. It even noted "No Investigator handoff required" — correctly reading its own verdict's implication.
  - Trace `research/research-01.json` parses and mirrors the block (same verdict, 3 evidence, workaround).
  - Not logged to `fixes/` — documented stdlib behavior teaches nothing a future search wouldn't re-find; correct curation restraint.
