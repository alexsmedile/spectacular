#!/bin/sh
set -eu
grep -q '^VALUE=OK$' src/value.txt
