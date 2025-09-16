package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Service configuration
type ServiceConfig struct {
	Name string
	Host string
	Port string
}

var services = map[string]ServiceConfig{
	"exercises": {Name: "exercises", Host: "exercises-service", Port: "8081"},
	"programs":  {Name: "exercises", Host: "exercises-service", Port: "8081"}, // Programs are handled by exercises service
	"docs":      {Name: "docs", Host: "docs", Port: "8082"},
	"logger":    {Name: "logger", Host: "logger", Port: "8083"},
}

func main() {
	// initialize logger
	InitLogger()
	defer logger.Sync()

	r := chi.NewRouter()

	// Middleware
	r.Use(RequestLogger())
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

func createProxyHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service, exists := services[serviceName]
		if !exists {
			logger.Error("Unknown service", zap.String("service", serviceName))
			http.NotFound(w, r)
			return
		}

		// Reconstruct the path for the target service
		// For docs and logger, remove /api/servicename to get /*
		// For exercises, keep the /api prefix since exercises service expects /api/*
		var targetPath string
		if serviceName == "exercises" {
			targetPath = r.URL.Path // Keep the full /api/* path
		} else {
			targetPath = strings.TrimPrefix(r.URL.Path, "/api/"+serviceName)
		}

		target := "http://" + service.Host + ":" + service.Port + targetPath

		proxyToService(w, r, target)
	}
}

func proxyToService(w http.ResponseWriter, r *http.Request, target string) {
	// Create a new request to the target service
	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		logger.Error("Failed to create request", zap.Error(err))
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Copy query parameters
	req.URL.RawQuery = r.URL.RawQuery

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to send request", zap.Error(err), zap.String("target", target))
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", zap.Error(err))
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Write response body
	w.Write(body)
}

// logging helpers (InitLogger, RequestLogger and responseWriter) live in logger.go
