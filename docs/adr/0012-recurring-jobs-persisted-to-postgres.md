# ADR 0012: Recurring jobs defined via cron expressions, persisted to Postgres

## Status
Accepted

## Context
Some work needs to happen on a repeating schedule (e.g. a daily digest
email) rather than being triggered by a producer. Hardcoding schedules
in application code would mean any change requires a redeploy.

## Decision
Add a `recurring_jobs` Postgres table storing cron expressions and job
templates. A separate `cmd/cron` process loads enabled definitions at
startup, registers them with `robfig/cron`, and enqueues a fresh job
instance each time a schedule fires.

## Consequences
- Schedules survive restarts and are decoupled from code — a future
  admin API could let schedules be added/disabled without a deploy.
- `cmd/cron` currently loads definitions once at startup; adding a new
  recurring job requires restarting the cron process to pick it up.
  Not yet dynamic/hot-reloaded — acceptable for now, a known
  improvement if this becomes a real pain point.
- Like the scheduler in ADR 0003, cmd/cron is a single point of failure
  for recurring jobs specifically — if it's down, scheduled jobs simply
  don't fire (no catch-up/backfill behavior for missed runs).
- Uses cron.WithSeconds() for sub-minute granularity in testing/dev;
  production schedules would typically use minute-or-coarser
  expressions.
