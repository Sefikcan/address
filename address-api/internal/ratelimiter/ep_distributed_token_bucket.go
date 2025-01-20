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
	lastRefillKey := fmt.Sprintf("%s:last_refill", key)

	// Current bucket size
	currentBucketSize, err := rtb.client.Get(ctx, bucketKey).Float64()
	if err == redis.Nil {
		currentBucketSize = MAX_BUCKET_SIZE
	} else if err != nil {
		fmt.Printf("Error getting bucket size for %s: %v\n", bucketKey, err)
		return false
	}

	// Last refill time
	now := time.Now().UnixNano()
	lastRefillTime, err := rtb.client.Get(ctx, lastRefillKey).Int64()
	if err == redis.Nil {
		lastRefillTime = now
	} else if err != nil {
		fmt.Printf("Error getting last refill timestamp for %s: %v\n", bucketKey, err)
		return false
	}

	// Calculate tokens to add based on time passed since last refill
	timePassed := float64(now - lastRefillTime)
	tokensToAdd := timePassed * REFILL_RATE / 1e9 // Convert nanoseconds to seconds
	newBucketSize := math.Min(currentBucketSize+tokensToAdd, MAX_BUCKET_SIZE)

	// Check if we have enough tokens for the request
	if newBucketSize >= tokens {
		newBucketSize -= tokens
		// Update bucket size and last refill time atomically
		_, err = rtb.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, bucketKey, newBucketSize, 0)
			pipe.Set(ctx, lastRefillKey, now, 0)
			return nil
		})
		if err != nil {
			fmt.Printf("Error updating bucket for %s: %v\n", bucketKey, err)
			return false
		}
		return true
	}

	// Even if we don't have enough tokens for the request, update the bucket size
	_, err = rtb.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, bucketKey, newBucketSize, 0)
		pipe.Set(ctx, lastRefillKey, now, 0)
		return nil
	})
	if err != nil {
		fmt.Printf("Error updating bucket without token use for %s: %v\n", bucketKey, err)
	}

	return false
}
