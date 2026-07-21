# ADR 0003: Durable delayed/retry jobs via a Redis sorted set, not an in-memory timer

## Status
Accepted (supersedes an earlier in-process approach)

## Context
Failed jobs need to retry after a backoff delay. The first implementation
used a goroutine with `time.After(delay)` to sleep, then re-enqueue. This
worked, but the retry existed only in that worker process's memory — if
the process crashed or restarted during the delay, the retry was lost
with no trace anywhere.

## Decision
Store delayed/retry jobs in a Redis sorted set (score = Unix run-at
timestamp). A separate `cmd/scheduler` process polls the set every
second and promotes due jobs into the pending queue.

## Consequences
- Retries survive worker crashes and restarts — state lives in Redis,
  not in a goroutine.
- Introduces a new failure mode to be aware of: the scheduler itself is
  now a dependency for retries to ever fire. If the scheduler is down,
  jobs already in the pending queue keep processing fine, but nothing
  currently in the delayed set will be promoted until it's back.
- The scheduler is a separate process by design (see ADR 0006) — job
  execution and job promotion are different responsibilities and scale
  independently.
- 1-second polling granularity is an accepted precision trade-off; a job
  scheduled for exactly `now` may sit up to ~1s before being promoted.
- `PromoteDueJobs` does an LPUSH then a ZREM as two separate Redis calls,
  not wrapped in a Lua script. A crash between the two calls could
  theoretically cause a job to be promoted but not removed from the
  delayed set, leading to a duplicate promotion later. Known and
  accepted for now; a Lua script would make this atomic if it ever
  becomes a real problem.
