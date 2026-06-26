-- Clears all Go transactional data. Master data (warehouses, hos_limit) is preserved.
DELETE FROM driver_bol_assignment;
DELETE FROM plan_bol_stop;
DELETE FROM plan_bol_record;
DELETE FROM breakdown_record;
DELETE FROM maintenance_record;
DELETE FROM equipment;
DELETE FROM hos_window;
DELETE FROM driver;
