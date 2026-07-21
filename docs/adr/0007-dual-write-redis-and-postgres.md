# ADR 0007: Dual-write to Redis and Postgres for queue state vs. history

## Status
Accepted

## Context
Redis holds live queue state ("what needs to run next"), but has no
built-in way to query historical job outcomes. A durable, queryable
audit trail needs a different kind of storage.

## Decision
Add a Postgres `job_history` table, written to alongside Redis at each
lifecycle transition (created, completed, failed, dead-lettered), via a
separate `store.Store` type distinct from `queue.Queue`.

## Consequences
- SQL queries like "all failed jobs of type X in the last hour" become
  trivial, versus awkward or impossible against Redis alone.
- Two separate write calls per transition means a real, accepted risk:
  if one write succeeds and the other fails, Redis and Postgres can
  disagree about a job's true state. No outbox pattern or transactional
  write is implemented to close this gap.
- Store and Queue are kept as separate types with a clear split of
  responsibility (Queue = what's next, Store = what happened), so
  either storage system could be swapped independently later without
  the other needing to change.
