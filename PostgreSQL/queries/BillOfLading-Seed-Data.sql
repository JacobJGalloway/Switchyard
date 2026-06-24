-- database: switchyard_logistics
-- 11 transactions, 1 header row each. Line entries are in LineEntry-Seed-Data.

INSERT INTO "BillsOfLading" ("PartitionKey", "TransactionId", "Status", "CustomerFirstName", "CustomerLastName", "City", "State") VALUES
('a1b2c3d4-a1b2c3d4', 'a1b2c3d4', 'Submitted', 'John',   'Smith',    'Chicago',      'IL'),
('b2c3d4e5-b2c3d4e5', 'b2c3d4e5', 'Submitted', 'Sarah',  'Johnson',  'Milwaukee',    'WI'),
('c3d4e5f6-c3d4e5f6', 'c3d4e5f6', 'Submitted', 'Mike',   'Davis',    'Indianapolis', 'IN'),
('d4e5f6a7-d4e5f6a7', 'd4e5f6a7', 'Submitted', 'Lisa',   'Martinez', 'Detroit',      'MI'),
('e5f6a7b8-e5f6a7b8', 'e5f6a7b8', 'Submitted', 'Tom',    'Wilson',   'Columbus',     'OH'),
('f6a7b8c9-f6a7b8c9', 'f6a7b8c9', 'Submitted', 'Emily',  'Brown',    'Cleveland',    'OH'),
('a7b8c9d0-a7b8c9d0', 'a7b8c9d0', 'Submitted', 'Carlos', 'Garcia',   'St. Louis',    'MO'),
('b8c9d0e1-b8c9d0e1', 'b8c9d0e1', 'Submitted', 'Karen',  'Lee',      'Kansas City',  'MO'),
('c9d0e1f2-c9d0e1f2', 'c9d0e1f2', 'Submitted', 'David',  'Taylor',   'Minneapolis',  'MN'),
('d0e1f2a3-d0e1f2a3', 'd0e1f2a3', 'Submitted', 'Rachel', 'Anderson', 'Madison',      'WI'),
('e1f2a3b4-e1f2a3b4', 'e1f2a3b4', 'Submitted', 'Marcus', 'Wright',   'Fort Wayne',   'IN')
ON CONFLICT DO NOTHING;
