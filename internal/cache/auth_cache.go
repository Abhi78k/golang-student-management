package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthCache struct {
	client *redis.Client
}

func NewAuthCache(client *redis.Client) *AuthCache {
	return &AuthCache{
		client: client,
	}
}

func (c *AuthCache) StoreRefreshToken(
	ctx context.Context,
	userID string,
	token string,
) error {

	key := "refresh:" + userID

	return c.client.Set(
		ctx,
		key,
		token,
		15*24*time.Hour,
	).Err()
}

func (c *AuthCache) GetRefreshToken(
	ctx context.Context,
	userID string,
) (string, error) {

	key := "refresh:" + userID

	return c.client.Get(
		ctx,
		key,
	).Result()
}

func (c *AuthCache) DeleteRefreshToken(
	ctx context.Context,
	userID string,
) error {

	key := "refresh:" + userID

	return c.client.Del(
		ctx,
		key,
	).Err()
}
