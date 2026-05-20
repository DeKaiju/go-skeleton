# go-skeleton

`go-skeleton` is a reusable Go backend template with a clean service layout and minimal sample auth flow.

## Features

- `cmd/` plus `service/server/` layering for CLI and HTTP services
- Gin, Cobra, GORM, and optional Redis integration
- Consistent logging, `traceId`, response formatting, and graceful shutdown
- Sample Web3 auth flow: `nonce -> SIWE login -> JWT auth`
- Self-contained HTTP and Redis test examples

## Project Layout

```text
.
├── cmd/                    # CLI entry points
├── config/                 # environment and config helpers
├── data/                   # models and DAO helpers
├── pkg/                    # shared infrastructure packages
├── service/server/         # HTTP service layers
│   ├── handler/
│   ├── logic/
│   ├── middleware/
│   └── route/
├── test/                   # sample tests
└── types/                  # shared constants and errors
```

## Quick Start

1. Copy the environment template:

```bash
cp .env.example .env
```

2. Start dependencies:

- MySQL
- Redis is optional and only needed for Redis-backed features such as the sample login session flow

3. Start the API server:

```bash
go run main.go server
```

The default listen address is `:3000`.

## Common Commands

```bash
go run main.go --help
go run main.go server -c .env -p 3000 -m dev
go run main.go cmd demo hello
go test ./...
```

## Template Rules

- `service/server/logic/` should not handle HTTP formatting details.
- `handler/` should only bind input, call logic, and format responses.
- `pkg/mysql` should propagate `traceId` into SQL logging.
- Redis should stay optional at startup. Only features that explicitly need Redis should depend on it at runtime.
- The template auto-migrates the sample `users` table to make first-run verification easier.
