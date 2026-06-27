-- database: switchyard_logistics
-- 11 transactions, 1 header row each. Line entries are in LineEntry-Seed-Data.
-- CommittedDate spread across last 30 days for analytics demo data.

INSERT INTO "BillsOfLading" ("PartitionKey", "TransactionId", "Status", "CustomerFirstName", "CustomerLastName", "City", "State", "CommittedDate") VALUES
('a1b2c3d4-a1b2c3d4', 'a1b2c3d4', 'Submitted', 'John',   'Smith',    'Chicago',      'IL', '2026-05-25'),
('b2c3d4e5-b2c3d4e5', 'b2c3d4e5', 'Submitted', 'Sarah',  'Johnson',  'Milwaukee',    'WI', '2026-05-28'),
('c3d4e5f6-c3d4e5f6', 'c3d4e5f6', 'Submitted', 'Mike',   'Davis',    'Indianapolis', 'IN', '2026-05-31'),
('d4e5f6a7-d4e5f6a7', 'd4e5f6a7', 'Submitted', 'Lisa',   'Martinez', 'Detroit',      'MI', '2026-06-03'),
('e5f6a7b8-e5f6a7b8', 'e5f6a7b8', 'Submitted', 'Tom',    'Wilson',   'Columbus',     'OH', '2026-06-06'),
('f6a7b8c9-f6a7b8c9', 'f6a7b8c9', 'Submitted', 'Emily',  'Brown',    'Cleveland',    'OH', '2026-06-09'),
('a7b8c9d0-a7b8c9d0', 'a7b8c9d0', 'Submitted', 'Carlos', 'Garcia',   'St. Louis',    'MO', '2026-06-12'),
('b8c9d0e1-b8c9d0e1', 'b8c9d0e1', 'Submitted', 'Karen',  'Lee',      'Kansas City',  'MO', '2026-06-15'),
('c9d0e1f2-c9d0e1f2', 'c9d0e1f2', 'Submitted', 'David',  'Taylor',   'Minneapolis',  'MN', '2026-06-18'),
('d0e1f2a3-d0e1f2a3', 'd0e1f2a3', 'Submitted', 'Rachel', 'Anderson', 'Madison',      'WI', '2026-06-21'),
('e1f2a3b4-e1f2a3b4', 'e1f2a3b4', 'Submitted', 'Marcus', 'Wright',   'Fort Wayne',   'IN', '2026-06-24')
ON CONFLICT DO NOTHING;
