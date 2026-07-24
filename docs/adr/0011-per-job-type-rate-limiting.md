# ADR 0011: Per-job-type rate limiting via token buckets

## Status
Accepted

## Context
Worker pool concurrency bounds how many jobs run in parallel, but says
nothing about the rate at which a specific job type hits a downstream
dependency (e.g. an email provider's API limit). A burst of queued
jobs of one type could still overwhelm a rate-limited external service
even with modest worker concurrency.

## Decision
Add a ratelimit package wrapping golang.org/x/time/rate, keyed by job
type. The worker pool calls Wait(ctx, job.Type) before invoking a
handler; job types with no configured limit pass through unrestricted.

## Consequences
- Rate limits are independent of worker concurrency — you can run 20
  workers but still cap send_email at 5/sec if that's the provider's
  actual limit.
- Limits are configured once at startup (SetLimit), not dynamically
  adjustable at runtime without a restart — acceptable for this
  project; a production system might expose this via the admin API
  (planned) for live tuning.
- If ctx is cancelled while waiting for a token (typically during
  shutdown), the job is re-enqueued rather than dropped or treated as
  a failed attempt, since the handler never actually ran.
- Rate limiting is per-process, not cluster-wide: running multiple
  worker nodes (ADR 0006) means each node's limiter operates
  independently, so the *effective* cluster-wide rate for a job type
  is roughly (configured rate) × (number of worker nodes). A true
  cluster-wide limit would need a Redis-backed token bucket instead of
  an in-memory one — noted as a gap, not implemented, since it adds
  real complexity for a benefit this project doesn't currently need.
