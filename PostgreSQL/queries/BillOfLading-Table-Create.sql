-- database: switchyard_logistics

CREATE TABLE IF NOT EXISTS "BillsOfLading" (
    "PartitionKey"       TEXT NOT NULL PRIMARY KEY,
    "TransactionId"      TEXT NOT NULL UNIQUE DEFAULT '',
    "Status"             TEXT NOT NULL DEFAULT 'Pending',
    "CustomerFirstName"  TEXT NOT NULL DEFAULT '',
    "CustomerLastName"   TEXT NOT NULL DEFAULT '',
    "City"               TEXT NOT NULL DEFAULT '',
    "State"              TEXT NOT NULL DEFAULT ''
);