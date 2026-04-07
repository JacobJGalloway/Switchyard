# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

Multi-service warehouse management system. Solution file: `WarehouseInventoryAPI-Claude.sln`.

**Services:**
1. `WarehouseInventoryAPI-Claude/` — Inventory API. Tracks warehouse inventory (Clothing, PPE, Tools) backed by SQLite via EF Core.
2. `WarehouseLogisticsAPI-Claude/` — Logistics API. Handles Bill of Lading creation and location-stop processing. References the Inventory project.
3. `WarehouseSalesUI-Claude/` — Client UI for inventory withdrawals. Not yet started (README placeholder only).

**Shared database:** `Sqlite 3 Implementation/WarehouseData.db3` — both APIs resolve this path relative to their content root at startup.

## Build & Run

```bash
dotnet build                                           # build solution
dotnet run --project WarehouseInventoryAPI-Claude      # run inventory API
dotnet run --project WarehouseLogisticsAPI-Claude      # run logistics API
dotnet test                                            # run all tests
dotnet test --filter "FullyQualifiedName~TestName"     # run a single test
```

## Architecture

### Namespace roots
- Inventory API: `WarehouseInventory_Claude`
- Logistics API: `WarehouseLogistics_Claude`

### Conventions
- File-scoped namespaces (C# 10+)
- Nullable enabled, implicit usings enabled, target framework net10.0
- Interfaces live in `Interfaces/` subdirectories with the namespace suffix `.Interfaces`
- Domain models implement their corresponding interface (e.g., `Tool : ITool`)

### WarehouseInventoryAPI-Claude layers
- `Models/` — domain models (Clothing, PPE, Tool) and their interfaces
- `Data/InventoryContext.cs` — EF Core DbContext; registered with SQLite in `Program.cs`
- `Data/Interfaces/` — `IUnitOfWork`, `IClothingRepository`, `IPPERepository`, `IToolRepository`
- `Data/Repositories/` — concrete repository implementations + `UnitOfWork.cs`
- `Controllers/` — `ClothingController`, `PPEController`, `ToolController`
- `Tests/Controllers/` — controller unit tests (Moq)
- `Tests/Repositories/` — repository unit tests (in-memory SQLite)

### WarehouseLogisticsAPI-Claude layers
- `Models/` — `BillOfLading`, `LineEntry`, `Store`, `Warehouse` and their interfaces
- `Data/LogisticsContext.cs` — EF Core DbContext (entity config TBD)
- `Data/Interfaces/` — `IBillOfLadingRepository`, `ILineEntryRepository`
- `Data/Repositories/` — `BillOfLadingRepository`, `LineEntryRepository`
- `Services/BillOfLadingService.cs` — service layer stub (currently commented out)
- `Services/Interfaces/IBillOfLadingService.cs` — service interface
- `Controllers/BillOfLadingController.cs` — BOL creation endpoint; writes a formatted `.txt` file to `~/Downloads`
- `Tests/` — empty; tests not yet written
