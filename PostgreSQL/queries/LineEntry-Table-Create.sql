-- database: switchyard_logistics

CREATE TABLE IF NOT EXISTS "LineEntries" (
    "PartitionKey"  TEXT    NOT NULL PRIMARY KEY,
    "TransactionId" TEXT    NOT NULL DEFAULT '',
    "LocationId"    TEXT    NOT NULL DEFAULT '',
    "SKUMarker"     TEXT    NOT NULL DEFAULT '',
    "Quantity"      INTEGER NOT NULL DEFAULT 0,
    "IsProcessed"   INTEGER NOT NULL DEFAULT 0,
    "ProcessedDate" TEXT,
    FOREIGN KEY ("TransactionId") REFERENCES "BillsOfLading"("TransactionId"),
    FOREIGN KEY ("SKUMarker") REFERENCES "SKUCatalog"("SKUMarker")
);