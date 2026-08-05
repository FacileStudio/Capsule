# Capsule — Deployment

How the image is built, what Compose brings up, how Traefik routes to it, and what an
upgrade risks.

## The image

`Dockerfile` is a three-stage build producing one distroless image.

1. `oven/bun:1` installs the client dependencies from the lockfile and runs `bun run build`,
   which `adapter-static` turns into a fallback SPA.
2. `golang:1.24-alpine` downloads modules, then builds a static binary with
   `CGO_ENABLED=0`, `-trimpath` and `-ldflags="-s -w"`.
3. `gcr.io/distroless/static-debian12:nonroot` receives the binary at `/api`, the built SPA
   at `/client`, `ENV CLIENT_DIR=/client`, and runs as `nonroot:nonroot`.

That `ENV` line is load-bearing. The `:nonroot` base sets `WorkingDir=/home/nonroot`, so a
relative `./client` would resolve there and the SPA would silently not be served — the API
would answer `/api/*` correctly while every page a human types 404s.

The image exposes `4000` and declares no volumes: all state is in Postgres.

## Compose topology

Two services, matching the suite's one-container/one-router/one-hostname rule.

```
dokploy-network ──▶ traefik ──▶ capsule-api:4000 ──▶ capsule-db:5432
                                                            │
                                                       db_data volume
```

| Service | Image | Notes |
|---|---|---|
| `capsule-db` | `postgres:16-alpine` | `expose: 5432` only, `pg_isready` healthcheck every 5s |
| `capsule-api` | built from `Dockerfile` | `expose: 4000`, no published port; waits for the database to be healthy |

Both services set `container_name`, which means **only one Capsule stack can run per Docker
host** — a second `docker compose up` with the same file collides on the name.

The one named volume is `capsule_db_data`. Losing it loses every unburned capsule, which is
a smaller disaster here than in most apps: capsules are ephemeral by design, and the oldest
useful row is at most 30 days old.

Database credentials are hardcoded in the Compose file (`capsule` / `capsule-internal-db`).
The database is never published, so this is a Compose-network credential rather than an
internet-facing secret. Change it if untrusted containers share the host network.

## Traefik

Labels on `capsule-api` declare two routers and one service, all on hostname
`capsule.facile.studio`:

| Router | Entrypoint | Behavior |
|---|---|---|
| `capsule-web` | `web` | Redirects to HTTPS through the `redirect-to-https@file` middleware |
| `capsule-secure` | `websecure` | TLS with the `letsencrypt` cert resolver |

Both point at `capsule-svc`, load-balancing to port `4000`. `traefik.docker.network` is
`dokploy-network`, declared external — Dokploy owns it.

One hostname serves everything. That is not only tidiness: it is what lets the CORS
allow-list stay empty, because the browser never makes a cross-origin request.

## Healthchecks

The container healthcheck is `/api healthcheck`, the binary invoking itself — `tronc`'s
`healthcheck.Handle` intercepts `os.Args` before anything else in `main`, which is the only
way to health-check a distroless image with no shell and no `curl`.

Over HTTP, `tronc/health` mounts `/health` and `/ready`, plus `/api/health` and `/api/ready`.
`/ready` pings the database; `/health` does not.

A green `/api/health` says the process is up. It says nothing about the SPA being served — if
`CLIENT_DIR` is wrong, health stays green while every page 404s. Check a real page.

## Deploying to la ruche

Deployment is managed through Dokploy at `gare.facile.studio`, which builds from the repo and
runs the Compose file. Prefer the `dokploy` CLI over SSH and raw `docker`.

Environment comes from Dokploy's environment editor. `.env.example` is the template of what
belongs there. In practice production needs almost nothing beyond what Compose already sets:
`DATABASE_URL`, `PORT` and `LOG_LEVEL` are inline, so only Journal credentials are worth
adding.

## Migrations

There is no migration step. `schemas.Migrate` runs GORM `AutoMigrate` over the single `Paste`
model at every boot.

- Adding a struct field adds a column. Removing one leaves the column behind, unread.
- Two instances booting at once both migrate. Deploy one at a time.
- Rolling back the image rolls back the binary and the SPA together, since they ship in the
  same image. The database does not roll back, which is safe as long as changes stay
  additive.

## Scaling, and why you should not

Capsule is stateless apart from Postgres, so more than one replica would work correctly for
reads and writes — the burn transaction takes `SELECT … FOR UPDATE`, so even a race across
replicas is safe.

Two things are per-process, though. The rate limiter is an in-memory map, so N replicas means
N times the effective create limit. And the cleanup goroutine runs in every replica, doing
the same work redundantly; harmless, but pointless. Vertical is the right direction here.

## What an upgrade risks

The failure that matters is the SPA not being served — a wrong `CLIENT_DIR`, or a client
build that failed while the Go build succeeded. Health stays green through both. After any
deploy, load the root page and seal a throwaway capsule end to end.

The second thing to watch is the cleanup goroutine. If it stops, expired capsules are still
refused at read time — expiry is checked in `GetContent`, not only by the sweeper — but
burned rows stop being purged and the table grows without bound. `cleanup: burned expired
pastes` and `cleanup: purged old burned pastes` are the log lines that prove it is alive.
