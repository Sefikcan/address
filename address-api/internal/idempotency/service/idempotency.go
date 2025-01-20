package idempotency

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type IdempotencyService interface {
	GetByKey(ctx context.Context, key string) (interface{}, error)
	Save(ctx context.Context, key string, result []byte) error
}

type idempotencyService struct {
	redisClient *redis.Client
}

func (i idempotencyService) GetByKey(ctx context.Context, key string) (interface{}, error) {
	result, err := i.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var response interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i idempotencyService) Save(ctx context.Context, key string, result []byte) error {
	expiration := time.Hour // TODO: Move configuration or another constant
	return i.redisClient.Set(ctx, key, result, expiration).Err()
}

func NewIdempotencyService(redisClient *redis.Client) IdempotencyService {
	return &idempotencyService{
		redisClient: redisClient,
	}
}
