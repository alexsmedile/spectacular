# Scenario A — draft a PLAN (autopilot)

You are operating the Spectacular skill for a project. The user says:

> "spectacular new — I want export to PDF. Our users keep asking to share their
> dashboards outside the app. Add a way to export a dashboard to PDF. It should
> be fast and look good. Probably a button on the dashboard toolbar."

Project context you have available (pretend you read these):
- The project PRD's Goals section contains: "G2 — A dashboard owner can share
  a snapshot of their dashboard with a non-user in under 30 seconds" and
  "G3 — p95 page interaction latency stays under 200ms".
- STACK.md says: TypeScript/React front end, Node API, no headless-browser
  dependency currently installed. Two plausible export approaches exist:
  server-side headless Chrome rendering, or client-side canvas capture
  (html2canvas-style). Each has real trade-offs (fidelity vs new heavy dep).
- The request summary you would stamp is: "Add a way to export a dashboard to
  PDF" (this is what goes in PLAN frontmatter `summary:`).

Produce the complete PLAN.md you would write for this request, plus whatever
else your instructions tell you to present alongside it. Follow ONLY the
reference docs you were given. Output the PLAN.md content verbatim in a fenced
block, then anything else your flow requires you to show.
