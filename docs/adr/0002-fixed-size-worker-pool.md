# ADR 0002: Fixed-size worker pool instead of one goroutine per job

## Status
Accepted

## Context
Jobs need to be processed concurrently, but processing many jobs at once
can overwhelm downstream dependencies (databases, external APIs) if
concurrency is unbounded.

## Decision
Run a fixed number of long-lived goroutines (`concurrency`), each looping
and pulling one job at a time, rather than spawning `go handler(job)` per
dequeued job.

## Consequences
- Concurrency is capped and predictable — a burst of 10,000 queued jobs
  never spawns 10,000 goroutines.
- Requires explicit shutdown coordination (context + WaitGroup, see ADR
  0005) since goroutines are long-lived, not fire-and-forget.
- Throughput is bounded by `concurrency`, meaning under sustained high
  load workers may lag behind producers — acceptable for this project;
  a production system would pair this with horizontal scaling (more
  worker processes, see ADR 0006) rather than raising concurrency
  unboundedly on one process.
