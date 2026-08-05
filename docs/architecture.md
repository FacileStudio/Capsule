# Capsule — Architecture

The crypto boundary stated exactly, then how a request moves through the system and what the
database holds.

## Runtime topology

```
Internet ──▶ Traefik ──▶ Go binary (:4000) ──┬──▶ /health, /ready   liveness + readiness
                                              ├──▶ /api/health, /ready
                                              ├──▶ /api/pastes/*    paste handlers
                                              ├──▶ /docs            OpenAPI + Scalar UI
                                              └──▶ /*               SPA catch-all
                                                              │
                                                        Postgres 16
                                                              │
                                              Journal (log shipping, optional)
```

One container, one router, one hostname. The client is built with `@sveltejs/adapter-static`
into a fallback SPA and copied into the API image at `/client`; `tronc/spa` serves it as the
catch-all. The browser calls `/api/*` on the same origin, so no request is ever cross-origin
and the CORS allow-list can stay empty.

There is **no server-side rendering and no server-side proxy**. Every page is client-rendered.
That is not a style preference: server-rendering a page whose URL fragment carries the
decryption key would put that key one templating mistake away from the server.

## The crypto boundary

All cryptography lives in `apps/client/src/lib/crypto.ts`. There is no crypto in the Go code
beyond `crypto/rand` for IDs and tokens and `crypto/subtle` for comparing the delete token.

### Sealing

1. `crypto.subtle.generateKey` produces a 256-bit AES-GCM key in the browser.
2. `encrypt` draws a fresh 12-byte IV from `crypto.getRandomValues`, encrypts the UTF-8
   plaintext, and returns standard base64 of `IV || ciphertext || tag`.
3. Without a password, the raw key is base64url-encoded and becomes the whole URL fragment.
4. With a password, `wrapContentKey` derives a wrapping key with PBKDF2-SHA256 — 600,000
   iterations, a fresh 32-byte salt — encrypts the raw content key under it with a fresh
   12-byte IV, and the fragment becomes `encryptedKey.salt.iv`, all base64url.
5. Only the base64 ciphertext is POSTed. The link is `https://<host>/<id>#<fragment>`.

### Opening

The reveal page reads `window.location.hash`, fetches metadata, asks for a password if the
capsule has one, unwraps or imports the key, requests the content, and decrypts in the
browser. `importKey` and `unwrapContentKey` both mark the resulting key non-extractable and
`['decrypt']`-only.

### What the server can see

Everything below is stored in plaintext in the `pastes` table or written to logs. Treat it
as visible to whoever runs the server.

| Item | Where |
|---|---|
| The ciphertext, and therefore the approximate plaintext **length** — AES-GCM adds no padding | `pastes.content` |
| Which syntax the author picked, e.g. `sql`, `bash` | `pastes.syntax` |
| Whether a password is set | `pastes.has_password` |
| Burn-after-read flag, expiry timestamp, max views, view count, creation time | `pastes` |
| The paste ID and the delete token, both generated server-side | `pastes` |
| Request metadata — method, path, status, `remote_addr` and `client_ip` | request log, shipped to Journal when it is configured |

### What the server cannot see

| Item | Why |
|---|---|
| The plaintext | Never leaves the browser unencrypted |
| The AES content key | Lives only in the URL fragment, which browsers do not send in the request line, `Referer`, or anywhere else |
| The password | Never transmitted; only used locally to derive a wrapping key |
| The wrapped key, its salt and its IV | Also fragment-only |

### The limits of that claim

Be honest about these. They are properties of browser-delivered end-to-end encryption, not
bugs, but "zero knowledge" without them is a half-truth.

- **The server ships the JavaScript that does the encryption.** A compromised server, or a
  malicious operator, can serve modified `crypto.ts` that exfiltrates the plaintext or a
  weakened key. Encryption in the browser protects against passive storage compromise and
  against a curious operator reading the database — it does not protect against an operator
  who is actively hostile at the moment you use the app. Self-hosting is what makes that
  threat model acceptable: you are the operator.
- **Password mode is offline-attackable by anyone who has the link.** The wrapped key, salt
  and IV all sit in the fragment. Someone holding the full URL can brute-force the password
  locally at PBKDF2-SHA256 600,000 iterations per guess. The password buys you "link alone
  is not enough", not "link plus weak password is safe".
