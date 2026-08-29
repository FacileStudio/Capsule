# Capsule — API

Every HTTP route the Go binary registers, generated from `apps/api/modules/`.

There is **no authentication of any kind** — no users, no sessions, no tokens beyond the
per-capsule delete token. Possession of a capsule ID is what authorizes reading it, and
possession of its delete token is what authorizes revoking it. Responses are JSON via
`httpjson.WriteJSON`, and errors share the tronc envelope:

```json
{ "error": { "code": "not_found", "message": "paste not found" } }
```

A live OpenAPI 3.1 spec is served at `GET /docs/openapi.json`, with a Scalar UI at `GET /docs` —
the same two paths as every other Facile backend. A test walks the live router and asserts that
every route is registered with complete typed request and response schemas.

## Health

`/health` and `/ready` are mounted twice, bare and under `/api`. `/ready` pings the database;
`/health` does not.

## Create a capsule

```
POST /api/pastes
```

Rate-limited to 30 requests per minute. `content` must already be ciphertext — the server
does not encrypt anything and will happily store whatever you send, including plaintext.

| Field | Type | Notes |
|---|---|---|
| `content` | string | Required, trimmed, non-empty. Rejected over `MAX_PASTE_SIZE` (default 1048576 bytes) with `resource_exhausted` |
| `burn_after_read` | bool | Defaults to `true` when omitted |
| `expires_in` | string | One of `1h`, `24h`, `7d`, `30d`. Anything else is a 400. Omitted means no expiry |
| `max_views` | int | Optional alternative to burn-after-read |
| `has_password` | bool | Metadata only — it tells the reveal page which flow to run. The server never sees the password |
| `syntax` | string | Shiki language label, stored in the clear |

```json
{
  "id": "cap_1f4a9c3e5b7d8a02",
  "delete_token": "…64 hex chars…",
  "expires_at": "2026-01-02T12:00:00Z",
  "created_at": "2026-01-01T12:00:00Z"
}
```

`201 Created`. The ID is `cap_` plus 8 random bytes as hex — 64 bits of entropy. The delete
token is 32 random bytes as hex, returned **exactly once**; it is never retrievable
afterwards, and the client does not persist it. Save it or lose the ability to revoke.

The decryption key is not in this response and never touches the server. The client appends
it to the URL as a fragment: `https://<host>/<id>#<key>`.

## Read metadata

```
GET /api/pastes/{id}
```

Safe to poll: it does not increment the view count and does not burn anything.

```json
{
  "id": "cap_1f4a9c3e5b7d8a02",
  "exists": true,
  "has_password": false,
  "syntax": "sql",
  "expires_at": "2026-01-02T12:00:00Z",
  "created_at": "2026-01-01T12:00:00Z"
}
```

A missing, burned or expired capsule returns `200` with `{ "id": "…", "exists": false }`, not
a 404. The three cases are deliberately indistinguishable — knowing that an ID *was* once
valid is itself information.

## Read content

```
POST /api/pastes/{id}/content
```

`POST`, not `GET`, because it mutates: this is the call that burns. No body, no auth.

```json
{ "content": "…base64 of IV || ciphertext || tag…" }
```

Inside one transaction holding `SELECT … FOR UPDATE` on the row, the server increments
`view_count` and then burns when `burn_after_read` is set, or when `view_count + 1` reaches
`max_views`. Burning sets `burned = true` and blanks `content`. Two simultaneous readers
therefore cannot both receive a burn-after-read capsule.

`404 not_found` when the capsule is missing, already burned, or past `expires_at` — expiry is
enforced here as well as by the background sweeper.

**Anyone holding the ID can call this**, with or without the decryption key, and on a
burn-after-read capsule that destroys the content. The returned ciphertext is useless without
the key, so this is a denial-of-read, not a disclosure.

## Revoke

```
DELETE /api/pastes/{id}
X-Delete-Token: <token from the create response>
```

Returns `{ "deleted": true }`. A missing header is `403`, and a wrong token is `403` after a
constant-time comparison. Revoking burns the capsule exactly as reading it would: `burned`
set, `content` blanked, row retained.

`X-Delete-Token` is added to the CORS allowed-headers list on top of the tronc defaults, so
the header survives a preflight if you ever do open the allow-list.

## Rate limiting

Only creation is limited: 30 requests per minute, keyed on `r.RemoteAddr`, answered with
`429` and a `Retry-After` header. Behind Traefik, `RemoteAddr` is the proxy's address rather
than the caller's, so in the deployed topology this behaves as one global limit shared by
everyone rather than a per-client one.

## What is never exposed

`schemas.Paste` tags both `Content` and `DeleteToken` as `json:"-"`, so serializing the model
directly cannot leak either. `content` reaches a client only through the dedicated
`ContentResponse`, and `delete_token` only in the create response.
