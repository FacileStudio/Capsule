# Changelog

All notable changes to this project are documented here. The format is
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Licence: Capsule is now licensed under `AGPL-3.0-only` (see `LICENSE`); before this, the repository carried no licence.
- build: tronc v0.17.1 to v0.21.0, porte v0.13.0 to v0.18.0 (`directory` v0.3.0 to v0.4.0), muse v0.16.1 to v0.18.0. Every entry in between is additive.

## [0.5.0] - 2026-11-08

### Added

- File encryption and sharing: attach a file up to 10 MB via the drag-and-drop zone (alone or alongside text), encrypted in the browser inside a backward-compatible `{"__capsule":1}` envelope.
- Reveal page renders image, audio and video previews with a download button showing the file name and size.

### Changed

- `MAX_PASTE_SIZE` defaults to 16 MiB (`16777216` bytes) and `POST /api/pastes` accepts bodies up to 32 MiB so encrypted files fit without hitting a 1 MiB body limit.

## [0.4.4] - 2026-11-07

### Changed

- build(client): bump muse to v0.16.1 and replace the deprecated `outline` button variant with `secondary`.
- build: `GET /version` is baked at build time from `apps/client/package.json` by `scripts/version.sh`. The `APP_VERSION` build arg is gone from the Dockerfile and compose file, so a stale Dokploy env can no longer show the wrong version. The response is unchanged.

## [0.4.3] - 2026-11-07

### Changed

- build(client): bump muse to v0.14.4 (from v0.14.2). Switching the language in the user menu re-renders the whole UI at once, with no reload. A single 100% DonutChart slice draws as a full ring. The Subscription row no longer passes an empty `plan`.

## [0.4.2] - 2026-11-07
### Security

- Registre tokens are sealed at rest. porte v0.12.1 to v0.13.0: `porte_identities.access_token` and `refresh_token` are encrypted with AES-256-GCM under `PORTE_TOKEN_KEY`. Existing rows are sealed once, in the background, after boot, and the count is logged. An invalid `PORTE_TOKEN_KEY` stops the boot. With the key unset nothing changes.
- Losing `PORTE_TOKEN_KEY` makes the stored tokens unreadable: people sign in again, nothing else is lost. Do not roll back to a build older than porte v0.13.0 after the first reseal.

## [0.4.1] - 2026-11-07

### Changed

- build(client): bump muse to v0.14.2. The user menu clamps itself inside the viewport, so the class override on `UserMenu` is gone.
- The user menu gains a Subscription row that opens the plans on Facile (`https://sso.facile.studio/spaces`).

### Fixed

- `GET /auth/me` returned an empty email for cookie sessions. It now reads the name, email and picture from the porte user row, like `GET /auth/session`, and also returns `avatar_url`.

## [0.4.0] - 2026-11-07

### Added

- `GET /api/auth/session` answers 200 with `{"user": null}` for an anonymous visitor and the name, email and Registre picture for a signed-in one, so the page can decide whether to show the user menu without a failed request.
- Signed-in people get muse's `UserMenu` (Account, Language, Theme, Help, Sign out) in the header, in place of the language and theme controls. Capsule has no settings page, so the Settings row is hidden.

### Changed

- build(client): bump muse to v0.14.1.
- Anonymous visitors see the same header and the same seal, reveal and revoke flows as before; the only difference is one extra `GET /api/auth/session` per page load.

## [0.3.1] - 2026-11-07

### Changed

- build(api): bump porte to v0.12.1. The GDPR erasure no longer fails on a schema without `porte_magic_links`.
- build(client): bump muse to v0.13.0. Form controls are borderless fills, `Select` is muse's own combobox and badges stay on one line.
- test(client): the e2e specs open the `Select` combobox and read or click its options instead of using `selectOption`.

## [0.3.0] - 2026-11-06

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
- The Vision analytics script loads `async` and only on the landing and revoke pages. A reveal page's path is the capsule id, which Vision would record.
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
