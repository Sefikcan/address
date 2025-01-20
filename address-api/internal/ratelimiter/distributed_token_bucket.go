package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type DistributedTokenBucket struct {
	client *redis.Client
}

// NewDistributedTokenBucket create new DistributedTokenBucket instance
func NewDistributedTokenBucket(redisClient *redis.Client) *DistributedTokenBucket {
	return &DistributedTokenBucket{
		client: redisClient,
	}
}

// AllowRequest controls whether a request will be accepted or not
func (dtb *DistributedTokenBucket) AllowRequest(ctx context.Context, key string, tokens float64) bool {
	// date convert nanotime
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// create redis key
	bucketKey := fmt.Sprintf("%s:bucket", key)
	lastRefillKey := fmt.Sprintf("%s:last_refill", key)

	// check lastrefill time
	lastRefill, err := dtb.client.Get(ctx, lastRefillKey).Int64()
	if err != nil {
		lastRefill = now
		dtb.client.Set(ctx, lastRefillKey, now, 0)
		dtb.client.Set(ctx, bucketKey, MAX_BUCKET_SIZE, 0)
	}

	// calculate amount of tokens to be added elapsed time
	elapsed := float64(now-lastRefill) / 1000
	tokensToAdd := REFILL_RATE * elapsed
	dtb.client.IncrByFloat(ctx, bucketKey, tokensToAdd)

	// control number of available token
	currentTokens, err := dtb.client.Get(ctx, bucketKey).Float64()
	if err != nil {
		return false
	}

	// accept the request if there are enough tokens
	if currentTokens >= tokens {
		newBucketSize := currentTokens - tokens
		dtb.client.Set(ctx, bucketKey, newBucketSize, 0)
		dtb.client.Set(ctx, lastRefillKey, now, 0)
		return true
	}

	// if token value less than zero, set bucket to zero
	if currentTokens < 0 {
		dtb.client.Set(ctx, bucketKey, 0, 0)
	}

	// update refill time, but reject the request
	dtb.client.Set(ctx, lastRefillKey, now, 0)
	return false
}
