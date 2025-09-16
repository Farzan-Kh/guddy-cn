package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger initializes the package-level logger
func InitLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// responseWriter is a small wrapper to capture status and size for logging
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// RequestLogger returns a middleware that logs requests using zap
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w}

			start := time.Now()
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote", r.RemoteAddr),
				zap.Int("status", rw.status),
				zap.Int("size", rw.size),
				zap.Time("start", start),
				zap.Duration("duration", duration),
			)
		})
	}
}
