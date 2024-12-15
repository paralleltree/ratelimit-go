package ratelimit

import (
	"sync"
	"time"
)

type tokenBucket struct {
	mu                sync.Mutex
	timeProvider      func() time.Time
	tokenCapacity     int
	replenishCount    int
	replenishInterval time.Duration
	tokenStore        map[string]int
	lastReplenishedAt map[string]time.Time
}

func NewTokenBucket(
	timeProvider func() time.Time,
	tokenCapacity int,
	replenishCount int,
	replenishInterval time.Duration,
) *tokenBucket {
	return &tokenBucket{
		mu:                sync.Mutex{},
		timeProvider:      timeProvider,
		tokenCapacity:     tokenCapacity,
		replenishCount:    replenishCount,
		replenishInterval: replenishInterval,
		tokenStore:        make(map[string]int),
		lastReplenishedAt: make(map[string]time.Time),
	}
}

func (b *tokenBucket) Consume(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tryInitialize(key)
	b.replenish(key)
	return b.tryDecrementToken(key)
}

func (b *tokenBucket) tryInitialize(key string) {
	if _, ok := b.lastReplenishedAt[key]; ok {
		return
	}
	b.lastReplenishedAt[key] = b.timeProvider()
	b.tokenStore[key] = b.tokenCapacity
}

func (b *tokenBucket) replenish(key string) {
	lastReplenishedAt := b.lastReplenishedAt[key]
	sinceLastReplenishment := b.timeProvider().Sub(lastReplenishedAt)
	replenishedTokenCount := b.replenishCount * int(sinceLastReplenishment/b.replenishInterval)
	currentTokenCount := b.tokenStore[key]
	b.tokenStore[key] = min(currentTokenCount+replenishedTokenCount, b.tokenCapacity)
	b.lastReplenishedAt[key] = lastReplenishedAt.Add(b.replenishInterval * (sinceLastReplenishment / b.replenishInterval))
}

func (b *tokenBucket) tryDecrementToken(key string) bool {
	current := b.tokenStore[key]
	if current == 0 {
		return false
	}
	b.tokenStore[key] = current - 1
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
