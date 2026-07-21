# ADR 0006: Multiple worker processes require no additional code

## Status
Accepted

## Context
Running this queue across multiple machines/processes (rather than one)
is a natural next step for a "distributed" job queue, and it's worth
being explicit about what does and doesn't need to change to support it.

## Decision
No sharding or leader election is implemented. Multiple `cmd/worker`
processes can already run against the same Redis instance safely,
because BRPOP against a shared list (ADR 0001) guarantees each job goes
to exactly one blocked consumer. Each worker is tagged with a `nodeID`
(env var or hostname) purely for log attribution across machines.

## Consequences
- Horizontal scaling of job processing is already possible with zero
  additional code — just run more `cmd/worker` instances.
- The scheduler (ADR 0003) remains a single point of failure for
  promoting delayed jobs; running multiple scheduler instances is safe
  (no duplicate promotion risk beyond what's noted in ADR 0003) but
  wasteful, since they'd race on the same 1s poll with only one doing
  useful work per tick.
- True leader election (e.g. a Redis `SET NX EX` lock so exactly one
  scheduler is active) and queue sharding (splitting jobs across
  multiple Redis instances) are known, well-understood techniques that
  were deliberately not implemented — they'd mostly exercise
  distributed-systems boilerplate rather than teach something new about
  this project specifically.
