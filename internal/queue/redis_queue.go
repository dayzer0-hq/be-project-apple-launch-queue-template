package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Key schema:
//   queue:{productID}                      → Sorted Set (score=joinTimestampMs, member=queueID)
//   entry:{queueID}                        → Hash { product_id, user_id, status, joined_at }
//   user_queue:{userID}:{productID}        → queueID (no TTL — cleared on claim/expire)

type RedisQueueStore struct{ rdb *redis.Client }

func NewRedisQueueStore(rdb *redis.Client) *RedisQueueStore {
	return &RedisQueueStore{rdb: rdb}
}

func (s *RedisQueueStore) Join(ctx context.Context, productID, userID string) (*Entry, error) {
	// TODO: Check user_queue:{userID}:{productID} — if exists, return existing entry
	// TODO: Generate queueID = "q_" + random hex
	// TODO: ZADD queue:{productID} NX score=UnixMilli() member=queueID
	// TODO: HSET entry:{queueID} product_id, user_id, status=waiting, joined_at=now
	// TODO: SET user_queue:{userID}:{productID} queueID (no TTL)
	// TODO: ZRANK queue:{productID} queueID → position (0-based, add 1 for display)
	return nil, ErrNotFound
}

func (s *RedisQueueStore) Position(ctx context.Context, queueID string) (*Entry, error) {
	// TODO: HGETALL entry:{queueID} → get product_id, user_id, status
	// TODO: ZRANK queue:{product_id} queueID → position
	// TODO: Estimate wait: position / 50 * 1 minute (50 users admitted per minute)
	return nil, ErrNotFound
}

func (s *RedisQueueStore) Claim(ctx context.Context, queueID string) (string, time.Time, error) {
	// TODO: HGET entry:{queueID} status → must be "ready"
	// TODO: Generate purchase_token = "pt_" + random hex
	// TODO: HSET entry:{queueID} status claimed purchase_token {token}
	// TODO: EXPIRE entry:{queueID} 15*60  (15 min claim window)
	return "", time.Time{}, ErrNotFound
}

func (s *RedisQueueStore) Advance(ctx context.Context, productID string, batchSize int) (int, error) {
	// TODO: ZRANGE queue:{productID} 0 batchSize-1 → get next N queue IDs
	// TODO: For each: HGET status → if "waiting", HSET status=ready
	// TODO: Return count of entries actually advanced
	return 0, nil
}
