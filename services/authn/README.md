Authn microservice

Environment variables:
- AUTHN_DB_DSN (required) Postgres DSN, e.g. "postgres://user:pass@host:5432/dbname?sslmode=disable"
- AUTHN_JWT_SECRET (required)
- PORT (optional, default 8080)

Endpoints:
- POST /signup {"email":"...","password":"..."} -> 201
- POST /login {"email":"...","password":"..."} -> {"token":"..."}
- GET /validate (Authorization: Bearer <token>) -> returns subject (user id)

Build:
- go build

Docker:
- docker build -t authn:local .

