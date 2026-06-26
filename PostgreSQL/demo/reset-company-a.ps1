# Switchyard Demo Reset — Company A (Monday Morning / Switchyard Brand)
# Clears all transactional data and reseeds to a Monday morning operational state.
# Master/reference data (Warehouses, SKUCatalog, Stores, HOS limits) is preserved.
#
# Usage:
#   .\reset-company-a.ps1
#   .\reset-company-a.ps1 -Container my-pg-container -PgUser myuser

param(
    [string]$Container = "switchyard-pg",
    [string]$PgUser    = "postgres"
)

$ErrorActionPreference = "Stop"
$here = $PSScriptRoot

function Invoke-Sql {
    param([string]$Database, [string]$File)
    Write-Host "  [$Database] $([System.IO.Path]::GetFileName($File))..."
    Get-Content $File | docker exec -i $Container psql -U $PgUser -d $Database -v ON_ERROR_STOP=1
    if ($LASTEXITCODE -ne 0) { throw "SQL error in $File — aborting." }
}

Write-Host ""
Write-Host "=== Switchyard Demo Reset — Company A ===" -ForegroundColor Cyan
Write-Host "    Brand: Switchyard (default)" -ForegroundColor DarkCyan
Write-Host "    State: Monday morning — 3 early departures, 2 trailers loading" -ForegroundColor DarkCyan
Write-Host ""

# ── Step 1: Clear Go whiteboard data ──────────────────────────────────────────
Write-Host "[1/4] Clearing Go whiteboard data..." -ForegroundColor Yellow
Invoke-Sql "switchyard-go" "$here\_clear_go.sql"

# ── Step 2: Clear .NET inventory ──────────────────────────────────────────────
Write-Host "[2/4] Clearing inventory data..." -ForegroundColor Yellow
Invoke-Sql "switchyard_inventory" "$here\_clear_inventory.sql"

# ── Step 3: Clear .NET logistics ──────────────────────────────────────────────
Write-Host "[3/4] Clearing logistics data..." -ForegroundColor Yellow
Invoke-Sql "switchyard_logistics" "$here\_clear_logistics.sql"

# ── Step 4: Reseed all three databases ────────────────────────────────────────
Write-Host "[4/4] Seeding Company A data..." -ForegroundColor Yellow
Invoke-Sql "switchyard-go"         "$here\seed_go_company_a.sql"
Invoke-Sql "switchyard_inventory"  "$here\seed_inventory_shared.sql"
Invoke-Sql "switchyard_logistics"  "$here\seed_logistics_shared.sql"

Write-Host ""
Write-Host "Company A reset complete." -ForegroundColor Green
Write-Host ""
Write-Host "  Next steps:" -ForegroundColor White
Write-Host "  1. Set VITE_CLIENT=switchyard  (or leave unset — Switchyard is the default)" -ForegroundColor Gray
Write-Host "  2. Log in with your primary Switchyard account" -ForegroundColor Gray
Write-Host "  3. Restart the UI dev server if VITE_CLIENT changed" -ForegroundColor Gray
Write-Host ""
