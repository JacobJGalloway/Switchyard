CREATE TABLE maintenance_record (
    id               UUID        PRIMARY KEY,
    equipment_id     UUID        NOT NULL REFERENCES equipment(id),
    description      TEXT        NOT NULL,
    scheduled_at     TIMESTAMPTZ NOT NULL,
    estimated_return TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ
);
