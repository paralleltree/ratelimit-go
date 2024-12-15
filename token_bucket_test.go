package ratelimit_test

import (
	"testing"
	"time"

	"github.com/paralleltree/ratelimit-go"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket(t *testing.T) {
	// arrange
	initTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	setTime, timeProvider := func(initTime time.Time) (func(time.Time), func() time.Time) {
		current := initTime
		return func(t time.Time) { current = t },
			func() time.Time { return current }
	}(initTime)
	tokenCapacity := 3
	// interval/count(1個あたりの補充数)の次元の単位にすると厳密にintervalを待たずして補充できる？
	// 今は10秒経たないと3つ補充されない
	// vs. 10s/3 = 3.3秒経ったら1個補充する
	replenishCount := 3
	replenishInterval := 10 * time.Second
	limiter := ratelimit.NewTokenBucket(timeProvider, tokenCapacity, replenishCount, replenishInterval)
	key := "testkey"

	// act/assert
	assert.True(t, limiter.Consume(key))
	assert.True(t, limiter.Consume(key))
	assert.True(t, limiter.Consume(key))

	// トークンを使いきったので通らない
	assert.False(t, limiter.Consume(key))

	setTime(initTime.Add(1 * time.Second))
	// まだ補充されないので通らない
	assert.False(t, limiter.Consume(key))

	setTime(initTime.Add(10 * time.Second))
	// 補充され通る
	assert.True(t, limiter.Consume(key))
	assert.True(t, limiter.Consume(key))
	assert.True(t, limiter.Consume(key))

	// トークンを使いきったので通らない
	assert.False(t, limiter.Consume(key))
}
