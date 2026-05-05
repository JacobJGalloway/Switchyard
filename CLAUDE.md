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

Read `ARCHITECTURE.md` before writing any code for the Go backend.

**Go project root:** `Switchyard-Go/`
**Module path:** `github.com/JacobJGalloway/switchyard-go`

### Build & Run (Go)

```bash
cd Switchyard-Go
go mod tidy              # resolve dependencies and generate go.sum (required first run)
go build ./...           # build all packages
go run ./cmd/main.go     # run the Go backend
go test ./...            # run all tests
```

### Priority reading order before touching Go code
1. Section 2 — CUD authority boundary. PlanBOL entities are Go's domain. Committed BOL stops are .NET's domain. These never cross.
2. Section 3 — The event handler is the single Go entry point. The M2M token lives here and nowhere else.
3. Section 4 — Business rules are hard constraints. Empty truck rule, 4-hour dead-head window, state-level HOS limits are enforced at the service layer. They are rejections, not warnings.
4. Section 5 — Read the full Kanban board design before touching whiteboard code. Column transitions are driven by assignment state.
5. Section 13 — The integrations adapter is the only place that calls the .NET system. No exceptions.

### Go implementation order
1. ✅ go.mod and project scaffold
2. ✅ Domain models `/internal/models`
3. ✅ Repository interfaces `/internal/repository/interfaces.go`
4. Migration files `/internal/migrations`
5. HOSLimit seed data
6. Integration adapters `/internal/integrations`
7. Event handler `/internal/events`
8. Service layer — route_planner and hos_service first
9. Whiteboard service
10. Notification service
11. Handlers `/internal/handlers`
12. Web templates `/web/templates`
