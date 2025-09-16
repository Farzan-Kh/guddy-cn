package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/handler"
	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/jwtjw"
	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/store"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5"
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

	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer dbConn.Close()

	if err := store.InitDB(dbConn); err != nil {
		log.Fatalf("init db: %v", err)
	}

	storeRepo := store.New(dbConn)
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
