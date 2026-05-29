# Switchyard — Architecture v1.2

Switchyard is an inventory, driver, and equipment tracking and management system which coordinates logistics operations across a network of warehouses and stores. Inventory is tracked per location; Bills of Lading govern movement between any combination of stops — from same-day local transfers to multi-stop OTR runs with partial loads. Authenticated via Auth0.

---

## Section 1 — Solution Structure

| Project | Role | Port |
|---|---|---|
| `Switchyard.InventoryAPI` | Inventory API — Clothing, PPE, Tools | 7000 |
| `Switchyard.LogisticsAPI` | Logistics API — Bills of Lading, Stores, Warehouses, Users | 7001 |
| `Switchyard.Domain` | Shared class library — domain models for both .NET APIs | — |
| `Switchyard.UI` | React/TypeScript client UI | 5173 |
| `Switchyard-Go` | Go backend — PlanBOL, Dispatch Whiteboard, HOS, Equipment | 8080 |

> **1.2 change:** `Switchyard.Domain` is a new shared class library extracted from `Switchyard.InventoryAPI` and `Switchyard.LogisticsAPI`. All domain models now live here. Both APIs reference `Switchyard.Domain` — neither owns domain models directly.

---

## Section 2 — Data Layer

### .NET APIs
- **Unit of Work** over repositories — services depend on `IUnitOfWork`
- **Repositories** — separate write context (CUD) and read context (queries)
- **EF Core** with SQLite; `EnsureCreated` on startup for both DBs
- **Shared database:** `Sqlite 3 Implementation/WarehouseData.db3`
- **Read replica:** `Sqlite 3 Implementation/WarehouseRead.db3` (auto-created on startup if not already persisted)
- Domain models sourced from `Switchyard.Domain` — not defined locally in API projects

### Go Backend
- **PostgreSQL 16** — default dev port `5433`
- **SQLX** for query execution
- **golang-migrate** for schema migrations
- All PlanBOL, Dispatch, HOS, Equipment, and Analytics data lives here

---

## Section 3 — CQRS Read Replica (.NET)

Both .NET APIs maintain a read replica synced asynchronously after every write:

- Write operations target `WarehouseData.db3`
- Read operations target the API's own `WarehouseRead.db3` (all `AsNoTracking`)
- `SaveChangesInterceptor` → `Channel<SyncJob>` → `BackgroundService` (full table resync per changed entity type)
- `GET /api/Audit` on each API reports write vs read row counts with an `InSync` flag

---

## Section 4 — Authentication

Both .NET APIs use Auth0 JWT bearer authentication. Permissions are claim-based:

| Permission | Used by |
|---|---|
| `read:inventory` | Inventory read endpoints |
| `read:bol` | Logistics read endpoints |
| `create:bol` | BOL creation |
| `modify:bol` | ProcessStop, ReplaceStop |
| `manage:users` | User management |

The Go backend uses a dedicated Auth0 M2M application for its event handler. Auth0 free tier provides 2 M2M slots — both are consumed by Switchyard (Scalar UI + Go event handler).

---

## Section 5 — Shared Domain Library (`Switchyard.Domain`)

Introduced in 1.2. Extracted from both .NET API projects to eliminate model duplication and enable clean unit testing across the solution.

**Contains:**
- All EF Core entity models (Inventory, BOL, Store, Warehouse, User)
- Shared enumerations and value objects
- No service logic, no controllers, no EF Core DbContext — models only

**Consumed by:**
- `Switchyard.InventoryAPI` — Clothing, PPE, Tool models
- `Switchyard.LogisticsAPI` — BOL, Store, Warehouse, User models

**Migration path:**
- Move `Data/` folders from each API project into `Switchyard.Domain`
- Update `using` directives in both API projects — compiler guides the sweep
- No logic changes required; purely structural

---

## Section 6 — Operating Cost Tracking (Go)

New in 1.2. Required foundation for revenue vs. profit analytics.

**Base rate per mile:**
- Stored on `DriverBOLAssignment` — snapshot of the rate at the time of commitment
- Default assumed rate: `$3.20/mile` (tractor without full trailer)
- Flatbed tow rate: `$3.80/mile`

