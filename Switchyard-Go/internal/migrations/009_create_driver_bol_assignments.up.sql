CREATE TABLE driver_bol_assignment (
    id                    UUID        PRIMARY KEY,
    driver_id             UUID        NOT NULL REFERENCES driver(id),
    plan_bol_id           UUID        NOT NULL REFERENCES plan_bol_record(id),
    equipment_id          UUID        NOT NULL REFERENCES equipment(id),
    assigned_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    departed_at           TIMESTAMPTZ,
    fulfilled_at          TIMESTAMPTZ,
    deadhead_confirmed_at TIMESTAMPTZ
);
