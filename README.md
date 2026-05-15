# Switchyard 1.1

Switchyard is an inventory, driver, and equipment tracking and management system which coordinates logistics operations across a network of warehouses and stores. Inventory is tracked per location; Bills of Lading govern movement between any combination of stops — from same-day local transfers to multi-stop OTR runs with partial loads. Authenticated via Auth0.

## Solution Structure

| Project | Role | Port |
|---|---|---|
| `Switchyard.InventoryAPI` | Inventory API — Clothing, PPE, Tools | 7000 |
| `Switchyard.LogisticsAPI` | Logistics API — Bills of Lading, Stores, Warehouses, Users | 7001 |
| `Switchyard.UI` | React/TypeScript client UI | 5173 |
| `Switchyard-Go` | Go backend — PlanBOL, Dispatch Whiteboard, HOS, Equipment | 8080 |

**Shared database:** `Sqlite 3 Implementation/WarehouseData.db3`
**Read replica:** `Sqlite 3 Implementation/WarehouseRead.db3` (auto-created on startup if not already persisted)

**Go backend documentation:** [`Switchyard-Go/README.md`](Switchyard-Go/README.md) — setup, environment variables, key architectural constraints, and API reference.

## Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| [.NET SDK](https://dotnet.microsoft.com/download) | 10.0 | InventoryAPI and LogisticsAPI |
| [Go](https://go.dev/dl/) | 1.25+ | Switchyard-Go backend |
| [Node.js](https://nodejs.org/) | 24+ | Switchyard.UI |
| [PostgreSQL](https://www.postgresql.org/download/) | 16 | Switchyard-Go backend — Docker or local install |
| [Auth0 account](https://auth0.com/) | — | Tenant + API resource + two M2M applications |

**SQLite** (used by both .NET APIs) is bundled with EF Core — no separate install required. The shared DB file lives at `Sqlite 3 Implementation/WarehouseData.db3` and is created on first startup.

**PostgreSQL:** Easiest to run via Docker (`postgres:16` image). Default dev port is **5433** — if another Postgres instance is already on 5432 (e.g. a Yearly Yields container), use 5433 to avoid the conflict. See `Switchyard-Go/.env.example` for the full connection string format.

**Auth0 M2M applications (free tier: 2 slots):** Switchyard uses both — one for the Scalar UI on the .NET APIs, one for the Go event handler. Confirm available M2M slots before setting up a new tenant. See [Auth0 Setup](#auth0-setup) for full configuration steps.

**SMTP:** Required for email notifications (HOS warnings, breakdown alerts, dead-head expiry timer). Any SMTP-accessible mail account works in dev. Notifications can be left unconfigured early on — fill in `SMTP_*` env vars when you need them.

**Go `.env` loading — known gotcha:** Viper does not auto-load `.env` files. Before running the Go backend, source your `.env` from `Switchyard-Go/` using the one-liner in `Switchyard-Go/README.md`, or set the vars directly in your shell session.

## Running the System

```bash
# APIs (run each in a separate terminal)
dotnet run --project Switchyard.InventoryAPI
dotnet run --project Switchyard.LogisticsAPI

# UI
cd Switchyard.UI
npm run dev

# Tests
dotnet test

# API docs (Scalar UI — while API is running)
# Inventory: https://localhost:7000/scalar/v1
# Logistics: https://localhost:7001/scalar/v1
```

## Architecture

### CQRS Read Replica
Both APIs maintain a read replica (`WarehouseInventoryRead.db3` / `WarehouseLogisticsRead.db3`) synced asynchronously after every write:
- Write operations target `WarehouseData.db3`
- Read operations target the API's own `WarehouseRead.db3` (all `AsNoTracking`)
- `SaveChangesInterceptor` → `Channel<SyncJob>` → `BackgroundService` (full table resync per changed entity type)
- `GET /api/Audit` on each API reports write vs read row counts with an `InSync` flag

### Data Layer Pattern
- **Unit of Work** over repositories — services depend on `IUnitOfWork`
- **Repositories** — separate write context (CUD) and read context (queries)
- **EF Core** with SQLite; `EnsureCreated` on startup for both DBs; initial full sync enqueued at startup

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
  }
}
```

## Wanted Features

### v1.2 — June
- [ ] Operating cost tracking — base rate per mile; roadside tow rates applied to resolved breakdown records; required foundation for revenue vs. profit analytics
- [ ] Advanced analytics and reporting — revenue vs. profit per BOL, driver, and warehouse; cost overlay on throughput charts; depends on operating cost tracking
- [ ] Analytics handler refactor — extract thin `AnalyticsQuerier` interface to enable unit testing; add testcontainers-go integration test suite against real PostgreSQL
- [ ] Returns — `return_depot` stop type on PlanBOLStop; constraint solver already accommodates the extension
- [ ] Mid-BOL transfer stops — `transfer` stop type; formal custody checkpoint for driver/equipment handoffs mid-route; requires `DriverBOLAssignment` restructuring
- [ ] Warehouse region attribute — `region` column on Warehouse model; replaces flat `WAREHOUSE_IDS` env list; new warehouses in a region picked up automatically without a config change
- [ ] White-label theming — "Industrial Cool" light and dark defaults; client DNS-scoped SCSS variable overrides
- [ ] Switchyard brand assets — logo, name, combined lockup, and "Powered by Switchyard" treatment (Light and Dark variants)
- [ ] Extract `Data/` folders to a shared class library — domain models separated from API projects
- [ ] Migrate .NET APIs from SQLite to PostgreSQL — consolidate onto the PostgreSQL stack already running for the Go backend

### Backlog
- [ ] Rolling refresh tokens for Auth0 sessions in place of fixed-expiry client secrets
- [ ] Read replica health endpoint — expose sync lag and InSync status
- [ ] Migrate from `EnsureCreated` to EF Core migrations for controlled schema evolution
- [ ] Extract User Management to a dedicated identity service when the data layer splits
- [ ] Scalar branding — Switchyard logo and name above the API title; currently blocked by Scalar's limited logo support in the .NET package
