# Architecture Decision Records

Records of significant architectural decisions made in this project,
in Context → Decision → Consequences format.

| ADR | Title |
|-----|-------|
| [0001](0001-redis-list-over-streams-for-queue.md) | Redis list over Streams for the job queue |
| [0002](0002-fixed-size-worker-pool.md) | Fixed-size worker pool |
| [0003](0003-durable-delayed-jobs-via-sorted-set.md) | Durable delayed/retry jobs via sorted set |
| [0004](0004-separate-dead-letter-queue.md) | Separate dead-letter queue |
| [0005](0005-graceful-shutdown-with-bounded-timeout.md) | Graceful shutdown with bounded timeout |
| [0006](0006-multi-node-workers-require-no-code-changes.md) | Multi-node workers require no code changes |
| [0007](0007-dual-write-redis-and-postgres.md) | Dual-write to Redis and Postgres |

New ADRs should follow [the template](0000-template.md) and be numbered
sequentially.
