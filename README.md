# goauthly-fiber-example

> **Example / test project** demonstrating how to integrate
> [goAuthly](https://github.com/Keksclan/goAuthly) with
> [Fiber v3](https://github.com/gofiber/fiber) and
> [goConfy](https://github.com/keksclan/goConfy) for configuration.

⚠️ This is **not** a production-ready product.
It uses production-grade structure patterns (clean layout, graceful shutdown,
structured errors) so you can copy the approach into real services, but the
project itself is intentionally minimal.

---

## Overview

This service shows:

- **JWT validation** (issuer, audience, expiry) via goAuthly.
- **Opaque-token introspection** as a fallback when a JWKS-only flow is not
  enough.
- Configuration through **goConfy** – a single YAML file with `{ENV:KEY:default}`
  macro expansion, so environment variables override YAML defaults.
- `AUTH_AUDIENCE=*` – wildcard means _any_ audience is accepted.
- **Required-header enforcement** (`AUTH_REQUIRED_HEADERS`) – the middleware
  rejects requests that lack the listed headers (e.g. `x-user-sub`).
- Structured JSON error responses with stable error codes.
- Request-ID propagation (generates a UUID when the caller does not send one).
- Graceful shutdown on `SIGINT` / `SIGTERM`.

---

## Project Structure

```
goauthly-fiber-example/
├── cmd/server/main.go                 # Entrypoint
├── internal/
│   ├── app/
│   │   ├── app.go                     # App init + goAuthly Engine setup
│   │   └── shutdown.go                # Graceful shutdown
│   ├── config/
│   │   ├── config.go                  # Config loading via goConfy
│   │   └── config.example.yml        # Example YAML (goConfy template)
│   ├── http/
│   │   ├── router.go                  # Route registration
│   │   ├── handlers/
│   │   │   ├── health.go             # /healthz, /readyz
│   │   │   └── me.go                 # /me (protected)
│   │   └── middleware/
│   │       ├── requestid.go          # X-Request-Id
│   │       ├── logger.go            # Structured request logging
│   │       └── auth.go              # goAuthly auth middleware
│   └── platform/errors/
│       └── http_errors.go            # Standard JSON error responses
├── scripts/
│   └── dev.sh                         # Dev helper (loads .env, runs server)
├── .editorconfig
├── .env.example
├── .gitignore
├── config.yml                         # Runtime config (gitignored)
├── CONTRIBUTING.md
├── docker-compose.yml
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
└── SECURITY.md
```

---

## Requirements

| Dependency | Version       |
|------------|---------------|
| **Go**     | 1.26+         |
| Docker     | optional      |

---

## Configuration

### How goConfy works

The application loads `config.yml` (or the path in `CONFIG_PATH`) through
**goConfy**. The YAML file may contain `{ENV:KEY:default}` macros that goConfy
replaces with the corresponding environment variable at load time. This means
you can keep sensible defaults in YAML and override them per-environment with
plain `ENV` vars.

### YAML template (config.yml)

```yaml
server:
  port: "{ENV:SERVER_PORT:8080}"
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 60s

auth:
  issuer: "{ENV:AUTH_ISSUER:}"
  audience: "{ENV:AUTH_AUDIENCE:*}"
  jwks_url: "{ENV:AUTH_JWKS_URL:}"
  introspection_url: "{ENV:AUTH_INTROSPECTION_URL:}"
  client_id: "{ENV:AUTH_CLIENT_ID:}"
  client_secret: "{ENV:AUTH_CLIENT_SECRET:}"
  required_headers: "{ENV:AUTH_REQUIRED_HEADERS:}"
```

### Environment variables

| Variable                 | Description                                                  | Default      |
|--------------------------|--------------------------------------------------------------|--------------|
| `SERVER_PORT`            | HTTP listen port                                             | `8080`       |
| `AUTH_ISSUER`            | Expected JWT `iss` claim                                     | –            |
| `AUTH_AUDIENCE`          | `*` = accept any audience; otherwise a comma-separated list  | `*`          |
| `AUTH_JWKS_URL`          | JWKS endpoint for JWT signature verification                 | –            |
| `AUTH_INTROSPECTION_URL` | OAuth 2 token introspection endpoint (optional)              | –            |
| `AUTH_CLIENT_ID`         | Client ID for introspection Basic Auth (optional)            | –            |
| `AUTH_CLIENT_SECRET`     | Client Secret for introspection Basic Auth (optional)        | –            |
| `AUTH_REQUIRED_HEADERS`  | Comma-separated required headers, e.g. `x-user-sub`         | –            |
| `CONFIG_PATH`            | Path to the YAML config file (optional)                      | `config.yml` |

#### `AUTH_AUDIENCE` behavior

| Value              | Effect                                             |
|--------------------|----------------------------------------------------|
| `*` (or empty)     | **Any** audience is accepted                       |
| `api,admin-portal` | Token `aud` must contain at least one listed value |

#### `AUTH_REQUIRED_HEADERS` behavior

When set (e.g. `x-user-sub`), the auth middleware returns **400 Bad Request**
for every authenticated request that does _not_ carry those headers. Multiple
headers are comma-separated:

```
AUTH_REQUIRED_HEADERS=x-user-sub,x-tenant-id
```

---

## Quickstart (Local)

```bash
# 1. (optional) Create a local .env with your auth-server URLs
cp .env.example .env
# Edit .env – set AUTH_ISSUER, AUTH_JWKS_URL, etc.

# 2. Fetch dependencies
go mod tidy

# 3. Run the server
go run ./cmd/server/

# Or use the dev script (auto-loads .env):
bash scripts/dev.sh
```

The server listens on `http://localhost:8080` by default.

---

## Quickstart (Docker)

```bash
# Adjust environment values in docker-compose.yml, then:
docker compose up --build
```

---

## Endpoints

| Method | Path       | Auth       | Description                                   |
|--------|------------|------------|-----------------------------------------------|
| GET    | `/healthz` | public     | Always returns `200 OK`                       |
| GET    | `/readyz`  | public     | `200` when config is loaded and engine is ready |
| GET    | `/me`      | Bearer ✓   | Returns `sub`, `scopes`, `aud`, `iss` from the token |

### Sample curl

```bash
# Health check
curl http://localhost:8080/healthz
# → {"status":"ok"}

# Readiness check
curl http://localhost:8080/readyz
# → {"status":"ready"}

# Protected endpoint – pass a valid Bearer token
curl -H "Authorization: Bearer <YOUR_TOKEN>" \
     -H "x-user-sub: user123" \
     http://localhost:8080/me
# → {"sub":"user123","scopes":["openid"],"aud":"my-client","iss":"https://auth.example.com","type":"jwt","source":"jwt"}

# Without a token → 401
curl http://localhost:8080/me
# → {"error":"unauthorized","code":"missing_authorization_header","message":"missing authorization header","request_id":"..."}
```

---

## How to Plug In Your Auth Server

Replace the placeholder URLs with your real authorization-server values:

```env
AUTH_ISSUER=https://your-auth-server.example.com
AUTH_JWKS_URL=https://your-auth-server.example.com/.well-known/jwks.json
AUTH_INTROSPECTION_URL=https://your-auth-server.example.com/oauth2/introspect
AUTH_CLIENT_ID=my-service
AUTH_CLIENT_SECRET=change-me
```

Token validation happens entirely through **goAuthly** – the example just wires
the engine into a Fiber middleware. See `internal/app/app.go` for the engine
setup and `internal/http/middleware/auth.go` for the middleware.

---

## Error Format

All auth errors follow a consistent JSON structure:

```json
{
  "error": "unauthorized",
  "code": "invalid_token",
  "message": "token invalid or expired",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

| Code                           | HTTP | Description                      |
|--------------------------------|------|----------------------------------|
| `missing_authorization_header` | 401  | No `Authorization` header        |
| `invalid_token`                | 401  | Token invalid or expired         |
| `missing_required_header`      | 400  | A required header is missing     |
| `forbidden`                    | 403  | Access denied                    |

---

## Troubleshooting

| Symptom | Common cause |
|---------|--------------|
| **401 – invalid_token** | Issuer (`iss`) mismatch between token and `AUTH_ISSUER` |
| **401 – invalid_token** | Token expired (`exp` in the past) |
| **401 – invalid_token** | Audience (`aud`) not in `AUTH_AUDIENCE` list |
| **401 – invalid_token** | JWKS endpoint unreachable or returning wrong keys |
| **400 – missing_required_header** | Request lacks a header listed in `AUTH_REQUIRED_HEADERS` |

---

## Makefile Targets

```
make build          # Build binary to ./bin/
make run            # Build + run
make dev            # Run via scripts/dev.sh (loads .env)
make test           # go test ./...
make docker-build   # Docker image build
make docker-up      # docker compose up --build -d
make docker-down    # docker compose down
make clean          # Remove ./bin/
```

---

## License

This project is released under the [MIT License](LICENSE).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
