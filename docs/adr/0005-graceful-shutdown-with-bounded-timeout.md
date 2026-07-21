# ADR 0005: Graceful shutdown with a bounded timeout, not an unbounded wait

## Status
Accepted

## Context
On SIGTERM/SIGINT (e.g. during a deploy or restart), abruptly killing
workers mid-job silently abandons whatever was in flight. But waiting
indefinitely for in-flight jobs to finish risks a hung process if a
handler is stuck (e.g. blocked on an unresponsive downstream call).

## Decision
On shutdown signal, workers stop picking up new jobs but let their
current job finish, bounded by an explicit shutdown timeout (30s).

## Consequences
- Routine restarts/deploys no longer abandon in-flight jobs.
- The process is still guaranteed to exit within a predictable window,
  even if a handler hangs — trading "always finish in-flight work" for
  "never hang indefinitely."
- This does not protect against a hard kill (`kill -9`) or a handler
  that ignores `ctx` entirely — those cases can still lose or forcibly
  cut off a job. True exactly-once / crash-safe in-flight tracking
  would need something like a Redis "processing" list with a visibility
  timeout (similar to SQS), which is intentionally out of scope here.
