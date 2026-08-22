#!/bin/sh
set -eu
grep -q '^AUTH_RESULT=OK$' src/auth.txt
grep -q '^EMPTY_TOKEN=REJECT$' src/auth.txt
