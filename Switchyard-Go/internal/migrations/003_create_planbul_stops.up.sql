CREATE TABLE plan_bol_stop (
    id             UUID        PRIMARY KEY,
    plan_bol_id    UUID        NOT NULL REFERENCES plan_bol_record(id),
    sequence       INTEGER     NOT NULL,
    location_id    TEXT        NOT NULL,
    stop_type      TEXT        NOT NULL,
    is_processed   BOOLEAN     NOT NULL DEFAULT false,
    processed_at   TIMESTAMPTZ,
    driver_log_ref TEXT
);
