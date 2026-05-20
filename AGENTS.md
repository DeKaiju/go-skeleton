# AGENTS.md - go-skeleton Development Guide

## Project Overview

Go service skeleton (module: `github.com/dekaiju/go-skeleton`, Go `1.26.0`) with a reusable backend project structure and no application-specific business coupling.

Current entry points:

- **API Server** (`go run main.go server`): Gin HTTP service with graceful shutdown, environment loading, DB/Redis initialization, and sample auth routes
- **Demo Scripts** (`go run main.go cmd demo hello` / `world`): Example Cobra script commands that show how to bootstrap shared dependencies for non-HTTP jobs

This repository is a template, not an application-specific backend. Keep shared infrastructure patterns, but do not introduce Rise-specific chains, vaults, fiat flows, or other business-coupled concepts into the skeleton.

## Development Commands

### Running the application

```bash
# prepare env
cp .env.example .env

# start API server
go run main.go server
go run main.go server -c .env -p 3000 -m dev
go run main.go server -c /path/to/.env -p 8080 -m prod

# run example scripts
go run main.go cmd demo hello
go run main.go cmd demo world
```

### Testing

```bash
# run all tests
go test ./...

# run a specific package
go test ./test/curl
go test ./test/redis

# compile a test package without running it
go test -c ./test/curl
go test -c ./test/redis
```

### Dependencies

```bash
go mod download
go mod tidy
go mod verify
```

## Architecture

### Project structure

The codebase follows a service-oriented structure with template-safe building blocks:

- `cmd/` - CLI entry points and bootstrapping
- `config/` - environment loading and typed configuration helpers
- `data/` - models and DAO helpers
- `pkg/` - reusable infrastructure packages such as JWT, Redis, MySQL, logging, HTTP client helpers, and response formatting
- `service/server/` - HTTP server code split into `handler`, `logic`, `middleware`, and `route`
- `types/` - shared constants and cross-package errors
- `test/` - example tests

### HTTP layering

Keep the server layer split consistently:

- `service/server/handler/` - bind request data, call logic, map errors to HTTP status codes
- `service/server/logic/` - business logic with no HTTP response formatting
- `service/server/middleware/` - auth, trace ID, request logging, panic recovery, CORS
- `service/server/route/` - route registration only

The current sample endpoints are:

- `GET /healthz`
- `GET /api/auth/nonce`
- `POST /api/auth/login`
- `GET /api/user/profile`

### Middleware

Global middleware currently applied in `service/server/route/route.go`:

- `gin.Logger()`
- `middleware.Cors()`
- `middleware.RecoverAtLast()`
- `middleware.RequestLogger()`
- `middleware.TraceId()`

Route-level middleware:

- `middleware.AuthRequired()` for authenticated endpoints

### Response format

Use `pkg/response`:

```go
response.Success(c, data)       // -> 200 {"data": ..., "message": ""}
response.Fail(c, httpCode, err) // -> httpCode {"data": "", "message": "..."}
```

### Data and infrastructure

- MySQL is initialized in `pkg/mysql/`
- Redis pool is initialized in `pkg/redis/`
- JWT creation and parsing live in `pkg/auth/`
- Request trace IDs are propagated into SQL logging through `pkg/mysql.GinContextToContext`
- Sample session state is stored in Redis
- Sample nonce state is stored in in-memory cache via `pkg/cache/`

### Bootstrapping

Server setup lives in `cmd/api/server.go`:

- resolves app root
- loads `.env`
- configures Gin mode
- configures logging
- validates required environment variables
- initializes MySQL and Redis
- auto-migrates the sample `data.User` model
- starts the HTTP server with graceful shutdown

Scripts that need shared infrastructure should use `cmd/scripts/common.Depend()` and `Release()`.

## Configuration Management

Environment helpers live in `config/app.go`. Prefer using:

- `config.GetEnv`
- `config.GetIntEnv`
- `config.GetBoolEnv`
- `config.MissingEnvs`
- `config.ResolveConfigPath`

Do not scatter repetitive `os.Getenv` parsing logic across handlers and service code when a small reusable config helper is appropriate.

Current required environment variables enforced at server startup:

- `JWT_SECRET`
- `DB_HOST`
- `DB_PORT`
- `DB_DATABASE`
- `DB_USERNAME`

## Code Style

### Imports

Use three import groups separated by blank lines: stdlib, external, internal (`github.com/dekaiju/go-skeleton/...`).

### Readability and locality

- Prefer code that can be understood in one pass during review.
- Do not extract one-off helpers for trivial logic, simple conditionals, or simple map access.
- Keep simple transformations and business checks inline near the call site.
- Prefer locality over abstraction for small logic.
- Keep the main flow explicit instead of hiding it behind thin wrappers.
- If a non-trivial package-private helper is extracted for a single main function, place it immediately below that main function unless it is broadly shared.
- When passing a multi-field struct literal into another call, prefer assigning it to a clearly named local variable first.
- Promote reusable header names, context keys, prefixes, and similar shared values into `types/constant.go`.

### Comments

- Keep or add comments at tricky control-flow points and fallback paths.
- Do not add comments that only restate the code.

### Error handling

- Add useful context when returning errors across package boundaries.
- Return errors early.
- In handlers, distinguish transport errors from business errors by mapping them to the correct HTTP status code.
- Do not collapse infrastructure failures into `401` or other user-facing auth errors unless that is actually the intended behavior.
- Use `log` for runtime logging and operational output whenever possible. Do not use `fmt` for logging when `log` is sufficient. Reserve `fmt` for non-logging string formatting and value construction.

### Context handling

- Logic and DAO functions should accept `context.Context` where possible.
- When code starts from a `gin.Context`, use `pkg/mysql.GinContextToContext(c)` to preserve request cancellation and attach `traceId`.

## Testing

- Standard `testing` is the default testing framework in this repository.
- `test/curl` uses `httptest` for self-contained HTTP client verification.
- `test/redis` uses `miniredis` instead of relying on an external Redis instance.
- Prefer hermetic tests over tests that depend on a manually started local service.

## Common Patterns

### Adding a new API endpoint

1. Create request and response DTOs in `service/server/logic/<domain>/`.
2. Implement the domain logic in `service/server/logic/<domain>/service.go`.
3. Add a handler in `service/server/handler/`.
4. Register the route in `service/server/route/route.go`.
5. Apply middleware explicitly at the route group or endpoint level.

### Adding a new DAO

1. Add or update the model in `data/model.go` or a new model file under `data/`.
2. Add context-aware query helpers in `data/`.
3. Keep DB access out of handlers.

### Adding a new script command

1. Register the command under `cmd/scripts/`.
2. If it needs DB or Redis, use `common.Depend()` and `common.Release()`.
3. Keep script bootstrapping consistent with server bootstrapping.

## Important Constraints

- The `service/server/logic/` layer must stay independent of HTTP response formatting.
- The skeleton should remain generic. Do not hardcode application-specific domains, chain IDs, vendor APIs, or company-specific operational rules into the template.
- Keep changes template-friendly: every addition should improve reuse, not just solve one application’s local needs.
- Auto-migration in `cmd/api/server.go` exists only for the sample model and template ergonomics. If a future project needs formal migrations, add them deliberately instead of quietly expanding `AutoMigrate`.
- Redis-backed session validation and in-memory nonce caching are sample auth mechanisms; keep extensions coherent with that flow unless the template’s auth architecture is intentionally changed.
- For non-trivial changes, define a concrete verification step before or during implementation, preferably a test or an exact command.
