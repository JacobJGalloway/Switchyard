CREATE TABLE breakdown_record (
    id             UUID        PRIMARY KEY,
    equipment_id   UUID        NOT NULL REFERENCES equipment(id),
    breakdown_type TEXT        NOT NULL,
    location_desc  TEXT,
    driver_id      UUID        REFERENCES driver(id),
    load_attached  BOOLEAN     NOT NULL DEFAULT false,
    reported_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at    TIMESTAMPTZ
);
