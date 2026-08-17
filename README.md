# Apple Store Launch Queue — DayZer0 Template

Starter for the **"Build Apple Store's product launch queue"** project at DayZer0.

## Prerequisites

- Go 1.22+
- Docker / Docker Compose (for Redis)

## Get started

```bash
docker compose up -d     # Redis on :6379
cp .env.example .env
go run ./cmd/server
```

## What you need to build

Implement all methods in `internal/queue/redis_queue.go`. The handlers are pre-written.

Key Redis patterns:
- `ZADD` with `NX` flag for fair ordering
- `ZRANK` for position lookup
- `HSET` / `HGET` / `EXPIRE` for entry state

## Submitting

Open a PR to `main`. Include `Closes #N` for every issue. Raj reviews automatically.
