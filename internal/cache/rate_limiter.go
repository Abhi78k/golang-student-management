package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

func (r *RateLimiter) Allow(
	ctx context.Context,
	userID string,
) (bool, error) {

	key := "rate:" + userID

	count, err := r.client.Incr(
		ctx,
		key,
	).Result()

	if err != nil {
		return false, err
	}

	if count > 3 {
		return false, nil
	}

	if count == 1 {
		err := r.client.Expire(
			ctx,
			key,
			time.Minute,
		).Err()

		if err != nil {
			return false, err
		}
	}

	return true, nil
}
