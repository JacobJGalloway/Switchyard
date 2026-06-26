ALTER TABLE driver_bol_assignment
    ADD COLUMN segment_start_stop_id UUID REFERENCES plan_bol_stop(id),
    ADD COLUMN segment_end_stop_id   UUID REFERENCES plan_bol_stop(id),
    ADD COLUMN transfer_reason       VARCHAR(32),
    ADD COLUMN notes                 TEXT,
    ADD COLUMN transferred_at        TIMESTAMPTZ;
