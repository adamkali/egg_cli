<img width="200" height="200" alt="egg" src="https://github.com/user-attachments/assets/cb9da171-2b36-4a74-8e0b-629c38bab2bc" />

# egg_cli

A fast and simple command line interface for creating fullstack web applications with the **egg** opinionated framework.

In truth, **egg** is a collection of tools and libraries assembled so that developers can get a production-ready Go backend up and running without drowning in boilerplate.

## What's included in every egg project?

- **REST API** — Echo-based routing with Swagger docs generated automatically
- **Database integration** — SQLC for type-safe queries, Goose for migrations (PostgreSQL)
- **Authentication** — JWT (symmetric) or Auth0 OIDC — your choice at generation time
- **Cache integration** — Redis
- **S3 integration** — MinIO-compatible
- **Docker setup** — multi-stage Dockerfile ready to go
- **Organised project structure** following Go best practices

Core libraries used:

| Library | Purpose |
|---|---|
| [Echo](https://echo.labstack.com/) | HTTP router and middleware |
| [SQLC](https://sqlc.dev/) | Type-safe SQL → Go code generation |
| [Goose](https://github.com/pressly/goose) | Database migrations |
| [echo-jwt](https://github.com/labstack/echo-jwt) | JWT middleware |
| [go-oidc](https://github.com/coreos/go-oidc) | Auth0 OIDC token verification |
| [MinIO Go SDK](https://github.com/minio/minio-go) | S3-compatible object storage |
| [go-redis](https://github.com/redis/go-redis) | Redis client |
| [echo-swagger](https://github.com/swaggo/echo-swagger) | Swagger UI |

---

## Installation

### From Go
```bash
go install github.com/adamkali/egg_cli@latest
```

### Build from source
```bash
git clone https://github.com/adamkali/egg_cli.git
cd egg_cli
make build
```

---

## Quick Start

### 1. Generate an example config

```bash
egg_cli generate example          # JWT authentication (default)
egg_cli generate example --auth0  # Auth0 OIDC authentication
```

This writes an `egg.yaml` in the current directory with sensible defaults. Edit it to fill in your real database URL, secrets, and (for Auth0) your tenant details.

### 2. Generate the project

```bash
egg_cli generate --config egg.yaml
```

egg_cli will:
1. Clone the matching template repo from GitHub (`egg-template-jwt` or `egg-template-auth0`)
2. Substitute all `__EGG_*__` placeholders with your config values
3. Run `go mod download`, `sqlc generate`, `goose up`, `swag init`, `go mod tidy`
4. Leave you with a fully-compiling, runnable project in `./<name>/`

### 3. Interactive init wizard

If you prefer a guided setup:

```bash
egg_cli init
```

The TUI wizard walks through every config section (project info, server, auth, database, S3) and generates both the `egg.yaml` and the project in one shot.

---

## Authentication options

### JWT (default)

The generated project validates tokens with a symmetric HS256 secret stored in your config. `AuthService` issues and checks tokens; tokens are persisted in a `tokens` DB table for revocation support.

### Auth0

Generated with `generate example --auth0`. The project uses the Auth0 OIDC discovery document to validate bearer tokens — no symmetric secret required. The auth flow:

| Endpoint | Method | Description |
|---|---|---|
| `/api/users/login` | GET | Redirects to Auth0 Universal Login |
| `/api/users/signup` | GET | Redirects to Auth0 with `screen_hint=signup` |
| `/api/users/callback` | GET | Exchanges auth code for tokens, returns user claims |

All other protected routes use the OIDC middleware to verify the `Authorization: Bearer <token>` header directly against Auth0's JWKS endpoint.

After generation, update these fields in `config/dev.yaml`:

```yaml
server:
  auth0:
    domain:        your-tenant.us.auth0.com
    audience:      https://api.example.com
    client_id:     your-auth0-client-id
    client_secret: your-auth0-client-secret
    callback_url:  http://localhost:8080/api/users/callback
```

---

## Commands

### `generate`

```bash
egg_cli generate --config egg.yaml
```

Generates a full project from a config file.

#### `generate example`

```bash
egg_cli generate example              # generates egg.yaml with JWT auth
egg_cli generate example --auth0      # generates egg.yaml with Auth0 auth
```

#### `generate dotenv`

Creates a `.env` file from a config file:

```bash
egg_cli generate dotenv                        # reads config/development.yaml
egg_cli generate dotenv path/to/config.yaml    # specific file
egg_cli generate dotenv --output .env.prod     # custom output path
```

---

### `init`

```bash
egg_cli init
```

Interactive TUI wizard — covers all config sections and generates the project on completion.

---

### `env`

```bash
egg_cli env
egg_cli env --env prod --input .env.prod
```

Builds a config YAML from `EGG_`-prefixed environment variables. Useful in CI or Docker environments where secrets live in env vars rather than config files.

Required variables:

| Variable | Description |
|---|---|
| `EGG_NAME` | Project name |
| `EGG_SEMVER` | Project version |
| `EGG_SERVER_PORT` | HTTP port |
| `EGG_AUTH_PROVIDER` | `jwt` or `auth0` |
| `EGG_SERVER_JWT` | JWT secret (jwt mode only) |
| `EGG_SERVER_FRONTEND_DIR` | Built frontend directory |
| `EGG_DATABASE_USERNAME` | Postgres username |
| `EGG_DATABASE_PASSWORD` | Postgres password |
| `EGG_DATABASE_DBNAME` | Database name |
| `EGG_DATABASE_PORT` | Postgres port |
| `EGG_DATABASE_HOST` | Postgres host |
| `EGG_CACHE_PASSWORD` | Redis password |
| `EGG_CACHE_HOST` | Redis host |
| `EGG_CACHE_PORT` | Redis port |
| `EGG_S3_ACCESS` | S3 / MinIO access key |
| `EGG_S3_SECRET` | S3 / MinIO secret key |
| `EGG_S3_URL` | S3 / MinIO endpoint |

---

### `version`

```bash
egg_cli version              # basic version string
egg_cli version --verbose    # full build details
egg_cli version --oneline    # single line
egg_cli version --compact    # version+hash+buildtime
```

---

## Generated project structure

```
my-project/
├── cmd/
│   ├── configuration/configuration.go
│   ├── db.go           # sqlc generate wrapper
│   ├── migrate.go      # goose up
│   ├── down.go         # goose down
│   ├── swag.go         # swag init wrapper
│   └── version.go
├── config/
│   └── dev.yaml
├── controllers/
│   ├── controller.go
│   ├── routes.go
│   └── user_controller.go
├── db/
│   ├── migrations/
│   ├── queries/
│   └── repository/     # generated by sqlc
├── docs/               # generated by swag
├── middlewares/configs/
│   ├── AuthConfig.go
│   └── StaticConfig.go
├── models/
│   ├── handlers/user_handlers/
│   ├── requests/
│   └── responses/
├── services/
│   ├── AuthService.go
│   ├── IAuthService.go
│   ├── UserService.go
│   ├── MockAuthService.go
│   └── ...
├── Dockerfile
├── Makefile
├── main.go
├── sqlc.yml
└── web/dist/           # drop your built frontend here
```

---

## Template repos

egg_cli generates projects by cloning one of these templates and substituting placeholders:

- **[egg-template-jwt](https://github.com/adamkali/egg-template-jwt)** — symmetric JWT auth
- **[egg-template-auth0](https://github.com/adamkali/egg-template-auth0)** — Auth0 OIDC auth

---

## Design philosophy

### Services and interfaces

Every service (`AuthService`, `UserService`, `MinioService`, …) implements a matching `IService` interface. Swap in a mock for tests, swap in a different implementation for different providers — handlers never care.

### Handlers

Handlers in `models/handlers/` follow a single pattern: `Handle()` does the work and returns itself; `JSON()` renders the response. This keeps logic and serialisation cleanly separated and makes A/B testing between two service implementations trivial.

### Echo ecosystem

egg rides on Echo's plugin ecosystem for JWT middleware, Swagger UI, and static file serving — no need to reinvent those wheels.

---

## Development

```bash
make build      # build binary with version info
make test       # go vet + staticcheck + go test
make tidy       # go mod tidy + verify
make clean      # remove build artifacts
make all        # tidy → test → build
```

End-to-end generation test (requires Docker):

```bash
bash scripts/test-generate.sh
```

Spins up a Postgres container, generates both a JWT and an Auth0 project, verifies no unreplaced placeholders remain, and confirms `go build ./...` passes for each.

---

## Contributing

1. Fork the repository
2. Work on the `dev` branch
3. Run `make test` and `bash scripts/test-generate.sh` before submitting
4. Open a pull request against `dev`

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.

## Support

- 🐛 [Issues](https://github.com/adamkali/egg_cli/issues)
- 💬 [Discussions](https://github.com/adamkali/egg_cli/discussions)
