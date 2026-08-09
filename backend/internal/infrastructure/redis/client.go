package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	options, err := redisOptions(cfg)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err = client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect to Redis: %w", err)
	}

	return client, nil
}

func redisOptions(cfg *config.Config) (*redis.Options, error) {
	if cfg.RedisURI != "" {
		options, err := redis.ParseURL(cfg.RedisURI)
		if err != nil {
			return nil, fmt.Errorf("parse REDIS_URI: %w", err)
		}
		return options, nil
	}

	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}, nil
}
