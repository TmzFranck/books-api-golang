package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	EXPIRE_TIME = 3600 // 1 hour
)

// AddJwtToBlacklist adds the given JWT to the blacklist in Redis with an expiration time
func AddJwtToBlacklist(ctx context.Context, redisClient *redis.Client, jwt string) error {
	err := redisClient.Set(ctx, "jwt", jwt, EXPIRE_TIME).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetJwtBlacklist returns true if the given JWT is in the blacklist, false otherwise
func GetJwtBlacklist(ctx context.Context, redisClient *redis.Client) (bool, error) {
	jwt, err := redisClient.Get(ctx, "jwt").Result()
	if err != nil {
		return false, err
	}
	return jwt != "", nil
}
