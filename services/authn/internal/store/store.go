package store

import (
	"context"
	"errors"
	"time"

	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserExists = errors.New("user exists")
var ErrUserNotFound = errors.New("user not found")

type Store struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func InitDB(ctx context.Context, db *pgxpool.Pool) error {
	q := `CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);`
	_, err := db.Exec(ctx, q)
	return err
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	now := time.Now().UTC()
	var id int64
	err := s.db.QueryRow(ctx, "INSERT INTO users(email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id", email, passwordHash, now).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return nil, ErrUserExists
			}
		}
		return nil, err
	}

	user := &models.User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now}
	return user, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := s.db.QueryRow(ctx, "SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email)
	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
