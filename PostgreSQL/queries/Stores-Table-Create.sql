-- database: switchyard_logistics

CREATE TABLE IF NOT EXISTS "Stores" (
    "PartitionKey"     TEXT NOT NULL PRIMARY KEY,
    "StoreId"          TEXT NOT NULL DEFAULT '',
    "BaseWarehouseId"  TEXT NOT NULL DEFAULT '',
    "City"             TEXT NOT NULL DEFAULT '',
    "State"            TEXT NOT NULL DEFAULT '',
    FOREIGN KEY ("BaseWarehouseId") REFERENCES "Warehouses"("WarehouseId")
);