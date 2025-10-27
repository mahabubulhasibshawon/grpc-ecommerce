package infrastructure

import (
	"context"
	"log"
	"time"

	"github.com/mahabubulhasibshawon/grpc-ecommerce.git/internal/adapters/redis"
)

func ConnectRedis(addr,
	username,
	password string,
	ttl time.Duration) *redis.Cache {
	cache := redis.NewCache(addr, username, password, 0, ttl)
	if err := cache.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully: ",cache)
	return cache
}
