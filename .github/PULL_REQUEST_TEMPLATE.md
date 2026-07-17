## What this PR does

## Issues closed

Closes #

## Checklist

- [ ] Join is idempotent (same user + product returns same queue_id)
- [ ] Position decreases as users ahead claim
- [ ] Claim fails with 409 when status != ready
- [ ] Advance is safe to call concurrently (no double-advancing)
- [ ] `go test -race ./...` passes