- **Anyone with the ID can destroy a capsule without reading it.**
  `POST /api/pastes/{id}/content` is unauthenticated: it returns the ciphertext and, on a
  burn-after-read capsule, burns it. No key needed. Guessing an ID is a 64-bit search, so
  this is not a realistic attack, but a leaked ID is enough to deny the recipient the
  secret.
- **Burning blanks the content; it does not delete the row.** Metadata survives until the
  purge, 30 days after creation.
- **IP addresses appear in the request log.** The tronc request logger records `remote_addr`
  and `client_ip` on every request, and ships them to Journal when Journal is configured.
  The content is opaque; who fetched it and when is not.
- **The rate limiter keys on `r.RemoteAddr`.** Behind Traefik that is the proxy's address,
  not the caller's, so the 30-per-minute create limit behaves as a global limit rather than
  a per-client one. It does not weaken the encryption; it does change what the limit means.

## Request lifecycle

1. `httpx.NewRouter` (tronc) applies request logging, panic recovery and CORS. The allowed
   headers list is the tronc default plus `X-Delete-Token`.
2. `health.Mount` registers `/health` and `/ready` twice, bare and under `/api`, with a
   database ping as the readiness check.
3. `/api/pastes` routes into the paste handlers. Only creation is rate-limited.
4. `/docs` serves the Scalar UI and `/docs/openapi.yaml` the spec.
5. Anything unmatched falls through to the SPA.

Server timeouts are generous relative to the rest of the suite, because a 1 MB paste has to
finish uploading: 5s read header, 30s read, 60s write, 120s idle. Shutdown is graceful with
a 10s budget on `SIGINT`/`SIGTERM`, and cancels the cleanup goroutine first.

## Data model

One table.

| Column | Type | Notes |
|---|---|---|
| `id` | `varchar(24)` primary key | `cap_` plus 8 random bytes hex — 64 bits of entropy |
| `content` | `text` | The base64 ciphertext. Never serialized into a JSON response by the model — the `json:"-"` tag makes leaking it an explicit act |
| `burn_after_read` | bool, default true | |
| `expires_at` | nullable timestamp | Null means no expiry |
| `max_views` | nullable int | An alternative to burn-after-read |
| `view_count` | int, default 0 | |
| `has_password` | bool, default false | Only tells the reveal page which flow to run |
| `syntax` | `varchar(50)` | Shiki language label |
| `delete_token` | `varchar(64)` unique | 32 random bytes hex, `json:"-"` |
| `burned` | bool, default false | |
| `created_at` | timestamp | |

`schemas.Migrate` is GORM `AutoMigrate` over this single model.

## Burning

`GetContent` runs inside a transaction that takes `SELECT … FOR UPDATE` on the row, so two
simultaneous readers cannot both receive the ciphertext of a burn-after-read capsule. Inside
the lock it increments `view_count`, and burns when `burn_after_read` is set or when
`view_count + 1` reaches `max_views`. Burning sets `burned = true` and `content = ''` in the
same update.

`Revoke` compares the supplied `X-Delete-Token` against the stored one with
`subtle.ConstantTimeCompare`, then burns the same way. The delete token is returned exactly
once, in the create response, and is never persisted client-side — closing the tab without
copying it means the capsule can no longer be revoked, only left to expire.

A cleanup goroutine started at boot ticks every five minutes and does two things: burns
capsules whose `expires_at` has passed and are not yet burned, and hard-deletes burned rows
created more than 30 days ago. Expiry is enforced at read time as well, so a stopped
cleanup goroutine cannot serve an expired capsule.

## Cross-app integration

Capsule is deliberately outside the suite's event bus. It has **no users, no sessions and no
authentication at all**, so it cannot supply the `actor_email` the `pool`/`enveloppe` event
contract keys on. Nothing about a capsule is attributable to a person by design, and adding
an actor identity would undo that.

The one integration is **Journal**: with both `JOURNAL_URL` and `JOURNAL_TOKEN` set, the
tronc logger is wrapped in `journal.NewHandler` and structured log lines ship to Journal —
including the `client_ip` field noted above.
