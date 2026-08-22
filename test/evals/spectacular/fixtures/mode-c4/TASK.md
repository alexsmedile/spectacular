Resume this interrupted authentication repair from the workspace state.

- Change only `src/auth.txt`.
- Replace `AUTH_RESULT=BUG` with `AUTH_RESULT=OK`.
- Preserve `EMPTY_TOKEN=REJECT`.
- Run `sh tests/check.sh` and report the verified result.
