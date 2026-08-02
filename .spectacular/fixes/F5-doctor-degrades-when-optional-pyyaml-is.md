---
type: fix
opened: 2026-08-01
verified: 2026-08-01
severity: medium
from_audit: null
debug_job: null
signature: doctor config.yaml invalid YAML ModuleNotFoundError no module named yaml optional PyYAML
related: []
---

# F5 — Doctor degrades when optional PyYAML is unavailable

## Problem
doctor treated ModuleNotFoundError as invalid config YAML, causing false workspace errors and blocking archive checks

## Intended behavior
When PyYAML is absent, doctor falls back to its structural project stanza check

## Root cause
the Python subprocess exit path did not distinguish missing optional dependency from yaml.YAMLError

## Fix
reserve exit 3 for missing PyYAML and keep yaml_ok true so the documented fallback runs

## Success criteria
doctor.test.sh and mutator.test.sh pass without PyYAML installed

## Verified by
bash tests/cli/doctor.test.sh; bash tests/cli/mutator.test.sh; bash tests/run.sh

## Signature
doctor config.yaml invalid YAML ModuleNotFoundError no module named yaml optional PyYAML
