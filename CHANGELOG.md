# Changelog

All notable changes to this project are documented here. The format is
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
