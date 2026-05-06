# Switchyard Go Backend

Go backend for the Switchyard dispatch and logistics platform.

## Prerequisites

- Go 1.25+
- PostgreSQL (connection string via `DATABASE_URL` environment variable)

## First-time setup

```bash
go mod tidy
```

> **Note — viper/gotenv incompatibility:** `github.com/subosito/gotenv` v1.6.0 reorganized
> its package structure and breaks viper's internal dotenv encoder. `go.mod` pins gotenv to
> v1.4.2 as an indirect dependency. If you ever see the error
> `module gotenv@latest found but does not contain package github.com/subosito/gotenv`
> after updating dependencies, run `go get github.com/subosito/gotenv@v1.4.2` before
> re-running `go mod tidy`.

## Build & Run

```bash
go build ./...          # compile all packages
go run ./cmd/main.go    # start the server (default port 8080)
go test ./...           # run all tests
```

## Environment

Copy `.env.example` to `.env` and fill in values before running locally.

| Variable         | Default | Description                        |
|------------------|---------|------------------------------------|
| `PORT`           | `8080`  | HTTP listen port                   |
| `DATABASE_URL`   | —       | PostgreSQL connection string       |
| `AUTH0_DOMAIN`   | —       | Auth0 tenant domain                |
| `AUTH0_CLIENT_ID`| —       | M2M client ID (event handler only) |
| `AUTH0_CLIENT_SECRET` | — | M2M client secret                 |
| `DOTNET_BASE_URL`| —       | Switchyard .NET API base URL       |
