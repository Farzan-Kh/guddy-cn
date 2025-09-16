package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Farzan-Kh/guddy-cn/services/authn/internal/models"
)

var ErrUserExists = errors.New("user exists")
var ErrUserNotFound = errors.New("user not found")

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func InitDB(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);`
	_, err := db.Exec(q)
	return err
}

func (s *Store) CreateUser(email, passwordHash string) (*models.User, error) {
	now := time.Now().UTC()
	var id int64
	err := s.db.QueryRow("INSERT INTO users(email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id", email, passwordHash, now).Scan(&id)
	if err != nil {
		return nil, ErrUserExists
	}
	u := &models.User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now}
	return u, nil
}

func (s *Store) GetByEmail(email string) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email)
	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
