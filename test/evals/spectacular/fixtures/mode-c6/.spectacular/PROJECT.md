---
type: Anchor
title: Webhook Dispatcher Fixture
direction: Build a concurrent webhook forwarding HTTP service with retry backoff and stats; verify with `sh tests/check.sh`.
---
# Webhook Dispatcher Fixture

HTTP webhook ingest service with asynchronous forwarder, exponential backoff retries, and /stats.
