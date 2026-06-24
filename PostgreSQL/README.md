# PostgreSQL

All query files are in the `queries/` sub-folder. Both .NET APIs and the Go backend run against PostgreSQL on port 5433 (`switchyard-pg` Docker container).

## Databases

| Database | Owner |
|---|---|
| `switchyard_inventory` | InventoryAPI (EF Core) |
| `switchyard_logistics` | LogisticsAPI (EF Core) |

## Setting Up from Scratch

### Phase 1 — Apply EF Core Migrations

EF Core owns `Clothing`, `PPE`, and `Tools` in `switchyard_inventory`. Run migrations first:

```bash
dotnet ef database update --project Switchyard.InventoryAPI --context InventoryContext
dotnet ef database update --project Switchyard.LogisticsAPI --context LogisticsContext
```

### Phase 2 — Create Non-EF Tables

Run these against the indicated database:

**switchyard_inventory:**
- `queries/SKUCatalog-Table-Create.sql`
- `queries/Warehouses-Table-Create.sql`
- `queries/WarehouseTransit-Table-Create.sql`

**switchyard_logistics:**
- `queries/BillOfLading-Table-Create.sql`
- `queries/Stores-Table-Create.sql`
- `queries/LineEntry-Table-Create.sql`

> `Clothing-Table-Create.sql`, `PPE-Table-Create.sql`, and `Tools-Table-Create.sql` are reference only — EF Core manages those tables.

### Phase 3 — Seed Data

Foreign key constraints apply; parent tables must be seeded before dependents.

**Step 1 — seed parent tables:**
- `queries/Warehouses-Seed-Data.sql` → `switchyard_inventory`
- `queries/SKUCatalog-Seed-Data.sql` → `switchyard_inventory`

**Step 2 — seed dependent tables:**
- `queries/Clothing-Seed-Data.sql` → `switchyard_inventory`
- `queries/PPE-Seed-Data.sql` → `switchyard_inventory`
- `queries/Tools-Seed-Data.sql` → `switchyard_inventory`
- `queries/WarehouseTransit-Seed-Data.sql` → `switchyard_inventory`
- `queries/BillOfLading-Seed-Data.sql` → `switchyard_logistics`
- `queries/Stores-Seed-Data.sql` → `switchyard_logistics`

**Step 3 — tables that depend on Step 2:**
- `queries/LineEntry-Seed-Data.sql` → `switchyard_logistics`

**Step 4 — backfill unit prices (after items are seeded):**
- `queries/UnitPrice-Backfill.sql` → `switchyard_inventory`

## Running Queries

```powershell
Get-Content "PostgreSQL/queries/<filename>.sql" | docker exec -i switchyard-pg psql -U postgres -d <database>
```
