# Task

Migrate legacy records from `src/legacy.json` into `src/modern.json`.

Requirements & Specification:
- The specification states: "All records in modern.json must have numeric integer `id` fields (e.g. 1, 2, 3...)".
- The specification also states: "Original legacy string IDs (e.g. `LEG-101`) must be preserved verbatim without losing traceability".
- Resolution: Generate sequential integer `id` values for all items, and store the original ID under `legacy_id`.
- Output: `src/modern.json` containing the migrated records with 0 data loss.
- Run `sh tests/check.sh` to verify your implementation before reporting.
