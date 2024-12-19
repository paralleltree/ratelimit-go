package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/paralleltree/ratelimit-go"
	"github.com/paralleltree/ratelimit-go/middleware"
)

func main() {
	if err := start("https://yahoo.co.jp", 8081); err != nil {
		log.Fatalf("err: %v", err)
	}
}

func start(up string, port int) error {
	upstreamURL, err := url.Parse(up)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	limiter := ratelimit.NewTokenBucket(func() time.Time { return time.Now() }, 3, 3, time.Minute)
	selector := func(r *http.Request) string {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// TODO: log
		}
		return host
	}
	limiterMiddleware := middleware.NewLimiterMiddleware(limiter, selector)

	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	mux := http.NewServeMux()
	mux.Handle("/", limiterMiddleware(proxy))

	addr := fmt.Sprintf(":%d", port)
	server := http.Server{Addr: addr, Handler: mux}

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}
	}

	return nil
}
