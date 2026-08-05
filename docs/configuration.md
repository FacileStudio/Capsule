# Capsule — Configuration

Every environment variable the API actually reads, taken from `apps/api/internal/env/env.go`
and the `tronc/env` and `tronc/spa` packages it builds on.

Capsule has a small configuration surface on purpose. It has no users, no sessions, no mail,
no OIDC and no object storage, so there is nothing to configure for any of them.

## Core

These come from `troncenv.LoadCore()`, shared by every Go app in the suite.

| Variable | Required | Default | What it does |
|---|---|---|---|
| `DATABASE_URL` | yes | — | Postgres connection string. `LoadCore` errors without it and the process exits |
| `PORT` | no | `8080` | HTTP listen port. Must parse as an integer in 1–65535 or startup fails |
| `LOG_LEVEL` | no | `info` | `debug`, `info`, `warn` or `error`. API requests log at info; health probes and static assets at debug |
| `APP_ENV` | no | `development` | `development`, `staging` or `production`. Never gates security behavior |
| `CORS_ALLOWED_ORIGINS` | no | — | Comma-separated allowed origins. Unset allows no cross-origin caller |
| `JOURNAL_URL` | no | — | Journal ingest URL. Shipping needs both this and the token |
| `JOURNAL_TOKEN` | no | — | Per-app Journal key |

`CORS_ALLOWED_ORIGINS` is the canonical name; `tronc/env` falls back in order to
`ALLOWED_ORIGINS`, `DOMAINS`, `DOMAIN`, `CORS_ORIGINS`, `TRUSTED_ORIGINS` and
`CLIENT_ORIGIN`.

**Leaving CORS unset is the correct configuration.** One binary serves the API and the SPA on
one hostname, so the browser never makes a cross-origin request. Opening the allow-list only
lets other origins script against your capsules.

Note that enabling Journal ships the request log, which carries `client_ip` and `user_agent`
for every `/api/*` call. Capsule content stays opaque; access patterns do not. Decide
accordingly.

## Capsule

| Variable | Required | Default | What it does |
|---|---|---|---|
| `MAX_PASTE_SIZE` | no | `1048576` | Maximum accepted `content` length in bytes. Must be a positive integer or startup fails |
| `CLIENT_DIR` | no | `./client` | Directory holding the built SPA. Read by `tronc/spa`. The Dockerfile pins `/client` explicitly, because the `:nonroot` distroless base sets `WorkingDir=/home/nonroot` and a relative path would resolve there — the API would answer `/api/*` while every page 404s |

`MAX_PASTE_SIZE` caps the **ciphertext**, not the plaintext. Base64 of `IV || ciphertext ||
tag` is roughly 4/3 of the plaintext plus 28 bytes, so the default 1 MiB accepts around
760 KB of text.

## Compose

`docker-compose.yml` hardcodes the database credentials rather than reading them from the
environment: user `capsule`, password `capsule-internal-db`, database `capsule`, and
`DATABASE_URL` set to
`postgres://capsule:capsule-internal-db@capsule-db:5432/capsule?sslmode=disable`. The
database service only `expose`s 5432 on the Compose network and is never published, so that
password is a network-internal credential, not a secret guarding the internet. Change it if
your host runs untrusted containers on the same network.

Only `LOG_LEVEL`, `JOURNAL_URL` and `JOURNAL_TOKEN` are read from the environment by the
Compose file.

## Variables that no longer exist

`README.md` and `CLAUDE.md` historically documented `ORIGIN` and `API_URL`, from when the
SvelteKit server ran as its own container and proxied to the API. It does not any more: the
client is built with `adapter-static` and served by the Go binary. Neither variable is read
by anything in this repo.

## The failure mode worth knowing

Configuration errors are fatal. `env.Load` returns an error, `main` logs
`failed to load config` and returns. A container restarting in a loop with one log line is
almost always a missing `DATABASE_URL`, a malformed `PORT`, or a non-numeric or negative
`MAX_PASTE_SIZE`.
