-- database: switchyard_inventory

CREATE TABLE IF NOT EXISTS "WarehouseTransit" (
    "PartitionKey"             TEXT    NOT NULL PRIMARY KEY,
    "OriginWarehouseId"        TEXT    NOT NULL,
    "DestinationWarehouseId"   TEXT    NOT NULL,
    "TransitDays"              INTEGER NOT NULL DEFAULT 1,
    UNIQUE ("OriginWarehouseId", "DestinationWarehouseId"),
    FOREIGN KEY ("OriginWarehouseId") REFERENCES "Warehouses"("WarehouseId"),
    FOREIGN KEY ("DestinationWarehouseId") REFERENCES "Warehouses"("WarehouseId")
);