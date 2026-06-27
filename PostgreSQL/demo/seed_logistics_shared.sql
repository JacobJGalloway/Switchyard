-- database: switchyard_logistics
-- Date-relative BOL and line entry seed. Shared between Company A and Company B demos.
-- CommittedDates are relative to NOW() to keep analytics charts populated and current.

-- ── Bills of Lading ───────────────────────────────────────────────────────────
INSERT INTO "BillsOfLading" ("PartitionKey", "TransactionId", "Status", "CustomerFirstName", "CustomerLastName", "City", "State", "CommittedDate")
VALUES
('bol-demo-01', 'bol-demo-01', 'Submitted', 'Marcus',   'Hendricks',  'Chicago',      'IL', NOW() - INTERVAL '28 days'),
('bol-demo-02', 'bol-demo-02', 'Submitted', 'Patricia',  'Osei',       'Indianapolis', 'IN', NOW() - INTERVAL '24 days'),
('bol-demo-03', 'bol-demo-03', 'Submitted', 'Daniel',    'Kowalczyk',  'Milwaukee',    'WI', NOW() - INTERVAL '21 days'),
('bol-demo-04', 'bol-demo-04', 'Submitted', 'Angela',    'Washington', 'Detroit',      'MI', NOW() - INTERVAL '17 days'),
('bol-demo-05', 'bol-demo-05', 'Submitted', 'James',     'Okafor',     'Louisville',   'KY', NOW() - INTERVAL '14 days'),
('bol-demo-06', 'bol-demo-06', 'Submitted', 'Sandra',    'Martínez',   'Madison',      'WI', NOW() - INTERVAL '10 days'),
('bol-demo-07', 'bol-demo-07', 'Submitted', 'Raymond',   'Gutierrez',  'Rockford',     'IL', NOW() - INTERVAL  '7 days'),
('bol-demo-08', 'bol-demo-08', 'Submitted', 'Diane',     'Halverson',  'South Bend',   'IN', NOW() - INTERVAL  '5 days'),
('bol-demo-09', 'bol-demo-09', 'Submitted', 'Thomas',    'Brierley',   'Ann Arbor',    'MI', NOW() - INTERVAL  '3 days'),
('bol-demo-10', 'bol-demo-10', 'Submitted', 'Carol',     'Metzger',    'Green Bay',    'WI', NOW() - INTERVAL  '1 days');

-- ── Line Entries ──────────────────────────────────────────────────────────────
-- Negative quantity = outbound/delivered. Positive = inbound transfer.
-- Older BOLs are fully processed; recent ones are partially processed.

-- BOL 01 — 28 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-01-le-01', 'bol-demo-01', 'WH001',   'CLTH001', -24, true, NOW() - INTERVAL '28 days'),
('bol-demo-01-le-02', 'bol-demo-01', 'ST0001',   'CLTH001', -12, true, NOW() - INTERVAL '27 days'),
('bol-demo-01-le-03', 'bol-demo-01', 'ST0002',   'CLTH001', -12, true, NOW() - INTERVAL '27 days'),
('bol-demo-01-le-04', 'bol-demo-01', 'WH001',   'SPPE001', -18, true, NOW() - INTERVAL '28 days'),
('bol-demo-01-le-05', 'bol-demo-01', 'ST0003',   'SPPE001',  -9, true, NOW() - INTERVAL '27 days'),
('bol-demo-01-le-06', 'bol-demo-01', 'ST0004',   'SPPE001',  -9, true, NOW() - INTERVAL '27 days');