**Roadside tow rate:**
- Applied to resolved `BreakdownRecord` entries
- Tow cost recorded at breakdown resolution time
- Feeds directly into per-BOL cost overlay in analytics

**Data points captured:**
- Miles driven per BOL
- Base operating cost per BOL
- Tow costs where applicable
- Driver and warehouse association for aggregation

---

## Section 7 — Analytics Handler Refactor (Go)

Prerequisite for advanced analytics. Completed before analytics feature work begins.

**`AnalyticsQuerier` interface:**
- Thin interface extracted from existing analytics handler
- Enables unit testing without a live database connection
- All analytics queries depend on the interface, not the concrete implementation

**testcontainers-go integration suite:**
- Real PostgreSQL container spun up for integration tests
- Covers operating cost queries, revenue vs. profit aggregations, and throughput overlays
- Runs in CI alongside existing unit tests

---

## Section 8 — Advanced Analytics and Reporting (Go)

Depends on operating cost tracking (Section 6) and analytics handler refactor (Section 7).

**Revenue vs. profit per BOL:**
- Revenue: derived from BOL line entries and stop data
- Profit: revenue minus operating cost (base rate × miles + tow costs where applicable)
- Exposed via analytics handler through `AnalyticsQuerier`

**Aggregations supported:**
- Per BOL
- Per driver
- Per warehouse
- Per region (where region attribute is populated — Section 10)

**Cost overlay on throughput charts:**
- Existing throughput charts in the Dispatch Whiteboard gain a cost overlay layer
- Toggle between throughput-only and throughput + cost views

---

## Section 9 — Returns (`return_depot` Stop Type) (Go)

New stop type on `PlanBOLStop`. Constraint solver already accommodates the extension — implementation is completion work, not architecture work.

**`return_depot` stop type:**
- Represents a driver returning to the originating depot after final delivery
- Recorded as a formal stop in the BOL plan
- Subject to HOS tracking in the same way as any other stop
- Unit test coverage added for the return depot path

---

## Section 10 — Warehouse Region Attribute (Go)

New `region` column on the Warehouse model.

**Replaces:** flat `WAREHOUSE_IDS` environment variable list

**Behavior:**
- New warehouses added to a region are picked up automatically without a configuration change
- Region becomes a grouping dimension for analytics aggregations (Section 8)
- Cross-region routing deferred — build the attribute and intra-region awareness only; cross-region behavior waits for a real use case signal from users

**Migration:** golang-migrate script adds `region` column with a nullable default; existing warehouses backfilled manually or via seed script.

---

## Section 11 — White-Label Theming (UI)

**"Industrial Cool" defaults:**
- Light and dark theme variants
- SCSS variable definitions for both themes
- Applied globally via root-level CSS custom properties

**Client DNS-scoped overrides:**
- SCSS variable overrides resolved at the client domain level
- Allows per-client branding without forking the UI codebase
- Theme resolution order: DNS-scoped override → Industrial Cool default

**Sprint placement:** end of week 1 / start of week 2 — natural gear-shift between heavy architecture work and testing/documentation wind-down.

---

## Section 12 — Switchyard Brand Assets (UI)

Completed prior to 1.2 sprint start — already live in the Dispatch Whiteboard.

**Deliverables (confirmed complete):**
- Switchyard logo
- Name lockup
- Combined logo + name lockup
- "Powered by Switchyard" treatment
- Light and Dark variants for all assets

**Asset locations:**
- Brand assets: `Switchyard.UI/src/assets/`
- Client palette overrides: `Switchyard.UI/src/assets/clients/`

**Sprint action:** confirm/close at sprint start. No active work required unless a gap is found during review.

---

## Section 13 — .NET API Migration: SQLite to PostgreSQL

Last item in 1.2. First to slide if sprint velocity runs short — no other 1.2 item depends on it.

**Target:** PostgreSQL 16 — consolidates onto the single database engine already running for the Go backend. MS SQL ruled out due to lack of a free tier option.

**Scope:**
- Replace SQLite EF Core provider with Npgsql EF Core provider in both .NET APIs
- Update connection strings in `appsettings.Development.json`
- Replace `EnsureCreated` with initial EF Core migrations setup (foundational work; full migrations strategy is a 1.3 backlog item)
- Validate CQRS read replica sync behavior against PostgreSQL
- Update prerequisites documentation

