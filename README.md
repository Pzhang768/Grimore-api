# Grimore-api

Go/Gin backend for Grimoire — orchestrates AI agent pipelines for job search automation.

## Prerequisites

- Go 1.22+
- PostgreSQL (local)

## Setup

```bash
cp .env.example .env
go mod download
go run cmd/server/main.go
```

Runs at `http://localhost:8080`.

### Database

```bash
createdb grimoire
go run cmd/migrate/main.go
```

## Environment Variables

| Variable | Description |
|---|---|
| `PORT` | Server port (default `8080`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `ANTHROPIC_API_KEY` | Claude API key |
| `SUPABASE_URL` | Supabase project URL |
| `SUPABASE_SERVICE_ROLE_KEY` | Supabase service role key |
| `STRIPE_SECRET_KEY` | Stripe secret key |
| `STRIPE_PRICE_ID` | Stripe subscription price ID |

## Commands

```bash
go run cmd/server/main.go   # Dev server
go test ./...               # Run all tests
go test ./... -cover        # With coverage
go vet ./...                # Static analysis
gofmt -w .                  # Format
```
