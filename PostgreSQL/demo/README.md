# Demo Reset Scripts

Two PowerShell scripts for resetting all three databases to a canned demo state.
Each script clears transactional data and reseeds from scratch. Reference/master data
(Warehouses, SKUCatalog, Stores, HOS state limits) is preserved.

## Quick Start

```powershell
# From the PostgreSQL/demo/ directory:
.\reset-company-a.ps1   # Switchyard brand — Monday morning state
.\reset-company-b.ps1   # Digital Parts brand — mid-week complexity state
```

Optional parameters (defaults match the standard local setup):

```powershell
.\reset-company-a.ps1 -Container switchyard-pg -PgUser postgres
```

## Company A — Switchyard (Monday Morning)

**VITE_CLIENT:** `switchyard` (or leave unset — Switchyard is the default brand)

**Login:** Your primary Switchyard Auth0 account

**Board state:**
| Column | Cards |
|--------|-------|
| Draft | 2 (just received) |
| Pending | 1 (route planning in progress) |
| Loading / Ready | 2 loading + 1 validated-ready |
| In Delivery | 3 in-transit (Tom, Keisha, Angela — all green HOS) |
| Delivered | — |
| Available | 4 drivers (Marcus, Diane, Ray, Pete, James, Carol) |
| Maintenance | — |

All drivers have fresh HOS (start of week). Two trailers are at the loading dock.

## Company B — Digital Parts (Mid-Week Complexity)

**VITE_CLIENT:** `digital-parts`

**Login:** Digital Parts demo account (username/password — not Google)

Create this user once via the Switchyard Users page (while logged in as yourself with
`VITE_CLIENT=switchyard`), then record the credentials below:

| Field | Value |
|-------|-------|
| Email | *(fill in after creation)* |
| Password | *(fill in after creation)* |
| Role | Admin |
| Warehouse | WH001 (or unassigned) |

The company brand is controlled entirely by `VITE_CLIENT` — no Auth0 company metadata
is needed. Any username/password user with an Admin role works as the Digital Parts demo login.

**Board state:**
| Column | Cards |
|--------|-------|
| Draft | 1 draft, 1 pending |
| Loading / Ready | 2 loading + 1 validated **⚠ long-wait warning** (36h old) |
| In Delivery | Tom (yellow HOS), Angela (green) |
| Delivered | Carol (dead-head window counting down) |
| Available — Resting | Marcus & Diane (weekly reset), Ray (daily 10h rest) |
| Available — Now | Sandra |
| Maintenance | TK-102 (scheduled), TK-103 (scheduled), TK-105 (depot breakdown), TR-2003 (DOT inspection), TR-2004 (refrigeration service) |

**Active alerts generated automatically:**
- `hos_warning` — Tom Brierley at 9.5h daily (approaching IL daily limit)
- `hos_weekly_limit` — Marcus Webb on weekly restart (62h)
- `hos_weekly_limit` — Diane Kowalski on weekly restart (61.5h)
- `hos_weekly_limit` — Ray Gutierrez on daily 10h rest (47h weekly)

## What Gets Cleared

| Database | Tables cleared |
|----------|---------------|
| `switchyard-go` | driver, hos_window, equipment, maintenance_record, breakdown_record, plan_bol_record, plan_bol_stop, driver_bol_assignment |
| `switchyard_inventory` | Clothing, PPE, Tools |
| `switchyard_logistics` | BillsOfLading, LineEntries |

**Preserved:** Warehouses, SKUCatalog, Stores, WarehouseTransit, HOS state limits, Go warehouse/HOS seed migrations.

## Shared Inventory Seed

Both companies load the same inventory (`seed_inventory_shared.sql`):
- ~120 Clothing items across WH001–WH003 (CLTH001, CLTH002, CLTH004, CLTH007)
- ~140 PPE items across WH001–WH003 (SPPE001, SPPE003, SPPE004, SPPE008)
- ~80 Tool items across WH001–WH003 (PWTL001, PWTL003, PWTL006)
- All items `Projected=false` (at warehouse), dates relative to run time

## Shared Logistics Seed

Both companies load the same BOL history (`seed_logistics_shared.sql`):
- 10 committed BOLs spanning the last 28 days
- Mixed completion states — older BOLs fully processed, recent ones partially open
- Gives analytics charts a rolling 30-day activity curve
