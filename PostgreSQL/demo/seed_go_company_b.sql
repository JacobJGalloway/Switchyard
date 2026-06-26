-- Go DB seed — Company B (Mid-Week Complexity / Digital Parts Brand)
-- database: switchyard-go
-- All timestamps are NOW()-relative. Board state:
--   Draft column:       1 draft, 1 pending
--   LoadingReady:       1 validated (36h old → IsLongWait warning), 2 loading
--   InDelivery:         2 in-transit (1 yellow HOS), 1 roadside breakdown w/ load
--   Delivered:          1 dead-head (completed, window counting down)
--   Available.Resting:  2 weekly rest, 1 daily rest
--   Available.Now:      1 available driver
--   Maintenance:        2 trucks scheduled, 2 trailers scheduled, 1 truck depot breakdown

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

-- ── HOS Windows — mid-week complexity ────────────────────────────────────────
-- Marcus (1): weekly reset in progress — mandated_stop_at set, 62h weekly used
-- Diane (2):  weekly reset in progress — mandated_stop_at set, 61h weekly used
-- Ray (3):    daily 10h rest — mandated_stop_at set
-- Sandra (4): available, moderate hours
-- Tom (5):    in-transit, yellow HOS (9.5h daily — approaching daily limit)
-- Keisha (6): loading BOL assigned, not departed — low hours, green
-- Pete (7):   loading BOL assigned, not departed — moderate hours, green
-- Angela (8): in-transit, green HOS
-- James (9):  validated BOL assigned — moderate hours
-- Carol (10): completed delivery, dead-head back (fulfilled_at set on assignment)
INSERT INTO hos_window (id, driver_id, window_start, daily_hours_used, weekly_hours_used, mandated_stop_at) VALUES
  ('b1000000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', DATE_TRUNC('week', NOW()) - INTERVAL '7 days', 11.0, 62.0, NOW() - INTERVAL '14 hours'),
  ('b1000000-0000-0000-0000-000000000002', 'a1000000-0000-0000-0000-000000000002', DATE_TRUNC('week', NOW()) - INTERVAL '7 days', 11.0, 61.5, NOW() - INTERVAL '10 hours'),
  ('b1000000-0000-0000-0000-000000000003', 'a1000000-0000-0000-0000-000000000003', DATE_TRUNC('week', NOW()),                    10.0, 47.0, NOW() - INTERVAL  '3 hours'),
  ('b1000000-0000-0000-0000-000000000004', 'a1000000-0000-0000-0000-000000000004', DATE_TRUNC('week', NOW()),                     3.0, 21.0, NULL),
  ('b1000000-0000-0000-0000-000000000005', 'a1000000-0000-0000-0000-000000000005', DATE_TRUNC('week', NOW()),                     9.5, 57.0, NULL),
  ('b1000000-0000-0000-0000-000000000006', 'a1000000-0000-0000-0000-000000000006', DATE_TRUNC('week', NOW()),                     4.0, 32.0, NULL),
  ('b1000000-0000-0000-0000-000000000007', 'a1000000-0000-0000-0000-000000000007', DATE_TRUNC('week', NOW()),                     5.5, 38.0, NULL),
  ('b1000000-0000-0000-0000-000000000008', 'a1000000-0000-0000-0000-000000000008', DATE_TRUNC('week', NOW()),                     4.0, 31.0, NULL),
  ('b1000000-0000-0000-0000-000000000009', 'a1000000-0000-0000-0000-000000000009', DATE_TRUNC('week', NOW()),                     6.0, 44.0, NULL),
  ('b1000000-0000-0000-0000-000000000010', 'a1000000-0000-0000-0000-000000000010', DATE_TRUNC('week', NOW()),                     7.5, 50.0, NULL);

-- ── Equipment — Trucks ────────────────────────────────────────────────────────
-- TK-102: scheduled maintenance    → Maintenance column
-- TK-103: scheduled maintenance    → Maintenance column
-- TK-104: assigned/departed (Tom)  → InDelivery
-- TK-105: depot breakdown, no load → Maintenance column
-- TK-106: assigned/departed (Angela) → InDelivery
INSERT INTO equipment (id, unit_id, equipment_type, home_warehouse_id, status) VALUES
  ('c1000000-0000-0000-0000-000000000001', 'TK-101', 'truck', 'WH001', 'available'),
  ('c1000000-0000-0000-0000-000000000002', 'TK-102', 'truck', 'WH001', 'maintenance'),
  ('c1000000-0000-0000-0000-000000000003', 'TK-103', 'truck', 'WH002', 'maintenance'),
  ('c1000000-0000-0000-0000-000000000004', 'TK-104', 'truck', 'WH003', 'assigned'),
  ('c1000000-0000-0000-0000-000000000005', 'TK-105', 'truck', 'WH002', 'breakdown'),
  ('c1000000-0000-0000-0000-000000000006', 'TK-106', 'truck', 'WH004', 'assigned');

