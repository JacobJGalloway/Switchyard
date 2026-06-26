<div align="center">
<img src="Switchyard.UI/src/assets/logo-full-name-light.png" />
</div>

# Switchyard 1.2

Switchyard is an inventory, driver, and equipment tracking and management system which coordinates logistics operations across a network of warehouses and stores. Inventory is tracked per location; Bills of Lading govern movement between any combination of stops — from same-day local transfers to multi-stop OTR runs with partial loads. Authenticated via Auth0.

## Solution Structure

| Project | Role | Port |
|---|---|---|
| `Switchyard.InventoryAPI` | Inventory API — Clothing, PPE, Tools | 7000 |
| `Switchyard.LogisticsAPI` | Logistics API — Bills of Lading, Stores, Warehouses, Users | 7001 |
| `Switchyard.Domain` | Shared class library — domain models for both .NET APIs | — |
| `Switchyard.UI` | React/TypeScript client UI | 5173 |
| `Switchyard-Go` | Go backend — PlanBOL, Dispatch Whiteboard, HOS, Equipment | 8080 |

**Go backend documentation:** [`Switchyard-Go/README.md`](Switchyard-Go/README.md) — setup, environment variables, key architectural constraints, and API reference.

## Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| [.NET SDK](https://dotnet.microsoft.com/download) | 10.0 | InventoryAPI and LogisticsAPI |
| [Go](https://go.dev/dl/) | 1.25+ | Switchyard-Go backend |
| [Node.js](https://nodejs.org/) | 24+ | Switchyard.UI |
| [PostgreSQL](https://www.postgresql.org/download/) | 16 | All backends — Docker or local install |
| [Auth0 account](https://auth0.com/) | — | Tenant + API resource + two M2M applications |

**PostgreSQL:** Used by both the Go backend and the .NET APIs. Easiest to run via Docker (`postgres:16` image). Default dev port is **5433** — if another Postgres instance is already on 5432, use 5433 to avoid the conflict. See `Switchyard-Go/.env.example` for the full connection string format. The .NET APIs connect to separate databases (`switchyard_inventory`, `switchyard_logistics`) on the same instance.

**Go Service Initialization:** Due to how the environmental variables are read in Go, the initial setup for Docker will need to be different if the image is not up and running on a container. Subsequent restarts with the container already running are a single line restart. See the README.md under Switchyard-Go for more details.

**Auth0 M2M applications (free tier: 2 slots):** Switchyard uses both — one for the Scalar UI on the .NET APIs, one for the Go event handler. Confirm available M2M slots before setting up a new tenant. See [Auth0 Setup](#auth0-setup) for full configuration steps.

**SMTP:** Required for email notifications (HOS warnings, breakdown alerts, dead-head expiry timer). Any SMTP-accessible mail account works in dev. Notifications can be left unconfigured early on — fill in `SMTP_*` env vars when you need them.

**Go `.env` loading — known gotcha:** Viper does not auto-load `.env` files. Before running the Go backend, source your `.env` from `Switchyard-Go/` using the one-liner in `Switchyard-Go/README.md`, or set the vars directly in your shell session.

## Architecture

### Agentic AI Development

Switchyard was built with [Claude Code](https://claude.ai/code) as an active pair programming and mentoring partner throughout every sprint. The model always knew the answer — but frequently chose not to give it directly. Instead it would surface the right question, flag the constraint that hadn't been named yet, or ask what the invariant was before suggesting the fix. The design emerged from understanding, not from prescription.

That dynamic shaped the architecture as much as the implementation. The CUD authority boundary between Go and .NET, the service-layer transport agnosticism, and the CQRS read-replica pattern were all reasoned out through guided back-and-forth rather than handed down as decisions. The goal was to own every piece of it — not just ship it.

### CQRS Read Replica
Both .NET APIs maintain a read replica synced asynchronously after every write:
- Write operations target the primary PostgreSQL database
- Read operations target the read replica database (all `AsNoTracking`)
- `SaveChangesInterceptor` → `Channel<SyncJob>` → `BackgroundService` (full table resync per changed entity type)
- `GET /api/Audit` on each API reports write vs read row counts with an `InSync` flag

### Data Layer Pattern
- **Unit of Work** over repositories — services depend on `IUnitOfWork`
- **Repositories** — separate write context (CUD) and read context (queries)
- **EF Core** with PostgreSQL (Npgsql); migrations applied on startup for write context; `EnsureCreated` for read replica
- **Switchyard.Domain** — shared class library containing all entity models and interfaces; neither API project owns domain models directly

### Auth
Both APIs use Auth0 JWT bearer authentication. Permissions are claim-based:

| Permission | Used by |
|---|---|
| `read:inventory` | Inventory read endpoints |
| `read:bol` | Logistics read endpoints |
| `create:bol` | BOL creation |
| `modify:bol` | ProcessStop, ReplaceStop |
| `manage:users` | User management |

## API Endpoints

### Inventory API (`/api`) — port 7000

| Method | Path | Description |
|---|---|---|
| GET | `/Clothing` | All clothing items |
| GET | `/Clothing/{skuId}` | By SKU |
| GET | `/Clothing/location/{locationId}` | By location |
| GET | `/Clothing/filter?locationId=&skuId=` | By location + SKU |
| POST | `/Clothing` | Add item |
| PUT | `/Clothing/{skuId}` | Full update by SKU |
| PATCH | `/Clothing/item/{partitionKey}` | Patch projected/unloadedDate |
| DELETE | `/Clothing/item/{partitionKey}` | Delete item |
| _(same shape for `/PPE` and `/Tool`)_ | | |
| GET | `/Audit` | Write vs read row counts |

### Logistics API (`/api`) — port 7001

| Method | Path | Description |
|---|---|---|
| GET | `/BillOfLading` | All BOLs |
| GET | `/BillOfLading/{transactionId}` | BOL + line entries |
| GET | `/BillOfLading/{transactionId}/line-entry` | Line entries only |
| POST | `/BillOfLading` | Create BOL, persist line entries, write `.txt` to Downloads |
| POST | `/BillOfLading/{transactionId}/process/{locationId}` | Mark location stop as processed |
| POST | `/BillOfLading/{transactionId}/replace-stop` | Move unprocessed stop to a new location |
| GET | `/Store` | All stores |
| GET | `/Warehouse` | All warehouses |
| GET | `/User` | All Auth0 users |
| POST | `/User` | Create Auth0 user + assign role |
| PATCH | `/User/{userId}/deactivate` | Block user (soft deactivate) |
| GET | `/Audit` | Write vs read row counts |

## Auth0 Setup

1. Create an API resource in Auth0 and set its identifier as `Auth0:Audience`
2. Set `Auth0:Authority` to your Auth0 domain (e.g. `https://your-tenant.auth0.com/`)
3. Add permissions to the API: `read:inventory`, `read:bol`, `create:bol`, `modify:bol`, `manage:users`
4. For user management, create an M2M application and grant it the Auth0 Management API with scopes:
   `read:users`, `create:users`, `update:users`, `read:roles`, `create:role_members`
5. Set credentials in `{API Project Name}/appsettings.Development.json` (gitignored):

```json
{
  "Auth0": {
    "Authority": "https://your-tenant.auth0.com/",
    "Audience": "your-api-audience",
    "ScalarClientId": "your-m2m-client-id",
    "ScalarClientSecret": "your-m2m-client-secret"
  },
  "ConnectionStrings": {
    "InventoryWrite": "Host=localhost;Port=5433;Database=switchyard_inventory;Username=postgres;Password=password",
    "InventoryRead": "Host=localhost;Port=5433;Database=switchyard_inventory_read;Username=postgres;Password=password"
  }
}
```

## Running the System

```bash
# APIs (run each in a separate terminal)
dotnet run --project Switchyard.InventoryAPI
dotnet run --project Switchyard.LogisticsAPI

# Go support services — start Postgres first (docker-compose.yml at project root)
docker compose up -d
cd Switchyard-Go
Get-Content .env | Where-Object { $_ -notmatch '^\s*#' -and $_ -match '=' } | ForEach-Object { $k,$v = $_ -split '=',2; Set-Item "Env:$($k.Trim())" $v.Trim() }
go run ./cmd/main.go

# UI
cd Switchyard.UI
npm run dev

# Unit Tests
dotnet test
go test ./...

# API docs (Scalar UI — while API is running)
# Inventory: https://localhost:7000/scalar/v1
# Logistics: https://localhost:7001/scalar/v1
```

## Wanted Features

### v1.3 — Next sprint
- [ ] Mid-BOL transfer stops — `transfer` stop type; formal custody checkpoint for driver/equipment handoffs mid-route; requires `DriverBOLAssignment` restructuring
- [X] Demo reset / reseed script — date-relative seed so the board always looks like a live operational day at demo time
- [X] Two-company demo seed — Company A (Monday morning, default brand) and Company B (mid-week complexity, client palette override)
- [X] Dispatch board dark mode nuance rework + favicon swap
- [X] SKU unit price — extend inventory model to hold unit price; foundation for analytics such as revenue vs. profit charts and other displays
- [X] SKU movement chart — warehouse filter, SKU filter, and chart dots all working; EF Core funcletizer bug bypassed by materializing before `Split()`
- [X] Explicit routes on all controllers — plural kebab-case across all 10 controllers in both APIs; UI callers and Go integration clients updated to match
- [X] Auth0 Google login — single-click through; `onRedirectCallback` + `isLoading` guard + `strictPort: true` in `vite.config.ts` to prevent port drift
- [X] Receive-delivery endpoint — flip `Projected = false` and set `LocationId` when a BOL stop is confirmed; currently missing, seed mimics it manually

### v1.4 Wanted Features - Demo Stable Hardening / pilot-client ready; order flexible relative to v1.3 completion
- [ ] Rolling refresh tokens for Auth0 sessions in place of fixed-expiry client secrets
- [ ] Color contrast audit (WCAG AA) — verify all text/bg combinations across light and dark themes
- [ ] ARIA compliance audit — board columns, cards, icon-only buttons, skip-nav
- [ ] Document dispatch board card border language in a design note — border present = status alert (danger/warn/ok variants); no border = clean/good status. Borders carry semantic weight and should not be used decoratively.
- [ ] CQRS read replica hardening — provision real Postgres read replicas for `switchyard_inventory` and `switchyard_logistics`; dev currently points read connections at the write database as a stand-in.

### Backlog
- [ ] Read replica health endpoint — expose sync lag and InSync status
- [ ] Extract User Management to a dedicated identity service when the data layer splits
- [ ] Scalar branding — Switchyard logo and name above the API title; currently blocked by Scalar's limited logo support in the .NET package
