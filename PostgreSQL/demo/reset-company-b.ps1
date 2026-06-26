# Switchyard Demo Reset - Company B (Mid-Week Complexity / Digital Parts Brand)
# Clears all transactional data and reseeds to a busy mid-week operational state.
# Master/reference data (Warehouses, SKUCatalog, Stores, HOS limits) is preserved.
#
# Usage:
#   .\reset-company-b.ps1
#   .\reset-company-b.ps1 -Container my-pg-container -PgUser myuser

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
    if ($LASTEXITCODE -ne 0) { throw "SQL error in $File - aborting." }
}

Write-Host ""
Write-Host "=== Switchyard Demo Reset - Company B ===" -ForegroundColor Cyan
Write-Host "    Brand: Digital Parts" -ForegroundColor DarkCyan
Write-Host "    State: Mid-week - 3 resting, 2 in-transit, 1 dead-head, equipment issues" -ForegroundColor DarkCyan
Write-Host ""

# Step 1: Clear Go whiteboard data
Write-Host "[1/4] Clearing Go whiteboard data..." -ForegroundColor Yellow
Invoke-Sql "switchyard" "$here\_clear_go.sql"

# Step 2: Clear .NET inventory
Write-Host "[2/4] Clearing inventory data..." -ForegroundColor Yellow
Invoke-Sql "switchyard_inventory" "$here\_clear_inventory.sql"

# Step 3: Clear .NET logistics
Write-Host "[3/4] Clearing logistics data..." -ForegroundColor Yellow
Invoke-Sql "switchyard_logistics" "$here\_clear_logistics.sql"

# Step 4: Reseed all three databases
Write-Host "[4/4] Seeding Company B data..." -ForegroundColor Yellow
Invoke-Sql "switchyard"            "$here\seed_go_company_b.sql"
Invoke-Sql "switchyard_inventory"  "$here\seed_inventory_shared.sql"
Invoke-Sql "switchyard_logistics"  "$here\seed_logistics_shared.sql"

Write-Host ""
Write-Host "Company B reset complete." -ForegroundColor Green
Write-Host ""
Write-Host "  Next steps:" -ForegroundColor White
Write-Host "  1. Set VITE_CLIENT=digital-parts in Switchyard.UI/.env" -ForegroundColor Gray
Write-Host "  2. Log in with the Digital Parts demo account (see PostgreSQL/demo/README.md)" -ForegroundColor Gray
Write-Host "  3. Restart the UI dev server so the new VITE_CLIENT takes effect" -ForegroundColor Gray
Write-Host ""
