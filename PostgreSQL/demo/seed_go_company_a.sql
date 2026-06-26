-- Go DB seed — Company A (Monday Morning / Switchyard Brand)
-- database: switchyard-go
-- All timestamps are NOW()-relative. Board state: 2 draft, 1 pending, 2 loading,
-- 1 validated-ready, 3 in-transit (early departures), 2 trailers loading.

-- ── Drivers ──────────────────────────────────────────────────────────────────
INSERT INTO driver (id, name, home_warehouse_id, auth0_user_id, license_state, is_active) VALUES
  ('a1000000-0000-0000-0000-000000000001', 'Marcus Webb',     'WH001', 'auth0|testdriver001', 'IL', true),
  ('a1000000-0000-0000-0000-000000000002', 'Diane Kowalski',  'WH001', 'auth0|testdriver002', 'IL', true),
  ('a1000000-0000-0000-0000-000000000003', 'Ray Gutierrez',   'WH002', 'auth0|testdriver003', 'IN', true),
  ('a1000000-0000-0000-0000-000000000004', 'Sandra Okafor',   'WH002', 'auth0|testdriver004', 'IN', true),
  ('a1000000-0000-0000-0000-000000000005', 'Tom Brierley',    'WH003', 'auth0|testdriver005', 'IL', true),
  ('a1000000-0000-0000-0000-000000000006', 'Keisha Drummond', 'WH003', 'auth0|testdriver006', 'IL', true),
  ('a1000000-0000-0000-0000-000000000007', 'Pete Halverson',  'WH004', 'auth0|testdriver007', 'MI', true),
  ('a1000000-0000-0000-0000-000000000008', 'Angela Torres',   'WH004', 'auth0|testdriver008', 'MI', true),
  ('a1000000-0000-0000-0000-000000000009', 'James Whitfield', 'WH005', 'auth0|testdriver009', 'WI', true),
  ('a1000000-0000-0000-0000-000000000010', 'Carol Metzger',   'WH005', 'auth0|testdriver010', 'WI', true);

-- ── HOS Windows — Monday morning, everyone fresh ──────────────────────────────
-- Available drivers: 1,2,3,4,7,9,10 (all low hours)
-- Departed this morning: 5 (Tom), 6 (Keisha), 8 (Angela)
INSERT INTO hos_window (id, driver_id, window_start, daily_hours_used, weekly_hours_used) VALUES
  ('b1000000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', DATE_TRUNC('week', NOW()),  0.5,  4.0),
  ('b1000000-0000-0000-0000-000000000002', 'a1000000-0000-0000-0000-000000000002', DATE_TRUNC('week', NOW()),  0.0,  0.0),
  ('b1000000-0000-0000-0000-000000000003', 'a1000000-0000-0000-0000-000000000003', DATE_TRUNC('week', NOW()),  0.0,  0.0),
  ('b1000000-0000-0000-0000-000000000004', 'a1000000-0000-0000-0000-000000000004', DATE_TRUNC('week', NOW()),  1.0,  7.5),
  ('b1000000-0000-0000-0000-000000000005', 'a1000000-0000-0000-0000-000000000005', DATE_TRUNC('week', NOW()),  3.0, 18.5),
  ('b1000000-0000-0000-0000-000000000006', 'a1000000-0000-0000-0000-000000000006', DATE_TRUNC('week', NOW()),  2.5, 15.0),
  ('b1000000-0000-0000-0000-000000000007', 'a1000000-0000-0000-0000-000000000007', DATE_TRUNC('week', NOW()),  0.0,  0.0),
  ('b1000000-0000-0000-0000-000000000008', 'a1000000-0000-0000-0000-000000000008', DATE_TRUNC('week', NOW()),  2.0, 12.0),
  ('b1000000-0000-0000-0000-000000000009', 'a1000000-0000-0000-0000-000000000009', DATE_TRUNC('week', NOW()),  0.0,  0.0),
  ('b1000000-0000-0000-0000-000000000010', 'a1000000-0000-0000-0000-000000000010', DATE_TRUNC('week', NOW()),  0.0,  2.5);

-- ── Equipment — Trucks ────────────────────────────────────────────────────────
-- TK-102 (Keisha), TK-104 (Tom), TK-105 (Angela) are assigned/departed.
INSERT INTO equipment (id, unit_id, equipment_type, home_warehouse_id, status) VALUES
  ('c1000000-0000-0000-0000-000000000001', 'TK-101', 'truck', 'WH001', 'available'),
  ('c1000000-0000-0000-0000-000000000002', 'TK-102', 'truck', 'WH001', 'assigned'),
  ('c1000000-0000-0000-0000-000000000003', 'TK-103', 'truck', 'WH002', 'available'),
  ('c1000000-0000-0000-0000-000000000004', 'TK-104', 'truck', 'WH003', 'assigned'),
  ('c1000000-0000-0000-0000-000000000005', 'TK-105', 'truck', 'WH004', 'assigned'),
  ('c1000000-0000-0000-0000-000000000006', 'TK-106', 'truck', 'WH005', 'available');

