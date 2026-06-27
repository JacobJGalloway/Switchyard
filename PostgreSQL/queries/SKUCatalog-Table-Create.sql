-- database: switchyard_inventory

CREATE TABLE IF NOT EXISTS "SKUCatalog" (
    "SKUMarker"  TEXT    NOT NULL PRIMARY KEY,
    "Category"   TEXT    NOT NULL DEFAULT '',
    "Type"       TEXT    NOT NULL DEFAULT '',
    "Color"      TEXT    NOT NULL DEFAULT '',
    "Size"       TEXT    NOT NULL DEFAULT '',
    "UnitPrice"  NUMERIC NOT NULL DEFAULT 0.0
);