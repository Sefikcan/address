package ratelimiter

import (
	"math"
	"sync"
	"time"
)

const (
	// MAX_BUCKET_SIZE token bucket max capacity
	MAX_BUCKET_SIZE float64 = 3
	// REFILL_RATE token bucket filling rate per second
	REFILL_RATE int = 1
)

type TokenBucket struct {
	currentBucketSize   float64
	lastRefillTimestamp int64
	mutex               sync.Mutex
}

// NewTokenBucket create new TokenBucket instance
func NewTokenBucket() *TokenBucket {
	return &TokenBucket{
		currentBucketSize:   MAX_BUCKET_SIZE,
		lastRefillTimestamp: getCurrentTimeInNanoseconds(),
	}
}

// getCurrentTimeInNanoseconds return time to nano second
func getCurrentTimeInNanoseconds() int64 {
	return time.Now().UnixNano()
}

// refill update bucket fill
func (tb *TokenBucket) refill() {
	now := getCurrentTimeInNanoseconds()
	tokensToAdd := float64(now-tb.lastRefillTimestamp) * float64(REFILL_RATE) / 1e9
	tb.currentBucketSize = math.Min(tb.currentBucketSize+tokensToAdd, MAX_BUCKET_SIZE)
	tb.lastRefillTimestamp = now
}

// AllowRequest determines whether a request will be accepted or not
func (tb *TokenBucket) AllowRequest(tokens float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()

	if tb.currentBucketSize >= tokens {
		tb.currentBucketSize -= tokens
		return true
	}

	return false
}
