package redis

import (
	"ModeAuth/internal/shared/config"
	"ModeAuth/pkg/logging"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func ConnectRedis(ctx context.Context) *redis.Client {
	log.Println(logging.INFO + "Connecting to Redis")

	cfg := config.GetRedis()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf(logging.FATAL+"Failed to connect to Redis: %v", err)
	}

	log.Println(logging.INFO + "Connection to Redis successful")

	return rdb
}
