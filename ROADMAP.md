# Capsule Roadmap

Zero-knowledge encrypted sharing. Every feature reinforces one principle:
**the server never knows what you shared.**

---

## Phase 1 — Complete the Foundation

_Ship what's half-built. Fill the gaps that make the MVP feel unfinished._

- [x] **Password Protection** — PBKDF2 key derivation (600k iterations),
  wraps content key with password-derived key. Server sees nothing.
- [x] **Dark Mode** — Toggle respects `prefers-color-scheme`, persists in
  `localStorage`. Toggles `dark` class on root.
- [x] **Mobile Polish** — `100dvh` viewport, `min-h-[44px]` touch targets,
  responsive button sizing (`h-12` mobile → `h-10`/`h-11` desktop).
- [x] **Test Suite** — Go unit + integration tests (service, handler,
  cleanup), Playwright E2E (seal → reveal → burn), crypto round-trip tests.
- [x] **CI Pipeline** — GitHub Actions: Go lint/build/test, client build +
  typecheck, Docker Compose validation. Runs on push to main and all PRs.

---

## Phase 2 — Developer Experience

_Make Capsule a tool developers reach for instinctively._

- [x] **CLI (`capsule`)** — `capsule seal`, `capsule reveal`, `capsule revoke`.
  Go binary, AES-256-GCM native (`crypto/aes` + `crypto/cipher`), pipe-friendly.
  Ship via GitHub Releases + Homebrew tap. → [FacileStudio/capsule-cli](https://github.com/FacileStudio/capsule-cli)
- [ ] **API Documentation** — OpenAPI 3.1 spec from Go handlers. Serve Scalar
  or Swagger UI at `/docs`. Include `curl` examples for every endpoint.
- [ ] **Self-Host Guide** — One-command deploy templates for Fly.io, Railway,
  Coolify, and bare `docker compose up`. Document every env var, TLS, and
  PostgreSQL requirements. _(docker-compose.yml exists, README has basics)_
- [ ] **Webhooks** — Optional callback URL on creation. Events: `paste.read`,
  `paste.burned`, `paste.expired`. HMAC-signed payloads, never content.
- [x] **Syntax Highlighting on Reveal** — Shiki with dual-theme (github-dark /
  github-light), lazy-loaded, all 12 syntax options supported.

---

## Phase 3 — Power Features

_Stretch beyond basic paste sharing into territory competitors don't cover._

- [ ] **Paste Collections (Bundles)** — Group multiple secrets under one link.
  Each item encrypted independently, reveal shows a list with copy buttons.
  One burn destroys the whole bundle.
- [ ] **QR Code Sharing** — Client-side QR on seal confirmation screen.
  Include a "print" layout for physical handoff.
- [ ] **Read Receipts** — Opt-in. Notify sender via email, webhook, or status
  page. Confirms *that* it was read, never *what*.
- [ ] **IP / Region Restrictions** — IP allowlist (CIDR), country allowlist
  (GeoIP at edge). Server enforces before returning ciphertext.
- [ ] **Custom Expiration** — Custom duration or specific datetime. Support
  very short TTLs (5 min) for real-time credential handoff.
- [ ] **Paste History (Local Only)** — `localStorage`-encrypted log of paste
  IDs + delete tokens. Revoke from a local dashboard. Never touches server.

---

## Phase 4 — Teams & Compliance

_Enterprise-adjacent features for teams that need audit trails without
sacrificing zero-knowledge._

- [ ] **Team Namespaces** — Shared namespace (`team.capsule.dev/...`), admin
  defaults (max expiry, mandatory burn, IP ranges). Admins see metadata only.
- [ ] **Audit Log** — Append-only event log per namespace. Exportable
  CSV/JSON. Zero plaintext stored.
- [ ] **SSO / SAML** — Okta, Google Workspace, Azure AD. Gates team dashboard
  and audit log, not individual paste sharing.
- [ ] **Custom Domains** — CNAME + auto TLS. Branded experience, same
  zero-knowledge guarantees.
- [ ] **Data Residency** — Region selection for self-hosters and hosted
  version. Document region-specific deployment.

---

## Phase 5 — Ecosystem

_Meet users where they already are._

- [ ] **Browser Extension** — Right-click → "Seal with Capsule" → encrypted
  link on clipboard. Chrome + Firefox, local encryption, no persistent perms.
- [ ] **Slack Integration** — `/capsule seal <secret>` → ephemeral response
  with encrypted link. Secret never in Slack history.
- [ ] **GitHub Action** — `facile/capsule-action@v1`. Outputs a Capsule URL
  for one-time credentials in CI logs.
- [ ] **SDK / Client Libraries** — `@facile/capsule` (TS/JS), `capsule-go`,
  `capsule-py`. Encryption + API calls in one package.
- [ ] **Raycast / Alfred Extension** — Keyboard shortcut → paste content →
  encrypted link. Local encryption, zero-knowledge.

---

## Non-Goals

Things Capsule deliberately does not do:

- **User accounts for basic sharing** — anonymous by default, always
- **Server-side search** — we can't search what we can't read
- **Paste editing after creation** — immutable by design; revoke and re-create
- **Long-term storage** — max 30 days; use a vault for permanent secrets
- **Rich text editing** — plain text and code; not a document editor
- **File sharing** — text-only is a feature; files bring storage costs, abuse vectors, and legal liability that don't fit a paste tool

---

## Versioning Note

This roadmap reflects current thinking as of 2025-05. Phases are directional,
not commitments. Features may shift between phases as priorities evolve.
