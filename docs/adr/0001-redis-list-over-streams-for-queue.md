# ADR 0001: Use a Redis list (LPUSH/BRPOP) instead of Redis Streams for the job queue

## Status
Accepted

## Context
The pending job queue needs a way to push work in and have workers pull it
out, with blocking semantics so workers don't busy-poll an empty queue.
Redis offers multiple primitives for this: a plain list with BRPOP, or
Streams with consumer groups.

## Decision
Use a Redis list with LPUSH (producer) and BRPOP (consumer).

## Consequences
- Simple mental model: FIFO list, blocking pop, nothing else to configure.
- BRPOP already gives safe multi-consumer behavior — multiple worker
  processes can BRPOP the same list and Redis guarantees each item goes
  to exactly one consumer, which is what made multi-node workers (ADR
  0006) work with zero additional code.
- Trade-off: no consumer groups, no replay, no per-consumer acknowledgment
  tracking. If a worker crashes mid-job after BRPOP has already removed
  the item, that job is gone unless other durability measures catch it
  (see ADR 0004 on delayed/retry jobs, which only covers *failed* jobs,
  not crashed-mid-processing ones).
- If this project ever needed replay, audit of exactly what was
  delivered to which consumer, or true exactly-once processing
  guarantees, migrating to Streams would be the natural next step —
  deferred deliberately since the added complexity wasn't justified
  for a learning project's core queue.
