package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/handler"
	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/jwtjw"
	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/store"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("starting authn service")

	// config via env
	// Postgres DSN, e.g. "postgres://user:pass@host:5432/dbname?sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL env is required")
	}

	jwtSecret := os.Getenv("AUTHN_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("AUTHN_JWT_SECRET env is required")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()
	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	storeRepo := store.New(pool)
	jwtSvc := jwtjw.New([]byte(jwtSecret), time.Hour*24)
	h := handler.New(storeRepo, jwtSvc)

	r := mux.NewRouter()
	r.HandleFunc("/signup", h.SignUp).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/validate", h.Validate).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
