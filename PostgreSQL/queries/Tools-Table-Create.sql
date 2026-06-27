-- database: switchyard_inventory
-- NOTE: Managed by EF Core (InventoryContext). Run "dotnet ef database update" instead of this file.
-- Kept as structural reference only.

CREATE TABLE IF NOT EXISTS "Tools" (
    "PartitionKey"       TEXT      NOT NULL PRIMARY KEY,
    "RowKey"             TEXT      NOT NULL DEFAULT '',
    "LocationId"         TEXT      NOT NULL DEFAULT '',
    "SKUMarker"          TEXT      NOT NULL DEFAULT '',
    "UnloadedDate"       TIMESTAMPTZ NOT NULL,
    "Projected"          INTEGER   NOT NULL DEFAULT 1,
    "UnitPrice"          NUMERIC   NOT NULL DEFAULT 0.0,
    "PriceCurrency"      TEXT      NOT NULL DEFAULT 'USD',
    "PriceEffectiveDate" TIMESTAMPTZ
);