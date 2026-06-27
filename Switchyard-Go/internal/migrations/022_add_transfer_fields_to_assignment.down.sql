ALTER TABLE driver_bol_assignment
    DROP COLUMN IF EXISTS segment_start_stop_id,
    DROP COLUMN IF EXISTS segment_end_stop_id,
    DROP COLUMN IF EXISTS transfer_reason,
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS transferred_at;
