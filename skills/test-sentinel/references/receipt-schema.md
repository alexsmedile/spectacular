# Test Sentinel Receipt Specification (`test-sentinel.receipt.v1`)

To enable autonomous agent verification and mechanical integration with Spectacular or external CI systems, Test Sentinel standardizes on a machine-readable JSON receipt.

---

## 1. Schema Definition

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "TestSentinelReceipt",
  "type": "object",
  "required": [
    "schema_version",
    "status",
    "tier",
    "command",
    "duration_ms",
    "commit"
  ],
  "properties": {
    "schema_version": {
      "type": "string",
      "const": "test-sentinel.receipt.v1"
    },
    "status": {
      "type": "string",
      "enum": ["pass", "fail"]
    },
    "tier": {
      "type": "string",
      "enum": ["tier0-preflight", "tier1-unit", "tier2-hardened", "tier3-acceptance", "custom"]
    },
    "command": {
      "type": "string",
      "description": "Exact shell command executed"
    },
    "duration_ms": {
      "type": "integer",
      "description": "Execution time in milliseconds"
    },
    "commit": {
      "type": "string",
      "description": "Git commit SHA tested against"
    },
    "flakes_detected": {
      "type": "integer",
      "description": "Number of non-deterministic run failures observed during repeat/stress checks (e.g. -count=N). Omitted on single-run passes."
    },
    "failures": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["test", "error"],
        "properties": {
          "test": { "type": "string" },
          "error": { "type": "string" },
          "file": { "type": "string" },
          "line": { "type": "integer" }
        }
      }
    }
  }
}
```

---

## 2. Examples

### Successful Single-Run Receipt
```json
{
  "schema_version": "test-sentinel.receipt.v1",
  "status": "pass",
  "tier": "tier1-unit",
  "command": "go test ./internal/auth/...",
  "duration_ms": 420,
  "commit": "3a4d2c270e763c5c05ef95b132e70fd6552090de",
  "failures": []
}
```

### Successful Concurrency Stress Receipt (with repeat checks)
```json
{
  "schema_version": "test-sentinel.receipt.v1",
  "status": "pass",
  "tier": "tier2-hardened",
  "command": "go test -race -count=20 ./internal/queue/...",
  "duration_ms": 3850,
  "commit": "3a4d2c270e763c5c05ef95b132e70fd6552090de",
  "flakes_detected": 0,
  "failures": []
}
```

### Failed Receipt with Diagnostics
```json
{
  "schema_version": "test-sentinel.receipt.v1",
  "status": "fail",
  "tier": "tier2-hardened",
  "command": "go test -race ./internal/auth/...",
  "duration_ms": 1120,
  "commit": "3a4d2c270e763c5c05ef95b132e70fd6552090de",
  "failures": [
    {
      "test": "TestRegression_TokenRefreshRace",
      "error": "DATA RACE: Write at 0x00c0000a6040 by goroutine 7: concurrent map read and map write",
      "file": "internal/auth/token_test.go",
      "line": 42
    }
  ]
}
```
