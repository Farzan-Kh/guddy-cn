package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool    *pgxpool.Pool
	Queriez *Queries
)

func InitFromEnv() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	Init(dsn)
}

func Init(dsn string) {
	var err error
	ctx := context.Background()

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	Queriez = New(pool)
}

func GetPool() *pgxpool.Pool {
	return pool
}

func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}
