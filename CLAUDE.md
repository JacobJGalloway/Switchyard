# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

Switchyard — multi-service warehouse management system. Solution file: `Switchyard.InventoryAPI.sln`.

**Services:**
1. `Switchyard.InventoryAPI/` — Inventory API. Tracks warehouse inventory (Clothing, PPE, Tools) backed by SQLite via EF Core.
2. `Switchyard.LogisticsAPI/` — Logistics API. Handles Bill of Lading creation and location-stop processing. References the Inventory project.
3. `Switchyard.UI/` — Client UI for inventory withdrawals. Not yet started (README placeholder only).

**Shared database:** `Sqlite 3 Implementation/WarehouseData.db3` — both APIs resolve this path relative to their content root at startup.

## Build & Run

```bash
dotnet build                                           # build solution
dotnet run --project Switchyard.InventoryAPI      # run inventory API
dotnet run --project Switchyard.LogisticsAPI      # run logistics API
dotnet test                                            # run all tests
dotnet test --filter "FullyQualifiedName~TestName"     # run a single test
```

## Architecture

### Namespace roots
- Inventory API: `Switchyard.InventoryAPI`
- Logistics API: `Switchyard.LogisticsAPI`

### Conventions
- File-scoped namespaces (C# 10+)
- Nullable enabled, implicit usings enabled, target framework net10.0
- Interfaces live in `Interfaces/` subdirectories with the namespace suffix `.Interfaces`
- Domain models implement their corresponding interface (e.g., `Tool : ITool`)

### Switchyard.InventoryAPI layers
- `Models/` — domain models (Clothing, PPE, Tool) and their interfaces
- `Data/InventoryContext.cs` — EF Core DbContext; registered with SQLite in `Program.cs`
- `Data/Interfaces/` — `IUnitOfWork`, `IClothingRepository`, `IPPERepository`, `IToolRepository`
- `Data/Repositories/` — concrete repository implementations + `UnitOfWork.cs`
- `Controllers/` — `ClothingController`, `PPEController`, `ToolController`
- `Tests/Controllers/` — controller unit tests (Moq)
- `Tests/Repositories/` — repository unit tests (in-memory SQLite)

### Switchyard.LogisticsAPI layers
- `Models/` — `BillOfLading`, `LineEntry`, `Store`, `Warehouse` and their interfaces
- `Data/LogisticsContext.cs` — EF Core DbContext (entity config TBD)
- `Data/Interfaces/` — `IBillOfLadingRepository`, `ILineEntryRepository`
- `Data/Repositories/` — `BillOfLadingRepository`, `LineEntryRepository`
- `Services/BillOfLadingService.cs` — service layer stub (currently commented out)
- `Services/Interfaces/IBillOfLadingService.cs` — service interface
- `Controllers/BillOfLadingController.cs` — BOL creation endpoint; writes a formatted `.txt` file to `~/Downloads`
- `Tests/` — empty; tests not yet written

---

## Switchyard Go Architecture (v1.1)

Read `Switchyard-Go/README.md` before writing any code for the Go backend. Key constraints and API reference live there.

**Go project root:** `Switchyard-Go/`
**Module path:** `github.com/JacobJGalloway/switchyard-go`

### Build & Run (Go)

```bash
cd Switchyard-Go
go mod tidy              # resolve dependencies and generate go.sum (required first run)
go build ./...           # build all packages
go run ./cmd/main.go     # run the Go backend (see README for env var loading)
go test ./...            # run all tests
```

### Key constraints (enforce without exception)
1. **CUD authority boundary.** `PlanBOLRecord`, `PlanBOLStop`, and `TruckInventorySnapshot` are Go's domain. Committed BOL stops and inventory writes are .NET's domain. This boundary never crosses.
2. **M2M token lives in `internal/events` only.** No other package requests, caches, or refreshes a token.
3. **`internal/integrations` is the only .NET caller.** No handler, service, or repository may call the .NET system directly.
4. **Business rules are rejections, not warnings.** Empty truck rule, 4-hour dead-head window, and state-level HOS limits are enforced at the service layer and return errors.
5. **Whiteboard column transitions are derived.** Board columns follow from `PlanBOLStatus` and assignment timestamps — not dispatcher input. Read `whiteboard_service.go` before touching board code.
