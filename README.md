# chirpy

A small HTTP server in Go, built in Boot.dev's "Learn HTTP Servers in Go" course: a JSON API for short posts with users,
sign-in and a payment-provider webhook.

- **Standard library routing** (`net/http` method patterns), JSON in and out, a metrics and reset admin pair.
- **Auth:** Argon2id password hashing, short-lived JWT access tokens, refresh tokens stored in Postgres with revoke, and
  a webhook endpoint checked by API key.
- **Data:** PostgreSQL with goose migrations (`sql/schema`) and sqlc-generated queries (`sql/queries`, `internal/database`).
- **Tests:** `go test ./...` (10 tests for hashing, JWT creation and expiry, bearer and API-key parsing).

## Run it

```bash
export DB_URL="postgres://USER:PASSWORD@localhost:5432/chirpy?sslmode=disable"
export JWT_SECRET="$(openssl rand -base64 48)"
export POLKA_KEY="any-test-key"
export PLATFORM=dev          # enables POST /admin/reset
goose -dir sql/schema postgres "$DB_URL" up
go run .
```

Endpoints: `GET /api/healthz` · `POST /api/users` · `PUT /api/users` · `POST /api/login` · `POST /api/refresh` ·
`POST /api/revoke` · `GET/POST /api/chirps` · `GET/DELETE /api/chirps/{chirpID}` · `POST /api/polka/webhooks` ·
`GET /admin/metrics` · `POST /admin/reset`
