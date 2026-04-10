# Warehouse Management System

A simplified retail logistics system demonstrating warehouse-to-store inventory transfers. Three warehouses each supply three stores within same-day transfer range. Inventory is tracked per location; Bills of Lading move it between them. Authenticated via Auth0.

## Solution Structure

| Project | Role | Port |
|---|---|---|
| `WarehouseInventoryAPI-Claude` | Inventory API — Clothing, PPE, Tools | 7000 |
| `WarehouseLogisticsAPI-Claude` | Logistics API — Bills of Lading, Stores, Warehouses, Users | 7001 |
| `WarehouseSalesUI-Claude` | React/TypeScript client UI | 5173 |

**Shared database:** `Sqlite 3 Implementation/WarehouseData.db3`
**Read replica:** `Sqlite 3 Implementation/WarehouseRead.db3` (auto-created on startup if not already persised)

## Running the System

```bash
# APIs (run each in a separate terminal)
dotnet run --project WarehouseInventoryAPI-Claude
dotnet run --project WarehouseLogisticsAPI-Claude

# UI
cd WarehouseSalesUI-Claude && npm run dev

# Tests
dotnet test

# API docs (Scalar UI — while API is running)
# Inventory: https://localhost:7000/scalar/v1
# Logistics: https://localhost:7001/scalar/v1
```

## Architecture

### CQRS Read Replica
Both APIs maintain a read replica (`WarehouseRead{WarehouseInventory/WarehouseLogistics}.db3`) synced asynchronously after every write:
- Write operations target `WarehouseData.db3`
- Read operations target appropriate `WarehouseRead.db3` (all `AsNoTracking`)
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

- [ ] Analytics endpoints — inventory over time by location (leverage `UnloadedDate` + projected flag)
- [ ] Shipping figures per warehouse — aggregate from processed BOL line entries
- [ ] Inventory withdrawal UI — `WarehouseSalesUI-Claude` currently shows Inventory Viewer; withdrawal flow not yet built
- [ ] BOL status history — audit trail of status transitions
- [ ] Read replica health endpoint — expose sync lag / InSync status
- [ ] Migrate from `EnsureCreated` to EF Core migrations for controlled schema evolution
- [ ] Extract User Management to a dedicated identity service when the data layer splits
- [ ] Sales UI for system-registered non-employee users (no assigned location, no role) — separate main menu surfacing inventory by type, plus a checkout workflow that reduces inventory quantities via line entries on a customer-facing BOL variant
- [ ] Scalar branding — company logo and name above the API title; currently blocked by Scalar's limited logo support in the .NET package
