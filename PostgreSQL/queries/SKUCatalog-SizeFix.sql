-- database: switchyard_inventory
-- Fix: Size column was created as NUMERIC, should be TEXT.

ALTER TABLE "SKUCatalog" ALTER COLUMN "Size" TYPE TEXT USING "Size"::TEXT;
