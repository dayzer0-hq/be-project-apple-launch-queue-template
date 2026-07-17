package queue

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAlreadyInQueue = errors.New("already_in_queue")
	ErrNotReady       = errors.New("not_ready")
	ErrNotFound       = errors.New("not_found")
	ErrExpired        = errors.New("expired")
)

const ClaimTTL = 15 * time.Minute

type Entry struct {
	QueueID    string    `json:"queue_id"`
	ProductID  string    `json:"product_id"`
	UserID     string    `json:"user_id"`
	Position   int64     `json:"position"`
	Status     string    `json:"status"` // "waiting" | "ready" | "claimed" | "expired"
	JoinedAt   time.Time `json:"joined_at"`
	EstWaitMin int       `json:"estimated_wait_minutes"`
}

// QueueStore is the interface handlers use.
// Implement it in queue/redis_queue.go using Redis sorted sets.
type QueueStore interface {
	// Join adds userID to the queue for productID.
	// Returns ErrAlreadyInQueue if a non-expired entry exists; returns that entry instead.
	// Score = Unix milliseconds (lower = earlier = better position).
	Join(ctx context.Context, productID, userID string) (*Entry, error)

	// Position returns the current queue entry and position.
	// Returns ErrNotFound if the entry doesn't exist.
	Position(ctx context.Context, queueID string) (*Entry, error)

	// Claim moves an entry from "ready" to "claimed" and sets a 15-min expiry.
	// Returns ErrNotReady if status != "ready".
	// Returns a purchase token on success.
	Claim(ctx context.Context, queueID string) (purchaseToken string, expiresAt time.Time, err error)

	// Advance moves the next batchSize "waiting" entries to "ready".
	// Called by the cron handler. Must be idempotent.
	Advance(ctx context.Context, productID string, batchSize int) (int, error)
}
