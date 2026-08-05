# Capsule

Zero-knowledge, end-to-end encrypted paste tool. Secrets self-destruct after reading.

## Tech Stack

- **API**: Go 1.24, Chi router, GORM, PostgreSQL 16, [`tronc`](https://github.com/FacileStudio/tronc) as the app chassis
- **Client**: SvelteKit 5 (Svelte 5 runes), Tailwind CSS 4, Vite 7, adapter-node
- **Encryption**: AES-256-GCM via Web Crypto API (client-side only)
- **Runtime**: Bun (client), Docker Compose for full stack
- **Gate**: `sh scripts/check.sh` via a pre-push hook. GitHub Actions is gone suite-wide

## Project Structure

```
apps/
  api/                   # Go backend
    main.go              # Entrypoint: Chi router, graceful shutdown, cleanup goroutine
    internal/
      cleanup/           # Background expired-paste cleanup (has tests)
      database/          # GORM database connection
      env/               # Capsule-only config, wrapping tronc/env
      middleware/        # Rate limiting (the rest comes from tronc)
    modules/
      docs/              # OpenAPI spec + Scalar UI at /docs
      pastes/            # Core paste CRUD: handler, service, router, types (has tests)
    schemas/             # GORM models + migrations
  client/                # SvelteKit frontend
    src/
      lib/
        crypto.ts        # AES-256-GCM encrypt/decrypt (has unit test)
        backend.ts       # API client
        highlight.ts     # Shiki syntax highlighting
        components/      # Svelte components (ThemeToggle, etc.)
      routes/
        (app)/           # Main pages: create, reveal ([id]), revoke
        api/[...path]/   # Proxy to Go API
        docs/            # Scalar API docs page
    tests/               # Playwright E2E tests
docker-compose.yml       # Full stack: postgres + api + client
```

## Key Commands

### API (from `apps/api/`)

```bash
go run .                 # Start dev server (port 4000)
go test ./...            # Run all tests (needs PostgreSQL for integration tests)
go vet ./...             # Lint
go build -o bin/api .    # Build binary
```

### Client (from `apps/client/`)

```bash
bun install              # Install dependencies
bun dev                  # Dev server (port 5173)
bun run build            # Production build
bun run check            # Svelte type checking
bun run test             # Vitest unit tests
bun run test:e2e         # Playwright E2E tests (starts dev server automatically)
```

### Full Stack

```bash
docker compose up -d     # Start everything (postgres + api + client)
docker compose down      # Stop and remove containers
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | **required** | PostgreSQL connection string. No default since the tronc adoption — an app that boots against a database that isn't there just 500s later |
| `APP_ENV` | `development` | `development`, `staging`, `production`. Never gates security behaviour |
| `PORT` | `8080` | API server port. Compose and Dokploy set `4000` |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `MAX_PASTE_SIZE` | `1048576` | Max paste size in bytes (default 1MB) |
| `CORS_ALLOWED_ORIGINS` | *(none — deny)* | Comma-separated CORS origins. Falls back to `ALLOWED_ORIGINS`. Unset denies every cross-origin caller, which is correct here: the client proxies server-side |
| `API_URL` | `http://localhost:4000` | Backend URL (used by client proxy) |
| `ORIGIN` | `http://localhost:3000` | Client origin (used by SvelteKit) |

## Conventions

- The server never sees plaintext. Encryption/decryption happens exclusively in the browser via Web Crypto API. The URL fragment carries the decryption key (never sent to server).
- Svelte 5 runes API (`$state`, `$props`, `$derived`, `$effect`) is enforced via vite plugin config.
- API tests use SQLite in-memory for unit tests and PostgreSQL for integration tests.
- Rate limiting on paste creation: 30 requests per minute.
- Background cleanup goroutine runs on startup to purge expired pastes.
- The client proxies API requests through SvelteKit's server routes (`/api/[...path]`) to avoid CORS in dev. The API is `expose`d, never published, so it is never called cross-origin.
- The error envelope, JSON helpers, logger, request logging, CORS, panic recovery, `/health` and `/ready` all come from `tronc`. Do not reintroduce local copies — fix them upstream.
- `/health` and `/ready` answer at both `/` and `/api`, so the same probe works through the edge.
- The container healthcheck re-executes the binary: `["CMD", "/app/api", "healthcheck"]`.
