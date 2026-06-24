-- Backfill UnitPrice and PriceCurrency on Clothing, PPE, and Tools tables
-- after the AddUnitPrice EF Core migration has been applied.
-- Prices are sourced from SKUCatalog and matched by SKUMarker.
-- Target: PostgreSQL (switchyard_inventory) â€” table and column names must be quoted.

-- Clothing
UPDATE "Clothing" SET "UnitPrice" = 28.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH001';
UPDATE "Clothing" SET "UnitPrice" = 45.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH002';
UPDATE "Clothing" SET "UnitPrice" = 79.95,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH003';
UPDATE "Clothing" SET "UnitPrice" = 24.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH004';
UPDATE "Clothing" SET "UnitPrice" = 67.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH005';
UPDATE "Clothing" SET "UnitPrice" = 32.00,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH006';
UPDATE "Clothing" SET "UnitPrice" = 89.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH007';
UPDATE "Clothing" SET "UnitPrice" = 72.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH008';
UPDATE "Clothing" SET "UnitPrice" = 54.95,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH009';
UPDATE "Clothing" SET "UnitPrice" = 22.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH010';
UPDATE "Clothing" SET "UnitPrice" = 65.00,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH011';
UPDATE "Clothing" SET "UnitPrice" = 31.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH012';
UPDATE "Clothing" SET "UnitPrice" = 14.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH013';
UPDATE "Clothing" SET "UnitPrice" = 12.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH014';
UPDATE "Clothing" SET "UnitPrice" = 26.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH015';
UPDATE "Clothing" SET "UnitPrice" = 58.00,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH016';
UPDATE "Clothing" SET "UnitPrice" = 84.95,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH017';
UPDATE "Clothing" SET "UnitPrice" = 29.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH018';
UPDATE "Clothing" SET "UnitPrice" = 47.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH019';
UPDATE "Clothing" SET "UnitPrice" = 33.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH020';
UPDATE "Clothing" SET "UnitPrice" = 23.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH021';
UPDATE "Clothing" SET "UnitPrice" = 13.99,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH022';
UPDATE "Clothing" SET "UnitPrice" = 76.00,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH023';
UPDATE "Clothing" SET "UnitPrice" = 25.50,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH024';
UPDATE "Clothing" SET "UnitPrice" = 34.00,  "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'CLTH025';

-- PPE
UPDATE "PPE" SET "UnitPrice" =  34.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE001';
UPDATE "PPE" SET "UnitPrice" =  42.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE002';
UPDATE "PPE" SET "UnitPrice" =  18.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE003';
UPDATE "PPE" SET "UnitPrice" =  12.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE004';
UPDATE "PPE" SET "UnitPrice" =  21.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE005';
UPDATE "PPE" SET "UnitPrice" =  38.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE006';
UPDATE "PPE" SET "UnitPrice" =  29.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE007';
UPDATE "PPE" SET "UnitPrice" =  44.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE008';
UPDATE "PPE" SET "UnitPrice" =  22.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE009';
UPDATE "PPE" SET "UnitPrice" =  11.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE010';
UPDATE "PPE" SET "UnitPrice" = 124.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE011';
UPDATE "PPE" SET "UnitPrice" = 124.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE012';
UPDATE "PPE" SET "UnitPrice" = 189.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE013';
UPDATE "PPE" SET "UnitPrice" =  89.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE014';
UPDATE "PPE" SET "UnitPrice" =   8.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE015';
UPDATE "PPE" SET "UnitPrice" =  67.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE016';
UPDATE "PPE" SET "UnitPrice" =   8.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE017';
UPDATE "PPE" SET "UnitPrice" =   9.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE018';
UPDATE "PPE" SET "UnitPrice" =  14.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE019';
UPDATE "PPE" SET "UnitPrice" =  32.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE020';
UPDATE "PPE" SET "UnitPrice" =  38.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE021';
UPDATE "PPE" SET "UnitPrice" =  15.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE022';
UPDATE "PPE" SET "UnitPrice" =  27.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE023';
UPDATE "PPE" SET "UnitPrice" = 134.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE024';
UPDATE "PPE" SET "UnitPrice" =   9.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'SPPE025';

-- Tools
UPDATE "Tools" SET "UnitPrice" = 149.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL001';
UPDATE "Tools" SET "UnitPrice" = 189.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL002';
UPDATE "Tools" SET "UnitPrice" =  89.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL003';
UPDATE "Tools" SET "UnitPrice" = 159.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL004';
UPDATE "Tools" SET "UnitPrice" = 134.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL005';
UPDATE "Tools" SET "UnitPrice" =  79.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL006';
UPDATE "Tools" SET "UnitPrice" = 114.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL007';
UPDATE "Tools" SET "UnitPrice" = 289.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL008';
UPDATE "Tools" SET "UnitPrice" = 199.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL009';
UPDATE "Tools" SET "UnitPrice" = 119.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL010';
UPDATE "Tools" SET "UnitPrice" = 124.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL011';
UPDATE "Tools" SET "UnitPrice" = 349.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL012';
UPDATE "Tools" SET "UnitPrice" = 449.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL013';
UPDATE "Tools" SET "UnitPrice" = 649.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL014';
UPDATE "Tools" SET "UnitPrice" = 549.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL015';
UPDATE "Tools" SET "UnitPrice" =  74.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL016';
UPDATE "Tools" SET "UnitPrice" = 524.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL017';
UPDATE "Tools" SET "UnitPrice" =  99.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL018';
UPDATE "Tools" SET "UnitPrice" =  64.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL019';
UPDATE "Tools" SET "UnitPrice" =  89.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL020';
UPDATE "Tools" SET "UnitPrice" = 139.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL021';
UPDATE "Tools" SET "UnitPrice" = 169.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL022';
UPDATE "Tools" SET "UnitPrice" =  95.00, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL023';
UPDATE "Tools" SET "UnitPrice" = 149.50, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL024';
UPDATE "Tools" SET "UnitPrice" = 269.99, "PriceCurrency" = 'USD' WHERE "SKUMarker" = 'PWTL025';