-- BOL 02 — 24 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-02-le-01', 'bol-demo-02', 'WH002',   'CLTH004', -20, true, NOW() - INTERVAL '24 days'),
('bol-demo-02-le-02', 'bol-demo-02', 'ST0007',   'CLTH004', -10, true, NOW() - INTERVAL '23 days'),
('bol-demo-02-le-03', 'bol-demo-02', 'ST0008',   'CLTH004', -10, true, NOW() - INTERVAL '23 days'),
('bol-demo-02-le-04', 'bol-demo-02', 'WH002',   'PWTL001', -12, true, NOW() - INTERVAL '24 days'),
('bol-demo-02-le-05', 'bol-demo-02', 'ST0009',   'PWTL001',  -6, true, NOW() - INTERVAL '23 days'),
('bol-demo-02-le-06', 'bol-demo-02', 'ST0010',   'PWTL001',  -6, true, NOW() - INTERVAL '23 days');

-- BOL 03 — 21 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-03-le-01', 'bol-demo-03', 'WH001',   'CLTH002', -15, true, NOW() - INTERVAL '21 days'),
('bol-demo-03-le-02', 'bol-demo-03', 'ST0001',   'CLTH002',  -8, true, NOW() - INTERVAL '20 days'),
('bol-demo-03-le-03', 'bol-demo-03', 'ST0002',   'CLTH002',  -7, true, NOW() - INTERVAL '20 days'),
('bol-demo-03-le-04', 'bol-demo-03', 'WH001',   'SPPE003', -20, true, NOW() - INTERVAL '21 days'),
('bol-demo-03-le-05', 'bol-demo-03', 'ST0005',   'SPPE003', -10, true, NOW() - INTERVAL '20 days'),
('bol-demo-03-le-06', 'bol-demo-03', 'ST0006',   'SPPE003', -10, true, NOW() - INTERVAL '20 days');

-- BOL 04 — 17 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-04-le-01', 'bol-demo-04', 'WH002',   'SPPE004', -30, true, NOW() - INTERVAL '17 days'),
('bol-demo-04-le-02', 'bol-demo-04', 'ST0013',   'SPPE004', -15, true, NOW() - INTERVAL '16 days'),
('bol-demo-04-le-03', 'bol-demo-04', 'ST0014',   'SPPE004', -15, true, NOW() - INTERVAL '16 days'),
('bol-demo-04-le-04', 'bol-demo-04', 'WH002',   'PWTL003', -10, true, NOW() - INTERVAL '17 days'),
('bol-demo-04-le-05', 'bol-demo-04', 'ST0015',   'PWTL003',  -5, true, NOW() - INTERVAL '16 days'),
('bol-demo-04-le-06', 'bol-demo-04', 'ST0016',   'PWTL003',  -5, true, NOW() - INTERVAL '16 days');

-- BOL 05 — 14 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-05-le-01', 'bol-demo-05', 'WH001',   'CLTH007', -16, true, NOW() - INTERVAL '14 days'),
('bol-demo-05-le-02', 'bol-demo-05', 'ST0001',   'CLTH007',  -8, true, NOW() - INTERVAL '13 days'),
('bol-demo-05-le-03', 'bol-demo-05', 'ST0003',   'CLTH007',  -8, true, NOW() - INTERVAL '13 days'),
('bol-demo-05-le-04', 'bol-demo-05', 'WH001',   'SPPE008', -24, true, NOW() - INTERVAL '14 days'),
('bol-demo-05-le-05', 'bol-demo-05', 'ST0002',   'SPPE008', -12, true, NOW() - INTERVAL '13 days'),
('bol-demo-05-le-06', 'bol-demo-05', 'ST0004',   'SPPE008', -12, true, NOW() - INTERVAL '13 days');

-- BOL 06 — 10 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-06-le-01', 'bol-demo-06', 'WH002',   'CLTH001', -18, true, NOW() - INTERVAL '10 days'),
('bol-demo-06-le-02', 'bol-demo-06', 'ST0007',   'CLTH001',  -9, true, NOW() - INTERVAL  '9 days'),
('bol-demo-06-le-03', 'bol-demo-06', 'ST0008',   'CLTH001',  -9, true, NOW() - INTERVAL  '9 days'),
('bol-demo-06-le-04', 'bol-demo-06', 'WH002',   'PWTL006', -12, true, NOW() - INTERVAL '10 days'),
('bol-demo-06-le-05', 'bol-demo-06', 'ST0011',   'PWTL006',  -6, true, NOW() - INTERVAL  '9 days'),
('bol-demo-06-le-06', 'bol-demo-06', 'ST0012',   'PWTL006',  -6, true, NOW() - INTERVAL  '9 days');

