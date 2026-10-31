# Failure patterns in generated and rapidly assembled code

Use this when: conducting a whole-codebase or vibe coding audit. These are investigation
prompts, not evidence that a project has a defect. Prioritize by actual usage and impact.

| Pattern | Probe and evidence needed |
|---|---|
| Disconnected feature | Trace a supported entrypoint to implementation and output consumer; distinguish planned work from advertised runtime behavior. |
| Dead code / unused computation | Follow symbol reads and dynamic dispatch; demonstrate an unused local result separately from module reachability. |
| Copy-and-modify drift | Compare AST/body clones and near-clones, including diverging validation; find the proper shared owner before extracting. |
| Test-shaped reassurance | Follow assertions to the promised outcome; use the test-credibility reference for counterexamples. |
| Happy-path-only validation | Try missing, malformed, empty, zero, negative, nonfinite, boundary and contradictory inputs at a real entrypoint. |
| Silent fallback / swallowed errors | Supply an explicitly invalid path or provider; inspect exit status and receipts for lost failure attribution. Optional absence differs from explicit misconfiguration. |
| Schema theater | Verify fields are consumed, unknown controls rejected when required, and Python/JS/types agree on meaning and units. |
| Claims exceeding evidence | Compare names/docs/status with what is actually measured; metadata and transcripts do not establish visual content or perceptual quality. |
| Stale caches / evidence | Change source bytes, parameters, model/version or revision while keeping an ID; verify output provenance and invalidation. |
| Broken state / partial writes | Fail between writes, retry after interruption and test concurrent writers; inspect atomicity, locks and recovery. |
| Duplicate external effects | Simulate timeout after remote success in an isolated adapter fixture; check idempotency and reconciliation before retry. |
| Paths / portability / security | Test traversal, symlinks, moved roots, offline optional providers and untrusted strings at real trust boundaries; never print secrets. |
| Wrong clocks / units / identities | Follow timebase, fps, duration, coordinates, asset IDs and hash bindings across adapters and serialization. |
| Unnecessary abstractions / competing owners | Identify real readers and writers; demonstrate behavioral divergence or maintenance cost rather than criticizing file count. |
| Dependency / environment illusion | Check declared dependencies, lockfiles, clean-import assumptions and generated sources without implicit installation. |
| Scope and legacy confusion | Check current requirements before marking deferred UI, preserved compatibility, experimental tools or purchased sources as defects. |
| Fail-open controls | Distinguish required missing config from optional absence; trace fallback values to the actual authorization, acceptance or write decision. |
| Boundary asymmetry | Trace producer IDs/paths/units through adapters to filesystem, HTML, network or subprocess sinks; a safe producer does not validate every supported direct caller. |
| Input/output aliasing | Test output equal to or containing an input, partial directory overlap, same file through symlinks and missing directory components before destructive writes. |
| API or dependency hallucination | Compare imports, method signatures and flags with the pinned, installed runtime; declarations/source availability alone do not prove the integration executes. |

Preserve documented single-writer and supported legacy constraints when refuting
candidates; report operational limitations separately from contract violations.
