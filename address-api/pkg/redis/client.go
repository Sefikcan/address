package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/sefikcan/address-api/pkg/config"
)

var Client *redis.Client

func InitializeRedis(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	_, err := Client.Ping(Client.Context()).Result()
	return err
}

func GetClient() *redis.Client {
	return Client
}
