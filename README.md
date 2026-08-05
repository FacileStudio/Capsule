# Capsule

Self-hosted encrypted paste. Write a secret, get a link, the content burns after it is read.

Content is encrypted in the browser with AES-256-GCM before it is sent. The key travels in
the URL fragment, which browsers never transmit, so the server stores ciphertext it has no
way to open. What the server does and does not see is spelled out precisely in
[docs/architecture.md](docs/architecture.md) — read it before relying on the guarantee.

Live at [capsule.facile.studio](https://capsule.facile.studio).

## What it does

- Encrypts in the browser with AES-256-GCM through the Web Crypto API, key generated
  client-side
- Puts the key in the URL fragment so it never reaches the server
- Burns after reading by default, in a single locking transaction, so a race cannot serve
  the same secret twice
- Expires on a schedule instead — 1 hour, 24 hours, 7 days or 30 days
- Optionally wraps the content key with a password through PBKDF2-SHA256, 600,000 iterations
- Returns a delete token so the author can revoke a capsule before anyone opens it
- Syntax-highlights the revealed text with Shiki, in the reader's browser
- Purges burned rows entirely 30 days after creation, on a background ticker

## Stack

| Layer | Tech |
|---|---|
| API | Go 1.24, Chi v5, GORM, PostgreSQL 16, [tronc](https://github.com/FacileStudio/tronc) v0.6.0 |
| Client | SvelteKit 2 (Svelte 5 runes), Tailwind CSS 4, Shiki, `adapter-static` |
| Crypto | Web Crypto API — AES-256-GCM, PBKDF2-SHA256 for password mode. Browser only |
| Deploy | Docker Compose, single distroless container behind Traefik |

## Quick start

```sh
cp .env.example .env
docker compose up -d --build
```

Compose starts `capsule-db` and `capsule-api`. The API listens on `4000` and serves both
`/api/*` and the built SPA. Migrations run on startup.

### Local development

```sh
mise run install
cd apps/api    && DATABASE_URL=… go run .    # :8080 unless PORT says otherwise
cd apps/client && bun run dev                # :5173
```

`DATABASE_URL` is required — the API refuses to start without it.

## Configuration

| Variable | What it does |
|---|---|
| `DATABASE_URL` | Postgres connection string. Required; no default |
| `PORT` | HTTP listen port, `8080` by default, `4000` in Compose |
| `MAX_PASTE_SIZE` | Ciphertext size cap in bytes, `1048576` by default |
| `LOG_LEVEL` | `debug`, `info`, `warn` or `error` |
| `CLIENT_DIR` | Directory holding the built SPA the binary serves |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allow-list. Unset denies every cross-origin caller, which is correct here: the SPA is same-origin |

Full reference: [docs/configuration.md](docs/configuration.md).

## Structure

```
apps/
  api/       Go backend — modules/pastes (the whole domain), modules/docs,
             internal/ (env, cleanup, database, middleware), schemas/
  client/    SvelteKit SPA — lib/crypto.ts is the entire crypto surface,
             built by adapter-static and served by the API binary
docs/        Architecture, configuration, development, deployment, API
scripts/     check.sh, the repository quality gate
```

## Documentation

| Doc | What's in it |
|---|---|
| [Architecture](docs/architecture.md) | The crypto boundary, request flow, data model |
| [Configuration](docs/configuration.md) | Every environment variable and default |
| [Development](docs/development.md) | Local setup, tests, the quality gate |
| [Deployment](docs/deployment.md) | Docker Compose, Dokploy, Traefik routing |
| [API](docs/api.md) | HTTP endpoints and payloads |

---

Part of the [Facile Suite](https://facile.studio) — self-hosted tools for creative studios
and freelancers. One login, zero cloud dependency.