-- ── Equipment — Trailers ──────────────────────────────────────────────────────
-- TR-2001/2002: assigned (loading dock — backing A004/A005 loading BOLs).
-- TR-2003/2004/2005: out with departed drivers.
INSERT INTO equipment (id, unit_id, equipment_type, home_warehouse_id, status) VALUES
  ('c1000000-0000-0000-0000-000000000011', 'TR-2001', 'trailer', 'WH001', 'assigned'),
  ('c1000000-0000-0000-0000-000000000012', 'TR-2002', 'trailer', 'WH001', 'assigned'),
  ('c1000000-0000-0000-0000-000000000013', 'TR-2003', 'trailer', 'WH002', 'assigned'),
  ('c1000000-0000-0000-0000-000000000014', 'TR-2004', 'trailer', 'WH002', 'assigned'),
  ('c1000000-0000-0000-0000-000000000015', 'TR-2005', 'trailer', 'WH003', 'assigned'),
  ('c1000000-0000-0000-0000-000000000016', 'TR-2006', 'trailer', 'WH004', 'available'),
  ('c1000000-0000-0000-0000-000000000017', 'TR-2007', 'trailer', 'WH005', 'available');

-- ── PlanBOL Records ───────────────────────────────────────────────────────────
INSERT INTO plan_bol_record (id, driver_id, originating_wh_id, status, created_at, submitted_at) VALUES
  -- Draft: received 2h ago, not yet claimed
  ('d1000000-0000-0000-0000-000000000001',
   'a1000000-0000-0000-0000-000000000001', 'WH001', 'draft',
   NOW() - INTERVAL '2 hours', NULL),
  -- Draft: received 4h ago
  ('d1000000-0000-0000-0000-000000000002',
   'a1000000-0000-0000-0000-000000000002', 'WH001', 'draft',
   NOW() - INTERVAL '4 hours', NULL),
  -- Pending: route planner working it since yesterday
  ('d1000000-0000-0000-0000-000000000003',
   'a1000000-0000-0000-0000-000000000003', 'WH002', 'plan-progress',
   NOW() - INTERVAL '20 hours', NULL),
  -- Loading: dock crew loading TR-2001 now
  ('d1000000-0000-0000-0000-000000000004',
   'a1000000-0000-0000-0000-000000000004', 'WH001', 'loading',
   NOW() - INTERVAL '22 hours', NULL),
  -- Loading: dock crew loading TR-2002 now
  ('d1000000-0000-0000-0000-000000000005',
   'a1000000-0000-0000-0000-000000000007', 'WH001', 'loading',
   NOW() - INTERVAL '18 hours', NULL),
  -- Validated: loaded and ready, waiting on driver + truck assignment
  ('d1000000-0000-0000-0000-000000000006',
   'a1000000-0000-0000-0000-000000000009', 'WH005', 'validated',
   NOW() - INTERVAL '26 hours', NULL),
  -- Submitted: Tom Brierley, departed 3h ago, 1 warehouse stop done
  ('d1000000-0000-0000-0000-000000000007',
   'a1000000-0000-0000-0000-000000000005', 'WH003', 'submitted',
   NOW() - INTERVAL '8 hours', NOW() - INTERVAL '3 hours'),
  -- Submitted: Keisha Drummond, departed 2.5h ago, 1 warehouse stop done
  ('d1000000-0000-0000-0000-000000000008',
   'a1000000-0000-0000-0000-000000000006', 'WH003', 'submitted',
   NOW() - INTERVAL '7 hours', NOW() - INTERVAL '2 hours 30 minutes'),
  -- Submitted: Angela Torres, departed 2h ago, 1 warehouse stop done
  ('d1000000-0000-0000-0000-000000000009',
   'a1000000-0000-0000-0000-000000000008', 'WH004', 'submitted',
   NOW() - INTERVAL '6 hours', NOW() - INTERVAL '2 hours');

-- ── PlanBOL Stops ─────────────────────────────────────────────────────────────
-- Loading BOL A004 — WH001 pickup + 2 store deliveries (not processed)
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000001', 'd1000000-0000-0000-0000-000000000004', 1, 'WH001',    'warehouse', '{"CLTH001":24,"SPPE001":18}', false),
  ('e1000000-0000-0000-0000-000000000002', 'd1000000-0000-0000-0000-000000000004', 2, 'STORE-101', 'store',     '{"CLTH001":12,"SPPE001":9}',  false),
  ('e1000000-0000-0000-0000-000000000003', 'd1000000-0000-0000-0000-000000000004', 3, 'STORE-102', 'store',     '{"CLTH001":12,"SPPE001":9}',  false);

