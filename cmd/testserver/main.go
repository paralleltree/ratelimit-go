package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/paralleltree/ratelimit-go"
	"github.com/paralleltree/ratelimit-go/middleware"
)

func main() {

	limiter := ratelimit.NewTokenBucket(func() time.Time { return time.Now() }, 3, 3, time.Minute)
	selector := func(r *http.Request) string {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// TODO: log
		}
		return host
	}
	limiterMiddleware := middleware.NewLimiterMiddleware(limiter, selector)

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	})

	mux := http.NewServeMux()
	mux.Handle("/", limiterMiddleware(mainHandler))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("unexpected error: %v", err)
		}
	}
}