-- BOL 07 — 7 days ago (fully processed)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-07-le-01', 'bol-demo-07', 'WH001',   'SPPE001', -20, true, NOW() - INTERVAL '7 days'),
('bol-demo-07-le-02', 'bol-demo-07', 'ST0005',   'SPPE001', -10, true, NOW() - INTERVAL '6 days'),
('bol-demo-07-le-03', 'bol-demo-07', 'ST0006',   'SPPE001', -10, true, NOW() - INTERVAL '6 days'),
('bol-demo-07-le-04', 'bol-demo-07', 'WH001',   'CLTH004', -14, true, NOW() - INTERVAL '7 days'),
('bol-demo-07-le-05', 'bol-demo-07', 'ST0001',   'CLTH004',  -7, true, NOW() - INTERVAL '6 days'),
('bol-demo-07-le-06', 'bol-demo-07', 'ST0002',   'CLTH004',  -7, true, NOW() - INTERVAL '6 days');

-- BOL 08 — 5 days ago (warehouse pickups done, store deliveries pending)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-08-le-01', 'bol-demo-08', 'WH002',   'CLTH002', -16, true,  NOW() - INTERVAL '5 days'),
('bol-demo-08-le-02', 'bol-demo-08', 'ST0013',   'CLTH002',  -8, false, NULL),
('bol-demo-08-le-03', 'bol-demo-08', 'ST0014',   'CLTH002',  -8, false, NULL),
('bol-demo-08-le-04', 'bol-demo-08', 'WH002',   'SPPE004', -12, true,  NOW() - INTERVAL '5 days'),
('bol-demo-08-le-05', 'bol-demo-08', 'ST0015',   'SPPE004',  -6, false, NULL);

-- BOL 09 — 3 days ago (warehouse pickup done, 1 of 2 store stops done)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-09-le-01', 'bol-demo-09', 'WH001',   'PWTL001', -12, true,  NOW() - INTERVAL '3 days'),
('bol-demo-09-le-02', 'bol-demo-09', 'ST0001',   'PWTL001',  -6, true,  NOW() - INTERVAL '2 days'),
('bol-demo-09-le-03', 'bol-demo-09', 'ST0003',   'PWTL001',  -6, false, NULL),
('bol-demo-09-le-04', 'bol-demo-09', 'WH001',   'CLTH007', -10, true,  NOW() - INTERVAL '3 days'),
('bol-demo-09-le-05', 'bol-demo-09', 'ST0002',   'CLTH007',  -5, true,  NOW() - INTERVAL '2 days'),
('bol-demo-09-le-06', 'bol-demo-09', 'ST0004',   'CLTH007',  -5, false, NULL);

-- BOL 10 — yesterday (just picked up, no store stops done)
INSERT INTO "LineEntries" ("PartitionKey", "TransactionId", "LocationId", "SKUMarker", "Quantity", "IsProcessed", "ProcessedDate")
VALUES
('bol-demo-10-le-01', 'bol-demo-10', 'WH002',   'SPPE001', -24, true,  NOW() - INTERVAL '1 days'),
('bol-demo-10-le-02', 'bol-demo-10', 'ST0007',   'SPPE001', -12, false, NULL),
('bol-demo-10-le-03', 'bol-demo-10', 'ST0008',   'SPPE001', -12, false, NULL),
('bol-demo-10-le-04', 'bol-demo-10', 'WH002',   'PWTL006',  -8, true,  NOW() - INTERVAL '1 days'),
('bol-demo-10-le-05', 'bol-demo-10', 'ST0009',   'PWTL006',  -4, false, NULL),
('bol-demo-10-le-06', 'bol-demo-10', 'ST0010',   'PWTL006',  -4, false, NULL);
