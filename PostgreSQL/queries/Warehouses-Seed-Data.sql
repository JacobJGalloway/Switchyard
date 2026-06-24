-- database: switchyard_inventory

INSERT INTO "Warehouses" ("WarehouseId", "City", "State") VALUES
('WH001', 'Chicago',       'IL'),
('WH002', 'Indianapolis',  'IN'),
('WH003', 'Milwaukee',     'WI'),
('WH004', 'Grand Rapids',  'MI'),
('WH005', 'Des Moines',    'IA'),
('WH006', 'Kansas City',   'MO'),
('WH007', 'Lexington',     'KY'),
('WH008', 'Memphis',       'TN'),
('WH009', 'Little Rock',   'AR')
ON CONFLICT DO NOTHING;