-- ── Equipment — Trailers ──────────────────────────────────────────────────────
-- TR-2003: scheduled maintenance → Maintenance column
-- TR-2004: scheduled maintenance → Maintenance column
-- TR-2005: assigned/departed (Tom, BOL-B006)
-- TR-2006: assigned/departed (Angela, BOL-B007)
-- TR-2007: assigned, dead-head return (Carol, BOL-B008)
INSERT INTO equipment (id, unit_id, equipment_type, home_warehouse_id, status) VALUES
  ('c1000000-0000-0000-0000-000000000011', 'TR-2001', 'trailer', 'WH001', 'available'),
  ('c1000000-0000-0000-0000-000000000012', 'TR-2002', 'trailer', 'WH001', 'available'),
  ('c1000000-0000-0000-0000-000000000013', 'TR-2003', 'trailer', 'WH002', 'maintenance'),
  ('c1000000-0000-0000-0000-000000000014', 'TR-2004', 'trailer', 'WH002', 'maintenance'),
  ('c1000000-0000-0000-0000-000000000015', 'TR-2005', 'trailer', 'WH003', 'assigned'),
  ('c1000000-0000-0000-0000-000000000016', 'TR-2006', 'trailer', 'WH004', 'assigned'),
  ('c1000000-0000-0000-0000-000000000017', 'TR-2007', 'trailer', 'WH005', 'assigned');

-- ── Maintenance Records ───────────────────────────────────────────────────────
INSERT INTO maintenance_record (id, equipment_id, description, scheduled_at, estimated_return) VALUES
  ('f1000000-0000-0000-0000-000000000001',
   'c1000000-0000-0000-0000-000000000002',
   'Scheduled 60k mile service — oil, filters, brake inspection',
   NOW() - INTERVAL '6 hours',
   NOW() + INTERVAL '2 days'),
  ('f1000000-0000-0000-0000-000000000002',
   'c1000000-0000-0000-0000-000000000003',
   'Scheduled 60k mile service — oil, filters, tire rotation',
   NOW() - INTERVAL '4 hours',
   NOW() + INTERVAL '1 day'),
  ('f1000000-0000-0000-0000-000000000003',
   'c1000000-0000-0000-0000-000000000013',
   'Annual DOT inspection — lighting, brakes, coupling',
   NOW() - INTERVAL '8 hours',
   NOW() + INTERVAL '1 day'),
  ('f1000000-0000-0000-0000-000000000004',
   'c1000000-0000-0000-0000-000000000014',
   'Refrigeration unit service and seal replacement',
   NOW() - INTERVAL '10 hours',
   NOW() + INTERVAL '3 days');

-- ── Breakdown Record ──────────────────────────────────────────────────────────
-- TK-105: depot breakdown, no load attached → goes to Maintenance column (not InDelivery)
INSERT INTO breakdown_record (id, equipment_id, breakdown_type, load_attached, reported_at) VALUES
  ('f2000000-0000-0000-0000-000000000001',
   'c1000000-0000-0000-0000-000000000005',
   'depot', false,
   NOW() - INTERVAL '5 hours');

