package main

import (
	"log"
	"net/http"

	"github.com/Farzan-kh/guddy-cn/exercises/internal/db"
	"github.com/Farzan-kh/guddy-cn/exercises/internal/router"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize database connection
	db.InitFromEnv()
	defer db.CloseDB()

	// Initialize router
	r := router.NewRouter()

	logger.Info("Exercises service starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
