# Capsule — Development

Local setup, the two processes, the four test suites, and the quality gate that guards
pushes.

## Prerequisites

| Tool | Version | Why |
|---|---|---|
| Go | 1.24 | `apps/api/go.mod` declares `go 1.24.0`; `mise.toml` pins the toolchain |
| Bun | 1.x | Client package manager and build runner |
| PostgreSQL | 16 | What Compose and production run |
| mise | any | Task runner for `install`, `check`, `format`, `hooks` |
| Docker | any | Only if you want the Compose database instead of a local one |

## Setup

```sh
mise run install       # bun install --frozen-lockfile in apps/client
mise run hooks         # git config core.hooksPath .githooks
```

Bring up a database:

```sh
docker compose up -d capsule-db
```

## Running

Two processes, two terminals.

```sh
cd apps/api
DATABASE_URL=postgres://capsule:capsule-internal-db@localhost:5432/capsule?sslmode=disable go run .
```

Note that the Compose database only `expose`s 5432 on the Compose network, so connecting to
it from the host means publishing the port yourself or running a local Postgres instead.

```sh
cd apps/client
bun run dev            # Vite on http://localhost:5173
```

In this mode the API serves no SPA — `tronc/spa` only mounts the catch-all when `CLIENT_DIR`
actually contains a build. Migrations run on every API start through `schemas.Migrate`, which
is GORM `AutoMigrate` over the single `Paste` model.

**There is no dev proxy.** `lib/backend.ts` hardcodes a same-origin `/api` base, and neither
`vite.config.ts` nor a `hooks.server.ts` forwards it anywhere — the repo has no
`hooks.server.ts` at all, which is deliberate. So `bun run dev` on its own renders the UI but
cannot reach the API. For a full local round trip, build the client and let the binary serve
it:

```sh
cd apps/client && bun run build
cd ../api && CLIENT_DIR=../client/build DATABASE_URL=… go run .
```

Then use the API's port, not `:5173`.

## Tests

Four suites, three runners.

```sh
cd apps/api    && go test ./...        # unit tests on in-memory SQLite
cd apps/client && bun run test         # vitest, crypto round-trip
cd apps/client && bun run test:e2e     # playwright, starts the dev server itself
```

- `apps/api/modules/pastes/` covers the service and the handlers. Most of it runs against
  an in-memory SQLite database. One test needs real Postgres for `SELECT … FOR UPDATE`
  semantics and **skips itself** unless `TEST_DATABASE_URL` is set. If you touch the burn
  transaction, set it and run the suite again — SQLite will not catch a locking regression.
- `apps/api/internal/cleanup/` covers the expiry burn and the 30-day purge.
- `apps/client/src/lib/crypto.test.ts` is the unit test for the crypto surface. Anything
  that touches `crypto.ts` needs a test here; a silent crypto regression is the one bug in
  this repo that nobody would notice from the UI, because a wrong key looks exactly like a
  corrupted link.
- `apps/client/tests/capsule.spec.ts` drives Chromium over 14 tests. Playwright starts
  `bun run dev` itself and reuses a running server outside CI. Because the dev server has no
  API behind it, these cover rendering and form state, not a real seal-and-reveal round
  trip. Point `PLAYWRIGHT_BASE_URL` at a running full stack if you want that.

## The quality gate

`scripts/check.sh` is the gate. It depends on nothing but `go` and, for the client half,
`bun`. It is a shell script rather than a mise task body on purpose: `mise run` resolves
every tool in the merged config before running anything, so one broken tool in your global
config would otherwise take the gate down with it.

```sh
sh scripts/check.sh              # gofmt -l, go vet, go test, then the client type-check
sh scripts/check.sh --go-only    # Go half only
sh scripts/check.sh --format     # rewrite Go sources in place
```

Equivalent mise tasks: `mise run check`, `mise run check-go`, `mise run format`.

The script resolves `go` and `gofmt` from `GOROOT` when it is set. mise exports `GOROOT` for
the pinned version but leaves an unrelated `go` earlier on `PATH`, and mixing the two
produces `compile: version "X" does not match go tool version "Y"`.

Neither `vitest` nor Playwright is part of the gate. Run them yourself when you touch the
client.

## The pre-push hook

`.githooks/pre-push` runs `scripts/check.sh --go-only`, not the full gate. Capsule's client
carries pre-existing `svelte-check` errors, so gating on the full check would block every
push until those are cleared. Run `sh scripts/check.sh` yourself for the whole picture. When
the client reaches zero errors, drop the `--go-only`.

Bypass once with `git push --no-verify`.

## Conventions

- The Go layout is `internal/` for infrastructure and `modules/` for domain code. `pastes`
  is the only domain module and splits into `router.go`, `handler.go`, `service.go`,
  `types.go`.
- Handlers return through `httpjson.WriteJSON` and `httpjson.WriteError` from tronc, and the
  error envelope, logger, request logging, CORS, panic recovery, `/health` and `/ready` all
  come from tronc. Do not reintroduce local copies — fix them upstream.
- `schemas.Paste` tags `Content` and `DeleteToken` as `json:"-"`. Keep it that way:
  serializing the model directly is then safe by default, and leaking either becomes an
  explicit act rather than an oversight.
- The client is Svelte 5 runes only, enforced through `dynamicCompileOptions` in
  `svelte.config.js`.
- **Never add server-side rendering or a server-side proxy to the client.** The decryption
  key lives in the URL fragment; SSR is the one change that could put it on the server.
