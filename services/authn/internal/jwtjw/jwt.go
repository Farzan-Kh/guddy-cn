package jwtjw

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalid = errors.New("invalid token")

type Service struct {
	secret []byte
	expire time.Duration
}

func New(secret []byte, expire time.Duration) *Service {
	return &Service{secret: secret, expire: expire}
}

func (s *Service) Generate(subject string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expire)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Service) Validate(tokenStr string) (string, error) {
	p := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, p, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return "", ErrInvalid
	}
	return p.Subject, nil
}
