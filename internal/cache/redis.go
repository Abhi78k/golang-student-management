package cache

import (
	"context"
	"project-1/internal/config"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {

	client := redis.NewClient(
		&redis.Options{
			Addr: cfg.RedisAddr,
		},
	)

	err := client.Ping(
		context.Background(),
	).Err()

	if err != nil {
		panic(err)
	}

	return client
}
