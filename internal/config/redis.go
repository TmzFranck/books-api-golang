package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper, address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}
