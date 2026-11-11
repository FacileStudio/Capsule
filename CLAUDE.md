# Capsule

Zero-knowledge, end-to-end encrypted paste tool. Secrets self-destruct after reading.

## Tech Stack

- **API**: Go 1.26, Chi router, GORM, PostgreSQL 16, [`tronc`](https://github.com/FacileStudio/tronc) as the app chassis, [`tronc/migrate`](https://github.com/FacileStudio/tronc/tree/main/migrate) for schema
- **Client**: SvelteKit 5 (Svelte 5 runes), Tailwind CSS 4, Vite 7, adapter-static (the Go binary serves the build)
- **Design system**: [`@facile/muse`](https://github.com/FacileStudio/muse), pinned to a tag. It owns the palette, the Goga faces, dark mode and the base layer
- **Encryption**: AES-256-GCM via Web Crypto API (client-side only)
- **Runtime**: Bun (client), Docker Compose for full stack
- **Gate**: `sh scripts/check.sh` and `filet gate .`, run before every push by the tracked `lefthook.yml` (`mise run hooks` installs it). `.github/workflows/ci.yml` runs gofmt, vet, race tests on Postgres, svelte-check, the client tests, filet v0.19.0 (`preset: facile` in `filet.yml`) and douane 0.7.0 (`-fail critical -scope prod`). `scripts/check.sh` runs the `filet` on `PATH`, which must be v0.19.0 or later: older ones refuse `preset: facile`

## Project Structure

```
apps/
  api/                   # Go backend
    main.go              # Entrypoint: Chi router, graceful shutdown, cleanup goroutine
    internal/
      cleanup/           # Background expired-paste cleanup (has tests)
      database/          # GORM database connection
      env/               # Capsule-only config, wrapping tronc/env
    modules/
      docs/              # OpenAPI spec + Scalar UI at /docs
      pastes/            # Core paste CRUD: handler, service, router, types (has tests)
    migrations/          # Ordered SQL, embedded and applied at boot. Owns the schema
    schemas/             # GORM models (the schema is migrations/, not these)
  client/                # SvelteKit frontend
    src/
      lib/
        crypto.ts        # AES-256-GCM encrypt/decrypt (has unit test)
        backend.ts       # API client over @facile/client's createApi (CSRF header, error envelope)
        highlight.ts     # Shiki syntax highlighting
        theme.svelte.ts  # @facile/client theme store under `capsule-theme`, mirrored by app.html
        components/      # Svelte components (ThemeToggle, etc.)
      routes/
        (app)/           # Main pages: create, reveal ([id]), revoke
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
bun run test:e2e         # Playwright E2E — needs a real API, see below
PLAYWRIGHT_BASE_URL=https://capsule.facile.studio bun run test:e2e   # run against a deployment
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
| `MAX_PASTE_SIZE` | `16777216` | Max paste size in bytes (default 16MB) |
| `CORS_ALLOWED_ORIGINS` | *(none — deny)* | Comma-separated CORS origins. Falls back to `ALLOWED_ORIGINS`. Unset denies every cross-origin caller, which is correct here: the client proxies server-side |
| `API_URL` | `http://localhost:4000` | Backend URL (used by client proxy) |
| `ORIGIN` | `http://localhost:3000` | Client origin (used by SvelteKit) |

## Conventions

- The server never sees plaintext. Encryption/decryption happens exclusively in the browser via Web Crypto API. The URL fragment carries the decryption key (never sent to server).
- Svelte 5 runes API (`$state`, `$props`, `$derived`, `$effect`) is enforced via vite plugin config.
- UI comes from muse: reach for its components and `fc-*` tokens before hand-rolling markup, and read its `CHARTE.md` first. `app.css` imports `@facile/muse/styles` and must not `@import 'tailwindcss'` a second time — tokens.css already does. The `@theme inline` block there only aliases the suite's semantic names onto `fc-*`; no colour is declared locally.
- `optimizeDeps.exclude: ['@facile/muse', '@facile/client']` in `vite.config.ts` is load-bearing: both ship uncompiled `.svelte.ts`, and without it `vite dev` refuses to start while `bun run build` and `bun run check` stay green.
- The fetch kernel is `@facile/client` (`createApi`), with no `onUnauthorized`: Capsule has no `/login` page, so the suite's `redirectToLogin` would send a 401 to a page that does not exist. No paste route answers 401 anyway.
- The theme writes **both** `.dark` and `.light` on `<html>` (`system` writes neither). muse flips its tokens from `prefers-color-scheme` scoped to `:root:not(.light)`, so only the `.light` class can force light on a dark OS.
- API tests need PostgreSQL. **Without `TEST_DATABASE_URL` the database-backed tests skip and `go test` still prints `ok` and exits 0** — 35 of 41 tests, as of 2026-08-07. `sh scripts/check.sh` passing is therefore not evidence the API is tested; export `TEST_DATABASE_URL` to actually run them.
- The client pages call `/api` and there is **no dev proxy**, so `vite dev` has no backend and the sealing E2E tests cannot pass locally. Point the suite at a real deployment instead: `PLAYWRIGHT_BASE_URL=https://capsule.facile.studio bun run test:e2e`. Setting that variable also suppresses the dev-server boot in `playwright.config.ts`. The consequence to keep in mind: run that way, E2E validates the *deployed* build, not your working tree — so it is a post-deploy check, not a pre-deploy one.
- Sign-in is Registre only, through porte, and exists for the CLI. Pastes stay anonymous: no route under `/api/pastes` needs a session, and a reveal link works without one. The web client has no login page.
- tronc's security headers are on and the CSP stays report-only. `security.go` adds the client's inline-script hashes and the Vision analytics origin to `script-src`; `connect-src 'self'` holds because every icon is bundled (`src/lib/icons.ts` registers the Capsule mark). Web Crypto and Shiki's regex engine need no `unsafe-eval`. The client's pages answer `Referrer-Policy: no-referrer`; the key lives in the URL fragment, which no Referer carries. Keep it that way: Vision loads only on `/` and `/revoke` (a reveal path is the capsule id); never put the key in a path or query, never `fetch` or log `location.href`.
- Auth glue is `porte/kit` (`modules/auth/store.go`): `kit.Mount` mounts porte's routes, `kit.RequireAuth` guards `/auth/me`, and a `kit.Store` over `porte_users` resolves a session to its profile and erases accounts. Sign-in itself still writes through porte/pg's user store, which owns `facile_id` and `email_verified`; never pass the kit store as `oidc.Deps.Users`. Capsule has no `users` table of its own.
- GDPR: `modules/auth/deletions.go` polls Registre's deletions with the app's OIDC client credentials; `kit.Store.EraseAccount` deletes the person's `porte_users` row (identities, sessions, login codes and magic links cascade), then porte's `EraseIdentity` runs. Pastes carry no owner column and an attachment lives inside its paste's ciphertext, so nothing links a paste to an account and there is no paste to erase by person. The cursor lives in `deletion_state`. Run it in one replica only.
- UI strings live in `apps/client/src/lib/i18n/messages/` as `[en, fr, es, de]` rows. es and de are machine translated and need native review. The locale comes from the `fc_locale` cookie, then the browser.
- Rate limiting on paste creation is tronc's `middleware.Ratelimit` (`newCreateLimiter` in `main.go`): 30 per visitor, refilling at 30 per minute, keyed on the host of `RemoteAddr` after tronc's real-IP rewrite. Never key on a forwarded header or `ClientIP`.
- MCP (`modules/mcp`, `mcp.go`): `/mcp` on tronc/mcp with `capsule_get_status` and `capsule_revoke` only, because the server cannot decrypt. Never add a tool that seals, opens or returns content or ciphertext, and never store a delete token: the revoke checks it at request time and the pending action names the capsule by id. The confirm page routes take the cookie session only.
- Background cleanup goroutine runs on startup to purge expired pastes.
- The client proxies API requests through SvelteKit's server routes (`/api/[...path]`) to avoid CORS in dev. The API is `expose`d, never published, so it is never called cross-origin.
- The error envelope, JSON helpers, logger, request logging, CORS, panic recovery, `/health` and `/ready` all come from `tronc`. Do not reintroduce local copies — fix them upstream.
- `/health` and `/ready` answer at both `/` and `/api`, so the same probe works through the edge.
- The container healthcheck re-executes the binary: `["CMD", "/app/api", "healthcheck"]`.
- The schema lives in `apps/api/migrations/`, applied at boot by `tronc/migrate`. Do not reintroduce `AutoMigrate` outside tests — `scripts/check.sh` fails on it. Adding a model field without a migration means the column does not exist and the query fails at runtime.
- Production was baselined (version 1 recorded as applied, never run) because its schema predates migrations. A fresh database runs `00001_baseline.sql` for real.