-- ── PlanBOL Records ───────────────────────────────────────────────────────────
INSERT INTO plan_bol_record (id, driver_id, originating_wh_id, status, created_at, submitted_at, fulfilled_at) VALUES
  -- Draft
  ('d1000000-0000-0000-0000-000000000001',
   'a1000000-0000-0000-0000-000000000001', 'WH001', 'draft',
   NOW() - INTERVAL '6 hours', NULL, NULL),
  -- Pending
  ('d1000000-0000-0000-0000-000000000002',
   'a1000000-0000-0000-0000-000000000002', 'WH001', 'plan-progress',
   NOW() - INTERVAL '30 hours', NULL, NULL),
  -- Validated — 36h old → IsLongWait=true (WARNING card on board)
  ('d1000000-0000-0000-0000-000000000003',
   'a1000000-0000-0000-0000-000000000009', 'WH005', 'validated',
   NOW() - INTERVAL '36 hours', NULL, NULL),
  -- Loading — Keisha assigned to load
  ('d1000000-0000-0000-0000-000000000004',
   'a1000000-0000-0000-0000-000000000006', 'WH003', 'loading',
   NOW() - INTERVAL '20 hours', NULL, NULL),
  -- Loading — Pete assigned to load
  ('d1000000-0000-0000-0000-000000000005',
   'a1000000-0000-0000-0000-000000000007', 'WH004', 'loading',
   NOW() - INTERVAL '16 hours', NULL, NULL),
  -- Submitted — Tom Brierley, yellow HOS, departed 12h ago
  ('d1000000-0000-0000-0000-000000000006',
   'a1000000-0000-0000-0000-000000000005', 'WH003', 'submitted',
   NOW() - INTERVAL '16 hours', NOW() - INTERVAL '12 hours', NULL),
  -- Submitted — Angela Torres, green HOS, departed 6h ago
  ('d1000000-0000-0000-0000-000000000007',
   'a1000000-0000-0000-0000-000000000008', 'WH004', 'submitted',
   NOW() - INTERVAL '10 hours', NOW() - INTERVAL '6 hours', NULL),
  -- Fulfilled — Carol Metzger, dead-head in progress (all stops done 2h ago)
  ('d1000000-0000-0000-0000-000000000008',
   'a1000000-0000-0000-0000-000000000010', 'WH005', 'fulfilled',
   NOW() - INTERVAL '30 hours', NOW() - INTERVAL '26 hours', NOW() - INTERVAL '2 hours');

-- ── PlanBOL Stops ─────────────────────────────────────────────────────────────

-- Validated BOL B003 (IsLongWait warning) — WH005 pickup + 3 store stops, none processed
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000001', 'd1000000-0000-0000-0000-000000000003', 1, 'WH005',    'warehouse', '{"CLTH001":20,"SPPE001":16}', false),
  ('e1000000-0000-0000-0000-000000000002', 'd1000000-0000-0000-0000-000000000003', 2, 'STORE-501', 'store',     '{"CLTH001":10,"SPPE001":8}',  false),
  ('e1000000-0000-0000-0000-000000000003', 'd1000000-0000-0000-0000-000000000003', 3, 'STORE-502', 'store',     '{"CLTH001":10,"SPPE001":8}',  false);

-- Loading BOL B004 (Keisha) — WH003 pickup + 2 store stops
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000004', 'd1000000-0000-0000-0000-000000000004', 1, 'WH003',    'warehouse', '{"CLTH004":18,"PWTL003":8}',  false),
  ('e1000000-0000-0000-0000-000000000005', 'd1000000-0000-0000-0000-000000000004', 2, 'STORE-303', 'store',     '{"CLTH004":9,"PWTL003":4}',   false),
  ('e1000000-0000-0000-0000-000000000006', 'd1000000-0000-0000-0000-000000000004', 3, 'STORE-304', 'store',     '{"CLTH004":9,"PWTL003":4}',   false);

-- Loading BOL B005 (Pete) — WH004 pickup + 2 store stops
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed) VALUES
  ('e1000000-0000-0000-0000-000000000007', 'd1000000-0000-0000-0000-000000000005', 1, 'WH004',    'warehouse', '{"SPPE003":24,"PWTL001":10}', false),
  ('e1000000-0000-0000-0000-000000000008', 'd1000000-0000-0000-0000-000000000005', 2, 'STORE-401', 'store',     '{"SPPE003":12,"PWTL001":5}',  false),
  ('e1000000-0000-0000-0000-000000000009', 'd1000000-0000-0000-0000-000000000005', 3, 'STORE-402', 'store',     '{"SPPE003":12,"PWTL001":5}',  false);

-- Submitted BOL B006 — Tom Brierley (yellow HOS) — WH003 done, stop 2 of 4 done, 2 remaining
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000010', 'd1000000-0000-0000-0000-000000000006', 1, 'WH003',    'warehouse', '{"CLTH007":20,"SPPE004":30}', true,  NOW() - INTERVAL '12 hours 30 minutes'),
  ('e1000000-0000-0000-0000-000000000011', 'd1000000-0000-0000-0000-000000000006', 2, 'STORE-301', 'store',     '{"CLTH007":10,"SPPE004":15}', true,  NOW() - INTERVAL  '8 hours'),
  ('e1000000-0000-0000-0000-000000000012', 'd1000000-0000-0000-0000-000000000006', 3, 'STORE-302', 'store',     '{"CLTH007":5,"SPPE004":8}',   false, NULL),
  ('e1000000-0000-0000-0000-000000000013', 'd1000000-0000-0000-0000-000000000006', 4, 'STORE-303', 'store',     '{"CLTH007":5,"SPPE004":7}',   false, NULL);

