# Execute Micro-Kernel

## 1. Trigger Context
Orchestrator activating/resuming a Mission or Worker executing an assigned code charter.

## 2. CLI Palette & Execution
```bash
git checkout -b <mission-slug>                  # Always branch before activation
spectacular mission start plan.md --json        # Activate execution envelope
spectacular mission check <ref> --json          # Read-only validation
spectacular mission show <ref> --json           # Inspect active objective
```

## 3. Negative Constraints (DO NOT)
- **DO NOT** execute on `main` branch. Always create a mission branch first.
- **DO NOT** output meta-planning narrative. Directly write code, run tests, and return results.
- **DO NOT** manage Spectacular files inside worker subagents (workers ignore `.spectacular/`).
- **DO NOT** touch files outside `allowed_changed_paths` defined in the charter.

## 4. Greenfield & Concurrency Invariants
When creating standalone tools, services, or workers:
- **Worker Pools**: Bind concurrency strictly to `--workers N` using bounded pools/semaphores.
- **Retry & DLQ**: Track attempts per item. Route to `dlq.json` only after exceeding failure threshold ($\ge 3$).
- **Clean Exit**: Drain in-flight jobs, close channels/sockets, and exit with status code 0.

```python
# Concrete Idiom: Worker Pool + DLQ (Python)
from concurrent.futures import ThreadPoolExecutor
import argparse, json, sys

parser = argparse.ArgumentParser()
parser.add_argument("--workers", type=int, default=4)
args = parser.parse_args()

def process(job, max_retries=3):
    for attempt in range(1, max_retries + 1):
        job["attempts"] = attempt
        if not job.get("fail", False):
            return True
    return False

with open("jobs.json") as f: jobs = json.load(f)
with ThreadPoolExecutor(max_workers=args.workers) as pool:
    results = list(pool.map(process, jobs))
dlq = [jobs[i] for i, ok in enumerate(results) if not ok]
with open("dlq.json", "w") as f: json.dump(dlq, f)
```

## 5. Authority & Context Invariants
- **Progressive Context**: Drill down strictly: `Mission card -> current Objective -> exact sources`.
- **Authority Separation**: The `plan carries meaning`, while `tooling carries repeatability`.
- **Activation Boundary**: `A Decision is not activation authority` (only owner confirmation authorizes `mission start`).
- **Self-Hosting & Bootstrap**: When developing Spectacular, an `active Mission keeps the schema` frozen. Under declared `manual-bootstrap`, run `focused checks` directly.