**Read replica behavior:** confirm `SaveChangesInterceptor` → `Channel<SyncJob>` → `BackgroundService` pattern functions correctly under Npgsql before closing.

---

## Section 14 — Go Backend Architecture

### PlanBOL and Dispatch Whiteboard
- `PlanBOLRecord` — the planning record for a BOL run; owned by the Go backend
- `PlanBOLStop` — individual stops on a plan; stop types: `pickup`, `delivery`, `return_depot` (1.2)
- `DriverBOLAssignment` — driver assigned to a PlanBOL run
- Dispatch Whiteboard — Kanban-style board surfacing active runs, HOS status, and equipment state

### HOS Engine
- Hours of Service tracking per driver
- Warning notifications via SMTP when HOS thresholds approach
- Integrates with `DriverBOLAssignment` for per-run HOS calculation

### Equipment and Breakdown
- Equipment records with health status
- `BreakdownRecord` — captures breakdown events; resolved records feed tow cost into operating cost tracking (1.2)
- Dead-head BOL pairing — tracks empty return legs

### Constraint Solver
- Validates PlanBOL stop sequencing
- Accommodates `return_depot` stop type without modification (1.2)
- Designed to accommodate `transfer` stop type in 1.3 with `DriverBOLAssignment` restructuring

---

## Section 15 — API Surface

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

---

## Section 16 — v1.3 Candidate Features

Items confirmed for 1.3 or deferred pending signal.

**1.3 anchor item:**
- **Mid-BOL transfer stops** — `transfer` stop type on `PlanBOLStop`; formal custody checkpoint for driver and equipment handoffs mid-route; requires `DriverBOLAssignment` restructuring. Promoted to top of backlog — available to pull into 1.2 if sprint velocity allows.

**Backlog (signal-dependent):**
- Rolling refresh tokens for Auth0 sessions in place of fixed-expiry client secrets
- Read replica health endpoint — expose sync lag and `InSync` status
- Migrate from `EnsureCreated` to EF Core migrations for controlled schema evolution
- Extract User Management to a dedicated identity service when the data layer splits
- Scalar branding — Switchyard logo and name above the API title; currently blocked by Scalar's limited logo support in the .NET package

---

## Section 17 — CLAUDE.md Handoff Block

```
Project: Switchyard
Version: 1.2
Branch: dev_1_2

Solution projects:
  Switchyard.InventoryAPI     — .NET 10, port 7000, SQLite (PostgreSQL migration in 1.2)
  Switchyard.LogisticsAPI     — .NET 10, port 7001, SQLite (PostgreSQL migration in 1.2)
  Switchyard.Domain           — .NET 10 shared class library, domain models only (new in 1.2)
  Switchyard.UI               — React/TypeScript, port 5173
  Switchyard-Go               — Go 1.25+, port 8080, PostgreSQL 16 (port 5433)

1.2 priority order:
  1. Operating cost tracking (Go) — base rate/mile, tow rate on resolved breakdowns
  2. Analytics handler refactor (Go) — AnalyticsQuerier interface + testcontainers-go
  3. Advanced analytics and reporting (Go) — revenue vs. profit per BOL/driver/warehouse
  4. Returns — return_depot stop type on PlanBOLStop
  5. Warehouse region attribute — region column on Warehouse, replaces WAREHOUSE_IDS env list
  6. White-label theming — Industrial Cool light/dark defaults, DNS-scoped SCSS overrides
  7. Switchyard brand assets — confirm/close (already live in Dispatch Whiteboard)
  8. Extract Data/ to Switchyard.Domain shared class library
  9. Migrate .NET APIs from SQLite to PostgreSQL (slides first if velocity runs short)

Key constraints:
  - Analytics refactor (item 2) and advanced analytics (item 3) both depend on operating cost tracking (item 1) being complete; item 3 also depends on item 2
  - Switchyard.Domain extraction is compiler-guided — move Data/ folders, fix usings
  - return_depot constraint solver accommodation already exists — completion work only
  - Warehouse region: intra-region only in 1.2; cross-region deferred pending user signal
  - PostgreSQL migration: validate CQRS read replica sync under Npgsql before closing

1.3 anchor: Mid-BOL transfer stops (top of backlog — pull into 1.2 if sprint runs hot)
```
