# Capsule

Ephemeral encrypted sharing. Secrets that self-destruct.

Capsule is a zero-knowledge, end-to-end encrypted paste tool. Share secrets, code snippets, or sensitive text — everything burns after reading.

## How it works

1. You write a secret
2. Capsule encrypts it in your browser (AES-256-GCM)
3. The encrypted blob is stored on the server — the key never leaves your browser
4. You get a link with the decryption key in the URL fragment
5. The recipient opens the link, decrypts in their browser
6. The paste is destroyed

The server never sees your plaintext. Zero knowledge, zero trust.

## Stack

- **API**: Go (Chi, GORM, PostgreSQL)
- **Client**: SvelteKit 5, Tailwind CSS 4
- **Encryption**: AES-256-GCM (Web Crypto API, client-side)

## Quick Start

### Docker Compose

```bash
docker compose up -d
```

Open http://localhost:3000

### Local Development

**API:**
```bash
cd apps/api
go run .
```

**Client:**
```bash
cd apps/client
bun install
bun dev
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://...` | PostgreSQL connection string |
| `PORT` | `4000` | API server port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `MAX_PASTE_SIZE` | `1048576` | Max paste size in bytes (1MB) |
| `ORIGIN` | `http://localhost:3000` | Client origin URL |
| `API_URL` | `http://localhost:4000` | Backend API URL (client-side proxy target) |

## Part of [Facile](https://facile.studio)
