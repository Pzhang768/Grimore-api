# CLAUDE.md — Grimore-api

Backend for Grimoire: a Go/Gin API that orchestrates AI agent pipelines for job search automation.

**Stack:** Go, Gin, PostgreSQL (GORM), Supabase Auth, Stripe, Railway

## Repository Structure

```
cmd/
└── server/
    └── main.go            # Entry point

internal/
├── handlers/              # Gin route handlers (HTTP layer only)
├── services/              # Business logic
├── agents/                # Individual agent implementations (fetcher, analyser, tailor)
├── ai/                    # Claude provider interface, implementation, and Coordinator (pipeline orchestration)
├── db/                    # Database connection and queries
├── models/                # GORM structs
└── middleware/            # JWT auth, logging, rate limiting

pkg/
└── sse/                   # SSE streaming utilities

config/
└── config.go              # Env vars and app config
```

## Commands

```bash
go run cmd/server/main.go   # Dev server on :8080
go test ./...               # Run all tests
go test ./... -cover        # Tests with coverage
go build ./cmd/server       # Build binary
go vet ./...                # Static analysis
gofmt -w .                  # Format
```

## Environment Variables (`.env`)

```
PORT=8080
DATABASE_URL=postgres://localhost:5432/grimoire
ANTHROPIC_API_KEY=
SUPABASE_URL=
SUPABASE_SERVICE_ROLE_KEY=
STRIPE_SECRET_KEY=
STRIPE_PRICE_ID=
```

## Agent Pipeline

```
JobListingFetcher → FitAnalyser → ResumeTailoringAgent → checkpoint (max 3 iterations)
```

Individual agents (fetcher, analyser, tailor) call `ai.Provider` via constructor injection — never import the Anthropic SDK directly.
The Coordinator in `ai/` owns execution order, context passing between agents, and checkpoint signalling.
SSE events stream to the frontend via `pkg/sse`.

## Database Schema (Key Tables)

- **users** — id, email, created_at
- **teams** — id, user_id, name, created_at
- **team_agents** — team_id, agent_type, context (JSON), position
- **runs** — id, team_id, status, iteration, created_at
- **run_events** — id, run_id, agent_type, event_type, content, created_at
- **deliverables** — id, run_id, type, content (JSON), created_at

## Code Rules

- **[Architecture](../rules/architecture.md)** — layered architecture, banned patterns, AI provider interface
- **[General](../rules/general.md)** — naming conventions, Go quality rules
- **[Testing](../rules/testing.md)** — 100% coverage, table-driven tests, mock AI provider
- **[CI & Git](../rules/ci.md)** — branch strategy, conventional commits, pre-commit hooks
- **[DRY](../rules/dry.md)** / **[KISS](../rules/kiss.md)** / **[SOLID](../rules/solid.md)** / **[YAGNI](../rules/yagni.md)**

## Security

- Never commit `.env` — copy `.env.example` to set up locally
- All routes require Supabase JWT; queries scoped to `user_id`
- Anthropic and Stripe keys never leave the backend
