package ratelimit

import "time"

type tokenBucket struct {
	timeProvider      func() time.Time
	tokenCapacity     int
	replenishCount    int
	replenishInterval time.Duration
}

func NewTokenBucket(
	timeProvider func() time.Time,
	tokenCapacity int,
	replenishCount int,
	replenishInterval time.Duration,
) *tokenBucket {
	return &tokenBucket{
		timeProvider:      timeProvider,
		tokenCapacity:     tokenCapacity,
		replenishCount:    replenishCount,
		replenishInterval: replenishInterval,
	}
}

func (b *tokenBucket) Consume(key string) bool {
	return true
}
