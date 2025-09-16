package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Service configuration and proxy logic moved to proxy.go

func main() {
	// initialize logger
	InitLogger()
	defer logger.Sync()

	r := chi.NewRouter()

	// Middleware
	r.Use(RequestLogger())
	// Initialize and register rate limiter (10 reqs per minute per user)
	if err := InitRateLimiter(); err != nil {
		logger.Error("failed to initialize rate limiter", zap.Error(err))
	} else {
		r.Use(RateLimitMiddleware)
	}
	r.Use(middleware.Recoverer)

	// Route all /api/* requests
	r.Route("/api", func(r chi.Router) {
		// Specific routes first (higher priority)
		r.HandleFunc("/docs/*", createProxyHandler("docs"))
		r.HandleFunc("/logger/*", createProxyHandler("logger"))

		// Catch-all for all other /api/* requests to exercises service
		r.HandleFunc("/*", createProxyHandler("exercises"))
	})

	logger.Info("Gateway starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// proxy functions are implemented in proxy.go

// logging helpers (InitLogger, RequestLogger and responseWriter) live in logger.go
