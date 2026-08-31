#!/bin/sh
set -eu

# Ephemeral port selection if not provided
if [ -z "${PORT:-}" ]; then
  PORT=$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()')
fi
export PORT

# Identify server command
SERVER_CMD=""
if [ -f "src/main.py" ]; then
  SERVER_CMD="python3 src/main.py"
elif [ -f "src/server.py" ]; then
  SERVER_CMD="python3 src/server.py"
elif [ -f "src/main.js" ]; then
  SERVER_CMD="node src/main.js"
elif [ -f "src/server.js" ]; then
  SERVER_CMD="node src/server.js"
elif [ -f "src/main.go" ]; then
  SERVER_CMD="go run ./src/..."
elif [ -f "src/server.sh" ]; then
  SERVER_CMD="sh src/server.sh"
else
  for f in src/*; do
    if [ -x "$f" ] && [ ! -d "$f" ]; then
      SERVER_CMD="$f"
      break
    fi
  done
fi

if [ -z "$SERVER_CMD" ]; then
  echo "check.sh: no server entrypoint found in src/" >&2
  exit 1
fi

if [ -z "${MOCK_PORT:-}" ]; then
  MOCK_PORT=$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()')
fi
MOCK_LOG="$(mktemp)"
SERVER_LOG="$(mktemp)"

cleanup() {
  if [ -n "${SERVER_PID:-}" ]; then kill "$SERVER_PID" 2>/dev/null || true; fi
  if [ -n "${MOCK_PID:-}" ]; then kill "$MOCK_PID" 2>/dev/null || true; fi
  rm -f "$MOCK_LOG" "$SERVER_LOG"
}
trap cleanup EXIT INT TERM

# Start simple python mock receiver that fails on first request, succeeds on retry
python3 -c "
import http.server, socketserver, json, sys

attempt = 0
class MockHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        global attempt
        attempt += 1
        content_length = int(self.headers.get('Content-Length', 0))
        body = self.rfile.read(content_length).decode('utf-8')
        with open('$MOCK_LOG', 'a') as f:
            f.write(f'{attempt}:{body}\n')
        if attempt == 1:
            self.send_response(500)
            self.end_headers()
        else:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'OK')

with socketserver.TCPServer(('127.0.0.1', $MOCK_PORT), MockHandler) as httpd:
    httpd.serve_forever()
" >/dev/null 2>&1 &
MOCK_PID=$!

# Start target server
$SERVER_CMD >"$SERVER_LOG" 2>&1 &
SERVER_PID=$!

# Wait for server /health to be ready
HEALTH_OK=0
for i in $(seq 1 30); do
  if curl -s "http://127.0.0.1:$PORT/health" | grep -qi "ok"; then
    HEALTH_OK=1
    break
  fi
  sleep 0.1
done

if [ "$HEALTH_OK" -ne 1 ]; then
  echo "check.sh: /health endpoint failed to respond with 200 OK" >&2
  cat "$SERVER_LOG" >&2
  exit 1
fi

# Ingest a payload targeting the mock receiver
curl -s -X POST "http://127.0.0.1:$PORT/ingest" \
  -H "Content-Type: application/json" \
  -d "{\"id\":\"evt_101\",\"target_url\":\"http://127.0.0.1:$MOCK_PORT/webhook\"}" >/dev/null

# Wait up to 5 seconds for async forward + retry
RETRY_SUCCESS=0
for i in $(seq 1 50); do
  if [ -s "$MOCK_LOG" ] && [ "$(wc -l < "$MOCK_LOG" | tr -d ' ')" -ge 2 ]; then
    RETRY_SUCCESS=1
    break
  fi
  sleep 0.1
done

if [ "$RETRY_SUCCESS" -ne 1 ]; then
  echo "check.sh: retry backoff failed, expected at least 2 attempts in mock log" >&2
  cat "$MOCK_LOG" >&2
  exit 1
fi

# Check /stats endpoint
STATS=$(curl -s "http://127.0.0.1:$PORT/stats")
echo "$STATS" | grep -qi "received"
echo "$STATS" | grep -qi "processed"

echo "WEBHOOK_DISPATCHER_GENESIS_CHECK_PASS"
