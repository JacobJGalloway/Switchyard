-- database: switchyard_logistics
-- Backfill CommittedDate for existing BOL seed rows.

UPDATE "BillsOfLading" SET "CommittedDate" = CASE "TransactionId"
  WHEN 'a1b2c3d4' THEN '2026-05-25'::date
  WHEN 'b2c3d4e5' THEN '2026-05-28'::date
  WHEN 'c3d4e5f6' THEN '2026-05-31'::date
  WHEN 'd4e5f6a7' THEN '2026-06-03'::date
  WHEN 'e5f6a7b8' THEN '2026-06-06'::date
  WHEN 'f6a7b8c9' THEN '2026-06-09'::date
  WHEN 'a7b8c9d0' THEN '2026-06-12'::date
  WHEN 'b8c9d0e1' THEN '2026-06-15'::date
  WHEN 'c9d0e1f2' THEN '2026-06-18'::date
  WHEN 'd0e1f2a3' THEN '2026-06-21'::date
  WHEN 'e1f2a3b4' THEN '2026-06-24'::date
END;
