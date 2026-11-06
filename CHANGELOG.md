# Changelog

All notable changes to this project are documented here. The format is
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- porte v0.5.2 to v0.12.0, `porte/directory` v0.3.0 and tronc v0.16.0 to v0.17.1.
- `@facile/muse` v0.8.0 to v0.12.1. The page header and footer are muse `LandingHeader` and
  `LandingFooter`; muse's own icons render from bundled paths and the Capsule mark is bundled in
  `src/lib/icons.ts`, so no icon is fetched from api.iconify.design and `iconify-icon` is gone.
- Every UI string is in French, English, Spanish and German (es and de need native review). The
  language follows the `fc_locale` cookie, then the browser, and a switch sits in the header. The
  "sealed ago" badge uses `Intl.RelativeTimeFormat`.
- tronc's security headers are on by default. The CSP stays report-only, with the client's
  inline-script hashes and the Vision analytics origin in `script-src`. The client pages answer
  `Referrer-Policy: no-referrer`; the decryption key stays in the URL fragment, which browsers
  never send.
- The Vision analytics script loads `async`, so a slow analytics host cannot hold up the app.
- `mise.toml` pins Go 1.26, matching `go.mod` and the Dockerfile.
- CI: `.github/workflows/ci.yml` runs gofmt, vet, race tests on Postgres, svelte-check, the
  client tests, the build and filet. It replaces `filet.yml`.

### Removed

- The unused client `/docs` route. The API reference is served by the Go binary at `/docs`.

### Added

- Migration `00003_launch`: `porte_users.locale`, `porte_login_codes.provider` and `.subject`
  (porte v0.10 login-code binding), `porte_magic_links` (porte v0.12's erase deletes from it) and
  `deletion_state`.
- GDPR deletion: the app polls Registre for erased people with its OIDC client credentials and
  deletes their account, sessions and login codes. Pastes hold no personal data. Idempotent.
- `GET /version`, answering the `APP_VERSION` build arg (`dev` by default).
- Tests: an anonymous reveal needs no session, the headers keep the fragment private and allow no
  `unsafe-eval`, deletion erases and a second call is a no-op.

## [0.2.0] - 2026-11-03

### Added

- `OIDC_CLI_AUDIENCE` mounts porte's device exchange at `POST /api/auth/oidc/device/exchange`,
  so suite CLIs sign in with the browser approval link from Registre's device grant, over SSH
  too. Set it to `facile-cli`; unset, the route stays absent and answers 404.

## [0.1.0] - 2026-11-03

First tagged release.

### Added

- Single sign-on through porte and Registre, so the `capsule` CLI can sign in with one browser
  step (`capsule login`). Routes: `/api/auth/config`, `/api/auth/oidc`, `/api/auth/oidc/exchange`,
  `/api/auth/me`, `/api/auth/logout`. Pastes stay anonymous and end-to-end encrypted.
- Migration `00002_porte` for the identity, session and login-code tables.

### Changed

- `@facile/muse` v0.4.0 to v0.8.0.
- The page header and footer no longer draw border lines, matching muse's borderless surfaces.
