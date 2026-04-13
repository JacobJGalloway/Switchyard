# SQLite 3 Implementation

All query files are located in the `queries/` sub-folder. The target database file is `WarehouseData.db3`.

> **Note:** You may notice `*Read.db3` files (e.g. `WarehouseInventoryRead.db3`) are greyed out in your editor — they are gitignored because each API creates and populates its own read replica on startup via `EnsureCreated()` and an initial sync job. Only `WarehouseData.db3` is tracked in source control.

## Setting Up from an Empty Database

Run the queries in the two phases below, in order. Each phase must complete before starting the next.

### Phase 1 — Create Tables

Run these files in any order:

- `queries/Warehouses-Table-Create.sqlite3-query`
- `queries/SKUCatalog-Table-Create.sqlite3-query`
- `queries/Clothing-Table-Create.sqlite3-query`
- `queries/PPE-Table-Create.sqlite3-query`
- `queries/Tools-Table-Create.sqlite3-query`
- `queries/BillOfLading-Table-Create.sqlite3-query`
- `queries/WarehouseTransit-Table-Create.sqlite3-query`
- `queries/Stores-Table-Create.sqlite3-query`
- `queries/LineEntry-Table-Create.sqlite3-query`

### Phase 2 — Seed Data

Foreign key constraints are enforced, so parent tables must be seeded before their dependents.

**Step 1** — seed parent tables first (in any order):

- `queries/Warehouses-Seed-Data.sql3-query`
- `queries/SKUCatalog-Seed-Data.sql3-query`

**Step 2** — seed dependent tables (in any order):

- `queries/Clothing-Seed-Data.sql3-query` _(depends on SKUCatalog)_
- `queries/PPE-Seed-Data.sql3-query` _(depends on SKUCatalog)_
- `queries/Tools-Seed-Data.sql3-query` _(depends on SKUCatalog)_
- `queries/BillOfLading-Seed-Data.sql3-query` _(depends on SKUCatalog)_
- `queries/WarehouseTransit-Seed-Data.sql3-query` _(depends on Warehouses)_
- `queries/Stores-Seed-Data.sql3-query` _(depends on Warehouses)_

**Step 3** — seed tables that depend on Step 2 results:

- `queries/LineEntry-Seed-Data.sqlite3query` _(depends on BillOfLading and SKUCatalog)_

## Running Queries

SQLite does not enforce foreign keys by default. Run this pragma once at the start of each connection before executing any queries:

```sql
PRAGMA foreign_keys = ON;
```

**sqlite3 CLI:**
```bash
sqlite3 WarehouseData.db3 < queries/<filename>
```

**DB Browser for SQLite:** Open `WarehouseData.db3`, paste the file contents into the Execute SQL tab, and run.
