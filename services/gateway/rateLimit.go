package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"context"

	goredis "github.com/redis/go-redis/v9"
	limiterpkg "github.com/ulule/limiter/v3"
	redstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// limiterInstance is initialized in InitRateLimiter
var limiterInstance *limiterpkg.Limiter

// InitRateLimiter initializes a sliding-window limiter with 10 requests per minute
// per-IP. It uses Redis as the backing store at redis:6379 so limits are shared
// across gateway replicas.
func InitRateLimiter() error {
	// 10 requests per minute -> formatted as "10-M"
	rate, err := limiterpkg.NewRateFromFormatted("10-M")
	if err != nil {
		return err
	}

	// Create go-redis client pointing to the k8s service name `redis:6379`.
	// The zero value context is fine for initialization.
	client := goredis.NewClient(&goredis.Options{
		Addr: "redis:6379",
	})

	// Test the connection briefly
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	store, err := redstore.NewStoreWithOptions(client, limiterpkg.StoreOptions{Prefix: "gateway_limiter"})
	if err != nil {
		return err
	}

	limiterInstance = limiterpkg.New(store, rate)
	return nil
}

// keyForRequest extracts an IP-based key from the request (no X-User-ID used).
func keyForRequest(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RateLimitMiddleware enforces the configured limiter per extracted key.
// It also sets common rate-limit headers: X-RateLimit-Limit and X-RateLimit-Remaining
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiterInstance == nil {
			next.ServeHTTP(w, r)
			return
		}

		key := keyForRequest(r)
		ctx, err := limiterInstance.Get(r.Context(), key)
		if err != nil {
			// On error, allow the request but don't fail open loudly
			next.ServeHTTP(w, r)
			return
		}

		// Set informative headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", ctx.Limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", ctx.Remaining))
		// Try to set reset header if available (some implementations provide Reset as time.Time)
		switch t := any(ctx.Reset).(type) {
		case time.Time:
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", t.Unix()))
		case int64:
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", t))
		default:
			// ignore
		}

		// If no remaining requests, respond with 429
		if ctx.Remaining <= 0 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