-- Submitted BOL B007 — Angela Torres — WH004 done, 1 of 3 store stops done, 2 remaining
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000014', 'd1000000-0000-0000-0000-000000000007', 1, 'WH004',    'warehouse', '{"CLTH002":16,"PWTL006":12}', true,  NOW() - INTERVAL '6 hours 30 minutes'),
  ('e1000000-0000-0000-0000-000000000015', 'd1000000-0000-0000-0000-000000000007', 2, 'STORE-401', 'store',     '{"CLTH002":8,"PWTL006":6}',   true,  NOW() - INTERVAL '3 hours'),
  ('e1000000-0000-0000-0000-000000000016', 'd1000000-0000-0000-0000-000000000007', 3, 'STORE-402', 'store',     '{"CLTH002":4,"PWTL006":3}',   false, NULL),
  ('e1000000-0000-0000-0000-000000000017', 'd1000000-0000-0000-0000-000000000007', 4, 'STORE-403', 'store',     '{"CLTH002":4,"PWTL006":3}',   false, NULL);

-- Fulfilled BOL B008 — Carol Metzger — all stops processed (dead-head)
INSERT INTO plan_bol_stop (id, plan_bol_id, sequence, location_id, stop_type, delivery_items, is_processed, processed_at) VALUES
  ('e1000000-0000-0000-0000-000000000018', 'd1000000-0000-0000-0000-000000000008', 1, 'WH005',    'warehouse', '{"SPPE001":20,"CLTH004":14}', true, NOW() - INTERVAL '25 hours'),
  ('e1000000-0000-0000-0000-000000000019', 'd1000000-0000-0000-0000-000000000008', 2, 'STORE-501', 'store',     '{"SPPE001":10,"CLTH004":7}',  true, NOW() - INTERVAL '10 hours'),
  ('e1000000-0000-0000-0000-000000000020', 'd1000000-0000-0000-0000-000000000008', 3, 'STORE-502', 'store',     '{"SPPE001":10,"CLTH004":7}',  true, NOW() - INTERVAL  '3 hours');

-- ── Assignments ───────────────────────────────────────────────────────────────
-- Loading assignments (DepartedAt NULL → shows driver overlay on LoadingReady card)
INSERT INTO driver_bol_assignment (id, driver_id, plan_bol_id, equipment_id, assigned_at) VALUES
  ('f1000000-0000-0000-0000-000000000001',
   'a1000000-0000-0000-0000-000000000009',
   'd1000000-0000-0000-0000-000000000003',
   'c1000000-0000-0000-0000-000000000001',
   NOW() - INTERVAL '2 hours'),
  ('f1000000-0000-0000-0000-000000000002',
   'a1000000-0000-0000-0000-000000000006',
   'd1000000-0000-0000-0000-000000000004',
   'c1000000-0000-0000-0000-000000000001',
   NOW() - INTERVAL '4 hours'),
  ('f1000000-0000-0000-0000-000000000003',
   'a1000000-0000-0000-0000-000000000007',
   'd1000000-0000-0000-0000-000000000005',
   'c1000000-0000-0000-0000-000000000001',
   NOW() - INTERVAL '3 hours');

-- In-transit and dead-head assignments (DepartedAt set)
INSERT INTO driver_bol_assignment (id, driver_id, plan_bol_id, equipment_id, assigned_at, departed_at, fulfilled_at) VALUES
  -- Tom Brierley → BOL B006 → TK-104, departed 12h ago (yellow HOS warning)
  ('f1000000-0000-0000-0000-000000000004',
   'a1000000-0000-0000-0000-000000000005',
   'd1000000-0000-0000-0000-000000000006',
   'c1000000-0000-0000-0000-000000000004',
   NOW() - INTERVAL '13 hours',
   NOW() - INTERVAL '12 hours',
   NULL),
  -- Angela Torres → BOL B007 → TK-106, departed 6h ago (green HOS)
  ('f1000000-0000-0000-0000-000000000005',
   'a1000000-0000-0000-0000-000000000008',
   'd1000000-0000-0000-0000-000000000007',
   'c1000000-0000-0000-0000-000000000006',
   NOW() - INTERVAL '7 hours',
   NOW() - INTERVAL '6 hours',
   NULL),
  -- Carol Metzger → BOL B008 → TK-101, dead-head (fulfilled 2h ago, deadhead_confirmed_at NULL)
  ('f1000000-0000-0000-0000-000000000006',
   'a1000000-0000-0000-0000-000000000010',
   'd1000000-0000-0000-0000-000000000008',
   'c1000000-0000-0000-0000-000000000001',
   NOW() - INTERVAL '27 hours',
   NOW() - INTERVAL '26 hours',
   NOW() - INTERVAL  '2 hours');
