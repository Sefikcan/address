package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

type EPDistributedTokenBucket struct {
	client *redis.Client
	ctx    context.Context
}

// NewEPDistributedTokenBucket create new DistributedTokenBucket instance
func NewEPDistributedTokenBucket(redisAddr string) *EPDistributedTokenBucket {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &EPDistributedTokenBucket{
		client: rdb,
		ctx:    context.Background(),
	}
}

// AllowRequest controls whether a request will be accepted or not
func (rtb *EPDistributedTokenBucket) EPAllowRequest(ctx context.Context, key, endpoint string, tokens float64) bool {
	bucketKey := fmt.Sprintf("%s:%s", key, endpoint)

	currentBucketSize, err := rtb.client.Get(ctx, bucketKey).Float64()
	if err == redis.Nil {
		currentBucketSize = MAX_BUCKET_SIZE
	} else if err != nil {
		fmt.Printf("Error getting bucket size for %s: %v\n", bucketKey, err)
		return false
	}

	now := time.Now().UnixNano()
	lastRefillTimestamp, err := rtb.client.HGet(rtb.ctx, bucketKey, "lastRefill").Int64()
	if err == redis.Nil {
		lastRefillTimestamp = now
	} else if err != nil {
		fmt.Printf("Error getting last refill timestamp for %s: %v\n", bucketKey, err)
		return false
	}

	tokensToAdd := float64(now-lastRefillTimestamp) * float64(REFILL_RATE) / 1e9
	newBucketSize := math.Min(currentBucketSize+tokensToAdd, MAX_BUCKET_SIZE)

	if newBucketSize >= tokens {
		newBucketSize -= tokens
		_, err := rtb.client.TxPipelined(rtb.ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(rtb.ctx, bucketKey, newBucketSize, 0)
			pipe.HSet(rtb.ctx, bucketKey, "lastRefill", now)
			return nil
		})
		if err != nil {
			fmt.Printf("Error updating bucket for %s: %v\n", bucketKey, err)
			return false
		}
		return true
	}

	_, err = rtb.client.TxPipelined(rtb.ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(rtb.ctx, bucketKey, newBucketSize, 0)
		pipe.HSet(rtb.ctx, bucketKey, "lastRefill", now)
		return nil
	})
	if err != nil {
		fmt.Printf("Error updating bucket without token use for %s: %v\n", bucketKey, err)
	}

	return false
}
