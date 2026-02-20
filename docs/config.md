# Configuration Reference

> See also: the [README – Configuration section](../README.md#configuration).

## How Configuration Is Loaded

1. The application calls `config.Load(configPath)` at startup.
2. **goConfy** reads the YAML file and expands `{ENV:KEY:default}` macros using
   the current environment variables.
3. After decoding, `Config.Normalize()` is called automatically — this parses
   `AUTH_REQUIRED_HEADERS` from a comma-separated string into a slice.
4. `Config.Validate()` applies defaults for any missing server settings (port,
   timeouts).

### Config file resolution order

| Priority | Source                       |
|----------|------------------------------|
| 1        | `configPath` argument        |
| 2        | `CONFIG_PATH` env variable   |
| 3        | `config.yml` in working dir  |

## Environment Variables

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

## Audience Rules

- **`AUTH_AUDIENCE=*`** (or empty) — any audience is accepted; the `aud` claim
  is not checked.
- **`AUTH_AUDIENCE=api,admin-portal`** — the token's `aud` must contain **at
  least one** of the listed values.

Internally, `AuthConfig.AudienceList()` returns `nil` for the wildcard case and
a `[]string` slice otherwise. `AudienceIsWildcard()` is a convenience helper.

## Required Headers

When `AUTH_REQUIRED_HEADERS` is set, the auth middleware checks each
authenticated request for the listed HTTP headers. If any are missing, the
request is rejected with **400 Bad Request** and error code
`missing_required_header`.

Example — require both `x-user-sub` and `x-tenant-id`:

```env
AUTH_REQUIRED_HEADERS=x-user-sub,x-tenant-id
```

## OAuth 2 Mode Selection

The goAuthly engine mode is determined automatically based on which URLs are
configured:

| `AUTH_JWKS_URL` | `AUTH_INTROSPECTION_URL` | Mode               |
|-----------------|--------------------------|--------------------|
| ✓               | ✗                        | JWT only           |
| ✗               | ✓                        | Opaque only        |
| ✓               | ✓                        | JWT + Opaque       |

See `internal/app/app.go` → `buildAuthEngine()` for the exact logic.
