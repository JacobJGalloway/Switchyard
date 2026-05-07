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

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | — | PostgreSQL connection string (required) |
| `LOGISTICS_BASE_URL` | — | Switchyard .NET Logistics API base URL |
| `INVENTORY_BASE_URL` | — | Switchyard .NET Inventory API base URL |
| `AUTH0_DOMAIN` | — | Auth0 tenant domain |
| `AUTH0_CLIENT_ID` | — | M2M client ID (event handler only) |
| `AUTH0_CLIENT_SECRET` | — | M2M client secret |
| `AUTH0_AUDIENCE` | — | Auth0 API audience identifier |
| `SMTP_HOST` | — | Outbound SMTP host for notifications |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | — | SMTP login user |
| `SMTP_PASS` | — | SMTP password |
| `DISPATCH_EMAIL` | — | Recipient address for dispatch alerts |
| `HOS_WARNING_THRESHOLD_HOURS` | `2.0` | Hours remaining before HOS warning fires |
| `DEADHEAD_WINDOW_HOURS` | `4.0` | Minimum hours between run fulfillment and dead-head pairing |
| `DEADHEAD_SEARCH_WINDOW_HOURS` | `2.0` | Look-ahead window when finding eligible dead-head pairings |
| `LOADING_AGE_THRESHOLD_HOURS` | `4.0` | Hours before a BOL stuck in `loading` triggers a dispatch alert |
| `DEFAULT_CYCLE_LABEL` | `60h/7d` | HOS cycle used when none is specified on a plan |
| `WAREHOUSE_IDS` | `WH001,...,WH009` | Comma-separated list of warehouse IDs the regional inventory endpoint fans out to. Add a new warehouse here whenever a warehouse is added to the network. *(v1.2: this flat list will be replaced by a `WAREHOUSE_REGIONS` param once the `region` field is added to the Warehouse model — new warehouses will be picked up automatically via the DB rather than requiring a config change.)* |
| `PORT` | `8080` | HTTP listen port |
| `LOG_LEVEL` | `info` | Log verbosity |
