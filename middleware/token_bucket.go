package middleware

import (
	"net/http"

	"github.com/paralleltree/ratelimit-go"
)

func NewLimiterMiddleware(limiter ratelimit.Limiter, selector func(*http.Request) string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := selector(r)
			if limiter.Consume(key) {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusTooManyRequests)
		})
	}
}
