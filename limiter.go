package ratelimit

type Limiter interface {
	Consume(key string) bool
}
