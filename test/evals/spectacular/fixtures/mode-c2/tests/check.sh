#!/bin/sh
set -eu
grep -q '^COMPONENT_A=OK$' src/a.txt
grep -q '^INTERFACE=v1$' src/a.txt
grep -q '^COMPONENT_B=OK$' src/b.txt
