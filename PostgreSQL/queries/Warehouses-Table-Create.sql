-- database: switchyard_inventory

CREATE TABLE IF NOT EXISTS "Warehouses" (
    "WarehouseId"  TEXT NOT NULL PRIMARY KEY,
    "City"         TEXT NOT NULL DEFAULT '',
    "State"        TEXT NOT NULL DEFAULT ''
);