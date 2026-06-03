ALTER TABLE driver_bol_assignment
    ADD COLUMN base_rate_per_mile NUMERIC(6,4) NOT NULL DEFAULT 3.20;

ALTER TABLE plan_bol_record
    ADD COLUMN miles_driven NUMERIC(8,2);
