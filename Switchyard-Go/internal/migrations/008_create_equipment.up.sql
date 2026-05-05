CREATE TABLE equipment (
    id                UUID        PRIMARY KEY,
    unit_id           TEXT        NOT NULL UNIQUE,
    equipment_type    TEXT        NOT NULL,
    home_warehouse_id TEXT        NOT NULL,
    status            TEXT        NOT NULL DEFAULT 'available',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