-- Loading BOL A005 — WH001 pickup + 2 store deliveries (not processed)
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000004', 'd1000000-0000-0000-0000-000000000005', 1, 'WH001',    'warehouse', '{"CLTH004":20,"PWTL001":12}', false),
  ('e1000000-0000-0000-0000-000000000005', 'd1000000-0000-0000-0000-000000000005', 2, 'STORE-103', 'store',     '{"CLTH004":10,"PWTL001":6}',  false),
  ('e1000000-0000-0000-0000-000000000006', 'd1000000-0000-0000-0000-000000000005', 3, 'STORE-104', 'store',     '{"CLTH004":10,"PWTL001":6}',  false);

-- Validated BOL A006 — WH005 pickup + 2 store deliveries (not processed)
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000007', 'd1000000-0000-0000-0000-000000000006', 1, 'WH005',    'warehouse', '{"CLTH002":16,"SPPE003":20}', false),
  ('e1000000-0000-0000-0000-000000000008', 'd1000000-0000-0000-0000-000000000006', 2, 'STORE-501', 'store',     '{"CLTH002":8,"SPPE003":10}',  false),
  ('e1000000-0000-0000-0000-000000000009', 'd1000000-0000-0000-0000-000000000006', 3, 'STORE-502', 'store',     '{"CLTH002":8,"SPPE003":10}',  false);

-- Submitted BOL A007 — Tom Brierley — WH003 done, 2 store stops remaining
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000010', 'd1000000-0000-0000-0000-000000000007', 1, 'WH003',    'warehouse', '{"CLTH007":16,"SPPE001":24}', true,  NOW() - INTERVAL '3 hours 30 minutes'),
  ('e1000000-0000-0000-0000-000000000011', 'd1000000-0000-0000-0000-000000000007', 2, 'STORE-301', 'store',     '{"CLTH007":8,"SPPE001":12}',  false, NULL),
  ('e1000000-0000-0000-0000-000000000012', 'd1000000-0000-0000-0000-000000000007', 3, 'STORE-302', 'store',     '{"CLTH007":8,"SPPE001":12}',  false, NULL);

-- Submitted BOL A008 — Keisha Drummond — WH003 done, 2 store stops remaining
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000013', 'd1000000-0000-0000-0000-000000000008', 1, 'WH003',    'warehouse', '{"SPPE004":30,"PWTL003":10}', true,  NOW() - INTERVAL '3 hours'),
  ('e1000000-0000-0000-0000-000000000014', 'd1000000-0000-0000-0000-000000000008', 2, 'STORE-303', 'store',     '{"SPPE004":15,"PWTL003":5}',  false, NULL),
  ('e1000000-0000-0000-0000-000000000015', 'd1000000-0000-0000-0000-000000000008', 3, 'STORE-304', 'store',     '{"SPPE004":15,"PWTL003":5}',  false, NULL);

-- Submitted BOL A009 — Angela Torres — WH004 done, 1 store stop remaining
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000016', 'd1000000-0000-0000-0000-000000000009', 1, 'WH004',    'warehouse', '{"CLTH002":14,"PWTL006":12}', true,  NOW() - INTERVAL '2 hours 30 minutes'),
  ('e1000000-0000-0000-0000-000000000017', 'd1000000-0000-0000-0000-000000000009', 2, 'STORE-401', 'store',     '{"CLTH002":7,"PWTL006":6}',   false, NULL),
  ('e1000000-0000-0000-0000-000000000018', 'd1000000-0000-0000-0000-000000000009', 3, 'STORE-402', 'store',     '{"CLTH002":7,"PWTL006":6}',   false, NULL);

-- ── Assignments ───────────────────────────────────────────────────────────────
INSERT INTO driver_bol_assignment (id, driver_id, plan_bol_id, equipment_id, assigned_at, departed_at) VALUES
  -- Tom Brierley → BOL A007 → TK-104, departed 3h ago
  ('f1000000-0000-0000-0000-000000000001',
   'a1000000-0000-0000-0000-000000000005',
   'd1000000-0000-0000-0000-000000000007',
   'c1000000-0000-0000-0000-000000000004',
   NOW() - INTERVAL '4 hours',
   NOW() - INTERVAL '3 hours'),
  -- Keisha Drummond → BOL A008 → TK-102, departed 2.5h ago
  ('f1000000-0000-0000-0000-000000000002',
   'a1000000-0000-0000-0000-000000000006',
   'd1000000-0000-0000-0000-000000000008',
   'c1000000-0000-0000-0000-000000000002',
   NOW() - INTERVAL '3 hours 30 minutes',
   NOW() - INTERVAL '2 hours 30 minutes'),
  -- Angela Torres → BOL A009 → TK-105, departed 2h ago
  ('f1000000-0000-0000-0000-000000000003',
   'a1000000-0000-0000-0000-000000000008',
   'd1000000-0000-0000-0000-000000000009',
   'c1000000-0000-0000-0000-000000000005',
   NOW() - INTERVAL '3 hours',
   NOW() - INTERVAL '2 hours');
