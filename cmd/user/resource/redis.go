package resource

import (
	"context"
	"fmt"
	stdLog "log"
	"userfc/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		stdLog.Fatalf("Failed to connect to redis: %v", err)
	}

	stdLog.Println("connected to redis : ")
	return RedisClient

}
