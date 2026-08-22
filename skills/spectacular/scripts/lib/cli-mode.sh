#!/bin/sh
# Shared read-only compatibility probe. Sets SPECTACULAR_CLI_MODE,
# SPECTACULAR_CLI_VERSION, and SPECTACULAR_CLI_DETAIL.

spectacular_cli_probe() {
  catalog=$1
  expected=$(sed -n 's/^[[:space:]]*"release_version": "\([^"]*\)",*$/\1/p' "$catalog" | head -1)
  SPECTACULAR_CLI_MODE=reduced
  SPECTACULAR_CLI_VERSION=unknown
  SPECTACULAR_CLI_DETAIL="ABSENT — governed execution unavailable"

  if ! command -v spectacular >/dev/null 2>&1; then
    return
  fi

  output=$(spectacular --version --json 2>/dev/null) || {
    SPECTACULAR_CLI_DETAIL="UNREADABLE — governed execution unavailable"
    return
  }
  schema=$(printf '%s\n' "$output" | sed -n 's/.*"schema_version":"\([^"]*\)".*/\1/p')
  actual=$(printf '%s\n' "$output" | sed -n 's/.*"version":"\([^"]*\)".*/\1/p')
  SPECTACULAR_CLI_VERSION=${actual:-unknown}

  if [ "$schema" != "spectacular.build-info.v1" ] || [ -z "$expected" ] || [ "$actual" != "$expected" ]; then
    SPECTACULAR_CLI_DETAIL="INCOMPATIBLE — found ${actual:-unknown}, need ${expected:-unknown}; governed execution unavailable"
    return
  fi

  SPECTACULAR_CLI_MODE=full
  SPECTACULAR_CLI_DETAIL="compatible — spectacular $actual"
}
